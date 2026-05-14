package schedule_test

import (
	"testing"

	"cronlens/internal/schedule"
)

func TestSuggestTemplates_EmptyQuery(t *testing.T) {
	results := schedule.SuggestTemplates("")
	if len(results) != len(schedule.CommonTemplates) {
		t.Errorf("expected all templates for empty query, got %d", len(results))
	}
}

func TestSuggestTemplates_ByName(t *testing.T) {
	results := schedule.SuggestTemplates("weekly")
	if len(results) == 0 {
		t.Fatal("expected at least one result for 'weekly'")
	}
	for _, r := range results {
		if r.Name == "" {
			t.Error("result name should not be empty")
		}
	}
}

func TestSuggestTemplates_ByDescription(t *testing.T) {
	results := schedule.SuggestTemplates("midnight")
	if len(results) == 0 {
		t.Fatal("expected at least one result for 'midnight'")
	}
}

func TestSuggestTemplates_NoMatch(t *testing.T) {
	results := schedule.SuggestTemplates("zzznomatch")
	if len(results) != 0 {
		t.Errorf("expected no results for nonsense query, got %d", len(results))
	}
}

func TestSuggestTemplates_CaseInsensitive(t *testing.T) {
	lower := schedule.SuggestTemplates("daily")
	upper := schedule.SuggestTemplates("DAILY")
	if len(lower) != len(upper) {
		t.Errorf("case sensitivity mismatch: lower=%d upper=%d", len(lower), len(upper))
	}
}

func TestMatchTemplate_ExactPrefix(t *testing.T) {
	name := schedule.MatchTemplate("hour")
	if name != "hourly" {
		t.Errorf("expected 'hourly', got %q", name)
	}
}

func TestMatchTemplate_NoMatch(t *testing.T) {
	name := schedule.MatchTemplate("zzz")
	if name != "" {
		t.Errorf("expected empty string for no match, got %q", name)
	}
}

func TestMatchTemplate_EmptyInput(t *testing.T) {
	name := schedule.MatchTemplate("")
	if name != "" {
		t.Errorf("expected empty string for empty input, got %q", name)
	}
}

func TestMatchTemplate_AmbiguousPrefix(t *testing.T) {
	// "weekly" prefix matches weekly-sunday and weekly-monday; should return one deterministically
	name := schedule.MatchTemplate("weekly")
	if name == "" {
		t.Error("expected a match for 'weekly' prefix")
	}
}
