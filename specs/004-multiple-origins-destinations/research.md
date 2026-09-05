# Phase 0 Research: Multiple Origins and Destinations

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per constitution
Principle V.

## 1. Generalize the shared core, don't duplicate it

**Decision**: Generalize `loadIncrementally` and `refine` (feature 003) to
operate over a flat `[]odPair` — one entry per agent, each carrying its own
`origin`, `destination`, and which `Demand` index it belongs to — instead
of assuming every agent shares one fixed pair. `Run` (feature 003's public
API) becomes a thin wrapper: it builds an `[]odPair` where every entry
repeats `Population.Origin`/`Destination`, and calls the same core.

**Rationale**: This is the identical pattern already used for
`agent.ShortestRoute`/`ShortestRouteAtVolumes` in feature 003 — the
generalization here is genuinely small (each agent looks up its own
origin/destination from a slice instead of reading one shared field), and
feature 003 already has strong, specific regression coverage for exactly
this code path: its equilibrium tests and, especially, the classic Braess
network's exact 65 → 80 result. Re-running those unmodified after the
refactor is a much stronger correctness signal than either a hand-audit of
duplicated code or a leap of faith that two copies stay in sync over time.

**Alternatives considered**: Duplicating `loadIncrementally`/`refine` into
new multi-demand-only versions, leaving feature 003's originals completely
untouched — considered specifically to minimize risk to a "done" feature,
but rejected: the two copies would differ by one line (where the
origin/destination comes from) and everything else — the marginal-cost
math, the epsilon comparison, the convergence loop — would need to be kept
in sync by hand forever after. That ongoing duplication risk is worse than
the one-time, verifiable risk of a small, tested generalization.

## 2. `RunDemands` is a new function, `Run` is untouched

**Decision**: Add `RunDemands(g *graph.Graph, mp MultiPopulation)
(MultiAssignmentResult, error)` as a new, separate entry point.
`Population`, `Run`, and `AssignmentResult` keep their exact existing
signatures and field sets from feature 003.

**Rationale**: Feature 003's own contract (`contracts/simulation-api.md`,
"Stability") anticipated exactly this: "Phase 4's multiple distinct
origin/destination pairs will need a richer population shape, expected to
be a new type ... rather than a breaking change to `Run`'s signature."
Every existing caller of `Run` (feature 003's tests, `cmd/graphcli`'s
Braess demonstration) keeps working with zero changes.

**Alternatives considered**: Adding a `Demands []Demand` field directly to
`Population` and reinterpreting `Origin`/`Destination`/`Size` as a
shorthand for a single-element case — rejected as a confusing dual-mode
struct (which fields are authoritative depends on which are set); a
separate `MultiPopulation` type is unambiguous.

## 3. Reporting: one overall result, plus one result per demand

**Decision**: `MultiAssignmentResult` holds the overall total/average
travel time across every agent, `Converged`, and `Rounds`, plus
`PerDemand []DemandResult` — one entry per input `Demand`, in the same
order, each with that demand's own routes, total, and average.

**Rationale**: Directly implements FR-005/FR-006 and User Story 2: a
shared network can be fine for one flow of traffic and worse for another,
and hiding that behind a single aggregate number would defeat the point of
having distinct demands at all. Keeping `PerDemand` in input order (rather
than, say, a map keyed by origin/destination) keeps a demand with zero
agents (FR-007) or a duplicate origin/destination pair (spec Edge Cases)
trivially representable — every declared `Demand` gets exactly one
`DemandResult`, regardless of size or duplication.

**Alternatives considered**: Only reporting the overall aggregate and
leaving per-demand breakdown to the caller (re-deriving it from `Routes`
by filtering on origin/destination) — rejected; the caller would need to
duplicate `RunDemands`'s own bookkeeping to get per-demand numbers, and
two demands sharing an origin/destination (spec Edge Cases) couldn't be
told apart this way even with route filtering.

## 4. Identifying which demand failed

**Decision**: `loadIncrementally`/`refine` return which agent index (if
any) a failure occurred at. `RunDemands` maps that back to a `Demand`
index via each agent's tagged `odPair.demandIndex` and wraps the
underlying error with that context: `"demand %d (origin=%q
destination=%q): %w"`. `errors.Is(err, agent.ErrNoRoute)` still succeeds
through the wrap.

**Rationale**: FR-008 explicitly requires identifying *which* declared
pair is affected when one has no route — the underlying
`agent.ErrNoRoute`-wrapping error already names the origin/destination
(feature 002), but not which of possibly several declared `Demand`s (by
index, in the caller's own list) it corresponds to, which matters most
when two demands happen to share an origin/destination.

**Alternatives considered**: Returning a partial `MultiAssignmentResult`
for the demands that *did* succeed, alongside an error for the one that
didn't — rejected, consistent with every earlier phase's fail-fast
philosophy (Phase 1-3 all propagate rather than silently produce a
degraded result); a partial multi-demand result would be a new, riskier
kind of "silent partial success" this project has deliberately avoided
everywhere else.

## 5. Load order across demands doesn't affect correctness

**Decision**: Build the flat `[]odPair` by iterating demands in the order
given, and within each demand, one agent at a time (i.e., all of demand 0's
agents load before any of demand 1's). No interleaving is implemented.

**Rationale**: Rosenthal's potential-game convergence guarantee (feature
003, decision #1) does not depend on the identity or arrival order of
players — it holds for *any* starting assignment, since best-response
refinement afterward strictly decreases the same potential function
regardless of how loading got there. Load order can only affect *how many*
refinement rounds are needed, never *whether* a valid equilibrium is
eventually found — and at this project's test/demo scale (SC-003's ≤15
nodes, ≥100 agents), that difference is immaterial.

**Alternatives considered**: Round-robin interleaving agents across
demands during loading (one from each demand, repeated) — a plausible
alternative warm start, rejected only for being unnecessary extra
complexity: it might shave a few refinement rounds off in some cases, but
nothing in this feature's success criteria needs that, and the simpler,
order-preserving approach is easier to reason about and to write
deterministic tests against.
