package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Entry represents a parsed cron schedule with metadata.
type Entry struct {
	Expression string
	Schedule   cron.Schedule
}

// Parse validates and parses a cron expression string.
// Supports standard 5-field and optional seconds (6-field) expressions.
func Parse(expr string) (*Entry, error) {
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	sched, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}

	return &Entry{
		Expression: expr,
		Schedule:   sched,
	}, nil
}

// NextN returns the next n execution times starting from the given time.
func (e *Entry) NextN(from time.Time, n int) []time.Time {
	if n <= 0 {
		return nil
	}

	times := make([]time.Time, 0, n)
	current := from

	for i := 0; i < n; i++ {
		next := e.Schedule.Next(current)
		if next.IsZero() {
			break
		}
		times = append(times, next)
		current = next
	}

	return times
}

// NextFrom returns the next execution time after the given time.
func (e *Entry) NextFrom(from time.Time) time.Time {
	return e.Schedule.Next(from)
}
