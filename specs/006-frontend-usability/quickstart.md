# Quickstart: Frontend Usability and Visual Design

Presentation-only feature — no Go changes, so this only touches `web/`'s
own validation path from feature 005's quickstart.

## Prerequisites

- Node 20+ and npm.
- Feature 005 already in place (`cd web && npm test` passing).

## Run the automated tests

```sh
cd web && npm test
```

**Expected outcome**: feature 005's existing suites (`network.test.ts`,
`useAgentAnimation.test.ts`, `runHistory.test.ts`) still pass unmodified,
plus a new `content/copy.test.ts` covering:

- `modeDescription` returns a distinct, non-empty description for each of
  the three modes, and a different one for "connect" before vs. after the
  first node is clicked.
- `emptyCanvasHint(0)` returns a message; `emptyCanvasHint(1)` returns
  `null`.
- `missingDemandHint(0)` returns a message; `missingDemandHint(1)` returns
  `null`.

## Build it

```sh
npm run build
```

**Expected outcome**: builds clean, same as feature 005.

## Walk through the Acceptance Scenarios by hand

1. **User Story 1**: open the app fresh (nothing loaded). Confirm, without
   clicking anything, you can see (a) a sentence explaining what the tool
   demonstrates and (b) guidance on the empty canvas suggesting a first
   action.
2. **User Story 2**: click "Load example." Confirm a legend is visible
   mapping each of the three node markers to a name/meaning. Switch
   between "Place node," "Draw road," and "Select / remove" and confirm
   the visible description text changes to describe that mode — and, in
   "Draw road" mode, changes again after clicking one node (first-click
   vs. second-click wording).
3. **User Story 3**: look over every control (buttons, the mode picker,
   the demand form, the node-type picker). Confirm they share one visual
   style — no raw bulleted radio lists, no unstyled default buttons.
   Then view the page in grayscale (e.g. browser DevTools' "Emulate
   vision deficiencies: Achromatopsia," or a screenshot desaturated
   afterward) and confirm: node types are still distinguishable (by
   shape/label), and a before/after comparison's meaning is still readable
   (by its wording), without relying on red vs. green.
4. **Edge cases**: trigger a run with no demand declared — confirm a
   plain-language explanation appears near the Run button. Trigger a run
   with an unroutable demand — confirm the resulting error is explained
   near where Run was pressed, not just a raw server message.

## What this proves

- `content/copy.ts` is the single place every explanatory string and its
  selection logic lives, and it's exercised by `go`-style unit tests the
  same way `model/network.ts` was in feature 005 (Principle IV, as far as
  text logic — not pixels — can be).
- No file under `graph/`, `agent/`, `simulation/`, or `cmd/server`
  changed — verifiable via `git diff --stat` touching only `web/`
  (Principle I).

## Out of scope here (see [spec.md](./spec.md#assumptions))

A dismissible banner or guided tour, a new UI framework or component
library, and any change to the three editing modes' underlying behavior
— all explicitly deferred or ruled out.
