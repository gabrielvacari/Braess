# Quickstart: Congestion and Braess's Paradox

Validates the feature end-to-end without any UI, per constitution
Principle IV — and this is the roadmap's explicit pre-frontend gate
(AGENTS.md).

## Prerequisites

- Go 1.26+ installed (`go version`).
- Features 001 (`graph`) and 002 (`agent`) already in place.

## Run the automated tests

```sh
go test ./...
```

**Expected outcome**: all tests pass, including:
- Feature 002's existing `agent` tests, unchanged (the `agent` extension
  in this feature is additive-only — see contracts/agent-api-extension.md).
- SC-001: across 3 distinct small networks with a known analytical
  equilibrium split, `simulation.Run`'s final assignment matches it within
  a small tolerance.
- SC-004: with `Population.MaxRounds == 0`, `Run` reports
  `Converged: false`.
- FR-010: `Population.Size == 0` returns a zero-cost result without error.
- The classic Braess network (`braess_test.go`): running the same
  4000-agent population with and without the extra road shows the
  with-road total/average travel time is strictly worse — this is the
  automated backing for SC-002, in addition to seeing it in the CLI below.

## Run the standalone CLI

```sh
go run ./cmd/graphcli
```

**Expected outcome** (User Story 3 / SC-002, SC-003): in addition to
earlier phases' output, the command now runs the classic four-node Braess
network twice — once without its extra connecting road, once with it —
and prints each run's total and average travel time, making explicit that
adding the road made the average trip *slower* for everyone (the
well-known 65 → 80 result). Completes in well under 5 seconds and exits
with status `0`.

## What this proves

- `simulation` imports only `agent`, `graph`, and the Go standard library
  — no UI/frontend dependency exists to even accidentally import
  (Principle I).
- Every agent's reconsideration during refinement is its own independent
  `agent.ShortestRouteAtVolumes` call against currently-observable
  volumes — verifiable by reading `simulation/run.go`, which holds no
  central "assign everyone optimally" step anywhere (Principle III).
- Braess's Paradox — adding a road can make a selfishly-routed network
  worse for everyone — is demonstrated in text, matching the literature's
  own numbers, before any frontend code exists (Principle II, IV). Per
  AGENTS.md's golden rule, this is what unlocks starting roadmap Phase 5
  (the web frontend).

## Out of scope here (see [spec.md](./spec.md#assumptions))

Multiple distinct origin/destination pairs (still one fixed pair, as in
Phase 2), a configurable/interactively-edited network (the Braess example
is fixed and hand-built), and any frontend — those belong to later roadmap
phases and their own specs.
