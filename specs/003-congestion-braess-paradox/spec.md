# Feature Specification: Congestion and Braess's Paradox

**Feature Branch**: `003-congestion-braess-paradox`

**Created**: 2026-09-04

**Status**: Draft

**Input**: User description: "Phase 3 of the roadmap (per AGENTS.md and constitution Principle II): congestion — travel time increasing with actual traffic volume from a population of agents sharing the fixed origin/destination pair from Phase 2, with agents iteratively re-routing in response to current congestion until the assignment stabilizes. Validate via console/log output that Braess's Paradox appears when a road is added to the classic four-node network — this is the project's central experimental payoff and the roadmap's explicit gate before any frontend work starts. Still no UI; still a single fixed origin/destination pair (multiple distinct pairs are Phase 4)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Travel time reflects real traffic (Priority: P1)

A developer runs several agents sharing one road network and sees each
road's travel time computed from how many of those agents are actually on
it right now — not from an arbitrary number chosen by whoever is calling
the system.

**Why this priority**: This is congestion's foundational mechanism. Without
travel time genuinely depending on how many agents are actually present,
"congestion" isn't real — it's still Phase 2's fixed baseline with extra
steps. Every later capability in this feature depends on this being true.

**Independent Test**: Can be fully tested by assigning a known set of
agents to specific roads and checking the computed travel times match what
each road's travel-time function returns for those exact counts, via
`go test` — no iterative re-routing or the Braess network needed yet.

**Acceptance Scenarios**:

1. **Given** a road currently used by a known number of agents, **When**
   its travel time is computed, **Then** it matches that road's
   travel-time function evaluated at that exact count.
2. **Given** two agents move from one road to another, **When** travel
   times are recomputed, **Then** the road they left reflects one fewer
   agent and the road they joined reflects one more.

---

### User Story 2 - Agents settle into a stable assignment (Priority: P2)

Many agents, all sharing the same origin and destination, repeatedly
reconsider their route in light of current congestion. Round after round,
fewer and fewer of them find a faster alternative, until the population
settles into a stable pattern — or the system gives up after a bounded
number of attempts and says so.

**Why this priority**: A single, one-shot routing pass against whatever
volume happens to exist isn't how real congestion behaves — it emerges
from everyone reacting to everyone else, repeatedly. This is the direct
implementation of the project's premise ("recalculates or reacts to
current road congestion") and of constitution Principle III at scale,
rather than the two- or three-agent toy case from Phase 2.

**Independent Test**: Can be fully tested by running the process on a
small network with a known correct equilibrium split between two routes
and confirming it converges to that split within a small tolerance, via
`go test` — no UI or the specific Braess network needed.

**Acceptance Scenarios**:

1. **Given** a population of agents re-routing round after round on a
   network with a known stable split between two routes, **When** the
   process runs, **Then** it reaches that split within a small tolerance
   and stops re-iterating once it has.
2. **Given** a network and population where no stable split emerges
   before a maximum number of rounds is reached, **When** the process
   runs, **Then** it stops at that round and clearly reports the result
   as not fully settled, rather than looping forever or claiming success.

---

### User Story 3 - Reproduce Braess's Paradox in the console (Priority: P3)

A developer runs the terminal command against the classic four-node
network twice — once as-is, once with one extra road added — and sees
printed proof that adding the road made the average trip *slower* for
everyone, not faster.

**Why this priority**: This is the actual reason the project exists, and
the roadmap's explicit gate: AGENTS.md requires validating the paradox in
text/logs before a single line of frontend code is written. Everything
built so far (the graph engine, shortest-path agents, real congestion) was
in service of being able to run this specific comparison.

**Independent Test**: Can be fully tested by running the terminal command
and checking that the printed total/average travel time for the
"with the extra road" run is strictly worse than the "without it" run —
no UI required.

**Acceptance Scenarios**:

1. **Given** the classic four-node network without its extra connecting
   road, **When** the population is run to a stable assignment, **Then**
   the command prints the resulting total and average travel time.
2. **Given** the same network with the extra road added, **When** the
   population is run to a stable assignment, **Then** the command prints a
   total and average travel time that is strictly worse than the
   without-the-extra-road result, and the output makes that comparison
   explicit.

---

