package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postQueueRun(t *testing.T, body string) (*http.Response, []byte) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/queue-run", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handleQueueRun(rec, req)
	resp := rec.Result()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	return resp, buf.Bytes()
}

// T011: a valid request returns 200 with queue samples and agent reports
// reflecting real, sustained queuing — mirroring queuesim's own
// two-road, opposite-phase scenario (queuesim/two_road_test.go) at a
// smaller scale, through the HTTP boundary this time.
func TestHandleQueueRun_ValidRequestReflectsRealQueuing(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"},{"id":"c1","type":"company"}],
		"edges": [
			{"id":"road-a","from":"h1","to":"c1","length":0,"capacity":0,"travelTime":{"type":"constant","value":1}},
			{"id":"road-b","from":"h1","to":"c1","length":0,"capacity":0,"travelTime":{"type":"constant","value":1}}
		],
		"signals": [
			{"id":"signal-a","edgeId":"road-a","greenDuration":5,"redDuration":5,"offset":5},
			{"id":"signal-b","edgeId":"road-b","greenDuration":5,"redDuration":5,"offset":0}
		],
		"demands": [{"origin":"h1","destination":"c1","count":150,"arrivalInterval":0.3}],
		"duration": 70,
		"tick": 0.25
	}`

	resp, raw := postQueueRun(t, body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", resp.StatusCode, raw)
	}

	var got QueueRunResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshaling response: %v; body = %s", err, raw)
	}

	if len(got.Agents) != 150 {
		t.Fatalf("len(Agents) = %d, want 150", len(got.Agents))
	}

	sawWait := false
	for _, a := range got.Agents {
		if a.WaitTime > 0 {
			sawWait = true
		}
	}
	if !sawWait {
		t.Fatal("no agent recorded any wait time — the queue mechanic never engaged")
	}

	sawQueue := false
	for _, s := range got.QueueSamples {
		if s.Length > 0 {
			sawQueue = true
		}
	}
	if !sawQueue {
		t.Fatal("no queue sample shows a nonzero length — scenario didn't exercise any congestion")
	}
}

// Feature 010, FR-001/FR-002: the response includes per-tick agent
// positions and each queue sample's signal phase, not just the
// feature-009 summary fields.
func TestHandleQueueRun_IncludesPositionsAndSignalPhase(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"},{"id":"c1","type":"company"}],
		"edges": [
			{"id":"road-a","from":"h1","to":"c1","length":0,"capacity":0,"travelTime":{"type":"constant","value":1}},
			{"id":"road-b","from":"h1","to":"c1","length":0,"capacity":0,"travelTime":{"type":"constant","value":1}}
		],
		"signals": [
			{"id":"signal-a","edgeId":"road-a","greenDuration":5,"redDuration":5,"offset":5},
			{"id":"signal-b","edgeId":"road-b","greenDuration":5,"redDuration":5,"offset":0}
		],
		"demands": [{"origin":"h1","destination":"c1","count":150,"arrivalInterval":0.3}],
		"duration": 70,
		"tick": 0.25
	}`

	resp, raw := postQueueRun(t, body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", resp.StatusCode, raw)
	}

	var got QueueRunResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshaling response: %v; body = %s", err, raw)
	}

	if len(got.Positions) == 0 {
		t.Fatal("Positions is empty, want at least one recorded position sample")
	}
	for _, p := range got.Positions {
		if p.Progress < 0 || p.Progress > 1 {
			t.Fatalf("position sample out of bounds: %+v", p)
		}
	}

	sawGreen, sawRed := false, false
	for _, s := range got.QueueSamples {
		if s.Green {
			sawGreen = true
		} else {
			sawRed = true
		}
	}
	if !sawGreen || !sawRed {
		t.Fatal("QueueSamples never shows both a green and a red sample")
	}

	// Every agent's ID appears among the position samples (unless it
	// never got a route at all, which doesn't happen in this scenario).
	agentIDs := make(map[int]bool)
	for _, p := range got.Positions {
		agentIDs[p.AgentID] = true
	}
	for _, a := range got.Agents {
		if !agentIDs[a.ID] {
			t.Fatalf("agent ID %d has no position samples", a.ID)
		}
	}
}

// FR-008 (T015): a demand with no route to its destination returns 422
// with a clear message, mirroring feature 005/008's own error-reporting
// convention.
func TestHandleQueueRun_NoRouteReturns422(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"},{"id":"c1","type":"company"}],
		"edges": [],
		"signals": [],
		"demands": [{"origin":"h1","destination":"c1","count":5,"arrivalInterval":1}],
		"duration": 10,
		"tick": 1
	}`

	resp, raw := postQueueRun(t, body)
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body = %s", resp.StatusCode, raw)
	}

	// writeError (handlers.go) always encodes an ErrorResponse regardless
	// of which endpoint called it — no separate QueueErrorResponse type
	// exists, since the shape would be identical.
	var errResp ErrorResponse
	if err := json.Unmarshal(raw, &errResp); err != nil {
		t.Fatalf("unmarshaling error response: %v; body = %s", err, raw)
	}
	for _, want := range []string{"h1", "c1"} {
		if !bytes.Contains([]byte(errResp.Error), []byte(want)) {
			t.Fatalf("error %q does not name the failing demand (h1 -> c1)", errResp.Error)
		}
	}
}

// A self-loop edge is rejected with 400, before it ever reaches
// queuesim.Run — mirroring handleRun's own FR-008 precedent.
func TestHandleQueueRun_SelfLoopRejected(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"}],
		"edges": [{"id":"e1","from":"h1","to":"h1","length":0,"capacity":0,"travelTime":{"type":"constant","value":5}}],
		"signals": [],
		"demands": [],
		"duration": 10,
		"tick": 1
	}`

	resp, raw := postQueueRun(t, body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", resp.StatusCode, raw)
	}
}

// Malformed JSON is rejected with 400.
func TestHandleQueueRun_MalformedJSON(t *testing.T) {
	resp, raw := postQueueRun(t, `{not json`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", resp.StatusCode, raw)
	}
}
