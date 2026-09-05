# Implementation Plan: Automatic House Demand

**Branch**: `012-auto-demand-assignment` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/012-auto-demand-assignment/spec.md`

## Summary

Removes Signals mode's manual demand form (origin/destination/count/
arrival-interval) and `QueueNetworkState.demands` entirely. In its
place, a new pure function derives every house's demand automatically
from the drawn network: 12 agents every 0.8s, targeting a company
chosen by a deterministic round-robin over the network's company
nodes (in creation order) — guaranteeing at least 2 distinct companies
receive traffic whenever 2 or more exist, and identical results across
repeated runs of the same network. The assignment is surfaced in the
toolbar (read-only) so a person can still see which company each house
was assigned to. `queuesim`, `cmd/server`, and the wire format of a
time-based demand are all unchanged — this only changes how the
request's `demands` array gets built on the client.

## Technical Context

**Language/Version**: TypeScript 5 / React 19, `web/` only

**Primary Dependencies**: None new

**Storage**: N/A

**Testing**: Vitest, for the new pure assignment function

**Target Platform**: Browser, Signals mode only

**Project Type**: Web application (existing `web/src/model`, `web/src/components`, `web/src/SignalApp.tsx`)

**Performance Goals**: N/A — same or smaller data volume than the manual form ever produced

**Constraints**: Zero change to `queuesim`, `cmd/server`, the `POST /api/queue-run` wire contract, or the Equilibrium mode (spec FR-008/FR-009)

**Scale/Scope**: `web/` only — one model file removed of its demand CRUD, one new pure-function module, `SignalToolbar.tsx` and `SignalApp.tsx` updated to match

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Status |
|---|---|---|
| I. Engine Independent of UI | Change is entirely inside `web/` — no engine or server package touched | PASS |
| II. Simulate Before You Visualize (NON-NEGOTIABLE) | No new simulation capability — this only changes how an existing, already-validated request field (`demands`) gets populated | PASS (N/A — no new engine capability) |
| III. Decentralized, Selfish Agent Behavior | Not applicable — deciding a demand's origin/destination is not a route choice; agents still find their own route to whichever destination they're assigned, exactly as before | PASS |
| IV. Validate in Text Before Validating Visually | The one new piece of logic (house-to-company assignment) is a pure function, validated by Vitest, independent of any rendering | PASS |
| V. Explain Reasoning, Not Just Code | research.md documents why round-robin (not random or hash-based) was chosen, and FR-010's visibility requirement is honored in the toolbar rather than left implicit | PASS |

No violations — Complexity Tracking not needed.

## Project Structure

### Documentation (this feature)

```text
specs/012-auto-demand-assignment/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md         # Phase 1 output
└── quickstart.md
```

No `contracts/` — no API/server contract changes.

### Source Code (repository root)

```text
web/src/
├── model/
│   ├── autoDemand.ts        # NEW: assignHousesToCompanies, autoTimeDemands
│   ├── autoDemand.test.ts   # NEW
│   ├── queueNetwork.ts      # MODIFIED: remove TimeDemand, demands field, add/removeTimeDemand
│   ├── queueNetwork.test.ts # MODIFIED: drop demand tests, adjust cascade test
│   └── queueApi.ts          # MODIFIED: toQueueRunRequest builds demands via autoDemand.ts
├── components/
│   └── SignalToolbar.tsx    # MODIFIED: demand form replaced by a read-only assignment list
└── SignalApp.tsx             # MODIFIED: drop demand state/handlers; "Run" gating uses autoDemand.ts
```

`queuesim/`, `cmd/server/`, `graph/`, `agent/`, `simulation/`,
`cmd/graphcli/`, `App.tsx`'s equilibrium half, `Toolbar.tsx`, and every
other equilibrium-mode file are untouched.

**Structure Decision**: Adds one pure-function module alongside
`queuePlayback.ts`/`roadVisuals.ts` in `model/`, matching this
project's established split (pure derivation logic separate from both
state-CRUD files and rendering components).

## Complexity Tracking

*No constitution violations — table not needed.*
