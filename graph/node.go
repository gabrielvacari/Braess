// Package graph is the Braess traffic simulator's engine: an in-memory,
// UI-independent model of a road network (nodes and edges) and the
// travel-time behavior of its roads under traffic. It has no dependency on
// any rendering or UI package (constitution Principle I) and is meant to be
// exercised directly, e.g. from tests or the cmd/graphcli terminal command.
package graph

// NodeType classifies what a Node represents in the road network.
type NodeType int

const (
	// Intersection is a point where roads meet; neither an origin nor a
	// destination.
	Intersection NodeType = iota
	// House is an origin — where agents spawn, from roadmap Phase 2 onward.
	House
	// Company is a destination — where agents are headed, from roadmap
	// Phase 2 onward.
	Company
)

// String returns a human-readable name for t, used by cmd/graphcli and in
// error messages.
func (t NodeType) String() string {
	switch t {
	case Intersection:
		return "intersection"
	case House:
		return "house"
	case Company:
		return "company"
	default:
		return "unknown"
	}
}

// Node is an intersection, house, or company in the road network.
type Node struct {
	ID   string
	Type NodeType
}
