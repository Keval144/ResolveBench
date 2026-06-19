package dns

import (
	"testing"
	"time"

	"dnsbench/internal/models"
)

func TestScoreResults_empty(t *testing.T) {
	scoreResults(nil)
	scoreResults([]models.ProviderResult{})
}

func TestScoreResults_singleProvider(t *testing.T) {
	results := []models.ProviderResult{
		{
			Name:        "Test DNS",
			AvgLatency:  25 * time.Millisecond,
			MinLatency:  10 * time.Millisecond,
			MaxLatency:  50 * time.Millisecond,
			OverallRate: 100.0,
		},
	}
	scoreResults(results)

	if results[0].Score == 0 {
		t.Error("expected non-zero score for successful provider")
	}
}

func TestScoreResults_scoring(t *testing.T) {
	results := []models.ProviderResult{
		{
			Name:        "Fast DNS",
			AvgLatency:  10 * time.Millisecond,
			MinLatency:  5 * time.Millisecond,
			MaxLatency:  20 * time.Millisecond,
			OverallRate: 100.0,
		},
		{
			Name:        "Slow DNS",
			AvgLatency:  100 * time.Millisecond,
			MinLatency:  50 * time.Millisecond,
			MaxLatency:  200 * time.Millisecond,
			OverallRate: 80.0,
		},
	}
	scoreResults(results)

	if results[0].Score <= results[1].Score {
		t.Error("expected faster DNS to score higher")
	}
}

func TestSortResults(t *testing.T) {
	results := []models.ProviderResult{
		{Name: "A", Score: 50.0},
		{Name: "B", Score: 90.0},
		{Name: "C", Score: 70.0},
	}
	sortResults(results)

	if results[0].Name != "B" {
		t.Errorf("expected B as first (score 90), got %s (%f)", results[0].Name, results[0].Score)
	}
	if results[1].Name != "C" {
		t.Errorf("expected C as second (score 70), got %s (%f)", results[1].Name, results[1].Score)
	}
	if results[2].Name != "A" {
		t.Errorf("expected A as third (score 50), got %s (%f)", results[2].Name, results[2].Score)
	}
}

func TestSortResults_tie(t *testing.T) {
	results := []models.ProviderResult{
		{Name: "A", Score: 80.0},
		{Name: "B", Score: 80.0},
	}
	sortResults(results)
}
