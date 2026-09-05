---

description: "Task list template for feature implementation"
---

# Tasks: Multiple Origins and Destinations

**Input**: Design documents from `/specs/004-multiple-origins-destinations/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/simulation-api-extension.md, quickstart.md. Depends on features 001-003, none of which this feature modifies destructively (feature 003's `simulation` package is extended non-breakingly).

**Tests**: Included, for the same reason as prior features — plan.md's Testing strategy relies on `go test ./...` to prove SC-001/SC-002/SC-003 and the FR-007/FR-008 paths, plus a regression check that feature 003's `Run`/`Population`/`AssignmentResult` and its exact Braess numbers are untouched.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

Generalizes `simulation/run.go` (feature 003) internally, adds
`simulation/demand.go` (new public API), and further extends
`cmd/graphcli/main.go`.

---

## Phase 1: Setup

- [X] T001 [P] Create `simulation/demand.go` per plan.md's Project Structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The shared multi-demand core and data shapes every user
story in this feature needs

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 Generalize `loadIncrementally` in `simulation/run.go` to accept a flat `[]odPair{origin, destination string; demandIndex int}` instead of a `Population`, returning `(routes []agent.Route, volumes map[string]float64, failedAt int, err error)` (contracts/simulation-api-extension.md)
- [X] T003 Generalize `refine` in `simulation/run.go` the same way, returning `(rounds int, converged bool, failedAt int, err error)` (depends on T002)
- [X] T004 Update `Run` in `simulation/run.go` to build a same-pair `[]odPair` from its `Population` and call the generalized core from T002/T003, preserving `Run`'s exact existing signature and behavior (depends on T003)
- [X] T005 Run `go test ./simulation/...` and confirm every feature-003 test — including the exact Braess 65 → 80 numbers in `braess_test.go` — still passes unmodified: the regression check for the T002-T004 refactor (depends on T004)
- [X] T006 [P] Define `Demand`, `MultiPopulation`, `DemandResult`, and `MultiAssignmentResult` in `simulation/demand.go` (data-model.md; contracts/simulation-api-extension.md)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Multiple flows share real congestion (Priority: P1) 🎯 MVP

**Goal**: Several distinct origin/destination demands, run together,
genuinely share congestion on any road they both use.

**Independent Test**: Build a network where two demands are forced
through one shared road, and confirm changing one demand's size measurably
changes the other's resulting travel time, via `go test
./simulation/...` — no UI involved.

### Implementation for User Story 1

- [X] T007 [US1] Implement `RunDemands` in `simulation/demand.go`: build the flat `[]odPair` (tagging each agent's `demandIndex`), call the generalized `loadIncrementally`/`refine` from Phase 2, compute the overall `TotalTravelTime`/`AverageTravelTime`, and wrap any failure with a `demandError` helper naming the specific `Demand` by index and origin/destination (FR-001, FR-002, FR-003, FR-004, FR-008; depends on T003, T006)

### Tests for User Story 1

- [X] T008 [P] [US1] Tests for `RunDemands` in `simulation/demand_test.go`, covering Acceptance Scenarios 1-2 (a shared bottleneck road's travel time reflects combined traffic from both demands; increasing one demand's size measurably changes the other's reported average — SC-001) and FR-008 (a demand with no route produces an error naming that specific demand) (depends on T007)

**Checkpoint**: `go test ./simulation/...` passes for shared congestion across demands — User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Results reported per pair, not just overall (Priority: P2)

**Goal**: Each demand's own total/average travel time is available
alongside one combined overall number.

**Independent Test**: Run two distinct demands together and confirm each
demand's own reported result reflects only that demand's own agents, via
`go test ./simulation/...` — no UI involved.

### Implementation for User Story 2

- [X] T009 [US2] Implement per-demand grouping in `RunDemands` (`simulation/demand.go`): populate each `DemandResult`'s `Routes`/`TotalTravelTime`/`AverageTravelTime` from the shared final assignment using each agent's `demandIndex`, and handle a zero-size demand cleanly (FR-005, FR-007; depends on T007)

### Tests for User Story 2

- [X] T010 [P] [US2] Tests for per-demand reporting in `simulation/demand_test.go`, covering Acceptance Scenarios 1-2 (each demand's own numbers reflect only its own agents; one overall total/average is also available) and FR-007 (a zero-agent demand reports a trivial result without affecting other demands) (depends on T009)

**Checkpoint**: `go test ./simulation/...` now covers shared congestion and per-demand reporting — User Stories 1 and 2 both independently pass.

---

## Phase 5: User Story 3 - Inspect multiple pairs from a terminal (Priority: P3)

**Goal**: The existing standalone terminal command, extended, prints each
demand's own result plus the overall result for a small multi-demand
network.

**Independent Test**: Run the terminal command and check its console
output and exit code — no browser or UI required.

### Implementation for User Story 3

- [X] T011 [US3] Extend `cmd/graphcli/main.go`: build a small network where at least 3 distinct house/company demands share a bottleneck road, run `simulation.RunDemands`, and print each demand's own total/average plus one overall total/average (FR-006, FR-009; depends on T009)

### Tests for User Story 3

- [X] T012 [P] [US3] Automated test in `simulation/demand_test.go` covering SC-002: across 3 distinct multi-demand networks, verify the final assignment is a genuine equilibrium — no agent, from any demand, could unilaterally switch route and strictly improve (depends on T009)

### Validation for User Story 3

- [X] T013 [US3] Run `go run ./cmd/graphcli` per quickstart.md and confirm the printed per-demand/overall results and a `0` exit code (SC-003; depends on T011)

**Checkpoint**: All three user stories are independently functional — roadmap Phase 4 is complete, closing out the simulation engine ahead of Phase 5's frontend.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T014 [P] Run `gofmt -l .` and `go vet ./...` across `simulation/` and the updated `cmd/graphcli/`; fix any findings
- [X] T015 Run the full quickstart.md validation end-to-end (`go test ./...` then `go run ./cmd/graphcli`) (depends on T005, T008, T010, T012, T013)
- [X] T016 [P] Add a doc comment to `simulation/demand.go` stating it generalizes feature 003's single-pair congestion to multiple demands sharing one network, and that every agent's reconsideration remains independent regardless of which demand it belongs to (constitution Principle III at multi-commodity scale)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories; includes the `Run` regression check (T005), which must be green before any `RunDemands` work builds on the same generalized core
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on User Story 1's `RunDemands` (T007) — per-demand grouping is additional logic inside the same function
- **User Story 3 (Phase 5)**: Depends on User Story 2's complete `RunDemands` (T009) — nothing meaningful to print or check equilibrium over without it
- **Polish (Phase 6)**: Depends on all three user stories being complete

As in prior features, later stories here build on top of earlier ones
rather than run fully in parallel — each remains independently *testable*
per its own Independent Test above.

### Within Each User Story

- Foundational (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green) before moving to the next

### Parallel Opportunities

- `T006` (new types) can be written in parallel with `T002`-`T004` (core refactor), since they touch different files — but `T007` needs both done
- `T008` (US1 tests) can be written in parallel with starting `T009` (US2 implementation), since `T009` only needs `T007` to exist, not US1's tests to be finished
- `T012` (equilibrium test) and `T011` (CLI) touch different files and can proceed in parallel once `T009` lands
- `T014` and `T016` (Polish) in parallel

---

## Parallel Example: Foundational Phase

```bash
# T006 (new types) can proceed alongside the T002-T004 core refactor:
Task: "Define Demand, MultiPopulation, DemandResult, MultiAssignmentResult in simulation/demand.go"
Task: "Generalize loadIncrementally/refine in simulation/run.go and update Run"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (includes the feature-003 regression check)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `go test ./simulation/...` passes, including the SC-001 cross-effect test
5. This alone proves congestion is genuinely shared across demands — the load-bearing new idea of roadmap Phase 4

### Incremental Delivery

1. Setup + Foundational → foundation ready, feature 003 confirmed regression-free
2. User Story 1 → validate independently (shared congestion proven)
3. User Story 2 → validate independently (per-demand reporting proven)
4. User Story 3 → validate independently (CLI demonstration — closing out roadmap Phase 4)
5. Polish → whole-feature formatting/vet/doc pass

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- This feature does not let agents choose among multiple destinations, and
  has no frontend, by design (spec.md Assumptions) — do not add them here;
  the frontend belongs to roadmap Phase 5 and its own `/speckit-specify` run
