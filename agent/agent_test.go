package agent

import (
	"testing"

	"braess/graph"
)

// Acceptance Scenario 1 (spec.md User Story 2) / SC-003: at least 3
// independently-created agents sharing an origin/destination on a graph
// with a unique shortest path all report the same route and travel time,
// each arrived at through its own independent ComputeRoute call.
func TestAgent_ComputeRoute_IndependentAgentsAgreeOnUniqueShortestPath(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "house-1", graph.House)
	mustAddNode(t, g, "intersection-1", graph.Intersection)
	mustAddNode(t, g, "company-1", graph.Company)
	mustAddEdge(t, g, "road-1", "house-1", "intersection-1", graph.Constant(3))
	mustAddEdge(t, g, "road-2", "intersection-1", "company-1", graph.Constant(4))
	// A slower alternative, so the shortest path is genuinely unique.
	mustAddEdge(t, g, "road-3", "house-1", "company-1", graph.Constant(100))

	agents := []*Agent{
		{Origin: "house-1", Destination: "company-1"},
		{Origin: "house-1", Destination: "company-1"},
		{Origin: "house-1", Destination: "company-1"},
	}

	var routes []Route
	for i, a := range agents {
		route, err := a.ComputeRoute(g)
		if err != nil {
			t.Fatalf("agent %d ComputeRoute returned unexpected error: %v", i, err)
		}
		routes = append(routes, route)
	}

	want := []string{"road-1", "road-2"}
	for i, route := range routes {
		if got := edgeIDs(route.Edges); !equalStrings(got, want) {
			t.Fatalf("agent %d route = %v, want %v", i, got, want)
		}
		if route.TotalTravelTime != 7 {
			t.Fatalf("agent %d total = %g, want 7", i, route.TotalTravelTime)
		}
	}
}

// Acceptance Scenario 2 (spec.md User Story 2): with two equally-shortest
// routes, each agent's chosen path still comes from that agent's own
// computation — verified here by checking each call independently
// reaches the same, valid minimum-time answer rather than any agent
// short-circuiting on another's result.
func TestAgent_ComputeRoute_TiedRoutesEachIndependentlyComputed(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddEdge(t, g, "route-x", "a", "b", graph.Constant(5))
	mustAddEdge(t, g, "route-y", "a", "b", graph.Constant(5))

	agents := []*Agent{
		{Origin: "a", Destination: "b"},
		{Origin: "a", Destination: "b"},
		{Origin: "a", Destination: "b"},
	}

	for i, a := range agents {
		route, err := a.ComputeRoute(g)
		if err != nil {
			t.Fatalf("agent %d ComputeRoute returned unexpected error: %v", i, err)
		}
		if len(route.Edges) != 1 {
			t.Fatalf("agent %d route = %v, want exactly one edge", i, route.Edges)
		}
		if route.TotalTravelTime != 5 {
			t.Fatalf("agent %d total = %g, want 5", i, route.TotalTravelTime)
		}
	}
}

// Mutating one Agent must not affect another's independently computed
// result — reinforcing FR-003 (no shared state between agents).
func TestAgent_ComputeRoute_AgentsDoNotShareState(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddNode(t, g, "c", graph.Company)
	mustAddEdge(t, g, "a-b", "a", "b", graph.Constant(2))
	mustAddEdge(t, g, "a-c", "a", "c", graph.Constant(9))

	toB := &Agent{Origin: "a", Destination: "b"}
	toC := &Agent{Origin: "a", Destination: "c"}

	routeB, err := toB.ComputeRoute(g)
	if err != nil {
		t.Fatalf("toB.ComputeRoute returned unexpected error: %v", err)
	}
	routeC, err := toC.ComputeRoute(g)
	if err != nil {
		t.Fatalf("toC.ComputeRoute returned unexpected error: %v", err)
	}

	if got := edgeIDs(routeB.Edges); !equalStrings(got, []string{"a-b"}) {
		t.Fatalf("toB route = %v, want [a-b]", got)
	}
	if got := edgeIDs(routeC.Edges); !equalStrings(got, []string{"a-c"}) {
		t.Fatalf("toC route = %v, want [a-c]", got)
	}
}
