package simulation

import (
	"errors"
	"strings"
	"testing"

	"braess/agent"
	"braess/graph"
)

// buildBottleneckNetwork builds a network where two houses (h1, h2) each
// feed into node A, a single shared road A->B is the only way across, and
// B fans out to two companies (c1, c2). Both house-1->company-1 and
// house-2->company-2 demands are forced through the same A->B edge — the
// shared bottleneck this file's tests exercise.
func buildBottleneckNetwork(t *testing.T, bottleneckFn graph.TravelTimeFunc) *graph.Graph {
	t.Helper()
	g := graph.New()
	for _, id := range []string{"h1", "h2", "A", "B", "c1", "c2"} {
		if _, err := g.AddNode(id, graph.Intersection); err != nil {
			t.Fatalf("AddNode(%q) returned unexpected error: %v", id, err)
		}
	}
	edges := []struct {
		id, from, to string
		tt           graph.TravelTimeFunc
	}{
		{"h1-A", "h1", "A", graph.Constant(0)},
		{"h2-A", "h2", "A", graph.Constant(0)},
		{"A-B", "A", "B", bottleneckFn},
		{"B-c1", "B", "c1", graph.Constant(0)},
		{"B-c2", "B", "c2", graph.Constant(0)},
	}
	for _, e := range edges {
		if _, err := g.AddEdge(e.id, e.from, e.to, 0, 0, e.tt); err != nil {
			t.Fatalf("AddEdge(%q) returned unexpected error: %v", e.id, err)
		}
	}
	return g
}

// Acceptance Scenario 1 (spec.md User Story 1): the shared road's travel
// time reflects agents from both demands, not just one considered alone.
func TestRunDemands_SharedRoadReflectsCombinedTraffic(t *testing.T) {
	g := buildBottleneckNetwork(t, graph.Linear(0, 1)) // time == volume

	mp := MultiPopulation{
		Demands: []Demand{
			{Origin: "h1", Destination: "c1", Size: 10},
			{Origin: "h2", Destination: "c2", Size: 15},
		},
		MaxRounds: DefaultMaxRounds,
	}

	result, err := RunDemands(g, mp)
	if err != nil {
		t.Fatalf("RunDemands returned unexpected error: %v", err)
	}
	if !result.Converged {
		t.Fatalf("RunDemands did not converge within %d rounds", DefaultMaxRounds)
	}

	// Reconstruct the A-B edge's final volume from every agent's route,
	// across both demands.
	bottleneckVolume := 0
	for _, dr := range result.PerDemand {
		for _, r := range dr.Routes {
			for _, e := range r.Edges {
				if e.ID == "A-B" {
					bottleneckVolume++
				}
			}
		}
	}
	if want := 25; bottleneckVolume != want {
		t.Fatalf("A-B volume = %d, want %d (10 + 15 agents from both demands)", bottleneckVolume, want)
	}
}

// SC-001 / Acceptance Scenario 2: increasing one demand's size measurably
// changes the other demand's own reported average travel time, proving
// the two flows genuinely share congestion rather than being simulated in
// isolation.
func TestRunDemands_OneDemandSizeAffectsTheOthersAverage(t *testing.T) {
	run := func(sizeH2 int) float64 {
		g := buildBottleneckNetwork(t, graph.Linear(0, 1))
		mp := MultiPopulation{
			Demands: []Demand{
				{Origin: "h1", Destination: "c1", Size: 10},
				{Origin: "h2", Destination: "c2", Size: sizeH2},
			},
			MaxRounds: DefaultMaxRounds,
		}
		result, err := RunDemands(g, mp)
		if err != nil {
			t.Fatalf("RunDemands returned unexpected error: %v", err)
		}
		if !result.Converged {
			t.Fatalf("RunDemands did not converge within %d rounds", DefaultMaxRounds)
		}
		return result.PerDemand[0].AverageTravelTime // h1 -> c1's own average
	}

	small := run(10)
	large := run(50)

	if large <= small {
		t.Fatalf("h1->c1's average travel time did not increase when h2->c2 grew: with small h2 (%g) vs large h2 (%g)", small, large)
	}
}

