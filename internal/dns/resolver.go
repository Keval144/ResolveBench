package dns

import (
	"context"
	"log"
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

func ResolveDomain(domain, dnsServer string) LookupResult {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: resolveTimeout}
			return d.DialContext(ctx, "udp", dnsServer+":53")
		},
	}

	start := time.Now()
	_, err := r.LookupHost(context.Background(), domain)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("DNS lookup failed for %s via %s: %v", domain, dnsServer, err)
		return LookupResult{Latency: elapsed, OK: false, Error: err}
	}
	return LookupResult{Latency: elapsed, OK: true}
}

func BatchResolve(domain, dnsServer string, count int, workers int) []LookupResult {
	results := make([]LookupResult, count)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	for i := range count {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			r := ResolveDomain(domain, dnsServer)
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
		lats = append(lats, r.Latency)
		if r.OK {
			success++
		}
	}

	sort.Slice(lats, func(i, j int) bool {
		return lats[i] < lats[j]
	})

	var totalNs int64
	for _, l := range lats {
		totalNs += l.Nanoseconds()
	}

	denom := len(lats)
	if denom < 1 {
		denom = 1
	}
	avg := time.Duration(totalNs / int64(denom))
	min := lats[0]
	max := lats[len(lats)-1]
	rate := (float64(success) / float64(len(results))) * 100

	return DomainStats{
		Domain:      domain,
		Latencies:   lats,
		AvgLatency:  avg,
		MinLatency:  min,
		MaxLatency:  max,
		SuccessRate: rate,
	}
}
