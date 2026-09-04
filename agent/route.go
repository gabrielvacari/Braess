// Package agent decides how a car gets from a fixed origin to a fixed
// destination over a *graph.Graph. It depends only on the graph package
// and the Go standard library — never on any UI/rendering package
// (constitution Principle I) — and every Agent computes its own route
// through its own independent call, with no shared or cached state
// between agents (constitution Principle III).
package agent

import "braess/graph"

// Route is the ordered sequence of edges an agent will take from its
// origin to its destination, together with the route's total travel time
// at the volume it was computed for. Edges is empty, and
// TotalTravelTime is 0, when the origin and destination are the same
// node.
type Route struct {
	Edges           []graph.Edge
	TotalTravelTime float64
}
