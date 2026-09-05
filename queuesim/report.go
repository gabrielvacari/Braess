package queuesim

// AgentReport is one spawned agent's outcome from a Run.
type AgentReport struct {
	// ID is a stable, spawn-order identifier (0-based), assigned to
	// every spawned agent regardless of whether it ever moves. Reports
	// are appended in *completion* order, not spawn order, so ID is what
	// lets a caller match a report back to that same agent's entries in
	// RunResult.Positions (feature 010, data-model.md).
	ID int
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
	// Green is whether this signal was green at Time (feature 010,
	// research.md decision #1) — recorded alongside Length rather than
	// in a separate array, since both already share the same per-tick,
	// per-signal cadence.
	Green bool
}

// PositionSample is one agent's location at one simulated tick: which
// road it occupies, how far along it, and whether it's currently queued
// rather than moving (feature 010, FR-001). A full Run produces one
// ordered sequence of these per agent, from spawn to arrival (or to the
// end of the run, if it never arrives) — enough for a client to
// reconstruct a full time-based playback without recomputing any
// simulation logic itself (constitution Principle III: this only
// records what queuesim's own routing/queuing already decided).
type PositionSample struct {
	Time float64
	// AgentID matches the AgentReport.ID this sample belongs to.
	AgentID int
	EdgeID  string
	// Progress is 0..1 along EdgeID; always 0 while Queued (the agent
	// hasn't started moving along the edge yet).
	Progress float64
	Queued   bool
}

// RunResult is everything a Run call reports.
type RunResult struct {
	QueueSamples []QueueSample
	// Agents is one entry per spawned agent, in spawn order.
	Agents []AgentReport
	// Positions is one entry per in-flight agent per tick, in the order
	// recorded (feature 010, FR-001) — not grouped or sorted further;
	// callers group by AgentID (see web/src/model/queuePlayback.ts).
	Positions []PositionSample
}
