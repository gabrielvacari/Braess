# Feature Specification: Web Frontend

**Feature Branch**: `005-web-frontend`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Phase 5 of the roadmap (per AGENTS.md and constitution Principle II, now satisfied by Phases 1-4): a browser-based visualization and editor for the road network validated underneath. Draw the map interactively, mark nodes as houses or companies, dynamically add/remove roads, and watch agents move in real time — built on the existing Go simulation engine (graph/agent/simulation), kept completely separate from it: the engine code does not move, a new server component exposes it to the browser, and the browser app lives in its own directory alongside the engine, per the confirmed backend/frontend architecture decision."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Watch the simulation happen on a map (Priority: P1)

A user opens the web app and sees the road network drawn as a map —
houses, companies, intersections, and roads — with cars visibly moving
along their computed routes over time, matching what the console/log
validation from earlier phases already proved.

**Why this priority**: This is the most fundamental payoff of the whole
project restated visually — the first proof that what was validated in
text also holds up as something a person can actually watch happen.
Nothing else in this phase matters if this doesn't work.

**Independent Test**: Can be fully tested by loading the page against an
already-computed network and population and confirming every node/edge is
drawn and every agent's icon visibly progresses along its route, reaching
its destination — no editing required yet.

**Acceptance Scenarios**:

1. **Given** a network and a population already computed by the engine,
   **When** the page loads, **Then** every node is drawn with a visual
   distinction between house, company, and intersection, and every road
   is drawn connecting its two nodes.
2. **Given** the same loaded network, **When** the simulation is playing,
   **Then** each agent's icon visibly moves along its assigned route's
   roads over time, arriving at its destination.

---

### User Story 2 - Draw and edit the network by hand (Priority: P2)

A user adds and removes nodes and roads directly on the map, and marks
any node as a house, a company, or a plain intersection.

**Why this priority**: This is the interactivity the project has always
been building toward — turning the tool from a fixed viewer into
something a person can actually experiment with, per the project's
original goal of drawing roads and marking houses/companies.

**Independent Test**: Can be fully tested by, starting from an empty or
preset map, placing a few nodes, marking their types, connecting them
with roads, and confirming the resulting network matches what was drawn —
independent of whether a simulation has run yet.

**Acceptance Scenarios**:

1. **Given** an empty map, **When** the user places a node and marks it
   as a house, **Then** that node exists in the network as a house.
2. **Given** two existing nodes, **When** the user draws a road between
   them, **Then** that road exists in the network connecting those two
   nodes.
3. **Given** an existing road or node, **When** the user removes it,
   **Then** it no longer exists in the network.

---

### User Story 3 - See how adding a road changes congestion, live (Priority: P3)

After editing the network — most notably, adding an extra road between
two points already connected another way — a user re-runs the population
and immediately sees, on the map and in the reported numbers, whether
that change made the average trip better or worse.

**Why this priority**: This is the interactive version of the project's
whole reason for existing — turning the fixed Braess's Paradox
demonstration from Phase 3 into something a person can reproduce (or fail
to reproduce, just as informatively) with their own hand-drawn networks.

**Independent Test**: Can be fully tested by drawing a small network,
running the population, noting the reported average travel time, adding
one more road between two already-connected points, re-running, and
confirming the new average is shown alongside the previous one.

**Acceptance Scenarios**:

1. **Given** a network with a population already run once, **When** the
   user adds a road and re-runs, **Then** the map and the reported
   total/average travel time reflect the new network, not the old one.
2. **Given** a before/after pair of runs on the same population, **When**
   both have been run, **Then** the user can see both results together
   well enough to tell whether the network change helped or hurt.

---

### Edge Cases

- What happens when the user tries to draw a road from a node to itself?
  This must be rejected — it doesn't correspond to any real road.
- What happens when the user tries to draw a second road between two
  nodes that already have one? This must be allowed, not blocked —
  parallel roads are exactly what adding a Braess's-Paradox-style
  shortcut looks like, and the engine has supported them since Phase 1.
- What happens when the user removes a node that still has roads attached
  to it? Those roads must be removed along with it — never left dangling,
  referencing a node that no longer exists.
