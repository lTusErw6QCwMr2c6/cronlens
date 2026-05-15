package schedule

import (
	"strings"
	"testing"
	"time"
)

var retryBase = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestComputeRetryWindows_Valid(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Backoff: 5 * time.Minute}
	res, err := ComputeRetryWindows("*/5 * * * *", retryBase, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(res.Windows))
	}
	if res.Windows[0].Attempt != 1 {
		t.Errorf("expected attempt 1, got %d", res.Windows[0].Attempt)
	}
	expectedDelay := 5 * time.Minute
	if res.Windows[0].Delay != expectedDelay {
		t.Errorf("expected delay %v, got %v", expectedDelay, res.Windows[0].Delay)
	}
	expectedTime := retryBase.Add(5 * time.Minute)
	if !res.Windows[0].ScheduledAt.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, res.Windows[0].ScheduledAt)
	}
}

func TestComputeRetryWindows_InvalidExpr(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 2, Backoff: time.Minute}
	_, err := ComputeRetryWindows("not-a-cron", retryBase, policy)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestComputeRetryWindows_ZeroAttempts(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 0, Backoff: time.Minute}
	_, err := ComputeRetryWindows("* * * * *", retryBase, policy)
	if err == nil {
		t.Fatal("expected error for zero attempts")
	}
}

func TestComputeRetryWindows_NegativeBackoff(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 1, Backoff: -time.Minute}
	_, err := ComputeRetryWindows("* * * * *", retryBase, policy)
	if err == nil {
		t.Fatal("expected error for negative backoff")
	}
}

func TestComputeRetryWindows_ResultFields(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 2, Backoff: 10 * time.Minute}
	res, err := ComputeRetryWindows("0 * * * *", retryBase, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expr != "0 * * * *" {
		t.Errorf("unexpected expr: %s", res.Expr)
	}
	if !res.Missed.Equal(retryBase) {
		t.Errorf("unexpected missed time")
	}
}

func TestFormatRetryResult(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 2, Backoff: 5 * time.Minute}
	res, _ := ComputeRetryWindows("*/5 * * * *", retryBase, policy)
	out := FormatRetryResult(res)
	if !strings.Contains(out, "Schedule") {
		t.Error("expected 'Schedule' in output")
	}
	if !strings.Contains(out, "Retries") {
		t.Error("expected 'Retries' in output")
	}
}

func TestRetrySummaryLine(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Backoff: 2 * time.Minute}
	res, _ := ComputeRetryWindows("* * * * *", retryBase, policy)
	line := RetrySummaryLine(res)
	if !strings.Contains(line, "3 retries") {
		t.Errorf("expected '3 retries' in summary, got: %s", line)
	}
}
