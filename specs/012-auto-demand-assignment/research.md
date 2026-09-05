# Phase 0 Research: Automatic House Demand

## Decision 1: Deterministic round-robin assignment, not random or hash-based

**Decision**: Assign the i-th house (in the order it appears in
`nodes`, i.e. creation order) to `companies[i % companies.length]`.

**Rationale**: Guarantees FR-004 by construction whenever there are 2+
houses and 2+ companies — the first two houses land on two different
companies, full stop, with no probability involved. It also satisfies
FR-006 (repeated runs of the same network agree) trivially, since it's
a pure function of the node list, which doesn't change between runs of
the same network.

**Alternatives considered**: A random assignment (re-rolled per run) —
rejected: fails FR-006 outright (a second run could reassign houses). A
hash of the house's id modulo company count — rejected: for a small
number of houses, nothing stops a hash from mapping several houses to
the same bucket by chance, which could violate FR-004 for exactly the
small networks (2-3 houses) this hand-drawn tool mostly produces; a
plain index-based round-robin has no such risk.

## Decision 2: A new `model/autoDemand.ts`, not an addition to `queueNetwork.ts`

**Decision**: The assignment function and the DTO-building function
(`assignHousesToCompanies`, `autoTimeDemands`) live in a new file, not
inside `queueNetwork.ts`.

**Rationale**: `queueNetwork.ts`'s own file header already states it
"contains no pathfinding or queuing logic" — it is deliberately just
state CRUD (add/remove a node/road). Deciding which company a house's
traffic goes to is a derived computation over that state, not an edit
to it, so it belongs in its own module, matching how `queuePlayback.ts`
and `roadVisuals.ts` already sit alongside `queueNetwork.ts` rather than
inside it.

**Alternatives considered**: Adding the function to `queueNetwork.ts`
directly — rejected as contradicting that file's own stated scope.

## Decision 3: The assignment is surfaced read-only in `SignalToolbar`, not a separate panel

**Decision**: `SignalToolbar`'s demand-form section is replaced by a
short, read-only list ("House 1 → Company 2", etc.), computed by
calling `assignHousesToCompanies(nodes)` directly from the `nodes` prop
`SignalToolbar` already receives — no new prop from `SignalApp`.

**Rationale**: FR-010 requires visibility, and the demand form's old
location is the natural place for it — a person configuring the network
looks there already. Computing it directly from the existing `nodes`
prop avoids introducing a new prop just to pass the same derivable
value down.

**Alternatives considered**: Showing the assignment only after a run, in
the results panel — rejected: FR-010 is about seeing the assignment
*before* running, since it's no longer something the person chose
themselves and might want to check beforehand.

## Decision 4: `toQueueRunRequest` calls `autoTimeDemands` internally, no new parameter

**Decision**: `toQueueRunRequest(state, duration, tick)` keeps its exact
signature; internally, its `demands` field is now built by calling
`autoTimeDemands(state.nodes)` instead of mapping `state.demands`.

**Rationale**: Keeps the wire-building function's public shape stable
for its one caller (`SignalApp.tsx`) — the only thing that changes is
where the `demands` array's values come from, exactly matching the
spec's framing ("only how those values get decided... instead of being
entered by the user").

**Alternatives considered**: Passing a pre-computed `demands` array into
`toQueueRunRequest` — rejected as a distinction without a difference;
inlining the call keeps one fewer thing `SignalApp.tsx` has to
orchestrate.
