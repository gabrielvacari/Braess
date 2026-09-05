# Contract: `web/src/model/network.ts` — roads (feature 005 → feature 007)

This replaces the relevant part of feature 005's implicit edge contract.
`cmd/server`'s `POST /api/run` contract (feature 005/006) is unchanged —
see [contracts/api-contract.md](../../005-web-frontend/contracts/api-contract.md).

## Types

```ts
export interface ClientRoad {
  id: string;
  a: string; // Node id — an endpoint, no direction
  b: string; // Node id — the other endpoint
  travelTime: TravelTimeSpec; // shared by both directions
}

export interface NetworkState {
  nodes: ClientNode[];
  roads: ClientRoad[]; // was `edges: ClientEdge[]`
  demands: ClientDemand[];
}
```

## Functions

```ts
// Adds a bidirectional road between two existing nodes. Rejects a
// self-loop (a === b) exactly as addEdge did (FR-008, feature 006).
function addRoad(state: NetworkState, a: string, b: string, travelTime: TravelTimeSpec): NetworkState | NetworkError;

// Removes one road — both its directions disappear together, since a
// road is a single record (FR-005).
function removeRoad(state: NetworkState, id: string): NetworkState;

// Unchanged behavior, reframed: cascades to every road (not "edge") and
// demand referencing the removed node (FR-005, feature 005).
function removeNode(state: NetworkState, id: string): NetworkState;

// Expands every road into its two directed edges, with deterministic
// ids `${road.id}-ab` (a -> b) and `${road.id}-ba` (b -> a), both
// carrying the road's own travelTime (FR-001, FR-002, FR-003).
function expandRoadsToDirectedEdges(roads: ClientRoad[]): DirectedEdge[];
```

## Invariants

- `expandRoadsToDirectedEdges` always returns exactly `2 * roads.length`
  entries — never fewer (a road can't fail to expand) or more.
- For every road, its two derived `DirectedEdge`s have identical
  `travelTime` (deep-equal to `road.travelTime`) and reversed
  `from`/`to` (`ab.from === ba.to === road.a`, `ab.to === ba.from === road.b`).
- Calling `expandRoadsToDirectedEdges` twice with the same `roads` value
  returns the same ids both times (determinism, research.md decision #3)
  — required so a request built from a `NetworkState` and a later
  response interpreted against that same `NetworkState` agree on what
  each edge id means.
- `addRoad`/`removeRoad`/`removeNode` retain every guarantee `addEdge`/
  `removeEdge`/`removeNode` had in feature 005/006, reframed around roads
  instead of edges.

## Compatibility

`cmd/server`, `graph`, `agent`, `simulation`, and `model/api.ts`'s
`RunRequest`/`RunResponse` wire shapes (feature 005) are unchanged — only
what `toRunRequest` builds *from* changes (it now calls
`expandRoadsToDirectedEdges` instead of mapping `edges` 1:1).
