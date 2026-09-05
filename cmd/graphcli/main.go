// Command graphcli builds a small sample road network with the graph
// engine, routes a few agents across it, and prints the full result to
// standard output. It exists to prove, per constitution Principle I and
// IV, that the engine and its agents work and are verifiable in text —
// completely standalone, with zero UI dependency.
package main

import (
	"fmt"
	"os"

	"braess/agent"
	"braess/graph"
	"braess/queuesim"
	"braess/simulation"
)

// fixedOrigin and fixedDestination are this roadmap phase's one
// externally-fixed origin/destination pair (spec 002-shortest-path-agents
// Assumptions) — multiple distinct houses/companies arrive in Phase 4.
const (
	fixedOrigin      = "house-1"
	fixedDestination = "company-2"
	simulatedAgents  = 3
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
	if err := printAgentRoutes(g); err != nil {
		return err
	}

	if err := printBraessParadox(); err != nil {
		return err
	}

	if err := printMultiDemand(); err != nil {
		return err
	}

	return printSignalQueuing()
}

// printAgentRoutes constructs several independent agents for the fixed
// origin/destination pair and prints each one's route. Each Agent value
// is separate and calls ComputeRoute itself — nothing here computes one
// route and copies it to the others (constitution Principle III, FR-003).
func printAgentRoutes(g *graph.Graph) error {
	fmt.Printf("Agents (%s -> %s):\n", fixedOrigin, fixedDestination)
	for i := 1; i <= simulatedAgents; i++ {
		a := &agent.Agent{Origin: fixedOrigin, Destination: fixedDestination}
		route, err := a.ComputeRoute(g)
		if err != nil {
			return fmt.Errorf("agent %d: computing route: %w", i, err)
		}

		fmt.Printf("  agent-%d: %s", i, fixedOrigin)
		for _, e := range route.Edges {
			fmt.Printf(" -[%s]-> %s", e.ID, e.To)
		}
		fmt.Printf(" (total_travel_time=%g)\n", route.TotalTravelTime)
	}
	return nil
}

// braessAgentCount matches the textbook Braess's Paradox example (4000
// drivers), so the printed result below reproduces its well-known numbers
// exactly (spec 003-congestion-braess-paradox, research.md decision #6).
const braessAgentCount = 4000

// printBraessParadox runs the classic four-node Braess's Paradox network
// twice — once without its extra connecting road, once with it — and
// prints both results plus an explicit comparison. This is the roadmap's
// pre-frontend validation gate (AGENTS.md, constitution Principle II):
// proving in text, before any UI exists, that adding a road can make a
// selfishly-routed network worse for everyone.
func printBraessParadox() error {
	population := simulation.Population{
		Origin:      "S",
		Destination: "T",
		Size:        braessAgentCount,
		MaxRounds:   simulation.DefaultMaxRounds,
	}

	fmt.Println("Braess's Paradox (classic four-node network, 4000 agents):")

	without, err := braessNetwork(false)
	if err != nil {
		return fmt.Errorf("building Braess network without shortcut: %w", err)
	}
	resultWithout, err := simulation.Run(without, population)
	if err != nil {
		return fmt.Errorf("running population without shortcut: %w", err)
	}
	fmt.Printf("  without the extra road: converged=%t rounds=%d total=%g average=%g\n",
		resultWithout.Converged, resultWithout.Rounds, resultWithout.TotalTravelTime, resultWithout.AverageTravelTime)

	with, err := braessNetwork(true)
	if err != nil {
		return fmt.Errorf("building Braess network with shortcut: %w", err)
	}
	resultWith, err := simulation.Run(with, population)
	if err != nil {
		return fmt.Errorf("running population with shortcut: %w", err)
	}
	fmt.Printf("  with the extra road:    converged=%t rounds=%d total=%g average=%g\n",
		resultWith.Converged, resultWith.Rounds, resultWith.TotalTravelTime, resultWith.AverageTravelTime)

	if resultWith.AverageTravelTime > resultWithout.AverageTravelTime {
		fmt.Printf("  => adding the road made the average trip WORSE (%g -> %g) — Braess's Paradox\n",
			resultWithout.AverageTravelTime, resultWith.AverageTravelTime)
	} else {
		fmt.Printf("  => adding the road did not make things worse here (%g -> %g)\n",
			resultWithout.AverageTravelTime, resultWith.AverageTravelTime)
	}
	return nil
}

