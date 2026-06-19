package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"dnsbench/internal/dns"

	"github.com/spf13/cobra"
)

var listDNSCmd = &cobra.Command{
	Use:   "list-dns",
	Short: "List available DNS providers",
	Long:  "Display all DNS providers that DnsBench can benchmark.",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "Provider\tPrimary DNS\tSecondary DNS\tUse Case")
		fmt.Fprintln(w, "--------\t-----------\t-------------\t--------")
		for _, p := range dns.Providers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.PrimaryDNS, p.SecondaryDNS, p.UseCase)
		}
		w.Flush()
	},
}
