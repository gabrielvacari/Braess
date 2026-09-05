# Contract Extension: `simulation` package (feature 003 → feature 004)

This is an additive, non-breaking extension of
[feature 003's `simulation-api.md`](../../003-congestion-braess-paradox/contracts/simulation-api.md)
contract, exercising the "expected to grow" clause in its Stability
section. `Population`, `Run`, and `AssignmentResult`'s existing
signatures and behavior are unchanged.

## New types

```go
package simulation

// Demand is one origin/destination pair together with how many agents
// travel between them in a given run.
type Demand struct {
    Origin, Destination string // Node IDs
    Size                 int    // number of agents; 0 is valid
}

// MultiPopulation is several Demands sharing one road network and its
// real congestion.
type MultiPopulation struct {
    Demands   []Demand
    MaxRounds int // same semantics as Population.MaxRounds
}

// DemandResult is one Demand's own outcome within a MultiPopulation run.
type DemandResult struct {
    Demand            Demand
    Routes            []agent.Route
    TotalTravelTime   float64
    AverageTravelTime float64
}

// MultiAssignmentResult is the outcome of running a MultiPopulation to a
// stable (or round-limited) shared assignment.
type MultiAssignmentResult struct {
    PerDemand         []DemandResult // same order as MultiPopulation.Demands
    TotalTravelTime   float64        // across every agent, every demand
    AverageTravelTime float64
    Converged         bool
    Rounds            int
}
```

## New function

```go
// RunDemands computes a shared route assignment for every agent across
// every demand in mp.Demands: agents are loaded (all of one demand, then
// the next, in mp.Demands order) against real accumulating traffic from
// every demand combined, then refined round by round exactly as Run does
// — every agent, regardless of which demand it belongs to, may switch to
// a strictly better route given everyone else's current choice, until no
// agent improves or mp.MaxRounds is reached.
//
// Returns an error, rather than swallowing it, if any agent's route
// computation fails, naming the specific Demand (by index and
// origin/destination) it belongs to — e.g. no route exists for that
// demand's origin/destination (wraps agent.ErrNoRoute), or an edge's
// travel-time function itself errors.
func RunDemands(g *graph.Graph, mp MultiPopulation) (MultiAssignmentResult, error)
```

## Internal generalization (not part of the public contract)

`Run` and `RunDemands` both delegate to the same internal
`loadIncrementally`/`refine` core, now parameterized by a flat per-agent
`[]odPair{origin, destination string; demandIndex int}` instead of one
shared pair — `Run` builds an `odPair` list where every entry repeats
`Population.Origin`/`Destination` (`demandIndex` unused). This is the same
"factor into one shared core" pattern feature 003 used for
`agent.ShortestRoute`/`ShortestRouteAtVolumes`.

## Compatibility

Feature 003's existing tests (`simulation/run_test.go`,
`simulation/braess_test.go`) must continue to pass unmodified, including
the exact 65 → 80 Braess numbers — this refactor changes internal
structure only, not `Run`'s exported signature, field set, or observable
behavior.
