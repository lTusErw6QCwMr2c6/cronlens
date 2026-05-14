package schedule

import (
	"testing"
)

func TestDescribe_ValidExpression(t *testing.T) {
	descs, err := Describe("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(descs) != 5 {
		t.Fatalf("expected 5 field descriptions, got %d", len(descs))
	}
	if descs[0].Field != "Minute" {
		t.Errorf("expected first field to be Minute, got %q", descs[0].Field)
	}
	if descs[0].Value != "*/5" {
		t.Errorf("expected value */5, got %q", descs[0].Value)
	}
	if descs[0].Meaning != "every 5 minute(s)" {
		t.Errorf("unexpected meaning: %q", descs[0].Meaning)
	}
}

func TestDescribe_InvalidExpression(t *testing.T) {
	_, err := Describe("not-a-cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestDescribe_WildcardMeaning(t *testing.T) {
	descs, err := Describe("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if descs[1].Meaning != "every hour" {
		t.Errorf("expected 'every hour', got %q", descs[1].Meaning)
	}
}

func TestDescribe_RangeMeaning(t *testing.T) {
	descs, err := Describe("0 9-17 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if descs[1].Meaning != "9 through 17" {
		t.Errorf("expected '9 through 17', got %q", descs[1].Meaning)
	}
}

func TestHumanReadable_EveryMinute(t *testing.T) {
	result, err := HumanReadable("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Every minute" {
		t.Errorf("expected 'Every minute', got %q", result)
	}
}

func TestHumanReadable_SpecificFields(t *testing.T) {
	result, err := HumanReadable("0 12 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty human readable string")
	}
}

func TestHumanReadable_Invalid(t *testing.T) {
	_, err := HumanReadable("bad expression here")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}
