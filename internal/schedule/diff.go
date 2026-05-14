package schedule

import (
	"fmt"
	"time"
)

// DiffResult holds the comparison between two cron schedules over a time window.
type DiffResult struct {
	ExprA      string
	ExprB      string
	OnlyInA    []time.Time
	OnlyInB    []time.Time
	InBoth     []time.Time
}

// Diff computes the symmetric difference between two cron expressions
// over the next n occurrences starting from a given time.
func Diff(exprA, exprB string, from time.Time, n int) (*DiffResult, error) {
	schedA, err := Parse(exprA)
	if err != nil {
		return nil, fmt.Errorf("invalid expression A: %w", err)
	}
	schedB, err := Parse(exprB)
	if err != nil {
		return nil, fmt.Errorf("invalid expression B: %w", err)
	}

	setA := make(map[time.Time]struct{})
	setB := make(map[time.Time]struct{})

	t := from
	for i := 0; i < n; i++ {
		t = schedA.Next(t)
		setA[t] = struct{}{}
	}

	t = from
	for i := 0; i < n; i++ {
		t = schedB.Next(t)
		setB[t] = struct{}{}
	}

	result := &DiffResult{ExprA: exprA, ExprB: exprB}

	for ts := range setA {
		if _, ok := setB[ts]; ok {
			result.InBoth = append(result.InBoth, ts)
		} else {
			result.OnlyInA = append(result.OnlyInA, ts)
		}
	}
	for ts := range setB {
		if _, ok := setA[ts]; !ok {
			result.OnlyInB = append(result.OnlyInB, ts)
		}
	}

	return result, nil
}
