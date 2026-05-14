package schedule_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/cronlens/internal/schedule"
)

func TestExport_JSON(t *testing.T) {
	out, err := schedule.Export("* * * * *", 3, schedule.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result schedule.ExportResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result.Entries))
	}
	if result.Expression != "* * * * *" {
		t.Errorf("unexpected expression %q", result.Expression)
	}
}

func TestExport_CSV(t *testing.T) {
	out, err := schedule.Export("0 9 * * 1", 2, schedule.FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// header + 2 data rows
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header+2), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "expression,") {
		t.Errorf("expected CSV header, got %q", lines[0])
	}
}

func TestExport_Text(t *testing.T) {
	out, err := schedule.Export("30 6 * * *", 5, schedule.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Schedule:") {
		t.Error("text output missing 'Schedule:' header")
	}
	if !strings.Contains(out, "30 6 * * *") {
		t.Error("text output missing expression")
	}
}

func TestExport_InvalidExpr(t *testing.T) {
	_, err := schedule.Export("not-a-cron", 3, schedule.FormatJSON)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	_, err := schedule.Export("* * * * *", 3, schedule.ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExportSummaryLine(t *testing.T) {
	out, err := schedule.Export("* * * * *", 4, schedule.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result schedule.ExportResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	line := schedule.ExportSummaryLine(result)
	if !strings.Contains(line, "4 occurrences") {
		t.Errorf("summary line missing count: %q", line)
	}
}
