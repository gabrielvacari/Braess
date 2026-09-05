# Feature Specification: Signal Queuing Frontend

**Feature Branch**: `009-signal-frontend`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Expose the queuesim engine capability (feature 008: signal-controlled queuing) in the web frontend. A user should be able to attach a traffic signal (green/red duration) to a road they draw, declare a time-based demand (agents arriving at a steady rate, not all at once) between two nodes, run the signal-controlled simulation, and see the result: each signal-controlled road's queue length changing over time (a chart or equivalent time-series visualization), and a summary of how many agents arrived vs. are still waiting, and average wait time. This is a new, distinct simulation mode from the existing equilibrium-based one (features 003/004/005) since queuesim models discrete time and real queuing rather than a static converged assignment — it does not replace the existing 'run the population' flow, it adds this as a separate capability. No change to graph/agent/simulation/queuesim/cmd/graphcli — this only adds a way to reach the already-validated queuesim engine from the browser, mirroring how feature 005 exposed simulation/RunDemands."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure a signal and a time-based demand (Priority: P1)

A user attaches a traffic signal (green and red durations) to a road they
draw, and declares a demand as agents arriving at a steady rate over
time, rather than only as a fixed total present from the start.

**Why this priority**: Without a way to configure these, there is nothing
new to run — this is the entry point for every other capability in this
feature.

**Independent Test**: Can be fully tested by drawing a road, attaching a
signal to it, declaring a time-based demand, and confirming both are
reflected in the visible network/demand state before anything is run.

**Acceptance Scenarios**:

1. **Given** a drawn road, **When** the user attaches a signal to it,
   **Then** that road's green and red durations are set and visible.
2. **Given** a network with at least two nodes, **When** the user
   declares a demand with an arrival rate, **Then** that demand is
   recorded with its rate, distinct from the existing fixed-count demand.

---

### User Story 2 - Run the signal-controlled simulation and see queues over time (Priority: P2)

After configuring signals and time-based demands, a user runs the
signal-controlled simulation and sees each signal-controlled road's queue
length change over time — enough to judge, by looking, whether one road
ends up more congested, the congestion oscillates, or it balances,
without the app presuming which outcome occurs.

**Why this priority**: This is the actual payoff — the same "make the
result observable, don't presume it" principle feature 008 established
for logs, now made visible.

**Independent Test**: Can be fully tested by running a small two-road,
opposite-signal scenario and confirming a queue-length-over-time
visualization renders for each signal-controlled road, reflecting the
underlying simulation's real output.

**Acceptance Scenarios**:

1. **Given** a network with at least one signal-controlled road and a
   time-based demand, **When** the user runs the signal-controlled
   simulation, **Then** a queue-length-over-time visualization appears
   for every signal-controlled road.
2. **Given** two signal-controlled roads configured with different
   timings, **When** the simulation runs, **Then** their visualizations
   are visibly distinguishable from one another.

---

### User Story 3 - See an arrival and waiting summary (Priority: P3)

After a signal-controlled run, a user sees how many agents arrived versus
are still waiting, and the average wait time, without needing to read
raw data to figure it out.

**Why this priority**: Matches feature 008's own requirement that an
agent still waiting must never be presented as having arrived — now
surfaced visually so a user can judge a configuration at a glance.

**Independent Test**: Can be fully tested by running a scenario where at
least one agent doesn't finish before the run ends, and confirming the
displayed summary distinguishes arrived from still-waiting agents and
shows an average wait time.

**Acceptance Scenarios**:

1. **Given** a completed signal-controlled run, **When** the user looks
   at the summary, **Then** the number of agents arrived and the number
   still waiting are both shown, along with the average wait time.

---

### Edge Cases

- What happens when the signal-controlled simulation is run with no
  signal-controlled roads configured? It must run and report normally
  (no queues form, all reported queue lengths are zero), not error.
- What happens when a user switches between this signal-controlled mode
  and the existing equilibrium-based "run the population" flow? Each
  mode's own configuration and results must stay distinct — one is never
  shown as if it were the other's, and configuring one does not silently
  affect the other.
- What happens when a signal-controlled run fails (e.g. a demand's
  destination is unreachable)? It must be reported clearly, the same way
  the existing run flow already reports its own failures.
- What happens when a road has no signal attached and is used in a
  signal-controlled run? It behaves as an always-open road (no queuing
  ever applies to it), consistent with feature 008's own engine
  behavior.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST let a user attach a traffic signal (a green
  duration and a red duration) to a road they draw.
- **FR-002**: The system MUST let a user declare a demand as agents
  arriving at a steady rate over time, rather than only as a fixed total
  present from the start.
- **FR-003**: The system MUST let a user run a signal-controlled
  simulation as a capability separate from the existing population-
  equilibrium run.
- **FR-004**: The system MUST display, after a signal-controlled run,
  each signal-controlled road's queue length changing over time.
- **FR-005**: The system MUST display, after a signal-controlled run, how
  many agents arrived versus are still waiting, and the average wait
  time.
- **FR-006**: The system MUST NOT alter or remove the existing
  equilibrium-based "run the population" capability — this feature adds
  a separate capability alongside it.
- **FR-007**: The system MUST clearly indicate which simulation mode
  (equilibrium-based or signal-controlled) a shown result belongs to, so
  the two are never confused.
- **FR-008**: The system MUST report a signal-controlled run's failure
  clearly, consistent with how the existing run flow already reports its
  own failures.

### Key Entities

- **Signal Configuration**: A road's green and red durations, attached
  when that road is drawn.
- **Time-Based Demand**: An origin/destination pair with an arrival rate
  and a total count, distinct from the existing fixed-count demand.
- **Queue Report**: The queue-length-over-time visualization for each
  signal-controlled road, plus the arrival/waiting summary, produced by
  one signal-controlled run.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can attach a signal to a road and see its
  configured durations reflected in the interface, without inspecting
  any code or console output.
- **SC-002**: After running a signal-controlled scenario, a queue-length-
  over-time visualization is shown for every signal-controlled road in
  the network.
- **SC-003**: After a run, the number of agents still waiting is visibly
  distinguishable from the number arrived, without the user needing to
  compute it manually from raw data.
- **SC-004**: A user can tell, at a glance, whether a shown result came
  from the equilibrium-based run or the signal-controlled run.
- **SC-005**: Running the same two-road, opposite-signal configuration
  twice with different signal durations produces visibly different
  queue-length-over-time results.

## Assumptions

- This feature only adds a new way to reach the already-validated
  `queuesim` engine capability from the browser — `graph`, `agent`,
  `simulation`, `queuesim`, and `cmd/graphcli` are not modified.
- The signal-controlled simulation and the existing equilibrium-based
  simulation (features 003-005) are separate, coexisting capabilities in
  the same application — neither replaces nor merges into the other.
- A road may have at most one signal attached in this feature; multiple
  signals per road, or signals not tied to a specific road, are out of
  scope.
- The specific visual form of "queue length over time" (e.g. a line
  chart, a sequence of bars) is a planning decision, not fixed by this
  spec, as long as it clearly conveys the time series.
- Configuration fields feature 008 supports but this spec doesn't call
  out (e.g. a custom discharge-rate override) may default sensibly
  rather than requiring their own control in a first version.
