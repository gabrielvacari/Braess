# Phase 0 Research: Web Frontend

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per constitution
Principle V. This feature introduces more genuinely new technical ground
than any prior one, so there's more to justify here than usual.

## 1. A stateless server, not a session-holding one

**Decision**: `POST /api/run` takes the *entire* network and demand set in
one request body and returns the full result in one response. The server
holds no state between requests — no session, no in-memory or persisted
"current graph."

**Rationale**: The engine itself is already a pure function of
`(graph, population) -> result` (features 001-004) — `Run`/`RunDemands`
take a `*graph.Graph` value and return a result, nothing about them
assumes or needs server-side state. Mirroring that at the API level keeps
`cmd/server` a thin translation layer (JSON ↔ Go types) instead of a
second thing this project has to keep correct (session lifecycle, one
graph per user/tab, cleanup). It also matches the spec's own Assumptions
directly: single user, single network, no cross-session persistence — a
stateful server would be solving problems this phase explicitly doesn't
have.

**Alternatives considered**: A CRUD-style API (`POST /nodes`,
`DELETE /edges/:id`, server-side graph mutated in place, a separate
`POST /run` reading that server-side state) — rejected: it would require
session/identity handling for zero benefit here (the browser already
needs to hold the network as UI state to draw it — duplicating that state
server-side would be two sources of truth to keep in sync for a
single-user tool).

## 2. One endpoint speaks `Demand`/`MultiPopulation` always

**Decision**: The API has exactly one computation endpoint, and it always
takes a list of demands (`[]Demand` shape), even for the simplest
single-pair case — there is no separate "simple" endpoint mirroring
`Population`/`Run`.

**Rationale**: `simulation.RunDemands` already subsumes the single-pair
case (a `MultiPopulation` with one `Demand` behaves exactly like `Run`
with the equivalent `Population` — feature 004 kept `Run` only for
backward compatibility with feature 003's existing callers, not because
the two cases need different treatment). One endpoint means one request/
response shape for the frontend to reason about, regardless of how many
origin/destination pairs the user has drawn.

**Alternatives considered**: Exposing both `Run` and `RunDemands` as
separate endpoints — rejected as needless API surface; nothing in the
spec asks for a single-pair fast path, and the frontend always works in
terms of "the demands currently drawn," which is naturally zero-or-more.

## 3. The browser never decides a route — only edits and animates

**Decision**: All routing/congestion computation happens exclusively in
the Go engine via `POST /api/run`. The frontend's model layer
(`web/src/model/`) only holds and edits plain node/edge/demand data and
maps API responses to animation state — it contains no pathfinding, no
cost functions, nothing that mirrors `agent`/`simulation`.

**Rationale**: Constitution Principle III ("decentralized, selfish agent
behavior") is a claim about how the *simulation* works, not about the
tool built on top of it — but that claim only stays true if this feature
doesn't quietly grow a second, JavaScript implementation of "what route
would an agent pick" for, say, a snappier live preview. Keeping that
temptation explicitly out of scope here (rather than discovering it as an
accidental duplication later) is worth stating as a hard boundary now,
while the codebase is small enough that it's still easy to keep it that
way.

**Alternatives considered**: A client-side "preview" shortest-path
computation for instant feedback while drawing, before hitting the real
server — rejected for this phase; not requested by any user story, and it
would be exactly the kind of logic duplication decision #3 exists to rule
out. `POST /api/run`'s round trip is fast enough (see Technical Context)
that a client-side preview isn't needed to hit SC-002's 5-second budget.

## 4. Responses group agents by distinct route, not one entry per agent

**Decision**: `RunDemands`'s response, per demand, reports each *distinct*
route actually used and how many agents took it (`routeGroups: [{edges,
count, travelTime}]`), rather than one entry per individual agent.

**Rationale**: At this project's Braess-style scale (thousands of
agents), returning one JSON object per agent would inflate the response
for no visual benefit — agents sharing an identical route are visually
indistinguishable anyway, and an equilibrium typically settles into only
a handful of distinct routes (feature 003's classic example: exactly one,
at 4000 agents). Grouping keeps the payload bounded by the number of
*distinct* routes, not the population size, and gives the frontend
exactly what FR-002 needs: something to animate, at a count it can choose
to represent (e.g. a capped, proportionate number of icons per group)
without ever losing information about how the population actually split.

**Alternatives considered**: One route per agent — rejected per the above;
it doesn't serve any acceptance scenario better and scales worse.
Reporting only aggregate numbers with no route detail at all — rejected;
FR-002 explicitly needs something to animate along actual roads, not just
a total.

## 5. Frontend toolchain: Vite + TypeScript + Vitest, plain React state

**Decision**: The `web/` app uses Vite (dev server/bundler), TypeScript,
and Vitest for the pure-logic tests; state is plain React
`useState`/`useReducer` — no external state-management library.

**Rationale**: Vite is the current standard lightweight React tooling
(fast dev server, minimal config) and ships first-class Vitest
integration, so testing "for free" the same way `go test` has been
throughout this project. TypeScript catches request/response shape
mismatches against the Go API contract at compile time — valuable
specifically because that contract (Section: contracts/api-contract.md)
is hand-maintained on both sides, not generated. A single-page, single-
network editor's state (a handful of nodes/edges/demands and one or two
past run results) doesn't need a dedicated state-management library; it's
the same "don't add a dependency the problem doesn't need" discipline
every prior feature's research.md has applied to the Go side, now applied
to the frontend.

**Alternatives considered**: Plain JavaScript (no TypeScript) — rejected;
the type safety is cheap to add up front and catches exactly the kind of
mismatch a hand-written API contract is prone to. Redux/Zustand/similar —
rejected as unjustified for this scale, matching decision reasoning used
throughout this project for the Go side (e.g. feature 002's research.md
rejecting `testify` as an unneeded dependency).

## 6. Editing during playback stops the animation (not: queues, not: live-patches it)

**Decision**: Starting an edit (add/remove node or road) while a previous
run's animation is playing simply stops that animation immediately. The
user must press "run" again to see the edited network's behavior.

**Rationale**: This is the simplest of the three options the spec's Edge
Case implicitly weighs (stop, pause-and-resume, or live-patch), and the
other two solve problems this phase doesn't have: "pause and resume"
implies resuming mid-animation still matters once the underlying network
has changed (it doesn't — the old result is now stale), and "live-patch
the in-progress animation" would require reconciling an old result against
a new network structurally, which is real complexity in service of a
scenario (uninterrupted viewing while editing) no user story asks for.

**Alternatives considered**: Disabling editing entirely while an
animation plays — rejected as more restrictive than necessary; simply
stopping the animation on edit is less surprising and doesn't block the
user from immediately starting to change the network.
