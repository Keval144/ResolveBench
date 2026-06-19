package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dnsbench",
	Short: "DNS benchmark tool with TUI",
	Long: `DnsBench benchmarks different DNS providers to help you find
the lowest-latency DNS for optimal routing.

It tests Cloudflare, Google DNS, Quad9, AdGuard DNS, and OpenDNS
by resolving popular domains and computing latency, reliability, and consistency scores.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(benchCmd)
	rootCmd.AddCommand(listDNSCmd)
}
