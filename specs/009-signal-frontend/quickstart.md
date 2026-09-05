# Quickstart: Signal Queuing Frontend

## Prerequisites

- Go 1.26+, Node 20+.
- Features 005-008 already in place.

## Run the automated tests

```sh
go test ./...
cd web && npm test
```

**Expected outcome**: all existing suites pass unmodified; new coverage
includes:
- `cmd/server/queuehandlers_test.go`: a valid `/api/queue-run` request
  returns `200` with queue samples and agent reports; a self-loop or
  malformed request returns `400`; an unroutable demand returns `422`.
- `web/src/model/queueNetwork.test.ts`: `addDirectedRoad` rejects a
  self-loop; `removeNodeFromQueueNetwork` cascades to roads and demands.
- `web/src/model/queueChart.test.ts`: `queueChartPoints` maps
  samples to in-bounds coordinates and preserves their time ordering.

## Run it for real

```sh
# Terminal 1
go run ./cmd/server

# Terminal 2
cd web && npm run dev
```

Open the app, switch to the **Signals** mode (separate from the default
Equilibrium mode — FR-007), and:

1. Place two nodes, draw a directed road between them, and attach a
   signal (e.g. green=5s, red=5s).
2. Add a time-based demand between those nodes with a steady arrival
   rate.
3. Press Run — a queue-length-over-time chart appears for the road's
   signal, plus an arrival/waiting summary.
4. Add a second directed road with an opposite-phase signal between the
   same two points (mirroring feature 008's own two-road experiment) and
   re-run — confirm both charts render and are visibly distinguishable.

## What this proves

- `cmd/server`'s new endpoint is the only new Go code, and it reuses
  `dto.go`'s existing node/edge handling rather than duplicating it
  (Principle I; research.md decision #5).
- Signals mode never shares state with the existing Equilibrium mode —
  switching modes shows a completely different network/toolbar/result,
  never a mix of the two (FR-006, FR-007).
- Every number and chart point came from `queuesim.Run`'s real,
  already-validated output — the browser never computes a queue or a
  route itself (Principle III boundary, consistent with feature 005).

## Out of scope here (see [spec.md](./spec.md#assumptions))

Sharing one drawn network between the two modes, multiple signals per
road, and any change to `queuesim`'s own engine behavior — all explicitly
deferred or ruled out.
