package schedule

import (
	"fmt"
	"time"
)

// SimulateResult holds the result of a cron simulation run.
type SimulateResult struct {
	Expr        string
	From        time.Time
	To          time.Time
	Occurrences []time.Time
	Count       int
}

// Simulate returns all occurrences of a cron expression within a time window.
// It returns at most maxResults occurrences.
func Simulate(expr string, from, to time.Time, maxResults int) (*SimulateResult, error) {
	if maxResults <= 0 {
		maxResults = 100
	}

	sched, err := Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("simulate: invalid expression %q: %w", expr, err)
	}

	if !to.After(from) {
		return nil, fmt.Errorf("simulate: 'to' must be after 'from'")
	}

	var occurrences []time.Time
	current := from

	for len(occurrences) < maxResults {
		next := sched.Next(current)
		if next.IsZero() || next.After(to) {
			break
		}
		occurrences = append(occurrences, next)
		current = next
	}

	return &SimulateResult{
		Expr:        expr,
		From:        from,
		To:          to,
		Occurrences: occurrences,
		Count:       len(occurrences),
	}, nil
}

// SimulateDuration is a convenience wrapper that simulates over a duration from a start time.
func SimulateDuration(expr string, from time.Time, d time.Duration, maxResults int) (*SimulateResult, error) {
	if d <= 0 {
		return nil, fmt.Errorf("simulate: duration must be positive")
	}
	return Simulate(expr, from, from.Add(d), maxResults)
}
