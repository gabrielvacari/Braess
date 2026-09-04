package agent

import (
	"errors"
	"testing"

	"braess/graph"
)

// SC-001: across at least 5 distinct test networks with a uniquely
// shortest path, ShortestRoute matches the mathematically correct
// minimum-travel-time route.
func TestShortestRoute_PicksMinimumTravelTime(t *testing.T) {
	tests := []struct {
		name      string
		build     func(t *testing.T) *graph.Graph
		from, to  string
		wantEdges []string // edge IDs, in order
		wantTotal float64
	}{
		{
			// Acceptance Scenario 1: a single possible path.
			name: "single path a-b-c",
			build: func(t *testing.T) *graph.Graph {
				g := graph.New()
				mustAddNode(t, g, "a", graph.Intersection)
				mustAddNode(t, g, "b", graph.Intersection)
				mustAddNode(t, g, "c", graph.Intersection)
				mustAddEdge(t, g, "ab", "a", "b", graph.Constant(3))
				mustAddEdge(t, g, "bc", "b", "c", graph.Constant(4))
				return g
			},
			from: "a", to: "c",
			wantEdges: []string{"ab", "bc"},
			wantTotal: 7,
		},
		{
			// Acceptance Scenario 2: multiple paths, different totals —
			// the direct, slower-looking edge is actually faster here.
			name: "direct edge beats two-hop path",
			build: func(t *testing.T) *graph.Graph {
				g := graph.New()
				mustAddNode(t, g, "a", graph.House)
				mustAddNode(t, g, "b", graph.Intersection)
				mustAddNode(t, g, "c", graph.Company)
				mustAddEdge(t, g, "ac", "a", "c", graph.Constant(5))
				mustAddEdge(t, g, "ab", "a", "b", graph.Constant(3))
				mustAddEdge(t, g, "bc", "b", "c", graph.Constant(4))
				return g
			},
			from: "a", to: "c",
			wantEdges: []string{"ac"},
			wantTotal: 5,
		},
		{
			name: "two-hop path beats a slow direct edge",
			build: func(t *testing.T) *graph.Graph {
				g := graph.New()
				mustAddNode(t, g, "a", graph.House)
				mustAddNode(t, g, "b", graph.Intersection)
				mustAddNode(t, g, "c", graph.Company)
				mustAddEdge(t, g, "ac", "a", "c", graph.Constant(45))
				mustAddEdge(t, g, "ab", "a", "b", graph.Constant(3))
				mustAddEdge(t, g, "bc", "b", "c", graph.Constant(4))
				return g
			},
			from: "a", to: "c",
			wantEdges: []string{"ab", "bc"},
			wantTotal: 7,
		},
		{
			name: "diamond network picks the cheaper side",
			build: func(t *testing.T) *graph.Graph {
				g := graph.New()
				for _, id := range []string{"a", "b", "c", "d"} {
					mustAddNode(t, g, id, graph.Intersection)
				}
				mustAddEdge(t, g, "ab", "a", "b", graph.Constant(1))
				mustAddEdge(t, g, "bd", "b", "d", graph.Constant(1))
				mustAddEdge(t, g, "ac", "a", "c", graph.Constant(1))
				mustAddEdge(t, g, "cd", "c", "d", graph.Constant(10))
				return g
			},
			from: "a", to: "d",
			wantEdges: []string{"ab", "bd"},
			wantTotal: 2,
		},
		{
			name: "five-node chain with a shortcut",
			build: func(t *testing.T) *graph.Graph {
				g := graph.New()
				for _, id := range []string{"a", "b", "c", "d", "e"} {
					mustAddNode(t, g, id, graph.Intersection)
				}
				mustAddEdge(t, g, "ab", "a", "b", graph.Constant(2))
				mustAddEdge(t, g, "bc", "b", "c", graph.Constant(2))
				mustAddEdge(t, g, "cd", "c", "d", graph.Constant(2))
				mustAddEdge(t, g, "de", "d", "e", graph.Constant(2))
				mustAddEdge(t, g, "ae", "a", "e", graph.Constant(100))
				mustAddEdge(t, g, "shortcut", "b", "e", graph.Constant(3))
				return g
			},
			from: "a", to: "e",
			wantEdges: []string{"ab", "shortcut"},
			wantTotal: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.build(t)
			route, err := ShortestRoute(g, tc.from, tc.to, 0)
			if err != nil {
				t.Fatalf("ShortestRoute(%q, %q) returned unexpected error: %v", tc.from, tc.to, err)
			}
			if got := edgeIDs(route.Edges); !equalStrings(got, tc.wantEdges) {
				t.Fatalf("ShortestRoute(%q, %q) edges = %v, want %v", tc.from, tc.to, got, tc.wantEdges)
			}
			if route.TotalTravelTime != tc.wantTotal {
				t.Fatalf("ShortestRoute(%q, %q) total = %g, want %g", tc.from, tc.to, route.TotalTravelTime, tc.wantTotal)
			}
		})
	}
}

