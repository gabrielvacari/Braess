# Quickstart: Graph Simulation Engine

Validates the feature end-to-end without any UI, per constitution
Principle IV.

## Prerequisites

- Go 1.26+ installed (`go version`).

## Run the automated tests

```sh
go test ./...
```

**Expected outcome**: all tests pass, including (per Success Criteria):
- SC-002: an edge's travel time at `volume == 0` equals its configured
  free-flow time.
- SC-003: across at least 5 increasing volume samples on the same edge, the
  computed travel time never decreases.
- FR-003 / FR-008 error cases: adding an edge with an unknown node ID, and
  computing travel time with a negative volume, both return an error.

## Run the standalone CLI

```sh
go run ./cmd/graphcli
```

**Expected outcome** (User Story 3 / SC-001): within seconds, the command
builds a small sample network (at least 4 nodes and 4 edges, mixing node
types and travel-time functions per [research.md](./research.md#2-travel-time-function-representation))
and prints, to stdout:
- every node's ID and type
- every edge's ID, `from`/`to` endpoints, length, capacity, and current
  travel time at a sample volume

exiting with status `0`.

## What this proves

- The engine builds and runs with `cmd/graphcli` importing only the `graph`
  package and the Go standard library — no UI/frontend dependency exists
  to even accidentally import (Principle I, SC-004).
- The graph's structure and its travel-time behavior are both verifiable
  from text output alone, before any frontend work starts (Principle II,
  Principle IV).

## Out of scope here (see [spec.md](./spec.md#assumptions))

No agents, no live/aggregated traffic volume, no persistence, no
frontend — those belong to later roadmap phases and their own specs.
