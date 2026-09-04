# Phase 1 Data Model: Graph Simulation Engine

Derived from spec [Key Entities](./spec.md#key-entities) and the Functional
Requirements. This describes the shape of the data, not the full
implementation.

## NodeType

An enum with three values:

| Value | Meaning |
|---|---|
| `Intersection` | A point where roads meet; neither an origin nor a destination |
| `House` | An origin — where agents spawn (used from roadmap Phase 2 onward) |
| `Company` | A destination — where agents are headed (used from roadmap Phase 2 onward) |

## Node

| Field | Type | Notes |
|---|---|---|
| `ID` | `string` | Unique within a `Graph`. Caller-supplied. |
| `Type` | `NodeType` | One of the three values above. |

**Validation rules**:
- `ID` must be non-empty.
- `ID` must be unique within the graph (`AddNode` fails on a duplicate ID —
  the spec doesn't require silently overwriting nodes, and doing so would
  make `Edge` references ambiguous).

## Edge (Road)

| Field | Type | Notes |
|---|---|---|
| `ID` | `string` | Unique within a `Graph`. |
| `From` | `string` (Node ID) | Must reference an existing node (FR-003). |
| `To` | `string` (Node ID) | Must reference an existing node (FR-003). |
| `Length` | `float64` | Descriptive; not used in travel-time computation unless a specific `TravelTimeFunc` references it. |
| `Capacity` | `float64` | Descriptive at this phase; available for travel-time functions that want to normalize volume by capacity. |
| `TravelTime` | `TravelTimeFunc` (`func(volume float64) (float64, error)`) | See [research.md](./research.md#2-travel-time-function-representation). Must return the free-flow time at `volume == 0` and a non-decreasing value as `volume` increases (FR-004); must error on `volume < 0` (FR-008). |

**Validation rules**:
- `ID` must be non-empty and unique within the graph.
- `From` and `To` must both already exist as nodes in the graph at the time
  `AddEdge` is called (FR-003).
- Multiple edges may share the same `(From, To)` pair — parallel/redundant
  routes are a first-class case, not an edge case to special-case away
  (FR-005, spec Edge Cases).
- `Length` and `Capacity` should be positive, but this phase does not
  mandate a specific unit system (none is implied by the spec).

**No state transitions**: nodes and edges are immutable once added in this
phase — there is no "remove edge" or "update capacity" requirement in the
spec. Later phases (interactive frontend, per AGENTS.md roadmap step 5) may
introduce mutation; not needed here.

## Graph

| Field | Type | Notes |
|---|---|---|
| (internal) nodes | `map[string]Node` | Keyed by Node ID. |
| (internal) edges | `map[string]Edge` | Keyed by Edge ID. |

**Operations** (see [contracts/graph-api.md](./contracts/graph-api.md) for
exact signatures):
- `AddNode` — FR-001
- `AddEdge` — FR-002, FR-003, FR-005
- `Nodes` / `Edges` — FR-007 (enumeration for inspection/logging/testing)
- `TravelTime(edgeID string, volume float64)` — FR-004, FR-008

## Relationships

```text
Graph 1 ── * Node
Graph 1 ── * Edge
Edge.From ──> Node.ID
Edge.To   ──> Node.ID
```

An `Edge` does not own its `Node`s; it references them by ID, resolved
against the owning `Graph` at edge-creation time.
