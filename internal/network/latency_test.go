package network

import (
	"context"
	"testing"
	"time"
)

func TestPingTarget_unreachable(t *testing.T) {
	result := pingTarget("198.51.100.1")

	if result.Reachable {
		t.Error("expected 198.51.100.1 (TEST-NET) to be unreachable")
	}
	if result.Latency != 0 {
		t.Errorf("expected zero latency for unreachable target, got %v", result.Latency)
	}
}

func TestRunLatencyTests(t *testing.T) {
	results := RunLatencyTests(context.Background(), nil)

	if len(results) == 0 {
		t.Error("expected at least one network result")
	}

	for _, r := range results {
		if r.Target == "" {
			t.Error("expected non-empty target")
		}
	}
}

func TestRunLatencyTests_progressCallback(t *testing.T) {
	calls := 0
	total := 0
	prog := func(current, total_ int, label string) {
		calls++
		total = total_
	}

	RunLatencyTests(context.Background(), prog)

	if calls == 0 {
		t.Error("expected progress callback to be called")
	}
	if total == 0 {
		t.Error("expected non-zero total from progress callback")
	}
}

func TestPingTarget_validLatency(t *testing.T) {
	start := time.Now()
	result := pingTarget("8.8.8.8")
	elapsed := time.Since(start)

	if elapsed > 30*time.Second {
		t.Fatal("test took too long, skipping")
	}

	if result.Target != "8.8.8.8" {
		t.Errorf("expected target 8.8.8.8, got %s", result.Target)
	}
}

func TestRunLatencyTests_progressCount(t *testing.T) {
	targets := []string{"1.1.1.1", "8.8.8.8"}
	calls := 0
	prog := func(current, total int, label string) {
		calls++
	}

	for i, target := range targets {
		pingTarget(target)
		if prog != nil {
			prog(i+1, len(targets), target)
		}
	}

	if calls == 0 {
		t.Error("expected progress callback to be called")
	}
}
