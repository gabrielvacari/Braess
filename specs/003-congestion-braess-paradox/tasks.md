---

description: "Task list template for feature implementation"
---

# Tasks: Congestion and Braess's Paradox

**Input**: Design documents from `/specs/003-congestion-braess-paradox/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/agent-api-extension.md, contracts/simulation-api.md, quickstart.md. Depends on feature 001 (`graph`, unmodified) and feature 002 (`agent`, extended non-breakingly).

**Tests**: Included, for the same reason as features 001-002 — plan.md's Testing strategy and quickstart.md's validation both rely on `go test ./...` to prove SC-001/SC-002/SC-004 and the FR-005/FR-006/FR-010 paths, plus a regression check that feature 002 is untouched behaviorally.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

New package `simulation/` (depends on `agent` + `graph` + stdlib), a
non-breaking extension of `agent/shortestpath.go` (feature 002), and a
further extension of `cmd/graphcli/main.go`.

---

## Phase 1: Setup

- [ ] T001 [P] Create the `simulation/` package directory per plan.md's Project Structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The per-edge-volume routing capability and shared data shapes
every user story in this feature needs

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T002 Factor `agent/shortestpath.go`'s Dijkstra into an internal `shortestRoute(g, from, to, volumeOf func(string) float64)` core; make `ShortestRoute(g, from, to, volume float64)` a thin wrapper over it with its exact existing signature and behavior preserved (contracts/agent-api-extension.md)
- [ ] T003 Add `ShortestRouteAtVolumes(g, from, to, volumes map[string]float64)` in `agent/shortestpath.go`, delegating to the core from T002 (contracts/agent-api-extension.md; depends on T002)
- [ ] T004 Run `go test ./agent/...` and confirm every feature-002 test still passes unmodified — regression check for the T002/T003 refactor (depends on T003)
- [ ] T005 [P] Define `Population` struct and `DefaultMaxRounds` const in `simulation/population.go` (data-model.md "Population")
- [ ] T006 [P] Define `AssignmentResult` struct in `simulation/run.go` (data-model.md "AssignmentResult")

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Travel time reflects real traffic (Priority: P1) 🎯 MVP

**Goal**: A road's travel time is computed from the actual number of
agents currently assigned to it.

**Independent Test**: Assign a known set of agents to specific roads and
check computed travel times match each road's function evaluated at that
exact count, via `go test ./simulation/...` — no iterative re-routing or
the Braess network needed yet.

### Implementation for User Story 1

- [ ] T007 [US1] Implement volume bookkeeping helpers in `simulation/run.go`: adding a route's edges to a `map[string]float64` volume tally, and removing them (FR-001; depends on T006)

### Tests for User Story 1

- [ ] T008 [P] [US1] Tests for the volume bookkeeping helpers and `agent.ShortestRouteAtVolumes` reflecting real counts in `simulation/run_test.go`, covering Acceptance Scenarios 1-2 (depends on T007)

**Checkpoint**: `go test ./simulation/...` passes for volume tracking — User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Agents settle into a stable assignment (Priority: P2)

**Goal**: A population of agents, loaded incrementally and then refined
round by round, reaches a stable route assignment (or clearly reports that
it didn't within the round budget).

**Independent Test**: Run the process on a small network with a known
correct equilibrium split and confirm it converges to that split within a
small tolerance, via `go test ./simulation/...` — no UI or the specific
Braess network needed.

### Implementation for User Story 2

- [ ] T009 [US2] Implement incremental loading in `simulation/run.go`: assign each agent in turn via `agent.ShortestRouteAtVolumes` against volume accumulated by agents already placed (depends on T003, T007)
- [ ] T010 [US2] Implement best-response refinement rounds in `simulation/run.go`: each round, for each agent, recompute its current route's cost and the best alternative at current volumes (excluding its own contribution), switch only on strict improvement beyond a small epsilon, and stop when a round makes zero switches or `MaxRounds` is reached (FR-003, FR-004, FR-005, FR-006; depends on T009)
- [ ] T011 [US2] Implement `Run(g, p) (AssignmentResult, error)` in `simulation/run.go`, tying together the `Population.Size == 0` short-circuit (FR-010), loading (T009), refinement (T010), and propagating any routing error rather than swallowing it (depends on T005, T006, T010)

### Tests for User Story 2

- [ ] T012 [P] [US2] Table-driven tests for `Run` in `simulation/run_test.go`, covering Acceptance Scenarios 1-2, SC-001 (3 distinct networks with a known equilibrium split, within tolerance), SC-004 (`MaxRounds: 0` reports `Converged: false`), FR-010 (`Size: 0`), and the "no route exists" edge case (depends on T011)

**Checkpoint**: `go test ./simulation/...` now covers volume tracking and full equilibrium-seeking — User Stories 1 and 2 both independently pass.

---

## Phase 5: User Story 3 - Reproduce Braess's Paradox in the console (Priority: P3)

**Goal**: Running the classic four-node network with and without its extra
road, printed, shows the extra road makes the average trip strictly worse.

**Independent Test**: Run the terminal command and check the printed
total/average travel time for "with the extra road" is strictly worse than
"without it" — no UI required.

### Implementation for User Story 3

- [ ] T013 [US3] Extend `cmd/graphcli/main.go`: add a function building the classic Braess network (nodes S/A/B/T; `Linear(0, 0.01)` on S→A and B→T; `Constant(45)` on S→B and A→T; optionally `Constant(0)` on A→B), run `simulation.Run` with a 4000-agent `Population{Origin: "S", Destination: "T", MaxRounds: simulation.DefaultMaxRounds}` for both variants, and print each result plus an explicit comparison line (FR-008, FR-009; depends on T011)
- [ ] T014 [P] [US3] Automated test in `simulation/braess_test.go` building both classic-network variants and asserting the with-shortcut total/average travel time is strictly greater than without (SC-002; depends on T011)

### Validation for User Story 3

- [ ] T015 [US3] Run `go run ./cmd/graphcli` per quickstart.md and confirm the printed comparison (the well-known 65 → 80 result) and a `0` exit code (SC-002, SC-003; depends on T013)

**Checkpoint**: All three user stories are independently functional — the roadmap's pre-frontend gate (AGENTS.md) is satisfied.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T016 [P] Run `gofmt -l .` and `go vet ./...` across `agent/`, `simulation/`, and the updated `cmd/graphcli/`; fix any findings
- [ ] T017 Run the full quickstart.md validation end-to-end (`go test ./...` then `go run ./cmd/graphcli`) (depends on T008, T012, T014, T015)
- [ ] T018 [P] Add a package doc comment to `simulation/run.go` stating the package's scope: depends only on `agent` + `graph` + stdlib (Principle I), and that convergence is guaranteed by potential-game theory, not just observed (Principle V, research.md decision #1)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories; includes the non-breaking `agent` extension, which must be regression-tested (T004) before anything in `simulation` builds on it
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on Foundational; also depends on User Story 1's volume bookkeeping (T007)
- **User Story 3 (Phase 5)**: Depends on User Story 2's `Run` (T011) — nothing to demonstrate without it
- **Polish (Phase 6)**: Depends on all three user stories being complete

As in features 001-002, later stories here are sequenced on top of earlier
ones rather than run fully in parallel — each remains independently
*testable* per its own Independent Test above.

### Within Each User Story

- Foundational (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green) before moving to the next

### Parallel Opportunities

- `T005`, `T006` (Foundational data shapes, separate files) in parallel, once `T004`'s regression check is green
- `T008` (US1 tests) can be written in parallel with starting `T009` (US2 loading) once `T007` lands
- `T013` (CLI) and `T014` (automated Braess test) touch different files and can proceed in parallel once `T011` lands
- `T016` and `T018` (Polish) in parallel

---

## Parallel Example: Foundational Phase (after T002-T004)

```bash
# Launch the two independent shared-shape tasks together:
Task: "Define the Population struct and DefaultMaxRounds in simulation/population.go"
Task: "Define the AssignmentResult struct in simulation/run.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (includes the `agent` extension + its regression check)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `go test ./simulation/...` passes for volume tracking
5. This alone proves travel time genuinely responds to real traffic — the load-bearing new idea of roadmap Phase 3

### Incremental Delivery

1. Setup + Foundational → foundation ready, feature 002 confirmed regression-free
2. User Story 1 → validate independently (real-traffic travel time proven)
3. User Story 2 → validate independently (equilibrium-seeking proven against known splits)
4. User Story 3 → validate independently (the actual Braess's Paradox demonstration — closing out roadmap Phase 3 and its pre-frontend gate)
5. Polish → whole-feature formatting/vet/doc pass

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- This feature has no multiple distinct origin/destination pairs and no
  frontend by design (spec.md Assumptions) — do not add them here; they
  belong to roadmap Phases 4 and 5 and their own `/speckit-specify` runs
