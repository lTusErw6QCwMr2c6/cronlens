package schedule

import (
	"strings"
	"testing"
)

func TestTagRegistry_AddAndGet(t *testing.T) {
	r := NewTagRegistry()
	if err := r.Add("daily", "0 9 * * *"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts, ok := r.Get("daily")
	if !ok {
		t.Fatal("expected to find tag 'daily'")
	}
	if ts.Expr != "0 9 * * *" {
		t.Errorf("expected expr '0 9 * * *', got %q", ts.Expr)
	}
}

func TestTagRegistry_AddDuplicateTag(t *testing.T) {
	r := NewTagRegistry()
	_ = r.Add("daily", "0 9 * * *")
	err := r.Add("daily", "0 10 * * *")
	if err == nil {
		t.Fatal("expected error for duplicate tag")
	}
}

func TestTagRegistry_AddInvalidExpr(t *testing.T) {
	r := NewTagRegistry()
	err := r.Add("bad", "not-a-cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestTagRegistry_AddEmptyTag(t *testing.T) {
	r := NewTagRegistry()
	err := r.Add("", "* * * * *")
	if err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestTagRegistry_Remove(t *testing.T) {
	r := NewTagRegistry()
	_ = r.Add("daily", "0 9 * * *")
	if err := r.Remove("daily"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", r.Len())
	}
}

func TestTagRegistry_RemoveNotFound(t *testing.T) {
	r := NewTagRegistry()
	err := r.Remove("ghost")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestTagRegistry_ListSorted(t *testing.T) {
	r := NewTagRegistry()
	_ = r.Add("zebra", "0 1 * * *")
	_ = r.Add("alpha", "0 2 * * *")
	_ = r.Add("middle", "0 3 * * *")
	list := r.List()
	if list[0].Tag != "alpha" || list[1].Tag != "middle" || list[2].Tag != "zebra" {
		t.Errorf("list not sorted: %v", list)
	}
}

func TestFormatTagRegistry_Empty(t *testing.T) {
	r := NewTagRegistry()
	out := FormatTagRegistry(r)
	if out != "(no tagged schedules)" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTagRegistry_NonEmpty(t *testing.T) {
	r := NewTagRegistry()
	_ = r.Add("hourly", "0 * * * *")
	out := FormatTagRegistry(r)
	if !strings.Contains(out, "hourly") {
		t.Errorf("expected tag in output, got: %q", out)
	}
}

func TestTagSummaryLine(t *testing.T) {
	ts := TaggedSchedule{Tag: "daily", Expr: "0 9 * * *"}
	line := TagSummaryLine(ts)
	if !strings.Contains(line, "daily") || !strings.Contains(line, "0 9 * * *") {
		t.Errorf("summary line missing expected content: %q", line)
	}
}
