package queuesim

import (
	"fmt"

	"braess/graph"
)

// defaultDischargeRate derives a discharge rate (agents/second) from an
// edge's own free-flow travel time when a Signal doesn't set one
// explicitly (spec Assumptions): one agent clears per free-flow-travel-
// time interval, so a slower road discharges its queue more slowly and a
// faster one more quickly — the same relationship the road already has
// to how long it takes one agent to use it.
func defaultDischargeRate(freeFlowTime float64) float64 {
	if freeFlowTime <= 0 {
		return 1
	}
	return 1 / freeFlowTime
}

// agentState is which of the two things an in-flight agent is doing.
type agentState int

const (
	moving agentState = iota
	queued
)

// inFlightAgent is queuesim's internal bookkeeping for one spawned agent
// still in the run — not part of the public contract (see report.go's
// AgentReport for what a caller actually sees).
type inFlightAgent struct {
	id          int
	demandIndex int
	route       []graph.Edge
	edgeIndex   int
	state       agentState
	remaining   float64 // seconds left on the current edge, if moving
	spawnTime   float64
	travelTime  float64
	waitTime    float64
}

// Run simulates signals and demands over g for the given duration, in
// steps of tick seconds (research.md decision #1). Agents spawn per
// their Demand's schedule, each computing its own route once at spawn
// (research.md decision #2) via queuesim's own wait-aware pathfinding,
// queuing at any signal-controlled edge on its route that is red when it
// arrives there, and discharging at that edge's signal's rate once
// green.
//
// Returns an error, rather than swallowing it, if a Demand references a
// node with no route to its destination, or an edge's travel-time
// function itself errors.
func Run(g *graph.Graph, signals []Signal, demands []Demand, duration, tick float64) (RunResult, error) {
	signalByEdge := make(map[string]Signal, len(signals))
	for _, s := range signals {
		signalByEdge[s.EdgeID] = s
	}

	freeFlow := make(map[string]float64)
	for _, e := range g.Edges() {
		tt, err := e.TravelTime(0)
		if err != nil {
			return RunResult{}, fmt.Errorf("queuesim: computing free-flow time for edge %q: %w", e.ID, err)
		}
		freeFlow[e.ID] = tt
	}

	dischargeRate := make(map[string]float64, len(signalByEdge))
	for edgeID, s := range signalByEdge {
		if s.DischargeRate > 0 {
			dischargeRate[edgeID] = s.DischargeRate
		} else {
			dischargeRate[edgeID] = defaultDischargeRate(freeFlow[edgeID])
		}
	}

	queues := make(map[string][]*inFlightAgent)
	dischargeCredit := make(map[string]float64, len(signalByEdge))

	// enterEdge places an agent onto route[edgeIndex]: moving immediately
	// if the edge isn't signal-controlled, or queued otherwise. A
	// signal-controlled edge's entry is *always* gated through its queue
	// and discharge-rate accumulator (step 3 below) — even with no
	// pre-existing backlog and a green light right now — so the road's
	// throughput capacity is enforced uniformly, not just when a queue
	// has already formed. In light traffic this still resolves within
	// the same or the next tick (the accumulator almost always has
	// credit available); it only becomes a visible wait once arrivals
	// exceed what the discharge rate can absorb.
	enterEdge := func(a *inFlightAgent, edgeIndex int, t float64) {
		a.edgeIndex = edgeIndex
		edgeID := a.route[edgeIndex].ID
		if _, controlled := signalByEdge[edgeID]; controlled {
			a.state = queued
			queues[edgeID] = append(queues[edgeID], a)
			return
		}
		a.state = moving
		a.remaining = freeFlow[edgeID]
	}

	// edgeCostAt prices an edge for route selection at simulated time t:
	// free-flow time, plus — if signal-controlled — a snapshot estimate
	// of the wait a new arrival would face right now (research.md
	// decision #2; not a prediction of conditions when the agent would
	// actually reach a *later* edge on its route).
	edgeCostAt := func(t float64) func(e graph.Edge) (float64, error) {
		return func(e graph.Edge) (float64, error) {
			base := freeFlow[e.ID]
			if s, controlled := signalByEdge[e.ID]; controlled {
				return base + s.estimatedWait(t, len(queues[e.ID]), dischargeRate[e.ID]), nil
			}
			return base, nil
		}
	}

	var inFlight []*inFlightAgent
	var result RunResult
	spawned := make([]int, len(demands))
	nextID := 0 // spawn-order agent identifier (feature 010, research.md decision #2)

	for t := 0.0; t < duration; t += tick {
		// 1. Spawn agents whose arrival schedule has come due.
		for di, d := range demands {
			for spawned[di] < d.Count && float64(spawned[di])*d.ArrivalInterval <= t {
				spawned[di]++
				id := nextID
				nextID++
				route, err := shortestRouteAt(g, d.Origin, d.Destination, edgeCostAt(t))
				if err != nil {
					return RunResult{}, fmt.Errorf("demand %d: %w", di, err)
				}
				a := &inFlightAgent{id: id, demandIndex: di, route: route, spawnTime: t}
				if len(route) == 0 {
					result.Agents = append(result.Agents, AgentReport{ID: id, DemandIndex: di, SpawnTime: t, Arrived: true})
					continue
				}
				enterEdge(a, 0, t)
				inFlight = append(inFlight, a)
			}
		}

		// 2. Advance every in-flight agent by one tick.
		stillInFlight := inFlight[:0]
		for _, a := range inFlight {
			if a.state == moving {
				// The whole tick counts as travel time even on the tick
				// an edge finishes partway through — a small, bounded
				// over-count inherent to fixed-tick simulation
				// (research.md decision #1), negligible for tick sizes
				// small relative to edge travel times.
				a.remaining -= tick
				a.travelTime += tick
				if a.remaining <= 0 {
					if a.edgeIndex+1 >= len(a.route) {
						// One last position sample at the moment of
						// arrival (feature 010, research.md decision #3)
						// — without it, the final edge of a trip would
						// never show the agent actually reaching its end.
						result.Positions = append(result.Positions, PositionSample{
							Time: t, AgentID: a.id, EdgeID: a.route[a.edgeIndex].ID, Progress: 1,
						})
						result.Agents = append(result.Agents, AgentReport{
							ID: a.id, DemandIndex: a.demandIndex, SpawnTime: a.spawnTime,
							TravelTime: a.travelTime, WaitTime: a.waitTime, Arrived: true,
						})
						continue
					}
					enterEdge(a, a.edgeIndex+1, t)
				}
			}
			if a.state == queued {
				a.waitTime += tick
			}
			stillInFlight = append(stillInFlight, a)
		}
		inFlight = stillInFlight

		// 3. Discharge queues for signals currently green, using a
		// fractional credit accumulator so a non-integer discharge rate
		// isn't rounded away (data-model.md).
		for edgeID, s := range signalByEdge {
			if !s.IsGreenAt(t) {
				continue
			}
			dischargeCredit[edgeID] += dischargeRate[edgeID] * tick
			for dischargeCredit[edgeID] >= 1 && len(queues[edgeID]) > 0 {
				next := queues[edgeID][0]
				queues[edgeID] = queues[edgeID][1:]
				dischargeCredit[edgeID]--
				next.state = moving
				next.remaining = freeFlow[edgeID]
			}
		}

		// 4. Sample every signal's queue length and phase this tick
		// (FR-004; feature 010 FR-002 adds Green).
		for _, s := range signals {
			result.QueueSamples = append(result.QueueSamples, QueueSample{
				Time: t, SignalID: s.ID, Length: len(queues[s.EdgeID]), Green: s.IsGreenAt(t),
			})
		}

		// 5. Record every remaining in-flight agent's position this tick
		// (feature 010, FR-001), reflecting any promotion the discharge
		// step above just made.
		for _, a := range inFlight {
			edgeID := a.route[a.edgeIndex].ID
			progress := 0.0
			if a.state == moving {
				// A zero-or-negative free-flow time (e.g. a linear
				// travel-time function whose freeFlow parameter is 0)
				// makes this edge instantaneous — 1-a.remaining/0 would
				// be NaN, which encoding/json cannot serialize at all
				// (silently producing an empty response body). Treat it
				// as already fully traversed instead.
				if freeFlow[edgeID] > 0 {
					progress = 1 - a.remaining/freeFlow[edgeID]
				} else {
					progress = 1
				}
			}
			result.Positions = append(result.Positions, PositionSample{
				Time: t, AgentID: a.id, EdgeID: edgeID, Progress: progress, Queued: a.state == queued,
			})
		}
	}

	// Anything still moving or queued when the run ends is reported as
	// not arrived (FR-008, SC-005) — never presented as having completed.
	for _, a := range inFlight {
		result.Agents = append(result.Agents, AgentReport{
			ID: a.id, DemandIndex: a.demandIndex, SpawnTime: a.spawnTime,
			TravelTime: a.travelTime, WaitTime: a.waitTime, Arrived: false,
		})
	}

	return result, nil
}
