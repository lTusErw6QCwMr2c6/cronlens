package schedule

import (
	"testing"
)

func TestResolveAlias_Known(t *testing.T) {
	expr, ok := ResolveAlias("@daily")
	if !ok {
		t.Fatal("expected @daily to resolve")
	}
	if expr != "0 0 * * *" {
		t.Errorf("expected '0 0 * * *', got %q", expr)
	}
}

func TestResolveAlias_Unknown(t *testing.T) {
	expr, ok := ResolveAlias("@unknown")
	if ok {
		t.Fatal("expected ok=false for unknown alias")
	}
	if expr != "@unknown" {
		t.Errorf("expected original value returned, got %q", expr)
	}
}

func TestLookupAlias_Found(t *testing.T) {
	a, err := LookupAlias("@hourly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Expression != "0 * * * *" {
		t.Errorf("unexpected expression %q", a.Expression)
	}
	if a.Description == "" {
		t.Error("expected non-empty description")
	}
}

func TestLookupAlias_NotFound(t *testing.T) {
	_, err := LookupAlias("@nope")
	if err == nil {
		t.Fatal("expected error for unknown alias")
	}
}

func TestListAliases_NotEmpty(t *testing.T) {
	aliases := ListAliases()
	if len(aliases) == 0 {
		t.Fatal("expected at least one alias")
	}
}

func TestListAliases_ContainsYearly(t *testing.T) {
	aliases := ListAliases()
	for _, a := range aliases {
		if a.Name == "@yearly" {
			return
		}
	}
	t.Error("expected @yearly in alias list")
}

func TestAliasForExpr_Found(t *testing.T) {
	a, ok := AliasForExpr("* * * * *")
	if !ok {
		t.Fatal("expected a match for '* * * * *'")
	}
	if a.Name != "@every_minute" {
		t.Errorf("expected @every_minute, got %q", a.Name)
	}
}

func TestAliasForExpr_NotFound(t *testing.T) {
	_, ok := AliasForExpr("5 4 * * 1")
	if ok {
		t.Error("expected no alias match for custom expression")
	}
}

func TestAnnuallyAndYearlySameExpression(t *testing.T) {
	yearly, _ := ResolveAlias("@yearly")
	annually, _ := ResolveAlias("@annually")
	if yearly != annually {
		t.Errorf("@yearly (%q) and @annually (%q) should share expression", yearly, annually)
	}
}
