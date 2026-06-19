package network

import (
	"net"
	"time"
	"resolvebench/internal/models"
	"resolvebench/internal/dns"
)

func RunLatencyTests(prog func(current, total int, label string)) []models.NetworkResult {
	targets := dns.NetworkTargets
	results := make([]models.NetworkResult, len(targets))

	for i, target := range targets {
		result := pingTarget(target)
		results[i] = result
		if prog != nil {
			prog(i+1, len(targets), target)
		}
	}

	return results
}

func pingTarget(target string) models.NetworkResult {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target+":80", 5*time.Second)
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