// braessNetwork builds the classic four-node Braess's Paradox network:
// S -> A and B -> T are congestible (time = volume/100); S -> B and
// A -> T are constant at 45. withShortcut adds the free A -> B edge that
// is the source of the paradox.
func braessNetwork(withShortcut bool) (*graph.Graph, error) {
	g := graph.New()
	for _, id := range []string{"S", "A", "B", "T"} {
		if _, err := g.AddNode(id, graph.Intersection); err != nil {
			return nil, err
		}
	}

	edges := []struct {
		id, from, to string
		tt           graph.TravelTimeFunc
	}{
		{"S-A", "S", "A", graph.Linear(0, 0.01)},
		{"S-B", "S", "B", graph.Constant(45)},
		{"A-T", "A", "T", graph.Constant(45)},
		{"B-T", "B", "T", graph.Linear(0, 0.01)},
	}
	for _, e := range edges {
		if _, err := g.AddEdge(e.id, e.from, e.to, 0, 0, e.tt); err != nil {
			return nil, err
		}
	}

	if withShortcut {
		if _, err := g.AddEdge("A-B", "A", "B", 0, 0, graph.Constant(0)); err != nil {
			return nil, err
		}
	}

	return g, nil
}

// printMultiDemand runs several distinct house/company demands together
// on one shared network and prints each demand's own result plus one
// overall result — proving, in text, that congestion is genuinely shared
// across distinct origin/destination pairs (roadmap Phase 4).
func printMultiDemand() error {
	g, err := multiDemandNetwork()
	if err != nil {
		return fmt.Errorf("building multi-demand network: %w", err)
	}

	mp := simulation.MultiPopulation{
		Demands: []simulation.Demand{
			{Origin: "house-A", Destination: "company-X", Size: 40},
			{Origin: "house-B", Destination: "company-Y", Size: 35},
			{Origin: "house-C", Destination: "company-Z", Size: 30},
		},
		MaxRounds: simulation.DefaultMaxRounds,
	}

	result, err := simulation.RunDemands(g, mp)
	if err != nil {
		return fmt.Errorf("running multi-demand population: %w", err)
	}

	fmt.Println("Multiple origins/destinations (shared bottleneck road):")
	for _, dr := range result.PerDemand {
		fmt.Printf("  %s -> %s (%d agents): total=%g average=%g\n",
			dr.Demand.Origin, dr.Demand.Destination, dr.Demand.Size, dr.TotalTravelTime, dr.AverageTravelTime)
	}
	fmt.Printf("  overall: converged=%t rounds=%d total=%g average=%g\n",
		result.Converged, result.Rounds, result.TotalTravelTime, result.AverageTravelTime)

	return nil
}

