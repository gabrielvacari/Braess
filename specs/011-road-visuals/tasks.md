---

description: "Task list template for feature implementation"
---

# Tasks: Top-Down Road Visuals

**Input**: Design documents from `/specs/011-road-visuals/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md. Depends on `NetworkCanvas.tsx`'s existing `roadColors`/`directed` props (features 009/010), reused unmodified in meaning.

**Tests**: Included — a Vitest unit test for the one new pure geometry function, per this project's established pattern of not unit-testing Konva rendering directly (visual verification is the manual quickstart walkthrough instead).

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Path Conventions

`web/src/model/roadVisuals.ts` (new) and `web/src/components/NetworkCanvas.tsx`
(modified). No other file changes — no engine, server, or other
frontend file is touched.

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: The one piece of new logic (directional-mark placement), needed before User Story 2 can render anything

- [X] T001 [P] Implement `RoadMark` and `roadChevronPositions` in `web/src/model/roadVisuals.ts` (data-model.md)
- [X] T002 [P] Vitest tests in `web/src/model/roadVisuals.test.ts`: marks are evenly spaced strictly between the two nodes' edges, point from a toward b, and no marks are returned for a road too short to fit one (depends on T001)

**Checkpoint**: The placement math is validated before any rendering uses it (constitution Principle IV)

---

## Phase 2: User Story 1 - See roads as roads, not lines (Priority: P1) 🎯 MVP

**Goal**: Every road (either mode) renders as a road-like strip with a dashed center line and rounded ends.

**Independent Test**: Draw a road in either mode and confirm it renders as a strip with a center line, not a thin line, with a smooth transition into each node.

### Implementation for User Story 1

- [X] T003 [US1] In `web/src/components/NetworkCanvas.tsx`, replace each road's plain `Line` with a thicker base strip (`lineCap="round"`, unchanged click/selection/`roadColors` logic) plus a second non-interactive dashed `Line` center marking on top, for both the `directed` and non-`directed` rendering paths (FR-001, FR-002, FR-003, FR-006, FR-007; research.md decisions #1, #2)

**Checkpoint**: Roads in both modes look like roads — User Story 1 is independently demoable.

---

## Phase 3: User Story 2 - See a one-way road's direction as pavement markings (Priority: P2)

**Goal**: A one-way road's direction is shown via repeated marks along its length, not one end arrowhead.

**Independent Test**: Draw a one-way road in Signals mode and confirm its direction is legible from marks distributed along its length; confirm a bidirectional road shows no such marks.

### Implementation for User Story 2

- [X] T004 [US2] When `directed` is true, render one small `Arrow`-based mark per `roadChevronPositions` result instead of the single end-arrowhead (FR-004; depends on T001, T003; research.md decisions #3, #4)
- [X] T005 [US2] Confirm by inspection that the non-`directed` rendering path never invokes `roadChevronPositions` or renders any mark (FR-005)

**Checkpoint**: One-way roads show pavement-style direction marks; bidirectional roads show none — User Stories 1 and 2 both independently pass.

---

## Phase 4: User Story 3 - Keep selection and signal state visible (Priority: P3)

**Goal**: Selection highlighting and the Signals-mode signal-state color still work on the new visual.

**Independent Test**: Select a road and confirm it's still visibly distinguished; run a Signals-mode scenario and confirm a signal-controlled road's green/red state is still clearly visible.

### Implementation for User Story 3

- [X] T006 [US3] Confirm by inspection that `selectedId`/`roadColors` still resolve to the base strip's color only (never the center line or direction marks, which stay a fixed neutral color) — FR-006, FR-007; research.md decision #2
- [X] T007 [US3] Manual walkthrough of quickstart.md steps 3-4: selection and signal green/red remain legible on the new road visual (SC-003)

**Checkpoint**: All three user stories are independently functional.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [X] T008 [P] Run `npx tsc --noEmit` and `npx oxlint` for `web/`; fix any findings
- [X] T009 Run `npm test` in `web/` and confirm every existing suite (App.tsx/SignalApp.tsx flows, all `model/*.test.ts`) still passes unmodified (SC-004)
- [X] T010 [P] Confirm via `git diff --stat` that only `web/src/components/NetworkCanvas.tsx` and `web/src/model/roadVisuals.ts`(+test) changed — no engine, server, or other frontend file (constitution Principle I; FR-008, FR-009)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: BLOCKS User Story 2 (needs `roadChevronPositions`); User Story 1 does not depend on it
- **User Story 1 (Phase 2)**: Can start immediately
- **User Story 2 (Phase 3)**: Depends on Foundational and on User Story 1's strip existing to draw marks onto
- **User Story 3 (Phase 4)**: A preservation check on Stories 1 and 2 — depends on both
- **Polish (Phase 5)**: Depends on all three user stories

### Parallel Opportunities

- `T001` and `T002` (Foundational) can be developed together, `T002` verifying `T001`
- `T008` and `T010` (Polish) in parallel

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Foundational (needed regardless, small)
2. Complete Phase 2: User Story 1
3. **STOP and VALIDATE**: roads look like roads in both modes
4. This alone already delivers the core of the request even before direction marks land

### Incremental Delivery

1. Foundational → placement math validated
2. User Story 1 → validate independently (roads look like roads)
3. User Story 2 → validate independently (one-way direction is legible via marks)
4. User Story 3 → validate independently (nothing prior was lost)
5. Polish → lint/build check + confirm the change stayed confined to the canvas

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- No configurable colors/spacing, and no motion on the road surface
  itself, by design (spec.md Assumptions) — do not add them here
