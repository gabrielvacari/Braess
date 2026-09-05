# Implementation Plan: Congestion and Braess's Paradox

**Branch**: `003-congestion-braess-paradox` | **Date**: 2026-09-04 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/003-congestion-braess-paradox/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Add a new `simulation` package (depending on `agent`, `graph`, and the Go
standard library only) that runs a `Population` of agents sharing one fixed
origin/destination to a stable route assignment via best-response dynamics:
agents are loaded incrementally against real per-edge traffic volume, then
repeatedly allowed to switch to a strictly better route given everyone
else's current choices, until nobody can improve (or a round limit is hit).
This requires a small, backward-compatible extension to feature 002's
`agent` package — a per-edge-volume routing entry point, since Phase 2's
`ShortestRoute` only understood one flat baseline volume for the whole
graph. Extend `cmd/graphcli` to run the classic four-node Braess network
with and without its extra road and print both outcomes side by side,
proving the paradox in text before any frontend work is allowed to start
(constitution Principle II).

## Technical Context

**Language/Version**: Go 1.26 (same module, `braess`, established in feature 001)

**Primary Dependencies**: Go standard library only — reuses `agent`'s existing `container/heap`-based Dijkstra; no new third-party dependency

**Storage**: N/A — in-memory, consistent with features 001-002

**Testing**: Go standard `testing`, table-driven — same style as `graph` and `agent`

**Target Platform**: Same as prior features (any Go toolchain target)

**Project Type**: Extends the existing single-module Go project with a new package (`simulation`), a small additive change to `agent`, and a further extension of the CLI

**Performance Goals**: SC-003 requires a population of ≥100 agents on a ≤10-node network to finish and print in under 5 seconds. Even the classic Braess example at its literature-standard scale (4000 agents, 4 nodes) is trivial for this design: each agent's reconsideration is one Dijkstra run on a tiny graph, and convergence is expected within a small number of rounds (see research.md decision #1)

**Constraints**: Zero UI/frontend dependency (FR-009); each agent's reconsideration must be its own independent computation against currently observable volumes, never a centrally computed optimum handed down (FR-003); a road's travel time must come from the real current count of agents on it (FR-001); the process must terminate — either by detecting stability or by hitting a bounded round limit — and must say clearly which one happened (FR-005, FR-006)

**Scale/Scope**: One fixed origin/destination pair (per Phase 2's constraint, carried forward); population sizes from 0 up to several thousand identical agents; the fixed classic Braess network (4 nodes) for User Story 3, plus small hand-built networks for testing the general equilibrium mechanism

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | `simulation` imports only `agent`, `graph`, and the Go standard library; `cmd/graphcli` depends on all three, never the reverse | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | This feature *is* the roadmap's explicit pre-frontend gate: AGENTS.md requires validating Braess's Paradox in text/logs before any UI work — User Story 3 exists specifically to satisfy that | PASS |
| III. Decentralized, Selfish Agent Behavior | Every reconsideration during refinement is one agent's own `agent.ShortestRouteAtVolumes` call against currently-observable volumes — never a centrally solved optimal assignment; this is best-response dynamics, the standard model of selfish behavior in congestion games | PASS |
| IV. Validate in Text Before Visually | `cmd/graphcli` (extended) and `go test ./...` remain the only verification surfaces | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions, including why this converges at all, are recorded in `research.md` | PASS |
| Technology Constraints | Go, standard library only, no frontend touched | PASS |
| Development Workflow | Branch `003-congestion-braess-paradox`; Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/003-congestion-braess-paradox/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   ├── agent-api-extension.md
│   └── simulation-api.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
go.mod                     # unchanged (module braess, go 1.26)

graph/                     # feature 001 — unchanged by this feature

agent/                     # feature 002 — EXTENDED, non-breaking
├── shortestpath.go          # internal Dijkstra core factored out; adds
│                              # ShortestRouteAtVolumes (per-edge volumes);
│                              # existing ShortestRoute(volume float64)
│                              # keeps its exact signature and behavior
└── (other files unchanged)

simulation/                 # NEW — depends on agent + graph + stdlib
├── population.go            # Population type, DefaultMaxRounds
├── run.go                    # AssignmentResult, Run (incremental load +
│                              # best-response refinement)
├── run_test.go                # SC-001 (3 networks), SC-004, FR-010
└── braess_test.go              # the classic 4-node network, both variants
                                # (automated check backing SC-002)

cmd/
└── graphcli/
    └── main.go               # EXTENDED — builds the classic Braess
                                # network with/without its extra road, runs
                                # simulation.Run on each, prints both
                                # results and the comparison
```

**Structure Decision**: `simulation` becomes the third layer in the
one-way dependency chain established across features 001-002 (`graph` →
`agent` → `simulation`), matching AGENTS.md's own framing: the graph is the
network, agents are individual decision-makers, and what this feature adds
— many agents reacting to each other over time — is a distinct concern
from either. `agent` is extended, not replaced: `ShortestRoute`'s existing
signature and behavior are preserved exactly (feature 002's own contract
already anticipated growth here), so feature 002's tests keep passing
unmodified.

## Complexity Tracking

*No violations — table intentionally omitted.*
