package schedule_test

import (
	"testing"
	"time"

	"github.com/cronlens/cronlens/internal/schedule"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestParse_ValidExpressions(t *testing.T) {
	cases := []string{
		"* * * * *",
		"0 9 * * 1-5",
		"*/15 * * * *",
		"@hourly",
		"@daily",
		"0 0 1 * *",
	}

	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			entry, err := schedule.Parse(expr)
			if err != nil {
				t.Fatalf("expected no error for %q, got: %v", expr, err)
			}
			if entry == nil {
				t.Fatal("expected non-nil entry")
			}
		})
	}
}

func TestParse_InvalidExpressions(t *testing.T) {
	cases := []string{
		"not a cron",
		"99 * * * *",
		"",
	}

	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			_, err := schedule.Parse(expr)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", expr)
			}
		})
	}
}

func TestNextN(t *testing.T) {
	entry, err := schedule.Parse("0 * * * *")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	times := entry.NextN(fixedTime, 5)
	if len(times) != 5 {
		t.Fatalf("expected 5 times, got %d", len(times))
	}

	for i := 1; i < len(times); i++ {
		if !times[i].After(times[i-1]) {
			t.Errorf("times[%d] is not after times[%d]", i, i-1)
		}
	}
}

func TestNextFrom(t *testing.T) {
	entry, err := schedule.Parse("0 12 * * *")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	next := entry.NextFrom(fixedTime)
	if next.Hour() != 12 || next.Minute() != 0 {
		t.Errorf("expected next at 12:00, got %v", next)
	}
}
