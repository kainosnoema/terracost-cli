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
	Use:   "plan [planfile]",
	Short: "Plan and calculate costs for a Terraform project or plan file",
	Long:  "",
	Args:  cobra.RangeArgs(0, 1),
	Run:   runPlan,
}

func runPlan(cmd *cobra.Command, args []string) {
	fmt.Println("Planning...")

	var tfPlan *terraform.PlanJSON
	var err error

	if len(args) > 0 {
		tfPlan, err = terraform.ShowPlan(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error running Terraform:", err.Error())
			return
		}
	} else {
		tfPlan, err = terraform.ExecPlan()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error running Terraform:", err.Error())
			return
		}
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
