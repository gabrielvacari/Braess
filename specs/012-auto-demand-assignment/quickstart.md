# Quickstart: Automatic House Demand

## Prerequisites

- Node 20+.
- Features 009-011 already in place.

## Run the automated tests

```sh
cd web && npm test
```

**Expected outcome**: all existing suites pass, minus the removed
demand-CRUD tests; new coverage includes `web/src/model/autoDemand.test.ts`:
- No companies -> no assignments/demands.
- One company -> every house assigned to it.
- Two or more houses and two or more companies -> at least two distinct
  companies appear among the assignments.
- Calling the function twice with the same nodes produces the same
  result (FR-006).

## Run it for real

```sh
# Terminal 1
go run ./cmd/server

# Terminal 2
cd web && npm run dev
```

1. Switch to Signals mode. Notice there is no demand form anymore.
2. Draw two houses and two companies, with roads connecting each house
   to at least one company. The toolbar shows which company each house
   has been assigned to (FR-010).
3. Press Run — confirm agents from both houses appear without ever
   having declared a demand.
4. Confirm the two houses' traffic reaches two different companies
   (SC-003), and that pressing Run again shows the same assignment
   (SC-004).

## What this proves

- Signals mode never asks a person to declare who goes where — traffic
  exists purely because houses and companies exist (SC-001).
- `queuesim` and `cmd/server` are untouched: `git diff --stat` shows no
  change under `queuesim/`, `cmd/server/`, or any equilibrium-mode file.

## Out of scope here (see [spec.md](./spec.md#assumptions))

Configurable rate or per-house overrides, and any change to the
Equilibrium mode's own demand form — both explicitly out of scope.
