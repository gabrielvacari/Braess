# Phase 0 Research: Congestion and Braess's Paradox

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per constitution
Principle V. This is the feature where the theoretical grounding matters
most, so several decisions here are backed by results from congestion-game
theory rather than just engineering taste.

## 1. Best-response dynamics converges here — this isn't a heuristic hope

**Decision**: Model the population as an atomic congestion game (Rosenthal,
1973): each of N identical agents picks one route, and a road's cost
depends only on how many agents are currently using it. Find a stable
assignment via best-response dynamics: repeatedly let one agent at a time
switch to a strictly better route given everyone else's current choice,
stopping when nobody can improve.

**Rationale**: This is not a heuristic that might happen to settle down —
Rosenthal proved that exactly this class of game is a *potential game*: a
single scalar function (the sum, over every road, of the cost each of its
1st, 2nd, ... nth user would pay) strictly decreases every time an agent
makes a strictly improving switch. Since the number of possible assignments
is finite, the potential can only decrease finitely many times, so
best-response dynamics is *guaranteed* to reach a pure Nash equilibrium in
finite steps — for any network, any (non-decreasing) travel-time functions,
any population size. That guarantee is exactly FR-005's requirement, not
just an empirical tendency, and it's why User Story 2's "reaches a known
stable split" is expected to hold in general, not just for the specific
networks tested.

**Alternatives considered**: A closed-form Wardrop (continuous-flow)
equilibrium solver — rejected: the population here is explicitly discrete
agents (matching AGENTS.md's "each car is an agent"), and the atomic/
discrete version is both truer to the project's premise and the one with
the clean finite-convergence guarantee above. Simulated annealing or other
general-purpose optimizers — rejected as unnecessary machinery for a
problem this well-structured.

## 2. Switch only on strict improvement

**Decision**: During refinement, an agent switches away from its current
route only if an alternative is strictly cheaper (beyond a small epsilon,
`1e-9`, to absorb floating-point noise) — never merely equal.

**Rationale**: This directly resolves the spec's tied-routes Edge Case.
Two routes tied at the same cost would otherwise cause agents to flip
back and forth between them forever with no actual improvement — the
exact "never converges" failure mode the spec calls out. Requiring strict
improvement makes ties inherently stable (nobody switches away from a tie)
with no extra bookkeeping, and it is also required for decision #1's
potential-function argument to guarantee termination — a dynamic that
allows sideways moves on ties isn't guaranteed to terminate the same way.

**Alternatives considered**: Detecting convergence by aggregate travel time
no longer changing (even if individual assignments still flip) — rejected
as an unnecessary second mechanism once switch-on-strict-improvement
already prevents the flipping it would be there to tolerate.

## 3. Incremental loading, then refine

**Decision**: Build the initial assignment by adding agents one at a time,
each computing its best route against the volume accumulated by the agents
already placed (`agent.ShortestRouteAtVolumes`, decision #4) — then run
best-response refinement rounds (decision #1) on top of that starting
point.

**Rationale**: Starting from "everyone independently picks their free-flow
shortest route" (Phase 2's behavior) would have every agent pile onto the
same edge in a first pass, which is a worse and less realistic starting
point that would need many more refinement rounds to unwind. Loading
agents one at a time against real, growing volume is a standard, simple
traffic-assignment warm start, and gets close to equilibrium immediately
for the symmetric and near-symmetric cases this project cares about.

**Alternatives considered**: Starting every agent from a shared free-flow
route (ignoring volume during init) — rejected per the above; it's not
wrong, just slower to converge and less representative of "agents reacting
to current congestion" from the very first agent onward, which is the
project's stated premise.

## 4. `agent` grows a per-edge-volume routing entry point

**Decision**: Add `agent.ShortestRouteAtVolumes(g, from, to, volumes
map[string]float64)`, factoring the existing Dijkstra implementation into a
shared internal core parameterized by a `func(edgeID string) float64`
volume lookup. The existing `ShortestRoute(g, from, to, volume float64)`
becomes a one-line wrapper (`uniform volume for every edge`) over that same
core — its signature and behavior are unchanged, so feature 002's tests and
callers (`Agent.ComputeRoute`) keep working exactly as before.

**Rationale**: Real congestion means different roads carry different
traffic — Phase 2's single flat `volume float64` parameter, applied
identically to every edge, cannot express that. Feature 002's own contract
(`contracts/agent-api.md`, "Stability") explicitly anticipated this:
"expected to grow, not break, in later phases." Factoring out one shared
Dijkstra core (rather than writing a second, separate pathfinding
implementation for `simulation` to use) keeps there being exactly one
place the algorithm can be wrong.

**Alternatives considered**: Writing pathfinding logic directly inside
`simulation` instead of extending `agent` — rejected; it would duplicate
Dijkstra and blur the layering decision 001/002 already established
(`agent` owns "how a route is chosen"; `simulation` owns "how many agents,
reacting to each other, end up choosing what").

## 5. `MaxRounds` is explicit, and `0` is a well-defined way to test FR-006

**Decision**: `Population.MaxRounds` is a required field (no hidden
"0 means use a default" magic) — the CLI sets it to `simulation.
DefaultMaxRounds`, and callers who want a different bound set it directly.
Setting it to `0` means "run the initial incremental loading only, perform
zero refinement rounds," which always reports `Converged: false` (removed
stability was never actually checked, so reporting anything else would be
overclaiming).

**Rationale**: This gives FR-006 ("not fully converged" reporting) a
precise, deterministic way to be exercised in tests (SC-004) without
needing to hand-craft a network that pathologically resists convergence —
which decision #1's finite-convergence guarantee makes hard to construct
on purpose. It's also just honest: with zero refinement rounds, stability
was never verified, so `Converged: false` is the conservative, correct
thing to report — not a test-only special case bolted on afterward.

**Alternatives considered**: `0` meaning "use `DefaultMaxRounds`" (the more
common Go convention for optional int fields) — rejected here specifically
because it would remove the one deterministic lever for testing the
"round limit reached" path; an explicit, always-honored `MaxRounds` is a
small enough field that requiring callers to set it is not a real burden.

## 6. The classic Braess network, at its literature scale

**Decision**: User Story 3's fixed network reproduces the textbook Braess's
Paradox example exactly: nodes S, A, B, T; edges S→A and B→T with
`Linear(0, 0.01)` (time = volume/100); edges S→B and A→T with `Constant(45)`;
4000 agents from S to T. The "before" network omits the extra A→B edge;
the "after" network adds it as `Constant(0)`.

**Rationale**: This is the well-known result every reader of Braess's
Paradox literature will recognize: without the extra road, the population
splits 2000/2000 and everyone pays 65; with it, decision #1's potential-
game guarantee predicts (and this plan's design was checked against) a
unique equilibrium where all 4000 agents route through the new edge and
everyone pays 80 — strictly worse for 100% of drivers, despite the network
gaining a free road. Reproducing the exact textbook numbers is strong,
checkable evidence the engine's congestion model is correct, not just
"some numbers went up."

**Alternatives considered**: A scaled-down, invented network with smaller
numbers — rejected for User Story 3 specifically: matching the famous
65 → 80 result is far more convincing (and independently checkable against
outside sources) than an arbitrary example would be. Smaller invented
networks are still used for the general-mechanism tests in
`run_test.go` (SC-001), where matching the literature isn't the point.
