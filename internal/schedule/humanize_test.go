package schedule

import (
	"testing"
	"time"
)

func baseTime() time.Time {
	// Monday 2024-01-15 10:00:00 UTC
	return time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
}

func TestHumanizeNext_InSeconds(t *testing.T) {
	from := baseTime()
	// fires every minute; next is at :01 — 60s away
	result, err := HumanizeNext("* * * * *", from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty humanized string")
	}
}

func TestHumanizeNext_InvalidExpr(t *testing.T) {
	_, err := HumanizeNext("invalid expr", baseTime())
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestHumanizeNext_Tomorrow(t *testing.T) {
	// fires daily at 08:00; from is 10:00 so next is tomorrow
	from := baseTime()
	result, err := HumanizeNext("0 8 * * *", from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "tomorrow at 08:00" {
		t.Errorf("expected 'tomorrow at 08:00', got %q", result)
	}
}

func TestHumanizeNext_SameDay(t *testing.T) {
	// fires at 12:00 daily; from is 10:00 — 2 hours away
	from := baseTime()
	result, err := HumanizeNext("0 12 * * *", from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "in 2 hours at 12:00" {
		t.Errorf("expected 'in 2 hours at 12:00', got %q", result)
	}
}

func TestHumanizeNextN_ReturnsMultiple(t *testing.T) {
	from := baseTime()
	results, err := HumanizeNextN("0 * * * *", from, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		if r == "" {
			t.Errorf("result[%d] is empty", i)
		}
	}
}

func TestHumanizeNextN_InvalidExpr(t *testing.T) {
	_, err := HumanizeNextN("bad", baseTime(), 3)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestHumanizeTime_InMinutes(t *testing.T) {
	from := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	target := from.Add(5 * time.Minute)
	result := humanizeTime(target, from)
	if result != "in 5 minutes" {
		t.Errorf("expected 'in 5 minutes', got %q", result)
	}
}

func TestHumanizeTime_InPast(t *testing.T) {
	from := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	past := from.Add(-1 * time.Minute)
	result := humanizeTime(past, from)
	if result != "in the past" {
		t.Errorf("expected 'in the past', got %q", result)
	}
}
