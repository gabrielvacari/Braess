# Feature Specification: Top-Down Road Visuals

**Feature Branch**: `011-road-visuals`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Redesign how roads render on the network canvas (both Equilibrium and Signals modes) so a road looks like an actual road viewed from above, instead of a thin line or a single arrowhead — a visibly road-like strip (asphalt-colored, thicker) with a dashed center line, rounded at the ends where it meets a node. For a one-way road (Signals mode), direction must still be visible, but shown as repeated small directional chevron marks along the road's surface (like real pavement arrow markings) instead of one arrowhead at the end. This is purely a rendering/visual change to the shared NetworkCanvas component — no change to the simulation engine, no new data, no change to how roads are drawn/edited/selected, and the existing signal-state color override (green/red) and selection highlight must still work, now applied to the road-like visual instead of a plain line."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See roads as roads, not lines (Priority: P1)

A person looking at the canvas, in either mode, sees each connection
between two nodes rendered as a recognizable road — a paved-looking
strip with a center line — instead of a plain thin line, so the whole
map reads like a road network at a glance.

**Why this priority**: This is the core of the request — the plain-line
rendering was explicitly reported as not visual enough; everything else
(direction markings, keeping selection/signal-state visible) only
matters once the base road visual itself exists.

**Independent Test**: Draw a road between two nodes in either mode and
confirm it renders as a road-like strip with a center line, not a thin
single-color line.

**Acceptance Scenarios**:

1. **Given** a road between two nodes, **When** it's rendered on the
   canvas, **Then** it appears as a strip wide enough to read as a road
   surface, with a dashed center line running along it.
2. **Given** that same road, **When** it meets a node at either end,
   **Then** the connection looks smooth/rounded rather than a hard,
   abrupt cut.

---

### User Story 2 - See a one-way road's direction as pavement markings (Priority: P2)

A person looking at a one-way road (Signals mode) sees its direction of
travel shown as repeated small directional markings along the road's
surface — like the arrows painted on a real one-way street — instead of
a single arrowhead at one end.

**Why this priority**: Directly replaces the single-arrowhead treatment
this project shipped most recently, per explicit feedback that it should
look like a real road rather than an arrow; depends on User Story 1's
road-strip visual already existing to paint the markings onto.

**Independent Test**: Draw a one-way road in Signals mode and confirm
its direction is shown via repeated markings along its length, visible
regardless of which part of the road is currently in view, not a single
arrowhead at one end.

**Acceptance Scenarios**:

1. **Given** a one-way road, **When** it's rendered, **Then** its
   direction of travel is legible from markings distributed along its
   length rather than concentrated at one end.
2. **Given** a bidirectional road (Equilibrium mode), **When** it's
   rendered, **Then** it shows no directional marking at all, since it
   isn't one-way.

---

### User Story 3 - Keep selection and signal state visible on the new look (Priority: P3)

A person can still tell whether a road is currently selected, and (in
Signals mode) whether its signal is green or red, at the same glance as
before — the more visual road treatment doesn't bury information that
already worked.

**Why this priority**: Preserves already-delivered value; lower
priority because it's a preservation check on Stories 1 and 2 rather
than new capability, but a regression here would be immediately
noticeable.

**Independent Test**: Select a road and confirm it's still visibly
distinguished; in Signals mode, run a scenario and confirm a
signal-controlled road's green/red state is still clearly visible on
the new road visual as playback advances.

**Acceptance Scenarios**:

1. **Given** a road is selected, **When** it's rendered, **Then** it is
   visibly distinguished from an unselected road.
2. **Given** a signal-controlled road during Signals-mode playback,
   **When** its signal is green or red, **Then** that state is clearly
   visible on the road itself.

---

### Edge Cases

- Two nodes placed very close together: the road's directional markings
  and center-line dashes must not visually overlap into an illegible
  smear.
- A network with many roads close together: the thicker road strips must
  not make the canvas harder to read than the previous thin lines did.
- A one-way road with no signal attached (Signals mode): it still shows
  its directional markings, in a neutral (non-green/red) state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The canvas MUST render every road as a road-like strip —
  visibly wider than a thin line, with a distinct surface appearance —
  rather than a single-color line.
- **FR-002**: Every road's strip MUST include a center-line marking
  distinguishing the road's surface from its edges.
- **FR-003**: A road MUST meet the node at each of its ends with a
  smooth, rounded transition rather than a hard-edged cut.
- **FR-004**: A one-way road (Signals mode) MUST show its direction of
  travel via repeated directional markings distributed along its length.
- **FR-005**: A bidirectional road (Equilibrium mode) MUST NOT show any
  directional marking.
- **FR-006**: The existing selection highlight for a road MUST remain
  clearly visible on the new road visual.
- **FR-007**: The existing signal-state color indication (green/red,
  Signals mode) MUST remain clearly visible on the new road visual.
- **FR-008**: This change MUST NOT alter how a road is drawn, selected,
  or removed — only how it appears once rendered.
- **FR-009**: This change MUST NOT alter any simulation engine behavior
  or data — it is a rendering-only change.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A first-time viewer can identify which shapes on the
  canvas are roads, versus other elements, purely from their
  road-like appearance, without consulting a legend.
- **SC-002**: A person can state a one-way road's direction of travel
  just from looking at it, without needing to read a separate
  explanation beyond what already exists for node types.
- **SC-003**: A person can still identify a selected road, and (in
  Signals mode) a road's current green/red signal state, as reliably as
  before this change.
- **SC-004**: Every existing automated test covering road drawing,
  selection, and removal continues to pass unmodified after this
  change.

## Assumptions

- This is a rendering change to the shared canvas component only
  (feature 009's precedent: one canvas, reused by both modes) — no new
  data is needed from the server or the engine, and no existing data
  shape changes.
- Exact colors (asphalt strip, center-line, directional markings) and
  the spacing/size of directional markings are visual-design decisions
  made during implementation, not user-configurable settings.
- No motion is added to the road surface itself (e.g. animated dashes);
  only agents (an existing, separate capability) move on top of it.
- The previous single-arrowhead treatment for one-way roads (introduced
  immediately prior to this feature) is fully replaced by the
  along-the-road marking style described here, not kept as an
  alternative.
