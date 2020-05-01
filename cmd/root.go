package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terracost",
	Short: "AWS cost estimation for Terraform projects.",
	Long:  "",
}

// Execute root command
func Execute() error {
	return rootCmd.Execute()
}
