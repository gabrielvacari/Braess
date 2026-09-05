package simulation

import (
	"errors"
	"testing"

	"braess/agent"
	"braess/graph"
)

// Acceptance Scenario 1 (spec.md User Story 1): a road's travel time
// matches its function evaluated at the exact number of agents on it.
func TestVolumeBookkeeping_TravelTimeReflectsRealCount(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddEdge(t, g, "road", "a", "b", graph.Linear(0, 1)) // time == volume

	volumes := make(map[string]float64)
	route := agent.Route{Edges: []graph.Edge{{ID: "road", From: "a", To: "b"}}}

	for i := 1; i <= 3; i++ {
		addToVolumes(volumes, route)
		got, err := g.TravelTime("road", volumes["road"])
		if err != nil {
			t.Fatalf("TravelTime returned unexpected error: %v", err)
		}
		if got != float64(i) {
			t.Fatalf("after %d agents, TravelTime(\"road\") = %g, want %g", i, got, float64(i))
		}
	}
}

// Acceptance Scenario 2 (spec.md User Story 1): agents moving from one
// road to another must be reflected as one fewer / one more.
func TestVolumeBookkeeping_MovingAgentsUpdatesBothRoads(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddEdge(t, g, "road-1", "a", "b", graph.Linear(0, 1))
	mustAddEdge(t, g, "road-2", "a", "b", graph.Linear(0, 1))

	route1 := agent.Route{Edges: []graph.Edge{{ID: "road-1", From: "a", To: "b"}}}
	route2 := agent.Route{Edges: []graph.Edge{{ID: "road-2", From: "a", To: "b"}}}

	volumes := make(map[string]float64)
	// Five agents start on road-1.
	for i := 0; i < 5; i++ {
		addToVolumes(volumes, route1)
	}
	if volumes["road-1"] != 5 || volumes["road-2"] != 0 {
		t.Fatalf("initial volumes = %v, want road-1=5 road-2=0", volumes)
	}

	// Two of them move to road-2.
	removeFromVolumes(volumes, route1)
	addToVolumes(volumes, route2)
	removeFromVolumes(volumes, route1)
	addToVolumes(volumes, route2)

	if volumes["road-1"] != 3 {
		t.Fatalf("road-1 volume after moving 2 agents away = %g, want 3", volumes["road-1"])
	}
	if volumes["road-2"] != 2 {
		t.Fatalf("road-2 volume after moving 2 agents in = %g, want 2", volumes["road-2"])
	}
}

// SC-001: across 3 distinct networks with two parallel roads between one
// origin and one destination, Run's final assignment is a genuine Nash
// equilibrium — no agent could unilaterally switch to the other road and
// strictly improve, using the same marginal-cost convention Run itself
// uses (research.md decision #1's guarantee, checked directly rather
// than by hand-deriving "the" expected split — a two-choice congestion
// game can legitimately have more than one valid equilibrium, so
// verifying the stability property itself is the robust check).
func TestRun_ReachesTwoEdgeEquilibrium(t *testing.T) {
	tests := []struct {
		name    string
		edge1Fn graph.TravelTimeFunc
		edge2Fn graph.TravelTimeFunc
		size    int
	}{
		{"symmetric slopes", graph.Linear(0, 1), graph.Linear(0, 1), 10},
		{"asymmetric slopes", graph.Linear(0, 1), graph.Linear(0, 2), 9},
		{"linear vs constant", graph.Linear(0, 1), graph.Constant(5), 8},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := graph.New()
			mustAddNode(t, g, "a", graph.House)
			mustAddNode(t, g, "b", graph.Company)
			mustAddEdge(t, g, "edge-1", "a", "b", tc.edge1Fn)
			mustAddEdge(t, g, "edge-2", "a", "b", tc.edge2Fn)

			result, err := Run(g, Population{Origin: "a", Destination: "b", Size: tc.size, MaxRounds: DefaultMaxRounds})
			if err != nil {
				t.Fatalf("Run returned unexpected error: %v", err)
			}
			if !result.Converged {
				t.Fatalf("Run did not converge within %d rounds (used %d)", DefaultMaxRounds, result.Rounds)
			}
			if len(result.Routes) != tc.size {
				t.Fatalf("len(Routes) = %d, want %d", len(result.Routes), tc.size)
			}

			assertTwoEdgeEquilibrium(t, g, result.Routes, "edge-1", "edge-2", tc.size)
		})
	}
}

