---

description: "Task list template for feature implementation"
---

# Tasks: Signal Queuing Frontend

**Input**: Design documents from `/specs/009-signal-frontend/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/queue-api-contract.md, quickstart.md. Depends on features 005 (`cmd/server`'s `NodeDTO`/`EdgeDTO`/`buildGraph`) and 008 (`queuesim`), both reused unmodified.

**Tests**: Included — Go contract tests for the new endpoint, Vitest for the new model and chart geometry, manual walkthrough for canvas/visual output, per every prior frontend feature's precedent.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Path Conventions

`cmd/server/` (new files, reusing `dto.go`'s existing types) and `web/src/`
(new, self-contained Signals-mode files). `graph/`, `agent/`,
`simulation/`, `queuesim/`, and `cmd/graphcli/` are not touched.

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: The DTOs and client model every user story builds on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T001 [P] Define `SignalDTO`, `TimeDemandDTO`, `QueueRunRequest` (embedding `dto.go`'s `NodeDTO`/`EdgeDTO`), `QueueSampleDTO`, `AgentReportDTO`, `QueueRunResponse` in `cmd/server/queuedto.go` (data-model.md)
- [X] T002 [P] Define `DirectedRoad`, `TimeDemand`, `QueueNetworkState`, and edit operations (`addDirectedRoad`, `removeDirectedRoad`, `removeNodeFromQueueNetwork`, `addTimeDemand`, `removeTimeDemand`) in `web/src/model/queueNetwork.ts`, reusing `ClientNode`/`addNode` from `model/network.ts` (data-model.md; contracts/queue-api-contract.md)
- [X] T003 [P] Define `QueueRunRequest`/`QueueRunResponse` TypeScript types, `toQueueRunRequest`, and a `runQueueSimulation` fetch wrapper in `web/src/model/queueApi.ts` (contracts/queue-api-contract.md)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 2: User Story 1 - Configure a signal and a time-based demand (Priority: P1) 🎯 MVP

**Goal**: A user can attach a signal to a drawn road and declare a
time-based demand, purely client-side — no server call needed yet.

**Independent Test**: Draw a road, attach a signal, declare a time-based
demand, and confirm both are reflected in the visible state.

### Implementation for User Story 1

- [X] T004 [US1] Implement `SignalToolbar.tsx`: node-type picker + placement, a "draw directed road" mode with green/red signal-duration inputs, and a time-based demand form (origin, destination, count, arrival interval) (FR-001, FR-002; depends on T002)
- [X] T005 [US1] Implement `SignalApp.tsx`: wires `QueueNetworkState`, `SignalToolbar`, and the reused `NetworkCanvas` (mapping `DirectedRoad.from`/`.to` to its `{a, b}` prop) for editing only — no run wiring yet (depends on T004)

### Tests for User Story 1

- [X] T006 [P] [US1] Vitest tests for `queueNetwork.ts` in `web/src/model/queueNetwork.test.ts`: `addDirectedRoad` rejects a self-loop; `removeNodeFromQueueNetwork` cascades to roads and demands referencing the node; `addTimeDemand`/`removeTimeDemand` work (depends on T002)

**Checkpoint**: A user can fully configure a signal-controlled scenario in the browser — User Story 1 is independently functional and demoable.

---

## Phase 3: User Story 2 - Run the signal-controlled simulation and see queues over time (Priority: P2)

**Goal**: Running a configured scenario shows each signal-controlled
road's queue length over time.

**Independent Test**: Run a small two-road, opposite-signal scenario and
confirm a queue-length-over-time visualization renders for each
signal-controlled road, reflecting real simulation output.

### Implementation for User Story 2

- [X] T007 [US2] Implement request mapping in `cmd/server/queuedto.go`: build `[]queuesim.Signal`/`[]queuesim.Demand` from `QueueRunRequest`, reusing the existing `buildGraph` for `Nodes`/`Edges` (research.md decision #5; depends on T001)
- [X] T008 [US2] Implement the `POST /api/queue-run` handler in `cmd/server/queuehandlers.go` (decode, map via T007, call `queuesim.Run`, encode `QueueRunResponse` or a `400`/`422` error) and register the route in `cmd/server/main.go` (contracts/queue-api-contract.md; depends on T007)
- [X] T009 [US2] Implement `QueueChart.tsx`: a pure `queueChartPoints(samples, width, height)` geometry function plus an SVG `<polyline>` render (research.md decision #4)
- [X] T010 [US2] Wire a "Run" action into `SignalApp.tsx`: call `runQueueSimulation` (T003) with the current `QueueNetworkState`, and render one `QueueChart` per signal-controlled road from the response (depends on T003, T008, T009)

### Tests for User Story 2

- [X] T011 [P] [US2] Go contract test in `cmd/server/queuehandlers_test.go`: a valid request returns `200` with queue samples and agent reports reflecting real queuing (mirroring feature 008's own two-road scenario at a smaller scale) (depends on T008)
- [X] T012 [P] [US2] Vitest test for `queueChartPoints` in `QueueChart.test.ts`: maps samples to in-bounds coordinates and preserves time ordering (depends on T009)

**Checkpoint**: Running a configured scenario shows real, per-road queue-length-over-time charts — User Stories 1 and 2 both independently pass.

---

## Phase 4: User Story 3 - See an arrival and waiting summary (Priority: P3)

**Goal**: After a run, arrived vs. still-waiting agents and average wait
time are visible without reading raw data.

**Independent Test**: Run a scenario where at least one agent doesn't
finish before the run ends, and confirm the summary distinguishes
arrived from still-waiting agents and shows an average wait time.

### Implementation for User Story 3

- [X] T013 [US3] Implement `QueueResultsPanel.tsx`: renders the `QueueChart`(s) plus an arrived/still-waiting count and average wait time, clearly labeled as a signal-controlled result (FR-005, FR-007; depends on T009)
- [X] T014 [US3] Wire `QueueResultsPanel` into `SignalApp.tsx`'s run flow, replacing the ad hoc chart rendering from T010 (depends on T010, T013)

### Tests for User Story 3

- [X] T015 [P] [US3] Go contract test in `cmd/server/queuehandlers_test.go`: a demand with no route to its destination returns `422` with a clear message (FR-008), mirroring feature 005/008's own error-reporting convention (depends on T008)

**Checkpoint**: All three user stories are independently functional.

---

## Phase 5: Mode Separation (FR-006, FR-007)

- [X] T016 Wire a top-level mode toggle in `App.tsx` ("Equilibrium" / "Signals"), rendering the existing equilibrium UI or `SignalApp` — without adding any field to the existing `NetworkState`/`Toolbar`/`ResultsPanel` (depends on T005, T014)
- [X] T017 Manual walkthrough of quickstart.md's scenario end to end: configure and run a two-road signal scenario, confirm switching modes never mixes equilibrium and signal state, and that two different signal timings produce visibly different charts (SC-005; depends on T016)

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T018 [P] Run `gofmt -l .` and `go vet ./...` for `cmd/server`; `npx tsc --noEmit` and `npx oxlint` for `web/`; fix any findings
- [X] T019 Run the full quickstart.md validation end to end: `go test ./...`, `npm test` (in `web/`), then the manual walkthrough (depends on T006, T011, T012, T015, T017)
- [X] T020 [P] Confirm via `git diff --stat` that `graph/`, `agent/`, `simulation/`, `queuesim/`, and `cmd/graphcli/` are unchanged, and that `dto.go`'s/`App.tsx`'s existing exported behavior for the equilibrium mode is untouched (constitution Principle I; FR-006)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: BLOCKS all user stories
- **User Story 1 (Phase 2)**: Depends on Foundational only — pure client-side configuration
- **User Story 2 (Phase 3)**: Depends on Foundational directly, and reuses User Story 1's `SignalApp`/`SignalToolbar` to have something to run
- **User Story 3 (Phase 4)**: Depends on User Story 2's chart/run wiring
- **Mode Separation (Phase 5)**: Depends on all three user stories being complete
- **Polish (Phase 6)**: Depends on Phase 5

### Parallel Opportunities

- `T001`, `T002`, `T003` (Foundational, independent files/codebases) in parallel
- `T007` (server mapping) and `T009` (chart component) in parallel once Foundational lands — different codebases
- `T011` and `T012` (US2 tests) in parallel
- `T018` and `T020` (Polish) in parallel

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Foundational
2. Complete Phase 2: User Story 1
3. **STOP and VALIDATE**: configure a signal and a time-based demand in the browser
4. This proves the configuration surface works before any server round trip is needed

### Incremental Delivery

1. Foundational → DTOs and client model ready
2. User Story 1 → validate independently (configuration works)
3. User Story 2 → validate independently (running shows real queue charts)
4. User Story 3 → validate independently (arrival/wait summary is legible)
5. Mode Separation → the two modes coexist without interfering
6. Polish → lint/build check + confirm zero changes to the engine or the existing equilibrium mode

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- No shared network between modes, no multiple signals per road, and no
  engine change, by design (spec.md Assumptions) — do not add them here
