# Feature Specification: Graph Simulation Engine

**Feature Branch**: `001-graph-engine`

**Created**: 2026-09-03

**Status**: Draft

**Input**: User description: "Phase 1 of the roadmap only (per constitution Principle II): a pure graph simulation engine in Go, no UI, no agents/cars yet. Nodes represent intersections, houses (origins), and companies (destinations). Edges represent roads with a length/capacity and a travel-time function. The engine must be runnable and verifiable standalone (e.g. via a terminal script/CLI), independent of any UI. This is the foundation the later agent/congestion/frontend phases will build on."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Build a road network graph (Priority: P1)

A developer defines a road network programmatically: intersections, houses
(origins), and companies (destinations) as nodes, and roads between them as
edges, each with a length, a capacity, and a travel-time function.

**Why this priority**: Every later phase (agents, congestion, frontend)
needs a graph to operate on. Without this, nothing else in the roadmap can
start.

**Independent Test**: Can be fully tested by constructing a small graph in
code and asserting its node/edge counts and attributes — no agents,
congestion logic, or UI required.

**Acceptance Scenarios**:

1. **Given** an empty graph, **When** a node is added with a declared type
   (intersection, house, or company), **Then** it appears in the graph's
   node list with that type.
2. **Given** a graph with two existing nodes, **When** an edge is added
   between them with a length and a capacity, **Then** the edge is
   retrievable and reports those attributes.

---

### User Story 2 - Compute travel time under traffic (Priority: P2)

A developer evaluates how long it takes to cross a given road for a given
amount of traffic currently on it — at zero traffic the time is the road's
free-flow time, and it does not get faster as more traffic is added.

**Why this priority**: This is the mechanism through which congestion, and
ultimately Braess's Paradox, will emerge in later phases. It must be
correct and independently verifiable before any agent relies on it.

**Independent Test**: Can be fully tested by supplying known volume values
to an edge's travel-time function and checking the returned times, with no
agent or UI involved.

**Acceptance Scenarios**:

1. **Given** an edge with zero traffic on it, **When** its travel time is
   computed, **Then** it equals the edge's configured free-flow time.
2. **Given** an edge, **When** its travel time is computed at increasing
   traffic volumes, **Then** the returned time never decreases as volume
   increases.

---

### User Story 3 - Inspect the graph from a terminal (Priority: P3)

A developer runs a standalone terminal command that builds (or loads) a
sample road network and prints its full structure — nodes, edges, and their
current travel times — with no frontend involved.

**Why this priority**: Proves, end to end, that the engine is genuinely
independent of any UI and can be validated in text before a single line of
frontend code exists.

**Independent Test**: Can be fully tested by running the command and
checking its console output and exit code — no browser or UI required.

**Acceptance Scenarios**:

1. **Given** a sample road network defined in code, **When** the terminal
   command runs, **Then** it prints every node's id and type, and every
   edge's endpoints, length, capacity, and current travel time, to
   standard output, exiting successfully.

---

### Edge Cases

- What happens when an edge is created referencing a node id that does not
  exist in the graph? The engine must reject the operation with a clear
  error instead of silently creating a dangling reference.
- How does the system handle two separate roads connecting the same pair of
  nodes (a parallel route)? This must be allowed and each edge must remain
  independently queryable — parallel/redundant routes are exactly what
  Braess's Paradox experiments need to add later.
- What happens when travel time is requested for a negative traffic volume?
  The engine must reject the request rather than return a misleading value.
- How does the system handle a node with no outgoing edges (e.g. a
  destination with nothing leaving it)? This must be representable without
  error.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The engine MUST allow creating nodes, each with a type of
  intersection, house (origin), or company (destination).
- **FR-002**: The engine MUST allow creating directed edges between two
  existing nodes, each edge carrying a length, a capacity, and a
  travel-time function of current traffic volume.
- **FR-003**: The engine MUST reject creating an edge that references a
  node id not present in the graph.
- **FR-004**: The engine MUST compute an edge's travel time for a given
  traffic volume, returning the edge's free-flow time at zero volume and a
  non-decreasing time as volume increases.
- **FR-005**: The engine MUST allow multiple distinct edges between the
  same pair of nodes, each independently queryable, so redundant/parallel
  routes can be represented.
- **FR-006**: The engine MUST be runnable and fully exercisable from a
  standalone terminal entry point, without importing or depending on any
  UI or rendering package.
- **FR-007**: The engine MUST expose a way to enumerate every node and
  every edge of a graph, for inspection, logging, and testing.
- **FR-008**: The engine MUST reject travel-time computation requests made
  with a negative traffic volume.

### Key Entities

- **Node**: An intersection, house (origin), or company (destination) in
  the road network. Has an identifier and a type.
- **Edge (Road)**: A directed connection between two nodes. Has a length,
  a capacity, and a travel-time function of current traffic volume.
- **Graph**: The full road network — the collection of all nodes and edges,
  queryable as a whole.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can define a network of at least 4 nodes and 4
  edges and see its full structure printed via the terminal command in
  under 5 seconds.
- **SC-002**: Travel time computed at zero traffic volume matches the
  edge's configured free-flow time in 100% of cases checked.
- **SC-003**: Across at least 5 increasing traffic-volume samples on the
  same edge, the computed travel time never decreases, in 100% of test
  runs.
- **SC-004**: The engine has zero UI/frontend dependencies, verifiable by
  inspecting its dependency manifest.

## Assumptions

- The exact shape of the travel-time function (e.g. linear, BPR-style) is
  not fixed by this spec — any function that is non-decreasing in traffic
  volume and returns the free-flow time at zero volume satisfies it. The
  precise formula is a planning/implementation decision and may evolve
  once agents start feeding it real traffic.
- "Traffic volume" in this phase is a value supplied directly by the
  caller (e.g. a test or the CLI), not yet produced by live agents — agent-
  driven traffic is out of scope until roadmap phase 2.
- Roads are represented as directed edges; a two-way street is represented
  as two edges (one per direction), each independently congestible —
  matching how traffic engineering typically models direction-specific
  congestion.
- This phase is in-memory only; no persistence/storage requirement is
  implied.
