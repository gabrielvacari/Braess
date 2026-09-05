# Feature Specification: Multiple Origins and Destinations

**Feature Branch**: `004-multiple-origins-destinations`

**Created**: 2026-09-04

**Status**: Draft

**Input**: User description: "Phase 4 of the roadmap (per AGENTS.md and constitution Principle II): multiple distinct origins and destinations — several houses and companies sharing one road network at the same time, instead of Phase 2/3's single fixed origin/destination pair. Congestion must be genuinely shared: a road used by one house-to-company flow of traffic must affect any other flow crossing that same road. Results must be reportable per origin/destination pair as well as overall. Still no UI; still validated via console/log output, building on Phase 3's congestion and equilibrium-seeking."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Multiple flows share real congestion (Priority: P1)

A developer declares several distinct origin/destination pairs — each
with its own number of agents — on one road network, and the system
settles all of them into a stable assignment together: a road used by
agents going from one house to one company is exactly as congested for
agents going from a different house to a different company, if they both
end up using it.

**Why this priority**: This is the entire point of the phase. Without
congestion genuinely shared across pairs, "multiple houses and companies"
would just be several independent, non-interacting copies of Phase 3 —
not a real shared network, and not able to show how one flow of traffic
affects another the way real cities do.

**Independent Test**: Can be fully tested by building a small network
where two distinct origin/destination pairs are forced to cross the same
road, and confirming that changing one pair's agent count measurably
changes the other pair's resulting travel time — via `go test`, no UI
involved.

**Acceptance Scenarios**:

1. **Given** two distinct origin/destination pairs whose shortest routes
   share one common road, **When** both pairs' agents are run together,
   **Then** that road's travel time reflects agents from both pairs, not
   just whichever pair is considered first.
2. **Given** the same shared-road setup, **When** one pair's agent count
   is increased and the simulation is re-run, **Then** the other pair's
   reported average travel time changes accordingly, showing the two
   flows are genuinely interacting through the shared road.

---

### User Story 2 - Results reported per pair, not just overall (Priority: P2)

After running several origin/destination pairs together, a developer can
see each pair's own total and average travel time individually, in
addition to one combined number for everybody.

**Why this priority**: A shared network can treat different flows of
traffic differently — one pair's trip might stay fast while another's
gets much worse. A single aggregate number would hide exactly the kind of
asymmetric effect this phase exists to make visible.

**Independent Test**: Can be fully tested by running two distinct pairs
together and confirming each pair's own reported result reflects only
that pair's own agents and routes — via `go test`, no UI involved.

**Acceptance Scenarios**:

1. **Given** two distinct origin/destination pairs run together, **When**
   results are reported, **Then** each pair has its own total and average
   travel time, computed only from that pair's own agents.
2. **Given** the same run, **When** results are reported, **Then** one
   overall total and average travel time across every agent, regardless
   of pair, is also available.

---

### User Story 3 - Inspect multiple pairs from a terminal (Priority: P3)

A developer runs the existing standalone terminal command, now extended
to declare a few distinct houses and companies, and sees each pair's
result plus the overall result printed to the console.

**Why this priority**: Keeps the project's verify-in-text-before-any-UI
discipline going for this phase's new capability, extending the same
terminal command from earlier phases rather than inventing a new
verification path.

**Independent Test**: Can be fully tested by running the command and
checking its console output and exit code — no browser or UI required.

**Acceptance Scenarios**:

1. **Given** a small network with a few distinct houses and companies,
   **When** the terminal command runs, **Then** it prints each
   origin/destination pair's own total/average travel time and one
   overall total/average, exiting successfully.

---

### Edge Cases

- What happens when a declared origin/destination pair has zero agents?
  The system must report a trivial (zero-cost) result for that pair
  rather than erroring.
- What happens when two declared pairs share the exact same origin and
  destination? Their agents must still genuinely share congestion (same
  as if they were one larger pair for congestion purposes), while each
  declared pair's own result is still reported separately.
- What happens when a declared pair has no route at all between its
  origin and destination? The system must report that clearly, identifying
  which specific pair is affected, rather than silently dropping that pair
  or crashing the entire run.
- What happens when two pairs' shortest routes overlap on some roads but
  not others? Only the shared roads' travel times are affected by both
  pairs; each pair's own unshared roads behave as in Phase 3.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST accept multiple distinct origin/destination
  pairs, each with its own agent count, for a single simulation run.
- **FR-002**: The system MUST compute each road's travel time from the
  combined traffic of every agent currently using it, regardless of which
  origin/destination pair that agent belongs to.
- **FR-003**: The system MUST let every agent, from any pair, recompute
  its route independently based on current shared congestion, continuing
  Phase 2/3's independent-agent behavior at multi-pair scale.
- **FR-004**: The system MUST run the same load-then-stabilize process
  established in Phase 3 across every pair's agents together, until the
  whole shared assignment stabilizes or a round limit is reached.
- **FR-005**: The system MUST report, for each origin/destination pair
  separately, that pair's own total and average travel time.
- **FR-006**: The system MUST also report one overall total and average
  travel time across every agent, regardless of pair.
- **FR-007**: The system MUST handle an origin/destination pair with zero
  agents without error, reporting a trivial result for that pair.
- **FR-008**: The system MUST clearly report, identifying the specific
  pair affected, when no route exists for one of the declared pairs,
  rather than silently dropping that pair or crashing the whole run.
- **FR-009**: The system MUST remain operable from the same standalone
  terminal entry point established in earlier phases, continuing to
  require zero UI/frontend dependency.

### Key Entities

- **Demand**: One origin/destination pair together with how many agents
  travel between them in a given run.
- **Multi-Pair Result**: The outcome of running several demands together —
  one overall total/average travel time across everyone, each demand's own
  total/average travel time, and whether the shared assignment converged.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In a network where two distinct pairs' shortest routes share
  a common road, increasing one pair's agent count measurably changes the
  other pair's reported average travel time, demonstrating real shared
  congestion rather than isolated, non-interacting simulations.
- **SC-002**: Across at least 3 distinct multi-pair networks, the final
  assignment is stable: no agent, from any pair, could unilaterally switch
  to a different route and strictly improve.
- **SC-003**: A run with at least 3 distinct origin/destination pairs and
  a combined population of at least 100 agents, on a network of at most 15
  nodes, completes and prints via the terminal command in under 5 seconds.
- **SC-004**: When a declared pair has no route between its origin and
  destination, the system reports that clearly, identifying the specific
  pair, in 100% of such cases, with no crash.

## Assumptions

- Each agent has one fixed origin/destination pair for its entire trip,
  drawn from the declared set of demands — agents do not choose among
  multiple possible destinations; that remains out of scope.
- Phase 3's classic Braess demonstration remains a single-pair scenario
  and is not required to be re-derived for multiple pairs here; this
  phase's own success criteria use their own dedicated small test
  networks.
- This phase remains in-memory, CLI/log-verified only, consistent with
  Phases 1-3; interactive network editing and visualization are still
  roadmap Phase 5.
- No new node/edge concepts are needed beyond what Phase 1 already
  provides (House/Company/Intersection node types already exist) — this
  phase is about using more than one House and Company at once, not about
  adding new graph primitives.
