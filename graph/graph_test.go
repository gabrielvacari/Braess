package graph

import "testing"

// Acceptance Scenario 1 (spec.md User Story 1): adding a node makes it
// appear in the graph's node list with its declared type.
func TestAddNode_AppearsInNodeList(t *testing.T) {
	g := New()

	if _, err := g.AddNode("h1", House); err != nil {
		t.Fatalf("AddNode returned unexpected error: %v", err)
	}

	nodes := g.Nodes()
	if len(nodes) != 1 {
		t.Fatalf("Nodes() = %d nodes, want 1", len(nodes))
	}
	if nodes[0].ID != "h1" || nodes[0].Type != House {
		t.Fatalf("Nodes()[0] = %+v, want {ID: h1, Type: House}", nodes[0])
	}
}

func TestAddNode_RejectsEmptyID(t *testing.T) {
	g := New()
	if _, err := g.AddNode("", Intersection); err == nil {
		t.Fatal("AddNode(\"\") = nil error, want error")
	}
}

func TestAddNode_RejectsDuplicateID(t *testing.T) {
	g := New()
	if _, err := g.AddNode("n1", Intersection); err != nil {
		t.Fatalf("first AddNode returned unexpected error: %v", err)
	}
	if _, err := g.AddNode("n1", House); err == nil {
		t.Fatal("second AddNode with duplicate id = nil error, want error")
	}
}

// Acceptance Scenario 2 (spec.md User Story 1): an edge between two
// existing nodes is retrievable and reports its length/capacity.
func TestAddEdge_RetrievableWithAttributes(t *testing.T) {
	g := New()
	mustAddNode(t, g, "a", House)
	mustAddNode(t, g, "b", Company)

	if _, err := g.AddEdge("a-b", "a", "b", 10, 100, Constant(5)); err != nil {
		t.Fatalf("AddEdge returned unexpected error: %v", err)
	}

	edges := g.Edges()
	if len(edges) != 1 {
		t.Fatalf("Edges() = %d edges, want 1", len(edges))
	}
	got := edges[0]
	if got.ID != "a-b" || got.From != "a" || got.To != "b" || got.Length != 10 || got.Capacity != 100 {
		t.Fatalf("Edges()[0] = %+v, want {ID: a-b, From: a, To: b, Length: 10, Capacity: 100}", got)
	}
}

// Edge Case (spec.md): an edge referencing a node id that doesn't exist
// must be rejected, not silently created.
func TestAddEdge_RejectsUnknownNodeReference(t *testing.T) {
	g := New()
	mustAddNode(t, g, "a", House)

	if _, err := g.AddEdge("a-b", "a", "does-not-exist", 1, 1, Constant(1)); err == nil {
		t.Fatal("AddEdge with unknown to-node = nil error, want error")
	}
	if _, err := g.AddEdge("x-a", "does-not-exist", "a", 1, 1, Constant(1)); err == nil {
		t.Fatal("AddEdge with unknown from-node = nil error, want error")
	}
}

func TestAddEdge_RejectsEmptyOrDuplicateID(t *testing.T) {
	g := New()
	mustAddNode(t, g, "a", House)
	mustAddNode(t, g, "b", Company)

	if _, err := g.AddEdge("", "a", "b", 1, 1, Constant(1)); err == nil {
		t.Fatal("AddEdge(\"\") = nil error, want error")
	}
	if _, err := g.AddEdge("a-b", "a", "b", 1, 1, Constant(1)); err != nil {
		t.Fatalf("first AddEdge returned unexpected error: %v", err)
	}
	if _, err := g.AddEdge("a-b", "a", "b", 1, 1, Constant(1)); err == nil {
		t.Fatal("second AddEdge with duplicate id = nil error, want error")
	}
}

// Edge Case (spec.md): two separate roads connecting the same pair of
// nodes must be allowed and independently queryable — this is exactly what
// Braess's Paradox experiments need: adding a redundant route.
func TestAddEdge_AllowsParallelEdgesBetweenSameNodes(t *testing.T) {
	g := New()
	mustAddNode(t, g, "a", House)
	mustAddNode(t, g, "b", Company)

	if _, err := g.AddEdge("road1", "a", "b", 10, 100, Constant(5)); err != nil {
		t.Fatalf("first AddEdge returned unexpected error: %v", err)
	}
	if _, err := g.AddEdge("road2", "a", "b", 20, 50, Constant(3)); err != nil {
		t.Fatalf("second parallel AddEdge returned unexpected error: %v", err)
	}

	edges := g.Edges()
	if len(edges) != 2 {
		t.Fatalf("Edges() = %d edges, want 2", len(edges))
	}
}

func mustAddNode(t *testing.T, g *Graph, id string, nt NodeType) {
	t.Helper()
	if _, err := g.AddNode(id, nt); err != nil {
		t.Fatalf("AddNode(%q) returned unexpected error: %v", id, err)
	}
}
