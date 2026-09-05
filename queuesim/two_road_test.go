package queuesim

import (
	"testing"

	"braess/graph"
)

// buildTwoRoadScenario builds the classic two-road, opposite-phase
// setup: one house, one company, two parallel roads between them (each
// bidirectional-in-spirit is out of scope here — feature 007's web-only
// bidirectionality doesn't apply to this engine-only feature), each
// governed by its own signal, the two signals always in opposite phase
// purely by sharing a cycle and being offset by one GreenDuration
// (research.md decision #4) — never both green, never both red.
func buildTwoRoadScenario(t *testing.T) (*graph.Graph, []Signal, []Demand) {
	t.Helper()

	g := graph.New()
	if _, err := g.AddNode("h1", graph.House); err != nil {
		t.Fatalf("AddNode returned unexpected error: %v", err)
	}
	if _, err := g.AddNode("c1", graph.Company); err != nil {
		t.Fatalf("AddNode returned unexpected error: %v", err)
	}
	if _, err := g.AddEdge("road-a", "h1", "c1", 0, 0, graph.Constant(1)); err != nil {
		t.Fatalf("AddEdge returned unexpected error: %v", err)
	}
	if _, err := g.AddEdge("road-b", "h1", "c1", 0, 0, graph.Constant(1)); err != nil {
		t.Fatalf("AddEdge returned unexpected error: %v", err)
	}

	signals := []Signal{
		{ID: "signal-a", EdgeID: "road-a", GreenDuration: 5, RedDuration: 5, Offset: 5, DischargeRate: 2},
		{ID: "signal-b", EdgeID: "road-b", GreenDuration: 5, RedDuration: 5, Offset: 0, DischargeRate: 2},
	}
	// Each road's long-run capacity is DischargeRate * (fraction of time
	// green) = 2 * 0.5 = 1/s; even always picking whichever road is
	// currently green tops out at 2/s (only one is ever green at once).
	// An arrival rate above that (here, ~3.33/s) is what actually forces
	// a sustained, sampleable backlog rather than one that clears within
	// the same tick it formed.
	demands := []Demand{
		{Origin: "h1", Destination: "c1", Count: 150, ArrivalInterval: 0.3},
	}

	return g, signals, demands
}

// twoRoadTick is the simulated-step size used for the two-road scenario
// — finer than the single-edge tests' 1s tick, so ArrivalInterval's 0.3s
// spawn schedule is resolved smoothly rather than bunching arrivals onto
// whichever tick boundary they happen to fall near.
const twoRoadTick = 0.25

// SC-003: the two-road run produces a distinguishable, measurable
// queue-length-over-time report for each road — enough to determine,
// from the data alone, whether one road ended up more congested, it
// oscillated, or it balanced. This test deliberately does not assert
// which of those outcomes occurs (spec Assumptions) — only that the data
// needed to answer the question exists and is well-formed.
func TestTwoRoadScenario_ProducesPerRoadQueueReport(t *testing.T) {
	g, signals, demands := buildTwoRoadScenario(t)

	result, err := Run(g, signals, demands, 70, twoRoadTick)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	avgQueue := make(map[string]float64)
	count := make(map[string]int)
	for _, s := range result.QueueSamples {
		avgQueue[s.SignalID] += float64(s.Length)
		count[s.SignalID]++
	}
	for _, sig := range signals {
		if count[sig.ID] == 0 {
			t.Fatalf("no queue samples recorded for %q", sig.ID)
		}
		avgQueue[sig.ID] /= float64(count[sig.ID])
	}

	// The report must at least distinguish "some congestion happened
	// somewhere" — a run where every sample is 0 for both roads would
	// mean the scenario never actually exercised a queue, telling us
	// nothing about the question this feature exists to make observable.
	if avgQueue["signal-a"] == 0 && avgQueue["signal-b"] == 0 {
		t.Fatal("both roads show zero average queue length — scenario didn't exercise any congestion")
	}

	t.Logf("average queue length: road-a=%.2f road-b=%.2f", avgQueue["signal-a"], avgQueue["signal-b"])
}

// FR-006: each agent's wait time is reported separately from travel
// time, and at least some agents actually experienced a nonzero wait
// (proving the queue mechanic engaged, not just passed through).
func TestTwoRoadScenario_ReportsPerAgentWaitSeparateFromTravel(t *testing.T) {
	g, signals, demands := buildTwoRoadScenario(t)

	result, err := Run(g, signals, demands, 70, twoRoadTick)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	if len(result.Agents) != demands[0].Count {
		t.Fatalf("len(Agents) = %d, want %d", len(result.Agents), demands[0].Count)
	}

	sawWait := false
	for _, a := range result.Agents {
		if a.WaitTime > 0 {
			sawWait = true
		}
		if a.WaitTime < 0 || a.TravelTime < 0 {
			t.Fatalf("agent report has a negative time: %+v", a)
		}
	}
	if !sawWait {
		t.Fatal("no agent recorded any wait time — the queue mechanic never engaged in this scenario")
	}
}

// SC-004 (User Story 3): changing a signal's configured durations
// changes the reported queue pattern.
func TestTwoRoadScenario_DifferentTimingsProduceDifferentQueuePatterns(t *testing.T) {
	g1, signals1, demands1 := buildTwoRoadScenario(t)
	resultA, err := Run(g1, signals1, demands1, 70, twoRoadTick)
	if err != nil {
		t.Fatalf("Run (config A) returned unexpected error: %v", err)
	}

	g2, signals2, demands2 := buildTwoRoadScenario(t)
	// A much longer green, much shorter red — a materially different
	// timing plan.
	signals2[0].GreenDuration, signals2[0].RedDuration = 9, 1
	signals2[0].Offset = 1
	signals2[1].GreenDuration, signals2[1].RedDuration = 9, 1
	signals2[1].Offset = 0
	resultB, err := Run(g2, signals2, demands2, 70, twoRoadTick)
	if err != nil {
		t.Fatalf("Run (config B) returned unexpected error: %v", err)
	}

	avg := func(result RunResult, signalID string) float64 {
		total, n := 0.0, 0
		for _, s := range result.QueueSamples {
			if s.SignalID == signalID {
				total += float64(s.Length)
				n++
			}
		}
		if n == 0 {
			return 0
		}
		return total / float64(n)
	}

	aBefore := avg(resultA, "signal-a") + avg(resultA, "signal-b")
	aAfter := avg(resultB, "signal-a") + avg(resultB, "signal-b")

	if aBefore == aAfter {
		t.Fatalf("changing signal timing did not change the reported queue pattern (both averaged %g)", aBefore)
	}
}
