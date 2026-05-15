package schedule

import (
	"fmt"
	"strings"
)

// FormatRetryResult returns a human-readable multi-line string describing
// all retry windows in the result.
func FormatRetryResult(r RetryResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Schedule : %s\n", r.Expr))
	sb.WriteString(fmt.Sprintf("Missed   : %s\n", FormatTime(r.Missed)))
	sb.WriteString(fmt.Sprintf("Policy   : %d attempts, %s backoff\n",
		r.Policy.MaxAttempts, FormatDuration(r.Policy.Backoff)))
	sb.WriteString("Retries  :\n")
	for _, w := range r.Windows {
		sb.WriteString(fmt.Sprintf("  [%d] %s (+%s)\n",
			w.Attempt,
			FormatTime(w.ScheduledAt),
			FormatDuration(w.Delay),
		))
	}
	return sb.String()
}

// RetrySummaryLine returns a single-line summary of a retry result.
func RetrySummaryLine(r RetryResult) string {
	if len(r.Windows) == 0 {
		return fmt.Sprintf("%s: no retries scheduled", r.Expr)
	}
	first := r.Windows[0]
	last := r.Windows[len(r.Windows)-1]
	return fmt.Sprintf("%s: %d retries from %s to %s",
		r.Expr,
		len(r.Windows),
		FormatTime(first.ScheduledAt),
		FormatTime(last.ScheduledAt),
	)
}