### Edge Cases

- What happens when zero agents are simulated? The system must report a
  trivial result (every road at its free-flow time, zero total travel
  time) rather than erroring.
- What happens when no route exists at all between the origin and
  destination for some agents? The system must report that clearly
  (consistent with Phase 2's "no route exists" behavior) rather than
  crashing or silently dropping those agents.
- What happens when two routes end up exactly tied, causing agents to keep
  switching between them every round without the individual assignment
  ever settling? The system must still recognize the population-level
  outcome (e.g., total or average travel time) as stable and stop, rather
  than treating that oscillation as "never converges."
- What happens when the maximum number of rounds is reached without the
  assignment settling? The reported result must be clearly marked as not
  fully converged, not presented as a final equilibrium.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST compute a road's travel time using the
  actual number of agents currently assigned to it, not a caller-supplied
  arbitrary baseline.
- **FR-002**: The system MUST support running a population of many agents
  — not just one or a few — sharing a single fixed origin/destination
  pair simultaneously.
- **FR-003**: The system MUST let each agent recompute its route based on
  current traffic volumes, independently of how other agents' individual
  decisions are computed (continuing Phase 2's Principle III behavior at
  population scale).
- **FR-004**: The system MUST repeat the recomputation process across
  multiple rounds, so agents' choices can react to each other's previous
  choices.
- **FR-005**: The system MUST detect when the population's assignment (or
  its aggregate outcome) has stabilized and stop iterating at that point.
- **FR-006**: The system MUST stop after a bounded maximum number of
  rounds even if stabilization is not detected, and clearly mark the
  result as not fully converged in that case.
- **FR-007**: The system MUST report the total and average travel time
  experienced across all agents at the final assignment.
- **FR-008**: The system MUST be able to run the same population against
  two versions of a network — one with, one without, an added road — and
  report both resulting total/average travel times so they can be
  compared directly.
- **FR-009**: The system MUST remain operable from the same standalone
  terminal entry point established in earlier phases, continuing to
  require zero UI/frontend dependency.
- **FR-010**: The system MUST handle a population of zero agents without
  error, reporting every road at its free-flow time and zero total travel
  time.

### Key Entities

- **Population**: The set of agents sharing one fixed origin/destination
  pair, together with the iterative process that re-routes them against
  current congestion until their assignment stabilizes or a round limit
  is reached.
- **Edge Volume**: The count of agents currently assigned to a given road
  at a point in that iterative process.
- **Assignment Result**: The final per-agent routes, the aggregate total
  and average travel time, and whether the process converged or instead
  stopped at the round limit.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Across at least 3 distinct networks with a known correct
  stable route split, the population's final assignment matches that split
  within a small tolerance.
- **SC-002**: Running the classic four-node network with its extra
  connecting road added produces a printed total/average travel time
  strictly greater than running it without that road — reproducing
  Braess's Paradox — verified via the terminal command.
- **SC-003**: A population of at least 100 agents on a network of at most
  10 nodes reaches a final result and prints it via the terminal command
  in under 5 seconds.
- **SC-004**: When the process does not detect a stable assignment within
  the maximum round count, the reported result is clearly marked as not
  fully converged in 100% of such cases, with no crash.

## Assumptions

- "Population" means many agent instances sharing the single fixed
  origin/destination pair established in Phase 2 — multiple distinct
  origin/destination pairs remain out of scope until roadmap Phase 4.
- The iterative re-routing process is a simplification of real traffic
  assignment (agents making best-response decisions round after round)
  rather than a closed-form equilibrium solver; the exact algorithmic
  detail (e.g., how many agents reconsider per round) is a planning
  decision, as long as it satisfies the stabilization/round-limit
  requirements above.
- The specific classic Braess network (four nodes, a mix of
  volume-sensitive and constant-time roads) used in User Story 3 is a
  fixed, hand-built example for this phase — not user-configurable;
  interactive graph editing arrives in roadmap Phase 5 (frontend).
- No new travel-time function shapes are required beyond what Phase 1
  already provides (`Linear`, `Constant`) — the classic Braess example
  specifically needs both side by side, which Phase 1 already supports.
- This phase remains in-memory, CLI/log-verified only, consistent with
  Phases 1 and 2.
