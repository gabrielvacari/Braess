# Contract: `POST /api/queue-run` (and client model operations)

Sibling to feature 005's `POST /api/run` — see
[005's api-contract.md](../../005-web-frontend/contracts/api-contract.md)
for that endpoint, which this one does not change.

## `POST /api/queue-run`

**Request body**:

```json
{
  "nodes": [{"id": "h1", "type": "house"}, {"id": "c1", "type": "company"}],
  "edges": [{
    "id": "road-a", "from": "h1", "to": "c1", "length": 0, "capacity": 0,
    "travelTime": {"type": "constant", "value": 1}
  }],
  "signals": [{"id": "signal-a", "edgeId": "road-a", "greenDuration": 5, "redDuration": 5, "offset": 5}],
  "demands": [{"origin": "h1", "destination": "c1", "count": 50, "arrivalInterval": 1}],
  "duration": 70,
  "tick": 0.5
}
```

**Success response** (`200 OK`):

```json
{
  "queueSamples": [{"time": 0, "signalId": "signal-a", "length": 0}, "..."],
  "agents": [{"demandIndex": 0, "spawnTime": 0, "travelTime": 1, "waitTime": 0, "arrived": true}, "..."]
}
```

**Error responses**:

- `400 Bad Request` — malformed JSON, a self-loop edge, or an edge/demand
  referencing an undeclared node (same validation `buildGraph` already
  performs for `/api/run`).
- `422 Unprocessable Entity` — `queuesim.Run` itself failed (e.g. a
  demand's destination is unreachable). Body: `{"error": "<message>"}`.

**Behavior**:

- Builds a `*graph.Graph` via the existing `buildGraph(RunRequest)` —
  `QueueRunRequest`'s `Nodes`/`Edges` fields are the same `NodeDTO`/
  `EdgeDTO` types, so no new graph-construction code is needed
  (research.md decision #5).
- Maps `Signals`/`Demands` to `[]queuesim.Signal`/`[]queuesim.Demand` and
  calls `queuesim.Run(g, signals, demands, duration, tick)`.
- Every request is independent — no state kept between requests, exactly
  like `/api/run`.

## Client model operations (`web/src/model/queueNetwork.ts`)

```ts
function addDirectedRoad(
  state: QueueNetworkState, from: string, to: string,
  travelTime: TravelTimeSpec, signal?: { greenDuration: number; redDuration: number } | null
): QueueNetworkState | { error: "self-loop" };

function removeDirectedRoad(state: QueueNetworkState, id: string): QueueNetworkState;

function removeNodeFromQueueNetwork(state: QueueNetworkState, id: string): QueueNetworkState;
// cascades to every road (from/to) and demand (origin/destination) referencing the node

function addTimeDemand(state: QueueNetworkState, origin: string, destination: string, count: number, arrivalInterval: number): QueueNetworkState;
function removeTimeDemand(state: QueueNetworkState, index: number): QueueNetworkState;
```

## `QueueChart`'s pure geometry function

```ts
// Maps a signal's QueueSampleDTO[] (already filtered to one signalId) to
// SVG point coordinates within a width x height viewport — pure, no DOM,
// directly testable.
function queueChartPoints(samples: QueueSampleDTO[], width: number, height: number): { x: number; y: number }[];
```

## Compatibility

`graph`, `agent`, `simulation`, `queuesim`, `cmd/graphcli`, and the
existing `/api/run` endpoint and its DTOs are unchanged — this contract
only adds `queuedto.go`'s new types and one new handler.
