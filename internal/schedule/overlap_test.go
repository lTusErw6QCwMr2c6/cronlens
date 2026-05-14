package schedule

import (
	"testing"
	"time"
)

func TestFindOverlaps_BothEveryMinute(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	window := 5 * time.Minute

	result, err := FindOverlaps("* * * * *", "* * * * *", from, window, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalFound == 0 {
		t.Error("expected overlaps for identical every-minute schedules")
	}
}

func TestFindOverlaps_Disjoint(t *testing.T) {
	// One fires at minute 0, other at minute 30
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	window := 59 * time.Minute

	result, err := FindOverlaps("0 * * * *", "30 * * * *", from, window, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalFound != 0 {
		t.Errorf("expected 0 overlaps, got %d", result.TotalFound)
	}
}

func TestFindOverlaps_MaxResults(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	window := 60 * time.Minute

	result, err := FindOverlaps("* * * * *", "* * * * *", from, window, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalFound > 3 {
		t.Errorf("expected at most 3 results, got %d", result.TotalFound)
	}
}

func TestFindOverlaps_InvalidExprA(t *testing.T) {
	from := time.Now()
	_, err := FindOverlaps("not valid", "* * * * *", from, time.Hour, 5)
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestFindOverlaps_InvalidExprB(t *testing.T) {
	from := time.Now()
	_, err := FindOverlaps("* * * * *", "also bad", from, time.Hour, 5)
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestFindOverlaps_ResultFields(t *testing.T) {
	from := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	result, err := FindOverlaps("0 12 * * *", "0 12 * * *", from, 25*time.Hour, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExprA != "0 12 * * *" {
		t.Errorf("ExprA mismatch: %s", result.ExprA)
	}
	if result.ExprB != "0 12 * * *" {
		t.Errorf("ExprB mismatch: %s", result.ExprB)
	}
	if result.TotalFound == 0 {
		t.Error("expected at least one overlap for daily schedules over 25h window")
	}
}
