package simulation

import (
	"testing"

	"braess/graph"
)

// braessAgentCount matches the textbook Braess's Paradox example (4000
// drivers), so the results below reproduce its well-known numbers
// exactly (research.md decision #6).
const braessAgentCount = 4000

// buildBraessNetwork builds the classic four-node Braess's Paradox
// network: S -> A and B -> T are congestible (time = volume/100); S -> B
// and A -> T are constant at 45. withShortcut adds the free A -> B edge
// that is the source of the paradox.
func buildBraessNetwork(withShortcut bool) (*graph.Graph, error) {
	g := graph.New()
	for _, id := range []string{"S", "A", "B", "T"} {
		if _, err := g.AddNode(id, graph.Intersection); err != nil {
			return nil, err
		}
	}

	edges := []struct {
		id, from, to string
		tt           graph.TravelTimeFunc
	}{
		{"S-A", "S", "A", graph.Linear(0, 0.01)},
		{"S-B", "S", "B", graph.Constant(45)},
		{"A-T", "A", "T", graph.Constant(45)},
		{"B-T", "B", "T", graph.Linear(0, 0.01)},
	}
	for _, e := range edges {
		if _, err := g.AddEdge(e.id, e.from, e.to, 0, 0, e.tt); err != nil {
			return nil, err
		}
	}

	if withShortcut {
		if _, err := g.AddEdge("A-B", "A", "B", 0, 0, graph.Constant(0)); err != nil {
			return nil, err
		}
	}

	return g, nil
}

// SC-002: adding the extra road to the classic Braess network makes the
// average travel time strictly worse for the population, reproducing the
// paradox — the well-known 65 -> 80 result at this population size.
func TestBraessParadox_AddingRoadMakesItWorse(t *testing.T) {
	population := Population{Origin: "S", Destination: "T", Size: braessAgentCount, MaxRounds: DefaultMaxRounds}

	without, err := buildBraessNetwork(false)
	if err != nil {
		t.Fatalf("buildBraessNetwork(false) returned unexpected error: %v", err)
	}
	resultWithout, err := Run(without, population)
	if err != nil {
		t.Fatalf("Run(without shortcut) returned unexpected error: %v", err)
	}
	if !resultWithout.Converged {
		t.Fatalf("Run(without shortcut) did not converge within %d rounds", DefaultMaxRounds)
	}

	with, err := buildBraessNetwork(true)
	if err != nil {
		t.Fatalf("buildBraessNetwork(true) returned unexpected error: %v", err)
	}
	resultWith, err := Run(with, population)
	if err != nil {
		t.Fatalf("Run(with shortcut) returned unexpected error: %v", err)
	}
	if !resultWith.Converged {
		t.Fatalf("Run(with shortcut) did not converge within %d rounds", DefaultMaxRounds)
	}

	if resultWith.AverageTravelTime <= resultWithout.AverageTravelTime {
		t.Fatalf("adding the road did not make things worse: without=%g, with=%g",
			resultWithout.AverageTravelTime, resultWith.AverageTravelTime)
	}

	// The textbook result: 65 without the shortcut, 80 with it.
	const tolerance = 0.5
	if abs(resultWithout.AverageTravelTime-65) > tolerance {
		t.Errorf("average travel time without shortcut = %g, want ~65", resultWithout.AverageTravelTime)
	}
	if abs(resultWith.AverageTravelTime-80) > tolerance {
		t.Errorf("average travel time with shortcut = %g, want ~80", resultWith.AverageTravelTime)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
