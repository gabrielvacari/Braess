// Package simulation runs a Population of agents sharing one fixed
// origin/destination to a stable route assignment: agents are loaded
// incrementally against real traffic, then repeatedly allowed to switch
// to a strictly better route given everyone else's current choice, until
// nobody can improve or a round limit is hit. It depends only on the
// agent and graph packages and the Go standard library — never on any
// UI/rendering package (constitution Principle I).
//
// Convergence isn't merely observed to happen — it's guaranteed. This
// models the population as an atomic congestion game (Rosenthal, 1973),
// which is a potential game: a single scalar strictly decreases on every
// strictly-improving switch, and since the number of possible assignments
// is finite, best-response dynamics reaches a pure Nash equilibrium in
// finite steps for any network and any population size (see research.md
// decision #1 of feature 003-congestion-braess-paradox). Every agent's
// reconsideration is its own independent computation against currently
// observable volumes — never a centrally solved optimal assignment
// (constitution Principle III).
package simulation

// DefaultMaxRounds is a generous safety cap on refinement rounds.
// Best-response dynamics on a congestion game is guaranteed to converge
// in finite steps; this cap exists for pathological inputs, not as the
// primary convergence mechanism.
const DefaultMaxRounds = 1000

// Population is many identical agents sharing one fixed origin and
// destination.
type Population struct {
	Origin, Destination string // Node IDs
	Size                int    // number of agents; 0 is valid
	// MaxRounds is the refinement-round budget. It is required, with no
	// implicit default: 0 means "run the initial incremental loading
	// only, perform zero refinement rounds," which always reports
	// Converged: false, since stability was never actually checked.
	MaxRounds int
}
