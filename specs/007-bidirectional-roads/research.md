# Phase 0 Research: Bidirectional Roads

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per
constitution Principle V.

## 1. The client's editable primitive becomes the road, not the edge

**Decision**: Replace `ClientEdge`/`NetworkState.edges` (one directed
connection) with `ClientRoad`/`NetworkState.roads` (one bidirectional
connection, `{ id, a, b, travelTime }` — no "from"/"to" direction).
Directed edges are no longer something the client model stores at all;
they're derived, at the API boundary, from roads.

**Rationale**: FR-004/FR-005 need a drawn road to be exactly one thing —
one line on screen, one thing to remove — and the engine's directed-edge
model (feature 001) can't represent that without an extra grouping
mechanism layered on top (e.g. a shared `roadId` field on two separate
edge records). Making the road itself the primitive means those
requirements hold *by construction*: there is only one array entry per
drawn road, so there is nothing to keep in sync and nothing to
accidentally remove only half of.

**Alternatives considered**: Keep `ClientEdge` as-is and add a `roadId`
field so two edges can be recognized as a pair — rejected; it would leave
two overlapping concepts (edge and road) permanently coexisting in the
client model for no benefit, since nothing in this frontend ever needs to
treat one direction of a road differently from the other (spec
Assumptions: both directions always share one configuration).

## 2. One shared function expands a road into its two directed edges

**Decision**: `expandRoadsToDirectedEdges(roads: ClientRoad[])` is the
single function that turns roads into the directed-edge shape the server
understands. Both `api.ts`'s `toRunRequest` (building the request) and
`App.tsx`'s route-animation path resolver (interpreting the response)
call it — neither reimplements the forward/reverse expansion itself.

**Rationale**: This is the same "one shared core, no second
implementation to drift out of sync" pattern the Go side has used
repeatedly (e.g. `agent.ShortestRoute`/`ShortestRouteAtVolumes` sharing
one Dijkstra core, feature 003). Here, the request builder and the
response interpreter both need to agree on exactly which synthetic edge
ids exist and which nodes they connect — if they computed that
independently, a future change to one could silently break the other.

**Alternatives considered**: Expanding inline wherever needed (once in
`api.ts`, once in `App.tsx`) — rejected for the drift risk above; two
near-identical five-line functions are a maintenance trap, not a
simplification.

## 3. Synthetic directed-edge ids are deterministic, not random

**Decision**: A road with id `road-3` expands to directed edges
`road-3-ab` (from `a` to `b`) and `road-3-ba` (from `b` to `a`) — derived
from the road's own id, not freshly generated each call.

**Rationale**: `expandRoadsToDirectedEdges` may be called more than once
for the same `NetworkState` (once to build a request, potentially again
later to resolve a response against the *current* roads). Deterministic
ids guarantee both calls agree on what a given directed edge id means,
which is what lets the animation code look up "which road did edge
`road-3-ab` belong to" without needing the server to echo back anything
beyond the edge ids it already returns.

**Alternatives considered**: Generating a fresh random id per direction
per call — rejected; it would make the request and the later response
interpretation potentially disagree about naming, since they're built at
different times from what should be the same conceptual mapping.

## 4. Both directions always share one travel-time configuration

**Decision**: `ClientRoad` has one `travelTime: TravelTimeSpec` field,
applied to both derived directed edges. There is no per-direction
override in this feature.

**Rationale**: Directly matches the spec's own wording ("apenas uma
conexão" — the user should only have to make one connection, implying one
configuration) and FR-003. It also keeps the existing "Road behavior"
picker in `Toolbar.tsx` valid completely unchanged — one choice, one
road, both directions.

**Alternatives considered**: Letting each direction be configured
separately — rejected as unrequested scope; nothing in the spec or the
reported problem asks for asymmetric roads, and it would double the
Toolbar's road-drawing form for no current benefit.

## 5. Zero server or engine change, confirmed by construction

**Decision**: `cmd/server`'s `POST /api/run` contract, and everything
under `graph`/`agent`/`simulation`, are untouched.

**Rationale**: The server already accepts an arbitrary `edges: EdgeDTO[]`
array (feature 005) — nothing about its contract assumes edges come in
matched forward/reverse pairs or any particular quantity per "road."
Sending two directed edges instead of one for a drawn connection is
purely a client-side decision about *what it sends*, not a new
capability the server needs. This is also exactly what FR-007 requires.

**Alternatives considered**: Adding a server-side "bidirectional" flag to
`EdgeDTO` that expands into two edges on the Go side — rejected; it would
duplicate the expansion logic in two languages for a purely
client-side-observable feature, and would violate this project's
long-standing rule (every prior feature) against touching the engine
without a reason the engine itself needs.
