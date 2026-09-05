# Feature Specification: Frontend Usability and Visual Design

**Feature Branch**: `006-frontend-usability`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Improve the usability of the Phase 5 web frontend (specs/005-web-frontend), which currently works functionally but was built with no visual design and no explanatory content: it has zero onboarding, no legend explaining what each node marker means, no instructions describing what each editing mode does, and unstyled default browser controls with no consistent visual hierarchy. This feature adds: explanatory content (what the app demonstrates, a legend for node types, contextual guidance for whichever mode is active, and empty-state guidance when there's nothing on the canvas yet) and a real visual design (consistent spacing, color, and typography, replacing default unstyled browser controls) so the app is understandable to a first-time visitor and looks like a finished tool rather than an unstyled prototype. No new simulation capability is added — this only makes the existing Phase 5 capabilities (load & watch, draw & edit, compare before/after) understandable and presentable."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand the tool and what to do first (Priority: P1)

A person opens the app for the first time, before anything is drawn, and
immediately understands what the tool demonstrates and what to do first —
without reading any code or documentation.

**Why this priority**: This is the single biggest gap reported: a blank
canvas and unlabeled controls give a first-time visitor no idea what the
tool is for or how to begin. Nothing else in this feature matters if this
first impression doesn't work.

**Independent Test**: Can be fully tested by opening the app with an
empty network and confirming a visible explanation of the app's purpose
and a clear first action are present, with no further interaction needed.

**Acceptance Scenarios**:

1. **Given** the app has just loaded with an empty network, **When** a
   visitor looks at the screen, **Then** a short explanation of what the
   tool demonstrates is visible without scrolling or clicking anything.
2. **Given** the empty network, **When** a visitor looks at the canvas
   area, **Then** it visibly suggests a first action (e.g. load an
   example, or start placing a node) rather than showing a blank space.

---

### User Story 2 - Understand what's on screen without guessing (Priority: P2)

A person looking at a drawn network can tell, from the map alone, which
nodes are houses, which are companies, and which are plain intersections
— and while editing, can tell what the currently selected mode will do
before clicking anything.

**Why this priority**: The map's visual encoding (a shape/color per node
type) only works if there's a key explaining it, and each editing mode
changes what a click does — invisibly, without a description of what's
about to happen.

**Independent Test**: Can be fully tested by loading a network and
confirming a legend maps every visual marker to its meaning, then
selecting each editing mode in turn and confirming a description of that
mode's effect is visible.

**Acceptance Scenarios**:

1. **Given** any node is visible on the map, **When** a visitor looks for
   an explanation of its marker, **Then** a legend showing all three node
   types and their markers is visible.
2. **Given** an editing mode is selected, **When** a visitor looks at the
   controls, **Then** a plain-language description of what clicking the
   canvas or a node will do in that mode is visible.

---

### User Story 3 - The tool looks and behaves like a finished product (Priority: P3)

The controls, layout, and typography are visually consistent and
deliberate rather than default/unstyled, so the tool reads as complete
and trustworthy rather than an in-progress prototype.

**Why this priority**: Default browser controls and inconsistent spacing
undermine confidence in the tool even when the underlying computation is
sound — this is specifically what was reported as wrong, and it's the
most polish-oriented (lowest-priority) of the three stories.

**Independent Test**: Can be fully tested by visually inspecting the
app's controls, spacing, and color usage for a single, deliberate style
used consistently throughout, with no unstyled default form controls
remaining.

**Acceptance Scenarios**:

1. **Given** any screen of the app, **When** a visitor looks at buttons,
   inputs, and panels, **Then** they share a consistent visual style, not
   the browser's unstyled defaults.
2. **Given** color is used to distinguish node types or convey a result
   (better/worse, converged/not), **When** a visitor looks at it, **Then**
   the same meaning is also conveyed through a label, icon, or text, not
   color alone.

---

### Edge Cases

- What happens when a run fails (e.g. an unroutable demand)? The failure
  must be explained in plain language near where the action was taken,
  not just shown as a raw technical error message.
- What happens when the network has nodes but no demand has been declared
  yet and "Run" is pressed? This must be explained clearly (e.g., what's
  missing and how to fix it), not left as a confusing empty result.
- What happens when the explanatory content (legend, mode description,
  onboarding text) would be shown on a very small viewport? It must
  remain readable and reachable rather than being cut off or hidden.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST display a short explanation of what the
  tool demonstrates, visible as soon as the app loads.
- **FR-002**: The system MUST show a legend mapping each node type's
  visual marker to its name and meaning, visible whenever nodes can be
  shown.
- **FR-003**: The system MUST describe, for whichever editing mode is
  currently active, what clicking the canvas or an existing node/road
  will do.
- **FR-004**: The system MUST show guidance in place of an empty canvas,
  suggesting a first action, when no nodes have been placed yet.
- **FR-005**: The system MUST present every interactive control (buttons,
  inputs, selects) with a consistent, deliberate visual style rather than
  unstyled browser defaults.
- **FR-006**: The system MUST use a consistent spacing and color scheme
  throughout every screen/panel of the app.
- **FR-007**: The system MUST convey any meaning currently carried by
  color alone (node type, better/worse comparison) through an additional
  label, icon, or text as well.
- **FR-008**: The system MUST explain a run failure (e.g. an unroutable
  demand) in plain language near the action that triggered it.
- **FR-009**: The system MUST explain why "Run" cannot proceed usefully
  (or what will happen) when no demand is declared yet, rather than
  silently doing nothing or returning a confusing empty result.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A first-time visitor can state what the tool demonstrates
  and what to do first within 10 seconds of the page loading, without
  reading code or documentation.
- **SC-002**: Every node type's marker is identifiable via a visible
  legend, with zero unexplained visual symbols on screen.
- **SC-003**: For each of the three editing modes, a plain-language
  description of its effect is visible whenever that mode is active.
- **SC-004**: No unstyled default browser control (raw bullet-style
  radio/checkbox lists, unstyled buttons) remains visible anywhere in the
  app.
- **SC-005**: Every place color conveys meaning (node type, result
  comparison, error) also conveys that meaning through text or an icon,
  verifiable by inspecting the app in grayscale.

## Assumptions

- This feature only changes presentation and explanatory content in the
  existing Phase 5 frontend (`specs/005-web-frontend`) — no new
  simulation capability, no new API surface, and no change to
  `graph`/`agent`/`simulation` or `cmd/server`.
- "Explanatory content" is static, hand-written text embedded in the UI —
  not a dynamic help system, guided tour, or tooltip framework.
- The existing three editing modes (place a node, draw a road, select/
  remove) and their underlying behavior are unchanged; this feature only
  makes their effects legible before and while using them.
- Visual design work stays within the project's existing technology
  (React + Konva.js, per the constitution) — no new UI framework or
  component library is introduced.
