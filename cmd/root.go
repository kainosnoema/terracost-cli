package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terracost",
	Short: "CLI for Terracost.",
	Long:  "",
}

// Execute root command
func Execute() error {
	return rootCmd.Execute()
}
