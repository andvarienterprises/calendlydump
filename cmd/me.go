/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"calendly/calendly"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Shoe information about the current user",
	Long: `Shows informaiton about the authorised user.
`,
	Run: func(cmd *cobra.Command, args []string) {
		doMe(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}

func doMe(cmd *cobra.Command, args []string) {
	c := calendly.NewCalendly()
	err := c.SetAPIKeyFromFile(authTokenFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	me, err := c.GetMe()

	if err != nil {
		log.Fatalf("me: %v", err)
	}

	for r := range me {
		fmt.Printf("%v: %v\n", r, me[r])
	}
}
