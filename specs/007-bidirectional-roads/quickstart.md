# Quickstart: Bidirectional Roads

Presentation/model-only feature — no Go changes, so this only extends
feature 005/006's `web/` validation path.

## Prerequisites

- Node 20+ and npm.
- Features 005-006 already in place (`cd web && npm test` passing).
- The Go API running (`go run ./cmd/server`) for the manual walkthrough.

## Run the automated tests

```sh
cd web && npm test
```

**Expected outcome**: `network.test.ts` covers, renamed around roads:

- `addRoad` adds one road; rejects a self-loop.
- `removeRoad` removes only the targeted road.
- `removeNode` cascades to every road (and demand) referencing it.
- `expandRoadsToDirectedEdges`: exactly two directed edges per road, one
  each direction, both carrying the road's `travelTime`; deterministic
  ids across repeated calls on the same input.

## Walk through the Acceptance Scenarios by hand

1. **User Story 1**: with the API running, draw one road between two
   nodes. Declare a demand from the *second* node to the *first* (the
   direction you didn't draw it in) and press Run — confirm it succeeds
   (no "no route exists" error), reproducing the exact scenario that
   motivated this feature.
2. **User Story 2**: select and remove that road. Declare the same
   demand again and press Run — confirm it now fails with a clear
   "no route exists" error (assuming no other road connects those two
   nodes), proving removal took both directions with it.
3. Draw a second, parallel road between the same two nodes (spec Edge
   Cases) — confirm it's accepted and the map still shows two lines (one
   per road), each still usable in both directions.

## What this proves

- The map shows one line per drawn road (SC-003), and both directions
  work or disappear together (SC-001, SC-002) — verifiable by eye per the
  walkthrough above.
- `git diff --stat` touches only files under `web/` — no change to
  `graph/`, `agent/`, `simulation/`, or any `cmd/` package (Principle I,
  FR-007).

## Out of scope here (see [spec.md](./spec.md#assumptions))

Per-direction road configuration, and any change to how the "draw road"
interaction itself is performed — both explicitly ruled out.
