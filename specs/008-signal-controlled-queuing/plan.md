# Implementation Plan: Signal-Controlled Queuing

**Branch**: `008-signal-controlled-queuing` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/008-signal-controlled-queuing/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Add a new `queuesim` package (depending only on `graph`, `agent`, and the
Go standard library — `simulation` is untouched) implementing a
fixed-tick, discrete-time simulation: signals gate whether a road is
passable; agents spawn over time per demand, choose a route once at spawn
using a snapshot estimate of current signal wait, and queue in arrival
order at any signal-controlled road they reach while it's red, discharging
at a defined rate once it turns green. A run reports each signal-controlled
road's queue length over time and each agent's total wait vs. travel time.
"Opposite phase" (FR-005) is achieved purely through signal configuration
(matching cycle lengths, complementary offset) — no special linking
construct is needed. Extend `cmd/graphcli` with the classic two-road
alternating-signal scenario, printed to the console (constitution
Principle IV) — a browser UI is explicitly deferred to a later feature.

## Technical Context

**Language/Version**: Go 1.26 (same module, `braess`, established in feature 001)

**Primary Dependencies**: Go standard library only — no new dependency; reuses nothing from `simulation` (a different computational model — tick-based, not equilibrium-based) but does reuse `graph.Graph`/`graph.Edge`/`graph.TravelTimeFunc` (feature 001) directly, and follows the same internal Dijkstra-core pattern `agent` established (feature 002/003) for its own, separate wait-aware route function.

**Storage**: N/A — in-memory, one `Run` call produces one in-memory `RunResult`; consistent with every prior feature.

**Testing**: Go standard `testing`, table-driven — signal phase math, queue discharge accounting, and the two-road scenario's queue/wait reports, mirroring the rigor of `simulation`'s own equilibrium tests (feature 003/004).

**Target Platform**: Same as every prior feature (any Go toolchain target).

**Project Type**: Extends the existing single-module Go project with a new package (`queuesim`) and a further extension of `cmd/graphcli`; no change to `graph`, `agent`, `simulation`, `cmd/server`, or `web/`.

**Performance Goals**: A run of the two-road scenario over a demonstrative duration (e.g. a few simulated minutes, at a 1-second tick) must complete and print in well under 5 seconds — trivial at this scale (a few hundred to a few thousand ticks, each doing O(agents-in-flight) work).

**Constraints**: No change to `graph`/`agent`/`simulation`/`cmd/server`/`web` (spec Assumptions); an agent's route is chosen once, at spawn, and not revised mid-trip (research.md decision #2); an agent already past a signal's queue point when it turns red is not interrupted (spec Assumptions); this feature does not model intersections with more than two mutually-exclusive movements, amber phases, or turns (spec Assumptions).

**Scale/Scope**: One new package (`queuesim`); the primary validated scenario is the classic two-road, opposite-phase setup with a steady agent arrival stream over a multi-cycle duration.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | `queuesim` depends only on `graph`, `agent`, and the Go standard library; `cmd/graphcli` depends on it, never the reverse; no UI/frontend code touched | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | This feature is itself an engine-only capability, validated via `cmd/graphcli`/`go test` before any browser work — a frontend for it is explicitly a separate, later feature (spec Assumptions, FR-009) | PASS |
| III. Decentralized, Selfish Agent Behavior | Each agent computes its own route independently at spawn, using only currently-observable signal state (a snapshot estimate) — no central scheduler assigns routes or queue positions | PASS |
| IV. Validate in Text Before Visually | `cmd/graphcli`'s new scenario prints queue-length-over-time and wait reports as text; `go test` covers signal phase math, discharge accounting, and the two-road scenario | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions — especially the fixed-tick model, spawn-time-snapshot routing, and "opposite phase via configuration" — are recorded in research.md | PASS |
| Technology Constraints | Go, standard library only | PASS |
| Development Workflow | Branch `008-signal-controlled-queuing`; Conventional Commits | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/008-signal-controlled-queuing/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── queuesim-api.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
go.mod                      # unchanged (module braess, go 1.26)
graph/, agent/, simulation/   # unchanged by this feature
cmd/server/, web/              # unchanged by this feature

queuesim/                    # NEW — depends on graph + agent + stdlib only
├── signal.go                  # Signal type, IsGreenAt(t), phase math
├── demand.go                   # Demand type (origin, destination, count, arrival interval)
├── route.go                     # wait-aware route selection: an internal Dijkstra variant
│                               # pricing a signal-controlled edge as free-flow + a snapshot
│                               # estimate of current wait, at spawn time only
├── run.go                        # the fixed-tick simulation loop: spawns agents over time,
│                               # advances in-flight agents, applies signal/queue rules
├── report.go                      # QueueSample, AgentReport, RunResult
├── signal_test.go
├── route_test.go
├── run_test.go
└── two_road_test.go                # the classic two-road, opposite-phase scenario

cmd/
└── graphcli/
    └── main.go                  # EXTENDED — the two-road scenario, printing queue-length-
                                  # over-time and wait reports
```

**Structure Decision**: A new, self-contained package rather than an
extension of `simulation` — feature 003/004's equilibrium/best-response
model and this feature's discrete-tick model are different computational
approaches to different questions ("what stable assignment do selfish
agents reach" vs. "how does a queue build and drain over real time"), and
forcing one into the other's internals would blur both. `queuesim` reuses
`graph`'s types directly and follows the same "one shared Dijkstra core"
pattern `agent` established, but implements its own — the two pathfinding
problems (aggregate-volume-priced vs. spawn-time-wait-priced) are
different enough that literally sharing `agent`'s internal `shortestRoute`
core isn't a natural fit (see research.md decision #3).

## Complexity Tracking

*No violations — table intentionally omitted.*
