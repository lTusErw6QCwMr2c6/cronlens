package schedule_test

import (
	"strings"
	"testing"

	"cronlens/internal/schedule"
)

func TestFindTemplate_Known(t *testing.T) {
	tmpl, ok := schedule.FindTemplate("hourly")
	if !ok {
		t.Fatal("expected to find 'hourly' template")
	}
	if tmpl.Expr != "0 * * * *" {
		t.Errorf("unexpected expr: %s", tmpl.Expr)
	}
	if tmpl.Description == "" {
		t.Error("expected non-empty description")
	}
}

func TestFindTemplate_Unknown(t *testing.T) {
	_, ok := schedule.FindTemplate("does-not-exist")
	if ok {
		t.Error("expected not to find unknown template")
	}
}

func TestListTemplates_NotEmpty(t *testing.T) {
	names := schedule.ListTemplates()
	if len(names) == 0 {
		t.Error("expected at least one template")
	}
	for _, n := range names {
		if n == "" {
			t.Error("template name should not be empty")
		}
	}
}

func TestListTemplates_ContainsCommon(t *testing.T) {
	names := schedule.ListTemplates()
	set := make(map[string]bool)
	for _, n := range names {
		set[n] = true
	}
	for _, expected := range []string{"hourly", "daily-midnight", "weekly-sunday", "monthly-first", "yearly"} {
		if !set[expected] {
			t.Errorf("expected template %q in list", expected)
		}
	}
}

func TestTemplateByExpr_Found(t *testing.T) {
	tmpl, ok := schedule.TemplateByExpr("0 * * * *")
	if !ok {
		t.Fatal("expected to find template for '0 * * * *'")
	}
	if tmpl.Name != "hourly" {
		t.Errorf("expected name 'hourly', got %q", tmpl.Name)
	}
}

func TestTemplateByExpr_NotFound(t *testing.T) {
	_, ok := schedule.TemplateByExpr("99 99 99 99 99")
	if ok {
		t.Error("expected not to find template for invalid expr")
	}
}

func TestFormatTemplate(t *testing.T) {
	tmpl := schedule.Template{Name: "hourly", Expr: "0 * * * *", Description: "Every hour at minute 0"}
	line := schedule.FormatTemplate(tmpl)
	if !strings.Contains(line, "hourly") {
		t.Error("expected name in formatted output")
	}
	if !strings.Contains(line, "0 * * * *") {
		t.Error("expected expr in formatted output")
	}
	if !strings.Contains(line, "Every hour") {
		t.Error("expected description in formatted output")
	}
}

func TestFormatTemplateList_Empty(t *testing.T) {
	out := schedule.FormatTemplateList([]schedule.Template{})
	if !strings.Contains(out, "No templates") {
		t.Error("expected 'No templates' message for empty list")
	}
}

func TestFormatTemplateList_NonEmpty(t *testing.T) {
	out := schedule.FormatTemplateList(schedule.CommonTemplates)
	if !strings.Contains(out, "NAME") {
		t.Error("expected header in table output")
	}
	if !strings.Contains(out, "hourly") {
		t.Error("expected 'hourly' in table output")
	}
}

func TestTemplateSummaryLine_Found(t *testing.T) {
	line := schedule.TemplateSummaryLine("hourly", true)
	if !strings.Contains(line, "found") {
		t.Error("expected 'found' in summary line")
	}
}

func TestTemplateSummaryLine_NotFound(t *testing.T) {
	line := schedule.TemplateSummaryLine("ghost", false)
	if !strings.Contains(line, "not found") {
		t.Error("expected 'not found' in summary line")
	}
}
