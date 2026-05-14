package schedule

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidationError holds a field name and a human-readable reason.
type ValidationError struct {
	Field  string
	Reason string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid cron field %q: %s", e.Field, e.Reason)
}

// ValidationResult contains all errors found in a cron expression.
type ValidationResult struct {
	Expression string
	Errors     []ValidationError
}

// Valid returns true when no validation errors were found.
func (r ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

// fieldSpec describes the allowed range for a single cron field.
type fieldSpec struct {
	name string
	min  int
	max  int
}

var fieldSpecs = []fieldSpec{
	{"minute", 0, 59},
	{"hour", 0, 23},
	{"day-of-month", 1, 31},
	{"month", 1, 12},
	{"day-of-week", 0, 6},
}

// Validate checks every field of a standard 5-part cron expression and
// returns a ValidationResult describing any problems found.
func Validate(expr string) ValidationResult {
	result := ValidationResult{Expression: expr}
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		result.Errors = append(result.Errors, ValidationError{
			Field:  "expression",
			Reason: fmt.Sprintf("expected 5 fields, got %d", len(parts)),
		})
		return result
	}
	for i, spec := range fieldSpecs {
		if errs := validateField(parts[i], spec); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
		}
	}
	return result
}

func validateField(token string, spec fieldSpec) []ValidationError {
	if token == "*" {
		return nil
	}
	// step syntax */n or base/n
	if strings.Contains(token, "/") {
		parts := strings.SplitN(token, "/", 2)
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return []ValidationError{{spec.name, fmt.Sprintf("invalid step value %q", parts[1])}}
		}
		if parts[0] != "*" {
			return validateField(parts[0], spec)
		}
		return nil
	}
	// list syntax
	if strings.Contains(token, ",") {
		var errs []ValidationError
		for _, item := range strings.Split(token, ",") {
			errs = append(errs, validateField(item, spec)...)
		}
		return errs
	}
	// range syntax
	if strings.Contains(token, "-") {
		parts := strings.SplitN(token, "-", 2)
		lo, err1 := strconv.Atoi(parts[0])
		hi, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return []ValidationError{{spec.name, fmt.Sprintf("non-numeric range %q", token)}}
		}
		if lo > hi {
			return []ValidationError{{spec.name, fmt.Sprintf("range start %d > end %d", lo, hi)}}
		}
		if lo < spec.min || hi > spec.max {
			return []ValidationError{{spec.name, fmt.Sprintf("range %d-%d out of bounds [%d,%d]", lo, hi, spec.min, spec.max)}}
		}
		return nil
	}
	// plain integer
	v, err := strconv.Atoi(token)
	if err != nil {
		return []ValidationError{{spec.name, fmt.Sprintf("non-numeric value %q", token)}}
	}
	if v < spec.min || v > spec.max {
		return []ValidationError{{spec.name, fmt.Sprintf("value %d out of bounds [%d,%d]", v, spec.min, spec.max)}}
	}
	return nil
}
