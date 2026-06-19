package dns

import (
	"context"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"resolvebench/internal/models"
)

const (
	lookupsPerDomain  = 10
	concurrentWorkers = 5
	latencyWeight     = 30.0
	reliabilityWeight = 60.0
	consistencyWeight = 10.0
)

type progressFn func(current, total int, label string)

func RunBenchmark(ctx context.Context, prog progressFn) []models.ProviderResult {
	results := make([]models.ProviderResult, len(Providers))
	var mu sync.Mutex
	var wg sync.WaitGroup
	totalTasks := len(Providers) * len(BenchDomains)
	var taskCount atomic.Int64

	for pi, provider := range Providers {
		select {
		case <-ctx.Done():
			wg.Wait()
			return results
		default:
		}

		wg.Add(1)
		go func(idx int, p Provider) {
			defer wg.Done()
			domains := make([]models.DomainResult, len(BenchDomains))

			for di, domain := range BenchDomains {
				lookups := BatchResolve(ctx, domain, p.PrimaryDNS, lookupsPerDomain, concurrentWorkers)
				stats := ComputeStats(domain, lookups)

				domains[di] = models.DomainResult{
					Domain:      stats.Domain,
					Latencies:   stats.Latencies,
					AvgLatency:  stats.AvgLatency,
					MinLatency:  stats.MinLatency,
					MaxLatency:  stats.MaxLatency,
					SuccessRate: stats.SuccessRate,
				}

				count := taskCount.Add(1)
				if prog != nil {
					prog(int(count), totalTasks, p.Name)
				}
			}

			var totalAvg int64
			var totalMin int64
			var totalMax int64
			var totalRate float64
			for _, d := range domains {
				totalAvg += d.AvgLatency.Nanoseconds()
				totalMin += d.MinLatency.Nanoseconds()
				totalMax += d.MaxLatency.Nanoseconds()
				totalRate += d.SuccessRate
			}

			providerResult := models.ProviderResult{
				Name:         p.Name,
				PrimaryDNS:   p.PrimaryDNS,
				SecondaryDNS: p.SecondaryDNS,
				UseCase:      p.UseCase,
				Domains:      domains,
				AvgLatency:   time.Duration(totalAvg / int64(len(domains))),
				MinLatency:   time.Duration(totalMin / int64(len(domains))),
				MaxLatency:   time.Duration(totalMax / int64(len(domains))),
				OverallRate:  totalRate / float64(len(domains)),
			}

			mu.Lock()
			results[idx] = providerResult
			mu.Unlock()
		}(pi, provider)
	}
	wg.Wait()

	scoreResults(results)
	sortResults(results)
	return results
}

func scoreResults(results []models.ProviderResult) {
	if len(results) == 0 {
		return
	}

	var maxLatencyNs float64
	var maxAvgNs float64
	for _, r := range results {
		if r.MaxLatency.Nanoseconds() > int64(maxLatencyNs) {
			maxLatencyNs = float64(r.MaxLatency.Nanoseconds())
		}
		if r.AvgLatency.Nanoseconds() > int64(maxAvgNs) {
			maxAvgNs = float64(r.AvgLatency.Nanoseconds())
		}
	}

	for i := range results {
		r := &results[i]
		latencyScore := 0.0
		if maxAvgNs > 0 {
			latencyScore = (1 - float64(r.AvgLatency.Nanoseconds())/maxAvgNs) * latencyWeight
		}
		reliabilityScore := (r.OverallRate / 100) * reliabilityWeight
		consistencyScore := 0.0
		if maxLatencyNs > 0 {
			variation := float64(r.MaxLatency.Nanoseconds()-r.MinLatency.Nanoseconds()) / maxLatencyNs
			consistencyScore = (1 - variation) * consistencyWeight
		}
		r.Score = math.Round(latencyScore+reliabilityScore+consistencyScore*100) / 100
	}
}

func sortResults(results []models.ProviderResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
}
