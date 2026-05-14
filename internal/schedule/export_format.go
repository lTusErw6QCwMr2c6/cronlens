package schedule

import "fmt"

// FormatExportResult returns a short human-readable summary of an export.
func FormatExportResult(expr string, count int, format ExportFormat) string {
	return fmt.Sprintf("Exported %d occurrences of %q as %s", count, expr, format)
}

// ExportSummaryLine returns a one-line description suitable for a status bar.
func ExportSummaryLine(result ExportResult) string {
	if len(result.Entries) == 0 {
		return fmt.Sprintf("No occurrences found for %q", result.Expression)
	}
	first := result.Entries[0]
	last := result.Entries[len(result.Entries)-1]
	return fmt.Sprintf(
		"%s — %d occurrences from %s to %s",
		result.Expression,
		len(result.Entries),
		FormatTime(first.Occurrence),
		FormatTime(last.Occurrence),
	)
}
