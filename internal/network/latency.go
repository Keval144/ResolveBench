package network

import (
	"context"
	"net"
	"resolvebench/internal/dns"
	"resolvebench/internal/models"
	"sort"
	"time"
)

func RunLatencyTests(ctx context.Context, prog func(current, total int, label string)) []models.NetworkResult {
	targets := dns.NetworkTargets
	results := make([]models.NetworkResult, len(targets))

	for i, target := range targets {
		select {
		case <-ctx.Done():
			return results
		default:
		}
		result := pingTarget(target)
		results[i] = result
		if prog != nil {
			prog(i+1, len(targets), target)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Reachable != results[j].Reachable {
			return results[i].Reachable
		}
		return results[i].Latency < results[j].Latency
	})

	return results
}

func pingTarget(target string) models.NetworkResult {
	port := "80"
	if net.ParseIP(target) != nil {
		port = "53"
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(target, port), 5*time.Second)
	if err != nil {
		return models.NetworkResult{
			Target:    target,
			Reachable: false,
		}
	}
	defer conn.Close()
	elapsed := time.Since(start)

	return models.NetworkResult{
		Target:    target,
		Latency:   elapsed,
		Reachable: true,
	}
}
