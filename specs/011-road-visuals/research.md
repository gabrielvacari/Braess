# Phase 0 Research: Top-Down Road Visuals

## Decision 1: A thick `Line` with rounded caps, not a custom polygon shape, for the road strip

**Decision**: Render the road's base strip as a Konva `Line` with a
larger `strokeWidth` and `lineCap="round"`, rather than building a
custom rectangle/polygon shape between the two node positions.

**Rationale**: Konva's `lineCap="round"` already produces the "rounded
where it meets a node" look FR-003 asks for, for free — no geometry
beyond what `NetworkCanvas` already computes (`a`, `b`). A hand-built
rounded rectangle would need to compute perpendicular offsets and arc
segments for a purely cosmetic gain the built-in cap already provides.

**Alternatives considered**: A custom `Shape` with manual path
construction — rejected as unnecessary complexity for the same visual
result.

## Decision 2: A second, non-interactive `Line` for the dashed center line

**Decision**: Draw the center-line marking as a second `Line` with the
same two endpoints, a fixed light color, `dash={[10, 8]}`, and
`listening={false}`, painted directly on top of the base strip.

**Rationale**: Konva's built-in `dash` prop on `Line` is exactly a
dashed center line; a second non-listening line avoids any change to
click/selection behavior (the base strip alone still owns
`onClick`/`hitStrokeWidth`), satisfying FR-008 (no change to how a road
is selected).

**Alternatives considered**: Encoding the dash directly on the base
strip itself — rejected because the base strip's color must still
switch for selection/signal-state (FR-006/FR-007), and a dashed stroke
color change would misleadingly recolor the lane marking too, not just
the road surface.

## Decision 3: Directional marks are small `Arrow` shapes, not rotated polygons

**Decision**: Each directional mark along a one-way road is a short
Konva `Arrow` (a tiny two-point line with `pointerLength`/`pointerWidth`
set small), oriented along the road's own direction vector — the same
primitive already used (and already visually verified live) for the
single end-arrowhead this feature replaces.

**Rationale**: `Arrow` already points correctly along whatever two
points it's given; reusing it removes any need to work out rotation
angles for a `RegularPolygon`-based triangle by hand, and reuses a
component this project already confirmed renders correctly.

**Alternatives considered**: A rotated `RegularPolygon` (3 sides) at
each mark position — rejected: correct rotation would need deriving the
polygon's own default vertex orientation from scratch, more room for a
visually-wrong result with no way to check it without a browser.

## Decision 4: Mark placement is a pure function, `roadChevronPositions`

**Decision**: A pure function in `web/src/model/roadVisuals.ts` computes
every directional mark's position and direction vector along a road,
given the two endpoints, the node radius (marks must not overlap a
node), and a fixed spacing — returning an empty array when the road is
too short to fit any mark between its two nodes' edges (spec Edge
Cases).

**Rationale**: Matches the project's established pattern (`queueChart.ts`,
`queuePlayback.ts`) of extracting anything computable into a pure,
independently-tested function rather than embedding the math inside JSX
— the one place in this feature with real logic (as opposed to prop
wiring) worth a unit test (constitution Principle IV).

**Alternatives considered**: Computing mark positions inline inside
`NetworkCanvas`'s render — rejected: the same math with no test coverage
and no reuse boundary.

## Decision 5: The dashed center line renders on every road, directed or not

**Decision**: FR-002 applies uniformly — both bidirectional
(Equilibrium) and one-way (Signals) roads get the base strip + dashed
center line; only the extra directional marks are conditional on
`directed`.

**Rationale**: The feature description explicitly describes the marks
as replacing the prior single arrowhead *in addition to* the road-strip
look, not a different visual for one-way roads altogether — keeping the
base look identical between modes is also simpler and more consistent
with `NetworkCanvas` staying one shared component (feature 009
research.md decision #3).

**Alternatives considered**: Omitting the center line for one-way roads
(reasoning that a single-lane one-way street has no need for a lane
divider) — rejected as contradicting the literal feature description and
adding a mode-conditional branch to FR-002 that the spec doesn't call
for.
