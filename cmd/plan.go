package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kainosnoema/terracost/cli/plan"
	"github.com/kainosnoema/terracost/cli/terraform"
)

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan [plan-file]",
	Short: "Calculate costs for a Terraform plan file",
	Long:  "",
	Run:   runPlan,
}

func runPlan(cmd *cobra.Command, args []string) {
	fmt.Println("Planning...")
	tfPlan, err := terraform.ExecPlan()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running Terraform:", err.Error())
		return
	}

	fmt.Println("Calculating...")
	resources, err := plan.Calculate(tfPlan)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error calculating:", err.Error())
		return
	}

	fmt.Println()
	plan.FormatTable(os.Stdout, resources)
}
