---

description: "Task list template for feature implementation"
---

# Tasks: Frontend Usability and Visual Design

**Input**: Design documents from `/specs/006-frontend-usability/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/copy-api.md, quickstart.md. Depends on feature 005 (`web/`), which this feature extends without changing its editing-mode behavior. No Go package is touched.

**Tests**: Included for `content/copy.ts`'s pure functions, the only genuinely testable logic this feature adds (plan.md Testing strategy). Visual/styling outcomes are verified manually via quickstart.md, consistent with feature 005's own precedent that canvas/visual work isn't practical to unit-test.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

All changes are under `web/src/`. No Go file is touched.

---

## Phase 1: Setup

- [X] T001 [P] Create the `web/src/content/` directory

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The copy module and design tokens every user story builds on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Define `APP_INTRO`, the `LegendEntry` type, and `NODE_TYPE_LEGEND` in `web/src/content/copy.ts` (data-model.md)
- [X] T003 [P] Implement `modeDescription`, `emptyCanvasHint`, and `missingDemandHint` in `web/src/content/copy.ts` (data-model.md; contracts/copy-api.md)
- [X] T004 [P] Rewrite `web/src/index.css`: design tokens (color, spacing, typography as CSS custom properties) replacing the leftover Vite starter theme, plus base styles for buttons, inputs, selects, and fieldsets (FR-005, FR-006)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Understand the tool and what to do first (Priority: P1) 🎯 MVP

**Goal**: A first-time visitor, before drawing anything, understands what
the tool demonstrates and what to do first.

**Independent Test**: Open the app with an empty network and confirm a
visible explanation of the app's purpose and a clear first action are
present, with no interaction needed.

### Implementation for User Story 1

- [X] T005 [US1] Implement `IntroBanner.tsx` rendering `APP_INTRO` (FR-001; depends on T002)
- [X] T006 [US1] Wire `emptyCanvasHint` into the canvas area (`NetworkCanvas.tsx` or an overlay in `App.tsx`) so it suggests a first action when `network.nodes.length === 0` (FR-004; depends on T003)
- [X] T007 [US1] Wire `missingDemandHint` near the "Run" button in `App.tsx`, and render `handleRun`'s caught errors in plain language near Run rather than as a raw message (FR-008, FR-009; depends on T003)
- [X] T008 [US1] Wire `IntroBanner` into `App.tsx`'s layout, always visible (depends on T005)

### Tests for User Story 1

- [X] T009 [P] [US1] Vitest tests for `emptyCanvasHint` and `missingDemandHint` in `web/src/content/copy.test.ts`: a message at count `0`, `null` at count `1` (depends on T003)

**Checkpoint**: Loading the app fresh shows a visible explanation and a clear first action — User Story 1 is independently functional and demoable.

---

## Phase 4: User Story 2 - Understand what's on screen without guessing (Priority: P2)

**Goal**: Every node marker is explained by a legend, and the active
editing mode's effect is described before clicking.

**Independent Test**: Load a network and confirm a legend maps every
marker to its meaning; select each editing mode in turn and confirm a
description of its effect is visible.

### Implementation for User Story 2

- [X] T010 [US2] Implement `Legend.tsx` rendering `NODE_TYPE_LEGEND` (FR-002; depends on T002)
- [X] T011 [US2] Wire `Legend` into `App.tsx`'s layout, visible whenever nodes can be shown (depends on T010)
- [X] T012 [US2] Wire `modeDescription` into `Toolbar.tsx`, rendering the active mode's description and updating it based on `connectingFrom` in "connect" mode (FR-003; depends on T003)
- [X] T013 [US2] Give the intersection node marker a text label (`"I"`) in `NetworkCanvas.tsx`'s `NODE_STYLE`, closing the one node type whose meaning relied on shape/color alone (FR-007)

### Tests for User Story 2

- [X] T014 [P] [US2] Vitest tests for `modeDescription` in `web/src/content/copy.test.ts`: a distinct, non-empty description for each of the three modes, and distinct pre-/post-first-click descriptions for "connect" (depends on T003)

**Checkpoint**: A legend and per-mode descriptions are visible — User Stories 1 and 2 both independently pass.

---

## Phase 5: User Story 3 - The tool looks and behaves like a finished product (Priority: P3)

**Goal**: Every control shares a consistent, deliberate visual style, and
color never carries meaning alone.

**Independent Test**: Visually inspect every control for a single
consistent style with no unstyled defaults; view the app in grayscale and
confirm node types and result comparisons remain understandable from text
alone.

### Implementation for User Story 3

- [X] T015 [US3] Restyle `Toolbar.tsx` using the T004 design tokens: replace raw fieldset/radio bullet lists with a styled, labeled control group (FR-005, FR-006; depends on T004)
- [X] T016 [US3] Restyle `NetworkCanvas.tsx`'s wrapping container (bordered card, background) and `App.tsx`'s overall layout (header, sidebar, main content) using the T004 design tokens (FR-005, FR-006; depends on T004)
- [X] T017 [US3] Restyle `ResultsPanel.tsx` and the "Load example"/"Run" buttons using the T004 design tokens, keeping the before/after comparison's wording (not just its color) prominent (FR-005, FR-006, FR-007; depends on T004)

### Validation for User Story 3

- [X] T018 [US3] Walk through quickstart.md's Acceptance Scenario 3 and the grayscale check by hand, confirming SC-004 (no unstyled default control remains) and SC-005 (color-coded meaning also reads from text/icon alone)

**Checkpoint**: All three user stories are independently functional — the frontend is explanatory and visually consistent, closing out the usability gap reported after feature 005.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T019 [P] Run `npx tsc --noEmit` and `npx oxlint` in `web/`; fix any findings
- [X] T020 Run the full quickstart.md validation end to end: `npm test`, `npm run build`, then the manual walkthrough (depends on T009, T014, T018)
- [X] T021 [P] Confirm via `git diff --stat` that no file under `graph/`, `agent/`, `simulation/`, or `cmd/server` changed (constitution Principle I)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on Foundational only — independent of User Story 1's components (different files: `IntroBanner`/canvas-hint vs. `Legend`/`Toolbar`)
- **User Story 3 (Phase 5)**: Depends on Foundational (T004's tokens) and touches the same files Users Stories 1-2 added content to, so it is sequenced after both
- **Polish (Phase 6)**: Depends on all three user stories being complete

Unlike some earlier features, User Stories 1 and 2 here are genuinely
independent of each other (disjoint files) and could be built in either
order or in parallel; User Story 3 necessarily comes last since it
restyles what both of them added.

### Within Each User Story

- Foundational (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green, or — for visual work —
  manually verified) before moving to the next

### Parallel Opportunities

- `T002`, `T003`, `T004` (Foundational, independent pieces of separate files) in parallel
- User Story 1 (`T005`-`T009`) and User Story 2 (`T010`-`T014`) can proceed in parallel — disjoint files
- `T015`, `T016`, `T017` (US3, different components) in parallel once T004 lands
- `T019` and `T021` (Polish) in parallel

---

## Parallel Example: Foundational Phase

```bash
Task: "Define APP_INTRO, LegendEntry, NODE_TYPE_LEGEND in web/src/content/copy.ts"
Task: "Implement modeDescription/emptyCanvasHint/missingDemandHint in web/src/content/copy.ts"
Task: "Rewrite web/src/index.css with design tokens and base control styles"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: open the app fresh and confirm the intro + empty-canvas guidance are visible
5. This alone fixes the single biggest reported gap: a visitor no longer faces an unexplained blank screen

### Incremental Delivery

1. Setup + Foundational → copy module and design tokens ready
2. User Story 1 → validate independently (onboarding works)
3. User Story 2 → validate independently (legend + mode descriptions work)
4. User Story 3 → validate independently (consistent visual style, closing out the feature)
5. Polish → lint/build check + confirm zero Go-side changes

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- No dismissible banner, guided tour, new UI framework/library, or change
  to editing-mode behavior, by design (spec.md Assumptions) — do not add
  them here
