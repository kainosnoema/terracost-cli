package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kainosnoema/terracost/cli/plan"
	"github.com/kainosnoema/terracost/cli/terraform"
)

func init() {
	rootCmd.AddCommand(estimateCmd)
}

var estimateCmd = &cobra.Command{
	Use:   "estimate [planfile]",
	Short: "Plan and estimate costs for a Terraform project or plan file",
	Long:  "",
	Args:  cobra.RangeArgs(0, 1),
	Run:   runEstimate,
}

func runEstimate(cmd *cobra.Command, args []string) {
	var tfPlan *terraform.PlanJSON
	var err error

	if len(args) > 0 {
		tfPlan, err = terraform.ShowPlan(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error running Terraform:", err.Error())
			return
		}
	} else {
		fmt.Println("Planning...")
		tfPlan, err = terraform.ExecPlan()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error running Terraform:", err.Error())
			return
		}
	}

	fmt.Println("Estimating...")
	resources, err := plan.Calculate(tfPlan)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error estimating:", err.Error())
		return
	}

	fmt.Println()
	plan.FormatTable(os.Stdout, resources)
}
