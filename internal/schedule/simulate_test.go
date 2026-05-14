package schedule

import (
	"testing"
	"time"
)

var simBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestSimulate_ValidExpression(t *testing.T) {
	// Every minute for 5 minutes
	result, err := Simulate("* * * * *", simBase, simBase.Add(5*time.Minute), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 5 {
		t.Errorf("expected 5 occurrences, got %d", result.Count)
	}
}

func TestSimulate_InvalidExpression(t *testing.T) {
	_, err := Simulate("bad expr", simBase, simBase.Add(time.Hour), 100)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestSimulate_ToBeforeFrom(t *testing.T) {
	_, err := Simulate("* * * * *", simBase, simBase.Add(-time.Minute), 100)
	if err == nil {
		t.Error("expected error when 'to' is before 'from'")
	}
}

func TestSimulate_MaxResults(t *testing.T) {
	// Every minute for 1 hour, but cap at 10
	result, err := Simulate("* * * * *", simBase, simBase.Add(time.Hour), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 10 {
		t.Errorf("expected 10 occurrences, got %d", result.Count)
	}
}

func TestSimulate_ResultFields(t *testing.T) {
	to := simBase.Add(time.Hour)
	result, err := Simulate("0 * * * *", simBase, to, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Expr != "0 * * * *" {
		t.Errorf("expected expr to be preserved")
	}
	if !result.From.Equal(simBase) {
		t.Errorf("expected From to match")
	}
	if !result.To.Equal(to) {
		t.Errorf("expected To to match")
	}
}

func TestSimulateDuration_Valid(t *testing.T) {
	result, err := SimulateDuration("* * * * *", simBase, 3*time.Minute, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 3 {
		t.Errorf("expected 3 occurrences, got %d", result.Count)
	}
}

func TestSimulateDuration_NegativeDuration(t *testing.T) {
	_, err := SimulateDuration("* * * * *", simBase, -time.Minute, 100)
	if err == nil {
		t.Error("expected error for negative duration")
	}
}
