package schedule

import (
	"strings"
	"testing"
)

func TestValidate_ValidExpressions(t *testing.T) {
	cases := []string{
		"* * * * *",
		"0 * * * *",
		"30 6 * * 1",
		"0 0 1 1 *",
		"*/15 * * * *",
		"0-30 8-17 * * 1-5",
		"0,15,30,45 * * * *",
		"59 23 31 12 6",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			r := Validate(expr)
			if !r.Valid() {
				t.Errorf("expected valid, got errors: %v", r.Errors)
			}
		})
	}
}

func TestValidate_WrongFieldCount(t *testing.T) {
	r := Validate("* * * *")
	if r.Valid() {
		t.Fatal("expected invalid for 4-field expression")
	}
	if r.Errors[0].Field != "expression" {
		t.Errorf("expected field=expression, got %q", r.Errors[0].Field)
	}
}

func TestValidate_OutOfRangeValues(t *testing.T) {
	cases := []struct {
		expr  string
		field string
	}{
		{"60 * * * *", "minute"},
		{"* 24 * * *", "hour"},
		{"* * 0 * *", "day-of-month"},
		{"* * * 13 *", "month"},
		{"* * * * 7", "day-of-week"},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			r := Validate(tc.expr)
			if r.Valid() {
				t.Fatalf("expected invalid for %q", tc.expr)
			}
			found := false
			for _, e := range r.Errors {
				if e.Field == tc.field {
					found = true
				}
			}
			if !found {
				t.Errorf("expected error for field %q, got %v", tc.field, r.Errors)
			}
		})
	}
}

func TestValidate_InvalidStep(t *testing.T) {
	r := Validate("*/0 * * * *")
	if r.Valid() {
		t.Fatal("expected invalid for step 0")
	}
}

func TestValidate_ReversedRange(t *testing.T) {
	r := Validate("30-10 * * * *")
	if r.Valid() {
		t.Fatal("expected invalid for reversed range")
	}
}

func TestFormatValidationResult_Valid(t *testing.T) {
	r := Validate("0 9 * * 1")
	out := FormatValidationResult(r)
	if !strings.Contains(out, "✓ valid") {
		t.Errorf("expected valid marker in output, got: %s", out)
	}
}

func TestFormatValidationResult_Invalid(t *testing.T) {
	r := Validate("99 * * * *")
	out := FormatValidationResult(r)
	if !strings.Contains(out, "✗ invalid") {
		t.Errorf("expected invalid marker in output, got: %s", out)
	}
}

func TestValidationSummaryLine(t *testing.T) {
	valid := ValidationSummaryLine(Validate("* * * * *"))
	if !strings.Contains(valid, "valid") {
		t.Errorf("unexpected summary: %s", valid)
	}
	invalid := ValidationSummaryLine(Validate("99 * * * *"))
	if !strings.Contains(invalid, "invalid") {
		t.Errorf("unexpected summary: %s", invalid)
	}
}