// multiDemandNetwork builds a network where three houses each feed into a
// shared hub, a single bottleneck road (hub-in -> hub-out) is the only way
// across, and the hub fans out to three companies. All three demands are
// forced through the same bottleneck, so it congests all of them together.
func multiDemandNetwork() (*graph.Graph, error) {
	g := graph.New()

	nodes := []struct {
		id string
		t  graph.NodeType
	}{
		{"house-A", graph.House}, {"house-B", graph.House}, {"house-C", graph.House},
		{"hub-in", graph.Intersection}, {"hub-out", graph.Intersection},
		{"company-X", graph.Company}, {"company-Y", graph.Company}, {"company-Z", graph.Company},
	}
	for _, n := range nodes {
		if _, err := g.AddNode(n.id, n.t); err != nil {
			return nil, err
		}
	}

	edges := []struct {
		id, from, to string
		tt           graph.TravelTimeFunc
	}{
		{"house-A-hub", "house-A", "hub-in", graph.Constant(0)},
		{"house-B-hub", "house-B", "hub-in", graph.Constant(0)},
		{"house-C-hub", "house-C", "hub-in", graph.Constant(0)},
		{"bottleneck", "hub-in", "hub-out", graph.Linear(0, 0.1)},
		{"hub-company-X", "hub-out", "company-X", graph.Constant(0)},
		{"hub-company-Y", "hub-out", "company-Y", graph.Constant(0)},
		{"hub-company-Z", "hub-out", "company-Z", graph.Constant(0)},
	}
	for _, e := range edges {
		if _, err := g.AddEdge(e.id, e.from, e.to, 0, 0, e.tt); err != nil {
			return nil, err
		}
	}

	return g, nil
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

// printSignalQueuing runs the classic two-road, opposite-phase signal
// scenario (feature 008) and prints each road's queue length over time
// plus an arrival/wait summary — proving, in text, that real queuing
// happens and letting the reader see for themselves whether one road
// ends up more congested, it oscillates, or it balances (spec
// Assumptions: this feature does not presume the answer).
func printSignalQueuing() error {
	g, signals, demands, err := signalQueuingNetwork()
	if err != nil {
		return fmt.Errorf("building signal queuing network: %w", err)
	}
	const duration, tick = 70, 0.25

	result, err := queuesim.Run(g, signals, demands, duration, tick)
	if err != nil {
		return fmt.Errorf("running signal queuing scenario: %w", err)
	}

	fmt.Println("Signal-controlled queuing (two roads, always-opposite signals):")

	avgQueue := make(map[string]float64)
	count := make(map[string]int)
	for _, s := range result.QueueSamples {
		avgQueue[s.SignalID] += float64(s.Length)
		count[s.SignalID]++
	}
	for _, s := range signals {
		if n := count[s.ID]; n > 0 {
			avgQueue[s.ID] /= float64(n)
		}
		fmt.Printf("  %s (edge %s): green=%gs red=%gs -> average queue length=%.2f\n",
			s.ID, s.EdgeID, s.GreenDuration, s.RedDuration, avgQueue[s.ID])
	}

	// A snapshot of each road's queue every 5 simulated seconds, enough
	// to see the shape of the curve rise and fall.
	fmt.Println("  queue length over time (every 5s):")
	for _, s := range signals {
		fmt.Printf("    %-10s", s.ID+":")
		for _, sample := range result.QueueSamples {
			if sample.SignalID != s.ID {
				continue
			}
			isSampleTick := sample.Time == float64(int(sample.Time/5))*5
			if isSampleTick {
				fmt.Printf(" %d", sample.Length)
			}
		}
		fmt.Println()
	}

	arrived, waiting := 0, 0
	var totalWait float64
	for _, a := range result.Agents {
		if a.Arrived {
			arrived++
		} else {
			waiting++
		}
		totalWait += a.WaitTime
	}
	fmt.Printf("  agents: spawned=%d arrived=%d still-waiting=%d average_wait=%.2f\n",
		len(result.Agents), arrived, waiting, totalWait/float64(len(result.Agents)))

	return nil
}

// signalQueuingNetwork builds the classic two-road, opposite-phase
// setup: one house, one company, two parallel roads between them, each
// governed by its own signal, always in opposite phase purely by sharing
// a cycle and being offset by one GreenDuration (queuesim research.md
// decision #4) — mirrors queuesim/two_road_test.go's fixture (a test
// file, so not importable from here — feature 003's cmd/graphcli
// established the same small-duplication precedent for its Braess
// network).
func signalQueuingNetwork() (*graph.Graph, []queuesim.Signal, []queuesim.Demand, error) {
	g := graph.New()
	for _, id := range []string{"h1", "c1"} {
		t := graph.House
		if id == "c1" {
			t = graph.Company
		}
		if _, err := g.AddNode(id, t); err != nil {
			return nil, nil, nil, err
		}
	}
	for _, edgeID := range []string{"road-a", "road-b"} {
		if _, err := g.AddEdge(edgeID, "h1", "c1", 0, 0, graph.Constant(1)); err != nil {
			return nil, nil, nil, err
		}
	}

	signals := []queuesim.Signal{
		{ID: "signal-a", EdgeID: "road-a", GreenDuration: 5, RedDuration: 5, Offset: 5, DischargeRate: 2},
		{ID: "signal-b", EdgeID: "road-b", GreenDuration: 5, RedDuration: 5, Offset: 0, DischargeRate: 2},
	}
	demands := []queuesim.Demand{
		{Origin: "h1", Destination: "c1", Count: 150, ArrivalInterval: 0.3},
	}

	return g, signals, demands, nil
}
