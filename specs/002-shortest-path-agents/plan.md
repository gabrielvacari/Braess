# Implementation Plan: Shortest-Path Agents

**Branch**: `002-shortest-path-agents` | **Date**: 2026-09-04 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-shortest-path-agents/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Add a new `agent` package, depending only on the existing `graph` package
(feature 001) and the Go standard library, that computes an agent's
shortest-time route between a fixed origin and destination via Dijkstra's
algorithm evaluated at a fixed baseline traffic volume (zero — no
congestion yet). Each `Agent` computes its own route through its own
independent call, with no shared/cached result across agents, directly
exercising constitution Principle III for the first time. Extend the
existing `cmd/graphcli` to spawn a few agents on a fixed route and print
each one's path and travel time, keeping verification in text (Principle
IV) with zero UI dependency (Principle I).

## Technical Context

**Language/Version**: Go 1.26 (same module, `braess`, established in feature 001)

**Primary Dependencies**: Go standard library only — `container/heap` for Dijkstra's priority queue, `errors`/`fmt` for error handling. No third-party graph/shortest-path library; the network sizes involved (tens of nodes) don't justify one, and it would be an unnecessary dependency in an engine that Principle I requires to stay self-contained.

**Storage**: N/A — in-memory, consistent with feature 001

**Testing**: Go standard `testing`, table-driven, `go test ./...` — same style as `graph`

**Target Platform**: Same as feature 001 (any Go toolchain target)

**Project Type**: Extends the existing single-module Go project with a new package (`agent`) and extends the existing CLI

**Performance Goals**: SC-002 requires completing and printing a route for a 10-node/15-edge network in under 5 seconds — trivial for Dijkstra at this scale (O((V+E) log V) with a binary heap); no special optimization needed

**Constraints**: Zero UI/frontend dependency (FR-007); each agent's route must come from its own independent computation, never a shared/cached answer (FR-003); parallel edges must be treated as distinct route options (FR-004); a missing route or a failed travel-time evaluation must surface as a clear error, never silently substituted (FR-005, FR-008)

**Scale/Scope**: Same small hand-built networks as feature 001; one fixed origin/destination pair per this phase (multiple distinct pairs are roadmap Phase 4)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | `agent` imports only `graph` and the Go standard library; `cmd/graphcli` depends on both, never the reverse | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | This is roadmap Phase 2 (shortest-path agents, still no congestion, still no UI), immediately following the completed Phase 1 | PASS |
| III. Decentralized, Selfish Agent Behavior | First phase where this is directly testable: FR-003 requires every `Agent.ComputeRoute` call to be independent, with no shared/central route table — enforced by giving `Agent` no package-level state to read from | PASS |
| IV. Validate in Text Before Visually | `cmd/graphcli` (extended) and `go test ./...` remain the only verification surfaces; no visual output exists | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions recorded in `research.md` below | PASS |
| Technology Constraints | Go, standard library only, no frontend touched | PASS |
| Development Workflow | Branch `002-shortest-path-agents` (Gitflow feature branch); Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/002-shortest-path-agents/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── agent-api.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
go.mod                    # unchanged (module braess, go 1.26)

graph/                    # feature 001 — unchanged by this feature

agent/                    # NEW — depends only on graph + stdlib (Principle I)
├── route.go              # Route type: ordered []graph.Edge + TotalTravelTime
├── shortestpath.go        # ShortestRoute(g, from, to, volume) via Dijkstra; ErrNoRoute
├── agent.go               # Agent type (Origin, Destination) + ComputeRoute method
├── shortestpath_test.go
└── agent_test.go

cmd/
└── graphcli/
    └── main.go            # EXTENDED — also builds a few agents on a fixed
                             # origin/destination and prints each route
```

**Structure Decision**: A new top-level package `agent`, mirroring
AGENTS.md's own architecture blocks ("1. Graph", "2. Agents") as separate
Go packages with a one-way dependency (`agent` → `graph`, never back) —
the same idiomatic-Go approach chosen in feature 001, now extended rather
than re-litigated. `graph` itself is not modified: `agent` builds the
adjacency it needs from `graph.Graph.Nodes()`/`Edges()`, which are already
public (feature 001, FR-007), keeping feature 001 stable while this
feature is purely additive.

## Complexity Tracking

*No violations — table intentionally omitted.*
