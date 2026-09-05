# Feature Specification: Signal-Controlled Queuing

**Feature Branch**: `008-signal-controlled-queuing`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Add signal-controlled queuing to the simulation engine: a road can be governed by a traffic signal that alternates between a green phase (passable) and a red phase (not passable), each with a configurable duration. Agents arriving at a signal-controlled road during red are held in a queue (not blocked/rerouted, not discarded) and released once it turns green, at a defined discharge rate, until the queue empties or it turns red again. This requires the simulation to advance through discrete time steps rather than resolving directly to one final static equilibrium (a departure from the existing Phase 3/4 model). Support at least two signal-controlled roads whose phases are explicitly linked to always be opposite each other (one green while the other is red), matching the classic 'alternating signal' setup, and report each road's queue length over time plus each agent's total waiting time, so it's possible to observe empirically whether one road ends up persistently more congested, the congestion oscillates, or it balances out — the spec does not presume which outcome occurs. Per constitution Principle II, this feature covers the simulation-engine capability only, validated via console/log output; a browser UI for configuring signal timing or visualizing queues is an explicitly separate, later feature."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Agents queue at a red signal and are released on green (Priority: P1)

A road can be governed by a signal that alternates between passable
(green) and not passable (red). Agents that arrive while it's red wait in
order rather than disappearing or being rerouted; once it turns green,
they proceed at a defined rate until the queue clears or it turns red
again.

**Why this priority**: This is the entire new mechanism — real waiting.
Without it, this feature is no different from the static congestion
already built in Phase 3. Every other capability in this feature depends
on queuing actually working correctly.

**Independent Test**: Can be fully tested by running a single
signal-controlled road on a fixed cycle with a steady stream of arriving
agents, and confirming — from console/log output — that no agent crosses
during red, and that queued agents leave in the order they arrived once
it turns green.

**Acceptance Scenarios**:

1. **Given** a signal-controlled road currently in its red phase, **When**
   an agent arrives at it, **Then** the agent is added to that road's
   queue rather than crossing or being dropped.
2. **Given** a signal-controlled road with agents waiting in its queue,
   **When** the signal turns green, **Then** agents leave the queue in
   arrival order, at the road's defined discharge rate, until the queue
   is empty or the signal turns red again.

---

### User Story 2 - Observe what alternating signals actually do to congestion (Priority: P2)

Two roads between the same two points are each governed by a signal, the
two signals always in opposite phase (one green exactly when the other is
red), swapping on a timer. Running many agents through this setup over
time produces a report of each road's queue length over time and each
agent's total waiting time — making it possible to see, from the data,
whether one road ends up persistently worse, the congestion alternates
between them, or it balances out.

**Why this priority**: This is the actual experiment that motivated the
feature. It doesn't assume an answer — it makes the answer observable,
the same experimental spirit the project has applied to Braess's Paradox
since Phase 3: build the mechanism faithfully, then see what happens.

**Independent Test**: Can be fully tested by running the two-road,
opposite-phase scenario for a fixed duration with a steady arrival rate,
and confirming the report shows each road's queue length over time and
each agent's total waiting time, whatever pattern that data turns out to
show.

**Acceptance Scenarios**:

1. **Given** two roads with signals always in opposite phase, **When**
   agents travel between the two points those roads connect over an
   extended run, **Then** a report shows each road's queue length at
   points across the run.
2. **Given** the same run, **When** it completes, **Then** each agent's
   total time spent waiting in a queue is available, separately from
   time spent actually moving.

---

### User Story 3 - Configure signal timing to run different experiments (Priority: P3)

A signal's green and red durations are set per scenario, not fixed in
code, so different timing plans can be tried and compared.

**Why this priority**: Without configurability, this is one fixed demo
rather than a tool for experimenting — the same reason earlier phases
(e.g. multiple demands, Phase 4) made their inputs configurable rather
than hardcoded.

**Independent Test**: Can be fully tested by running the same two-road
scenario twice with two different green/red duration configurations and
confirming the reported queue patterns differ between the two runs.

**Acceptance Scenarios**:

1. **Given** a signal's green and red durations are set to one
   configuration, **When** the scenario is run and then re-run with
   different durations, **Then** the two runs' reported queue patterns
   are different, reflecting the different timing.

