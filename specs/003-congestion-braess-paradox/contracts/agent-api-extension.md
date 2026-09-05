# Contract Extension: `agent` package (feature 002 → feature 003)

This is an additive, non-breaking extension of
[feature 002's `graph-api.md`](../../002-shortest-path-agents/contracts/agent-api.md)
contract, exercising the "expected to grow" clause in its Stability
section. `ShortestRoute`'s existing signature and behavior are unchanged.

## New function

```go
package agent

// ShortestRouteAtVolumes computes the minimum-total-travel-time route from
// from to to in g, evaluating each edge's travel time at its own current
// traffic volume: volumes[edge.ID], or 0 for any edge not present in
// volumes. This is what feature 003's population-level congestion uses,
// since different roads carry different real traffic simultaneously —
// unlike ShortestRoute's single flat baseline.
//
// Same error/edge-case behavior as ShortestRoute: a zero-edge route if
// from == to; an error wrapping ErrNoRoute if no path exists; propagates
// (never swallows) any error an edge's TravelTimeFunc produces; parallel
// edges are each considered as distinct options.
func ShortestRouteAtVolumes(g *graph.Graph, from, to string, volumes map[string]float64) (Route, error)
```

## Internal refactor (not part of the public contract)

`ShortestRoute` and `ShortestRouteAtVolumes` both delegate to one internal
Dijkstra core parameterized by a per-edge volume lookup function, so there
is exactly one pathfinding implementation to reason about:

```go
func shortestRoute(g *graph.Graph, from, to string, volumeOf func(edgeID string) float64) (Route, error)

func ShortestRoute(g *graph.Graph, from, to string, volume float64) (Route, error) {
    return shortestRoute(g, from, to, func(string) float64 { return volume })
}

func ShortestRouteAtVolumes(g *graph.Graph, from, to string, volumes map[string]float64) (Route, error) {
    return shortestRoute(g, from, to, func(id string) float64 { return volumes[id] })
}
```

## Compatibility

Feature 002's existing tests (`agent/shortestpath_test.go`,
`agent/agent_test.go`) must continue to pass unmodified — this refactor
changes internal structure only, not `ShortestRoute`'s exported signature
or observable behavior.
