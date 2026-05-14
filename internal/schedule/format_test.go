package schedule_test

import (
	"testing"
	"time"

	"github.com/cronlens/cronlens/internal/schedule"
)

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0s"},
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m 30s"},
		{1*time.Hour + 5*time.Minute, "1h 5m"},
		{25 * time.Hour, "1d 1h"},
		{-time.Minute, "overdue"},
	}

	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			got := schedule.FormatDuration(tc.duration)
			if got != tc.expected {
				t.Errorf("FormatDuration(%v) = %q, want %q", tc.duration, got, tc.expected)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	ts := time.Date(2024, 6, 1, 14, 30, 0, 0, time.UTC)
	got := schedule.FormatTime(ts)
	expected := "2024-06-01 14:30:00 UTC"
	if got != expected {
		t.Errorf("FormatTime() = %q, want %q", got, expected)
	}
}

func TestSummaryLine(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	next := now.Add(2 * time.Hour)

	line := schedule.SummaryLine("0 * * * *", next, now)
	if line == "" {
		t.Error("expected non-empty summary line")
	}
	if len(line) < 10 {
		t.Errorf("summary line too short: %q", line)
	}
}
