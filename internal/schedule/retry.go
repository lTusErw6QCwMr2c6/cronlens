package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// RetryPolicy defines how missed or failed executions should be retried.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

// RetryWindow represents a single retry attempt with its scheduled time.
type RetryWindow struct {
	Attempt     int
	ScheduledAt time.Time
	Delay       time.Duration
}

// RetryResult holds the outcome of computing retry windows for a schedule.
type RetryResult struct {
	Expr    string
	Missed  time.Time
	Policy  RetryPolicy
	Windows []RetryWindow
}

// ComputeRetryWindows returns retry attempt times for a missed execution
// starting from the missed time, using the given retry policy and the
// cron schedule to anchor the next valid slot after each backoff.
func ComputeRetryWindows(expr string, missed time.Time, policy RetryPolicy) (RetryResult, error) {
	if policy.MaxAttempts <= 0 {
		return RetryResult{}, fmt.Errorf("MaxAttempts must be greater than zero")
	}
	if policy.Backoff < 0 {
		return RetryResult{}, fmt.Errorf("Backoff must be non-negative")
	}

	_, err := cron.ParseStandard(expr)
	if err != nil {
		return RetryResult{}, fmt.Errorf("invalid expression %q: %w", expr, err)
	}

	windows := make([]RetryWindow, 0, policy.MaxAttempts)
	current := missed
	for i := 1; i <= policy.MaxAttempts; i++ {
		delay := time.Duration(i) * policy.Backoff
		scheduled := current.Add(delay)
		windows = append(windows, RetryWindow{
			Attempt:     i,
			ScheduledAt: scheduled,
			Delay:       delay,
		})
		current = scheduled
	}

	return RetryResult{
		Expr:    expr,
		Missed:  missed,
		Policy:  policy,
		Windows: windows,
	}, nil
}
