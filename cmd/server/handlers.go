package main

import (
	"encoding/json"
	"net/http"

	"braess/simulation"
)

// handleRun implements POST /api/run: decode a RunRequest, run it through
// the unchanged graph/agent/simulation engine, and encode the result —
// or a clear error, per contracts/api-contract.md.
func handleRun(w http.ResponseWriter, r *http.Request) {
	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	g, err := buildGraph(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	mp := simulation.MultiPopulation{
		Demands:   buildDemands(req),
		MaxRounds: maxRounds(req),
	}

	result, err := simulation.RunDemands(g, mp)
	if err != nil {
		// The request was well-formed; RunDemands itself failed (e.g. no
		// route for a declared demand, naming it — feature 004, FR-008).
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toRunResponse(result))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Encoding failures here would mean a bug in our own response types;
	// nothing useful to do but let the client see a truncated body.
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
