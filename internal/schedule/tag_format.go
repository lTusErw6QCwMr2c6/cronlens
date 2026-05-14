package schedule

import (
	"fmt"
	"strings"
	"time"
)

// FormatTaggedSchedule returns a single-line summary for a tagged schedule,
// showing the tag, expression, and the next occurrence from now.
func FormatTaggedSchedule(ts TaggedSchedule) string {
	next, err := NextOccurrence(ts.Expr, time.Now())
	if err != nil {
		return fmt.Sprintf("[%s] %s  (invalid expression)", ts.Tag, ts.Expr)
	}
	until := time.Until(next)
	return fmt.Sprintf("[%s] %s  next: %s (%s)",
		ts.Tag,
		ts.Expr,
		FormatTime(next),
		FormatDuration(until),
	)
}

// FormatTagRegistry renders all entries in a registry as a formatted block.
func FormatTagRegistry(r *TagRegistry) string {
	entries := r.List()
	if len(entries) == 0 {
		return "(no tagged schedules)"
	}
	var sb strings.Builder
	for i, ts := range entries {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(FormatTaggedSchedule(ts))
	}
	return sb.String()
}

// TagSummaryLine returns a compact one-liner for display in a list view.
func TagSummaryLine(ts TaggedSchedule) string {
	desc := HumanReadable(ts.Expr)
	return fmt.Sprintf("%-20s  %-20s  %s", ts.Tag, ts.Expr, desc)
}
