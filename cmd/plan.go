package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/kainosnoema/terracost/cli/plan"
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
	tfPlan, err := plan.ExecTerraform()
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

	data := [][]string{}
	hourlyTotal := 0.0
	monthlyTotal := 0.0
	for _, resource := range resources {
		hourlyCost, _ := strconv.ParseFloat(resource.Price.Dimensions[0].PricePerUnit, 32)
		monthlyCost := hourlyCost * 730
		hourlyTotal += hourlyCost
		monthlyTotal += monthlyCost
		data = append(data, []string{
			resource.Type,
			resource.Name,
			resource.ServiceCode,
			resource.UsageOperation,
			"$" + strconv.FormatFloat(hourlyCost, 'f', 3, 32),
			"$" + strconv.FormatFloat(monthlyCost, 'f', 2, 32),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Resource",
		"Name",
		"Service",
		"Usage Operation",
		"Hourly Cost",
		"Monthly Cost",
	})
	table.SetFooter([]string{
		"",
		"",
		"",
		"Total",
		"$" + strconv.FormatFloat(hourlyTotal, 'f', 3, 32),
		"$" + strconv.FormatFloat(monthlyTotal, 'f', 2, 32),
	})
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
