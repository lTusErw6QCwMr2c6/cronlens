package schedule

import (
	"fmt"
	"strings"
	"time"
)

// HumanizeNext returns a human-friendly string describing when a cron expression
// will next fire relative to now, e.g. "in 5 minutes", "in 2 hours", "tomorrow at 14:30".
func HumanizeNext(expr string, from time.Time) (string, error) {
	next, err := NextOccurrence(expr, from)
	if err != nil {
		return "", err
	}
	return humanizeTime(next, from), nil
}

// HumanizeNextN returns human-friendly strings for the next n occurrences.
func HumanizeNextN(expr string, from time.Time, n int) ([]string, error) {
	occurrences, err := NextOccurrences(expr, from, n)
	if err != nil {
		return nil, err
	}
	results := make([]string, len(occurrences))
	for i, t := range occurrences {
		results[i] = humanizeTime(t, from)
	}
	return results, nil
}

func humanizeTime(t, from time.Time) string {
	dur := t.Sub(from)
	if dur < 0 {
		return "in the past"
	}

	seconds := int(dur.Seconds())
	minutes := int(dur.Minutes())
	hours := int(dur.Hours())
	days := hours / 24

	switch {
	case seconds < 60:
		if seconds == 1 {
			return "in 1 second"
		}
		return fmt.Sprintf("in %d seconds", seconds)
	case minutes < 60:
		if minutes == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", minutes)
	case hours < 24:
		if hours == 1 {
			return fmt.Sprintf("in 1 hour at %s", t.Format("15:04"))
		}
		return fmt.Sprintf("in %d hours at %s", hours, t.Format("15:04"))
	case days == 1:
		return fmt.Sprintf("tomorrow at %s", t.Format("15:04"))
	case days < 7:
		return fmt.Sprintf("%s at %s", strings.ToLower(t.Weekday().String()), t.Format("15:04"))
	default:
		return fmt.Sprintf("%s at %s", t.Format("Jan 2"), t.Format("15:04"))
	}
}
