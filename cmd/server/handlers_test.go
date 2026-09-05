package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postRun(t *testing.T, body string) (*http.Response, []byte) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handleRun(rec, req)
	resp := rec.Result()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	return resp, buf.Bytes()
}

// A valid request returns 200 with route groups correctly reflecting the
// engine's result (data-model.md "RunResponse").
func TestHandleRun_ValidRequest(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"},{"id":"c1","type":"company"}],
		"edges": [{"id":"e1","from":"h1","to":"c1","length":0,"capacity":0,
			"travelTime":{"type":"linear","freeFlow":0,"slope":1}}],
		"demands": [{"origin":"h1","destination":"c1","size":5}]
	}`

	resp, raw := postRun(t, body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", resp.StatusCode, raw)
	}

	var got RunResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshaling response: %v; body = %s", err, raw)
	}
	if !got.Converged {
		t.Fatalf("Converged = false, want true")
	}
	if len(got.PerDemand) != 1 {
		t.Fatalf("len(PerDemand) = %d, want 1", len(got.PerDemand))
	}
	d := got.PerDemand[0]
	if d.TotalTravelTime != 25 || d.AverageTravelTime != 5 {
		t.Fatalf("demand result = %+v, want total=25 average=5", d)
	}
	if len(d.RouteGroups) != 1 || d.RouteGroups[0].Count != 5 || d.RouteGroups[0].TravelTime != 5 {
		t.Fatalf("route groups = %+v, want one group of 5 agents at travel time 5", d.RouteGroups)
	}
}

// FR-008: a self-loop edge is rejected with 400, before it ever reaches
// the graph package.
func TestHandleRun_SelfLoopRejected(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"}],
		"edges": [{"id":"e1","from":"h1","to":"h1","length":0,"capacity":0,
			"travelTime":{"type":"constant","value":5}}],
		"demands": []
	}`

	resp, raw := postRun(t, body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", resp.StatusCode, raw)
	}
}

// Malformed JSON is rejected with 400.
func TestHandleRun_MalformedJSON(t *testing.T) {
	resp, raw := postRun(t, `{not json`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", resp.StatusCode, raw)
	}
}

// FR-009: a demand with no route between its origin and destination
// returns 422 with a message naming that specific demand — a well-formed
// request that RunDemands itself couldn't satisfy.
func TestHandleRun_NoRouteReturns422NamingTheDemand(t *testing.T) {
	body := `{
		"nodes": [{"id":"h1","type":"house"},{"id":"c1","type":"company"}],
		"edges": [],
		"demands": [{"origin":"h1","destination":"c1","size":3}]
	}`

	resp, raw := postRun(t, body)
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body = %s", resp.StatusCode, raw)
	}

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
