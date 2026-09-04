package graph

import "fmt"

// TravelTimeFunc computes the travel time for an edge at a given traffic
// volume. Implementations MUST return the edge's free-flow time at
// volume == 0, MUST be non-decreasing as volume increases, and MUST return
// an error for volume < 0.
//
// A single global formula is deliberately not imposed here: the classic
// Braess's Paradox network needs edges with different congestion behavior
// side by side (see research.md, decision #2) — some roads slow down as
// more cars use them, others don't. Linear and Constant below are the
// initial, swappable building blocks for that; nothing prevents adding
// more later (e.g. a BPR-style curve) without changing this type.
type TravelTimeFunc func(volume float64) (float64, error)

// Linear returns a TravelTimeFunc where time = freeFlow + slope*volume.
// With slope == 0 this behaves like Constant(freeFlow).
func Linear(freeFlow, slope float64) TravelTimeFunc {
	return func(volume float64) (float64, error) {
		if volume < 0 {
			return 0, fmt.Errorf("graph: traffic volume must be >= 0, got %g", volume)
		}
		return freeFlow + slope*volume, nil
	}
}

// Constant returns a TravelTimeFunc that always returns t, regardless of
// traffic volume — for roads whose travel time does not depend on how many
// cars are using them.
func Constant(t float64) TravelTimeFunc {
	return func(volume float64) (float64, error) {
		if volume < 0 {
			return 0, fmt.Errorf("graph: traffic volume must be >= 0, got %g", volume)
		}
		return t, nil
	}
}
