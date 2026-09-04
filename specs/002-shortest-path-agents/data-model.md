# Phase 1 Data Model: Shortest-Path Agents

Derived from spec [Key Entities](./spec.md#key-entities) and the Functional
Requirements. Builds on, but does not modify, feature 001's `graph.Node` /
`graph.Edge` / `graph.Graph` / `graph.TravelTimeFunc`.

## Route

| Field | Type | Notes |
|---|---|---|
| `Edges` | `[]graph.Edge` | Ordered from origin to destination. Empty when origin equals destination (FR-006). |
| `TotalTravelTime` | `float64` | Sum of `Edges`' travel times at the volume the route was computed for. `0` when `Edges` is empty. |

**Validation rules**:
- `Edges` must form a connected path: each edge's `From` equals the
  previous edge's `To` (or the origin, for the first edge).
- `TotalTravelTime` is derived, not independently settable — it is always
  the sum of evaluating each edge in `Edges` at the volume the route was
  computed for.

## Agent

| Field | Type | Notes |
|---|---|---|
| `Origin` | `string` (Node ID) | Fixed for this phase (FR-001). |
| `Destination` | `string` (Node ID) | Fixed for this phase (FR-001). |

**Operations** (see [contracts/agent-api.md](./contracts/agent-api.md)):
- `ComputeRoute(g *graph.Graph) (Route, error)` — FR-001, FR-002, FR-003.
  Each call is a fresh, independent computation; `Agent` holds no other
  state and no package-level cache exists, which is what makes FR-003
  ("never a route copied or imposed from outside") true by construction
  rather than by convention.

**No state transitions**: an `Agent`'s `Origin`/`Destination` are fixed for
this phase (per spec Assumptions); `ComputeRoute` does not mutate the
`Agent`, so calling it repeatedly (e.g. from several goroutines, though
concurrency itself is out of scope) always re-derives the same answer from
the same graph rather than drifting.

## Relationships

```text
Agent.Origin, Agent.Destination ──> graph.Node.ID (resolved against a *graph.Graph passed to ComputeRoute)
Route.Edges                     ──> []graph.Edge (elements borrowed from graph.Graph.Edges())
```

An `Agent` does not own a `Graph` — it is handed one each time
`ComputeRoute` is called, keeping `agent` a pure function of
`(graph, origin, destination)` with no hidden state, again in service of
FR-003.
