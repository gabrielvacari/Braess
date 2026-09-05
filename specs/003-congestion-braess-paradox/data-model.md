# Phase 1 Data Model: Congestion and Braess's Paradox

Derived from spec [Key Entities](./spec.md#key-entities) and the Functional
Requirements. Builds on feature 001's `graph.*` and feature 002's
`agent.Route`/`agent.Agent`, and the new `agent.ShortestRouteAtVolumes`
(see [contracts/agent-api-extension.md](./contracts/agent-api-extension.md)).

## Population

| Field | Type | Notes |
|---|---|---|
| `Origin` | `string` (Node ID) | Fixed for this phase, same constraint as feature 002. |
| `Destination` | `string` (Node ID) | Fixed for this phase. |
| `Size` | `int` | Number of identical agents sharing this origin/destination. `0` is valid (FR-010). |
| `MaxRounds` | `int` | Refinement-round budget; `0` means "incremental loading only, no refinement" (research.md decision #5). Required — no implicit default. |

**Validation rules**:
- `Size` must be `>= 0`.
- `MaxRounds` must be `>= 0`.

## Edge Volume (internal, not exported)

Not a public type — represented internally as `map[string]float64` (Edge
ID → current agent count on that edge) while `Run` executes. Exists only
transiently during a `Run` call; nothing about it is part of the public
contract.

## AssignmentResult

| Field | Type | Notes |
|---|---|---|
| `Routes` | `[]agent.Route` | One entry per agent, `len(Routes) == Population.Size`. |
| `TotalTravelTime` | `float64` | Sum of every agent's own route's travel time, evaluated at the *final* per-edge volumes. |
| `AverageTravelTime` | `float64` | `TotalTravelTime / Size`; `0` when `Size == 0`. |
| `Converged` | `bool` | `true` only if a refinement round found zero agents wanting to switch before `MaxRounds` was reached. `false` if `MaxRounds == 0`, or the round limit was hit with agents still switching (FR-006). |
| `Rounds` | `int` | How many refinement rounds actually ran (`0` to `MaxRounds`). Does not count the initial incremental-loading pass. |

**Validation rules**:
- `TotalTravelTime` and `AverageTravelTime` are derived from `Routes` at
  the final volumes — never independently settable.
- `Converged == true` implies `Rounds <= MaxRounds` and that the last
  round found no improving agent; `Converged == false` implies either
  `MaxRounds == 0` or the loop exhausted `MaxRounds` rounds still finding
  improvements.

## Relationships

```text
Population.Origin, Population.Destination ──> graph.Node.ID
AssignmentResult.Routes[i]                 ──> agent.Route (one per agent, i in [0, Population.Size))
agent.Route.Edges                          ──> []graph.Edge (as in feature 002)
```

`Run(g *graph.Graph, p Population) (AssignmentResult, error)` is the one
operation this feature adds at the population level — see
[contracts/simulation-api.md](./contracts/simulation-api.md). It has no
persistent state of its own: an `AssignmentResult` is a snapshot of one
`Run` call, not a live object that mutates.

## No state transitions beyond a single `Run` call

A `Population` value is immutable input; `Run` does not mutate it. Nothing
in this feature keeps state between separate `Run` calls — running the
same `Population` twice on the same `*graph.Graph` performs the full
incremental-load-then-refine process again from scratch (deterministically
reaching the same result, since routing tie-breaks are deterministic per
feature 002's own Assumptions).
