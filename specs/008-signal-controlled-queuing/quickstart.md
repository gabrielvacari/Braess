# Quickstart: Signal-Controlled Queuing

Validates the feature end-to-end without any UI, per constitution
Principle IV — this feature is explicitly engine-only (spec Assumptions).

## Prerequisites

- Go 1.26+ installed.
- Features 001-002 (`graph`, `agent`) already in place. `simulation`,
  `cmd/server`, and `web/` are unaffected and not required for this
  feature's own validation.

## Run the automated tests

```sh
go test ./queuesim/...
```

**Expected outcome**: all tests pass, including:
- `Signal.IsGreenAt`: correct green/red at phase boundaries, and two
  signals configured with a one-`GreenDuration` offset are never green at
  the same time (SC test for research.md decision #4).
- Queue discharge accounting: a queue drains at the configured rate once
  green, and any remainder correctly carries into the next red/green
  cycle rather than being dropped (spec Edge Cases).
- The classic two-road, opposite-phase scenario (`two_road_test.go`):
  produces a `QueueSample` time series for each road and an `AgentReport`
  per spawned agent, with total wait time reported separately from travel
  time (SC-003).
- No `AgentReport` shows `Arrived: true` for an agent that crossed a
  signal-controlled edge while it was red (SC-002) — checked by
  reconstructing the arrival time at that edge and comparing against
  `Signal.IsGreenAt`.
- An agent still queued/moving when the run's `duration` elapses is
  reported with `Arrived: false` (SC-005).

## Run the standalone CLI

```sh
go run ./cmd/graphcli
```

**Expected outcome**: in addition to earlier phases' output, the command
now runs the classic two-road, opposite-phase signal scenario and prints:
- Each road's queue length at intervals across the run (enough to see the
  shape of the curve — rising during red, falling during green).
- Total agents spawned, and how many arrived vs. were still
  waiting/moving at the end.
- Each road's total and average wait time across the agents that used it.

Completes in well under 5 seconds and exits with status `0`.

## What this proves

- `queuesim` imports only `graph` and the Go standard library — no UI, no
  `agent`, no `simulation` dependency (Principle I; research.md decision
  #3).
- Every printed number came from the tick-based simulation actually
  running — no agent ever crosses a red signal, and every queue's rise
  and fall is driven by real spawn/discharge accounting, not a scripted
  outcome (Principle IV).
- Whatever pattern the two-road experiment's printed report shows —
  persistent asymmetry, oscillation, or balance — is the actual simulated
  result, not asserted in advance (spec Assumptions, matching the
  project's experimental treatment of congestion since Phase 3).

## Out of scope here (see [spec.md](./spec.md#assumptions))

Any browser UI for configuring signals or visualizing queues, mid-trip
route replanning, intersections with more than two movements, and amber
phases — all explicitly deferred or ruled out.
