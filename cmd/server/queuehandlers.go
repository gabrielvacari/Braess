package main

import (
	"encoding/json"
	"net/http"

	"braess/queuesim"
)

// handleQueueRun implements POST /api/queue-run: decode a
// QueueRunRequest, run it through the unchanged queuesim engine (feature
// 008), and encode the result — or a clear error, mirroring handleRun's
// own 400/422 convention (contracts/queue-api-contract.md).
func handleQueueRun(w http.ResponseWriter, r *http.Request) {
	var req QueueRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	g, err := buildQueueGraph(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := queuesim.Run(g, buildSignals(req), buildTimeDemands(req), req.Duration, req.Tick)
	if err != nil {
		// The request was well-formed; queuesim.Run itself failed (e.g. no
		// route for a declared demand — FR-008).
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toQueueRunResponse(result))
}
