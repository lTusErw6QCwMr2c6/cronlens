package schedule

import (
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

// FieldDescription holds a human-readable description of a single cron field.
type FieldDescription struct {
	Field string
	Value string
	Meaning string
}

// Describe returns a human-readable breakdown of a cron expression.
func Describe(expr string) ([]FieldDescription, error) {
	_, err := cron.ParseStandard(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	fields := []string{"Minute", "Hour", "Day of Month", "Month", "Day of Week"}
	descriptions := make([]FieldDescription, 5)

	for i, part := range parts {
		descriptions[i] = FieldDescription{
			Field:   fields[i],
			Value:   part,
			Meaning: describeField(fields[i], part),
		}
	}

	return descriptions, nil
}

// HumanReadable returns a single-line natural language description of the expression.
func HumanReadable(expr string) (string, error) {
	descs, err := Describe(expr)
	if err != nil {
		return "", err
	}

	parts := make([]string, 0, len(descs))
	for _, d := range descs {
		if d.Value != "*" {
			parts = append(parts, fmt.Sprintf("%s: %s", d.Field, d.Meaning))
		}
	}

	if len(parts) == 0 {
		return "Every minute", nil
	}
	return strings.Join(parts, ", "), nil
}

func describeField(field, value string) string {
	if value == "*" {
		return fmt.Sprintf("every %s", strings.ToLower(field))
	}
	if strings.HasPrefix(value, "*/") {
		step := strings.TrimPrefix(value, "*/")
		return fmt.Sprintf("every %s %s(s)", step, strings.ToLower(field))
	}
	if strings.Contains(value, "-") {
		parts := strings.SplitN(value, "-", 2)
		return fmt.Sprintf("%s through %s", parts[0], parts[1])
	}
	if strings.Contains(value, ",") {
		return fmt.Sprintf("at %s", value)
	}
	return fmt.Sprintf("at %s", value)
}
