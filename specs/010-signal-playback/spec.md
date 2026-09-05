# Feature Specification: Signal-Controlled Live Playback

**Feature Branch**: `010-signal-playback`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "Extend the queuesim engine to record each agent's position over time (which edge, whether queued or moving, progress along the edge) at every simulation tick, so a client can reconstruct a full time-based playback — not just the final summary it produces today. Then, in the web frontend's Signals mode, replace (or complement) the current chart-only result with a real, animated visualization: agents actually move across the canvas from origin to destination, pausing when queued at a red signal and proceeding when it turns green, exactly mirroring the equilibrium mode's existing car-animation experience. The road itself must visually indicate its signal's current state (e.g. colored green/red) as the animation plays through simulated time. This is an engine extension (new capability, must be validated in text/tests before any UI per the constitution's Principle II) plus a frontend change to Signals mode only — the equilibrium mode's existing animation and the queuesim engine's existing Run/Signal/Demand behavior must not change, only be extended with additive data."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Watch agents move through a signal-controlled network (Priority: P1)

A person runs a signal-controlled scenario and, instead of only reading a
queue-length chart, watches each agent actually travel from its origin to
its destination — moving along the road while free to go, and visibly
stopping and waiting when its road's signal turns against it, then
resuming the moment it turns green.

**Why this priority**: This is the whole reason the feature exists — the
chart-only result was explicitly reported as insufficient; seeing traffic
build up and dissipate in front of you is what makes congestion (and,
eventually, Braess's Paradox) legible at a glance, exactly like the
equilibrium mode already does for its own scenarios.

**Independent Test**: Configure the classic two opposite-signal roads
between one origin and one destination, run it, and confirm that agents
visibly queue on one road while moving on the other, then swap roles when
the signals flip.

**Acceptance Scenarios**:

1. **Given** a completed signal-controlled run, **When** the person
   presses play, **Then** every agent from that run appears at its origin
   and moves toward its destination along the road(s) it actually used.
2. **Given** an agent whose road is currently red, **When** the animation
   reaches the moment that agent tried to enter that road, **Then** the
   agent visibly stops and waits at the queue instead of continuing
   through the road.
3. **Given** an agent waiting at a red signal, **When** the animation's
   simulated clock reaches the moment that signal turns green and the
   agent is granted entry, **Then** the agent resumes moving.

---

### User Story 2 - See each signal's live state as the animation plays (Priority: P1)

While watching the animation, a person can tell at a glance whether a
given road's signal is currently green or red, without reading a
separate chart, so the reason agents are stopped (or moving) is always
visible on the road itself.

**Why this priority**: Without this, a person sees agents stop and start
with no visible cause — the signal state is the explanation for the
whole behavior this feature exists to show. It is inseparable from User
Story 1 in practice, but is called out on its own because it can be
verified independently of whether agent movement itself looks right.

**Independent Test**: Run a scenario with one signal-controlled road,
play the animation, and confirm the road's visual state changes from
green to red and back on schedule, matching the signal's configured
durations.

**Acceptance Scenarios**:

1. **Given** a signal-controlled road, **When** its signal is green at
   the animation's current simulated time, **Then** the road is shown in
   its "green" visual state.
2. **Given** the same road, **When** its signal turns red, **Then** the
   road's visual state changes to "red" at that same simulated moment.
3. **Given** a road with no signal attached, **When** the animation
   plays, **Then** that road never shows a green/red state — it reads as
   unconstrained, consistent with how it behaves in the simulation.

---

### User Story 3 - Keep the existing summary alongside the animation (Priority: P3)

After a run, a person can still see the queue-length-over-time chart and
the arrived/still-waiting/average-wait summary introduced by the
previous feature, presented together with the new animation rather than
lost when it was added.

**Why this priority**: Lower priority because it preserves already-
delivered value rather than adding new value, but dropping it would be a
regression a person would notice immediately after this feature ships.

**Independent Test**: Run a scenario, confirm the animation plays, and
confirm the same queue chart and arrival/wait summary from the prior
feature are still visible on the same screen.

**Acceptance Scenarios**:

1. **Given** a completed run, **When** the result is shown, **Then**
   both the live animation and the existing chart/summary are visible
   together, not one in place of the other.

---

### Edge Cases

- An agent that is still queued or still moving when the run's duration
  ends: the animation ends with that agent visible at its last known
  position, not vanished or thrown to an error state.
- A road with a signal that a person removes or edits after a run has
  already completed: the just-played animation continues to reflect the
  scenario as it was actually run, not the edited version.
- Replaying the same completed run's animation a second time (without
  re-running the simulation) produces the exact same movement every
  time.