// FR-008: a demand with no route between its origin and destination must
// produce an error naming that specific demand.
func TestRunDemands_NoRouteNamesTheDemand(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "h1", graph.House)
	mustAddNode(t, g, "c1", graph.Company)
	mustAddNode(t, g, "h2", graph.House)
	mustAddNode(t, g, "c2", graph.Company)
	mustAddEdge(t, g, "h1-c1", "h1", "c1", graph.Constant(5))
	// h2 -> c2 has no edge at all.

	mp := MultiPopulation{
		Demands: []Demand{
			{Origin: "h1", Destination: "c1", Size: 3},
			{Origin: "h2", Destination: "c2", Size: 3},
		},
		MaxRounds: DefaultMaxRounds,
	}

	_, err := RunDemands(g, mp)
	if !errors.Is(err, agent.ErrNoRoute) {
		t.Fatalf("RunDemands with an unroutable demand = %v, want an error wrapping agent.ErrNoRoute", err)
	}
	msg := err.Error()
	for _, want := range []string{"1", "h2", "c2"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("RunDemands error %q does not clearly identify demand 1 (h2 -> c2)", msg)
		}
	}
}

// Acceptance Scenario 1-2 (spec.md User Story 2): each demand's own
// result reflects only its own agents, and one overall result is also
// available.
func TestRunDemands_PerDemandResultsAreIsolated(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "h1", graph.House)
	mustAddNode(t, g, "c1", graph.Company)
	mustAddNode(t, g, "h2", graph.House)
	mustAddNode(t, g, "c2", graph.Company)
	mustAddEdge(t, g, "h1-c1", "h1", "c1", graph.Constant(10))
	mustAddEdge(t, g, "h2-c2", "h2", "c2", graph.Constant(20))

	mp := MultiPopulation{
		Demands: []Demand{
			{Origin: "h1", Destination: "c1", Size: 4},
			{Origin: "h2", Destination: "c2", Size: 6},
		},
		MaxRounds: DefaultMaxRounds,
	}

	result, err := RunDemands(g, mp)
	if err != nil {
		t.Fatalf("RunDemands returned unexpected error: %v", err)
	}
	if !result.Converged {
		t.Fatal("RunDemands did not converge")
	}

	d0, d1 := result.PerDemand[0], result.PerDemand[1]
	if len(d0.Routes) != 4 || d0.TotalTravelTime != 40 || d0.AverageTravelTime != 10 {
		t.Fatalf("demand 0 (h1->c1) = %+v, want 4 routes, total 40, average 10", d0)
	}
	if len(d1.Routes) != 6 || d1.TotalTravelTime != 120 || d1.AverageTravelTime != 20 {
		t.Fatalf("demand 1 (h2->c2) = %+v, want 6 routes, total 120, average 20", d1)
	}

	// Independent, unrelated edges: neither demand's own numbers should
	// leak into the other's.
	if d0.TotalTravelTime+d1.TotalTravelTime != result.TotalTravelTime {
		t.Fatalf("sum of per-demand totals (%g) != overall total (%g)", d0.TotalTravelTime+d1.TotalTravelTime, result.TotalTravelTime)
	}
	if want := (40.0 + 120.0) / 10.0; result.AverageTravelTime != want {
		t.Fatalf("overall average = %g, want %g", result.AverageTravelTime, want)
	}
}

// FR-007: a demand with zero agents reports a trivial result without
// error and without affecting other demands.
func TestRunDemands_ZeroSizeDemandDoesNotAffectOthers(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "h1", graph.House)
	mustAddNode(t, g, "c1", graph.Company)
	mustAddNode(t, g, "h2", graph.House)
	mustAddNode(t, g, "c2", graph.Company)
	mustAddEdge(t, g, "h1-c1", "h1", "c1", graph.Constant(10))
	// h2 -> c2 has no edge — proving the zero-size demand never actually
	// needs to route, since it must not error even though it's unroutable.
	// (No edge is added between h2 and c2 on purpose.)

	mp := MultiPopulation{
		Demands: []Demand{
			{Origin: "h1", Destination: "c1", Size: 5},
			{Origin: "h2", Destination: "c2", Size: 0},
		},
		MaxRounds: DefaultMaxRounds,
	}

	result, err := RunDemands(g, mp)
	if err != nil {
		t.Fatalf("RunDemands returned unexpected error: %v", err)
	}

	zero := result.PerDemand[1]
	if len(zero.Routes) != 0 || zero.TotalTravelTime != 0 || zero.AverageTravelTime != 0 {
		t.Fatalf("zero-size demand result = %+v, want all zero/empty", zero)
	}

	other := result.PerDemand[0]
	if len(other.Routes) != 5 || other.TotalTravelTime != 50 {
		t.Fatalf("non-zero demand result = %+v, want 5 routes, total 50", other)
	}
}

