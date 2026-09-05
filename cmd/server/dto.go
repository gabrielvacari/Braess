// Command server is a stateless HTTP API that translates JSON in and out
// of the existing graph/agent/simulation types. It is the only Go code
// in this project that knows HTTP or JSON exist — graph, agent, and
// simulation are unmodified by it and have no dependency on it
// (constitution Principle I). It holds no session state: every request
// carries the whole network and demand set it needs (research.md
// decision #1, feature 005-web-frontend).
package main

import (
	"fmt"

	"braess/agent"
	"braess/graph"
	"braess/simulation"
)

// NodeDTO is the JSON shape of a graph.Node.
type NodeDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "house" | "company" | "intersection"
}

// TravelTimeDTO is the JSON shape of a graph.TravelTimeFunc: a named
// shape plus whichever parameters that shape uses. Only the fields the
// selected Type uses are read; the others are ignored rather than
// rejected.
type TravelTimeDTO struct {
	Type     string  `json:"type"` // "linear" | "constant"
	FreeFlow float64 `json:"freeFlow,omitempty"`
	Slope    float64 `json:"slope,omitempty"`
	Value    float64 `json:"value,omitempty"`
}

// EdgeDTO is the JSON shape of a graph.Edge.
type EdgeDTO struct {
	ID         string        `json:"id"`
	From       string        `json:"from"`
	To         string        `json:"to"`
	Length     float64       `json:"length"`
	Capacity   float64       `json:"capacity"`
	TravelTime TravelTimeDTO `json:"travelTime"`
}

// DemandDTO is the JSON shape of a simulation.Demand.
type DemandDTO struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Size        int    `json:"size"`
}

// RunRequest is the POST /api/run request body.
type RunRequest struct {
	Nodes     []NodeDTO   `json:"nodes"`
	Edges     []EdgeDTO   `json:"edges"`
	Demands   []DemandDTO `json:"demands"`
	MaxRounds int         `json:"maxRounds,omitempty"`
}

// RouteGroupDTO is one distinct route within a demand's result, and how
// many of that demand's agents took exactly it (research.md decision #4:
// grouping keeps the response bounded by the number of distinct routes,
// not the population size).
type RouteGroupDTO struct {
	EdgeIDs    []string `json:"edgeIds"`
	Count      int      `json:"count"`
	TravelTime float64  `json:"travelTime"`
}

// DemandResultDTO is one demand's own outcome within a RunResponse.
type DemandResultDTO struct {
	Origin            string          `json:"origin"`
	Destination       string          `json:"destination"`
	Size              int             `json:"size"`
	RouteGroups       []RouteGroupDTO `json:"routeGroups"`
	TotalTravelTime   float64         `json:"totalTravelTime"`
	AverageTravelTime float64         `json:"averageTravelTime"`
}

// RunResponse is the POST /api/run success response body.
type RunResponse struct {
	Converged         bool              `json:"converged"`
	Rounds            int               `json:"rounds"`
	TotalTravelTime   float64           `json:"totalTravelTime"`
	AverageTravelTime float64           `json:"averageTravelTime"`
	PerDemand         []DemandResultDTO `json:"perDemand"`
}

// ErrorResponse is the body of a non-2xx POST /api/run response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// nodeType maps a NodeDTO's string Type to a graph.NodeType.
func nodeType(s string) (graph.NodeType, error) {
	switch s {
	case "house":
		return graph.House, nil
	case "company":
		return graph.Company, nil
	case "intersection":
		return graph.Intersection, nil
	default:
		return 0, fmt.Errorf("unknown node type %q", s)
	}
}

// travelTimeFunc maps a TravelTimeDTO to a graph.TravelTimeFunc.
func travelTimeFunc(dto TravelTimeDTO) (graph.TravelTimeFunc, error) {
	switch dto.Type {
	case "linear":
		return graph.Linear(dto.FreeFlow, dto.Slope), nil
	case "constant":
		return graph.Constant(dto.Value), nil
	default:
		return nil, fmt.Errorf("unknown travelTime type %q", dto.Type)
	}
}

