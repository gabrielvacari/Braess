# Implementation Plan: Frontend Usability and Visual Design

**Branch**: `006-frontend-usability` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/006-frontend-usability/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Rework `web/`'s presentation layer only: a real design-token-based visual
system (CSS custom properties for color/spacing/typography, replacing the
leftover Vite starter theme), plus three new pieces of explanatory
content — an always-visible intro banner (US1), a node-type legend and
per-mode description text (US2), and consistent styling with redundant
text/icons everywhere color currently carries meaning alone (US3). All
explanatory copy and the small amount of logic selecting which message to
show live in one new, pure, Vitest-tested module
(`web/src/content/copy.ts`), separated from the JSX that renders it.
Nothing in `graph`/`agent`/`simulation`/`cmd/server` changes, and none of
the three editing modes' underlying behavior changes — only how legible
they are.

## Technical Context

**Language/Version**: TypeScript 5.x (existing `web/` project, unchanged Node/Vite toolchain)

**Primary Dependencies**: None new. Plain CSS custom properties for the design system — no Tailwind, CSS-in-JS library, or component library — continuing every prior feature's "don't add a dependency the problem doesn't need" discipline, now applied to visual design.

**Storage**: N/A — no new state persists anywhere; the intro banner and legend are always-visible static content, not dismissible (spec Assumptions: no dynamic help system, and this project's frontend has no cross-session persistence per feature 005's own Assumptions).

**Testing**: Vitest for `content/copy.ts`'s pure functions (`modeDescription`, `emptyCanvasHint`, `missingDemandHint`) — the only genuinely testable logic this feature adds. Visual/CSS outcomes (consistent styling, legend/banner actually rendering, grayscale color-redundancy check) are verified manually via quickstart.md's walkthrough of the three Acceptance Scenarios and SC-004/SC-005, consistent with feature 005's own precedent that canvas/visual output isn't practical to unit-test.

**Target Platform**: Same as feature 005 — any modern desktop browser via the existing Vite dev server / static build.

**Project Type**: Frontend-only enhancement — no Go changes, no new npm packages.

**Performance Goals**: None new; this is a presentation change at the same scale as feature 005.

**Constraints**: No change to `graph`/`agent`/`simulation`/`cmd/server` (spec Assumptions); no new UI framework or component library (spec Assumptions, constitution Technology Constraints already name React + Konva.js); the three editing modes' behavior stays exactly as feature 005 built it — only their legibility changes; explanatory content must stay usable at a small viewport (spec Edge Cases).

**Scale/Scope**: One new content module, two new small components (`IntroBanner`, `Legend`), edits to three existing components (`Toolbar`, `NetworkCanvas`, `ResultsPanel`) and `App.tsx`, and a full rewrite of `index.css`.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | No file under `graph/`, `agent/`, `simulation/`, or `cmd/server` is touched | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | Refines the already-permitted visualize step (Phase 5); doesn't skip or reorder anything | PASS |
| III. Decentralized, Selfish Agent Behavior | Unaffected — no routing/congestion logic anywhere in this feature | N/A |
| IV. Validate in Text Before Visually | New pure copy-selection logic is Vitest-tested; visual/styling outcomes verified manually via quickstart.md, per feature 005's own established precedent for canvas/visual work | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions recorded in research.md | PASS |
| Technology Constraints | Still React + Konva.js; no new dependency added | PASS |
| Development Workflow | Branch `006-frontend-usability`; Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/006-frontend-usability/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── copy-api.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
web/src/
├── content/                  # NEW
│   ├── copy.ts                 # explanatory strings + modeDescription/emptyCanvasHint/missingDemandHint
│   └── copy.test.ts             # Vitest tests for the three pure functions
├── components/
│   ├── IntroBanner.tsx           # NEW — US1: what the tool demonstrates, always visible
│   ├── Legend.tsx                 # NEW — US2: node type marker -> name/meaning
│   ├── Toolbar.tsx                 # EXTENDED — renders modeDescription(mode) (US2); restyled (US3)
│   ├── NetworkCanvas.tsx            # EXTENDED — empty-state overlay via emptyCanvasHint (US1); every
│   │                                # node keeps/gains a text label so type isn't color-only (US3/FR-007);
│   │                                # wrapped in a styled, bordered container
│   └── ResultsPanel.tsx              # EXTENDED — restyled; better/worse already has text, kept explicit (US3/FR-007)
├── App.tsx                    # EXTENDED — layout restructure (header/sidebar/main), wires IntroBanner,
│                                # Legend, and missingDemandHint near the Run button (US1 Edge Case)
└── index.css                    # REWRITTEN — design tokens (color/spacing/typography) replacing the
                                  # leftover Vite starter theme; base styles for buttons/inputs/selects/
                                  # fieldsets so no unstyled control remains (US3, SC-004)
```

**Structure Decision**: Presentation-only change confined to `web/`;
`cmd/server` and every Go package are untouched. A new `content/` folder
holds copy and its pure selection logic separately from the components
that render it — the same "keep pure logic separate from rendering"
split `web/src/model/` already established in feature 005 (`network.ts`'s
edit operations vs. the components that call them), applied here to
explanatory text instead of network data.

## Complexity Tracking

*No violations — table intentionally omitted.*
