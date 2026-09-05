package queuesim

import "testing"

func TestSignal_IsGreenAt(t *testing.T) {
	s := Signal{GreenDuration: 10, RedDuration: 5} // 15s cycle

	tests := []struct {
		time float64
		want bool
	}{
		{0, true},
		{5, true},
		{9.999, true},
		{10, false}, // boundary: green ends exactly at GreenDuration
		{12, false},
		{14.999, false},
		{15, true}, // next cycle starts
		{30, true},
		{35, true},  // 35 mod 15 = 5, still within this cycle's green window
		{40, false}, // 40 mod 15 = 10, into this cycle's red window
	}

	for _, tc := range tests {
		if got := s.IsGreenAt(tc.time); got != tc.want {
			t.Errorf("IsGreenAt(%g) = %v, want %v", tc.time, got, tc.want)
		}
	}
}

// research.md decision #4: two signals sharing a cycle, offset by one
// GreenDuration, are never green at the same time — "opposite" falls out
// of configuration alone, no linking field needed.
func TestSignal_OppositePairNeverBothGreen(t *testing.T) {
	a := Signal{GreenDuration: 10, RedDuration: 10} // 20s cycle
	b := Signal{GreenDuration: 10, RedDuration: 10, Offset: 10}

	for tm := 0.0; tm < 100; tm += 0.5 {
		if a.IsGreenAt(tm) && b.IsGreenAt(tm) {
			t.Fatalf("both signals green at t=%g", tm)
		}
	}

	// And between them, something is always green (no dead gap either).
	for tm := 0.0; tm < 100; tm += 0.5 {
		if !a.IsGreenAt(tm) && !b.IsGreenAt(tm) {
			t.Fatalf("neither signal green at t=%g", tm)
		}
	}
}

func TestSignal_TimeUntilNextGreen(t *testing.T) {
	s := Signal{GreenDuration: 10, RedDuration: 5}

	if got := s.timeUntilNextGreen(3); got != 0 {
		t.Fatalf("timeUntilNextGreen(3) = %g, want 0 (already green)", got)
	}
	if got := s.timeUntilNextGreen(10); got != 5 {
		t.Fatalf("timeUntilNextGreen(10) = %g, want 5", got)
	}
	if got := s.timeUntilNextGreen(13); got != 2 {
		t.Fatalf("timeUntilNextGreen(13) = %g, want 2", got)
	}
}
