# Phase 1 Data Model: Signal Queuing Frontend

Two related models: the **server's DTOs** (`cmd/server`, Go, mapping
to/from `queuesim` types) and the **client's Signals-mode model**
(`web/src/model/queueNetwork.ts`, TypeScript). See
[contracts/queue-api-contract.md](./contracts/queue-api-contract.md) for
the exact JSON shape connecting them.

## Server DTOs (`cmd/server/queuedto.go`)

### SignalDTO

| Field | Type | Maps to |
|---|---|---|
| `ID` | `string` | `queuesim.Signal.ID` |
| `EdgeID` | `string` | `queuesim.Signal.EdgeID` |
| `GreenDuration` | `float64` | `queuesim.Signal.GreenDuration` |
| `RedDuration` | `float64` | `queuesim.Signal.RedDuration` |
| `Offset` | `float64` | `queuesim.Signal.Offset` (default `0` if omitted) |

### TimeDemandDTO

| Field | Type | Maps to |
|---|---|---|
| `Origin`, `Destination` | `string` | `queuesim.Demand.Origin`/`Destination` |
| `Count` | `int` | `queuesim.Demand.Count` |
| `ArrivalInterval` | `float64` | `queuesim.Demand.ArrivalInterval` |

### QueueRunRequest

| Field | Type | Notes |
|---|---|---|
| `Nodes` | `[]NodeDTO` | Reused from `dto.go` (feature 005) unchanged. |
| `Edges` | `[]EdgeDTO` | Reused from `dto.go` unchanged — directed, one entry per road (no expansion, unlike the equilibrium endpoint — research.md decision #1). |
| `Signals` | `[]SignalDTO` | |
| `Demands` | `[]TimeDemandDTO` | |
| `Duration`, `Tick` | `float64` | Passed straight to `queuesim.Run`. |

### QueueSampleDTO / AgentReportDTO / QueueRunResponse

Mirror `queuesim.QueueSample`/`AgentReport`/`RunResult` field-for-field
(see [contracts/queuesim-api.md](../008-signal-controlled-queuing/contracts/queuesim-api.md)) —
no reshaping needed.

### QueueErrorResponse

| Field | Type | Notes |
|---|---|---|
| `Error` | `string` | Same convention as `dto.go`'s `ErrorResponse`. |

## Client model (`web/src/model/queueNetwork.ts`)

### DirectedRoad

| Field | Type | Notes |
|---|---|---|
| `id` | `string` | |
| `from`, `to` | `string` (`ClientNode.id`, reused from `model/network.ts`) | Directed — order matters (research.md decision #1). |
| `travelTime` | `TravelTimeSpec` (reused type) | |
| `signal` | `{ greenDuration: number; redDuration: number } \| null` | Attached when drawing (FR-001); `null` means an always-open road. |

### TimeDemand

| Field | Type | Notes |
|---|---|---|
| `origin`, `destination` | `string` (`ClientNode.id`) | |
| `count` | `number` | |
| `arrivalInterval` | `number` | Seconds between spawns (FR-002). |

### QueueNetworkState

| Field | Type | Notes |
|---|---|---|
| `nodes` | `ClientNode[]` | Reused type from `model/network.ts` — a node means the same thing in both modes. |
| `roads` | `DirectedRoad[]` | |
| `demands` | `TimeDemand[]` | |

### Edit operations (pure functions, unit-tested)

- `addNode` — reused directly from `model/network.ts` (works on
  `ClientNode[]`, mode-agnostic).
- `addDirectedRoad(state, from, to, travelTime, signal?) -> QueueNetworkState | NetworkError` —
  rejects a self-loop (`from === to`), mirroring `addRoad`'s rule.
- `removeDirectedRoad(state, id) -> QueueNetworkState`
- `removeNodeFromQueueNetwork(state, id) -> QueueNetworkState` — cascades
  to roads and demands referencing the node, mirroring `removeNode`.
- `addTimeDemand` / `removeTimeDemand` — mirror `addDemand`/`removeDemand`.

## Relationships

```text
DirectedRoad.from/.to      ──> ClientNode.id
TimeDemand.origin/.destination ──> ClientNode.id
DirectedRoad.signal        ──> (client-local; expands to a SignalDTO keyed by this road's id, at request-build time)
QueueSampleDTO.SignalID    ──> SignalDTO.ID ──> the DirectedRoad it was attached to
```

## No state persists between runs

Same as every prior feature: `QueueNetworkState` is page-lifetime client
state (no cross-session persistence); the server endpoint is stateless.
