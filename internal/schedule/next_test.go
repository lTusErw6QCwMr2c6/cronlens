package schedule

import (
	"testing"
	"time"
)

var testBase = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestNextOccurrences(t *testing.T) {
	times, err := NextOccurrences("0 * * * *", testBase, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 occurrences, got %d", len(times))
	}
	expected := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	if !times[0].Equal(expected) {
		t.Errorf("first occurrence: got %v, want %v", times[0], expected)
	}
}

func TestNextOccurrences_InvalidExpr(t *testing.T) {
	_, err := NextOccurrences("invalid", testBase, 3)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestNextOccurrence(t *testing.T) {
	next, err := NextOccurrence("30 14 * * *", testBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("got %v, want %v", next, expected)
	}
}

func TestOccurrenceWindow(t *testing.T) {
	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC)

	// Every hour on the hour — expect 6 hits: 00:00 excluded (next after from), 01:00..05:00
	times, err := OccurrenceWindow("0 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 5 {
		t.Errorf("expected 5 occurrences in window, got %d", len(times))
	}
}

func TestOccurrenceWindow_InvalidExpr(t *testing.T) {
	from := time.Now()
	to := from.Add(time.Hour)
	_, err := OccurrenceWindow("bad expr", from, to)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestOccurrenceWindow_EmptyRange(t *testing.T) {
	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := from // zero-length window
	times, err := OccurrenceWindow("* * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 0 {
		t.Errorf("expected 0 occurrences in empty window, got %d", len(times))
	}
}
