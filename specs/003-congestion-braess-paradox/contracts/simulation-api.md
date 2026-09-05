# Contract: `simulation` package public API

Like `graph` and `agent`, this is a library — its contract is the exported
Go API. `simulation` depends on `agent`, `graph`, and the Go standard
library only.

## Types

```go
package simulation

import "braess/agent"

// DefaultMaxRounds is a generous safety cap on refinement rounds. Best-
// response dynamics on a congestion game is guaranteed to converge in
// finite steps (research.md decision #1); this cap exists for pathological
// inputs, not as the primary convergence mechanism.
const DefaultMaxRounds = 1000

// Population is many identical agents sharing one fixed origin and
// destination.
type Population struct {
    Origin, Destination string // Node IDs
    Size                int    // number of agents; 0 is valid
    MaxRounds           int    // refinement-round budget; 0 = loading only, no refinement
}

// AssignmentResult is the outcome of running a Population to a stable (or
// round-limited) route assignment.
type AssignmentResult struct {
    Routes            []agent.Route
    TotalTravelTime   float64
    AverageTravelTime float64
    Converged         bool
    Rounds            int
}
```

## Functions

```go
// Run computes a route assignment for p on g: agents are loaded one at a
// time against real accumulating traffic (each via
// agent.ShortestRouteAtVolumes), then refined round by round — each
// round, every agent may switch to a strictly better route given
// everyone else's current choice — until no agent improves or
// p.MaxRounds is reached.
//
// Returns an error (propagated, not swallowed) if any agent's route
// computation fails — e.g. no route exists between p.Origin and
// p.Destination (wraps agent.ErrNoRoute), or an edge's travel-time
// function itself errors.
func Run(g *graph.Graph, p Population) (AssignmentResult, error)
```

## Error handling

Same convention as `graph` and `agent`: a Go `error` as the last return
value, no panics, no swallowed errors from lower layers.

## Stability

This contract covers roadmap Phase 3 only. `Population` is deliberately
minimal (one fixed origin/destination, a size, a round budget) — Phase 4's
multiple distinct origin/destination pairs will need a richer population
shape, expected to be a new type or an extension of this one rather than a
breaking change to `Run`'s signature.