- What happens when the user requests a run whose declared origin or
  destination no longer exists in the current network (e.g., removed
  after being set)? This must be reported clearly, not crash — the same
  no-route philosophy already established for a missing route in earlier
  phases, now surfaced visually.
- What happens if the user edits the network while a previous run's
  animation is still playing? The in-progress animation must stop or
  pause rather than visually corrupt itself; a fresh run must be
  explicitly requested afterward to see the edited network's behavior.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST render a road network's nodes and roads
  visually, with houses, companies, and intersections visually
  distinguishable from one another.
- **FR-002**: The system MUST animate each agent traveling its computed
  route, showing visible progress along the network over time, from its
  origin to its destination.
- **FR-003**: The system MUST let a user add a node to the network from
  within the browser, choosing its type (house, company, or
  intersection).
- **FR-004**: The system MUST let a user add a road between two existing
  nodes from within the browser.
- **FR-005**: The system MUST let a user remove a node or a road from the
  network from within the browser, and MUST NOT leave a road referencing
  a node that no longer exists.
- **FR-006**: The system MUST let a user run the population (agents
  traveling between declared origin/destination pairs) against the
  currently drawn network from within the browser, and see the
  resulting total and average travel time.
- **FR-007**: The system MUST let a user compare the results of two runs
  on the same population — e.g. before and after a network edit — well
  enough to tell whether the change helped or hurt overall travel time.
- **FR-008**: The system MUST reject drawing a road from a node to
  itself.
- **FR-009**: The system MUST clearly report, rather than crash, when a
  run is requested for an origin or destination that no longer exists in
  the current network.
- **FR-010**: The system MUST keep the browser responsive while a run is
  being computed, rather than appearing frozen, for networks and
  populations within this project's already-established performance
  scale (Phase 3/4's own success criteria).

### Key Entities

- **Network (client view)**: The browser's picture of the road network
  being drawn/edited — nodes and roads, kept in sync with what the
  engine actually holds.
- **Node**: A house, company, or intersection placed on the map.
- **Road**: A connection between two nodes drawn on the map.
- **Run Result**: The outcome of running a population against the current
  network — total/average travel time, kept around long enough to compare
  against a later run.
- **Agent (visual)**: The moving representation of one engine-computed
  route, animated from origin to destination.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every node on the map has a distinct, consistent visual
  marker per type (house, company, intersection), used for no other
  purpose, verifiable by inspection.
- **SC-002**: From loading a network and population to seeing agents
  visibly moving, no more than 5 seconds pass, for networks and
  populations within this project's already-established scale (Phase 3/4
  success criteria).
- **SC-003**: A user can go from an empty map to a runnable network (at
  least one house, one company, and one road connecting them) using only
  the on-screen drawing tools, in under 2 minutes.
- **SC-004**: After adding a road to an already-run network and
  re-running, both the before and the after average travel time remain
  visible together, without navigating away or losing the earlier result.
- **SC-005**: Removing a node also removes every road attached to it —
  verifiable afterward by the absence of any road referencing a missing
  node.

## Assumptions

- The browser reaches the existing Go engine (`graph`/`agent`/
  `simulation`) through a new server component; the engine packages
  themselves are not modified by this phase (per the confirmed
  architecture: engine code stays where it is, a server exposes it, the
  browser app lives in its own directory). The exact communication
  mechanism is a planning decision, not fixed by this spec.
- A run computes a population's full result in one pass, as the engine
  already does (Phases 3-4) — the browser animates each agent along its
  already-known, already-computed route rather than requiring the engine
  to simulate step-by-step in real time. This keeps the same one-shot
  computation model the engine already has; no live, tick-by-tick
  server-driven simulation is required for this phase.
- Editing the network while an animation is playing stops or pauses that
  animation; a fresh run must be explicitly requested to see the edited
  network's behavior. Live-updating an in-progress animation from
  concurrent edits is out of scope.
- Only one network and one population are being viewed/edited at a time
  in this phase — no multi-tab or multi-user collaborative editing.
- Persisting a hand-drawn network across browser sessions (e.g. to a file
  or database) is out of scope for this phase; the focus is on
  draw → run → observe within one session.
- This phase adds visual confirmation on top of, not instead of, the
  console/log validation established by constitution Principle IV —
  `cmd/graphcli` remains valid and continues to work unchanged.
