# Implementation Plan: Signal Queuing Frontend

**Branch**: `009-signal-frontend` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/009-signal-frontend/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Add a second, stateless server endpoint (`POST /api/queue-run`, reusing
`cmd/server/dto.go`'s existing `NodeDTO`/`EdgeDTO`/`buildGraph`) that
translates JSON to/from `queuesim.Run`, and a fully separate "Signals"
mode in the web app — its own small network model, toolbar, canvas
reuse, and results panel — kept deliberately independent of the existing
equilibrium-based mode (FR-006/FR-007) rather than woven into its shared
state. `graph`, `agent`, `simulation`, `queuesim`, and `cmd/graphcli` are
untouched.

## Technical Context

**Language/Version**: Go 1.26 (server) + TypeScript 5.x (web/), both unchanged toolchains.

**Primary Dependencies**: None new. The queue-length-over-time visualization is a small hand-drawn inline SVG line chart, not a charting library — continuing every prior feature's "don't add a dependency the problem doesn't need" discipline, at a scale (a handful of signals, a few dozen samples each) far below where a library would earn its cost.

**Storage**: N/A — the new endpoint is stateless, exactly like the existing one (feature 005).

**Testing**: Go `testing` for the new handler's contract (mirroring `cmd/server/handlers_test.go`); Vitest for the new model's pure functions (edit operations, request/response mapping, and the chart's pure geometry-computation function). Canvas/visual rendering verified manually via quickstart.md, per every prior frontend feature's precedent.

**Target Platform**: Unchanged.

**Project Type**: Extends `cmd/server` (one new endpoint) and `web/` (a new, self-contained "Signals" mode); no other package is touched.

**Performance Goals**: None new — `queuesim.Run` already completes well under a second at this project's demo scale (feature 008).

**Constraints**: No change to `graph`/`agent`/`simulation`/`queuesim`/`cmd/graphcli` (spec Assumptions); the signal-controlled mode must never be visually or functionally confusable with the equilibrium-based mode (FR-007); a road in Signals mode represents one direction only — deliberately not the bidirectional `ClientRoad` feature 007 introduced for the other mode (see research.md decision #1).

**Scale/Scope**: One new server endpoint; one new frontend "mode" with its own model, toolbar, and results panel, reusing only `NetworkCanvas` (rendering) and `ClientNode`/`nodeDisplayLabel` (node data) from the existing code.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | No file under `graph/`, `agent/`, `simulation/`, `queuesim/`, or `cmd/graphcli` is touched | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | `queuesim` was already validated in text (feature 008) before this feature exposes it visually | PASS |
| III. Decentralized, Selfish Agent Behavior | Unaffected — no routing/queuing logic exists in this feature; it only renders `queuesim.Run`'s already-computed output | PASS |
| IV. Validate in Text Before Visually | The new server handler is `go test`-covered; the new frontend model and chart geometry are Vitest-covered; only pixel rendering itself is manual, per established precedent | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions — especially keeping Signals mode fully separate rather than merged, and directed (not bidirectional) roads — are recorded in research.md | PASS |
| Technology Constraints | Still React + Konva.js (canvas reused) and Go; no new dependency | PASS |
| Development Workflow | Branch `009-signal-frontend`; Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/009-signal-frontend/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── queue-api-contract.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
cmd/server/
├── queuedto.go              # NEW: SignalDTO, TimeDemandDTO, QueueRunRequest
│                            # (reuses NodeDTO/EdgeDTO), QueueSampleDTO,
│                            # AgentReportDTO, QueueRunResponse, mapping to/from
│                            # queuesim types
├── queuehandlers.go           # NEW: POST /api/queue-run handler
├── queuehandlers_test.go        # NEW: contract tests
└── main.go                       # UPDATED: registers the new route

web/src/
├── model/
│   ├── queueNetwork.ts          # NEW: DirectedRoad, TimeDemand, QueueNetworkState,
│   │                            # edit operations (reuses ClientNode from network.ts)
│   ├── queueNetwork.test.ts
│   ├── queueApi.ts               # NEW: QueueRunRequest/QueueRunResponse types + fetch wrapper
│   └── queueApi.test.ts
├── components/
│   ├── QueueChart.tsx             # NEW: pure geometry function + small SVG line chart
│   ├── QueueChart.test.ts          # NEW: geometry function tests
│   ├── SignalToolbar.tsx            # NEW: node/road(+signal)/time-demand editing controls
│   └── QueueResultsPanel.tsx         # NEW: charts + arrival/wait summary, labeled "Signal-controlled"
├── SignalApp.tsx                  # NEW: the self-contained Signals-mode screen
└── App.tsx                        # UPDATED: a top-level mode toggle rendering the
                                    # existing equilibrium UI or SignalApp
```

**Structure Decision**: Signals mode is a parallel, self-contained slice
— its own model, toolbar, and results panel — rather than extending the
existing `NetworkState`/`Toolbar`/`ResultsPanel`. `NetworkCanvas` is
reused as-is for rendering (it only needs node positions and `{id, a, b}`
line endpoints, which a directed road's `from`/`to` satisfy without
changing that component at all); everything else that's genuinely
different (directed, signal-bearing roads; time-based demands; a
time-series result) gets its own code. This keeps FR-006/FR-007
(never confusable, never interfering) true by construction rather than
by careful state-sharing discipline inside one shared model.

## Complexity Tracking

*No violations — table intentionally omitted.*
