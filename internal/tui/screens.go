package tui

import (
	"fmt"
	"resolvebench/internal/dns"
	"strings"
)

var resolveBenchArt = TitleStyle.Render(
	" 888888ba                             dP                    888888ba                             dP       \n" +
		" 88    `8b                            88                    88    `8b                            88       \n" +
		"a88aaaa8P' .d8888b. .d8888b. .d8888b. 88 dP   .dP .d8888b. a88aaaa8P' .d8888b. 88d888b. .d8888b. 88d888b. \n" +
		" 88   `8b. 88ooood8 Y8ooooo. 88'  `88 88 88   d8' 88ooood8  88   `8b. 88ooood8 88'  `88 88'  `\"\" 88'  `88 \n" +
		" 88     88 88.  ...       88 88.  .88 88 88 .88'  88.  ...  88    .88 88.  ... 88    88 88.  ... 88    88 \n" +
		" dP     dP `88888P' `88888P' `88888P' dP 8888P'   `88888P'  88888888P `88888P' dP    dP `88888P' dP    dP \n" +
		"                                                                                                           \n" +
		"                                                                                                           ",
)

func welcomeView() string {
	var b strings.Builder
	b.WriteString(resolveBenchArt)
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("Benchmark & Compare DNS Providers"))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("Find the fastest, most reliable DNS for your network"))
	b.WriteString("\n\n")

	b.WriteString(HighlightStyle.Render("DNS Providers to test:"))
	b.WriteString("\n")
	for _, p := range dns.Providers {
		b.WriteString("  ")
		b.WriteString(SuccessStyle.Render("▸"))
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%-12s", p.Name))
		b.WriteString(" ")
		b.WriteString(DimStyle.Render("(" + p.UseCase + ")"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(DimStyle.Render("Benchmarking 10 domains across 5 providers"))
	b.WriteString("\n\n")
	b.WriteString(HighlightStyle.Render("Press Enter"))
	b.WriteString(DimStyle.Render("  to start  ·  q / Esc to quit"))

	return b.String()
}

func (m *Model) runningView() string {
	var b strings.Builder
	b.WriteString(resolveBenchArt)
	b.WriteString("\n")

	totalDNS := len(dns.Providers) * len(dns.BenchDomains)
	currentDNS := m.progress
	if currentDNS > totalDNS {
		currentDNS = totalDNS
	}

	if currentDNS < totalDNS {
		b.WriteString(fmt.Sprintf("Benchmarking Domains: %d/%d\n", currentDNS, totalDNS))
	} else {
		b.WriteString(SuccessStyle.Render(fmt.Sprintf("Benchmarking Domains: %d/%d", totalDNS, totalDNS)))
		b.WriteString(fmt.Sprintf("  Network Latency: %d/%d\n",
			m.progress-totalDNS,
			m.progressMax-totalDNS,
		))
	}

	b.WriteString("\n")
	b.WriteString("  ")
	b.WriteString(progressBar(m.progress, m.progressMax, 30))
	b.WriteString("  ")
	b.WriteString(DimStyle.Render(fmt.Sprintf("%d/%d", m.progress, m.progressMax)))
	b.WriteString("\n")

	if m.progressLbl != "" {
		b.WriteString("  ")
		b.WriteString(DimStyle.Render(m.progressLbl))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(DimStyle.Render("Press q or Esc to quit"))
	return b.String()
}

func progressBar(current, total, width int) string {
	if total == 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return HighlightStyle.Render(bar)
}

func (m *Model) resultsView() string {
	tfmt := "%-2s %-14s %10s %10s %10s %8s %5s"
	var b strings.Builder
	b.WriteString(resolveBenchArt)
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("DNS Benchmark Results"))
	b.WriteString("\n\n")

	b.WriteString(TableHeader.Render(fmt.Sprintf(tfmt,
		"#", "Provider", "Avg", "Min", "Max", "Success", "Score",
	)))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render(strings.Repeat("─", 65)))
	b.WriteString("\n")

	for i, r := range m.results {
		line := fmt.Sprintf(tfmt,
			fmt.Sprintf("%d.", i+1),
			r.Name,
			formatDur(r.AvgLatency),
			formatDur(r.MinLatency),
			formatDur(r.MaxLatency),
			fmt.Sprintf("%.0f%%", r.OverallRate),
			fmt.Sprintf("%.1f", r.Score),
		)
		switch {
		case r.Score >= 80:
			b.WriteString(SuccessStyle.Render(line))
		case r.Score >= 60:
			b.WriteString(line)
		default:
			b.WriteString(DimStyle.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(BorderStyle.Render(m.recommendationView()))
	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("Enter — Network Tests  ·  Space — Restart Now  ·  q — Quit"))
	return b.String()
}

func (m *Model) recommendationView() string {
	return m.renderRecommendation()
}

func (m *Model) renderRecommendation() string {
	if len(m.results) < 2 {
		return "Insufficient data for recommendations"
	}

	var b strings.Builder
	primary := m.results[0]
	secondary := m.results[1]

	b.WriteString(PrimaryStyle.Render(fmt.Sprintf("  ★ Primary DNS:   %s (%s)", primary.Name, primary.PrimaryDNS)))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("     Use Case: %s\n", DimStyle.Render(primary.UseCase)))
	b.WriteString(fmt.Sprintf("     Latency:  %s  |  Success: %.0f%%\n",
		formatDur(primary.AvgLatency), primary.OverallRate))

	b.WriteString("\n")
	b.WriteString(SecondaryStyle.Render(fmt.Sprintf("  ◆ Secondary DNS: %s (%s)", secondary.Name, secondary.SecondaryDNS)))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("     Use Case: %s\n", DimStyle.Render(secondary.UseCase)))
	b.WriteString(fmt.Sprintf("     Latency:  %s  |  Success: %.0f%%",
		formatDur(secondary.AvgLatency), secondary.OverallRate))

	return b.String()
}

func (m *Model) networkView() string {
	tfmt := "%-3s %-18s %12s %14s"
	sep := strings.Repeat("─", 50)
	var b strings.Builder
	b.WriteString(resolveBenchArt)
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("Network Latency Tests"))
	b.WriteString("\n\n")

	b.WriteString(TableHeader.Render(fmt.Sprintf(tfmt,
		"#", "Target", "Latency", "Status",
	)))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render(sep))
	b.WriteString("\n")

	for i, r := range m.networkRezs {
		rank := fmt.Sprintf("%d.", i+1)
		status := "✓ reachable"
		style := SuccessStyle
		if !r.Reachable {
			status = "✗ unreachable"
			style = ErrorStyle
		}
		latency := formatDur(r.Latency)
		if !r.Reachable {
			latency = "-"
		}
		line := fmt.Sprintf(tfmt, rank, r.Target, latency, status)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	if len(m.results) >= 2 {
		b.WriteString("\n")
		b.WriteString(BorderStyle.Render(m.renderRecommendation()))
	}

	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("Enter — Back to Results  ·  Space — Restart Now  ·  q — Quit"))
	return b.String()
}
