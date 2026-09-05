# Phase 1 Data Model: Multiple Origins and Destinations

Derived from spec [Key Entities](./spec.md#key-entities) and the
Functional Requirements. Builds on feature 003's `Population`/`Run`/
`AssignmentResult` (all unchanged) and feature 002's `agent.Route`.

## Demand

| Field | Type | Notes |
|---|---|---|
| `Origin` | `string` (Node ID) | A house, typically, but any node id is accepted (consistent with earlier phases' node-agnostic routing). |
| `Destination` | `string` (Node ID) | A company, typically. |
| `Size` | `int` | Number of agents traveling this pair. `0` is valid (FR-007). |

**Validation rules**:
- `Size` must be `>= 0`.
- Two `Demand`s may share the same `Origin`/`Destination` (spec Edge
  Cases) — nothing deduplicates them; their agents still pool together for
  congestion purposes because they resolve to the same `origin`/
  `destination` when building the internal per-agent list.

## MultiPopulation

| Field | Type | Notes |
|---|---|---|
| `Demands` | `[]Demand` | The declared origin/destination pairs for this run. |
| `MaxRounds` | `int` | Same semantics as `Population.MaxRounds` (feature 003): required, `0` means "loading only, no refinement," always reporting `Converged: false`. |

## DemandResult

| Field | Type | Notes |
|---|---|---|
| `Demand` | `Demand` | Echoes the input this result corresponds to. |
| `Routes` | `[]agent.Route` | This demand's own agents' final routes, `len(Routes) == Demand.Size`. |
| `TotalTravelTime` | `float64` | Sum of this demand's own agents' route costs at the final shared volumes. |
| `AverageTravelTime` | `float64` | `TotalTravelTime / Demand.Size`; `0` when `Size == 0`. |

## MultiAssignmentResult

| Field | Type | Notes |
|---|---|---|
| `PerDemand` | `[]DemandResult` | Same length and order as the input `MultiPopulation.Demands`. |
| `TotalTravelTime` | `float64` | Sum of every agent's route cost, across every demand, at the final shared volumes. |
| `AverageTravelTime` | `float64` | `TotalTravelTime / (sum of all Demands' Size)`; `0` when that sum is `0`. |
| `Converged` | `bool` | Same semantics as `AssignmentResult.Converged` (feature 003), now over the whole shared, multi-demand assignment. |
| `Rounds` | `int` | Refinement rounds actually run, shared across all demands (there is one refinement process, not one per demand). |

**Validation rules**:
- `TotalTravelTime`/`AverageTravelTime` (both the overall ones and each
  `DemandResult`'s own) are derived from final routes and volumes — never
  independently settable.
- `sum(PerDemand[i].TotalTravelTime for all i) == TotalTravelTime` — the
  overall total is exactly the sum of every demand's own total, since
  every agent belongs to exactly one demand.

## Relationships

```text
MultiPopulation.Demands[i]        ──> Demand
DemandResult.Demand               ──> the Demand it corresponds to (same index, same run)
DemandResult.Routes[j]            ──> agent.Route (feature 002), one of this demand's own agents
MultiAssignmentResult.PerDemand[i]──> DemandResult
```

Internally (not exported — see
[contracts/simulation-api-extension.md](./contracts/simulation-api-extension.md)),
every agent across every demand is represented as one `odPair{origin,
destination, demandIndex}` in a single flat list that the generalized
`loadIncrementally`/`refine` operate over — this is what makes congestion
genuinely shared: all agents, regardless of demand, read and write the
same `volumes map[string]float64`.

## No state transitions beyond a single `RunDemands` call

Same as feature 003: a `MultiPopulation` value is immutable input;
`RunDemands` does not mutate it, and nothing persists between separate
calls.
