package schedule

import (
	"fmt"
	"strings"
	"time"
)

// FormatOverlapResult returns a human-readable summary of an OverlapResult.
func FormatOverlapResult(r *OverlapResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Overlap analysis: [%s] vs [%s]\n", r.ExprA, r.ExprB))

	if r.TotalFound == 0 {
		sb.WriteString("  No overlapping executions found in the given window.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("  Found %d overlapping execution(s):\n", r.TotalFound))
	for i, t := range r.Overlaps {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, FormatTime(t)))
	}
	return sb.String()
}

// OverlapSummaryLine returns a compact single-line summary.
func OverlapSummaryLine(r *OverlapResult, now time.Time) string {
	if r.TotalFound == 0 {
		return fmt.Sprintf("%s ∩ %s → no overlaps", r.ExprA, r.ExprB)
	}
	next := r.Overlaps[0]
	until := next.Sub(now)
	if until < 0 {
		until = 0
	}
	return fmt.Sprintf("%s ∩ %s → %d overlap(s), next in %s",
		r.ExprA, r.ExprB, r.TotalFound, FormatDuration(until))
}
