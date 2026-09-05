# Contract: `POST /api/run` (and client model operations)

Unlike features 001-004 (Go library contracts), this feature's primary
contract is an HTTP JSON API between two separate codebases — plus a
small set of pure client-side functions worth contracting explicitly
since they carry FR-005/FR-008's rules.

## `POST /api/run`

**Request body** (`application/json`, see data-model.md for full field
tables):

```json
{
  "nodes": [{"id": "h1", "type": "house"}],
  "edges": [{
    "id": "e1", "from": "h1", "to": "c1", "length": 0, "capacity": 0,
    "travelTime": {"type": "linear", "freeFlow": 0, "slope": 0.01}
  }],
  "demands": [{"origin": "h1", "destination": "c1", "size": 10}],
  "maxRounds": 1000
}
```

**Success response** (`200 OK`):

```json
{
  "converged": true,
  "rounds": 3,
  "totalTravelTime": 130,
  "averageTravelTime": 13,
  "perDemand": [{
    "origin": "h1", "destination": "c1", "size": 10,
    "totalTravelTime": 130, "averageTravelTime": 13,
    "routeGroups": [{"edgeIds": ["e1"], "count": 10, "travelTime": 13}]
  }]
}
```

**Error responses**:

- `400 Bad Request` — malformed JSON, an edge referencing an undeclared
  node, or a self-loop edge (`from == to`, FR-008) — caught before ever
  calling into `graph`/`simulation`.
- `422 Unprocessable Entity` — the request was well-formed but
  `simulation.RunDemands` itself failed: no route for a declared demand
  (wraps `agent.ErrNoRoute`, FR-009) or a travel-time evaluation error.
  Body: `{"error": "<message from RunDemands, already naming the specific demand>"}`.

**Behavior**:

- Builds a fresh `*graph.Graph` from `nodes`/`edges` (via `graph.New()`,
  `AddNode`, `AddEdge`), a `[]simulation.Demand` from `demands`, and calls
  `simulation.RunDemands` with `MaxRounds` (or `simulation.
  DefaultMaxRounds` if `0`/omitted — data-model.md "RunRequest").
- Every request is independent — no state is kept between requests
  (research.md decision #1).
- `RunResponse.PerDemand[i].RouteGroups` groups `simulation.
  DemandResult.Routes` by identical `[]graph.Edge` id sequence, counting
  how many agents share each distinct sequence (research.md decision #4).

## Client model operations (`web/src/model/network.ts`)

These are pure functions (no I/O), unit-tested with Vitest — the
authoritative rules for what the UI is allowed to let the user do,
independent of any rendering concern:

```ts
function addNode(state: NetworkState, type: NodeType, x: number, y: number): NetworkState;

function addEdge(
  state: NetworkState, from: string, to: string, travelTime: TravelTimeSpec
): NetworkState | { error: "self-loop" }; // FR-008: from === to is rejected here, never reaches the server

function removeNode(state: NetworkState, id: string): NetworkState;
// FR-005: also removes every edge with edge.from === id || edge.to === id,
// and every demand with demand.origin === id || demand.destination === id.

function removeEdge(state: NetworkState, id: string): NetworkState;
```

## Compatibility

`graph`, `agent`, and `simulation` (features 001-004) are read-only
dependencies of `cmd/server` — this contract introduces no changes to any
of their exported APIs.
