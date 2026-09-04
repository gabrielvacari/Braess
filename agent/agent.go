package agent

import "braess/graph"

// Agent is a car with a fixed origin and destination. It holds no other
// state: every ComputeRoute call is a fresh, independent computation
// (constitution Principle III) — there is no package-level cache and
// Agent does not remember a previously computed Route.
type Agent struct {
	Origin      string // Node ID
	Destination string // Node ID
}

// ComputeRoute computes a's own route over g, independently of any other
// Agent. It is equivalent to ShortestRoute(g, a.Origin, a.Destination, 0)
// — this phase's baseline volume (no congestion modeling yet).
func (a *Agent) ComputeRoute(g *graph.Graph) (Route, error) {
	return ShortestRoute(g, a.Origin, a.Destination, 0)
}