// buildGraph builds a *graph.Graph from a RunRequest's nodes and edges.
// It rejects a self-loop edge (From == To) before it ever reaches
// graph.AddEdge, per FR-008 — graph itself doesn't forbid self-loops,
// since that's a web-editing concern, not a graph-structure one.
func buildGraph(req RunRequest) (*graph.Graph, error) {
	g := graph.New()

	for _, n := range req.Nodes {
		t, err := nodeType(n.Type)
		if err != nil {
			return nil, fmt.Errorf("node %q: %w", n.ID, err)
		}
		if _, err := g.AddNode(n.ID, t); err != nil {
			return nil, err
		}
	}

	for _, e := range req.Edges {
		if e.From == e.To {
			return nil, fmt.Errorf("edge %q: a road cannot connect a node to itself", e.ID)
		}
		tt, err := travelTimeFunc(e.TravelTime)
		if err != nil {
			return nil, fmt.Errorf("edge %q: %w", e.ID, err)
		}
		if _, err := g.AddEdge(e.ID, e.From, e.To, e.Length, e.Capacity, tt); err != nil {
			return nil, err
		}
	}

	return g, nil
}

// buildDemands maps a RunRequest's demands to []simulation.Demand.
func buildDemands(req RunRequest) []simulation.Demand {
	demands := make([]simulation.Demand, len(req.Demands))
	for i, d := range req.Demands {
		demands[i] = simulation.Demand{Origin: d.Origin, Destination: d.Destination, Size: d.Size}
	}
	return demands
}

// maxRounds resolves RunRequest.MaxRounds to the round budget RunDemands
// should use: an explicit, non-zero value is honored as given; 0 (the
// zero value for an omitted JSON field) falls back to
// simulation.DefaultMaxRounds, an HTTP-layer convenience so a typical
// client never has to think about it (unlike simulation.MultiPopulation.
// MaxRounds itself, which is always honored literally — see data-model.md
// "RunRequest").
func maxRounds(req RunRequest) int {
	if req.MaxRounds == 0 {
		return simulation.DefaultMaxRounds
	}
	return req.MaxRounds
}

// toRunResponse maps a simulation.MultiAssignmentResult to a RunResponse,
// grouping each demand's routes by identical edge-ID sequence.
func toRunResponse(result simulation.MultiAssignmentResult) RunResponse {
	perDemand := make([]DemandResultDTO, len(result.PerDemand))
	for i, dr := range result.PerDemand {
		perDemand[i] = DemandResultDTO{
			Origin:            dr.Demand.Origin,
			Destination:       dr.Demand.Destination,
			Size:              dr.Demand.Size,
			RouteGroups:       groupRoutes(dr.Routes),
			TotalTravelTime:   dr.TotalTravelTime,
			AverageTravelTime: dr.AverageTravelTime,
		}
	}

	return RunResponse{
		Converged:         result.Converged,
		Rounds:            result.Rounds,
		TotalTravelTime:   result.TotalTravelTime,
		AverageTravelTime: result.AverageTravelTime,
		PerDemand:         perDemand,
	}
}

// groupRoutes groups routes by identical edge-ID sequence, preserving the
// order each distinct sequence was first seen in, so the response is
// deterministic given a deterministic input order.
func groupRoutes(routes []agent.Route) []RouteGroupDTO {
	var groups []RouteGroupDTO
	index := make(map[string]int) // edge-ID-sequence key -> index into groups

	for _, r := range routes {
		ids := make([]string, len(r.Edges))
		for i, e := range r.Edges {
			ids[i] = e.ID
		}
		key := fmt.Sprint(ids)

		if i, ok := index[key]; ok {
			groups[i].Count++
			continue
		}

		index[key] = len(groups)
		groups = append(groups, RouteGroupDTO{
			EdgeIDs:    ids,
			Count:      1,
			TravelTime: r.TotalTravelTime,
		})
	}

	return groups
}
