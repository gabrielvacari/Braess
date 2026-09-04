---

description: "Task list template for feature implementation"
---

# Tasks: Shortest-Path Agents

**Input**: Design documents from `/specs/002-shortest-path-agents/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/agent-api.md, quickstart.md. Depends on feature 001 (`graph` package), which is not modified by this feature.

**Tests**: Included, for the same reason as feature 001 — plan.md's Testing strategy and quickstart.md's validation steps both rely on `go test ./...` to prove SC-001/SC-003/SC-004 and the FR-005/FR-006/FR-008 error/edge-case paths.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

New package `agent/` (depends only on `graph` + stdlib, per plan.md's
"Structure Decision"), plus an extension of the existing
`cmd/graphcli/main.go`.

---

## Phase 1: Setup

- [ ] T001 [P] Create the `agent/` package directory per plan.md's Project Structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared shapes both user stories 1 and 2 produce/consume

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T002 [P] Define the `Route` struct (`Edges []graph.Edge`, `TotalTravelTime float64`) in `agent/route.go` (data-model.md "Route")
- [ ] T003 [P] Define the `ErrNoRoute` sentinel error in `agent/shortestpath.go` (contracts/agent-api.md)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Compute an agent's shortest route (Priority: P1) 🎯 MVP

**Goal**: Given an origin, a destination, and a graph, compute the
minimum-total-travel-time route between them at a fixed baseline volume.

**Independent Test**: Build a small graph with a known correct shortest
path and check `ShortestRoute`'s result matches it, via `go test
./agent/...` — no other agents, congestion, or UI involved.

### Implementation for User Story 1

- [ ] T004 [US1] Implement `ShortestRoute(g, from, to, volume)` (Dijkstra via `container/heap`) in `agent/shortestpath.go` (FR-001, FR-002, FR-004, FR-005, FR-006, FR-008; depends on T002, T003)

### Tests for User Story 1

- [ ] T005 [P] [US1] Table-driven tests for `ShortestRoute` in `agent/shortestpath_test.go`, covering Acceptance Scenarios 1-2, the "no route", "origin == destination", "parallel edges", and "travel-time evaluation error" Edge Cases, and SC-001 (≥5 distinct graphs) (depends on T004)

**Checkpoint**: `go test ./agent/...` passes for route computation — User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Multiple agents decide independently (Priority: P2)

**Goal**: Several `Agent`s sharing the same origin/destination each compute
their own route through their own independent call.

**Independent Test**: Create several independent `Agent`s for the same
origin/destination and confirm each one's route comes from its own
`ComputeRoute` call, via `go test ./agent/...` — no UI involved.

### Implementation for User Story 2

- [ ] T006 [US2] Define the `Agent` struct (`Origin`, `Destination`) and implement `ComputeRoute(g)` in `agent/agent.go` (FR-003; depends on T004 — calls `ShortestRoute` with `volume = 0`)

### Tests for User Story 2

- [ ] T007 [P] [US2] Table-driven tests for `Agent.ComputeRoute` in `agent/agent_test.go`, covering Acceptance Scenarios 1-2 and SC-003 (≥3 independently-created agents sharing an origin/destination yield identical, independently-computed routes) (depends on T006)

**Checkpoint**: `go test ./agent/...` now covers both route computation and independent agent behavior — User Stories 1 and 2 both independently pass.

---

## Phase 5: User Story 3 - Inspect agent routing from a terminal (Priority: P3)

**Goal**: The existing standalone terminal command, extended, prints each
simulated agent's route and travel time.

**Independent Test**: Run the command and check its console output and
exit code — no browser or UI required.

### Implementation for User Story 3

- [ ] T008 [US3] Extend `cmd/graphcli/main.go`: add a fixed origin/destination on the existing sample network, construct 2-3 independent `agent.Agent` values, call `ComputeRoute` on each, and print each one's ordered route (edge IDs and endpoints) and total travel time (FR-007; depends on T006)

### Validation for User Story 3

- [ ] T009 [US3] Run `go run ./cmd/graphcli` per quickstart.md and confirm the printed agent routes and a `0` exit code (SC-002) (depends on T008)

**Checkpoint**: All three user stories are independently functional — the feature is complete for roadmap Phase 2.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T010 [P] Run `gofmt -l .` and `go vet ./...` across `agent/` and the updated `cmd/graphcli/`; fix any findings
- [ ] T011 Run the full quickstart.md validation end-to-end (`go test ./...` then `go run ./cmd/graphcli`) (depends on T005, T007, T009)
- [ ] T012 [P] Add a package doc comment to `agent/agent.go` stating the package's scope: depends only on `graph` + stdlib (Principle I), and that every `Agent` computes its route independently (Principle III)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on Foundational; `T006` also depends on `T004` (it calls `ShortestRoute`)
- **User Story 3 (Phase 5)**: Depends on Foundational, and on User Story 2 (`Agent`/`ComputeRoute`) to have something to print
- **Polish (Phase 6)**: Depends on all three user stories being complete

As in feature 001, User Story 3 (the CLI) is sequenced after the story it
demonstrates rather than run fully in parallel with it — but each story
remains independently *testable* per its own Independent Test above.

### Within Each User Story

- Foundational types (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green) before moving to the next

### Parallel Opportunities

- `T002`, `T003` (Foundational, separate files) in parallel
- `T005` (US1 tests) can be written in parallel with starting `T006` (US2 implementation) once `T004` lands, since `T006` only needs `ShortestRoute` to exist, not US1's tests to be finished
- `T010` and `T012` (Polish) in parallel

---

## Parallel Example: Foundational Phase

```bash
# Launch the two independent shared-shape tasks together:
Task: "Define the Route struct in agent/route.go"
Task: "Define the ErrNoRoute sentinel error in agent/shortestpath.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `go test ./agent/...` passes for route computation
5. This alone proves shortest-path routing works — a real, demoable slice of roadmap Phase 2 ("simple agents going from A to B with shortest-path routing, no congestion yet")

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. User Story 1 → validate independently (route computation proven correct)
3. User Story 2 → validate independently (independent-agent behavior proven)
4. User Story 3 → validate independently (standalone CLI proves the whole feature end to end, closing out roadmap Phase 2)
5. Polish → whole-feature formatting/vet/doc pass

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- This feature has no congestion modeling and no multiple distinct origin/
  destination pairs by design (spec.md Assumptions) — do not add them here;
  they belong to roadmap Phases 3 and 4 and their own `/speckit-specify` runs
