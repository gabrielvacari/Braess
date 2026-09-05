# Quickstart: Multiple Origins and Destinations

Validates the feature end-to-end without any UI, per constitution
Principle IV — the last simulation-engine phase before roadmap Phase 5's
frontend.

## Prerequisites

- Go 1.26+ installed (`go version`).
- Features 001-003 already in place.

## Run the automated tests

```sh
go test ./...
```

**Expected outcome**: all tests pass, including:
- Feature 003's existing `simulation` tests, unchanged, including the
  exact 65 → 80 Braess numbers — the regression check for this feature's
  `loadIncrementally`/`refine` generalization (contracts/
  simulation-api-extension.md "Compatibility").
- SC-001: in a network where two demands' shortest routes are forced
  through one shared road, increasing one demand's size measurably
  changes the other demand's own reported average travel time.
- SC-002: across 3 distinct multi-demand networks, the final assignment is
  a genuine equilibrium — no agent, from any demand, could unilaterally
  switch route and strictly improve.
- FR-007: a demand with `Size: 0` reports a trivial result without error
  and without affecting other demands.
- FR-008: a demand with no route between its origin and destination
  produces an error naming that specific demand.

## Run the standalone CLI

```sh
go run ./cmd/graphcli
```

**Expected outcome** (User Story 3 / SC-003): in addition to earlier
phases' output, the command now runs a small network where at least 3
distinct house/company pairs share a bottleneck road, and prints each
pair's own total/average travel time plus one overall total/average —
completing in well under 5 seconds — and exits with status `0`.

## What this proves

- The new code in `simulation` touches only `agent`, `graph`, and the Go
  standard library — no UI/frontend dependency (Principle I).
- A road's travel time genuinely reflects combined traffic from every
  demand using it — verifiable both by the SC-001 cross-effect test and
  by reading `simulation/demand.go`, where every agent, regardless of
  demand, reads and writes the same shared `volumes` map (Principle III
  extended to multiple commodities).
- Feature 003's behavior is provably unchanged by this feature's
  refactor, since its own tests — numeric Braess result included — pass
  without modification.

## Out of scope here (see [spec.md](./spec.md#assumptions))

Agents choosing among multiple possible destinations (every agent's
origin/destination is still fixed, just drawn from more than one declared
pair now), interactive network editing, and any frontend — those belong
to roadmap Phase 5 and its own spec.
