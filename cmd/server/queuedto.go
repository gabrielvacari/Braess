package main

import (
	"fmt"

	"braess/graph"
	"braess/queuesim"
)

// SignalDTO is the JSON shape of a queuesim.Signal.
type SignalDTO struct {
	ID            string  `json:"id"`
	EdgeID        string  `json:"edgeId"`
	GreenDuration float64 `json:"greenDuration"`
	RedDuration   float64 `json:"redDuration"`
	Offset        float64 `json:"offset,omitempty"`
}

// TimeDemandDTO is the JSON shape of a queuesim.Demand.
type TimeDemandDTO struct {
	Origin          string  `json:"origin"`
	Destination     string  `json:"destination"`
	Count           int     `json:"count"`
	ArrivalInterval float64 `json:"arrivalInterval"`
}

// QueueRunRequest is the POST /api/queue-run request body. Nodes and
// Edges reuse dto.go's existing NodeDTO/EdgeDTO exactly — a node or a
// directed road means the same thing to both endpoints, and this
// endpoint's roads are not expanded from a bidirectional client concept
// the way /api/run's are (research.md decision #1, #5): one Edge here is
// one directed road, exactly as queuesim itself models it.
type QueueRunRequest struct {
	Nodes    []NodeDTO       `json:"nodes"`
	Edges    []EdgeDTO       `json:"edges"`
	Signals  []SignalDTO     `json:"signals"`
	Demands  []TimeDemandDTO `json:"demands"`
	Duration float64         `json:"duration"`
	Tick     float64         `json:"tick"`
}

// QueueSampleDTO mirrors queuesim.QueueSample field-for-field.
type QueueSampleDTO struct {
	Time     float64 `json:"time"`
	SignalID string  `json:"signalId"`
	Length   int     `json:"length"`
	// Green mirrors QueueSample.Green (feature 010, FR-002).
	Green bool `json:"green"`
}

// AgentReportDTO mirrors queuesim.AgentReport field-for-field.
type AgentReportDTO struct {
	// ID mirrors AgentReport.ID (feature 010) — matches PositionSampleDTO.AgentID.
	ID          int     `json:"id"`
	DemandIndex int     `json:"demandIndex"`
	SpawnTime   float64 `json:"spawnTime"`
	TravelTime  float64 `json:"travelTime"`
	WaitTime    float64 `json:"waitTime"`
	Arrived     bool    `json:"arrived"`
}

// PositionSampleDTO mirrors queuesim.PositionSample field-for-field
// (feature 010, FR-001).
type PositionSampleDTO struct {
	Time     float64 `json:"time"`
	AgentID  int     `json:"agentId"`
	EdgeID   string  `json:"edgeId"`
	Progress float64 `json:"progress"`
	Queued   bool    `json:"queued"`
}

// QueueRunResponse is the POST /api/queue-run success response body.
type QueueRunResponse struct {
	QueueSamples []QueueSampleDTO    `json:"queueSamples"`
	Agents       []AgentReportDTO    `json:"agents"`
	Positions    []PositionSampleDTO `json:"positions"`
}

// buildSignals maps a QueueRunRequest's signals to []queuesim.Signal.
func buildSignals(req QueueRunRequest) []queuesim.Signal {
	signals := make([]queuesim.Signal, len(req.Signals))
	for i, s := range req.Signals {
		signals[i] = queuesim.Signal{
			ID:            s.ID,
			EdgeID:        s.EdgeID,
			GreenDuration: s.GreenDuration,
			RedDuration:   s.RedDuration,
			Offset:        s.Offset,
		}
	}
	return signals
}

// buildTimeDemands maps a QueueRunRequest's demands to []queuesim.Demand.
func buildTimeDemands(req QueueRunRequest) []queuesim.Demand {
	demands := make([]queuesim.Demand, len(req.Demands))
	for i, d := range req.Demands {
		demands[i] = queuesim.Demand{
			Origin:          d.Origin,
			Destination:     d.Destination,
			Count:           d.Count,
			ArrivalInterval: d.ArrivalInterval,
		}
	}
	return demands
}

// toQueueRunResponse maps a queuesim.RunResult to a QueueRunResponse.
func toQueueRunResponse(result queuesim.RunResult) QueueRunResponse {
	samples := make([]QueueSampleDTO, len(result.QueueSamples))
	for i, s := range result.QueueSamples {
		samples[i] = QueueSampleDTO{Time: s.Time, SignalID: s.SignalID, Length: s.Length, Green: s.Green}
	}
	agents := make([]AgentReportDTO, len(result.Agents))
	for i, a := range result.Agents {
		agents[i] = AgentReportDTO{
			ID: a.ID, DemandIndex: a.DemandIndex, SpawnTime: a.SpawnTime,
			TravelTime: a.TravelTime, WaitTime: a.WaitTime, Arrived: a.Arrived,
		}
	}
	positions := make([]PositionSampleDTO, len(result.Positions))
	for i, p := range result.Positions {
		positions[i] = PositionSampleDTO{Time: p.Time, AgentID: p.AgentID, EdgeID: p.EdgeID, Progress: p.Progress, Queued: p.Queued}
	}
	return QueueRunResponse{QueueSamples: samples, Agents: agents, Positions: positions}
}

// buildQueueGraph builds a *graph.Graph from a QueueRunRequest's nodes
// and edges by delegating to the existing buildGraph (dto.go, feature
// 005) — QueueRunRequest.Nodes/.Edges are the same NodeDTO/EdgeDTO types,
// so no new graph-construction logic is needed (research.md decision #5).
func buildQueueGraph(req QueueRunRequest) (*graph.Graph, error) {
	g, err := buildGraph(RunRequest{Nodes: req.Nodes, Edges: req.Edges})
	if err != nil {
		return nil, fmt.Errorf("building network: %w", err)
	}
	return g, nil
}
