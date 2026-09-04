# Feature Specification: Shortest-Path Agents

**Feature Branch**: `002-shortest-path-agents`

**Created**: 2026-09-04

**Status**: Draft

**Input**: User description: "Phase 2 of the roadmap only (per constitution Principle II): simple agents (cars) that go from a fixed origin node to a fixed destination node using shortest-path routing over the existing graph engine, with no congestion modeling yet (edge travel time is evaluated at a fixed/zero traffic volume, not affected by other agents). Each agent independently computes its own route to minimize its own travel time (constitution Principle III — decentralized, selfish behavior — even though with no congestion yet, shortest path is trivially each agent's selfish choice). This still has no UI: agent routing and the resulting path/travel time must be verifiable via console/log output or a terminal command, building on the existing graph engine and its cmd/graphcli. Multiple agents may be simulated for the same fixed origin/destination pair, but multiple distinct origins/destinations (houses/companies) are out of scope until roadmap phase 4."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Compute an agent's shortest route (Priority: P1)

A developer asks the system for the best route between a fixed origin node
and a fixed destination node on the existing road network, and gets back
the exact sequence of roads to take and the total time that route takes.

**Why this priority**: This is the core new capability of this phase — the
direct extension of the Phase 1 graph engine into actual decision-making.
Nothing else in this feature has meaning without it.

**Independent Test**: Can be fully tested by building a small graph with a
known correct shortest path and checking the computed route matches it —
no other agents, congestion, or UI involved.

**Acceptance Scenarios**:

1. **Given** a graph with a single possible path from the origin to the
   destination, **When** a route is computed, **Then** it returns that path
   together with its total travel time.
2. **Given** a graph with multiple possible paths between the origin and
   the destination that have different total travel times, **When** a
   route is computed, **Then** it returns the path with the minimum total
   travel time.

---

### User Story 2 - Multiple agents decide independently (Priority: P2)

Several agents are all placed on the same fixed origin/destination pair.
Each one works out its own route on its own — nothing computes one answer
once and hands it to all of them.

**Why this priority**: This is where the project's decentralized-agent
premise gets validated for the first time, even in this simplified,
congestion-free case: each agent's choice must genuinely come from its own
computation, so that when congestion is introduced in the next phase,
agents already behave independently rather than needing to be rewired.

**Independent Test**: Can be fully tested by creating several independent
agents for the same origin/destination and confirming each one's route
comes from its own computation — no UI involved.

**Acceptance Scenarios**:

1. **Given** three independently-created agents sharing the same
   origin/destination on a graph with one unique shortest path, **When**
   each agent computes its route, **Then** each reports the same path and
   travel time, each arrived at through its own independent computation.
2. **Given** a graph where two paths between the origin and destination
   have equal total travel time, **When** multiple agents compute their
   routes, **Then** each agent's chosen path is the direct result of that
   agent's own computation (not copied from another agent).

---

### User Story 3 - Inspect agent routing from a terminal (Priority: P3)

A developer runs the existing standalone terminal command, now extended to
place one or more agents on a fixed origin/destination, and sees each
agent's chosen route and travel time printed to the console.

**Why this priority**: Proves, end to end, that agent routing is verifiable
in text before any frontend work starts — extending the same terminal
command from the previous phase rather than inventing a new verification
path.

**Independent Test**: Can be fully tested by running the command and
checking its console output and exit code — no browser or UI required.

**Acceptance Scenarios**:

1. **Given** the existing sample network plus a fixed origin and
   destination, **When** the terminal command runs, **Then** it prints,
   for each simulated agent, its ordered route and total travel time,
   exiting successfully.

---

### Edge Cases

- What happens when no path exists between the origin and the destination?
  The system must report this clearly (e.g. "no route found") rather than
  crashing or returning a misleading empty/zero-time result.
- What happens when the origin and destination are the same node? The
  system must return a trivial route with no edges and zero travel time,
  not an error.
- How does the system handle multiple parallel roads between the same pair
  of nodes (supported since Phase 1)? Each parallel edge must be
  considered as a distinct route option, and the one with the lower travel
  time must be preferred.
- What happens if evaluating a road's travel time itself fails? The
  failure must be reported, not silently treated as if that road did not
  exist or were instantly fast.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST compute, for a given origin node and
  destination node in a graph, the route (ordered sequence of edges) that
  minimizes total travel time, evaluated at a fixed baseline traffic
  volume (zero, since no congestion exists yet in this phase).
- **FR-002**: The system MUST report, alongside a computed route, its
  total travel time (the sum of its edges' travel times at the baseline
  volume).
- **FR-003**: The system MUST allow creating multiple independent agents
  for the same origin/destination pair, each of which computes its own
  route through its own independent computation — never a route copied or
  imposed from outside that agent.
- **FR-004**: The system MUST treat multiple parallel edges between the
  same node pair as distinct route options when computing the shortest
  route, selecting whichever has the lower travel time.
- **FR-005**: The system MUST report a clear "no route exists" result when
  no path connects the origin to the destination.
- **FR-006**: The system MUST return a zero-edge, zero-travel-time route
  when the origin and destination are the same node.
- **FR-007**: The system MUST be operable from the same standalone
  terminal entry point established in the previous phase, continuing to
  require zero UI/frontend dependency.
- **FR-008**: The system MUST propagate an error rather than silently
  substituting a value if evaluating an edge's travel time fails during
  route computation.

### Key Entities

- **Agent**: A car with a fixed origin node and a fixed destination node,
  capable of computing its own route.
- **Route**: The ordered sequence of edges an agent will take from its
  origin to its destination, together with the route's total travel time
  at the baseline volume.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Across at least 5 distinct test networks with a uniquely
  shortest path between two given nodes, the computed route matches the
  mathematically correct minimum-travel-time path in 100% of cases.
- **SC-002**: Route computation for a network of at least 10 nodes and 15
  edges completes and prints via the terminal command in under 5 seconds.
- **SC-003**: Running route computation for at least 3 independently
  created agents sharing the same origin/destination on an unchanged
  network yields identical routes and travel times in 100% of runs, each
  traced to that agent's own independent computation.
- **SC-004**: When no route exists between the given origin and
  destination, the system reports that clearly in 100% of such cases, with
  no crash and no misleading zero/empty result.

## Assumptions

- "Fixed/zero traffic volume" means every route computation in this phase
  evaluates each edge's travel-time behavior at its free-flow value —
  congestion (travel time responding to how many agents share a road)
  is out of scope until roadmap Phase 3, per constitution Principle II.
- Routes are chosen by minimizing total travel time only in this phase,
  not distance — travel time is the quantity later phases (congestion,
  Braess's Paradox) will operate on.
- Only one externally-fixed origin/destination pair is exercised in this
  phase; "multiple agents" refers to multiple agent instances sharing that
  one pair, not multiple distinct houses/companies — those arrive in
  roadmap Phase 4.
- Tie-breaking between two equally-shortest routes may be deterministic
  (e.g. by a consistent ordering rule); this spec does not require
  randomized tie-breaking, since nothing yet depends on agents diverging
  when routes are tied.
- This phase remains in-memory only, consistent with Phase 1.
