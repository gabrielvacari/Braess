// Command graphcli builds a small sample road network with the graph
// engine and prints its full structure to standard output. It exists to
// prove, per constitution Principle I and IV, that the engine works and is
// verifiable in text — completely standalone, with zero UI dependency.
package main

import (
	"fmt"
	"os"

	"braess/graph"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "graphcli:", err)
		os.Exit(1)
	}
}

func run() error {
	g, err := sampleNetwork()
	if err != nil {
		return err
	}

	printNodes(g)
	printEdges(g)
	return nil
}

// sampleNetwork builds a small network mixing every NodeType and both
// TravelTimeFunc shapes: houses and companies connected through an
// intersection, with one congestible road (Linear) and one road whose
// travel time doesn't depend on volume (Constant) — the two ingredients
// the classic Braess's Paradox network needs, once agents exist to use
// them.
func sampleNetwork() (*graph.Graph, error) {
	g := graph.New()

	nodes := []struct {
		id string
		t  graph.NodeType
	}{
		{"house-1", graph.House},
		{"intersection-1", graph.Intersection},
		{"company-1", graph.Company},
		{"company-2", graph.Company},
	}
	for _, n := range nodes {
		if _, err := g.AddNode(n.id, n.t); err != nil {
			return nil, fmt.Errorf("building sample network: %w", err)
		}
	}

	edges := []struct {
		id, from, to     string
		length, capacity float64
		travelTime       graph.TravelTimeFunc
	}{
		{"road-1", "house-1", "intersection-1", 5, 50, graph.Linear(2, 0.1)},
		{"road-2", "intersection-1", "company-1", 8, 100, graph.Constant(6)},
		{"road-3", "intersection-1", "company-2", 6, 80, graph.Linear(3, 0.2)},
		{"road-4", "house-1", "company-2", 12, 40, graph.Constant(9)},
	}
	for _, e := range edges {
		if _, err := g.AddEdge(e.id, e.from, e.to, e.length, e.capacity, e.travelTime); err != nil {
			return nil, fmt.Errorf("building sample network: %w", err)
		}
	}

	return g, nil
}

func printNodes(g *graph.Graph) {
	fmt.Println("Nodes:")
	for _, n := range g.Nodes() {
		fmt.Printf("  %-16s type=%s\n", n.ID, n.Type)
	}
}

func printEdges(g *graph.Graph) {
	const sampleVolume = 10.0

	fmt.Println("Edges:")
	for _, e := range g.Edges() {
		tt, err := g.TravelTime(e.ID, sampleVolume)
		if err != nil {
			fmt.Printf("  %-8s %s -> %-16s length=%-4g capacity=%-4g travel_time=<error: %v>\n",
				e.ID, e.From, e.To, e.Length, e.Capacity, err)
			continue
		}
		fmt.Printf("  %-8s %s -> %-16s length=%-4g capacity=%-4g travel_time(volume=%g)=%g\n",
			e.ID, e.From, e.To, e.Length, e.Capacity, sampleVolume, tt)
	}
}
