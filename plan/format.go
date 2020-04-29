package plan

import (
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// FormatTable takes an io.Writer (such as os.Stdout) and writes a nicely formatted cost table
func FormatTable(writer io.Writer, resources []Resource) {
	tableData := [][]string{}
	hourlyTotal := 0.0
	monthlyTotal := 0.0

	for _, res := range resources {
		if len(res.Prices[0].Dimensions) == 0 {
			tableData = append(tableData, []string{
				res.Type,
				res.Name,
				res.Prices[0].ServiceCode,
				res.Prices[0].UsageOperation,
				"(unknown)",
				"(unknown)",
			})
			continue
		}

		hourlyCost, _ := strconv.ParseFloat(res.Prices[0].Dimensions[0].PricePerUnit, 32)
		monthlyCost := hourlyCost * 730
		hourlyTotal += hourlyCost
		monthlyTotal += monthlyCost

		tableData = append(tableData, []string{
			res.Type,
			res.Name,
			res.Prices[0].ServiceCode,
			res.Prices[0].UsageOperation,
			"$" + strconv.FormatFloat(hourlyCost, 'f', 3, 32),
			"$" + strconv.FormatFloat(monthlyCost, 'f', 2, 32),
		})
	}

	table := tablewriter.NewWriter(writer)
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
	table.AppendBulk(tableData)
	table.Render()
}
