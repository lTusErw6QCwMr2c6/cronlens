package history

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExportedEntry is the serialisable form of a history entry.
type ExportedEntry struct {
	Expression string    `json:"expression"`
	UsedAt     time.Time `json:"used_at"`
	UseCount   int       `json:"use_count"`
}

// ExportHistory serialises the current history as JSON.
func (h *History) ExportJSON() (string, error) {
	entries := h.List()
	out := make([]ExportedEntry, len(entries))
	for i, e := range entries {
		out[i] = ExportedEntry{
			Expression: e.Expression,
			UsedAt:     e.UsedAt,
			UseCount:   e.UseCount,
		}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("history export: %w", err)
	}
	return string(b), nil
}

// ExportText returns a plain-text summary of the history.
func (h *History) ExportText() string {
	entries := h.List()
	if len(entries) == 0 {
		return "No history entries.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "History (%d entries):\n", len(entries))
	for i, e := range entries {
		fmt.Fprintf(&sb, "  %2d. %-25s  used %dx  last %s\n",
			i+1, e.Expression, e.UseCount, e.UsedAt.Format(time.RFC3339))
	}
	return sb.String()
}

// ImportJSON loads exported JSON entries into the history, adding each one.
func (h *History) ImportJSON(data string) error {
	var entries []ExportedEntry
	if err := json.Unmarshal([]byte(data), &entries); err != nil {
		return fmt.Errorf("history import: %w", err)
	}
	for _, e := range entries {
		h.Add(e.Expression)
	}
	return nil
}
