# Phase 0 Research: Signal Queuing Frontend

No `[NEEDS CLARIFICATION]` markers remained in the Technical Context, so
this covers the technology/design decisions made and why — per
constitution Principle V.

## 1. Signals-mode roads are directed, not the bidirectional `ClientRoad`

**Decision**: Signals mode introduces its own `DirectedRoad` type
(`{id, from, to, travelTime, signal?}`) rather than reusing feature 007's
`ClientRoad` (bidirectional).

**Rationale**: A traffic signal inherently governs one direction of
travel — that's the entire premise of the original motivating question
("one road green, the other red"). `queuesim.Signal` itself is tied to
one `graph.Edge` (one direction) by design (feature 008). Forcing
Signals-mode roads to be bidirectional (feature 007's model) would mean
either attaching one signal to *both* directions of a drawn road
(misrepresenting a real intersection, where each direction typically has
its own light) or inventing a second signal per direction — both add
complexity in service of a bidirectionality feature 007 solved for a
different problem (avoiding "no route exists" confusion in the
equilibrium mode's free-choice routing, which doesn't apply here: a
Signals-mode demand's route is computed once from real, directed roads,
and a one-way road is exactly what the feature is meant to explore).

**Alternatives considered**: Reusing `ClientRoad` and attaching one
`Signal` to its "forward" direction only, implicitly leaving the reverse
direction always open — rejected as a confusing, undocumented asymmetry;
a fresh, explicitly directed type is clearer than reusing a bidirectional
one with an implicit exception.

## 2. Signals mode is a fully separate slice, not woven into shared state

**Decision**: Signals mode has its own `QueueNetworkState` (nodes, its
own directed roads, its own time-based demands) and its own toolbar and
results panel — not fields added to the existing `NetworkState`/
`Toolbar`/`ResultsPanel`. The two modes do not share a drawn network.

**Rationale**: FR-006/FR-007 require the two capabilities to never
interfere with or be confused for one another. Achieving that inside one
shared state object would require careful, ongoing discipline (which
fields belong to which mode, ensuring a signal-mode-only field never
leaks into an equilibrium request, and vice versa) — a fully separate
slice makes it structurally impossible for one mode's data to
accidentally reach the other's request builder, the same "make the
constraint true by construction" reasoning several earlier features have
applied (e.g. feature 007's road-as-single-record design). Nothing in the
spec's Acceptance Scenarios requires sharing one network across modes;
each user story describes configuring and running one self-contained
scenario.

**Alternatives considered**: One shared `NetworkState` with optional,
mode-specific fields (a road's optional `signal`, an optional
`arrivalInterval` on a demand) — rejected: it would make `toRunRequest`
and a new `toQueueRunRequest` both need to know which fields to ignore
for their own mode, and a bidirectional `ClientRoad` still wouldn't fit
Signals mode's directed model (decision #1) without an awkward
reinterpretation.

## 3. `NetworkCanvas` is reused for rendering; nothing else is

**Decision**: Signals mode calls the existing `NetworkCanvas` component
for drawing nodes and road lines (mapping a `DirectedRoad`'s `from`/`to`
to the `{a, b}` shape it already expects), but has its own toolbar and
results components.

**Rationale**: `NetworkCanvas` only needs node positions and a pair of
endpoint ids per line to render — it has no opinion about directionality
or bidirectionality, so reusing it costs nothing and avoids a second
canvas-rendering implementation to maintain. Everything Signals mode
actually does differently (attaching a signal while drawing, a
time-based demand form, a queue-length-over-time result) doesn't live in
`NetworkCanvas` at all, so there's no risk of the two modes' distinct
concerns bleeding into shared rendering code.

**Alternatives considered**: A second, Signals-specific canvas component
— rejected as needless duplication of working, mode-agnostic rendering
code.

## 4. The queue-length-over-time chart is a small hand-drawn SVG, not a library

**Decision**: `QueueChart.tsx` computes simple point coordinates from a
`QueueSample[]` array (a pure, tested function) and renders them as an
inline SVG `<polyline>` — no charting library dependency.

**Rationale**: Continues every prior feature's "don't add a dependency
the problem doesn't need" discipline (feature 002 rejected `testify`,
feature 006 rejected Tailwind, feature 007 needed nothing new) — a line
chart with a handful of series and a few dozen points each is well within
what ~30 lines of SVG-generation code can do cleanly, and keeping the
geometry computation as a pure function makes it directly Vitest-testable
without needing to render anything (Principle IV).

**Alternatives considered**: A charting library (e.g. Recharts, Chart.js)
— rejected as unjustified weight for this scale, and it would be the
project's first UI dependency beyond `react-konva` itself.

## 5. The server reuses `dto.go`'s `NodeDTO`/`EdgeDTO`/`buildGraph` as-is

**Decision**: `QueueRunRequest` embeds the *same* `NodeDTO`/`EdgeDTO`
types feature 005 already defined, and the new handler calls the
existing (unexported) `buildGraph` helper to construct the `*graph.Graph`
— only `Signals []SignalDTO` and `Demands []TimeDemandDTO` are new types.

**Rationale**: A node and a directed edge mean exactly the same thing to
both endpoints — there's no reason for a second, parallel definition of
"what a node/edge looks like on the wire," and `buildGraph`'s existing
self-loop rejection and node/edge validation apply identically here.
This is a small, safe reuse within the same Go package (not a change to
`dto.go`'s existing exported behavior — feature 005/006's contract tests
keep passing unmodified).

**Alternatives considered**: A fully separate `QueueRunRequest` node/edge
shape — rejected as pure duplication of an already-correct, already-
tested mapping.
