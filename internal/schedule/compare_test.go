package schedule

import (
	"testing"
	"time"
)

func TestCompare_InvalidExprA(t *testing.T) {
	_, err := Compare("not-valid", "* * * * *", time.Now(), 5)
	if err == nil {
		t.Fatal("expected error for invalid exprA, got nil")
	}
}

func TestCompare_InvalidExprB(t *testing.T) {
	_, err := Compare("* * * * *", "not-valid", time.Now(), 5)
	if err == nil {
		t.Fatal("expected error for invalid exprB, got nil")
	}
}

func TestCompare_IdenticalExpressions(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := Compare("*/5 * * * *", "*/5 * * * *", from, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Overlaps) != 5 {
		t.Errorf("expected 5 overlaps for identical expressions, got %d", len(result.Overlaps))
	}
	if len(result.OnlyInA) != 0 {
		t.Errorf("expected 0 onlyInA, got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 0 {
		t.Errorf("expected 0 onlyInB, got %d", len(result.OnlyInB))
	}
}

func TestCompare_DisjointExpressions(t *testing.T) {
	// every even hour vs every odd hour — unlikely to overlap in first 5 results
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := Compare("0 0,2,4,6,8 * * *", "0 1,3,5,7,9 * * *", from, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Overlaps) != 0 {
		t.Errorf("expected 0 overlaps for disjoint expressions, got %d", len(result.Overlaps))
	}
	if len(result.OnlyInA) != 5 {
		t.Errorf("expected 5 onlyInA, got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 5 {
		t.Errorf("expected 5 onlyInB, got %d", len(result.OnlyInB))
	}
}

func TestCompare_ResultFields(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := Compare("* * * * *", "*/2 * * * *", from, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExprA != "* * * * *" {
		t.Errorf("unexpected ExprA: %s", result.ExprA)
	}
	if result.ExprB != "*/2 * * * *" {
		t.Errorf("unexpected ExprB: %s", result.ExprB)
	}
	if len(result.NextA) != 10 {
		t.Errorf("expected 10 nextA entries, got %d", len(result.NextA))
	}
	if len(result.NextB) != 10 {
		t.Errorf("expected 10 nextB entries, got %d", len(result.NextB))
	}
}
