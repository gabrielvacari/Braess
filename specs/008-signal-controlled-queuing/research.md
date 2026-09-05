# Phase 0 Research: Signal-Controlled Queuing

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, but
this feature required more genuinely new design decisions than any prior
one — the user's own clarification (real queuing, not just blocking)
committed this to a real time-stepped simulation, a first for this
project. Per constitution Principle V, every load-bearing decision is
recorded here.

## 1. Fixed-tick simulation, not discrete-event

**Decision**: The simulation advances in fixed steps of `dt` simulated
seconds (e.g. 1s), from `t=0` to a configured `duration`. Every tick,
every in-flight agent's remaining time-on-edge is decremented by `dt`,
signal phases are evaluated at `t`, and queues are discharged.

**Rationale**: The alternative — a discrete-event simulation that jumps
directly from one meaningful event (an agent finishing an edge, a signal
changing phase) to the next — is more precise (no tick-granularity error)
but is a materially bigger piece of engineering (an event queue, event
ordering/tie-breaking rules, etc.) for a benefit this feature doesn't need:
`dt` can simply be chosen small relative to the shortest relevant duration
(a signal's green/red length, an edge's travel time), making the
granularity error negligible. Fixed-tick simulation is also the more
common, more easily explained approach for exactly this class of problem
(it's essentially a cellular-automaton-style traffic model, in the
tradition of Nagel-Schreckenberg-style simulations) — easier to reason
about, test, and extend later than an event queue would be to build
correctly the first time.

**Alternatives considered**: Discrete-event simulation — rejected for the
complexity-vs-benefit reason above, for this feature's scope.

## 2. A route is chosen once, at spawn, from a snapshot — never revised

**Decision**: When an agent spawns, it computes its entire route once,
pricing each candidate signal-controlled edge using that signal's state
*at the spawn instant* (a snapshot: is it green or red right now, and how
long until it changes / how long is the current queue). The agent then
follows that fixed route for its whole trip; it does not re-plan once
underway, even if conditions along its remaining route change before it
gets there.

**Rationale**: This matches the spec's own Assumptions directly ("an
agent already past a signal's queue point ... is not interrupted") and is
the natural, minimal generalization of how every earlier phase's agents
decided: once (Phase 2), given currently-observable conditions (Phase
3/4) — never continuously replanning mid-trip. It also keeps this
feature's core routing algorithm the same class of problem
`agent`/`simulation` already solve (a one-shot shortest path against a
priced graph), rather than a full time-dependent shortest-path problem
(where an edge's cost depends on *when* you'd arrive at it, which the
agent can't know exactly without already having decided how it gets
there — a well-known, substantially harder problem). For this feature's
primary scenario — one road directly connecting the two points a demand
travels between — the route is exactly one edge, so the "spawn-time
snapshot" and "state at actual arrival" are identical anyway; the
simplification only introduces approximation error for a longer,
multi-signal route, which this feature's own scope (spec Assumptions: no
complex intersections) doesn't emphasize.

**Alternatives considered**: A full time-dependent shortest path
(predicting each downstream signal's phase at the agent's *predicted*
arrival time, recursively) — rejected as materially harder for a benefit
this feature's own target scenario doesn't need (see above); continuous
mid-trip re-planning — rejected per the spec's own Assumption ruling it
out.

## 3. `queuesim` implements its own pathfinding — doesn't reuse `agent`'s core

**Decision**: `queuesim/route.go` has its own Dijkstra-style function,
pricing each edge as `graph.Edge.TravelTime(0)` (free-flow) plus, for a
signal-controlled edge, an estimated current wait. It does not call into
`agent.ShortestRoute`/`ShortestRouteAtVolumes` or their shared internal
core.

**Rationale**: `agent`'s shared core (feature 003) is parameterized by a
*volume* — a single number per edge, fed into that edge's
`graph.TravelTimeFunc`. A signal's wait estimate isn't a volume and
doesn't flow through `graph.TravelTimeFunc` at all (a red light doesn't
make the road's congestion function return a bigger number; it makes the
road briefly impassable, which is a different kind of fact about an
edge). Bending `agent`'s volume-shaped core to also carry signal-wait
information would compromise a clean, already-shipped, already-tested
piece of code for a use case it wasn't designed for. A second, small,
self-contained Dijkstra in `queuesim` — following the same shape and
lazy-deletion-heap technique `agent` already validated — costs little and
keeps both packages simple to reason about independently.

