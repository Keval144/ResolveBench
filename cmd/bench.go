package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"dnsbench/internal/dns"
	"dnsbench/internal/models"
	"dnsbench/internal/network"
	"dnsbench/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type jsonOutput struct {
	Providers  []models.ProviderResult `json:"providers"`
	Network    []models.NetworkResult  `json:"network"`
}

var (
	jsonFlag  bool
	dnsFlag   string
)

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Run DNS benchmark",
	Long: `Run the DNS benchmark and display results.

By default launches the TUI interactive mode.
Use --json for non-interactive JSON output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonFlag {
			return runBenchJSON()
		}
		return runBenchTUI()
	},
}

func init() {
	benchCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output results as JSON")
	benchCmd.Flags().StringVar(&dnsFlag, "dns", "", "Comma-separated custom DNS servers (e.g. 1.1.1.1,8.8.8.8)")
}

func runBenchTUI() error {
	m := tui.NewModel()
	p := tea.NewProgram(&m, tea.WithAltScreen())
	m.Program = p

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func runBenchJSON() error {
	var results []models.ProviderResult
	var netResults []models.NetworkResult

	if dnsFlag != "" {
		customServers := strings.Split(dnsFlag, ",")
		results = dns.RunBenchmarkCustom(customServers, nil)
	} else {
		results = dns.RunBenchmarkCustom(nil, nil)
	}

	netResults = network.RunLatencyTests(nil)

	output := jsonOutput{
		Providers: results,
		Network:   netResults,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}
