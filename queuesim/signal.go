// Package queuesim is a discrete-time, tick-based simulation: signals
// gate whether a road is passable, agents spawn over simulated time per
// demand, and queue in arrival order at any signal-controlled road they
// reach while it is red. It depends only on the graph package and the Go
// standard library — never on agent, simulation, or any UI/rendering
// package (constitution Principle I). It models a different question
// than simulation's equilibrium-seeking (feature 003/004): not "what
// stable assignment do selfish agents reach," but "how does a queue
// build and drain over real time." Whether the classic two-road,
// opposite-phase experiment produces a persistently congested road, an
// oscillating pattern, or a balanced outcome is this package's
// experimental result, not an assumption baked into it (spec 008
// Assumptions).
package queuesim

import "math"

// Signal governs one road (a graph.Edge, by ID), alternating between a
// green phase (passable) and a red phase (not passable). Two Signals are
// "always opposite" purely by being configured with equal
// GreenDuration+RedDuration and an Offset one GreenDuration apart —
// there is no linking field; it falls out of the phase arithmetic below
// (research.md decision #4).
type Signal struct {
	ID            string
	EdgeID        string
	GreenDuration float64
	RedDuration   float64
	// Offset shifts the cycle: how many seconds into the cycle the
	// signal already is at simulated t=0.
	Offset float64
	// DischargeRate is how many queued agents can leave per simulated
	// second once green. 0 means "not explicitly configured": a default
	// derived from the edge's own free-flow travel time applies instead
	// (spec Assumptions) — see defaultDischargeRate in run.go.
	DischargeRate float64
}

// cycle is the signal's full green+red period.
func (s Signal) cycle() float64 {
	return s.GreenDuration + s.RedDuration
}

// IsGreenAt reports whether s is in its green phase at simulated time t.
func (s Signal) IsGreenAt(t float64) bool {
	c := s.cycle()
	if c <= 0 {
		return false
	}
	phase := math.Mod(t+s.Offset, c)
	if phase < 0 {
		phase += c
	}
	return phase < s.GreenDuration
}

// timeUntilNextGreen returns how many seconds from t until s next becomes
// green (0 if already green at t).
func (s Signal) timeUntilNextGreen(t float64) float64 {
	if s.IsGreenAt(t) {
		return 0
	}
	c := s.cycle()
	phase := math.Mod(t+s.Offset, c)
	if phase < 0 {
		phase += c
	}
	return c - phase
}

// estimatedWait is a snapshot estimate of how long a new arrival at t
// would wait to use s's edge, given queueLength agents already waiting
// and a dischargeRate (agents released per second once green). This is
// an approximation used only to choose a route at spawn time (research.md
// decision #2), not a guarantee: it doesn't account for how queueLength
// or the phase will have changed by the time an agent with earlier edges
// to traverse first actually arrives — exact when the edge being priced
// is the first (or only) edge of the route.
func (s Signal) estimatedWait(t float64, queueLength int, dischargeRate float64) float64 {
	wait := s.timeUntilNextGreen(t)
	if dischargeRate > 0 {
		wait += float64(queueLength) / dischargeRate
	}
	return wait
}
