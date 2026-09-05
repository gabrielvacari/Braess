# Contract: `web/src/content/copy.ts` public API

Unlike features 001-005, this feature has no HTTP or Go contract — its
only genuinely testable surface is this module's exported functions and
constants.

## Constants

```ts
export const APP_INTRO: string; // what the tool demonstrates — FR-001
export const NODE_TYPE_LEGEND: LegendEntry[]; // exactly 3 entries, one per NodeType — FR-002
```

## Functions

```ts
/**
 * What clicking the canvas or a node will do in the given mode right
 * now. For "connect" mode, connectingFrom distinguishes "click the
 * first node" from "click the second node to connect it" (FR-003).
 */
function modeDescription(mode: EditMode, connectingFrom: string | null): string;

/**
 * A suggested first action when the canvas has no nodes yet, or null
 * once it does (FR-004).
 */
function emptyCanvasHint(nodeCount: number): string | null;

/**
 * What's missing (and why Run can't proceed usefully) when no demand is
 * declared yet, or null once at least one exists (FR-009).
 */
function missingDemandHint(demandCount: number): string | null;
```

## Invariants

- `NODE_TYPE_LEGEND` has exactly one entry per `NodeType` ("house",
  "company", "intersection") — never more, never fewer.
- `NODE_TYPE_LEGEND[i].label` for each type matches
  `NetworkCanvas`'s internal `NODE_STYLE[type].label` exactly — the two
  are maintained by hand as the same three short strings; a mismatch
  would mean the legend describes a marker the canvas doesn't actually
  draw.
- `modeDescription` returns a non-empty string for every `EditMode` value
  and both `connectingFrom` states — there is no mode with no
  description.
- `emptyCanvasHint`/`missingDemandHint` return `null` exactly when their
  condition (`nodeCount > 0` / `demandCount > 0`) holds — never a message
  once the condition that motivated it is no longer true.

## Compatibility

`web/src/model/network.ts` and `web/src/model/api.ts` (feature 005) are
read-only dependencies of this module (for the `NodeType`/`EditMode`
types) — this contract introduces no changes to either.
