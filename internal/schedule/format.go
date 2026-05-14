package schedule

import (
	"fmt"
	"strings"
	"time"
)

// FormatDuration returns a human-readable string representing
// the duration until the next execution.
func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "overdue"
	}

	d = d.Round(time.Second)

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

// FormatTime formats a time.Time for display in the TUI.
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}

// SummaryLine returns a compact one-line summary for a cron entry.
func SummaryLine(expr string, next time.Time, now time.Time) string {
	until := next.Sub(now)
	return fmt.Sprintf("%-25s  next: %-28s  in: %s",
		expr,
		FormatTime(next),
		FormatDuration(until),
	)
}
