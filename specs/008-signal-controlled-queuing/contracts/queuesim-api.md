# Contract: `queuesim` package public API

Like `graph`/`agent`/`simulation`, this is a library — its contract is
the exported Go API. `queuesim` depends on `graph`, the Go standard
library, and nothing else (not `agent`, not `simulation` — research.md
decision #3).

## Types

```go
package queuesim

// Signal governs one road (graph.Edge), alternating between passable
// (green) and not (red). Two Signals are "always opposite" purely by
// sharing GreenDuration+RedDuration and an Offset one GreenDuration
// apart — there is no linking field.
type Signal struct {
    ID            string
    EdgeID        string
    GreenDuration float64
    RedDuration   float64
    Offset        float64
}

// IsGreenAt reports whether s is in its green phase at simulated time t.
func (s Signal) IsGreenAt(t float64) bool

// Demand describes agents spawning over time for one origin/destination
// pair — unlike simulation.Demand (feature 003/004), which has no time
// dimension.
type Demand struct {
    Origin, Destination string
    Count               int
    ArrivalInterval      float64
}

// AgentReport is one spawned agent's outcome.
type AgentReport struct {
    DemandIndex int
    SpawnTime   float64
    TravelTime  float64 // time spent moving
    WaitTime    float64  // time spent queued (FR-006)
    Arrived     bool      // false if still moving/queued when the run ended (FR-008)
}

// QueueSample is one signal's queue length at one point in time.
type QueueSample struct {
    Time     float64
    SignalID string
    Length   int
}

// RunResult is everything a Run call reports.
type RunResult struct {
    QueueSamples []QueueSample
    Agents       []AgentReport
}
```

## Functions

```go
// Run simulates signals and demands over g for the given duration, in
// steps of tick seconds. Agents spawn per their Demand's schedule, each
// computing its own route once at spawn (research.md decision #2) via
// queuesim's own wait-aware pathfinding (not agent's), queuing at any
// signal-controlled edge on its route that is red when it arrives there,
// and discharging at that signal's edge's own free-flow-derived rate
// once green (spec Assumptions: default discharge rate).
//
// Returns an error if a Demand or Signal references a node/edge that
// doesn't exist in g, or if an edge's travel-time evaluation itself
// errors — propagated, not swallowed, matching every prior feature's
// error-handling convention.
func Run(g *graph.Graph, signals []Signal, demands []Demand, duration, tick float64) (RunResult, error)
```

## Error handling

Same convention as every prior feature: a Go `error` as the last return
value, no panics, no swallowed errors.

## Stability

This contract covers this feature only. A browser-facing API for
configuring signals/demands and streaming `RunResult` data (analogous to
`cmd/server`'s `POST /api/run` for `simulation`) is explicitly a separate,
later feature (spec Assumptions, constitution Principle II) — nothing here
is designed against yet, since that shape should be driven by that
feature's own spec.
