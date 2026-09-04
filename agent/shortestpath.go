package agent

import (
	"container/heap"
	"errors"
	"fmt"

	"braess/graph"
)

// ErrNoRoute indicates no path connects the given origin to the given
// destination. Check with errors.Is.
var ErrNoRoute = errors.New("agent: no route exists")

// ShortestRoute computes the minimum-total-travel-time route from `from`
// to `to` in g, evaluating every edge's travel time at the given volume
// via Dijkstra's algorithm (edge weights — travel times at a fixed
// baseline volume — are assumed non-negative; see research.md decision
// #4). This phase's only caller, Agent.ComputeRoute, always passes 0.
//
// Returns a zero-edge, zero-time Route if from == to.
// Returns an error wrapping ErrNoRoute if no path exists.
// Returns, rather than swallows, any error an edge's TravelTimeFunc
// itself produces while the route is being computed — a broken edge must
// not be silently treated as absent or instantly fast.
// Parallel edges between the same node pair are each considered as a
// distinct option, since adjacency is built from every edge in the graph
// without deduplicating by (From, To).
func ShortestRoute(g *graph.Graph, from, to string, volume float64) (Route, error) {
	if from == to {
		return Route{}, nil
	}

	adjacency := make(map[string][]graph.Edge)
	for _, e := range g.Edges() {
		adjacency[e.From] = append(adjacency[e.From], e)
	}

	dist := map[string]float64{from: 0}
	prevEdge := make(map[string]graph.Edge)
	visited := make(map[string]bool)

	pq := &priorityQueue{{node: from, dist: 0}}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*pqItem)
		if visited[current.node] {
			continue
		}
		visited[current.node] = true
		if current.node == to {
			break
		}

		for _, e := range adjacency[current.node] {
			tt, err := g.TravelTime(e.ID, volume)
			if err != nil {
				return Route{}, fmt.Errorf("agent: computing travel time for edge %q: %w", e.ID, err)
			}
			newDist := dist[current.node] + tt
			if d, ok := dist[e.To]; !ok || newDist < d {
				dist[e.To] = newDist
				prevEdge[e.To] = e
				heap.Push(pq, &pqItem{node: e.To, dist: newDist})
			}
		}
	}

	if !visited[to] {
		return Route{}, fmt.Errorf("%w: from %q to %q", ErrNoRoute, from, to)
	}

	edges := make([]graph.Edge, 0)
	for cur := to; cur != from; {
		e, ok := prevEdge[cur]
		if !ok {
			// Unreachable given visited[to] is true, but fail loudly
			// rather than return a silently incomplete route.
			return Route{}, fmt.Errorf("agent: could not reconstruct route from %q to %q", from, to)
		}
		edges = append(edges, e)
		cur = e.From
	}
	for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
		edges[i], edges[j] = edges[j], edges[i]
	}

	return Route{Edges: edges, TotalTravelTime: dist[to]}, nil
}

// pqItem is one entry in the Dijkstra priority queue.
type pqItem struct {
	node string
	dist float64
}

// priorityQueue is a container/heap min-heap of pqItem ordered by dist.
// Stale entries (a node pushed more than once as shorter distances are
// found) are skipped on Pop via the `visited` check in ShortestRoute,
// the standard lazy-deletion approach to Dijkstra with container/heap.
type priorityQueue []*pqItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*pqItem)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
