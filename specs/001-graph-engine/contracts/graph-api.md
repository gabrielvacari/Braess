# Contract: `graph` package public API

This is a library, not a network service — its "contract" is the exported
Go API that later phases (agents, frontend) and this phase's own CLI and
tests depend on. Signatures below are the contract; internal representation
is free to change as long as these hold.

## Types

```go
package graph

type NodeType int

const (
    Intersection NodeType = iota
    House
    Company
)

type Node struct {
    ID   string
    Type NodeType
}

// TravelTimeFunc computes the travel time for an edge at a given traffic
// volume. Implementations MUST return the edge's free-flow time at
// volume == 0, MUST be non-decreasing as volume increases, and MUST
// return an error for volume < 0.
type TravelTimeFunc func(volume float64) (float64, error)

type Edge struct {
    ID         string
    From, To   string // Node IDs
    Length     float64
    Capacity   float64
    TravelTime TravelTimeFunc
}

type Graph struct {
    // unexported fields
}
```

## Constructors

```go
func New() *Graph

// Linear returns a TravelTimeFunc where time = freeFlow + slope*volume.
func Linear(freeFlow, slope float64) TravelTimeFunc

// Constant returns a TravelTimeFunc that always returns t, regardless of volume.
func Constant(t float64) TravelTimeFunc
```

## Graph methods

```go
// AddNode adds a node to the graph.
// Errors if id is empty or already present (FR-001).
func (g *Graph) AddNode(id string, t NodeType) (Node, error)

// AddEdge adds a directed edge between two existing nodes.
// Errors if id is empty or already present, or if from/to do not
// reference existing nodes (FR-002, FR-003).
// Multiple edges between the same (from, to) pair are allowed (FR-005).
func (g *Graph) AddEdge(id, from, to string, length, capacity float64, tt TravelTimeFunc) (Edge, error)

// Nodes returns every node in the graph, in unspecified order (FR-007).
func (g *Graph) Nodes() []Node

// Edges returns every edge in the graph, in unspecified order (FR-007).
func (g *Graph) Edges() []Edge

// TravelTime evaluates the travel time of the edge identified by edgeID
// at the given traffic volume. Errors if edgeID does not exist or if
// volume < 0 (FR-004, FR-008).
func (g *Graph) TravelTime(edgeID string, volume float64) (float64, error)
```

## Error handling

All fallible operations return a Go `error` as the last return value
(idiomatic Go) rather than panicking — callers (the CLI, tests, and future
agent code) are expected to check it. No sentinel error values are
mandated by this contract; wrapped, descriptive errors (`fmt.Errorf`) are
sufficient to satisfy FR-003 and FR-008's "reject with a clear error"
requirement.

## Stability

This contract covers roadmap Phase 1 only. It is expected to grow (not
break) in later phases — e.g. agents will need a way to read current
traffic volume per edge, which isn't required yet since volume is
caller-supplied in this phase per the spec's Assumptions.
