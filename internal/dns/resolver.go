package dns

import (
	"context"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	resolveTimeout = 3 * time.Second
)

type LookupResult struct {
	Latency time.Duration
	OK      bool
	Error   error
}

func ResolveDomain(ctx context.Context, domain, dnsServer string) LookupResult {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(dialCtx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: resolveTimeout}
			return d.DialContext(dialCtx, "udp", dnsServer+":53")
		},
	}

	start := time.Now()
	_, err := r.LookupHost(ctx, domain)
	elapsed := time.Since(start)

	if err != nil {
		return LookupResult{Latency: elapsed, OK: false, Error: err}
	}
	return LookupResult{Latency: elapsed, OK: true}
}

func BatchResolve(ctx context.Context, domain, dnsServer string, count int, workers int) []LookupResult {
	results := make([]LookupResult, count)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	for i := range count {
		select {
		case <-ctx.Done():
			wg.Wait()
			return results
		default:
		}

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()
			r := ResolveDomain(ctx, domain, dnsServer)
			mu.Lock()
			results[idx] = r
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return results
}

type DomainStats struct {
	Domain      string
	Latencies   []time.Duration
	AvgLatency  time.Duration
	MinLatency  time.Duration
	MaxLatency  time.Duration
	SuccessRate float64
}

func ComputeStats(domain string, results []LookupResult) DomainStats {
	var lats []time.Duration
	success := 0

	for _, r := range results {
		if r.OK {
			lats = append(lats, r.Latency)
			success++
		}
	}

	rate := (float64(success) / float64(len(results))) * 100

	if len(lats) == 0 {
		return DomainStats{
			Domain:      domain,
			SuccessRate: rate,
		}
	}

	sort.Slice(lats, func(i, j int) bool {
		return lats[i] < lats[j]
	})

	var totalNs int64
	for _, l := range lats {
		totalNs += l.Nanoseconds()
	}

	avg := time.Duration(totalNs / int64(len(lats)))
	min := lats[0]
	max := lats[len(lats)-1]

	return DomainStats{
		Domain:      domain,
		Latencies:   lats,
		AvgLatency:  avg,
		MinLatency:  min,
		MaxLatency:  max,
		SuccessRate: rate,
	}
}
