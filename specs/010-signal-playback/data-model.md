# Phase 1 Data Model: Signal-Controlled Live Playback

All additions are purely additive to the shapes feature 008/009 already
established — nothing below removes or renames an existing field.

## Engine (`queuesim`)

### `PositionSample` (new)

| Field | Type | Notes |
|---|---|---|
| `Time` | `float64` | The simulated tick this sample was recorded at. |
| `AgentID` | `int` | Matches the `AgentReport.ID` this sample belongs to (research.md decision #2). |
| `EdgeID` | `string` | The road the agent occupies at this tick. |
| `Progress` | `float64` | `0..1` fraction along the edge; always `0` while `Queued`. |
| `Queued` | `bool` | `true` if the agent is waiting at this edge's signal rather than moving along it. |

### `QueueSample` (extended)

Adds one field to the existing struct:

| Field | Type | Notes |
|---|---|---|
| `Green` | `bool` | **New.** Whether `SignalID`'s signal is green at `Time` (`Signal.IsGreenAt(Time)`), FR-002. |

### `AgentReport` (extended)

Adds one field to the existing struct:

| Field | Type | Notes |
|---|---|---|
| `ID` | `int` | **New.** Spawn-order identifier, 0-based, unique per spawned agent (research.md decision #2). |

### `RunResult` (extended)

Adds one field to the existing struct:

| Field | Type | Notes |
|---|---|---|
| `Positions` | `[]PositionSample` | **New.** One entry per in-flight agent per tick, in the order recorded (FR-001). |

## Server (`cmd/server`)

Mirrors the engine types field-for-field, exactly like `queuedto.go`
already mirrors `queuesim.Signal`/`queuesim.Demand`.

### `PositionSampleDTO` (new)

```go
type PositionSampleDTO struct {
	Time     float64 `json:"time"`
	AgentID  int     `json:"agentId"`
	EdgeID   string  `json:"edgeId"`
	Progress float64 `json:"progress"`
	Queued   bool    `json:"queued"`
}
```

### `QueueSampleDTO` (extended)

Adds `Green bool `json:"green"`` .

### `AgentReportDTO` (extended)

Adds `ID int `json:"id"`` .

### `QueueRunResponse` (extended)

Adds `Positions []PositionSampleDTO `json:"positions"`` .

## Client (`web/src/model`)

### `queueApi.ts` (extended)

Mirrors the server DTOs exactly — `PositionSampleDTO`, `QueueSampleDTO.green`,
`AgentReportDTO.id`, `QueueRunResponse.positions` — same convention as
every prior DTO mirror in this file.

### `queuePlayback.ts` (new, pure functions)

| Function | Signature | Behavior |
|---|---|---|
| `groupPositionsByAgent` | `(positions: PositionSampleDTO[]) => Map<number, PositionSampleDTO[]>` | Groups a run's flat positions by `agentId`, each group kept in recorded (time) order. |
| `positionAtSimTime` | `(samples: PositionSampleDTO[], edgesById: Map<string, {from: string; to: string}>, nodesById: Map<string, {x: number; y: number}>, time: number) => {x: number; y: number} \| null` | Finds the last sample with `time <= arg`, places the agent at its edge's `from` node if `queued`, or interpolates `from → to` by `progress` if moving. Returns `null` before the agent's first sample (not yet spawned). |
| `signalGreenAtSimTime` | `(samples: QueueSampleDTO[], signalId: string, time: number) => boolean` | Finds the last sample for `signalId` with `time <= arg` (falling back to the earliest sample for that signal if `time` precedes all of them) and returns its `green` value. |

## Notes

- `PositionSample`/`QueueSample.Green` describe *what already happened*
  during `Run` — no new decision-making, consistent with Principle III
  (nothing here changes how an agent chooses a route).
- Every existing field on `QueueSample`, `AgentReport`, and `RunResult`
  keeps its exact current meaning; this feature only adds to each.