// Edge Case (spec.md): no path exists between origin and destination.
func TestShortestRoute_NoRouteExists(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	// no edges at all

	_, err := ShortestRoute(g, "a", "b", 0)
	if !errors.Is(err, ErrNoRoute) {
		t.Fatalf("ShortestRoute with no connecting edges = %v, want an error wrapping ErrNoRoute", err)
	}
}

// FR-006 / Edge Case: origin and destination are the same node.
func TestShortestRoute_OriginEqualsDestination(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)

	route, err := ShortestRoute(g, "a", "a", 0)
	if err != nil {
		t.Fatalf("ShortestRoute(a, a) returned unexpected error: %v", err)
	}
	if len(route.Edges) != 0 {
		t.Fatalf("ShortestRoute(a, a) edges = %v, want empty", route.Edges)
	}
	if route.TotalTravelTime != 0 {
		t.Fatalf("ShortestRoute(a, a) total = %g, want 0", route.TotalTravelTime)
	}
}

// FR-004 / Edge Case: parallel edges between the same node pair are
// distinct options; the cheaper one must be preferred.
func TestShortestRoute_ParallelEdgesPicksCheaper(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)
	mustAddEdge(t, g, "slow-road", "a", "b", graph.Constant(20))
	mustAddEdge(t, g, "fast-road", "a", "b", graph.Constant(5))

	route, err := ShortestRoute(g, "a", "b", 0)
	if err != nil {
		t.Fatalf("ShortestRoute returned unexpected error: %v", err)
	}
	if got := edgeIDs(route.Edges); !equalStrings(got, []string{"fast-road"}) {
		t.Fatalf("ShortestRoute edges = %v, want [fast-road]", got)
	}
	if route.TotalTravelTime != 5 {
		t.Fatalf("ShortestRoute total = %g, want 5", route.TotalTravelTime)
	}
}

// FR-008 / Edge Case: a failing travel-time evaluation must be
// propagated, not silently treated as absent or instantly fast.
func TestShortestRoute_PropagatesTravelTimeError(t *testing.T) {
	g := graph.New()
	mustAddNode(t, g, "a", graph.House)
	mustAddNode(t, g, "b", graph.Company)

	broken := func(volume float64) (float64, error) {
		return 0, errors.New("boom")
	}
	mustAddEdge(t, g, "broken-road", "a", "b", broken)

	_, err := ShortestRoute(g, "a", "b", 0)
	if err == nil {
		t.Fatal("ShortestRoute over a broken edge = nil error, want error")
	}
	if errors.Is(err, ErrNoRoute) {
		t.Fatalf("ShortestRoute over a broken edge = %v, want the underlying error, not ErrNoRoute", err)
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

func edgeIDs(edges []graph.Edge) []string {
	ids := make([]string, len(edges))
	for i, e := range edges {
		ids[i] = e.ID
	}
	return ids
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
