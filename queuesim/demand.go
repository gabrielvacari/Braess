package queuesim

// Demand describes agents spawning over simulated time for one
// origin/destination pair — unlike simulation.Demand (feature 003/004),
// which has no time dimension (every agent is already present at the
// start of that equilibrium computation). A steady arrival process is
// what makes a queue mean anything: without one, every agent would need
// to cross a signal-controlled road in the very first tick, an
// unrealistic burst rather than a stream (research.md decision #5).
type Demand struct {
	Origin, Destination string
	// Count is the total number of agents this demand spawns over a run.
	Count int
	// ArrivalInterval is the simulated seconds between one spawn and the
	// next for this demand.
	ArrivalInterval float64
}
