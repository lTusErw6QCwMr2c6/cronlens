package schedule

import (
	"time"
)

// NextOccurrences returns the next n execution times for a cron expression
// starting from the given time t.
func NextOccurrences(expr string, t time.Time, n int) ([]time.Time, error) {
	s, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	return s.NextN(t, n), nil
}

// NextOccurrence returns the single next execution time for a cron expression
// starting from the given time t.
func NextOccurrence(expr string, t time.Time) (time.Time, error) {
	s, err := Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return s.Next(t), nil
}

// TimeUntilNext returns the duration until the next execution of the given
// cron expression relative to now.
func TimeUntilNext(expr string) (time.Duration, error) {
	now := time.Now()
	next, err := NextOccurrence(expr, now)
	if err != nil {
		return 0, err
	}
	return next.Sub(now), nil
}

// OccurrenceWindow returns all executions of a cron expression within
// the time window [from, to).
func OccurrenceWindow(expr string, from, to time.Time) ([]time.Time, error) {
	s, err := Parse(expr)
	if err != nil {
		return nil, err
	}

	var results []time.Time
	current := from
	for {
		next := s.Next(current)
		if next.IsZero() || !next.Before(to) {
			break
		}
		results = append(results, next)
		current = next
	}
	return results, nil
}
