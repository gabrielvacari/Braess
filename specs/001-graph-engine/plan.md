# Implementation Plan: Graph Simulation Engine

**Branch**: `001-graph-engine` | **Date**: 2026-09-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-graph-engine/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Build a standalone, in-memory Go package that models a road network as a
graph: nodes (intersections, houses, companies) and directed edges (roads)
each carrying a length, a capacity, and a swappable travel-time function of
current traffic volume. Ship it alongside a minimal terminal command that
constructs a sample network and prints its full structure — proving the
engine works and is verifiable in text with zero UI dependency, per
constitution Principles I, II, and IV.

## Technical Context

**Language/Version**: Go 1.26 (matches the locally installed toolchain; constitution mandates Go for the simulation engine)

**Primary Dependencies**: Go standard library only (`fmt`, `errors`, `flag`) — no third-party packages; nothing here needs more, and pulling in a dependency for a graph + CLI would violate the spirit of Principle I (keep the engine self-contained) for no benefit

**Storage**: N/A — in-memory only, per spec Assumptions

**Testing**: Go standard `testing` package, table-driven tests, run via `go test ./...`

**Target Platform**: Any platform the Go toolchain targets (developer machine / CI); no OS-specific code

**Project Type**: Single-module Go library + a thin CLI entry point (`cmd/`)

**Performance Goals**: Not a driver at this scale — SC-001 only requires printing a 4-node/4-edge network in under 5 seconds, which is trivial for in-memory graph operations. Graph construction and lookups should be O(1)/O(number of edges) — no algorithmic work is needed to hit this goal, just avoiding accidentally quadratic code.

**Constraints**: Zero UI/frontend dependencies (FR-006, SC-004); no negative-volume travel-time evaluation (FR-008); parallel edges between the same node pair must be independently addressable (FR-005)

**Scale/Scope**: Small hand-built or programmatically-built networks (tens of nodes/edges) — enough to reproduce the classic four-node Braess's Paradox example in a later phase; no scale requirement beyond that is implied by this spec

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | Package `graph` has no import of any UI/rendering package; `cmd/graphcli` depends on `graph`, never the reverse | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | This plan covers roadmap Phase 1 only — pure graph, no agents, no frontend — matching required order | PASS |
| III. Decentralized, Selfish Agent Behavior | Not yet applicable (no agents in this phase); the engine does not perform any central route optimization, it only models roads/travel-time so future agents can decide for themselves | N/A this phase |
| IV. Validate in Text Before Visually | `cmd/graphcli` and `go test ./...` are the only verification surfaces defined here; no visual/graphical output exists yet | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions and their rationale are recorded in `research.md` below | PASS |
| Technology Constraints | Engine written in Go, no frontend framework touched | PASS |
| Development Workflow | Work happens on branch `001-graph-engine` (Gitflow feature branch); commits will follow Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/001-graph-engine/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
│   └── graph-api.md
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
go.mod                   # module braess, go 1.26

graph/                   # the engine — importable, zero UI dependency (Principle I)
├── node.go              # Node type (id, NodeType: Intersection/House/Company)
├── edge.go              # Edge type (id, from, to, length, capacity, TravelTimeFunc)
├── traveltime.go         # TravelTimeFunc type + Linear/Constant constructors (see research.md)
├── graph.go              # Graph type: AddNode, AddEdge, Nodes, Edges, TravelTime
├── node_test.go
├── edge_test.go
└── graph_test.go

cmd/
└── graphcli/              # standalone terminal entry point (FR-006, User Story 3)
    └── main.go            # builds/loads a sample network, prints nodes/edges/travel times
```

**Structure Decision**: A single Go module using idiomatic Go layout rather
than the template's generic `src/`/`tests/` split: the engine lives in
package `graph` at the repo root (importable as a library, and the natural
home for `_test.go` files colocated with the code they test, per Go
convention) and the standalone CLI required by FR-006 lives under
`cmd/graphcli`, which depends on `graph` — never the other way around, so
Principle I (Engine Independent of UI) is structurally enforced, not just
promised. No `backend/`/`frontend/` split applies yet since no frontend
exists in this phase.

## Complexity Tracking

*No violations — table intentionally omitted.*
