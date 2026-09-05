package queuesim

import (
	"container/heap"
	"fmt"

	"braess/graph"
)

// shortestRouteAt computes the minimum-cost route from `from` to `to` in
// g, where each edge's cost comes from edgeCost. This is queuesim's own
// routing — a variant of the same lazy-deletion-heap Dijkstra technique
// `agent` uses (feature 002/003) — rather than a reuse of agent's core,
// since the cost priced here (free-flow time plus an estimated signal
// wait) isn't a volume fed through graph.TravelTimeFunc (research.md
// decision #3).
func shortestRouteAt(g *graph.Graph, from, to string, edgeCost func(e graph.Edge) (float64, error)) ([]graph.Edge, error) {
	if from == to {
		return nil, nil
	}

	adjacency := make(map[string][]graph.Edge)
	for _, e := range g.Edges() {
		adjacency[e.From] = append(adjacency[e.From], e)
	}

	dist := map[string]float64{from: 0}
	prevEdge := make(map[string]graph.Edge)
	visited := make(map[string]bool)

	pq := &routePQ{{node: from, dist: 0}}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*routePQItem)
		if visited[current.node] {
			continue
		}
		visited[current.node] = true
		if current.node == to {
			break
		}

		for _, e := range adjacency[current.node] {
			cost, err := edgeCost(e)
			if err != nil {
				return nil, fmt.Errorf("queuesim: computing cost for edge %q: %w", e.ID, err)
			}
			newDist := dist[current.node] + cost
			if d, ok := dist[e.To]; !ok || newDist < d {
				dist[e.To] = newDist
				prevEdge[e.To] = e
				heap.Push(pq, &routePQItem{node: e.To, dist: newDist})
			}
		}
	}

	if !visited[to] {
		return nil, fmt.Errorf("queuesim: no route exists from %q to %q", from, to)
	}

	edges := make([]graph.Edge, 0)
	for cur := to; cur != from; {
		e, ok := prevEdge[cur]
		if !ok {
			return nil, fmt.Errorf("queuesim: could not reconstruct route from %q to %q", from, to)
		}
		edges = append(edges, e)
		cur = e.From
	}
	for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
		edges[i], edges[j] = edges[j], edges[i]
	}
	return edges, nil
}

// routePQItem is one entry in the Dijkstra priority queue.
type routePQItem struct {
	node string
	dist float64
}

// routePQ is a container/heap min-heap of routePQItem ordered by dist.
type routePQ []*routePQItem

func (pq routePQ) Len() int            { return len(pq) }
func (pq routePQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq routePQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *routePQ) Push(x interface{}) { *pq = append(*pq, x.(*routePQItem)) }
func (pq *routePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
