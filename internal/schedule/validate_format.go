package schedule

import (
	"fmt"
	"strings"
)

// FormatValidationResult renders a ValidationResult as a multi-line string
// suitable for display in the terminal UI.
func FormatValidationResult(r ValidationResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expression : %s\n", r.Expression))
	if r.Valid() {
		sb.WriteString("Status     : ✓ valid\n")
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("Status     : ✗ invalid (%d error(s))\n", len(r.Errors)))
	for i, e := range r.Errors {
		sb.WriteString(fmt.Sprintf("  [%d] field=%s — %s\n", i+1, e.Field, e.Reason))
	}
	return sb.String()
}

// ValidationSummaryLine returns a compact single-line summary, e.g. for
// status bars or list views.
func ValidationSummaryLine(r ValidationResult) string {
	if r.Valid() {
		return fmt.Sprintf("%s → valid", r.Expression)
	}
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Reason)
	}
	return fmt.Sprintf("%s → invalid: %s", r.Expression, strings.Join(msgs, "; "))
}
