# Implementation Plan: Multiple Origins and Destinations

**Branch**: `004-multiple-origins-destinations` | **Date**: 2026-09-04 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-multiple-origins-destinations/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Generalize feature 003's single-origin/destination `Population`/`Run` to
support several distinct origin/destination `Demand`s sharing one road
network and its real congestion. Internally, `loadIncrementally` and
`refine` are generalized to operate over a flat list of (origin,
destination) pairs — one per agent, tagged with which demand it belongs
to — rather than assuming every agent shares the same pair; `Run` becomes
a thin wrapper over this shared core (exactly as feature 003 did for
`agent.ShortestRoute`), and a new `RunDemands` entry point exposes the
general case, reporting both an overall result and each demand's own.
Extend `cmd/graphcli` with a small network where multiple demands are
forced to share a bottleneck road, printing per-demand and overall results
to prove, in text, that congestion is genuinely shared (constitution
Principle II, IV) before this closes out the simulation-engine side of the
roadmap ahead of Phase 5's frontend.

## Technical Context

**Language/Version**: Go 1.26 (same module, `braess`, established in feature 001)

**Primary Dependencies**: Go standard library only — no new dependency; reuses `agent.ShortestRouteAtVolumes` from feature 003

**Storage**: N/A — in-memory, consistent with features 001-003

**Testing**: Go standard `testing`, table-driven — same style as prior features; feature 003's existing tests (including the exact Braess 65 → 80 numbers) serve as the regression check for this feature's core refactor

**Target Platform**: Same as prior features (any Go toolchain target)

**Project Type**: Extends the existing single-module Go project — generalizes internals of `simulation` (non-breaking to its public `Population`/`Run`), adds new public API to the same package, and further extends the CLI

**Performance Goals**: SC-003 requires ≥3 demands totaling ≥100 agents on a ≤15-node network to finish and print in under 5 seconds — trivial at this scale per feature 003's own performance analysis (each agent's reconsideration is one small Dijkstra run)

**Constraints**: Zero UI/frontend dependency (FR-009); a road's travel time must reflect combined traffic from every demand using it, not one demand at a time (FR-002); each agent's reconsideration remains its own independent computation regardless of which demand it belongs to (FR-003); a demand with no route must be reported clearly, naming that demand, without aborting silently or crashing (FR-008); feature 003's existing `Population`/`Run` behavior and numeric results (the Braess demonstration) must not regress

**Scale/Scope**: Several (at least 3, per SC-003) distinct origin/destination demand pairs, sharing a modestly sized network (≤15 nodes for the performance target, though nothing structurally limits it further); each agent still has exactly one fixed origin/destination for its whole trip (spec Assumptions)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | `simulation`'s new code touches only `agent`, `graph`, and the Go standard library; `cmd/graphcli` depends on it, never the reverse | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | Roadmap Phase 4 — still no UI; this is the last simulation-engine phase the roadmap requires before Phase 5's frontend | PASS |
| III. Decentralized, Selfish Agent Behavior | Every agent's reconsideration, regardless of which demand it belongs to, remains one independent `agent.ShortestRouteAtVolumes` call against shared observable volumes — no central multi-commodity optimizer is introduced | PASS |
| IV. Validate in Text Before Visually | `cmd/graphcli` (extended) and `go test ./...` remain the only verification surfaces | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions, especially the refactor-vs-duplicate trade-off, are recorded in `research.md` | PASS |
| Technology Constraints | Go, standard library only, no frontend touched | PASS |
| Development Workflow | Branch `004-multiple-origins-destinations`; Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/004-multiple-origins-destinations/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── simulation-api-extension.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
go.mod                     # unchanged (module braess, go 1.26)

graph/                     # feature 001 — unchanged by this feature
agent/                     # feature 002/003 — unchanged by this feature

simulation/                # feature 003 — EXTENDED, non-breaking
├── population.go            # unchanged (Population, DefaultMaxRounds)
├── run.go                    # loadIncrementally/refine generalized to a
│                              # shared internal core over []odPair;
│                              # Run becomes a thin wrapper, same
│                              # signature/behavior as feature 003
├── demand.go                  # NEW: Demand, DemandResult,
│                              # MultiAssignmentResult, MultiPopulation,
│                              # RunDemands
├── run_test.go                 # unchanged tests + regression assertions
├── braess_test.go               # unchanged (still passes exactly — the
│                              # regression check for this refactor)
└── demand_test.go                # NEW: SC-001 (shared-congestion cross-
                                # effect), SC-002 (equilibrium check,
                                # general topology), FR-007, FR-008

cmd/
└── graphcli/
    └── main.go               # EXTENDED — a small network with ≥3
                                # demands sharing a bottleneck road,
                                # printing per-demand and overall results
```

**Structure Decision**: Generalize once, inside `simulation`, rather than
duplicating `loadIncrementally`/`refine` for the multi-demand case — the
same "one shared core, thin wrappers on top" pattern already used for
`agent.ShortestRoute`/`ShortestRouteAtVolumes` in feature 003. The
generalization here is small (indexing each agent's own origin/destination
from a per-agent list instead of one fixed pair shared by all), and
feature 003's own tests — including the exact Braess 65 → 80 numbers —
serve as a strong, already-written regression check that the refactor
preserves behavior exactly.

## Complexity Tracking

*No violations — table intentionally omitted.*
