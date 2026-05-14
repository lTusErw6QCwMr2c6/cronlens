package schedule

import (
	"fmt"
	"sort"
	"strings"
)

// FormatDiffResult returns a human-readable multi-line string
// summarizing the diff between two cron schedules.
func FormatDiffResult(d *DiffResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Diff: [%s] vs [%s]\n", d.ExprA, d.ExprB))
	sb.WriteString(fmt.Sprintf("  Shared occurrences : %d\n", len(d.InBoth)))
	sb.WriteString(fmt.Sprintf("  Only in A (%s) : %d\n", d.ExprA, len(d.OnlyInA)))
	sb.WriteString(fmt.Sprintf("  Only in B (%s) : %d\n", d.ExprB, len(d.OnlyInB)))

	if len(d.OnlyInA) > 0 {
		sort.Slice(d.OnlyInA, func(i, j int) bool { return d.OnlyInA[i].Before(d.OnlyInA[j]) })
		sb.WriteString("  [A only]\n")
		for _, t := range d.OnlyInA {
			sb.WriteString(fmt.Sprintf("    - %s\n", FormatTime(t)))
		}
	}

	if len(d.OnlyInB) > 0 {
		sort.Slice(d.OnlyInB, func(i, j int) bool { return d.OnlyInB[i].Before(d.OnlyInB[j]) })
		sb.WriteString("  [B only]\n")
		for _, t := range d.OnlyInB {
			sb.WriteString(fmt.Sprintf("    - %s\n", FormatTime(t)))
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

// DiffSummaryLine returns a single-line summary of the diff result.
func DiffSummaryLine(d *DiffResult) string {
	return fmt.Sprintf(
		"%s vs %s — shared: %d, only-A: %d, only-B: %d",
		d.ExprA, d.ExprB,
		len(d.InBoth), len(d.OnlyInA), len(d.OnlyInB),
	)
}
