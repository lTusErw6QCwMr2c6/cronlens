package schedule

import (
	"strings"
	"testing"
	"time"
)

var conflictBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestFindConflicts_BothEveryMinute(t *testing.T) {
	r, err := FindConflicts("* * * * *", "* * * * *", conflictBase, 5*time.Minute, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.HasConflict {
		t.Error("expected conflicts for identical every-minute schedules")
	}
	if r.Count == 0 {
		t.Error("expected at least one conflict")
	}
}

func TestFindConflicts_Disjoint(t *testing.T) {
	// Every even vs every odd hour — within a 1-hour window they won't overlap
	r, err := FindConflicts("0 0 * * *", "0 12 * * *", conflictBase, 6*time.Hour, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.HasConflict {
		t.Error("expected no conflicts for disjoint daily schedules in 6h window")
	}
}

func TestFindConflicts_InvalidExprA(t *testing.T) {
	_, err := FindConflicts("bad expr", "* * * * *", conflictBase, time.Hour, 5)
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestFindConflicts_InvalidExprB(t *testing.T) {
	_, err := FindConflicts("* * * * *", "bad expr", conflictBase, time.Hour, 5)
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestFindConflicts_MaxResults(t *testing.T) {
	r, err := FindConflicts("* * * * *", "* * * * *", conflictBase, time.Hour, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count > 3 {
		t.Errorf("expected at most 3 results, got %d", r.Count)
	}
}

func TestFormatConflictResult_NoConflict(t *testing.T) {
	r := ConflictResult{ExprA: "0 0 * * *", ExprB: "0 12 * * *", HasConflict: false}
	out := FormatConflictResult(r)
	if !strings.Contains(out, "No conflicts") {
		t.Errorf("expected 'No conflicts' in output, got: %s", out)
	}
}

func TestConflictSummaryLine_WithConflict(t *testing.T) {
	r, _ := FindConflicts("* * * * *", "* * * * *", conflictBase, 5*time.Minute, 5)
	line := ConflictSummaryLine(r)
	if !strings.Contains(line, "conflict") {
		t.Errorf("expected 'conflict' in summary line, got: %s", line)
	}
}
