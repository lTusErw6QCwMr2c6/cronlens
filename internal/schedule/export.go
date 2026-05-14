package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExportFormat defines the supported export formats.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// ExportEntry represents a single scheduled occurrence for export.
type ExportEntry struct {
	Expression string    `json:"expression"`
	Occurrence time.Time `json:"occurrence"`
	Formatted  string    `json:"formatted"`
	Relative   string    `json:"relative"`
}

// ExportResult holds the full export payload.
type ExportResult struct {
	Expression  string        `json:"expression"`
	GeneratedAt time.Time     `json:"generated_at"`
	Entries     []ExportEntry `json:"entries"`
}

// Export generates the next n occurrences of expr starting from now
// and serialises them in the requested format.
func Export(expr string, n int, format ExportFormat) (string, error) {
	times, err := NextOccurrences(expr, n)
	if err != nil {
		return "", fmt.Errorf("export: %w", err)
	}

	now := time.Now()
	entries := make([]ExportEntry, len(times))
	for i, t := range times {
		entries[i] = ExportEntry{
			Expression: expr,
			Occurrence: t,
			Formatted:  FormatTime(t),
			Relative:   FormatDuration(t.Sub(now)),
		}
	}

	result := ExportResult{
		Expression:  expr,
		GeneratedAt: now,
		Entries:     entries,
	}

	switch format {
	case FormatJSON:
		return exportJSON(result)
	case FormatCSV:
		return exportCSV(result)
	case FormatText:
		return exportText(result)
	default:
		return "", fmt.Errorf("export: unsupported format %q", format)
	}
}

func exportJSON(r ExportResult) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func exportCSV(r ExportResult) (string, error) {
	var sb strings.Builder
	sb.WriteString("expression,occurrence,relative\n")
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "%s,%s,%s\n", e.Expression, e.Formatted, e.Relative)
	}
	return sb.String(), nil
}

func exportText(r ExportResult) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Schedule: %s\n", r.Expression)
	fmt.Fprintf(&sb, "Generated: %s\n\n", FormatTime(r.GeneratedAt))
	for i, e := range r.Entries {
		fmt.Fprintf(&sb, "  %2d. %s  (%s)\n", i+1, e.Formatted, e.Relative)
	}
	return sb.String(), nil
}
