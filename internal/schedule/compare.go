package schedule

import (
	"time"
)

// CompareResult holds the comparison between two cron schedules.
type CompareResult struct {
	ExprA      string
	ExprB      string
	NextA      []time.Time
	NextB      []time.Time
	Overlaps   []time.Time
	OnlyInA    []time.Time
	OnlyInB    []time.Time
}

// Compare returns the next n occurrences for two cron expressions and
// identifies overlapping execution times within the given window.
func Compare(exprA, exprB string, from time.Time, n int) (*CompareResult, error) {
	nextA, err := NextOccurrences(exprA, from, n)
	if err != nil {
		return nil, err
	}

	nextB, err := NextOccurrences(exprB, from, n)
	if err != nil {
		return nil, err
	}

	setA := make(map[time.Time]struct{}, len(nextA))
	for _, t := range nextA {
		setA[t] = struct{}{}
	}

	setB := make(map[time.Time]struct{}, len(nextB))
	for _, t := range nextB {
		setB[t] = struct{}{}
	}

	var overlaps, onlyInA, onlyInB []time.Time

	for _, t := range nextA {
		if _, ok := setB[t]; ok {
			overlaps = append(overlaps, t)
		} else {
			onlyInA = append(onlyInA, t)
		}
	}

	for _, t := range nextB {
		if _, ok := setA[t]; !ok {
			onlyInB = append(onlyInB, t)
		}
	}

	return &CompareResult{
		ExprA:    exprA,
		ExprB:    exprB,
		NextA:    nextA,
		NextB:    nextB,
		Overlaps: overlaps,
		OnlyInA:  onlyInA,
		OnlyInB:  onlyInB,
	}, nil
}