**Alternatives considered**: Generalizing `agent`'s core to accept an
arbitrary per-edge cost function instead of a volume — rejected; it would
be a non-trivial, behavior-risking change to a feature-002/003 file for
the sole benefit of this feature, when a fresh ~40-line function achieves
the same result with zero risk to already-shipped code (echoing feature
004's own reasoning for why it *did* generalize `simulation`'s core: there,
the generalization was a one-line-per-call-site change to the *same*
volume-shaped model; here, the underlying cost model itself is different).

## 4. "Opposite phase" is pure configuration, not a data relationship

**Decision**: A `Signal` is `{ID, EdgeID, GreenDuration, RedDuration,
Offset}`. `IsGreenAt(t)` computes `phase = (t + Offset) mod (GreenDuration
+ RedDuration)` and returns `phase < GreenDuration`. Two signals are
"always opposite" (FR-005) simply by being configured with the same
`GreenDuration`/`RedDuration` and an `Offset` differing by exactly one
`GreenDuration` (or, symmetrically, one `RedDuration`) — there is no
"linked to" field or runtime relationship between two `Signal` values.

**Rationale**: This is simpler than it looks: two independent signals
that happen to share a cycle length and a complementary offset are
*already* always in opposite phase, by the definition of the phase
function — no additional bookkeeping is needed to keep them "in sync"
(there's nothing to keep in sync; they're each just independently
evaluating the same deterministic function of `t`). This also
generalizes for free to configurations FR-005 doesn't strictly require
(e.g. three signals at 120°-offset thirds of a cycle) without the data
model needing to change.

**Alternatives considered**: A `PairedWith *Signal` field with an
"opposite" flag, resolved at evaluation time — rejected as unnecessary
indirection for something plain arithmetic on two independently-configured
values already achieves.

## 5. Demand gains a time dimension: agents arrive over the run, not all at once

**Decision**: `queuesim.Demand` is `{Origin, Destination string, Count
int, ArrivalInterval float64}` — one agent spawns every `ArrivalInterval`
simulated seconds, up to `Count` total. This differs from `simulation.
Demand` (feature 003/004), which has no time dimension at all (every
agent is already "present" from the start of that equilibrium computation).

**Rationale**: The spec's own Acceptance Scenarios describe "a steady
stream of arriving agents" — a run with no arrival process wouldn't have
anything meaningful for a queue to do (all agents would need to route
through a signal-controlled edge in the very first tick, an unrealistic
burst rather than a stream). A steady arrival process is also what makes
User Story 2's central question — does one road end up persistently
worse, does it oscillate, or does it balance — actually possible to
observe over a run's duration, rather than settling near-instantly the
way an equilibrium computation does.

**Alternatives considered**: Reusing `simulation.Demand` unmodified
(no time dimension, so all agents "arrive" at `t=0`) — rejected; it
would make queuing degenerate into a single, uninteresting burst rather
than a sustained flow, and wouldn't exercise the actual research question
User Story 2 asks about.

## 6. Reporting: samples over time, not a single final number

**Decision**: `RunResult` holds `QueueSamples []QueueSample`
(`{Time, SignalID, Length}`, recorded once per tick per signal) and
`Agents []AgentReport` (`{DemandIndex, TravelTime, WaitTime, Arrived
bool}`), one entry per spawned agent.

**Rationale**: FR-004/FR-006/SC-003 all require seeing *how a quantity
changed over the run*, not just its end state — the entire point of User
Story 2 is to see the shape of the queue-length curve for each road, which
requires the time series, not a final snapshot. Per-agent wait vs. travel
time (kept separate, FR-006) is what lets a caller compute both individual
experience and aggregate statistics without the report making that choice
for them.

**Alternatives considered**: Reporting only summary statistics (e.g. max
queue length, average wait) — rejected; SC-003 explicitly requires being
able to tell *from the data* whether congestion persisted, oscillated, or
balanced, which a single summary number can't distinguish (a road that's
briefly very congested twice and a road that's moderately congested
throughout could have the same average).
