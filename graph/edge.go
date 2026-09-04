package graph

// Edge is a directed road connecting two nodes.
type Edge struct {
	ID   string
	From string // Node ID
	To   string // Node ID

	Length   float64
	Capacity float64

	// TravelTime computes this edge's travel time at a given traffic
	// volume. See TravelTimeFunc for the contract implementations must
	// satisfy.
	TravelTime TravelTimeFunc
}
