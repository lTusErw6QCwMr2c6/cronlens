package schedule

import (
	"fmt"
	"time"
)

// ConflictResult holds the result of a conflict check between two cron expressions.
type ConflictResult struct {
	ExprA      string
	ExprB      string
	Conflicts  []time.Time
	Count      int
	HasConflict bool
}

// FindConflicts checks whether two cron expressions fire at the same time
// within a given window, returning up to maxResults overlapping times.
func FindConflicts(exprA, exprB string, from time.Time, window time.Duration, maxResults int) (ConflictResult, error) {
	if maxResults <= 0 {
		maxResults = 10
	}

	schedA, err := Parse(exprA)
	if err != nil {
		return ConflictResult{}, fmt.Errorf("expression A: %w", err)
	}

	schedB, err := Parse(exprB)
	if err != nil {
		return ConflictResult{}, fmt.Errorf("expression B: %w", err)
	}

	to := from.Add(window)
	var conflicts []time.Time

	t := from
	for len(conflicts) < maxResults {
		nextA := schedA.Next(t)
		if nextA.IsZero() || nextA.After(to) {
			break
		}
		nextB := schedB.Next(nextA.Add(-time.Second))
		if !nextB.IsZero() && nextB.Equal(nextA) {
			conflicts = append(conflicts, nextA)
		}
		t = nextA
	}

	return ConflictResult{
		ExprA:       exprA,
		ExprB:       exprB,
		Conflicts:   conflicts,
		Count:       len(conflicts),
		HasConflict: len(conflicts) > 0,
	}, nil
}
