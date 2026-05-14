package schedule_test

import (
	"testing"
	"time"

	"github.com/user/cronlens/internal/schedule"
)

func TestNextInTimezones_BasicUTC(t *testing.T) {
	from := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	result, err := schedule.NextInTimezones("30 14 * * *", []string{"UTC"}, from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Zones) != 1 {
		t.Fatalf("expected 1 zone, got %d", len(result.Zones))
	}
	if result.Zones[0].Next.Hour() != 14 || result.Zones[0].Next.Minute() != 30 {
		t.Errorf("expected next at 14:30, got %v", result.Zones[0].Next)
	}
}

func TestNextInTimezones_MultipleZones(t *testing.T) {
	from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	zones := []string{"UTC", "America/New_York", "Europe/London"}
	result, err := schedule.NextInTimezones("0 9 * * 1", zones, from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Zones) != 3 {
		t.Fatalf("expected 3 zones, got %d", len(result.Zones))
	}
	for _, z := range result.Zones {
		if z.Next.IsZero() {
			t.Errorf("zone %s has zero next time", z.ZoneName)
		}
		if z.Offset == "" {
			t.Errorf("zone %s has empty offset", z.ZoneName)
		}
	}
}

func TestNextInTimezones_InvalidExpression(t *testing.T) {
	from := time.Now()
	_, err := schedule.NextInTimezones("not-a-cron", []string{"UTC"}, from)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestNextInTimezones_InvalidTimezone(t *testing.T) {
	from := time.Now()
	_, err := schedule.NextInTimezones("* * * * *", []string{"Not/AZone"}, from)
	if err == nil {
		t.Error("expected error for invalid timezone, got nil")
	}
}

func TestNextInTimezones_OffsetFormat(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := schedule.NextInTimezones("* * * * *", []string{"Asia/Kolkata"}, from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	offset := result.Zones[0].Offset
	if len(offset) == 0 {
		t.Error("expected non-empty offset string")
	}
	if offset[:3] != "UTC" {
		t.Errorf("expected offset to start with UTC, got %q", offset)
	}
}
