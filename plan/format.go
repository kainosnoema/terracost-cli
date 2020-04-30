package plan

import (
	"io"
	"strconv"

	"github.com/kainosnoema/terracost/cli/prices"
	"github.com/olekukonko/tablewriter"
)

type pricingTable struct {
	tableData         [][]string
	hourlyTotalBefore float64
	hourlyTotalAfter  float64
	monthlyTotalDelta float64
}

// FormatTable takes an io.Writer (such as os.Stdout) and writes a nicely formatted cost table
func FormatTable(writer io.Writer, resources []Resource) {
	pricing := pricingTable{}

	for _, res := range resources {
		// unable to find prices
		if len(res.Before) == 0 && len(res.After) == 0 {
			pricing.tableData = append(pricing.tableData, []string{
				res.Address,
				res.Action,
				"?",
				"?",
				"?",
				"?",
				"?",
			})
			continue
		}

		if res.Action != "no-op" {
			for priceChange := range res.Before {
				addTableRow(&pricing, res, priceChange)
			}
		}

		for priceChange := range res.After {
			addTableRow(&pricing, res, priceChange)
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{
		"Resource",
		"Action",
		"Service",
		"Usage Operation",
		"Hourly Before",
		"Hourly After",
		"Monthly Delta",
	})
	table.SetFooter([]string{
		"",
		"",
		"",
		"Total",
		"$" + strconv.FormatFloat(pricing.hourlyTotalBefore, 'f', 3, 32),
		"$" + strconv.FormatFloat(pricing.hourlyTotalAfter, 'f', 3, 32),
		"$" + strconv.FormatFloat(pricing.monthlyTotalDelta, 'f', 3, 32),
	})
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	table.SetBorder(false)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
	table.AppendBulk(pricing.tableData)
	table.Render()
}

func addTableRow(pricing *pricingTable, res Resource, priceID prices.PriceID) {
	beforePrice := res.Before[priceID]
	price := res.After[priceID]

	var hourlyBefore, hourlyAfter float64
	if beforePrice != nil && len(beforePrice.Dimensions) > 0 {
		hourlyBefore, _ = strconv.ParseFloat(beforePrice.Dimensions[0].PricePerUnit, 32)
	}

	if price != nil && len(price.Dimensions) > 0 {
		hourlyAfter, _ = strconv.ParseFloat(price.Dimensions[0].PricePerUnit, 32)
	} else if beforePrice != nil {
		price = beforePrice
	}
	monthlyDelta := (hourlyAfter - hourlyBefore) * 730

	pricing.hourlyTotalBefore += hourlyBefore
	pricing.hourlyTotalAfter += hourlyAfter
	pricing.monthlyTotalDelta += monthlyDelta

	formattedBefore := "-"
	if hourlyBefore > 0 {
		formattedBefore = "$" + strconv.FormatFloat(hourlyBefore, 'f', 3, 32)
	}
	formattedAfter := "-"
	if hourlyAfter > 0 {
		formattedAfter = "$" + strconv.FormatFloat(hourlyAfter, 'f', 3, 32)
	}
	formattedMonthlyDelta := "$" + strconv.FormatFloat(monthlyDelta, 'f', 3, 32)

	pricing.tableData = append(pricing.tableData, []string{
		res.Address,
		res.Action,
		price.ServiceCode,
		price.UsageOperation,
		formattedBefore,
		formattedAfter,
		formattedMonthlyDelta,
	})
}
