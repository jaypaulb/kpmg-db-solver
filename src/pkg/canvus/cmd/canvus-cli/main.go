package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "canvus-cli",
		Short: "Canvus CLI - Manage Canvus resources from the command line",
		Long:  `A command-line tool for interacting with the Canvus API.`,
	}

	rootCmd.AddCommand(cleanupCmd)
	// TODO: Add more subcommands (list, create, delete, upload, etc.)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup test users and folders (for development/testing)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cleanup command not yet implemented.")
	},
}
