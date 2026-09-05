package simulation

import (
	"fmt"

	"braess/agent"
	"braess/graph"
)

// Demand is one origin/destination pair together with how many agents
// travel between them in a given run.
type Demand struct {
	Origin, Destination string // Node IDs
	Size                int    // number of agents; 0 is valid
}

// MultiPopulation is several Demands sharing one road network and its
// real congestion — the generalization of Population to more than one
// fixed origin/destination pair.
type MultiPopulation struct {
	Demands []Demand
	// MaxRounds has the same semantics as Population.MaxRounds: required,
	// with 0 meaning "loading only, no refinement" (always Converged: false).
	MaxRounds int
}

// DemandResult is one Demand's own outcome within a MultiPopulation run.
type DemandResult struct {
	Demand            Demand
	Routes            []agent.Route
	TotalTravelTime   float64
	AverageTravelTime float64
}

// MultiAssignmentResult is the outcome of running a MultiPopulation to a
// stable (or round-limited) shared assignment.
type MultiAssignmentResult struct {
	// PerDemand has the same length and order as the MultiPopulation's
	// Demands.
	PerDemand []DemandResult
	// TotalTravelTime and AverageTravelTime are computed across every
	// agent, from every demand, combined.
	TotalTravelTime   float64
	AverageTravelTime float64
	Converged         bool
	Rounds            int
}

// RunDemands computes a shared route assignment for every agent across
// every demand in mp.Demands: agents are loaded (all of one demand, then
// the next, in mp.Demands order — research.md decision #5: load order
// doesn't affect whether a valid equilibrium is eventually found) against
// real accumulating traffic from every demand combined, then refined
// round by round exactly as Run does — every agent, regardless of which
// demand it belongs to, may switch to a strictly better route given
// everyone else's current choice, until no agent improves or
// mp.MaxRounds is reached.
//
// Returns an error, rather than swallowing it, if any agent's route
// computation fails, naming the specific Demand (by index and
// origin/destination) it belongs to (FR-008) — e.g. no route exists for
// that demand's origin/destination (wraps agent.ErrNoRoute), or an
// edge's travel-time function itself errors.
func RunDemands(g *graph.Graph, mp MultiPopulation) (MultiAssignmentResult, error) {
	ods := buildODPairs(mp.Demands)
	perDemand := make([]DemandResult, len(mp.Demands))
	for i, d := range mp.Demands {
		perDemand[i] = DemandResult{Demand: d}
	}

	if len(ods) == 0 {
		return MultiAssignmentResult{PerDemand: perDemand, Converged: true}, nil
	}

	routes, volumes, failedAt, err := loadIncrementally(g, ods)
	if err != nil {
		return MultiAssignmentResult{}, demandError(mp.Demands, ods, failedAt, err)
	}

	rounds, converged, failedAt, err := refine(g, ods, mp.MaxRounds, routes, volumes)
	if err != nil {
		return MultiAssignmentResult{}, demandError(mp.Demands, ods, failedAt, err)
	}

	// Refresh each route's own TotalTravelTime against the final volumes
	// (same reasoning as Run) and, at the same time, group each agent's
	// final route into its own Demand's result (FR-005) — every agent's
	// ods[i].demandIndex says exactly which perDemand entry it belongs to.
	total := 0.0
	for i, r := range routes {
		cost, cerr := routeCost(g, r, volumes)
		if cerr != nil {
			return MultiAssignmentResult{}, demandError(mp.Demands, ods, i, cerr)
		}
		routes[i].TotalTravelTime = cost
		total += cost

		di := ods[i].demandIndex
		perDemand[di].Routes = append(perDemand[di].Routes, routes[i])
		perDemand[di].TotalTravelTime += cost
	}
	for i := range perDemand {
		if perDemand[i].Demand.Size > 0 {
			perDemand[i].AverageTravelTime = perDemand[i].TotalTravelTime / float64(perDemand[i].Demand.Size)
		}
	}

	return MultiAssignmentResult{
		PerDemand:         perDemand,
		TotalTravelTime:   total,
		AverageTravelTime: total / float64(len(ods)),
		Converged:         converged,
		Rounds:            rounds,
	}, nil
}

// buildODPairs flattens demands into one []odPair per agent, tagged with
// which demand (by index into demands) it belongs to.
func buildODPairs(demands []Demand) []odPair {
	var ods []odPair
	for di, d := range demands {
		for j := 0; j < d.Size; j++ {
			ods = append(ods, odPair{origin: d.Origin, destination: d.Destination, demandIndex: di})
		}
	}
	return ods
}

// demandError wraps err with which Demand (by index and origin/
// destination) the failing agent — ods[agentIndex] — belongs to (FR-008).
// errors.Is(err, agent.ErrNoRoute) still succeeds through the wrap.
func demandError(demands []Demand, ods []odPair, agentIndex int, err error) error {
	if agentIndex < 0 || agentIndex >= len(ods) {
		return err
	}
	d := demands[ods[agentIndex].demandIndex]
	return fmt.Errorf("demand %d (origin=%q destination=%q): %w", ods[agentIndex].demandIndex, d.Origin, d.Destination, err)
}
