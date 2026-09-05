# Quickstart: Top-Down Road Visuals

## Prerequisites

- Node 20+.
- Features 005-010 already in place.

## Run the automated tests

```sh
cd web && npm test
```

**Expected outcome**: all existing suites pass unmodified; new coverage
includes `web/src/model/roadVisuals.test.ts`: `roadChevronPositions`
places marks evenly between the two nodes' edges, points them a -> b,
and returns an empty array for a road too short to fit one.

## Run it for real

```sh
cd web && npm run dev
```

(The Go server isn't needed for this feature — no request/response
shape changes — but run it too if you want to actually execute a
scenario.)

1. In Equilibrium mode, draw a road between two nodes: it now renders
   as a road-like strip with a dashed center line and rounded ends,
   instead of a thin line. No directional marks appear (it's
   bidirectional).
2. Switch to Signals mode, draw a one-way road: it shows the same strip
   and center line, plus small directional arrow marks repeated along
   its length showing which way it runs — replacing the single
   end-arrowhead from feature 010.
3. Select a road in either mode: confirm it's still clearly highlighted.
4. In Signals mode, run a scenario: confirm a signal-controlled road's
   green/red state is still clearly visible on the (now thicker) strip.

## What this proves

- The change is confined to `web/src/components/NetworkCanvas.tsx` and
  a new `web/src/model/roadVisuals.ts` — `git diff --stat` shows no
  other file touched (constitution Principle I: not even adjacent to
  the engine).
- Selection and signal-state indication (feature 009/010) still work,
  proven by the existing tests for `App.tsx`/`SignalApp.tsx`'s edit and
  run flows passing unmodified.

## Out of scope here (see [spec.md](./spec.md#assumptions))

Configurable colors or marking spacing, and any animation of the road
surface itself — both explicitly deferred as implementation-only visual
choices.
