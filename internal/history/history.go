package history

import (
	"sync"
	"time"
)

// Entry represents a single cron expression that was inspected.
type Entry struct {
	Expression string
	AddedAt    time.Time
	LastUsed   time.Time
	UseCount   int
}

// History maintains an in-memory list of recently used cron expressions.
type History struct {
	mu      sync.Mutex
	entries []*Entry
	maxSize int
}

// New creates a new History with the given maximum size.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 20
	}
	return &History{
		entries: make([]*Entry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add records a cron expression in the history.
// If the expression already exists, its usage count and last-used time are updated.
func (h *History) Add(expr string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	for _, e := range h.entries {
		if e.Expression == expr {
			e.LastUsed = now
			e.UseCount++
			return
		}
	}

	if len(h.entries) >= h.maxSize {
		// Remove the oldest entry (first in slice).
		h.entries = h.entries[1:]
	}

	h.entries = append(h.entries, &Entry{
		Expression: expr,
		AddedAt:    now,
		LastUsed:   now,
		UseCount:   1,
	})
}

// List returns a copy of all history entries, most recently used first.
func (h *History) List() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]Entry, len(h.entries))
	for i, e := range h.entries {
		out[len(h.entries)-1-i] = *e
	}
	return out
}

// Clear removes all entries from the history.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}

// Len returns the current number of entries.
func (h *History) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.entries)
}
