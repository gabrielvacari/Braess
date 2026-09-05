package queuesim

// AgentReport is one spawned agent's outcome from a Run.
type AgentReport struct {
	// DemandIndex is which input Demand (by index) this agent belongs to.
	DemandIndex int
	// SpawnTime is the simulated time this agent entered the run.
	SpawnTime float64
	// TravelTime is total time spent actually moving (summed across the
	// agent's route), excluding any time spent waiting in a queue.
	TravelTime float64
	// WaitTime is total time spent queued at any signal (FR-006).
	WaitTime float64
	// Arrived is false if the agent was still moving or queued when the
	// run's duration elapsed (FR-008, SC-005) — never presented as
	// having completed its trip in that case.
	Arrived bool
}

// QueueSample is one signal's queue length at one point in simulated
// time — one is recorded per signal per tick (FR-004), giving the full
// time series a queue-length chart needs, not just a final snapshot
// (research.md decision #6).
type QueueSample struct {
	Time     float64
	SignalID string
	Length   int
}

// RunResult is everything a Run call reports.
type RunResult struct {
	QueueSamples []QueueSample
	// Agents is one entry per spawned agent, in spawn order.
	Agents []AgentReport
}
