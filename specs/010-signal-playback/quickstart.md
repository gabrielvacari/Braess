# Quickstart: Signal-Controlled Live Playback

## Prerequisites

- Go 1.26+, Node 20+.
- Features 005-009 already in place.

## Run the automated tests

```sh
go test ./...
cd web && npm test
```

**Expected outcome**: all existing suites pass unmodified; new coverage includes:
- `queuesim/*_test.go`: position samples are recorded for every in-flight
  agent every tick, with `Progress` in `[0, 1]`, `Queued` matching the
  agent's actual state, and a final `Progress: 1` sample on the tick an
  agent arrives; `QueueSample.Green` matches `Signal.IsGreenAt(Time)`.
- `cmd/server/queuehandlers_test.go`: a valid `/api/queue-run` request's
  response includes non-empty `positions` and `queueSamples[].green`.
- `web/src/model/queuePlayback.test.ts`: `positionAtSimTime` interpolates
  along an edge while moving, holds at the edge's start while queued, and
  returns `null` before an agent has spawned; `signalGreenAtSimTime`
  matches the nearest recorded sample.

## Run it for real

```sh
# Terminal 1
go run ./cmd/server

# Terminal 2
cd web && npm run dev
```

Open the app, switch to **Signals** mode, and:

1. Draw the classic two-road, opposite-signal scenario (or load whatever
   you already had), and press Run.
2. Alongside the existing queue chart and arrival/wait summary, a
   playback now appears: agents move across the canvas, and each
   signal-controlled road is colored green or red matching its current
   simulated state.
3. Watch an agent stop at a red road and resume the moment it turns
   green — confirm the pause is visible, not just implied by the chart.
4. Press "restart" — confirm the same run replays identically without a
   new request to the server.

## What this proves

- `queuesim.Run`'s existing output (`QueueSamples`, `Agents`) is
  unchanged in meaning — only extended (Principle I/II compliance,
  confirmed by `go test ./...` passing unmodified for every pre-existing
  assertion).
- The animation is a pure readout of already-computed engine data — the
  browser never decides an agent's route or a signal's phase itself
  (Principle III boundary).
- The equilibrium mode's own animation and network are untouched
  (`git diff --stat` shows no changes under `animation/useAgentAnimation.ts`,
  `animation/AgentTicker.tsx`, or the `EquilibriumApp` portion of `App.tsx`).

## Out of scope here (see [spec.md](./spec.md#assumptions))

A scrubber/seek control beyond play/pause/restart, and any change to how
many signals a road may have — both explicitly deferred.
