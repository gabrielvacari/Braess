package simulation

import (
	"braess/agent"
	"braess/graph"
)

// epsilon absorbs floating-point noise when comparing route costs: an
// agent switches only if an alternative is cheaper by more than this,
// never on an exact or near-exact tie (research.md decision #2).
const epsilon = 1e-9

// AssignmentResult is the outcome of running a Population to a stable (or
// round-limited) route assignment.
type AssignmentResult struct {
	Routes            []agent.Route
	TotalTravelTime   float64
	AverageTravelTime float64
	// Converged is true only if a refinement round found zero agents
	// wanting to switch before MaxRounds was reached. It is false if
	// MaxRounds == 0, or the round limit was hit with agents still
	// switching.
	Converged bool
	// Rounds is how many refinement rounds actually ran (0 to MaxRounds).
	// It does not count the initial incremental-loading pass.
	Rounds int
}

// addToVolumes records one more agent using every edge in route — the
// real-traffic bookkeeping that makes travel time depend on who is
// actually on a road (FR-001), rather than a caller-chosen baseline.
func addToVolumes(volumes map[string]float64, route agent.Route) {
	for _, e := range route.Edges {
		volumes[e.ID]++
	}
}

// removeFromVolumes undoes addToVolumes for route — used when an agent is
// about to be reconsidered, so its own presence doesn't bias the
// evaluation of its current route against alternatives.
func removeFromVolumes(volumes map[string]float64, route agent.Route) {
	for _, e := range route.Edges {
		volumes[e.ID]--
	}
}

// marginalVolumes returns, for every edge in g, one more than however
// many *other* agents are currently on it. This is the atomic-congestion-
// game convention (Rosenthal, 1973): the load on an edge counts every
// player using it, including whichever one is currently deciding whether
// to use it — so evaluating a candidate route against "others + 1"
// (rather than "others") is what makes a route's reported cost the cost
// that agent would actually experience by choosing it.
func marginalVolumes(g *graph.Graph, others map[string]float64) map[string]float64 {
	edges := g.Edges()
	marginal := make(map[string]float64, len(edges))
	for _, e := range edges {
		marginal[e.ID] = others[e.ID] + 1
	}
	return marginal
}

// routeCost sums route's edges' travel times at volumes as given — used
// with volumes that already include the route's own agent, to get that
// agent's actual currently-experienced cost.
func routeCost(g *graph.Graph, route agent.Route, volumes map[string]float64) (float64, error) {
	total := 0.0
	for _, e := range route.Edges {
		tt, err := g.TravelTime(e.ID, volumes[e.ID])
		if err != nil {
			return 0, err
		}
		total += tt
	}
	return total, nil
}

// loadIncrementally assigns each of p.Size agents a route in turn, each
// against the real traffic accumulated by the agents already placed
// (via marginalVolumes, so the agent being placed counts as its route's
// +1th user). This is a realistic warm start: agents react to congestion
// from the very first one, rather than all independently picking the
// same free-flow shortest route before any volume exists at all
// (research.md decision #3).
func loadIncrementally(g *graph.Graph, p Population) ([]agent.Route, map[string]float64, error) {
	routes := make([]agent.Route, p.Size)
	volumes := make(map[string]float64)

	for i := 0; i < p.Size; i++ {
		route, err := agent.ShortestRouteAtVolumes(g, p.Origin, p.Destination, marginalVolumes(g, volumes))
		if err != nil {
			return nil, nil, err
		}
		routes[i] = route
		addToVolumes(volumes, route)
	}

	return routes, volumes, nil
}

// refine runs up to p.MaxRounds best-response rounds over routes/volumes,
// mutating both in place: each round, every agent's current route is
// compared against the best alternative given everyone else's current
// choice, and the agent switches only on strict improvement (research.md
// decision #2 — this is what makes tied routes stable instead of
// flip-flopping forever). A round with zero switches means the population
// has reached a stable assignment (guaranteed to happen in finite rounds
// for this class of game — research.md decision #1) and refinement stops
// there; otherwise it stops at p.MaxRounds, not yet converged.
func refine(g *graph.Graph, p Population, routes []agent.Route, volumes map[string]float64) (rounds int, converged bool, err error) {
	for round := 1; round <= p.MaxRounds; round++ {
		changed := false

		for i, route := range routes {
			currentCost, cerr := routeCost(g, route, volumes)
			if cerr != nil {
				return round, false, cerr
			}

			removeFromVolumes(volumes, route)
			alt, aerr := agent.ShortestRouteAtVolumes(g, p.Origin, p.Destination, marginalVolumes(g, volumes))
			if aerr != nil {
				addToVolumes(volumes, route)
				return round, false, aerr
			}

			if alt.TotalTravelTime < currentCost-epsilon {
				routes[i] = alt
				addToVolumes(volumes, alt)
				changed = true
			} else {
				addToVolumes(volumes, route)
			}
		}

		if !changed {
			return round, true, nil
		}
	}

	return p.MaxRounds, false, nil
}

// Run computes a route assignment for p on g: agents are loaded one at a
// time against real accumulating traffic (loadIncrementally), then
// refined round by round (refine) until no agent improves or
// p.MaxRounds is reached.
//
// Returns an error, rather than swallowing it, if any agent's route
// computation fails — e.g. no route exists between p.Origin and
// p.Destination (wraps agent.ErrNoRoute), or an edge's travel-time
// function itself errors.
func Run(g *graph.Graph, p Population) (AssignmentResult, error) {
	if p.Size == 0 {
		return AssignmentResult{Converged: true}, nil
	}

	routes, volumes, err := loadIncrementally(g, p)
	if err != nil {
		return AssignmentResult{}, err
	}

	rounds, converged, err := refine(g, p, routes, volumes)
	if err != nil {
		return AssignmentResult{}, err
	}

	// Refresh each route's own TotalTravelTime against the final volumes:
	// a route's cost when it was assigned or last switched can be stale
	// by the time later agents finish moving, and FR-007 wants the cost
	// "experienced ... at the final assignment," not at assignment time.
	total := 0.0
	for i, r := range routes {
		cost, err := routeCost(g, r, volumes)
		if err != nil {
			return AssignmentResult{}, err
		}
		routes[i].TotalTravelTime = cost
		total += cost
	}

	return AssignmentResult{
		Routes:            routes,
		TotalTravelTime:   total,
		AverageTravelTime: total / float64(p.Size),
		Converged:         converged,
		Rounds:            rounds,
	}, nil
}