---

### Edge Cases

- What happens when arrivals during a red phase outpace what the
  discharge rate can clear during the following green phase? The
  remaining queue must carry over into the next cycle, not be silently
  cleared or reset.
- What happens to an agent whose route never gets a green phase before
  the simulated run ends (e.g. a misconfiguration where a signal never
  turns green)? It must be clearly reported as still waiting at the end
  of the run, not presented as having arrived.
- What happens to an agent already past a signal's queue point and
  moving when that signal turns red? It is not interrupted mid-trip —
  only agents not yet through the signal at that moment must wait.
- What happens when a signal's discharge rate isn't explicitly
  configured? A reasonable default, derived from the road's own normal
  (non-signal) travel-time behavior, applies.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST let a road be governed by a signal that
  alternates between a green phase (passable) and a red phase (not
  passable), each with a configurable duration.
- **FR-002**: The system MUST hold (queue) an agent that arrives at a
  signal-controlled road during its red phase, rather than letting it
  pass or discarding it.
- **FR-003**: The system MUST release queued agents, in the order they
  arrived, once the signal turns green, at a defined per-road discharge
  rate, until the queue is empty or the signal returns to red.
- **FR-004**: The system MUST report, over the course of a simulated run,
  how each signal-controlled road's queue length changed over time.
- **FR-005**: The system MUST support at least two signal-controlled
  roads whose phases are explicitly linked to always be opposite each
  other.
- **FR-006**: The system MUST report each agent's total time spent
  waiting in a queue, separate from time spent actually moving.
- **FR-007**: The system MUST allow a signal's green and red durations to
  be configured per scenario, not fixed in code.
- **FR-008**: The system MUST clearly report an agent still waiting at
  the end of a simulated run as still waiting, not as having completed
  its trip.
- **FR-009**: This capability MUST be validated via console/log output
  before any browser-based configuration or visualization is built for
  it — a frontend for signals is a separate, later feature.

### Key Entities

- **Signal**: Governs one road; has a green duration and a red duration,
  and may be linked to another signal to always be in the opposite phase.
- **Queue**: The agents currently waiting at a signal-controlled road's
  entrance, in arrival order.
- **Simulated Run**: A run of this feature advances through discrete time
  steps rather than resolving directly to one final assignment, unlike
  Phase 3/4's equilibrium-based runs.
- **Wait Report**: The output of a run — each signal-controlled road's
  queue length over time, and each agent's total waiting time.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In a single signal-controlled road scenario, the reported
  queue length reaches zero by the end of any green phase long enough to
  clear the arrivals from the preceding red phase.
- **SC-002**: Across an entire run, no agent's report shows it having
  crossed a signal-controlled road during that road's red phase.
- **SC-003**: The two-alternating-signal scenario produces a
  distinguishable, measurable queue-length-over-time report for each
  road, sufficient to determine from the data alone whether one road
  ended up more congested, the congestion alternated, or it balanced.
- **SC-004**: Changing a signal's configured durations changes the
  reported queue pattern, verified across at least two distinct timing
  configurations on the same scenario.
- **SC-005**: An agent that never gets released within the simulated run
  duration is reported as still waiting, not as arrived, in 100% of such
  cases.

## Assumptions

- This spec covers the simulation-engine capability only (extending
  `graph`/`agent`/`simulation`) — a browser UI for configuring signal
  timing or visualizing queues is a distinct, later feature, consistent
  with every earlier roadmap phase's order: simulate and validate in text
  before visualizing (constitution Principle II).
- An agent already past a signal's queue point and moving when it turns
  red is not interrupted mid-trip; only agents not yet through it at that
  moment must wait.
- A signal's discharge rate (how many queued agents can leave per unit of
  simulated time once green) defaults to a value derived from the road's
  own normal, non-signal travel-time behavior when not explicitly
  configured.
- This spec does not require modeling intersections with more than two
  mutually-exclusive movements, yellow/amber phases, or turning
  movements — only a road (or a small, explicitly linked set of roads)
  alternating between passable and not-passable on a timer.
- Whether the two-road alternating-signal experiment produces a
  persistently more-congested road, an oscillating pattern, or a balanced
  outcome is the empirical result this feature makes observable — not a
  predetermined property the feature must guarantee.
