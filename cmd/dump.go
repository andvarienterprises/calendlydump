package cmd

import (
	"calendly/calendly"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump calendly data in CSV format",
	Long: `Dump will export your calendly data in CSV format.
`,
	Run: func(cmd *cobra.Command, args []string) {
		doDump(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}

func doDump(cmd *cobra.Command, args []string) {
	c := calendly.NewCalendly()
	err := c.SetAPIKeyFromFile(authTokenFile)
	if err != nil {
		log.Fatalf("dump: %v", err)
	}
	events, err := c.GetEvents()

	if err != nil {
		log.Fatalf("dump: %v", err)
	}

	fmt.Print(calendly.EventsCSV(events))
}
