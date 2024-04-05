package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	authTokenFile string
	rootCmd       = &cobra.Command{
		Use:   "calendlydump",
		Short: "calendlydump is a CLI tool to spit out calendly event and invitee data",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&authTokenFile, "calendly-auth-token", "calendly.key", "The calendly authentication key filename")
}
