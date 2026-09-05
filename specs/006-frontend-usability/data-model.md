# Phase 1 Data Model: Frontend Usability and Visual Design

This feature introduces no persisted or transmitted data — no field here
is sent to `cmd/server` or stored anywhere. What follows are the shapes
of the static content and the small amount of state-dependent text
selection logic, all in `web/src/content/copy.ts` (see
[contracts/copy-api.md](./contracts/copy-api.md)).

## LegendEntry

| Field | Type | Notes |
|---|---|---|
| `type` | `NodeType` (from `model/network`) | Which node type this entry describes. |
| `label` | `string` | The short text already rendered inside that node's marker on the canvas (e.g. `"H"`) — kept in sync with `NetworkCanvas`'s `NODE_STYLE`, so the legend and the canvas never disagree about what a marker means. |
| `name` | `string` | The human-readable name (e.g. `"House"`). |
| `description` | `string` | One short sentence of what this node type is for. |

`NODE_TYPE_LEGEND: LegendEntry[]` is a fixed, hand-written array — one
entry per `NodeType`, always all three, always in the same order.

## Mode description

Not a stored entity — a pure function:

```ts
function modeDescription(mode: EditMode, connectingFrom: string | null): string
```

Given the currently active editing mode (and, for "connect" mode,
whether a first node has already been clicked), returns the sentence
describing what clicking the canvas or a node will do next. Every
`EditMode` value has exactly one corresponding description; "connect"
has two (before/after the first click), covering Edge Case: the user
must always be able to tell what the *next* click does, not just what
the mode is called.

## Empty-canvas hint

```ts
function emptyCanvasHint(nodeCount: number): string | null
```

Returns a suggested first action when `nodeCount === 0`, and `null`
(render nothing extra) otherwise.

## Missing-demand hint

```ts
function missingDemandHint(demandCount: number): string | null
```

Returns an explanation of what's missing when `demandCount === 0`, and
`null` otherwise — covers the spec's "Run pressed with no demands
declared" Edge Case.

## Relationships

```text
NODE_TYPE_LEGEND[i].type   ──> NodeType (model/network.ts, feature 005, unchanged)
NODE_TYPE_LEGEND[i].label  ──> NetworkCanvas's NODE_STYLE[type].label (kept identical by hand — see contracts/copy-api.md)
modeDescription(mode, ...) ──> EditMode (components/Toolbar.tsx, feature 005, unchanged)
```

No new types are added to `model/network.ts` or `model/api.ts` — this
feature only adds `content/copy.ts` alongside them.
