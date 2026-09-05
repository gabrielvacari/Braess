# Implementation Plan: Bidirectional Roads

**Branch**: `007-bidirectional-roads` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/007-bidirectional-roads/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Rework `web/`'s client-side network model around **roads** (bidirectional,
matching a user's mental model and the spec's own vocabulary) instead of
directed **edges** (the engine's model, feature 001). A road is the only
thing a user draws, sees, or removes; it expands into exactly two directed
edges — forward and reverse, sharing one travel-time configuration — at
the one point the client talks to the existing `POST /api/run` API. No
change to `graph`/`agent`/`simulation` or `cmd/server`: they already
accept an arbitrary list of directed edges, so two edges per drawn road is
nothing new to them (FR-007).

## Technical Context

**Language/Version**: TypeScript 5.x (existing `web/` project, unchanged Node/Vite toolchain)

**Primary Dependencies**: None new.

**Storage**: N/A — unchanged from features 005/006.

**Testing**: Vitest for the new `expandRoadsToDirectedEdges` function (produces exactly two directed edges per road, forward and reverse, sharing the road's travel-time) and the renamed `addRoad`/`removeRoad` operations (self-loop rejection, cascade-on-node-removal — same guarantees `addEdge`/`removeEdge` had, now reframed around roads). End-to-end proof that a route is found in either direction after one server round trip is verified manually via quickstart.md, consistent with prior features' precedent for anything requiring the live server.

**Target Platform**: Unchanged — same browser target, same Vite dev server / static build.

**Project Type**: Frontend-only enhancement — no Go changes, no new npm packages.

**Performance Goals**: None new.

**Constraints**: No change to `graph`/`agent`/`simulation`/`cmd/server`/`cmd/graphcli` (spec Assumptions, FR-007); both directions of a road share identical behavior (spec Assumptions); the existing "draw road" interaction (click a node, then a second node) is unchanged — only what it produces changes.

**Scale/Scope**: A rename-and-reshape of the client model (`ClientEdge`/`edges` → `ClientRoad`/`roads`) touching `model/network.ts`, `model/api.ts`, `components/NetworkCanvas.tsx`, and `App.tsx`; `components/Toolbar.tsx` is unaffected (it never referenced edges directly).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | No file under `graph/`, `agent/`, `simulation/`, or `cmd/` is touched (FR-007) | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | Refines the already-permitted visualize step; no reordering | PASS |
| III. Decentralized, Selfish Agent Behavior | Unaffected — this changes what edges *exist* in a request, never how a route is chosen | N/A |
| IV. Validate in Text Before Visually | `expandRoadsToDirectedEdges` and the renamed edit operations are Vitest-tested; the live find-a-route-either-direction proof is manual via quickstart.md, consistent with prior precedent | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions recorded in research.md | PASS |
| Technology Constraints | Still React + Konva.js; no new dependency | PASS |
| Development Workflow | Branch `007-bidirectional-roads`; Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/007-bidirectional-roads/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── client-model-api.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
web/src/
├── model/
│   ├── network.ts              # RENAMED/RESHAPED: ClientEdge/edges -> ClientRoad/roads;
│   │                            # addEdge -> addRoad, removeEdge -> removeRoad (same
│   │                            # self-loop/cascade guarantees); NEW
│   │                            # expandRoadsToDirectedEdges(roads) -> two directed
│   │                            # edges per road, deterministic synthetic ids
│   ├── network.test.ts          # UPDATED to the new names/shapes; NEW tests for
│   │                            # expandRoadsToDirectedEdges's bidirectionality
│   └── api.ts                   # UPDATED: toRunRequest expands roads via
│                                 # expandRoadsToDirectedEdges instead of mapping edges 1:1
├── components/
│   ├── NetworkCanvas.tsx         # UPDATED: `roads` prop, one Line per road (FR-004),
│   │                             # onRoadClick replaces onEdgeClick
│   └── Toolbar.tsx                # unaffected — never referenced edges directly
└── App.tsx                    # UPDATED: exampleNetwork uses roads; "connect" mode
                                # calls addRoad; selection/removal renamed to road;
                                # edgeIdsToPath resolves synthetic directed-edge ids
                                # back to node positions via expandRoadsToDirectedEdges
```

**Structure Decision**: Presentation/model-only change confined to
`web/`; no Go package is touched. Rather than bolting a "grouping id"
onto the existing directed-edge model, the client's *only* editable
network primitive becomes the road — directed edges become a derived,
API-boundary-only concept, generated in exactly one shared function
(research.md decision #2) so the API request builder and the
route-animation path resolver can never disagree about what a road
expands into.

## Complexity Tracking

*No violations — table intentionally omitted.*
