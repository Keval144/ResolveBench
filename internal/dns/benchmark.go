package dns

import (
	"math"
	"sort"
	"sync"
	"time"

	"dnsbench/internal/models"
)

const (
	lookupsPerDomain  = 10
	concurrentWorkers = 5
)

type progressFn func(current, total int, label string)

func RunBenchmark(prog progressFn) []models.ProviderResult {
	results := make([]models.ProviderResult, len(Providers))
	var mu sync.Mutex
	var wg sync.WaitGroup
	totalTasks := len(Providers) * len(BenchDomains)
	taskCount := 0

	for pi, provider := range Providers {
		wg.Add(1)
		go func(idx int, p Provider) {
			defer wg.Done()
			domains := make([]models.DomainResult, len(BenchDomains))

			for di, domain := range BenchDomains {
				lookups := BatchResolve(domain, p.PrimaryDNS, lookupsPerDomain, concurrentWorkers)
				stats := ComputeStats(domain, lookups)

				domains[di] = models.DomainResult{
					Domain:      stats.Domain,
					Latencies:   stats.Latencies,
					AvgLatency:  stats.AvgLatency,
					MinLatency:  stats.MinLatency,
					MaxLatency:  stats.MaxLatency,
					SuccessRate: stats.SuccessRate,
				}

				mu.Lock()
				taskCount++
				if prog != nil {
					prog(taskCount, totalTasks, p.Name)
				}
				mu.Unlock()
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

func RunBenchmarkCustom(customServers []string, prog progressFn) []models.ProviderResult {
	providers := Providers
	if len(customServers) > 0 {
		providers = make([]Provider, len(customServers))
		for i, s := range customServers {
			providers[i] = Provider{
				Name:        s,
				PrimaryDNS:  s,
				SecondaryDNS: s,
				UseCase:     "Custom",
				Description: "User-specified DNS server",
			}
		}
	}

	original := Providers
	Providers = providers
	results := RunBenchmark(prog)
	Providers = original
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
			latencyScore = (1 - float64(r.AvgLatency.Nanoseconds())/maxAvgNs) * 40
		}
		reliabilityScore := (r.OverallRate / 100) * 40
		consistencyScore := 0.0
		if maxLatencyNs > 0 {
			variation := float64(r.MaxLatency.Nanoseconds()-r.MinLatency.Nanoseconds()) / maxLatencyNs
			consistencyScore = (1 - variation) * 20
		}
		r.Score = math.Round(latencyScore+reliabilityScore+consistencyScore*100) / 100
	}
}

func sortResults(results []models.ProviderResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
}
