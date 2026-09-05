package queuesim

import (
	"testing"

	"braess/graph"
)

func buildSingleEdgeNetwork(t *testing.T, freeFlow float64) *graph.Graph {
	t.Helper()
	g := graph.New()
	if _, err := g.AddNode("h1", graph.House); err != nil {
		t.Fatalf("AddNode returned unexpected error: %v", err)
	}
	if _, err := g.AddNode("c1", graph.Company); err != nil {
		t.Fatalf("AddNode returned unexpected error: %v", err)
	}
	if _, err := g.AddEdge("road", "h1", "c1", 0, 0, graph.Constant(freeFlow)); err != nil {
		t.Fatalf("AddEdge returned unexpected error: %v", err)
	}
	return g
}

func findQueueSample(result RunResult, signalID string, at float64) (QueueSample, bool) {
	for _, s := range result.QueueSamples {
		if s.SignalID == signalID && s.Time == at {
			return s, true
		}
	}
	return QueueSample{}, false
}

// Acceptance Scenario 1-2 (spec.md User Story 1) / SC-002: no agent
// crosses a signal-controlled road during its red phase.
func TestRun_NoAgentCrossesDuringRed(t *testing.T) {
	g := buildSingleEdgeNetwork(t, 2)
	// Red for [0,5), green for [5,10), repeating.
	signal := Signal{ID: "s1", EdgeID: "road", GreenDuration: 5, RedDuration: 5, Offset: 5, DischargeRate: 2}
	demand := Demand{Origin: "h1", Destination: "c1", Count: 5, ArrivalInterval: 1}

	result, err := Run(g, []Signal{signal}, []Demand{demand}, 20, 1)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	arrivedCount := 0
	for _, a := range result.Agents {
		if !a.Arrived {
			continue
		}
		arrivedCount++
		enterTime := a.SpawnTime + a.WaitTime
		if !signal.IsGreenAt(enterTime) {
			t.Fatalf("agent (spawn=%g wait=%g) entered the road at t=%g, which is red", a.SpawnTime, a.WaitTime, enterTime)
		}
	}
	if arrivedCount == 0 {
		t.Fatal("no agents arrived — test scenario didn't exercise anything")
	}
}

// SC-001: a green phase long enough to clear the preceding red's
// arrivals fully drains the queue.
func TestRun_QueueDrainsDuringLongEnoughGreen(t *testing.T) {
	g := buildSingleEdgeNetwork(t, 2)
	signal := Signal{ID: "s1", EdgeID: "road", GreenDuration: 5, RedDuration: 5, Offset: 5, DischargeRate: 2}
	// 5 arrivals, all during the first red phase [0,5).
	demand := Demand{Origin: "h1", Destination: "c1", Count: 5, ArrivalInterval: 1}

	result, err := Run(g, []Signal{signal}, []Demand{demand}, 12, 1)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	// t=9 is the last tick of the first green phase [5,10) — with 5
	// queued and a discharge rate of 2/s over 5 seconds of green
	// (capacity for 10), the queue must be empty well before it ends.
	sample, ok := findQueueSample(result, "s1", 9)
	if !ok {
		t.Fatal("no queue sample recorded at t=9")
	}
	if sample.Length != 0 {
		t.Fatalf("queue length at t=9 (end of green) = %d, want 0", sample.Length)
	}
}

// Edge Case (spec.md): arrivals outpacing the discharge rate carry the
// remainder into the next cycle rather than being dropped or silently
// reset when the signal returns to red.
func TestRun_QueueCarriesOverWhenGreenIsntEnoughToDrainIt(t *testing.T) {
	g := buildSingleEdgeNetwork(t, 2)
	// Red for [0,8), green for only [8,10) — a short green, and a slow
	// discharge rate, so 8 queued arrivals can't fully clear. Offset ==
	// GreenDuration is what makes the cycle start in red: phase(t) =
	// (t+Offset) mod cycle is 0 (start of green) exactly at t =
	// RedDuration.
	signal := Signal{ID: "s1", EdgeID: "road", GreenDuration: 2, RedDuration: 8, Offset: 2, DischargeRate: 1}
	demand := Demand{Origin: "h1", Destination: "c1", Count: 8, ArrivalInterval: 1}

	result, err := Run(g, []Signal{signal}, []Demand{demand}, 20, 1)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	endOfGreen, ok := findQueueSample(result, "s1", 9)
	if !ok {
		t.Fatal("no queue sample recorded at t=9")
	}
	if endOfGreen.Length == 0 {
		t.Fatal("queue fully drained despite arrivals outpacing the discharge rate — scenario didn't exercise carry-over")
	}

	// No discharge happens again until the next green phase (t=18); the
	// queue must sit unchanged through the whole second red phase,
	// proving it carried over rather than resetting.
	stillRed, ok := findQueueSample(result, "s1", 17)
	if !ok {
		t.Fatal("no queue sample recorded at t=17")
	}
	if stillRed.Length != endOfGreen.Length {
		t.Fatalf("queue length changed from %d (t=9) to %d (t=17) with no discharge happening — did it get reset?",
			endOfGreen.Length, stillRed.Length)
	}
}

// SC-005: an agent still queued when the run ends is reported as not
// arrived, never as having completed its trip.
func TestRun_StillQueuedAtEndIsNotArrived(t *testing.T) {
	g := buildSingleEdgeNetwork(t, 2)
	// Green for only 1 second per 1,000,000-second cycle — effectively
	// never green again within this run's duration.
	signal := Signal{ID: "s1", EdgeID: "road", GreenDuration: 1, RedDuration: 1_000_000, Offset: 1}
	demand := Demand{Origin: "h1", Destination: "c1", Count: 1, ArrivalInterval: 1}

	result, err := Run(g, []Signal{signal}, []Demand{demand}, 10, 1)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	if len(result.Agents) != 1 {
		t.Fatalf("len(Agents) = %d, want 1", len(result.Agents))
	}
	got := result.Agents[0]
	if got.Arrived {
		t.Fatal("Arrived = true, want false — the signal never turned green during the run")
	}
	if got.WaitTime <= 0 {
		t.Fatalf("WaitTime = %g, want > 0 (the agent spent the whole run queued)", got.WaitTime)
	}
	if got.TravelTime != 0 {
		t.Fatalf("TravelTime = %g, want 0 (the agent never got to move)", got.TravelTime)
	}
}
