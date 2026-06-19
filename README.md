# DnsBench

A cross-platform TUI application that benchmarks different DNS providers to help you find the lowest-latency DNS for optimal routing.

## Features

- **TUI Interface**: Beautiful [Bubble Tea](https://github.com/charmbracelet/bubbletea) terminal UI with real-time progress
- **Multi-Provider Testing**: Benchmarks Cloudflare, Google DNS, Quad9, AdGuard DNS, and OpenDNS
- **Multi-Site Coverage**: Resolves 10 popular domains (GitHub, Google, Reddit, etc.) per provider
- **Scoring System**: Weighted score (40% latency, 40% reliability, 20% consistency)
- **Smart Recommendations**: Optimal primary/secondary DNS selection
- **Network Latency Tests**: TCP ping to major DNS servers and endpoints
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Requirements

- Go 1.22+

## Installation

```bash
git clone https://github.com/keval144/dnsbench.git
cd dnsbench
make build
```

Or manually:

```bash
go mod tidy
go build -o dnsbench
```

## Usage

### Start Benchmarking

```bash
dnsbench bench
```

### View Available DNS Servers

```bash
dnsbench list-dns
```

### Non-Interactive Mode

```bash
dnsbench bench --json
```

### Custom DNS Servers

```bash
dnsbench bench --dns 1.1.1.1,8.8.8.8
```

## UI Flow

1. **Welcome Screen**: App introduction with providers listed
2. **Run Test**: Press Enter to start benchmarking
3. **Progress Display**: Real-time progress bar with provider labels
4. **Results Dashboard**: Sorted table with Avg/Min/Max latency, success rate, score
5. **Network Tests**: TCP ping results for selected targets
6. **Recommendations**: Primary and secondary DNS suggestions

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

## Scoring

- **Latency (40%)**: Normalized average response time
- **Reliability (40%)**: Percentage of successful resolutions
- **Consistency (20%)**: Lower variance between min/max latency

## Project Structure

```
├── main.go                     # Entry point
├── internal/
│   ├── dns/
│   │   ├── providers.go        # DNS provider definitions
│   │   ├── resolver.go         # DNS lookup & stats computation
│   │   ├── resolver_test.go    # Tests for resolver
│   │   ├── benchmark.go        # Benchmark orchestration & scoring
│   │   └── benchmark_test.go   # Tests for benchmark
│   ├── models/
│   │   └── models.go           # Shared data types
│   ├── network/
│   │   ├── latency.go          # TCP network latency tests
│   │   └── latency_test.go     # Tests for latency
│   └── tui/
│       ├── tui.go              # Bubble Tea model & update loop
│       ├── screens.go          # TUI screen views
│       └── styles.go           # Lip Gloss style definitions
```

## Troubleshooting

- **Terminal Issues**: Ensure you're using a modern terminal that supports ANSI colors
- **Import Errors**: Run `go mod tidy` to ensure all dependencies are downloaded
- **Firewall**: DNS benchmarks require outbound UDP on port 53; network tests use TCP port 80

## License

MIT
