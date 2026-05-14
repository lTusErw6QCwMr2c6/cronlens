package schedule

import (
	"strings"
	"testing"
	"time"
)

var diffBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestDiff_InvalidExprA(t *testing.T) {
	_, err := Diff("bad", "* * * * *", diffBase, 5)
	if err == nil {
		t.Fatal("expected error for invalid expression A")
	}
}

func TestDiff_InvalidExprB(t *testing.T) {
	_, err := Diff("* * * * *", "bad", diffBase, 5)
	if err == nil {
		t.Fatal("expected error for invalid expression B")
	}
}

func TestDiff_IdenticalExpressions(t *testing.T) {
	result, err := Diff("* * * * *", "* * * * *", diffBase, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.OnlyInA) != 0 || len(result.OnlyInB) != 0 {
		t.Errorf("expected no exclusive occurrences for identical expressions")
	}
	if len(result.InBoth) != 10 {
		t.Errorf("expected 10 shared occurrences, got %d", len(result.InBoth))
	}
}

func TestDiff_DisjointExpressions(t *testing.T) {
	// Even hours vs odd hours — these should rarely collide
	result, err := Diff("0 0,2,4,6,8,10,12,14,16,18,20,22 * * *", "0 1,3,5,7,9,11,13,15,17,19,21,23 * * *", diffBase, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.InBoth) != 0 {
		t.Errorf("expected no shared occurrences for disjoint expressions, got %d", len(result.InBoth))
	}
}

func TestDiff_ResultFields(t *testing.T) {
	result, err := Diff("0 * * * *", "30 * * * *", diffBase, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExprA != "0 * * * *" {
		t.Errorf("unexpected ExprA: %s", result.ExprA)
	}
	if result.ExprB != "30 * * * *" {
		t.Errorf("unexpected ExprB: %s", result.ExprB)
	}
}

func TestFormatDiffResult(t *testing.T) {
	result, err := Diff("* * * * *", "0 * * * *", diffBase, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := FormatDiffResult(result)
	if !strings.Contains(out, "Diff:") {
		t.Errorf("expected output to contain 'Diff:', got: %s", out)
	}
}

func TestDiffSummaryLine(t *testing.T) {
	result, err := Diff("* * * * *", "* * * * *", diffBase, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := DiffSummaryLine(result)
	if !strings.Contains(line, "shared:") {
		t.Errorf("expected summary to contain 'shared:', got: %s", line)
	}
}
