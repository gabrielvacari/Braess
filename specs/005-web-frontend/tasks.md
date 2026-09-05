---

description: "Task list template for feature implementation"
---

# Tasks: Web Frontend

**Input**: Design documents from `/specs/005-web-frontend/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api-contract.md, quickstart.md. Depends on features 001-004 (`graph`, `agent`, `simulation`), none of which this feature modifies.

**Tests**: Included. `go test` covers the new `cmd/server` API contract; Vitest covers the frontend's pure (non-canvas) logic — network editing rules and route/animation mapping. Canvas rendering and animation itself are verified manually via quickstart.md's Acceptance-Scenario walkthrough, per plan.md's Testing strategy (not practical to unit-test pixels).

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

Two codebases, per plan.md: `cmd/server/` (Go, joins the existing
module) and `web/` (a separate npm project). Neither touches `graph/`,
`agent/`, or `simulation/`.

---

## Phase 1: Setup

- [X] T001 [P] Create the `cmd/server/` directory
- [X] T002 [P] Initialize the `web/` npm project: `package.json` (react, react-dom, konva, react-konva, vite, vitest, typescript), `tsconfig.json`, `vite.config.ts`, `index.html`, and a stub `src/main.tsx`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The API contract (both sides) and the client edit model
every user story builds on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 [P] Define server DTOs (`NodeDTO`, `TravelTimeDTO`, `EdgeDTO`, `DemandDTO`, `RunRequest`, `RouteGroupDTO`, `DemandResultDTO`, `RunResponse`, `ErrorResponse`) in `cmd/server/dto.go` (data-model.md)
- [X] T004 Implement request → engine mapping in `cmd/server/dto.go`: build a `*graph.Graph` from `RunRequest.Nodes`/`Edges` and a `[]simulation.Demand` from `RunRequest.Demands`; reject a self-loop edge (`From == To`) before it reaches `graph.AddEdge` (FR-008; depends on T003)
- [X] T005 Implement response mapping in `cmd/server/dto.go`: `simulation.MultiAssignmentResult` → `RunResponse`, grouping each demand's `Routes` by identical edge-ID sequence into `RouteGroupDTO` (research.md decision #4; depends on T003)
- [X] T006 Implement the `POST /api/run` handler in `cmd/server/handlers.go`: decode `RunRequest`, apply the T004/T005 mappings, call `simulation.RunDemands`, encode `RunResponse`, and map failures to `400` (malformed/self-loop) or `422` (an `agent.ErrNoRoute`-wrapping or other `RunDemands` failure) with an `ErrorResponse` body (contracts/api-contract.md; depends on T004, T005)
- [X] T007 Wire `http.ServeMux` and `http.Server` in `cmd/server/main.go`, registering the T006 handler at `POST /api/run` (depends on T006)
- [X] T008 [P] Define client types (`ClientNode`, `ClientEdge`, `ClientDemand`, `NetworkState`) and edit operations (`addNode`, `addEdge`, `removeNode`, `removeEdge`) in `web/src/model/network.ts`, rejecting a self-loop in `addEdge` and cascading node removal to its edges/demands in `removeNode` (data-model.md; contracts/api-contract.md; FR-005, FR-008)
- [X] T009 [P] Define `RunRequest`/`RunResponse` TypeScript types and a `runSimulation()` fetch wrapper in `web/src/model/api.ts` (contracts/api-contract.md)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Watch the simulation happen on a map (Priority: P1) 🎯 MVP

**Goal**: Load an already-computed network and population, render it, and
animate each agent traveling its route.

**Independent Test**: Load the page against a preset network/population
and confirm every node/edge is drawn and every agent's icon visibly
progresses along its route to its destination — no editing required yet.

### Implementation for User Story 1

- [X] T010 [US1] Implement `NetworkCanvas.tsx` (`react-konva` `Stage`/`Layer`): render each node with a shape/color distinct per type (FR-001, SC-001) and each edge as a line between its two nodes
- [X] T011 [US1] Implement a pure position-at-time function plus the `useAgentAnimation` hook in `web/src/animation/useAgentAnimation.ts`: given a `RouteGroupDTO` (edges + travel time), compute an agent icon's position along the route at a given elapsed time (FR-002)
- [X] T012 [US1] Wire `App.tsx`: a "load example" action building a preset network/population, calling `runSimulation()` (T009), and rendering the result through `NetworkCanvas` (T010) with animated agents (T011)

### Tests for User Story 1

- [X] T013 [P] [US1] Go contract test in `cmd/server/handlers_test.go`: a valid `RunRequest` returns `200` with correctly grouped `RouteGroups` (depends on T006)
- [X] T014 [P] [US1] Vitest test for the pure position-at-time function in `web/src/animation/useAgentAnimation.test.ts`: position at `t=0` is the route's start, at `t=travelTime` is its end, and progresses monotonically across intermediate edges (depends on T011)

**Checkpoint**: Loading a preset network shows a rendered, animated simulation — User Story 1 is independently functional and demoable.

---

## Phase 4: User Story 2 - Draw and edit the network by hand (Priority: P2)

**Goal**: Add/remove nodes and roads, and mark node types, directly on
the map.

**Independent Test**: Starting from an empty map, place nodes, mark
types, connect them with roads, then remove a road and a node, confirming
each edit is reflected immediately — independent of whether a simulation
has run.

### Implementation for User Story 2

- [X] T015 [US2] Implement `Toolbar.tsx`: choose a node type before placing one, a "draw road" mode, and remove-selected actions
- [X] T016 [US2] Wire `NetworkCanvas.tsx` for interactive editing: click-to-place a node (using the Toolbar's selected type), click two nodes in sequence to draw a road between them, click-to-select plus a remove action for a node or road (depends on T010, T015)

### Tests for User Story 2

- [X] T017 [P] [US2] Vitest tests for `network.ts`'s edit operations (T008) covering Acceptance Scenarios 1-3: placing a house/company/intersection, drawing a road between two nodes, removing a road, removing a node cascades to its roads (and any demand referencing it), and a self-loop attempt is rejected

**Checkpoint**: A user can build a network from scratch entirely by hand — User Story 2 is independently functional and demoable.

---

## Phase 5: User Story 3 - See how adding a road changes congestion, live (Priority: P3)

**Goal**: Declare demands, run the population against the hand-drawn
network, edit it, re-run, and compare both results.

**Independent Test**: Draw a small network, run it, note the average
travel time, add a road, re-run, and confirm both averages remain visible
together.

### Implementation for User Story 3

- [X] T018 [US3] Implement a demand editor (origin, destination, size) within `Toolbar.tsx` or a sibling component, backed by `ClientDemand` (T008)
- [X] T019 [US3] Implement `ResultsPanel.tsx`: display the current and previous `RunResponse`'s total/average travel time side by side (FR-007, SC-004)
- [X] T020 [US3] Wire the "Run" action in `App.tsx`: call `runSimulation()` with the current `NetworkState` + demands, move the prior result into "previous" (`ResultsPanel`), stop any in-progress animation before starting a new one on edit (research.md decision #6), and surface a `422` error clearly rather than crashing (FR-009; depends on T009, T012, T016, T018, T019)

### Tests for User Story 3

- [X] T021 [P] [US3] Go contract test in `cmd/server/handlers_test.go`: a demand with no route between its origin and destination returns `422` with an `ErrorResponse` naming that specific demand (FR-009; depends on T006)
- [X] T022 [P] [US3] Vitest test for the current/previous result state transition (T020's history logic, extracted as a pure reducer/function) in `web/src/model`

### Validation for User Story 3

- [X] T023 [US3] Walk through quickstart.md's three Acceptance-Scenario demonstrations end to end (load & watch; draw & edit; before/after comparison), confirming SC-001 through SC-005 by hand

**Checkpoint**: All three user stories are independently functional — roadmap Phase 5 is complete, and the project's original goal (draw roads, watch agents, see the paradox) is fully realized end to end.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T024 [P] Run `gofmt -l .` / `go vet ./...` for the new `cmd/server` code, and `tsc --noEmit` (plus any configured linter) for `web/`; fix any findings
- [X] T025 Run the full quickstart.md validation end to end: `go test ./...`, `cd web && npm test`, then the manual walkthrough (depends on T013, T014, T017, T021, T022, T023)
- [X] T026 [P] Add doc comments: a `cmd/server` package doc noting it's the only Go code touching HTTP/JSON and depends on `graph`/`agent`/`simulation` without modifying them (Principle I); a top-of-file note in `web/src/model/network.ts` stating the client model never computes a route or a travel time itself (Principle III boundary, research.md decision #3)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (T006, T009) only
- **User Story 2 (Phase 4)**: Depends on Foundational (T008) and on User Story 1's `NetworkCanvas` (T010) to extend
- **User Story 3 (Phase 5)**: Depends on User Story 1's run pipeline (T012) and User Story 2's editing (T016) — there's nothing to re-run or compare without both
- **Polish (Phase 6)**: Depends on all three user stories being complete

As in every prior feature, later stories here build on top of earlier
ones rather than run fully in parallel — each remains independently
*testable* and *demoable* per its own Independent Test above.

### Within Each User Story

- Foundational (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green, or — for canvas work —
  manually verified) before moving to the next

### Parallel Opportunities

- `T001`, `T002` (Setup) in parallel
- `T003` (server DTOs) and `T008`/`T009` (client model/API types) in parallel — different codebases entirely
- `T013` and `T014` (US1 tests, different codebases) in parallel
- `T021` and `T022` (US3 tests, different codebases) in parallel
- `T024` and `T026` (Polish) in parallel

---

## Parallel Example: Foundational Phase

```bash
# Server-side and client-side foundational work can proceed together:
Task: "Define server DTOs in cmd/server/dto.go"
Task: "Define client types and edit operations in web/src/model/network.ts"
Task: "Define RunRequest/RunResponse TS types and fetch wrapper in web/src/model/api.ts"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: load a preset network in the browser and watch agents move
5. This alone is the headline payoff of the whole roadmap — the engine's already-proven behavior, finally visible

### Incremental Delivery

1. Setup + Foundational → API contract and client model ready on both sides
2. User Story 1 → validate independently (load & watch works)
3. User Story 2 → validate independently (draw & edit works)
4. User Story 3 → validate independently (before/after comparison works — closing out the roadmap)
5. Polish → whole-feature formatting/vet/lint/doc pass

---

## Notes

- [P] tasks touch different files (often different codebases entirely) with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow) — note two toolchains are now involved (`go` and `npm`)
- No persistence, no multi-user/multi-tab editing, and no live tick-by-tick server-pushed simulation, by design (spec.md Assumptions) — do not add them here
