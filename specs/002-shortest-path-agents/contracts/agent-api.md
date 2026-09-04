# Contract: `agent` package public API

Like feature 001's `graph` package, this is a library — its contract is
the exported Go API. `agent` depends on `graph` (feature 001) and the Go
standard library only.

## Types

```go
package agent

import "braess/graph"

// Route is the ordered sequence of edges an agent will take from its
// origin to its destination, and the route's total travel time at the
// volume it was computed for.
type Route struct {
    Edges           []graph.Edge
    TotalTravelTime float64
}

// Agent is a car with a fixed origin and destination. It holds no other
// state: every ComputeRoute call is an independent computation (FR-003).
type Agent struct {
    Origin      string // Node ID
    Destination string // Node ID
}
```

## Errors

```go
// ErrNoRoute indicates no path connects the given origin to the given
// destination (FR-005). Check with errors.Is.
var ErrNoRoute = errors.New("agent: no route exists")
```

## Functions

```go
// ShortestRoute computes the minimum-total-travel-time route from `from`
// to `to` in g, evaluating every edge's travel time at the given volume
// (this phase's only caller, Agent.ComputeRoute, always passes 0 — see
// research.md decision #3).
//
// Returns a zero-edge, zero-time Route if from == to (FR-006).
// Returns an error wrapping ErrNoRoute if no path exists (FR-005).
// Returns (and does not swallow) any error an edge's TravelTimeFunc
// itself produces while the route is being computed (FR-008).
// Parallel edges between the same node pair are each considered as a
// distinct option (FR-004).
func ShortestRoute(g *graph.Graph, from, to string, volume float64) (Route, error)
```

## Methods

```go
// ComputeRoute computes a's own route over g, independently of any other
// Agent (FR-003). Equivalent to ShortestRoute(g, a.Origin, a.Destination, 0).
func (a *Agent) ComputeRoute(g *graph.Graph) (Route, error)
```

## Error handling

Same convention as `graph`: fallible operations return a Go `error` as the
last return value; no panics. `ErrNoRoute` is the one sentinel this
contract introduces, specifically so callers (the CLI, tests, and later
phases) can distinguish "no route exists" from any other failure via
`errors.Is(err, agent.ErrNoRoute)`.

## Stability

This contract covers roadmap Phase 2 only. `ShortestRoute` taking `volume`
as a parameter (rather than a hardcoded `0`) is deliberate so Phase 3
(congestion) can reuse it with a real, non-zero volume without breaking
this signature — see research.md decision #3.
