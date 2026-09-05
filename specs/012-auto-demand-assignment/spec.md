# Feature Specification: Automatic House Demand

**Feature Branch**: `012-auto-demand-assignment`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "In Signals mode, remove manual demand configuration entirely (no more choosing origin/destination/count/arrival-interval by hand) and replace it with an automatic rule: every house node always generates a demand of exactly 12 agents arriving every 0.8 simulated seconds, targeting a company node. The code — not the user — decides which company each house's agents go to. This assignment does not need to be balanced/equal across companies, but it must never send every house to the same single company when more than one company exists in the network (i.e., when there are 2 or more company nodes, at least 2 distinct companies must receive traffic). This only affects Signals mode's demand handling — it does not change the Equilibrium mode's demand form, does not change queuesim's engine, and does not change the wire format of a time-based demand (origin, destination, count, arrivalInterval) sent to the server — only how those values get decided (automatically, in code, from the drawn network) instead of being entered by the user."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Every house sends traffic automatically (Priority: P1)

A person building a Signals-mode network draws houses and companies and
roads between them, but never has to declare who goes where or how
many — every house automatically sends 12 agents every 0.8 simulated
seconds toward some company, the moment the network is run.

**Why this priority**: This is the core of the request — the manual
demand form is explicitly what should disappear; everything else
(which company gets picked, staying spread across more than one) only
matters once traffic is flowing automatically at all.

**Independent Test**: Draw two houses and one company, connect them
with roads, run the scenario, and confirm both houses' agents appear
without ever declaring a demand by hand.

**Acceptance Scenarios**:

1. **Given** a network with at least one house and one company,
   **When** it's run, **Then** every house's agents appear in the
   result without the person having configured any demand.
2. **Given** a network with a house but no company at all, **When** it's
   run, **Then** that house generates no demand (there is nowhere for it
   to be assigned to go).

---

### User Story 2 - Traffic spreads across more than one company (Priority: P1)

When a network has more than one company, a person can see that houses'
traffic doesn't all funnel into a single one of them — the assignment of
which company each house targets is decided automatically, and
guarantees more than one company is actually used.

**Why this priority**: Explicitly called out as a hard requirement, not
a nice-to-have — a network with multiple companies where everything
still piles onto just one wouldn't demonstrate anything different from
having drawn only one company, defeating the purpose of drawing several.

**Independent Test**: Draw at least two houses and two companies, run
the scenario, and confirm at least two distinct companies each receive
at least one house's traffic.

**Acceptance Scenarios**:

1. **Given** a network with 2 or more companies, **When** it's run,
   **Then** at least 2 distinct companies each have at least one house
   assigned to them.
2. **Given** a network with exactly 1 company, **When** it's run,
   **Then** every house is assigned to that one company (there is no
   other option).
3. **Given** the same network run twice without being edited in
   between, **When** comparing the two runs, **Then** each house is
   assigned to the same company both times (the assignment doesn't
   change on its own between runs).

---

### Edge Cases

- A network with houses but zero companies: no demand exists to run;
  the person sees a clear explanation rather than a confusing error.
- A network with companies but zero houses: nothing to run, same
  treatment.
- A house with no route at all to the company it was assigned: this
  surfaces as the same "no route" error this mode already reports for a
  manually-declared demand — the assignment itself doesn't try to avoid
  unreachable companies.
- Adding or removing a house or company changes what gets assigned;
  the person is not expected to reason about the assignment rule
  themselves — Success Criteria SC-002 covers making the choice visible
  instead.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Signals mode MUST NOT offer any way to manually declare a
  demand (no origin, destination, count, or arrival-interval entry).
- **FR-002**: Every house node MUST automatically generate a demand of
  exactly 12 agents, arriving every 0.8 simulated seconds, the moment
  the network is run.
- **FR-003**: Each house's demand MUST target exactly one company node,
  chosen automatically rather than by the person building the network.
- **FR-004**: When a network has 2 or more company nodes, the
  automatic assignment MUST result in at least 2 distinct companies
  each receiving at least one house's traffic.
- **FR-005**: The automatic assignment does not need to distribute
  houses evenly across companies — only the minimum spread required by
  FR-004.
- **FR-006**: Running the same, unedited network more than once MUST
  produce the same house-to-company assignment every time.
- **FR-007**: A house with no company available to target (a network
  with zero companies) MUST generate no demand, rather than erroring or
  targeting nothing meaningful.
- **FR-008**: This change MUST NOT alter the Equilibrium mode's demand
  form or behavior in any way.
- **FR-009**: This change MUST NOT alter `queuesim`'s engine or the
  wire format of a time-based demand sent to the server — only how
  Signals mode decides what to send.
- **FR-010**: The person building the network MUST be able to see which
  company each house has been assigned to, since it's no longer
  something they chose themselves.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A person can go from an empty Signals-mode canvas to a
  running scenario with visible traffic without ever filling out a
  demand form.
- **SC-002**: A person can identify which company any given house is
  sending its agents to, without needing to infer it from the run
  result alone.
- **SC-003**: In a network with 2 or more companies, every run shows
  traffic reaching at least 2 distinct companies.
- **SC-004**: Repeating a run of the same, unedited network never
  changes which company a house was sending its agents to.

## Assumptions

- Only house nodes generate demand; company and intersection nodes
  never do, matching the existing meaning of those node types.
- The 12-agents-per-0.8-seconds rate is a fixed constant for every
  house, not configurable per house or globally — consistent with the
  request that these numbers no longer be entered by hand at all.
- "Does not need to be balanced" is read as: the assignment rule is free
  to be as simple as it likes (e.g. deterministic round-robin) as long
  as it satisfies FR-004; it is not required to account for company
  capacity, distance, or any other weighting.
- The assignment is recomputed from the current network state every
  time it's needed (e.g., after an edit), not stored as separate data a
  person could hand-override — there is no manual escape hatch back to
  per-demand configuration.
- FR-004's "at least 2 distinct companies receive traffic" is only
  achievable when there are at least 2 houses to distribute in the
  first place, since each house's traffic goes to exactly one company
  (FR-003). With exactly 1 house and 2+ companies, that single house's
  traffic reaching only 1 company is expected, not a violation of
  FR-004.