- A scenario with no signal-controlled roads at all: the animation still
  plays (agents just never queue), and no road ever shows a green/red
  state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The simulation engine MUST record, at every simulation
  tick, each agent's location (which road it is on, and how far along
  it) and whether the agent is currently moving or queued.
- **FR-002**: The simulation engine MUST record, at every simulation
  tick, each signal-controlled road's current state (green or red).
- **FR-003**: The Signals-mode result MUST animate every agent from that
  run moving along the road(s) it actually used, at a pace derived from
  the recorded simulation time, including a visible pause for any tick
  during which that agent was queued.
- **FR-004**: The Signals-mode result MUST visually indicate, on each
  signal-controlled road, whether its signal is currently green or red
  as of the animation's current simulated time, updating as playback
  advances.
- **FR-005**: A road with no signal attached MUST NOT show a green/red
  state at any point during playback.
- **FR-006**: The animation MUST be replayable from the already-computed
  run result without requiring the simulation to be run again.
- **FR-007**: This feature MUST NOT change the existing queue-length
  samples or the existing arrived/wait agent summary already produced by
  the prior feature — the new per-tick position and signal-state data is
  additive alongside them.
- **FR-008**: The Signals-mode result MUST continue to show the existing
  queue-length chart and the arrived/still-waiting/average-wait summary
  together with the new animation.
- **FR-009**: This feature MUST NOT alter the equilibrium mode's
  existing network, run flow, or car animation in any way.

### Key Entities

- **Agent Position Sample**: One agent's recorded state at one
  simulation tick — which road it occupies, its progress along that
  road, and whether it is moving or queued. A full run produces one
  ordered sequence of these per agent, from spawn to arrival (or to the
  end of the run, if it never arrives).
- **Signal Phase Sample**: One signal-controlled road's recorded state
  (green or red) at one simulation tick. A full run produces one ordered
  sequence of these per signal-controlled road.
- **Playback**: The client-side act of stepping through a run's
  recorded Agent Position Samples and Signal Phase Samples in order,
  rendering each agent's position and each road's signal state as of the
  current point in simulated time. Introduces no new decision-making —
  every position and every signal state it shows already happened during
  the run.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A person watching the animation can identify, without
  reading any chart or number, which of two competing roads is
  currently more congested, just by observing where agents are
  stacking up.
- **SC-002**: A person can correctly state whether a given road's
  signal is green or red at any moment during playback just by looking
  at the road, with no separate legend lookup needed beyond the
  color/state key already established for node types.
- **SC-003**: Replaying the same completed run's animation twice
  produces identical agent movement both times.
- **SC-004**: Two runs of the same scenario with different signal
  timings produce visibly different animations (agents queue on
  different roads, or for different durations), without needing to
  re-read the underlying chart to notice the difference.
- **SC-005**: The existing queue chart and arrival/wait summary remain
  visible and correct after this feature ships, with zero loss of
  information a person had access to before.

## Assumptions

- The chart-only result is complemented, not replaced: the animation is
  the primary new view, but the existing queue-length chart and
  arrival/wait summary (feature 009) remain on the same screen.
- Playback pace mirrors the equilibrium mode's existing convention
  (simulated time mapped to a comfortable, fixed real-time pace for
  watching) rather than strict real-time (a scenario's simulated
  duration may be much longer than is pleasant to watch at 1:1 speed).
- One signal per road remains true (feature 009's existing assumption,
  unchanged) — this feature does not introduce multiple signals per
  road.
- The client renders an agent's position at each tick and smoothly
  interpolates between consecutive ticks for a fluid appearance, the
  same technique the equilibrium mode's animation already uses between
  its route waypoints — this feature does not require sub-tick
  simulation precision.
- Play/pause/replay controls for the animation are the minimum needed to
  satisfy "replayable" (FR-006); a full scrubber/seek control is not
  required for this feature.
- Data volume (one position sample per agent per tick, for the run's
  full duration) is acceptable for the hand-drawn, exploratory scenarios
  this tool targets (tens to a few hundred agents, run durations of
  tens to low hundreds of simulated seconds) — no pagination or
  down-sampling is required at this scale.