// SC-002: across 3 distinct multi-demand networks, the final assignment
// is a genuine equilibrium — no agent, from any demand, could
// unilaterally switch route and strictly improve. Verified generically
// (works for any topology), the same way feature 003's equilibrium tests
// checked the single-demand case.
func TestRunDemands_ReachesGlobalEquilibrium(t *testing.T) {
	tests := []struct {
		name string
		g    func(t *testing.T) *graph.Graph
		mp   MultiPopulation
	}{
		{
			name: "two demands sharing one bottleneck",
			g:    func(t *testing.T) *graph.Graph { return buildBottleneckNetwork(t, graph.Linear(0, 1)) },
			mp: MultiPopulation{
				Demands: []Demand{
					{Origin: "h1", Destination: "c1", Size: 12},
					{Origin: "h2", Destination: "c2", Size: 8},
				},
				MaxRounds: DefaultMaxRounds,
			},
		},
		{
			name: "three demands, independent routes",
			g: func(t *testing.T) *graph.Graph {
				g := graph.New()
				pairs := []string{"h1", "c1", "h2", "c2", "h3", "c3"}
				for _, id := range pairs {
					if _, err := g.AddNode(id, graph.Intersection); err != nil {
						t.Fatalf("AddNode(%q) returned unexpected error: %v", id, err)
					}
				}
				mustAddEdge(t, g, "h1-c1", "h1", "c1", graph.Linear(0, 1))
				mustAddEdge(t, g, "h2-c2", "h2", "c2", graph.Linear(0, 2))
				mustAddEdge(t, g, "h3-c3", "h3", "c3", graph.Constant(7))
				return g
			},
			mp: MultiPopulation{
				Demands: []Demand{
					{Origin: "h1", Destination: "c1", Size: 5},
					{Origin: "h2", Destination: "c2", Size: 5},
					{Origin: "h3", Destination: "c3", Size: 5},
				},
				MaxRounds: DefaultMaxRounds,
			},
		},
		{
			name: "three demands, two of them sharing a bottleneck and choices",
			g: func(t *testing.T) *graph.Graph {
				g := buildBottleneckNetwork(t, graph.Linear(0, 1))
				// A third demand with its own two parallel routes,
				// untouched by the bottleneck.
				mustAddNode(t, g, "h3", graph.House)
				mustAddNode(t, g, "c3", graph.Company)
				mustAddEdge(t, g, "h3-c3-fast", "h3", "c3", graph.Linear(0, 1))
				mustAddEdge(t, g, "h3-c3-slow", "h3", "c3", graph.Linear(0, 2))
				return g
			},
			mp: MultiPopulation{
				Demands: []Demand{
					{Origin: "h1", Destination: "c1", Size: 10},
					{Origin: "h2", Destination: "c2", Size: 10},
					{Origin: "h3", Destination: "c3", Size: 9},
				},
				MaxRounds: DefaultMaxRounds,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g(t)
			result, err := RunDemands(g, tc.mp)
			if err != nil {
				t.Fatalf("RunDemands returned unexpected error: %v", err)
			}
			if !result.Converged {
				t.Fatalf("RunDemands did not converge within %d rounds", tc.mp.MaxRounds)
			}
			assertGlobalEquilibrium(t, g, result)
		})
	}
}

// assertGlobalEquilibrium verifies that, for the final routes across
// every demand in result, no agent could strictly improve by switching to
// a different route — using the same marginal (including-self) cost
// convention RunDemands itself uses. Works for any network topology.
func assertGlobalEquilibrium(t *testing.T, g *graph.Graph, result MultiAssignmentResult) {
	t.Helper()

	volumes := make(map[string]float64)
	for _, dr := range result.PerDemand {
		for _, r := range dr.Routes {
			addToVolumes(volumes, r)
		}
	}

	for _, dr := range result.PerDemand {
		for _, r := range dr.Routes {
			currentCost, err := routeCost(g, r, volumes)
			if err != nil {
				t.Fatalf("routeCost returned unexpected error: %v", err)
			}

			removeFromVolumes(volumes, r)
			alt, err := agent.ShortestRouteAtVolumes(g, dr.Demand.Origin, dr.Demand.Destination, marginalVolumes(g, volumes))
			if err != nil {
				addToVolumes(volumes, r)
				t.Fatalf("ShortestRouteAtVolumes returned unexpected error: %v", err)
			}
			addToVolumes(volumes, r)

			if alt.TotalTravelTime < currentCost-epsilon {
				t.Fatalf("not an equilibrium: an agent on demand %q->%q (cost %g) could improve to %g",
					dr.Demand.Origin, dr.Demand.Destination, currentCost, alt.TotalTravelTime)
			}
		}
	}
}
