# Implementation Plan: Web Frontend

**Branch**: `005-web-frontend` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/005-web-frontend/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Add a stateless Go HTTP API (`cmd/server`) that translates JSON in and out
of the existing `graph`/`agent`/`simulation` types — one endpoint,
`POST /api/run`, that takes a whole network + demand set and returns a
`simulation.RunDemands` result — with zero changes to those packages. Add
a separate React + Konva.js single-page app (`web/`) that owns the network
being drawn/edited as client-side state, calls that endpoint to run it,
and animates each distinct route the response reports. All routing and
congestion computation stays server-side in Go; the browser never
re-implements simulation logic, only renders and edits.

## Technical Context

**Language/Version**: Go 1.26 (server, unchanged module) + TypeScript 5.x / Node 20+ (frontend, new, separate npm project)

**Primary Dependencies**:
- Server: Go standard library only (`net/http`'s `ServeMux`, `encoding/json`) — no router or web framework, consistent with every prior feature's "stdlib first" decisions.
- Frontend: `react`, `react-dom`, `konva`, `react-konva` (the constitution's named stack, via its idiomatic React bindings), `vite` (dev server/bundler), `vitest` (unit tests for pure, non-canvas logic).

**Storage**: N/A — the server is stateless (no session, no database); the browser holds the in-progress network as page-lifetime state only (spec Assumptions: no cross-session persistence).

**Testing**: Go standard `testing` for the new `cmd/server` handlers (request/response JSON contract, error mapping); Vitest for the frontend's pure logic (network edit operations, request/response mapping); canvas rendering and animation verified manually via quickstart.md's walkthrough of the spec's Acceptance Scenarios (not practical to unit-test pixels).

**Target Platform**: Server runs anywhere the Go toolchain does (developer machine); frontend runs in any modern desktop browser, served by Vite's dev server locally.

**Project Type**: Web application — Go backend (`cmd/server`, using the unchanged `graph`/`agent`/`simulation`) + a separate frontend project (`web/`), per the confirmed architecture decision.

**Performance Goals**: SC-002's 5-second budget (load → agents visibly moving) is dominated by network/render overhead, not computation — `simulation.RunDemands` already completes in well under a second at this project's established scale (Phase 3's 4000-agent Braess run, Phase 4's SC-003 ≤15-node/≥100-agent case). The main risk to that budget is payload size for large populations, addressed below.

**Constraints**: Zero change to `graph`/`agent`/`simulation` (spec Assumptions); the browser must never recompute routes or congestion itself (constitution Principle III stays true only if all decision-making logic stays server-side); a road from a node to itself must be rejected (FR-008); removing a node must cascade to its roads (FR-005); an unroutable demand must be reported clearly, not crash (FR-009).

**Scale/Scope**: Same network/population scale as Phases 3-4 (≤15 nodes for the performance target, populations up to a few thousand agents for the Braess-style case); single user, single network, single browser tab at a time (spec Assumptions).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Engine Independent of UI | `graph`/`agent`/`simulation` are not modified; `cmd/server` depends on them (never the reverse) and is the only Go code that knows JSON/HTTP exist; `web/` doesn't touch Go at all except over HTTP | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | This is the first feature allowed to touch UI — Phases 1-4 already validated the engine in text, satisfying the gate | PASS |
| III. Decentralized, Selfish Agent Behavior | Unaffected: this feature adds no new routing/congestion logic anywhere. Guarded by design — see "Constraints" above and research.md decision #3 (the browser is a thin edit/render/animate shell; `RunDemands` remains the only place a route is ever decided) | PASS |
| IV. Validate in Text Before Visually | The new server endpoint is covered by `go test` exactly like every prior feature; the frontend's pure (non-canvas) logic is covered by Vitest; only pixel rendering/animation itself is verified by hand, which is inherent to "the visualize step" and not a gap in the earlier, already-satisfied text validation | PASS |
| V. Explain Reasoning, Not Just Code | Design decisions — the stateless server, route-grouping for large populations, TypeScript/Vite choice — are recorded in research.md | PASS |
| Technology Constraints | Frontend is React + Konva.js (via `react-konva`), as specified; server stays in Go | PASS |
| Development Workflow | Branch `005-web-frontend`; Conventional Commits; two toolchains now involved (`go test`/`go build` and `npm test`/`npm run build`), both documented in quickstart.md | PASS |

No violations. Complexity Tracking table is not needed.

## Project Structure

### Documentation (this feature)

```text
specs/005-web-frontend/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/             # Phase 1 output (/speckit-plan command)
│   └── api-contract.md
└── tasks.md               # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
go.mod                     # unchanged (module braess, go 1.26)
graph/, agent/, simulation/  # unchanged by this feature

cmd/
├── graphcli/                # unchanged
└── server/                  # NEW: stateless HTTP API
    ├── main.go                # http.Server wiring, ServeMux routes
    ├── dto.go                  # JSON request/response types + mapping to/from graph/agent/simulation types
    ├── handlers.go              # POST /api/run handler
    └── handlers_test.go          # contract tests: valid run, validation errors, no-route errors

web/                        # NEW: separate npm project (not part of the Go module)
├── package.json
├── tsconfig.json
├── vite.config.ts
├── index.html
└── src/
    ├── main.tsx
    ├── App.tsx
    ├── components/
    │   ├── NetworkCanvas.tsx        # react-konva Stage/Layer: nodes, roads, agent icons
    │   ├── Toolbar.tsx               # add node (+ type), add road, run, demand editor
    │   └── ResultsPanel.tsx           # current + previous run's total/average, side by side
    ├── model/
    │   ├── network.ts                 # client-side Node/Edge/Demand types + edit operations
    │   ├── network.test.ts             # Vitest: add/remove node/edge, reject self-loop, cascade delete
    │   └── api.ts                       # POST /api/run request/response types + fetch wrapper
    └── animation/
        └── useAgentAnimation.ts          # maps a route group + travel time to an animated position over time
```

**Structure Decision**: Per the confirmed architecture: the Go module and
its three packages stay exactly where they are (zero churn, zero risk to
already-tested code); `cmd/server` is a new sibling to `cmd/graphcli`
inside the same module, following the same "engine stays UI-agnostic,
only `cmd/*` knows about its delivery mechanism" pattern already
established since feature 001; `web/` is a wholly separate npm project
living alongside, not inside, the Go module.

## Complexity Tracking

*No violations — table intentionally omitted.*
