# Feature Specification: Bidirectional Roads

**Feature Branch**: `007-bidirectional-roads`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Eu preciso que a ligação entre dois nós sejam bilaterias ou seja, A -> B -> A. E que no frontend eu precise apenas fazer uma conexão." (I need the connection between two nodes to be bidirectional — A -> B and B -> A. And in the frontend I should only need to make one connection.)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Drawing one road connects both directions (Priority: P1)

A user connects two nodes with a single action, and afterward agents can
travel between them in either direction — matching how a real road works.

**Why this priority**: This is the direct cause of a confusing "no route
exists" error a user just hit: they drew what they understood to be "a
road" between two nodes, then declared a trip in the direction they
didn't happen to draw it, and got an error that looked like a bug. This
is the core problem the feature exists to fix.

**Independent Test**: Can be fully tested by drawing one road between two
nodes, then declaring a demand in each direction between them, and
confirming both find a route — no server or engine change involved.

**Acceptance Scenarios**:

1. **Given** two existing, unconnected nodes, **When** the user draws a
   road between them, **Then** agents can be routed from the first to the
   second.
2. **Given** the same drawn road, **When** a demand is declared from the
   second node to the first, **Then** a route is found — the direction
   the road happened to be drawn in does not matter.

---

### User Story 2 - Removing a road removes both directions together (Priority: P2)

Removing a drawn road removes the connection completely, in both
directions — never leaving a one-way remnant a user didn't know existed.

**Why this priority**: A user thinks of what they drew as "one road," not
two separate one-way connections — removal has to match that mental
model, or the network could silently keep behaving as if a road were
still there.

**Independent Test**: Can be fully tested by drawing a road, removing it,
and confirming neither direction between those two nodes routes
successfully anymore (assuming no other road connects them).

**Acceptance Scenarios**:

1. **Given** a road drawn between two nodes, **When** the user removes
   it, **Then** no route exists between those two nodes in either
   direction (unless some other road independently connects them).

---

### Edge Cases

- What happens when two nodes already have one or more roads drawn
  between them and the user draws another? This must still be allowed
  (parallel roads, supported since the underlying engine's first phase)
  — the new road is its own independent, still-bidirectional connection,
  not merged with the existing one(s).
- What happens to tools that build one-directional roads directly,
  outside this web frontend (e.g. the project's existing terminal
  command)? They are unaffected — this feature only changes what
  drawing one road in the web frontend produces.
- What happens when the user tries to connect a node to itself? Unchanged
  from existing behavior — still rejected, regardless of directionality.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST create a connection usable in both
  directions between two nodes when the user performs one "draw road"
  action.
- **FR-002**: The system MUST allow an agent to be routed in either
  direction between two nodes connected by a drawn road.
- **FR-003**: The system MUST apply the same road-behavior configuration
  (chosen at the time of drawing) to both directions of a drawn road —
  a single choice covers the whole connection.
- **FR-004**: The system MUST display a drawn road as a single visual
  connection between the two nodes, not as two separate overlapping
  lines.
- **FR-005**: The system MUST remove both directions of a drawn road
  together when the user removes it, leaving no one-way remnant.
- **FR-006**: The system MUST continue to allow multiple, independent
  roads to be drawn between the same two nodes, each one still
  bidirectional on its own.
- **FR-007**: This change MUST NOT alter how any other part of the
  project (its simulation engine or its terminal command) builds or uses
  one-directional roads directly — it only changes what one "draw road"
  action in the web frontend produces.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: After drawing one road between two nodes, a demand declared
  in either direction between them finds a route, verified across
  multiple draw/demand-direction combinations.
- **SC-002**: After removing a drawn road, no route remains possible in
  either direction between those two nodes (assuming no other road
  independently connects them).
- **SC-003**: The map shows exactly one visual line per drawn road,
  regardless of it representing a two-way connection underneath.
- **SC-004**: Drawing multiple roads between the same two nodes still
  takes one user action per road, each independently usable in both
  directions.

## Assumptions

- Both directions of a drawn road always share identical behavior (the
  same road-behavior configuration) — this feature does not add a way to
  give each direction of one drawn road different behavior.
- This is a web-frontend-only change: the underlying simulation engine
  and any tool that already builds one-directional roads directly (e.g.
  the project's terminal command) are unaffected and continue to support
  one-directional roads in their own contexts — see FR-007.
- The existing "draw road" interaction (click a starting node, then an
  ending node) is unchanged; this feature changes what that interaction
  produces, not how a user performs it.
