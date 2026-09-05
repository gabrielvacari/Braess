# Implementation Plan: Signal-Controlled Live Playback

**Branch**: `010-signal-playback` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/010-signal-playback/spec.md`

## Summary

Today, `queuesim.Run` returns only aggregate results (per-agent totals,
per-signal queue lengths over time) — enough for a chart, not enough to
replay what actually happened. This feature adds two purely additive,
per-tick recordings to `queuesim.Run`'s output — each agent's location
and moving/queued state, and each signal's green/red phase — validated
in Go tests first (Principle II), then a Signals-mode frontend playback
that reads those recordings back and renders real agent movement plus
live road-signal coloring, alongside the existing queue chart and
arrival/wait summary (not replacing them).

## Technical Context

**Language/Version**: Go 1.26 (engine, `cmd/server`); TypeScript 5 / React 19 (`web/`)

**Primary Dependencies**: `braess/graph`, `braess/queuesim` (engine — no new dependency); `react-konva` (already used by `NetworkCanvas`) — no new frontend library

**Storage**: N/A — stateless, in-memory per request/response, matching every prior feature

**Testing**: Go `testing` (table-driven, `queuesim` + `cmd/server`); Vitest (pure functions in `web/src/model`)

**Target Platform**: Same as existing project — Go binary + static-served React SPA, browser client

**Project Type**: Web application (existing `cmd/server` + `web/` split, features 005-009)

**Performance Goals**: Smooth (60fps-capable) playback for the hand-drawn scenario scale already established (tens to a few hundred agents, run durations of tens to low hundreds of simulated seconds) — no new performance goal beyond what feature 009's canvas already handles

**Constraints**: Zero change to `queuesim.Run`'s existing fields/behavior (`QueueSamples`, `Agents` summary semantics) or to the equilibrium mode's existing animation — this feature only adds new fields/output, per spec FR-007/FR-009

**Scale/Scope**: Same scale as feature 008/009 (the classic two-road demo: ~150 agents, 70s duration, 0.25s tick → ~280 ticks × up to a few hundred in-flight agents per tick)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Status |
|---|---|---|
| I. Engine Independent of UI | New recordings live in `queuesim` (Go, no UI/rendering dependency); `cmd/server` DTOs translate them, same as every prior feature | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | The new per-tick recordings are a `queuesim` capability, validated via Go tests (asserting positions/phases are recorded correctly) *before* any frontend playback code is written — same order as feature 008 → 009 | PASS |
| III. Decentralized, Selfish Agent Behavior | No routing/decision logic changes — agents still choose routes exactly as `queuesim.Run` already decided; this feature only records what already happened, it decides nothing new | PASS |
| IV. Validate in Text Before Validating Visually | New recording behavior is exercised via Go tests (`go test ./queuesim/...`) before the frontend playback exists | PASS |
| V. Explain Reasoning, Not Just Code | research.md documents the reasoning for each design choice (why extend `QueueSample` rather than add a new type, why a single shared clock rather than per-agent tickers, etc.) | PASS |

No violations — Complexity Tracking not needed.

## Project Structure

### Documentation (this feature)

```text
specs/010-signal-playback/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md         # Phase 1 output
├── contracts/
│   └── playback-contract.md
├── quickstart.md
└── tasks.md              # /speckit-tasks output (not this command)
```

### Source Code (repository root)

```text
queuesim/
├── report.go             # MODIFIED: +PositionSample type, +RunResult.Positions,
│                          #   +QueueSample.Green, +AgentReport.ID
├── run.go                # MODIFIED: record positions/phases each tick (additive)
└── run_test.go, two_road_test.go  # MODIFIED/EXTENDED: assert the new recordings

cmd/server/
├── queuedto.go            # MODIFIED: +PositionSampleDTO, +QueueRunResponse.Positions,
│                          #   +QueueSampleDTO.Green, +AgentReportDTO.ID
└── queuehandlers_test.go  # MODIFIED: assert positions/green come through the API

web/src/
├── model/
│   ├── queueApi.ts         # MODIFIED: mirror the new DTO fields/types
│   └── queuePlayback.ts    # NEW: pure functions — positionAtSimTime, signalGreenAtSimTime
│   └── queuePlayback.test.ts
├── animation/
│   └── usePlaybackClock.ts # NEW: single shared simulated-time clock (play/pause/restart)
├── components/
│   ├── NetworkCanvas.tsx   # MODIFIED: +optional roadColors prop (backward compatible)
│   └── PlaybackControls.tsx # NEW: play/pause/restart + "t = Xs / Ys" readout
└── SignalApp.tsx           # MODIFIED: wires the clock, positions, and road colors into
                             #   NetworkCanvas, alongside the existing QueueResultsPanel
```

`graph/`, `agent/`, `simulation/`, `cmd/graphcli/`, and every equilibrium-mode file
(`App.tsx`'s `EquilibriumApp`, `Toolbar.tsx`, `ResultsPanel.tsx`, `animation/useAgentAnimation.ts`,
`animation/AgentTicker.tsx`) are untouched.

**Structure Decision**: Extends the existing `queuesim` + `cmd/server` + `web/src` layout
from features 008/009 — no new top-level directory. New frontend logic lives alongside its
siblings by kind (`model/` for pure data functions, `animation/` for the playback clock,
`components/` for rendering), matching the project's established file organization.

## Complexity Tracking

*No constitution violations — table not needed.*
