package schedule

import (
	"fmt"
	"sort"
	"strings"
)

// TaggedSchedule associates a human-friendly label with a cron expression.
type TaggedSchedule struct {
	Tag  string
	Expr string
}

// TagRegistry holds a named collection of tagged schedules.
type TagRegistry struct {
	entries []TaggedSchedule
}

// NewTagRegistry creates an empty TagRegistry.
func NewTagRegistry() *TagRegistry {
	return &TagRegistry{}
}

// Add registers a new tagged schedule. Returns an error if the tag already exists
// or the expression is invalid.
func (r *TagRegistry) Add(tag, expr string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return fmt.Errorf("tag must not be empty")
	}
	for _, e := range r.entries {
		if e.Tag == tag {
			return fmt.Errorf("tag %q already exists", tag)
		}
	}
	if _, err := Parse(expr); err != nil {
		return fmt.Errorf("invalid expression for tag %q: %w", tag, err)
	}
	r.entries = append(r.entries, TaggedSchedule{Tag: tag, Expr: expr})
	return nil
}

// Remove deletes a tag from the registry. Returns an error if not found.
func (r *TagRegistry) Remove(tag string) error {
	for i, e := range r.entries {
		if e.Tag == tag {
			r.entries = append(r.entries[:i], r.entries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("tag %q not found", tag)
}

// Get returns the TaggedSchedule for the given tag.
func (r *TagRegistry) Get(tag string) (TaggedSchedule, bool) {
	for _, e := range r.entries {
		if e.Tag == tag {
			return e, true
		}
	}
	return TaggedSchedule{}, false
}

// List returns all entries sorted alphabetically by tag.
func (r *TagRegistry) List() []TaggedSchedule {
	out := make([]TaggedSchedule, len(r.entries))
	copy(out, r.entries)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Tag < out[j].Tag
	})
	return out
}

// Len returns the number of registered tags.
func (r *TagRegistry) Len() int {
	return len(r.entries)
}
