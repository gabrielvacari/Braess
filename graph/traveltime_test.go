package graph

import "testing"

// SC-002: travel time at volume 0 equals the edge's configured free-flow
// time.
func TestLinear_FreeFlowAtZeroVolume(t *testing.T) {
	tt := Linear(10, 0.5)
	got, err := tt(0)
	if err != nil {
		t.Fatalf("Linear(10, 0.5)(0) returned unexpected error: %v", err)
	}
	if got != 10 {
		t.Fatalf("Linear(10, 0.5)(0) = %g, want 10", got)
	}
}

func TestConstant_FreeFlowAtZeroVolume(t *testing.T) {
	tt := Constant(45)
	got, err := tt(0)
	if err != nil {
		t.Fatalf("Constant(45)(0) returned unexpected error: %v", err)
	}
	if got != 45 {
		t.Fatalf("Constant(45)(0) = %g, want 45", got)
	}
}

// SC-003: across increasing volume samples, travel time never decreases.
func TestLinear_NonDecreasingAcrossVolumeSamples(t *testing.T) {
	tt := Linear(10, 2)
	volumes := []float64{0, 1, 2, 3, 4}

	var prev float64
	for i, v := range volumes {
		got, err := tt(v)
		if err != nil {
			t.Fatalf("Linear(10, 2)(%g) returned unexpected error: %v", v, err)
		}
		if i > 0 && got < prev {
			t.Fatalf("Linear(10, 2)(%g) = %g, which is less than the previous sample %g", v, got, prev)
		}
		prev = got
	}
}

// SC-003, applied to a constant edge: "non-decreasing" trivially holds
// since the value never changes.
func TestConstant_NonDecreasingAcrossVolumeSamples(t *testing.T) {
	tt := Constant(7)
	volumes := []float64{0, 1, 2, 3, 4}

	var prev float64
	for i, v := range volumes {
		got, err := tt(v)
		if err != nil {
			t.Fatalf("Constant(7)(%g) returned unexpected error: %v", v, err)
		}
		if i > 0 && got < prev {
			t.Fatalf("Constant(7)(%g) = %g, which is less than the previous sample %g", v, got, prev)
		}
		prev = got
	}
}

// FR-008: negative traffic volume must be rejected.
func TestLinear_RejectsNegativeVolume(t *testing.T) {
	tt := Linear(10, 1)
	if _, err := tt(-1); err == nil {
		t.Fatal("Linear(10, 1)(-1) = nil error, want error")
	}
}

func TestConstant_RejectsNegativeVolume(t *testing.T) {
	tt := Constant(10)
	if _, err := tt(-1); err == nil {
		t.Fatal("Constant(10)(-1) = nil error, want error")
	}
}

// FR-004, FR-008, exercised through Graph.TravelTime.
func TestGraph_TravelTime(t *testing.T) {
	g := New()
	mustAddNode(t, g, "a", House)
	mustAddNode(t, g, "b", Company)
	if _, err := g.AddEdge("a-b", "a", "b", 10, 100, Linear(5, 1)); err != nil {
		t.Fatalf("AddEdge returned unexpected error: %v", err)
	}

	got, err := g.TravelTime("a-b", 3)
	if err != nil {
		t.Fatalf("TravelTime(\"a-b\", 3) returned unexpected error: %v", err)
	}
	if want := 8.0; got != want {
		t.Fatalf("TravelTime(\"a-b\", 3) = %g, want %g", got, want)
	}

	if _, err := g.TravelTime("a-b", -1); err == nil {
		t.Fatal("TravelTime(\"a-b\", -1) = nil error, want error")
	}

	if _, err := g.TravelTime("does-not-exist", 0); err == nil {
		t.Fatal("TravelTime for unknown edge id = nil error, want error")
	}
}
