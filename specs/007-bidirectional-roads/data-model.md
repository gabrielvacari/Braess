# Phase 1 Data Model: Bidirectional Roads

Replaces `ClientEdge`/`NetworkState.edges` (feature 005) with
`ClientRoad`/`NetworkState.roads`. `ClientNode`, `ClientDemand`, and every
server-side type (`cmd/server/dto.go`) are unchanged.

## ClientRoad

| Field | Type | Notes |
|---|---|---|
| `id` | `string` | Client-generated. |
| `a`, `b` | `string` (`ClientNode.id`) | The two endpoints — no direction; `a`/`b` order has no meaning (unlike feature 005's `from`/`to`). |
| `travelTime` | `TravelTimeSpec` | Shared by both directions (spec Assumptions, research.md decision #4). |

**Validation rules**:
- `a !== b` (FR-008 of feature 006 continues to apply, reframed: a road
  cannot connect a node to itself, checked in `addRoad`).
- Multiple `ClientRoad`s may share the same `{a, b}` pair (parallel
  roads, FR-006) — not deduplicated.

## NetworkState (updated)

| Field | Type | Notes |
|---|---|---|
| `nodes` | `ClientNode[]` | Unchanged. |
| `roads` | `ClientRoad[]` | Replaces `edges`. |
| `demands` | `ClientDemand[]` | Unchanged. |

## Derived: directed edges (API-boundary only, not stored)

```ts
function expandRoadsToDirectedEdges(roads: ClientRoad[]): DirectedEdge[]

interface DirectedEdge {
  id: string;       // `${road.id}-ab` or `${road.id}-ba` — deterministic (research.md decision #3)
  from: string;      // road.a for "-ab", road.b for "-ba"
  to: string;         // road.b for "-ab", road.a for "-ba"
  travelTime: TravelTimeSpec; // == road.travelTime, both directions
}
```

Every `ClientRoad` expands to exactly two `DirectedEdge`s. This is not a
new stored entity — it exists only where the client talks to the server
(`toRunRequest`) and where it interprets the server's response back into
node positions (`App.tsx`'s path resolver).

## Relationships

```text
ClientRoad.a, ClientRoad.b        ──> ClientNode.id
expandRoadsToDirectedEdges(roads) ──> DirectedEdge[] (2 per road)
DirectedEdge                       ──> EdgeDTO (model/api.ts, sent to POST /api/run — feature 005, unchanged shape)
RunResponse.perDemand[].routeGroups[].edgeIds[] ──> DirectedEdge.id (resolved back to a ClientRoad's a/b for animation)
```

## Edit operations (updated)

- `addRoad(state, a, b, travelTime) -> NetworkState | NetworkError` —
  replaces `addEdge`; same self-loop rejection.
- `removeRoad(state, id) -> NetworkState` — replaces `removeEdge`.
- `removeNode(state, id) -> NetworkState` — cascades to every road with
  `road.a === id || road.b === id` (was `edge.from`/`edge.to`), and to
  every demand referencing the node, unchanged from feature 005.

## No state transitions beyond a single edit

Same as prior features: `NetworkState` is plain, immutable-per-call data;
every operation returns a new value rather than mutating in place.
