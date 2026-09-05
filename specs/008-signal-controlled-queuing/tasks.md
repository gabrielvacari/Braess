---

description: "Task list template for feature implementation"
---

# Tasks: Signal-Controlled Queuing

**Input**: Design documents from `/specs/008-signal-controlled-queuing/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/queuesim-api.md, quickstart.md. Depends on feature 001 (`graph`) only — `agent`, `simulation`, `cmd/server`, and `web/` are untouched.

**Tests**: Included — `go test` covers signal phase math, queue discharge accounting, the no-crossing-on-red invariant, and the two-road scenario's reportability, mirroring every prior engine feature's rigor. Console validation via `cmd/graphcli` per constitution Principle IV.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

New package `queuesim/` (depends on `graph` + stdlib only), plus a further
extension of `cmd/graphcli/main.go`. No other file is touched.

---

## Phase 1: Setup

- [X] T001 [P] Create the `queuesim/` package directory

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The core types and the wait-aware routing every user story's
simulation loop needs

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Define `Signal` and `IsGreenAt(t)` in `queuesim/signal.go` (data-model.md; contracts/queuesim-api.md; research.md decision #4)
- [X] T003 [P] Define `Demand` in `queuesim/demand.go` (data-model.md; research.md decision #5)
- [X] T004 [P] Define `AgentReport`, `QueueSample`, and `RunResult` in `queuesim/report.go` (data-model.md; research.md decision #6)
- [X] T005 Implement `queuesim`'s own wait-aware shortest-route function in `queuesim/route.go`: prices each edge as free-flow travel time, plus — for a signal-controlled edge — an estimated wait from that signal's state at the given time (research.md decisions #2, #3; depends on T002)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Agents queue at a red signal and are released on green (Priority: P1) 🎯 MVP

**Goal**: A signal-controlled road holds arriving agents during red and
releases them, in order, at a defined rate once green.

**Independent Test**: Run a single signal-controlled road on a fixed
cycle with a steady stream of arriving agents; confirm from the report
that no agent crosses during red, and queued agents leave in arrival
order once green.

### Implementation for User Story 1

- [X] T006 [US1] Implement the tick-based spawn/advance loop in `queuesim/run.go`: spawn agents per each `Demand.ArrivalInterval` up to its `Count`, and advance every in-flight moving agent's remaining edge time by `tick` each step (depends on T003, T004, T005)
- [X] T007 [US1] Implement signal queue join/discharge in `queuesim/run.go`: an agent finishing an edge and about to enter a signal-controlled edge that's red joins that edge's FIFO queue; each tick, while green, discharge queued agents using a fractional accumulator at a default rate derived from the edge's own free-flow travel time (spec Assumptions; depends on T006)
- [X] T008 [US1] Record each spawned agent's `AgentReport` (spawn time, travel time, wait time, `Arrived`) and one `QueueSample` per signal per tick into `RunResult` as the loop executes (FR-004, FR-006, FR-008; depends on T007)

### Tests for User Story 1

- [X] T009 [P] [US1] Table-driven tests for `Signal.IsGreenAt` in `queuesim/signal_test.go`: correct at phase boundaries; two signals offset by one `GreenDuration` are never green at the same time (depends on T002)
- [X] T010 [P] [US1] Tests in `queuesim/run_test.go` for a single signal-controlled edge with steady arrivals: no agent's crossing time falls within a red phase (SC-002); a green phase long enough to clear the preceding red's arrivals fully drains the queue (SC-001); arrivals outpacing the discharge rate carry the remainder into the next cycle rather than dropping it (Edge Case; depends on T008)
- [X] T011 [P] [US1] Test in `queuesim/run_test.go`: an agent still queued or moving when `duration` elapses is reported with `Arrived: false` (SC-005; depends on T008)

**Checkpoint**: `go test ./queuesim/...` passes for the core queue mechanic — User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Observe what alternating signals actually do to congestion (Priority: P2)

**Goal**: Two roads with always-opposite signals, run with a steady
agent stream, produce a report showing each road's queue-length-over-time
and wait experience — without presuming which pattern emerges.

**Independent Test**: Run the two-road, opposite-phase scenario for a
fixed duration with a steady arrival rate; confirm the report shows each
road's queue length over time and each agent's total wait time.

### Implementation for User Story 2

- [X] T012 [US2] Build the classic two-road, opposite-phase scenario (two nodes, two signal-controlled roads configured per research.md decision #4, one `Demand` with a steady arrival stream) as a reusable fixture in `queuesim/two_road_test.go` (depends on T002, T003)
- [X] T013 [US2] Extend `cmd/graphcli/main.go` to run this scenario via `queuesim.Run` and print each road's queue length at intervals across the run, plus an arrival/wait summary (FR-009, constitution Principle IV; depends on T012)

### Tests for User Story 2

- [X] T014 [P] [US2] Test asserting the two-road run's `RunResult` yields a computable, per-road queue-length time series and per-road average wait time — sufficient to tell from the data whether one road ended up more congested, it oscillated, or it balanced, without the test itself asserting which (SC-003; depends on T012)
- [X] T015 [P] [US2] Test asserting `AgentReport.WaitTime` and `.TravelTime` are both populated and distinguishable per agent (FR-006; depends on T012)

**Checkpoint**: `go test ./queuesim/...` and the CLI both demonstrate the experiment end to end — User Stories 1 and 2 both independently pass.

---

## Phase 5: User Story 3 - Configure signal timing to run different experiments (Priority: P3)

**Goal**: A signal's green/red durations are set per scenario, so
different timing plans can be compared.

**Independent Test**: Run the same two-road scenario twice with two
different green/red duration configurations and confirm the reported
queue patterns differ.

### Tests for User Story 3

- [X] T016 [P] [US3] Test running the two-road scenario (T012) twice with two distinct `GreenDuration`/`RedDuration` configurations and asserting the resulting queue-length patterns differ (SC-004; depends on T012) — no new implementation needed: `Signal`'s fields are already caller-configured since T002

**Checkpoint**: All three user stories are independently functional — the signal-controlled queuing mechanism is complete and validated in text, per constitution Principle II, ready for a future frontend feature to build on.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T017 [P] Run `gofmt -l .` and `go vet ./...` across `queuesim/` and the updated `cmd/graphcli/`; fix any findings
- [X] T018 Run the full quickstart.md validation end to end: `go test ./queuesim/...` then `go run ./cmd/graphcli` (depends on T009, T010, T011, T014, T015, T016)
- [X] T019 [P] Confirm via `git diff --stat` that `graph/`, `agent/`, `simulation/`, `cmd/server/`, and `web/` are all unchanged (constitution Principle I; plan.md Constraints)
- [X] T020 [P] Add a package doc comment to `queuesim/run.go` (or a `doc.go`) stating the package's scope: depends only on `graph` + stdlib, models discrete simulated time, and does not presume the two-road experiment's outcome (constitution Principle V)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on Foundational directly, and reuses User Story 1's queue mechanic (T007/T008) to actually produce meaningful data
- **User Story 3 (Phase 5)**: Depends on User Story 2's fixture (T012) — it's purely a test proving existing configurability, not new code
- **Polish (Phase 6)**: Depends on all three user stories being complete

As in every prior feature, later stories here build on top of earlier
ones rather than run fully in parallel — each remains independently
*testable* per its own Independent Test above.

### Within Each User Story

- Foundational (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green) before moving to the next

### Parallel Opportunities

- `T002`, `T003`, `T004` (Foundational, independent files) in parallel
- `T009`, `T010`, `T011` (US1 tests, independent test functions) in parallel once their shared dependency (T008) lands
- `T014` and `T015` (US2 tests) in parallel once T012 lands
- `T017` and `T019` (Polish) in parallel

---

## Parallel Example: Foundational Phase

```bash
Task: "Define Signal and IsGreenAt in queuesim/signal.go"
Task: "Define Demand in queuesim/demand.go"
Task: "Define AgentReport, QueueSample, RunResult in queuesim/report.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `go test ./queuesim/...` proves a single signal-controlled road correctly queues and releases agents
5. This alone proves the core new mechanic — real waiting — the load-bearing idea this whole feature exists to add

### Incremental Delivery

1. Setup + Foundational → wait-aware routing and core types ready
2. User Story 1 → validate independently (the queue/release mechanic is correct)
3. User Story 2 → validate independently (the actual experiment is observable, whatever it shows)
4. User Story 3 → validate independently (timing is genuinely configurable, not hardcoded)
5. Polish → whole-feature formatting/vet/doc pass, confirming zero changes outside `queuesim`/`cmd/graphcli`

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- No browser UI, no mid-trip route replanning, and no multi-movement
  intersections, by design (spec.md Assumptions) — do not add them here;
  a frontend for this capability is a distinct, later feature and its own
  `/speckit-specify` run, per constitution Principle II
