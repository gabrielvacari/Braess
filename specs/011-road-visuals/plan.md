# Implementation Plan: Top-Down Road Visuals

**Branch**: `011-road-visuals` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/011-road-visuals/spec.md`

## Summary

`NetworkCanvas` currently renders every road as a single thin `Line`
(plain lines in Equilibrium mode, `Arrow`-with-one-arrowhead in Signals
mode, added by feature 010). This feature replaces that with a
road-like visual shared by both modes: a thick, rounded-cap strip with a
dashed center line, plus — for one-way roads only — small directional
arrow marks repeated along the strip instead of one end arrowhead.
Selection highlighting and the Signals-mode green/red override continue
to color the strip itself. Purely a rendering change to one component;
no engine, API, or edit-behavior change.

## Technical Context

**Language/Version**: TypeScript 5 / React 19, `web/` only

**Primary Dependencies**: `react-konva` (already used — `Line`, `Arrow`; no new library)

**Storage**: N/A

**Testing**: Vitest, for the one new pure geometry function this introduces

**Target Platform**: Browser (Konva canvas), both existing app modes

**Project Type**: Web application (existing `web/src/components/NetworkCanvas.tsx`)

**Performance Goals**: No new goal — same node/road counts feature 009/010 already render smoothly (tens of roads, hundreds of agents)

**Constraints**: Zero change to `graph`/`agent`/`simulation`/`queuesim`/`cmd/server`; zero change to how roads are drawn, selected, or removed (spec FR-008/FR-009); the `roadColors` and `directed` props `NetworkCanvas` already exposes are extended, not replaced, so `App.tsx`/`SignalApp.tsx`'s existing calls keep working

**Scale/Scope**: One component file, one new small pure-function module

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Status |
|---|---|---|
| I. Engine Independent of UI | Change is entirely inside `web/src/components/` — no engine package touched | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | No new simulation capability is introduced; this only changes how already-validated data (roads, selection, signal state) is drawn | PASS (N/A — no new capability to simulate first) |
| III. Decentralized, Selfish Agent Behavior | Not applicable — no agent/routing logic touched | PASS |
| IV. Validate in Text Before Validating Visually | The one new piece of logic (chevron placement math) is a pure function, validated by a Vitest unit test, independent of any rendering | PASS |
| V. Explain Reasoning, Not Just Code | research.md documents why each visual choice (Line vs. custom shape, chevron-via-Arrow, etc.) was made | PASS |

No violations — Complexity Tracking not needed.

## Project Structure

### Documentation (this feature)

```text
specs/011-road-visuals/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md         # Phase 1 output (the one new pure function)
└── quickstart.md
```

No `contracts/` — this feature adds no API, server, or client-model
contract; it only changes what `NetworkCanvas` draws from data it
already receives.

### Source Code (repository root)

```text
web/src/
├── model/
│   ├── roadVisuals.ts       # NEW: pure function — roadChevronPositions
│   └── roadVisuals.test.ts
└── components/
    └── NetworkCanvas.tsx    # MODIFIED: road rendering only
```

`App.tsx`, `SignalApp.tsx`, `SignalToolbar.tsx`, `Toolbar.tsx`, and every
non-canvas component are untouched — both callers already pass
`roadColors`/`directed` as needed; neither call site changes.

**Structure Decision**: Extends the existing shared `NetworkCanvas`
(feature 009 research.md decision #3) with a new pure-geometry sibling
module in `model/`, matching the project's established split (pure logic
in `model/`, rendering in `components/`) already used by
`queueChart.ts`/`QueueChart.tsx` and `queuePlayback.ts`.

## Complexity Tracking

*No constitution violations — table not needed.*
