---

description: "Task list template for feature implementation"
---

# Tasks: Graph Simulation Engine

**Input**: Design documents from `/specs/001-graph-engine/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/graph-api.md, quickstart.md

**Tests**: Included. plan.md's Testing strategy (Go `testing`, table-driven) and quickstart.md's validation steps both rely on `go test ./...` to prove SC-002/SC-003 and the FR-003/FR-008 error paths, so test tasks are part of the deliverable, not optional.

**Organization**: Tasks are grouped by user story (spec.md) to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- File paths are exact, per plan.md's Project Structure

## Path Conventions

Single Go module, idiomatic layout (see plan.md "Structure Decision"):
`graph/` (the engine) and `cmd/graphcli/` (the standalone CLI) at repo root.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization

- [X] T001 Create `go.mod` at repo root (`module braess`, `go 1.26`)
- [X] T002 [P] Create the directory skeleton `graph/` and `cmd/graphcli/` per plan.md's Project Structure
- [X] T003 [P] Add Go build artifacts (e.g. a compiled `graphcli` binary) to root `.gitignore`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types every user story builds on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 [P] Define `NodeType` enum and `Node` struct in `graph/node.go` (data-model.md "NodeType", "Node")
- [X] T005 [P] Define the `TravelTimeFunc` type in `graph/traveltime.go` (contracts/graph-api.md)
- [X] T006 [P] Define the `Edge` struct in `graph/edge.go`, referencing `TravelTimeFunc` (data-model.md "Edge (Road)")
- [X] T007 Define the `Graph` struct (internal `nodes`/`edges` maps) and `New()` constructor in `graph/graph.go` (depends on T004, T006)

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Build a road network graph (Priority: P1) 🎯 MVP

**Goal**: Create nodes and edges programmatically and enumerate them back.

**Independent Test**: Construct a small graph in a test and assert node/edge
counts and attributes via `go test ./graph/...` — no agents, congestion
logic, or UI involved.

### Implementation for User Story 1

- [X] T008 [P] [US1] Implement `Graph.AddNode` in `graph/graph.go` (FR-001; rejects empty or duplicate node IDs)
- [X] T009 [US1] Implement `Graph.AddEdge` in `graph/graph.go` (FR-002, FR-003, FR-005; depends on T008 — validates `From`/`To` exist, allows parallel edges between the same pair, rejects empty/duplicate edge IDs)
- [X] T010 [US1] Implement `Graph.Nodes()` and `Graph.Edges()` enumeration in `graph/graph.go` (FR-007; depends on T008, T009)

### Tests for User Story 1

- [X] T011 [P] [US1] Table-driven tests for `AddNode`/`AddEdge`/`Nodes`/`Edges` in `graph/graph_test.go`, covering spec Acceptance Scenarios 1-2 and the "unknown node reference" and "parallel edges" Edge Cases (depends on T008-T010)

**Checkpoint**: `go test ./graph/...` passes for graph construction — User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Compute travel time under traffic (Priority: P2)

**Goal**: Evaluate an edge's travel time for a given traffic volume: free-flow
at zero, never decreasing as volume grows.

**Independent Test**: Feed known volume values to an edge's travel-time
function via `go test ./graph/...` and check the returned times — no agent
or UI involved.

### Implementation for User Story 2

- [X] T012 [P] [US2] Implement `Linear` and `Constant` `TravelTimeFunc` constructors in `graph/traveltime.go` (research.md decision #2)
- [X] T013 [US2] Implement `Graph.TravelTime(edgeID, volume)` in `graph/graph.go` (FR-004, FR-008; depends on T009 for edge lookup and T012 for the functions it evaluates)

### Tests for User Story 2

- [X] T014 [P] [US2] Table-driven tests for `Linear`, `Constant`, and `Graph.TravelTime` in `graph/traveltime_test.go`, covering SC-002 (free-flow time at volume 0), SC-003 (non-decreasing across ≥5 increasing volume samples), and FR-008 (negative volume returns an error) (depends on T012, T013)

**Checkpoint**: `go test ./graph/...` now covers both graph construction and travel-time behavior — User Stories 1 and 2 both independently pass.

---

## Phase 5: User Story 3 - Inspect the graph from a terminal (Priority: P3)

**Goal**: A standalone terminal command builds a sample network and prints
its full structure, proving the engine needs no UI.

**Independent Test**: Run the command and check its console output and exit
code — no browser or UI required.

### Implementation for User Story 3

- [X] T015 [US3] Implement `cmd/graphcli/main.go`: build a sample network (≥4 nodes mixing `NodeType`s, ≥4 edges mixing `Linear`/`Constant` travel-time functions) and print every node's ID/type and every edge's endpoints/length/capacity/current travel time to stdout (FR-006; depends on T008, T009, T010, T012, T013)

### Validation for User Story 3

- [X] T016 [US3] Run `go run ./cmd/graphcli` per quickstart.md and confirm the printed output and a `0` exit code (SC-001) (depends on T015)

**Checkpoint**: All three user stories are independently functional — the engine is complete for roadmap Phase 1.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Whole-engine checks that span every story

- [X] T017 [P] Run `gofmt -l .` and `go vet ./...` across `graph/` and `cmd/graphcli/`; fix any findings
- [X] T018 Run the full quickstart.md validation end-to-end (`go test ./...` then `go run ./cmd/graphcli`) (depends on T011, T014, T016)
- [X] T019 [P] Add a package doc comment to `graph/graph.go` and `cmd/graphcli/main.go` stating each package's scope, noting that `graph` has zero UI dependency (constitution Principle I)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational only
- **User Story 2 (Phase 4)**: Depends on Foundational; `T013` also depends on `T009` (edge lookup) from User Story 1
- **User Story 3 (Phase 5)**: Depends on Foundational, and on User Story 1 (`AddNode`/`AddEdge`/enumeration) and User Story 2 (`TravelTime`, `Linear`/`Constant`) to have something to print
- **Polish (Phase 6)**: Depends on all three user stories being complete

Because User Story 3 (the CLI) exercises the full engine to prove Principle
I end to end, it is sequenced after Stories 1 and 2 here — unlike a typical
spec-kit feature, these stories are not fully parallel with each other, but
each remains independently *testable* per its own Independent Test above.

### Within Each User Story

- Foundational types (Phase 2) before any story's implementation
- Implementation before that story's tests
- Story complete (implementation + tests green) before moving to the next

### Parallel Opportunities

- `T002`, `T003` (Setup) in parallel
- `T004`, `T005`, `T006` (Foundational struct/type definitions, separate files) in parallel
- `T008` (US1) can start as soon as Foundational is done; `T012` (US2, independent of AddNode/AddEdge) can be worked on in parallel with `T008`-`T010`
- `T011` and `T014` (tests) are only parallel with each other once their respective implementation tasks land

---

## Parallel Example: Foundational Phase

```bash
# Launch the three independent type-definition tasks together:
Task: "Define NodeType enum and Node struct in graph/node.go"
Task: "Define the TravelTimeFunc type in graph/traveltime.go"
Task: "Define the Edge struct in graph/edge.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (blocks everything else)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `go test ./graph/...` passes for graph construction
5. This alone proves the graph data structure works — a real, demoable slice per constitution Principle II's "pure graph in code, no UI" step

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. User Story 1 → validate independently (graph construction proven)
3. User Story 2 → validate independently (travel-time behavior proven)
4. User Story 3 → validate independently (standalone CLI proves zero-UI-dependency, closing out roadmap Phase 1 entirely)
5. Polish → whole-engine formatting/vet/doc pass

---

## Notes

- [P] tasks touch different files with no dependency between them
- [Story] labels trace every implementation task back to spec.md
- Commit after each task or logical group, using Conventional Commits (constitution Development Workflow)
- This feature has no persistence, no agents, and no frontend by design (spec.md Assumptions) — do not add them here; they belong to later roadmap phases and their own `/speckit-specify` runs
