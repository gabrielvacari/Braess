# Quickstart: Shortest-Path Agents

Validates the feature end-to-end without any UI, per constitution
Principle IV. Assumes feature 001 (`graph`) is already in place.

## Prerequisites

- Go 1.26+ installed (`go version`).

## Run the automated tests

```sh
go test ./...
```

**Expected outcome**: all tests pass, including (per Success Criteria):
- SC-001: across several small test graphs, `ShortestRoute` returns the
  mathematically correct minimum-travel-time route.
- SC-003: three independently-created `Agent`s sharing an origin/
  destination each compute the same route via their own `ComputeRoute`
  call.
- SC-004: a graph with no path between the given origin and destination
  makes `ShortestRoute` return an error satisfying
  `errors.Is(err, agent.ErrNoRoute)`.
- FR-006: origin == destination returns a zero-edge, zero-time route.
- FR-008: an edge whose travel-time evaluation errors surfaces that error
  from `ShortestRoute`, rather than being silently skipped.

## Run the standalone CLI

```sh
go run ./cmd/graphcli
```

**Expected outcome** (User Story 3 / SC-002): in addition to feature 001's
node/edge listing, the command now also prints, for a fixed
origin/destination pair on the sample network, each simulated agent's
route (its ordered edges) and total travel time — completing in well under
5 seconds — and exits with status `0`.

## What this proves

- `agent` imports only `graph` and the Go standard library — no UI/
  frontend dependency exists to even accidentally import (Principle I).
- Each printed agent's route came from that agent's own `ComputeRoute`
  call — verifiable by reading `cmd/graphcli`'s source, which constructs
  several independent `Agent` values and calls `ComputeRoute` on each,
  with no shared route cache anywhere in the call path (Principle III,
  FR-003).
- Route correctness is verified in text (via `go test`) before any
  frontend work starts (Principle IV).

## Out of scope here (see [spec.md](./spec.md#assumptions))

No congestion (travel time still evaluated at a fixed baseline volume), no
multiple distinct origin/destination pairs, no persistence, no frontend —
those belong to later roadmap phases and their own specs.
