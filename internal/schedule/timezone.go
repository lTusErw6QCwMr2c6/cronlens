package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// TimezoneResult holds next occurrences of a cron expression in multiple timezones.
type TimezoneResult struct {
	Expression string
	Zones      []ZoneOccurrence
}

// ZoneOccurrence represents the next occurrence of a schedule in a specific timezone.
type ZoneOccurrence struct {
	ZoneName string
	Location *time.Location
	Next      time.Time
	Offset    string
}

// NextInTimezones returns the next occurrence of the given cron expression
// evaluated in each of the provided IANA timezone names.
func NextInTimezones(expr string, zoneNames []string, from time.Time) (*TimezoneResult, error) {
	result := &TimezoneResult{
		Expression: expr,
		Zones:      make([]ZoneOccurrence, 0, len(zoneNames)),
	}

	for _, name := range zoneNames {
		loc, err := time.LoadLocation(name)
		if err != nil {
			return nil, fmt.Errorf("unknown timezone %q: %w", name, err)
		}

		parser := cron.NewParser(
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
		)
		sched, err := parser.Parse(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid expression %q: %w", expr, err)
		}

		localFrom := from.In(loc)
		next := sched.Next(localFrom)

		_, offsetSecs := next.Zone()
		hours := offsetSecs / 3600
		minutes := (offsetSecs % 3600) / 60
		offsetStr := fmt.Sprintf("UTC%+03d:%02d", hours, abs(minutes))

		result.Zones = append(result.Zones, ZoneOccurrence{
			ZoneName: name,
			Location: loc,
			Next:      next,
			Offset:   offsetStr,
		})
	}

	return result, nil
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
