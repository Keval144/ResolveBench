package models

import "time"

type Screen string

const (
	WelcomeScreen Screen = "welcome"
	RunningScreen Screen = "running"
	ResultsScreen Screen = "results"
	NetworkScreen Screen = "network"
)

type DomainResult struct {
	Domain      string
	Latencies   []time.Duration
	AvgLatency  time.Duration
	MinLatency  time.Duration
	MaxLatency  time.Duration
	SuccessRate float64
}

type ProviderResult struct {
	Name         string
	PrimaryDNS   string
	SecondaryDNS string
	UseCase      string
	Domains      []DomainResult
	AvgLatency   time.Duration
	MinLatency   time.Duration
	MaxLatency   time.Duration
	OverallRate  float64
	Score        float64
}

type NetworkResult struct {
	Target    string
	Latency   time.Duration
	Reachable bool
}
