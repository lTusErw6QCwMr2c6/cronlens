package schedule

import (
	"time"

	"github.com/robfig/cron/v3"
)

// OverlapResult holds information about two schedules that share execution times.
type OverlapResult struct {
	ExprA      string
	ExprB      string
	Overlaps   []time.Time
	TotalFound int
}

// FindOverlaps returns times within the given window where both schedules fire.
// It checks up to maxResults overlapping occurrences.
func FindOverlaps(exprA, exprB string, from time.Time, window time.Duration, maxResults int) (*OverlapResult, error) {
	schedA, err := cron.ParseStandard(exprA)
	if err != nil {
		return nil, &ParseError{Expr: exprA, Err: err}
	}
	schedB, err := cron.ParseStandard(exprB)
	if err != nil {
		return nil, &ParseError{Expr: exprB, Err: err}
	}

	end := from.Add(window)
	result := &OverlapResult{
		ExprA:    exprA,
		ExprB:    exprB,
		Overlaps: []time.Time{},
	}

	// Collect all times for A within window
	timesA := map[int64]bool{}
	t := from
	for t.Before(end) {
		next := schedA.Next(t)
		if next.IsZero() || !next.Before(end) {
			break
		}
		timesA[next.Unix()] = true
		t = next
	}

	// Walk schedule B and check for matches
	t = from
	for t.Before(end) && (maxResults <= 0 || result.TotalFound < maxResults) {
		next := schedB.Next(t)
		if next.IsZero() || !next.Before(end) {
			break
		}
		if timesA[next.Unix()] {
			result.Overlaps = append(result.Overlaps, next)
			result.TotalFound++
		}
		t = next
	}

	return result, nil
}

// ParseError is returned when a cron expression cannot be parsed.
type ParseError struct {
	Expr string
	Err  error
}

func (e *ParseError) Error() string {
	return "invalid cron expression " + e.Expr + ": " + e.Err.Error()
}
