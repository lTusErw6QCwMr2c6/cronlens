package schedule

import (
	"fmt"
	"strings"
)

// FormatConflictResult returns a human-readable multi-line string describing
// the conflict check result.
func FormatConflictResult(r ConflictResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Conflict check: [%s] vs [%s]\n", r.ExprA, r.ExprB))
	if !r.HasConflict {
		sb.WriteString("  No conflicts found in the given window.\n")
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("  %d conflict(s) found:\n", r.Count))
	for i, t := range r.Conflicts {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, FormatTime(t)))
	}
	return sb.String()
}

// ConflictSummaryLine returns a single-line summary of the conflict result.
func ConflictSummaryLine(r ConflictResult) string {
	if !r.HasConflict {
		return fmt.Sprintf("[%s] and [%s]: no conflicts", r.ExprA, r.ExprB)
	}
	return fmt.Sprintf("[%s] and [%s]: %d conflict(s), first at %s",
		r.ExprA, r.ExprB, r.Count, FormatTime(r.Conflicts[0]))
}
