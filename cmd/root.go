package cmd

import (
	"fmt"
	"os"

	"resolvebench/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "resolvebench",
	Short: "DNS benchmark tool with TUI",
	Long: `ResolveBench benchmarks different DNS providers to help you find
the fastest, most reliable DNS resolver for your network.

It tests Cloudflare, Google DNS, Quad9, AdGuard DNS, and OpenDNS
by resolving popular domains and computing latency, reliability, and consistency scores.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		m := tui.NewModel()
		p := tea.NewProgram(&m, tea.WithAltScreen())
		m.Program = p

		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listDNSCmd)
}
