# Phase 1 Data Model: Signal-Controlled Queuing

All new types live in the new `queuesim` package. `graph.Graph`/
`graph.Edge`/`graph.TravelTimeFunc` (feature 001) are used as-is;
`agent`/`simulation` are not referenced by this package at all
(research.md decision #3).

## Signal

| Field | Type | Notes |
|---|---|---|
| `ID` | `string` | |
| `EdgeID` | `string` | The `graph.Edge.ID` this signal governs. |
| `GreenDuration` | `float64` | Seconds passable, per cycle. |
| `RedDuration` | `float64` | Seconds not passable, per cycle. |
| `Offset` | `float64` | Seconds into the cycle at simulated `t=0` (phase shift). |

**Derived**: `IsGreenAt(t float64) bool` — see research.md decision #4.
Two `Signal`s are "always opposite" purely by having equal
`GreenDuration`+`RedDuration` and an `Offset` differing by one
`GreenDuration` — not a stored relationship (FR-005).

**Validation rules**:
- `GreenDuration > 0` and `RedDuration > 0` (a signal that's never green
  is representable — see Edge Cases — by setting `RedDuration` very large
  relative to the run's `duration`, not by `GreenDuration == 0`, which
  would make the cycle-length arithmetic degenerate).
- `EdgeID` must reference an edge that exists in the `*graph.Graph` a
  `Run` is called with.

## Demand

| Field | Type | Notes |
|---|---|---|
| `Origin`, `Destination` | `string` (Node ID) | |
| `Count` | `int` | Total agents this demand spawns over the run. `0` is valid (spawns none). |
| `ArrivalInterval` | `float64` | Seconds between one spawn and the next for this demand. |

**Validation rules**:
- `ArrivalInterval > 0` when `Count > 0` (otherwise there's no
  well-defined spawn schedule).

## AgentReport

| Field | Type | Notes |
|---|---|---|
| `DemandIndex` | `int` | Which input `Demand` this agent belongs to. |
| `SpawnTime` | `float64` | Simulated time this agent entered the run. |
| `TravelTime` | `float64` | Total time spent actually moving (sum across its route's edges), excluding waiting. |
| `WaitTime` | `float64` | Total time spent queued at any signal (FR-006). |
| `Arrived` | `bool` | `false` if still in the run (moving or queued) when the run's `duration` elapsed (FR-008, SC-005). |

## QueueSample

| Field | Type | Notes |
|---|---|---|
| `Time` | `float64` | Simulated time this sample was taken. |
| `SignalID` | `string` | Which signal's queue this measures. |
| `Length` | `int` | Number of agents currently queued at that signal at `Time`. |

One `QueueSample` is recorded per signal per tick (FR-004), giving the
full time series a queue-length chart needs.

## RunResult

| Field | Type | Notes |
|---|---|---|
| `QueueSamples` | `[]QueueSample` | Every signal, every tick. |
| `Agents` | `[]AgentReport` | One per spawned agent, in spawn order. |

## Internal (not exported): in-flight agent state

Not part of the public contract — `run.go`'s own bookkeeping during a
`Run` call:

```text
inFlightAgent {
  demandIndex int
  route       []graph.Edge   // computed once at spawn (research.md decision #2)
  edgeIndex   int             // which edge of route it's currently on
  state       moving | queued
  remaining   float64          // seconds left on the current edge, if moving
  queuedAt    float64            // simulated time it joined a queue, if queued
  spawnTime, travelTime, waitTime float64
}
```

Each signal-controlled edge also has an internal FIFO queue of waiting
`inFlightAgent`s and a discharge accumulator (a fractional "credit" of how
many agents the discharge rate allows to leave this tick, carried across
ticks so a non-integer rate doesn't get rounded away — research.md's
tick-based model implies this bookkeeping).

## Relationships

```text
Signal.EdgeID          ──> graph.Edge.ID
Demand.Origin/.Destination ──> graph.Node.ID
AgentReport.DemandIndex ──> index into the []Demand passed to Run
QueueSample.SignalID    ──> Signal.ID
inFlightAgent.route     ──> []graph.Edge (computed via queuesim's own route.go, research.md decision #3)
```

## No state persists between `Run` calls

Same as every prior feature: a `Run` call is a pure function of
`(*graph.Graph, []Signal, []Demand, duration, tick)` returning one
`RunResult` — nothing is kept afterward.
