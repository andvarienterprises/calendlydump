package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	authTokenFile string
	rootCmd       = &cobra.Command{
		Use:   "calendly",
		Short: "calendly is a CLI tool to manage your Calendly account.",
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&authTokenFile, "calendly-auth-token", "calendly.key", "The calendly authentication key filename")
}
