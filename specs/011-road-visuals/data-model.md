# Phase 1 Data Model: Top-Down Road Visuals

No new data enters or leaves the app — this feature draws differently
from data `NetworkCanvas` already receives (`nodes`, `roads`,
`roadColors`, `directed`, `selectedId`). The one new artifact is a pure
rendering-geometry helper.

## `web/src/model/roadVisuals.ts` (new)

### `RoadMark`

| Field | Type | Notes |
|---|---|---|
| `x` | `number` | Canvas x of this mark's center. |
| `y` | `number` | Canvas y of this mark's center. |
| `dx` | `number` | Unit direction vector x-component (a -> b). |
| `dy` | `number` | Unit direction vector y-component (a -> b). |

### `roadChevronPositions`

```ts
function roadChevronPositions(
  a: { x: number; y: number },
  b: { x: number; y: number },
  nodeRadius: number,
  spacing: number,
): RoadMark[];
```

Returns one `RoadMark` per directional mark, evenly spaced along the
segment strictly between `a`'s and `b`'s node radii (so no mark overlaps
either node), in order from `a` toward `b`. Returns an empty array when
the road is too short to fit even one mark between the two node edges
(spec Edge Cases: two nodes placed very close together).

## `NetworkCanvasProps` (extended — no field removed or renamed)

No new props. `roadColors` and `directed` (feature 010) are unchanged in
meaning; `directed` now additionally controls whether
`roadChevronPositions` marks are rendered, on top of continuing to
control the (now road-strip-shaped) rendering path.
