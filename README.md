# ResolveBench

A cross-platform TUI application that benchmarks different DNS providers to help you find the fastest, most reliable DNS resolver for your network.

## Features

- **TUI Interface**: Beautiful [Bubble Tea](https://github.com/charmbracelet/bubbletea) terminal UI with real-time progress
- **Multi-Provider Testing**: Benchmarks Cloudflare, Google DNS, Quad9, AdGuard DNS, and OpenDNS
- **Multi-Site Coverage**: Resolves 10 popular domains per provider
- **Direct DNS Queries**: Bypasses OS resolver cache, queries each provider's DNS server directly over UDP
- **Scoring System**: Weighted score (60% reliability, 30% latency, 10% consistency)
- **Raw Metrics**: Reports min/avg/max latency and success rate alongside the composite score
- **Smart Recommendations**: Optimal primary/secondary DNS selection based on DNS resolution performance
- **Network Latency Tests** *(informational)*: TCP ping to major DNS servers and endpoints
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Requirements

- Go 1.22+

## Installation

```bash
npm i -g resolvebench
```

Or manually (requires Go 1.22+):

```bash
git clone https://github.com/keval144/resolvebench.git
cd resolvebench
make build
```

## Usage

```bash
# Run the benchmark TUI
resolvebench

# List available DNS providers
resolvebench list-dns

# Show help
resolvebench -h
```
## DNS Providers Tested

| Provider | Primary | Secondary | Use Case |
|---|---|---|---|
| Cloudflare | 1.1.1.1 | 1.0.0.1 | Speed and Privacy |
| Google DNS | 8.8.8.8 | 8.8.4.4 | Reliability |
| Quad9 | 9.9.9.9 | 149.112.112.112 | Security |
| AdGuard DNS | 94.140.14.14 | 94.140.15.15 | Ad Blocking |
| OpenDNS | 208.67.222.222 | 208.67.220.220 | Content Filtering |

## Test Sites

github.com, google.com, reddit.com, wikipedia.org, cloudflare.com,
microsoft.com, amazon.com, stackoverflow.com, youtube.com, openai.com

## Methodology

ResolveBench measures DNS **resolution performance** — not routing quality or download speed.

- Each provider's DNS server is queried directly via UDP/53 using Go's `net.Resolver` with a custom dialer, bypassing the OS stub resolver and local cache
- 10 popular domains are resolved per provider, with 10 lookups per domain (5 concurrent workers)
- For each provider we report: min/avg/max latency, success rate, and a composite score
- Raw metrics are provided alongside the score so you can make your own judgment

## Scoring

The composite score is a quick comparison aid. Raw min/avg/max latency and success rate are displayed alongside for your own interpretation.

- **Reliability (60%)**: Percentage of successful resolutions — highest weight for production use
- **Latency (30%)**: Normalized average response time
- **Consistency (10%)**: Lower variance between min/max latency

## Project Structure

```
├── cli.js                     # NPM wrapper to launch Go binary
├── main.go                    # Entry point
├── cmd/
│   ├── root.go                # Cobra root command (launches TUI)
│   └── list-dns.go            # List DNS providers subcommand
├── internal/
│   ├── dns/
│   │   ├── providers.go       # DNS provider definitions & benchmark domains
│   │   ├── resolver.go        # Direct DNS lookup via UDP/53
│   │   ├── resolver_test.go
│   │   ├── benchmark.go       # Benchmark orchestration & scoring (60/30/10)
│   │   └── benchmark_test.go
│   ├── models/
│   │   └── models.go          # Shared data types
│   ├── network/
│   │   ├── latency.go         # TCP network latency tests (informational)
│   │   └── latency_test.go
│   └── tui/
│       ├── tui.go             # Bubble Tea model & update loop
│       ├── screens.go         # TUI screen views
│       └── styles.go          # Lip Gloss style definitions
├── Makefile
├── package.json
├── go.mod / go.sum
└── README.md
```

## Limitations

- **DNS resolution ≠ routing quality**: Low DNS latency doesn't guarantee optimal CDN routing or fast downloads
- **A records only**: Tests resolve A records via `LookupHost`; does not test AAAA, MX, or TXT records
- **Popular domains**: Test sites are heavily cached, globally distributed, and CDN-optimized — results may differ for less-common domains
- **UDP/53 only**: No DoH (DNS-over-HTTPS) or DoT (DNS-over-TLS) support yet
- **Geographic bias**: Results reflect performance from your location; a provider fast in one region may differ drastically in another

## Troubleshooting

- **Terminal Issues**: Ensure you're using a modern terminal that supports ANSI colors
- **Import Errors**: Run `go mod tidy` to ensure all dependencies are downloaded
- **Firewall**: DNS benchmarks require outbound UDP on port 53; network tests use TCP port 80

## License

MIT
