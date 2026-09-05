# Phase 1 Data Model: Web Frontend

Two related but distinct models exist here: the **server's DTOs**
(`cmd/server`, Go, mapping to/from `graph`/`agent`/`simulation` types) and
the **client's edit model** (`web/src/model`, TypeScript, what the user is
actually drawing). See [contracts/api-contract.md](./contracts/api-contract.md)
for the exact JSON shape connecting them.

## Server DTOs (`cmd/server/dto.go`)

### NodeDTO

| Field | Type | Maps to |
|---|---|---|
| `ID` | `string` | `graph.Node.ID` |
| `Type` | `string`: `"house"` \| `"company"` \| `"intersection"` | `graph.NodeType` |

### TravelTimeDTO

| Field | Type | Notes |
|---|---|---|
| `Type` | `string`: `"linear"` \| `"constant"` | Selects which `graph` constructor to call. |
| `FreeFlow` | `float64` | Used by `"linear"` → `graph.Linear(FreeFlow, Slope)`. |
| `Slope` | `float64` | Used by `"linear"` only. |
| `Value` | `float64` | Used by `"constant"` → `graph.Constant(Value)`. |

**Validation**: `Type` must be one of the two known values; the fields it
doesn't use are ignored rather than rejected (keeps the DTO forgiving of
a client that always sends both shapes' fields).

### EdgeDTO

| Field | Type | Maps to |
|---|---|---|
| `ID` | `string` | `graph.Edge.ID` |
| `From`, `To` | `string` | `graph.Edge.From`/`To` (must reference declared `NodeDTO.ID`s) |
| `Length`, `Capacity` | `float64` | `graph.Edge.Length`/`Capacity` (descriptive only, per feature 001) |
| `TravelTime` | `TravelTimeDTO` | `graph.Edge.TravelTime` |

**Validation**: `From == To` is rejected (FR-008) before it ever reaches
`graph.AddEdge` (which doesn't itself forbid self-loops, since that's a
web-editing concern, not a graph-structure one).

### DemandDTO

| Field | Type | Maps to |
|---|---|---|
| `Origin`, `Destination` | `string` | `simulation.Demand.Origin`/`Destination` |
| `Size` | `int` | `simulation.Demand.Size` |

### RunRequest

| Field | Type | Notes |
|---|---|---|
| `Nodes` | `[]NodeDTO` | |
| `Edges` | `[]EdgeDTO` | |
| `Demands` | `[]DemandDTO` | |
| `MaxRounds` | `int` | `0` or omitted → server uses `simulation.DefaultMaxRounds` (only an explicit, non-zero value overrides it — unlike `Population.MaxRounds`/`MultiPopulation.MaxRounds`, which are always honored literally, the DTO's zero-value-means-default handling is an HTTP-layer convenience so a typical client never has to think about it). |

### RouteGroupDTO

| Field | Type | Notes |
|---|---|---|
| `EdgeIDs` | `[]string` | The distinct route's edges, in order. |
| `Count` | `int` | How many agents (within this demand) took exactly this route. |
| `TravelTime` | `float64` | This route's travel time at the final assignment. |

### DemandResultDTO

| Field | Type | Maps to |
|---|---|---|
| `Origin`, `Destination`, `Size` | as `DemandDTO` | Echoes the input demand. |
| `RouteGroups` | `[]RouteGroupDTO` | Derived from `simulation.DemandResult.Routes`, grouped by identical edge sequence (research.md decision #4). |
| `TotalTravelTime`, `AverageTravelTime` | `float64` | `simulation.DemandResult.TotalTravelTime`/`AverageTravelTime`. |

### RunResponse

| Field | Type | Maps to |
|---|---|---|
| `Converged` | `bool` | `simulation.MultiAssignmentResult.Converged` |
| `Rounds` | `int` | `simulation.MultiAssignmentResult.Rounds` |
| `TotalTravelTime`, `AverageTravelTime` | `float64` | `simulation.MultiAssignmentResult.TotalTravelTime`/`AverageTravelTime` |
| `PerDemand` | `[]DemandResultDTO` | `simulation.MultiAssignmentResult.PerDemand`, regrouped |

### ErrorResponse

| Field | Type | Notes |
|---|---|---|
| `Error` | `string` | Human-readable message — for a no-route failure, includes the demand-identifying text `simulation.RunDemands` already produces (feature 004, FR-008). |

## Client edit model (`web/src/model/network.ts`)

### ClientNode

| Field | Type | Notes |
|---|---|---|
| `id` | `string` | Client-generated (e.g. a short random id) when the user places a node. |
| `type` | `"house" \| "company" \| "intersection"` | |
| `x`, `y` | `number` | Canvas position — purely a client concern, never sent to or from the server. |

### ClientEdge

| Field | Type | Notes |
|---|---|---|
| `id` | `string` | Client-generated. |
| `from`, `to` | `string` (`ClientNode.id`) | |
| `travelTime` | `{ type: "linear"; freeFlow: number; slope: number } \| { type: "constant"; value: number }` | Chosen (with a sensible default) when the road is drawn; editable afterward. |

### ClientDemand

| Field | Type | Notes |
|---|---|---|
| `origin`, `destination` | `string` (`ClientNode.id`) | |
| `size` | `number` | |

### Edit operations (pure functions over the above, unit-tested — see contracts/api-contract.md's "Client model operations")

- `addNode(state, type, x, y) -> state'`
- `addEdge(state, from, to, travelTime) -> state' | Error("self-loop")`
- `removeNode(state, id) -> state'` (cascades: also removes every edge referencing `id`, and any demand whose origin/destination is `id` — FR-005)
- `removeEdge(state, id) -> state'`

### RunHistory (client-only, not sent to the server)

| Field | Type | Notes |
|---|---|---|
| `previous` | `RunResponse \| null` | The prior run's result, kept for the before/after comparison (FR-007, SC-004). |
| `current` | `RunResponse \| null` | The latest run's result. |

## Relationships

```text
ClientNode.id            ──> ClientEdge.from / .to, ClientDemand.origin / .destination
NodeDTO.ID / EdgeDTO.*    ──> (1:1, request-time mapping from Client* to *DTO)
RunResponse.PerDemand[i]  ──> DemandDTO[i] (same order as sent)
RouteGroupDTO.EdgeIDs[j]  ──> EdgeDTO.ID (resolved back to a ClientEdge for animation)
```

No entity here is persisted beyond the current page load (spec
Assumptions) — `RunHistory` is the only thing carried across a "run"
action within one session, and it is plain in-memory React state.
