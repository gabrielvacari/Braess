# Phase 0 Research: Shortest-Path Agents

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per constitution
Principle V.

## 1. Algorithm: Dijkstra via `container/heap`

**Decision**: Use Dijkstra's algorithm with a binary min-heap
(`container/heap` from the standard library) to compute the minimum-total-
travel-time route from origin to destination.

**Rationale**: Every edge's travel time at the baseline (volume = 0) is a
non-negative free-flow time (see decision #3), which is exactly Dijkstra's
precondition. The networks in scope are tens of nodes/edges — `O((V+E) log
V)` with a binary heap is far more than fast enough for SC-002 (< 5s for
10 nodes/15 edges), and the standard library already provides the heap
primitive, so no third-party graph library is justified.

**Alternatives considered**: Bellman-Ford — handles negative weights, which
this domain doesn't have (a road can't have negative travel time); rejected
as unnecessary generality. A third-party graph library (e.g. `gonum/graph`)
— rejected per the same reasoning as feature 001's research.md: it would
be a dependency this small, self-contained engine doesn't need.

## 2. Package boundary: new `agent` package, `graph` untouched

**Decision**: Create a new top-level package `agent` that depends on
`graph` (and the standard library) and nothing else. `agent` builds its own
adjacency structure from `graph.Graph.Nodes()`/`Edges()` rather than
`graph` growing pathfinding methods itself.

**Rationale**: AGENTS.md's own architecture explicitly separates "1. Graph
(simulation engine)" from "2. Agents (cars)" as independent concerns; the
graph "should not depend on the UI" and, by the same logic, shouldn't need
to know how routes are chosen — that's the agents' decision-making, not the
network's structure. Keeping them as separate packages with a one-way
dependency means feature 001 stays completely stable while this feature is
purely additive, and later swapping or extending the routing strategy
(Phase 3's congestion-aware routing) won't touch `graph` at all.

**Alternatives considered**: Adding a `Graph.ShortestPath(...)` method
directly on `graph.Graph` — rejected because it would blur "what the road
network *is*" (Phase 1's concern) with "how an agent *decides*" (this
phase's and Phase 3's concern), and would make the graph package's surface
grow indefinitely as routing strategies evolve.

## 3. Baseline volume as a parameter, not a hardcoded zero

**Decision**: `ShortestRoute(g *graph.Graph, from, to string, volume float64) (Route, error)`
takes the traffic volume it evaluates edges at as an explicit parameter.
`Agent.ComputeRoute` (this phase's only caller) always passes `0`.

**Rationale**: The spec's Assumptions section is explicit that "fixed/zero
traffic volume" is a Phase 2 simplification, not a permanent rule — Phase 3
introduces real congestion. Taking `volume` as a parameter now costs
nothing (this phase's single caller just passes `0`) and avoids reshaping
the function's signature later; hardcoding `0` inside the algorithm would
require every future caller to be rewritten in Phase 3 for no benefit today.

**Alternatives considered**: Hardcoding the baseline volume to `0` inside
`ShortestRoute` — rejected for the forward-compatibility reason above. This
is a small, free choice today, not scope creep: this feature's behavior is
still fully pinned to volume `0` by `Agent.ComputeRoute`, per FR-001.

## 4. Free-flow times are assumed non-negative

**Decision**: Document, as an implementation assumption backing decision
#1, that every `TravelTimeFunc` used with `agent` returns a non-negative
value at `volume == 0` (as `graph`'s own `Linear`/`Constant` constructors
naturally do for any reasonable free-flow time).

**Rationale**: Dijkstra requires non-negative edge weights. `graph`
(feature 001) doesn't itself forbid a negative free-flow time — nothing in
its contract needs to, since "free-flow time" being non-negative is a
domain fact (a road never takes negative time), not a graph-structure
concern. Enforcing it in `agent` instead of `graph` keeps that concern where
it's actually used.

**Alternatives considered**: Validating non-negativity inside `graph`
itself — rejected; it isn't `graph`'s job to enforce a routing
algorithm's precondition, and doing so there would be exactly the kind of
scope-blur decision #2 rejects.
