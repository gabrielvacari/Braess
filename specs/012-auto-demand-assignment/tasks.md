---

description: "Task list template for feature implementation"
---

# Tasks: Automatic House Demand

**Input**: Design documents from `/specs/012-auto-demand-assignment/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md. Depends on `queueNetwork.ts`/`queueApi.ts`/`SignalToolbar.tsx`/`SignalApp.tsx` (features 009-011), modified in place.

**Tests**: Included — Vitest for the new pure assignment function, per this project's established pattern.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2)

## Path Conventions

`web/src/model/` (new `autoDemand.ts`+test, modified `queueNetwork.ts`,
`queueNetwork.test.ts`, `queueApi.ts`) and `web/src/components/SignalToolbar.tsx`,
`web/src/SignalApp.tsx` (modified). No engine, server, or equilibrium-mode file.

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: The assignment logic every user story depends on

- [X] T001 [P] Implement `AUTO_DEMAND_COUNT`, `AUTO_DEMAND_ARRIVAL_INTERVAL`, `HouseAssignment`, `assignHousesToCompanies`, and `autoTimeDemands` in `web/src/model/autoDemand.ts` (data-model.md)
- [X] T002 [P] Vitest tests in `web/src/model/autoDemand.test.ts`: no companies -> empty; one company -> every house assigned to it; 2+ houses and 2+ companies -> at least 2 distinct companies used; calling twice with the same nodes agrees (FR-004, FR-006, FR-007; depends on T001)

**Checkpoint**: The assignment rule is validated before anything renders or requests against it

---

## Phase 2: User Story 1 - Every house sends traffic automatically (Priority: P1) 🎯 MVP

**Goal**: Signals mode has no demand form; every house automatically generates traffic when run.

**Independent Test**: Draw two houses and one company, connect them with roads, run the scenario, and confirm both houses' agents appear without ever declaring a demand by hand.

### Implementation for User Story 1

- [X] T003 [US1] Remove `TimeDemand`, `QueueNetworkState.demands`, `addTimeDemand`, `removeTimeDemand` from `web/src/model/queueNetwork.ts`; drop the `demands` line from `removeNodeFromQueueNetwork`'s cascade (FR-001; data-model.md "Removed")
- [X] T004 [US1] Update `web/src/model/queueNetwork.test.ts`: remove the demand-CRUD tests, and update the cascade test to only assert roads (depends on T003)
- [X] T005 [US1] Update `web/src/model/queueApi.ts`'s `toQueueRunRequest` to build `demands` via `autoTimeDemands(state.nodes)` instead of mapping `state.demands` (FR-002, FR-003; depends on T001, T003; research.md decision #4)
- [X] T006 [US1] Remove the demand-form section and its props (`demands`, `onAddDemand`, `onRemoveDemand`) from `web/src/components/SignalToolbar.tsx` (FR-001; depends on T003)
- [X] T007 [US1] Remove demand-related state/handlers from `web/src/SignalApp.tsx`; base the "Run" button's disabled state and its hint text on `autoTimeDemands(network.nodes).length === 0` instead of `network.demands.length === 0` (FR-007; depends on T005, T006)

**Checkpoint**: Running a network with houses and companies produces traffic with no demand form anywhere — User Story 1 is independently functional.

---

## Phase 3: User Story 2 - Traffic spreads across more than one company (Priority: P1)

**Goal**: The automatic assignment is visible, and demonstrably spreads across more than one company when more than one exists.

**Independent Test**: Draw at least two houses and two companies, run the scenario, and confirm at least two distinct companies each receive at least one house's traffic; confirm the toolbar shows which company each house was assigned to before running.

### Implementation for User Story 2

- [X] T008 [US2] In `web/src/components/SignalToolbar.tsx`, render `assignHousesToCompanies(nodes)` as a short read-only list (e.g. "House 1 → Company 2"), replacing the removed demand form's screen space (FR-010; depends on T001, T006; research.md decision #3)

**Checkpoint**: Both user stories are independently functional.

---

## Phase 4: Polish & Cross-Cutting Concerns

- [X] T009 [P] Run `npx tsc --noEmit` and `npx oxlint` for `web/`; fix any findings
- [X] T010 Run `npm test` in `web/` and confirm every remaining suite passes (SC-001 through SC-004 covered by T002; T004 confirms no orphaned demand-CRUD tests remain)
- [X] T011 [P] Confirm via `git diff --stat` that only `web/src/model/{autoDemand.ts,autoDemand.test.ts,queueNetwork.ts,queueNetwork.test.ts,queueApi.ts}`, `web/src/components/SignalToolbar.tsx`, and `web/src/SignalApp.tsx` changed — no engine, server, or equilibrium-mode file (constitution Principle I; FR-008, FR-009)
- [X] T012 Manual walkthrough of quickstart.md end to end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: BLOCKS both user stories
- **User Story 1 (Phase 2)**: Depends on Foundational
- **User Story 2 (Phase 3)**: Depends on Foundational and reuses User Story 1's `SignalToolbar` edits (same file, same area)
- **Polish (Phase 4)**: Depends on both user stories

### Parallel Opportunities

- `T001` and `T002` (Foundational) developed together, `T002` verifying `T001`
- `T009` and `T011` (Polish) in parallel

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Foundational
2. Complete Phase 2: User Story 1
3. **STOP and VALIDATE**: a network runs with automatic traffic, no demand form
4. User Story 2's visibility list is additive on top of an already-working MVP

### Incremental Delivery

1. Foundational → assignment rule validated
2. User Story 1 → validate independently (automatic traffic works)
3. User Story 2 → validate independently (spread is visible and guaranteed)
4. Polish → lint/build check + confirm the change stayed confined to `web/`

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- No configurable rate, no per-house override, and no change to the
  Equilibrium mode's demand form, by design (spec.md Assumptions) — do
  not add them here
