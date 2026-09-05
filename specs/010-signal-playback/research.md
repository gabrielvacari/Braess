# Phase 0 Research: Signal-Controlled Live Playback

## Decision 1: Extend `QueueSample` with `Green bool` rather than a new `SignalPhaseSample` type

**Decision**: Add a `Green bool` field to the existing `queuesim.QueueSample`
struct (`{Time, SignalID, Length, Green}`) instead of introducing a
separate per-tick signal-phase type.

**Rationale**: `QueueSample` is already recorded at exactly the cadence
FR-002 needs — once per signal, per tick. A parallel `SignalPhaseSample`
array would duplicate that same iteration and carry redundant
`Time`/`SignalID` pairs for no benefit. `Signal.IsGreenAt(t)` is already
computed internally as part of the discharge step, so recording it
alongside the queue length it's already producing is a one-line addition
with no new bookkeeping.

**Alternatives considered**: A separate `SignalPhaseSample`/
`RunResult.Phases` array — rejected as pure duplication of `QueueSample`'s
existing per-tick-per-signal shape.

## Decision 2: Correlate `PositionSample` to `AgentReport` via a new `AgentReport.ID`

**Decision**: Add `ID int` to `AgentReport` (spawn order, 0-based,
assigned once per spawned agent regardless of whether it ever moves) and
`AgentID int` to the new `PositionSample`, so a client can group a run's
flat `[]PositionSample` by `AgentID` and match each group back to its
`AgentReport`.

**Rationale**: `AgentReport`s are appended in *completion* order (arrived
agents as they finish, then leftover still-in-flight agents at the end)
— not spawn order — so `DemandIndex` alone doesn't uniquely identify one
agent among several sharing a demand. A stable ID assigned at spawn is
the minimum needed to reconstruct "this trace belongs to this report,"
and it's additive: existing field-named composite literals
(`AgentReport{DemandIndex: ..., ...}`) are unaffected by a new field.

**Alternatives considered**: Embedding the full position trace *inside*
`AgentReport` (e.g. `AgentReport.Positions []PositionSample`) — rejected
because `RunResult.QueueSamples` already establishes the convention of a
flat, time-ordered array at the `RunResult` level rather than nested per
entity; keeping `Positions` flat and separate is consistent with that.

## Decision 3: Record one `PositionSample` per in-flight agent per tick, including a final sample on the arrival tick

**Decision**: During each tick, after the discharge step (so a
just-promoted agent's state is accurate for that tick), record one
`PositionSample{Time, AgentID, EdgeID, Progress, Queued}` for every agent
still in the run — `Progress` is `1 - remaining/freeFlowTime` while
moving (0 while queued, since a queued agent hasn't started the edge
yet), `Queued` mirrors its `agentState`. An agent that arrives *during*
a tick gets one extra sample at `Progress: 1` on the edge it just
finished, recorded before it's removed from the in-flight set — without
it, the last edge of a trip would have no sample showing the agent ever
reaching the end of it.

**Rationale**: This is the minimum data a client needs to place an agent
at any queried simulated time: which road, how far along it, and whether
it's stopped waiting rather than moving (FR-001). Progress is already
implicit in the existing `remaining`/`freeFlow` bookkeeping `Run` keeps
for its own travel-time accounting — reading it out is free.

**Alternatives considered**: Recording (x, y) canvas coordinates
directly in the engine — rejected: `queuesim` has no node-position
concept at all (Principle I; positions are a UI/web concern from feature
005 onward), and a client can always derive (x, y) from `EdgeID` +
`Progress` + the node positions it already has from the request it sent.

## Decision 4: A single shared playback clock, not one animation-hook instance per agent

**Decision**: The frontend introduces one `usePlaybackClock` hook holding
one simulated-time value for the whole run, advanced each animation
frame and clamped to `[0, duration]`. Every agent's canvas position and
every signal-controlled road's color are then pure functions of that one
current time, recomputed each frame.

**Rationale**: The equilibrium mode's `AgentTicker`/`useAgentAnimation`
(feature 005) intentionally gives each agent its own independent
real-time tween, because each equilibrium route is animated on its own
made-up wall-clock schedule with no shared clock to synchronize against.
Signals-mode agents are the opposite: their movements are *causally
linked* through shared queues and shared signals recorded against one
simulated timeline, so a single shared clock is both simpler and more
correct — there is exactly one "current simulated time" for the whole
scene, matching how the engine itself produced the data.

**Alternatives considered**: Reusing `useAgentAnimation`/`AgentTicker`
per agent by converting each agent's samples into a synthetic path —
rejected: it would need one instance per spawned agent (up to a few
hundred), each independently reading the wall clock, only to reconstruct
by other means the single shared timeline that already exists; a single
clock is both less code and the more accurate model of what's actually
being played back.

## Decision 5: `NetworkCanvas` gains one optional `roadColors` prop, not a Signals-only fork

**Decision**: Add an optional `roadColors?: Record<string, string>` prop
to `NetworkCanvasProps`. When a road's id is a key, its line uses that
color instead of the default gray (selection highlighting still takes
precedence). Omitted entirely, behavior is identical to today.

**Rationale**: `NetworkCanvas` is already shared, render-only, and takes
whatever `roads`/`agents` it's given (research.md decision #3, feature
009) — the equilibrium mode simply never passes `roadColors`, so it is
untouched by this addition. Forking a parallel canvas component for
Signals mode would duplicate the node/road/agent rendering feature 009
already established as shared.

**Alternatives considered**: A second canvas component specific to
Signals mode — rejected as duplicating rendering logic for no reason
other than one extra prop.

## Decision 6: Keep the existing chart/summary; add playback alongside, not instead

**Decision**: `QueueResultsPanel` (feature 009) is unchanged. The
animation and its play/pause/restart controls are new elements added to
`SignalApp`'s layout above or beside it, per spec Assumption
("complemented, not replaced").

**Rationale**: FR-008 and User Story 3 require both to coexist; nothing
about the chart/summary was reported as wrong, only insufficient on its
own.

**Alternatives considered**: Merging the chart into a tab alongside the
animation — rejected as an unnecessary interaction cost for a feature
whose whole point is showing both at once for comparison.
