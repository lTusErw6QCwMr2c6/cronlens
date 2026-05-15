package schedule

import (
	"strings"
	"testing"
	"time"
)

func TestFormatConflictResult_WithConflicts(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 1, 0, 0, time.UTC)
	r := ConflictResult{
		ExprA:       "* * * * *",
		ExprB:       "* * * * *",
		Conflicts:   []time.Time{ts},
		Count:       1,
		HasConflict: true,
	}
	out := FormatConflictResult(r)
	if !strings.Contains(out, "1 conflict") {
		t.Errorf("expected '1 conflict' in output, got: %s", out)
	}
	if !strings.Contains(out, "* * * * *") {
		t.Errorf("expected expressions in output, got: %s", out)
	}
}

func TestConflictSummaryLine_NoConflict(t *testing.T) {
	r := ConflictResult{
		ExprA:       "0 0 * * *",
		ExprB:       "0 12 * * *",
		HasConflict: false,
		Count:       0,
	}
	line := ConflictSummaryLine(r)
	if !strings.Contains(line, "no conflicts") {
		t.Errorf("expected 'no conflicts', got: %s", line)
	}
	if !strings.Contains(line, "0 0 * * *") {
		t.Errorf("expected exprA in summary, got: %s", line)
	}
}

func TestConflictSummaryLine_ShowsFirstConflict(t *testing.T) {
	ts := time.Date(2024, 6, 15, 8, 30, 0, 0, time.UTC)
	r := ConflictResult{
		ExprA:       "30 8 * * *",
		ExprB:       "30 8 * * *",
		Conflicts:   []time.Time{ts},
		Count:       1,
		HasConflict: true,
	}
	line := ConflictSummaryLine(r)
	if !strings.Contains(line, "1 conflict") {
		t.Errorf("expected count in summary, got: %s", line)
	}
	if !strings.Contains(line, "first at") {
		t.Errorf("expected 'first at' in summary, got: %s", line)
	}
}
