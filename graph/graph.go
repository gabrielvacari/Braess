package graph

import "fmt"

// Graph is a road network: the full collection of nodes and edges,
// queryable as a whole. The zero value is not ready to use — construct one
// with New.
type Graph struct {
	nodes map[string]Node
	edges map[string]Edge
}

// New returns an empty, ready-to-use Graph.
func New() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string]Edge),
	}
}

// AddNode adds a node with the given id and type to the graph.
// It errors if id is empty or already present.
func (g *Graph) AddNode(id string, t NodeType) (Node, error) {
	if id == "" {
		return Node{}, fmt.Errorf("graph: node id must not be empty")
	}
	if _, exists := g.nodes[id]; exists {
		return Node{}, fmt.Errorf("graph: node %q already exists", id)
	}
	n := Node{ID: id, Type: t}
	g.nodes[id] = n
	return n, nil
}

// AddEdge adds a directed edge between two existing nodes. It errors if id
// is empty or already present, or if from/to do not reference existing
// nodes. Multiple edges may share the same (from, to) pair — parallel
// routes are a first-class case, needed to model adding a redundant road.
func (g *Graph) AddEdge(id, from, to string, length, capacity float64, tt TravelTimeFunc) (Edge, error) {
	if id == "" {
		return Edge{}, fmt.Errorf("graph: edge id must not be empty")
	}
	if _, exists := g.edges[id]; exists {
		return Edge{}, fmt.Errorf("graph: edge %q already exists", id)
	}
	if _, exists := g.nodes[from]; !exists {
		return Edge{}, fmt.Errorf("graph: edge %q references unknown from-node %q", id, from)
	}
	if _, exists := g.nodes[to]; !exists {
		return Edge{}, fmt.Errorf("graph: edge %q references unknown to-node %q", id, to)
	}
	e := Edge{ID: id, From: from, To: to, Length: length, Capacity: capacity, TravelTime: tt}
	g.edges[id] = e
	return e, nil
}

// Nodes returns every node in the graph, in unspecified order.
func (g *Graph) Nodes() []Node {
	nodes := make([]Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

// Edges returns every edge in the graph, in unspecified order.
func (g *Graph) Edges() []Edge {
	edges := make([]Edge, 0, len(g.edges))
	for _, e := range g.edges {
		edges = append(edges, e)
	}
	return edges
}

// TravelTime evaluates the travel time of the edge identified by edgeID at
// the given traffic volume. It errors if edgeID does not exist; the edge's
// own TravelTimeFunc is responsible for erroring on a negative volume.
func (g *Graph) TravelTime(edgeID string, volume float64) (float64, error) {
	e, exists := g.edges[edgeID]
	if !exists {
		return 0, fmt.Errorf("graph: edge %q does not exist", edgeID)
	}
	return e.TravelTime(volume)
}
