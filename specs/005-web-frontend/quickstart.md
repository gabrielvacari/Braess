# Quickstart: Web Frontend

This is the first feature with a genuinely visual payoff — some of it is
still verified in text (per constitution Principle IV, as far as it can
be), and the rest is verified by eye, walking through spec.md's own
Acceptance Scenarios.

## Prerequisites

- Go 1.26+ (`go version`).
- Node 20+ and npm (`node --version`).
- Features 001-004 already in place and passing (`go test ./...`).

## Run the automated tests

```sh
go test ./...              # unchanged packages + the new cmd/server contract tests
cd web && npm install && npm test   # Vitest: network.ts edit operations, api.ts mapping
```

**Expected outcome**: all Go tests still pass (features 001-004
untouched); the new `cmd/server` tests cover a valid run, a self-loop
rejection, and a no-route demand producing a `422` naming the demand;
`web/`'s Vitest suite covers `addNode`/`addEdge`/`removeNode`/`removeEdge`
including the self-loop rejection and cascade-delete rules (FR-005,
FR-008).

## Run it for real

```sh
# Terminal 1 — the API (defaults to :8090; override with PORT=... if that's taken)
go run ./cmd/server

# Terminal 2 — the web app
cd web && npm install && npm run dev
```

Open the printed local URL (Vite's default is `http://localhost:5173`). The
web app calls the API at `http://localhost:8090` by default — override with
a `VITE_API_BASE` env var (or a `web/.env.local`) if you changed `PORT`.

## Walk through the Acceptance Scenarios by hand

1. **User Story 1** (load and watch): use the app's "load example" action
   (or draw the small network below) and press Run — every node should be
   visually distinguishable by type, every road drawn, and agent icons
   should visibly travel from origin to destination.
2. **User Story 2** (draw and edit): starting from an empty canvas, place
   two nodes, mark one a house and one a company, draw a road between
   them, then delete that road and one of the nodes — confirm each step
   is reflected immediately.
3. **User Story 3** (before/after): draw a small network with two routes
   between one house and one company (e.g. mirroring feature 003's
   classic Braess shape), run it, note the average, add one more road,
   run again, and confirm both averages remain visible together (SC-004)
   — this is the interactive version of Phase 3's 65 → 80 result.

## What this proves

- `cmd/server` is the only new Go code, and it depends on `graph`/`agent`/
  `simulation` without changing any of them — verifiable by `git diff`
  touching no files under those three packages (Principle I).
- Every number the UI displays came from one `POST /api/run` call into
  the unchanged, already-validated engine — the browser never computes a
  route or a travel time itself (Principle III, research.md decision #3).
- This is the first feature allowed to add UI at all, per Principle II —
  and the fact that it's just wiring a browser onto four already-proven
  Go packages, rather than re-deriving any simulation behavior, is the
  payoff of having done Phases 1-4 in the required order.

## Out of scope here (see [spec.md](./spec.md#assumptions))

Saving a network across browser sessions, multiple simultaneous users or
tabs editing the same network, and live tick-by-tick server-pushed
simulation (as opposed to animating an already-computed result) — all
explicitly deferred.