// SC-004: with MaxRounds: 0, Run performs the initial loading only and
// must report Converged: false (stability was never checked).
func TestRun_MaxRoundsZeroReportsNotConverged(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddEdge(t, g, "edge-1", "a", "b", graph.Linear(0, 1))
	mustAddEdge(t, g, "edge-2", "a", "b", graph.Linear(0, 1))

	result, err := Run(g, Population{Origin: "a", Destination: "b", Size: 10, MaxRounds: 0})
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if result.Converged {
		t.Fatal("Run with MaxRounds: 0 reported Converged: true, want false")
	}
	if result.Rounds != 0 {
		t.Fatalf("Run with MaxRounds: 0 reported Rounds = %d, want 0", result.Rounds)
	}
}

// FR-010: a population of zero agents must not error, and must report a
// trivial (zero-cost) result.
func TestRun_ZeroAgents(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddEdge(t, g, "edge-1", "a", "b", graph.Constant(5))

	result, err := Run(g, Population{Origin: "a", Destination: "b", Size: 0, MaxRounds: DefaultMaxRounds})
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if len(result.Routes) != 0 {
		t.Fatalf("Routes = %v, want empty", result.Routes)
	}
	if result.TotalTravelTime != 0 || result.AverageTravelTime != 0 {
		t.Fatalf("TotalTravelTime=%g AverageTravelTime=%g, want both 0", result.TotalTravelTime, result.AverageTravelTime)
	}
	if !result.Converged {
		t.Fatal("Run with zero agents reported Converged: false, want true")
	}
}

// Edge Case (spec.md): no route exists between origin and destination —
// Run must propagate that, not crash or silently return an empty result.
func TestRun_NoRouteExists(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	// no edges at all

	_, err := Run(g, Population{Origin: "a", Destination: "b", Size: 5, MaxRounds: DefaultMaxRounds})
	if !errors.Is(err, agent.ErrNoRoute) {
		t.Fatalf("Run with no connecting edges = %v, want an error wrapping agent.ErrNoRoute", err)
	}
}

// assertTwoEdgeEquilibrium verifies that, for a graph with exactly two
// parallel single-edge routes between an origin and destination, the
// given routes form a Nash equilibrium: no agent on either edge could
// strictly improve by switching to the other, using the same marginal
// (including-self) cost convention Run itself uses.
func assertTwoEdgeEquilibrium(t *testing.T, g *graph.Graph, routes []agent.Route, edge1, edge2 string, n int) {
	t.Helper()

	v1, v2 := 0, 0
	for _, r := range routes {
		if len(r.Edges) != 1 {
			t.Fatalf("route has %d edges, want exactly 1 (single-edge network)", len(r.Edges))
		}
		switch r.Edges[0].ID {
		case edge1:
			v1++
		case edge2:
			v2++
		default:
			t.Fatalf("route uses unexpected edge %q", r.Edges[0].ID)
		}
	}
	if v1+v2 != n {
		t.Fatalf("v1(%d) + v2(%d) = %d, want %d", v1, v2, v1+v2, n)
	}

	cost1, err := g.TravelTime(edge1, float64(v1))
	if err != nil {
		t.Fatalf("TravelTime(%q, %d) returned unexpected error: %v", edge1, v1, err)
	}
	cost2, err := g.TravelTime(edge2, float64(v2))
	if err != nil {
		t.Fatalf("TravelTime(%q, %d) returned unexpected error: %v", edge2, v2, err)
	}

	if v1 > 0 {
		// Would an edge-1 agent improve by switching to edge-2, joining
		// as its (v2+1)th user?
		altCost, err := g.TravelTime(edge2, float64(v2+1))
		if err != nil {
			t.Fatalf("TravelTime(%q, %d) returned unexpected error: %v", edge2, v2+1, err)
		}
		if altCost < cost1-epsilon {
			t.Fatalf("not an equilibrium: an edge-1 agent (cost %g) could improve to %g on edge-2", cost1, altCost)
		}
	}
	if v2 > 0 {
		altCost, err := g.TravelTime(edge1, float64(v1+1))
		if err != nil {
			t.Fatalf("TravelTime(%q, %d) returned unexpected error: %v", edge1, v1+1, err)
		}
		if altCost < cost2-epsilon {
			t.Fatalf("not an equilibrium: an edge-2 agent (cost %g) could improve to %g on edge-1", cost2, altCost)
		}
	}
}

func mustAddNode(t *testing.T, g *graph.Graph, id string, nt graph.NodeType) {
	t.Helper()
	if _, err := g.AddNode(id, nt); err != nil {
		t.Fatalf("AddNode(%q) returned unexpected error: %v", id, err)
	}
}

func mustAddEdge(t *testing.T, g *graph.Graph, id, from, to string, tt graph.TravelTimeFunc) {
	t.Helper()
	if _, err := g.AddEdge(id, from, to, 1, 1, tt); err != nil {
		t.Fatalf("AddEdge(%q) returned unexpected error: %v", id, err)
	}
}
