package dns

import (
	"testing"
	"time"
)

func TestComputeStats_allSuccess(t *testing.T) {
	results := []LookupResult{
		{Latency: 10 * time.Millisecond, OK: true},
		{Latency: 20 * time.Millisecond, OK: true},
		{Latency: 30 * time.Millisecond, OK: true},
	}

	stats := ComputeStats("example.com", results)

	if stats.Domain != "example.com" {
		t.Errorf("expected domain example.com, got %s", stats.Domain)
	}
	if stats.SuccessRate != 100.0 {
		t.Errorf("expected success rate 100%%, got %.1f%%", stats.SuccessRate)
	}
	if stats.AvgLatency != 20*time.Millisecond {
		t.Errorf("expected avg latency 20ms, got %v", stats.AvgLatency)
	}
	if stats.MinLatency != 10*time.Millisecond {
		t.Errorf("expected min latency 10ms, got %v", stats.MinLatency)
	}
	if stats.MaxLatency != 30*time.Millisecond {
		t.Errorf("expected max latency 30ms, got %v", stats.MaxLatency)
	}
}

func TestComputeStats_partialSuccess(t *testing.T) {
	results := []LookupResult{
		{Latency: 10 * time.Millisecond, OK: true},
		{Latency: 20 * time.Millisecond, OK: false},
		{Latency: 30 * time.Millisecond, OK: true},
	}

	stats := ComputeStats("example.com", results)

	if stats.SuccessRate != 66.66666666666666 {
		t.Errorf("expected success rate ~66.67%%, got %f", stats.SuccessRate)
	}
}

func TestComputeStats_singleResult(t *testing.T) {
	results := []LookupResult{
		{Latency: 42 * time.Millisecond, OK: true},
	}

	stats := ComputeStats("single.com", results)

	if stats.AvgLatency != 42*time.Millisecond {
		t.Errorf("expected avg 42ms, got %v", stats.AvgLatency)
	}
	if stats.MinLatency != 42*time.Millisecond {
		t.Errorf("expected min 42ms, got %v", stats.MinLatency)
	}
	if stats.MaxLatency != 42*time.Millisecond {
		t.Errorf("expected max 42ms, got %v", stats.MaxLatency)
	}
	if stats.SuccessRate != 100.0 {
		t.Errorf("expected 100%% success rate, got %f", stats.SuccessRate)
	}
}

func TestComputeStats_allFailed(t *testing.T) {
	results := []LookupResult{
		{Latency: 5 * time.Second, OK: false},
		{Latency: 5 * time.Second, OK: false},
	}

	stats := ComputeStats("fail.com", results)

	if stats.SuccessRate != 0.0 {
		t.Errorf("expected 0%% success rate, got %f", stats.SuccessRate)
	}
}
