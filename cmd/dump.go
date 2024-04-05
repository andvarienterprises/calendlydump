package cmd

import (
	"github.com/andvarienterprises/calendlydump/calendly"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump calendly data in CSV format",
	Long: `Dump will print all calendly events in CSV format.
`,
	Run: func(cmd *cobra.Command, args []string) {
		doDump(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}

func doDump(_ *cobra.Command, args []string) {
	if len(args) > 1 {
		log.Fatal("dump takes 0 or 1 arguments")
	}
	to_dump := "events"
	if len(args) == 1 {
		switch args[0] {
		case "invites":
			to_dump = "invites"
		case "events":
		default:
			log.Fatal("invalid argument to dump, need 'invites' or 'events'")
		}
	}

	c := calendly.NewCalendly()
	err := c.SetAPIKeyFromFile(authTokenFile)
	if err != nil {
		log.Fatalf("dump: %v", err)
	}
	events, err := c.GetEvents()
	if err != nil {
		log.Fatalf("dump: %v", err)
	}

	if to_dump == "events" {
		fmt.Print(calendly.EventsCSV(events))
	} else {
		// Populate each 'invitees' field first.
		for _, e := range events {
			err := c.PopulateInviteesForEvent(e)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Print(calendly.EventInvitesCSV(events))
	}
}
