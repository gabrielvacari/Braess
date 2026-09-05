---

description: "Task list template for feature implementation"
---

# Tasks: Bidirectional Roads

**Input**: Design documents from `/specs/007-bidirectional-roads/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/client-model-api.md, quickstart.md. Depends on features 005-006 (`web/`), which this feature renames/reshapes without changing the Go side.

**Tests**: Included for `network.ts`'s renamed edit operations and the new `expandRoadsToDirectedEdges` function. The live find-a-route-either-direction proof needs the actual server, so it's verified manually via quickstart.md, consistent with prior precedent.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

All changes are under `web/src/`. No Go file is touched.

---

## Phase 1: Setup

- [X] T001 Run `cd web && npm test` to confirm the pre-refactor baseline (feature 006's suite) passes — the regression floor this feature's rename/reshape must not break

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The road-based model and its expansion into directed edges —
everything else builds on this

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 In `web/src/model/network.ts`: rename `ClientEdge` → `ClientRoad` (fields `a`/`b` replacing `from`/`to`), `NetworkState.edges` → `roads`, `addEdge` → `addRoad`, `removeEdge` → `removeRoad` (same self-loop rejection); update `removeNode`'s cascade to check `road.a`/`road.b` (data-model.md)
- [X] T003 Implement `expandRoadsToDirectedEdges(roads: ClientRoad[])` in `web/src/model/network.ts`: two directed edges per road (`${road.id}-ab`, `${road.id}-ba`), both carrying `road.travelTime` (data-model.md; contracts/client-model-api.md; depends on T002)
- [X] T004 Update `web/src/model/api.ts`'s `toRunRequest` to build its `edges` array via `expandRoadsToDirectedEdges(state.roads)` instead of mapping `state.edges` 1:1 (depends on T003)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Drawing one road connects both directions (Priority: P1) 🎯 MVP

**Goal**: One "draw road" action lets agents be routed in either
direction between the two nodes.

**Independent Test**: Draw one road between two nodes, declare demands in
each direction, and confirm both find a route.

### Implementation for User Story 1

- [X] T005 [US1] Update `App.tsx`'s "connect" mode handler to call `addRoad` instead of `addEdge` (depends on T002)
- [X] T006 [US1] Update `App.tsx`'s `exampleNetwork()` to build `roads: ClientRoad[]` instead of `edges: ClientEdge[]` (depends on T002)
- [X] T007 [US1] Update `App.tsx`'s `edgeIdsToPath` to resolve a `RouteGroupDTO`'s edge ids back to node positions via `expandRoadsToDirectedEdges(network.roads)` instead of `network.edges` (FR-001, FR-002; depends on T003)

### Tests for User Story 1

- [X] T008 [P] [US1] Vitest tests for `expandRoadsToDirectedEdges` in `network.test.ts`: exactly two directed edges per road, correct `from`/`to` reversal, shared `travelTime`, deterministic ids across repeated calls on the same input (contracts/client-model-api.md invariants; depends on T003)
- [X] T009 [P] [US1] Update `network.test.ts`'s existing `addEdge`/`removeEdge` tests to `addRoad`/`removeRoad` (rename, same assertions; depends on T002)

**Checkpoint**: `npm test` passes; a road drawn in one click-pair routes agents in either direction — User Story 1 is independently functional and demoable.

---

## Phase 4: User Story 2 - Removing a road removes both directions together (Priority: P2)

**Goal**: Removing a drawn road removes connectivity in both directions
in one action.

**Independent Test**: Draw a road, remove it, and confirm neither
direction between those two nodes routes successfully anymore (absent
another road).

### Implementation for User Story 2

- [X] T010 [US2] Update `App.tsx`'s road-removal wiring: rename `handleEdgeClick` → `handleRoadClick`, call `removeRoad`, and rename the selection kind from "edge" to "road" (depends on T002)
- [X] T011 [US2] Update `NetworkCanvas.tsx`: `edges` prop → `roads`, render one `<Line>` per road using `road.a`/`road.b` node positions (FR-004), rename `onEdgeClick` → `onRoadClick` (depends on T002)

### Tests for User Story 2

- [X] T012 [P] [US2] Update/rename the `removeNode` cascade test in `network.test.ts` to assert against `roads` (depends on T002)

### Validation

- [X] T013 Walk through quickstart.md's three manual scenarios by hand: draw + reverse-direction demand succeeds; remove + same demand now fails clearly; a second parallel road between the same two nodes is still accepted and independently bidirectional (SC-001 through SC-004; depends on T005-T011)

**Checkpoint**: All user stories are independently functional — a drawn road behaves as one bidirectional thing to draw, use, and remove, closing out the confusion feature 006 surfaced.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [X] T014 [P] Run `npx tsc --noEmit` and `npx oxlint` in `web/`; fix any findings
- [X] T015 Run the full quickstart.md validation end to end: `npm test`, `npm run build`, then the manual walkthrough (depends on T008, T009, T012, T013)
- [X] T016 [P] Confirm via `git diff --stat` that no file under `graph/`, `agent/`, `simulation/`, or any `cmd/` package changed (FR-007, constitution Principle I)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS both user stories (the model rename/reshape everything else is built on)
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on Foundational (T002) directly; independent of User Story 1's specific tasks (different concerns: drawing/animating vs. selecting/removing), though both touch `App.tsx`
- **Polish (Phase 5)**: Depends on both user stories being complete

### Within Each User Story

- Foundational (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green, or — for the live-route
  proof — manually verified) before moving to the next

### Parallel Opportunities

- `T008` and `T009` (US1 tests, same file but independent test blocks) can be written together
- `T010` and `T011` (US2, different files — `App.tsx` vs. `NetworkCanvas.tsx`) in parallel
- `T014` and `T016` (Polish) in parallel

---

## Parallel Example: User Story 2

```bash
Task: "Update App.tsx's road-removal wiring to call removeRoad"
Task: "Update NetworkCanvas.tsx to accept roads and render one Line per road"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (confirm the pre-refactor baseline passes)
2. Complete Phase 2: Foundational (the road model + expansion function)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: draw a road, declare a demand in the direction you *didn't* draw it, confirm it routes
5. This alone fixes the exact confusing error that motivated this feature

### Incremental Delivery

1. Setup + Foundational → road-based model ready
2. User Story 1 → validate independently (bidirectional routing works)
3. User Story 2 → validate independently (bidirectional removal works)
4. Polish → lint/build check + confirm zero Go-side changes

---

## Notes

- [P] tasks touch different files (or independent blocks within one file) with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- Per-direction road configuration and any change to the draw-road
  interaction itself are explicitly out of scope (spec.md Assumptions) —
  do not add them here
