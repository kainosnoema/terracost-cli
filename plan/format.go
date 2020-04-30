package plan

import (
	"io"
	"strconv"

	"github.com/leekchan/accounting"
	"github.com/olekukonko/tablewriter"

	"github.com/kainosnoema/terracost/cli/prices"
)

var money = &accounting.Accounting{
	Symbol: "$", Precision: 3, Format: "%s%v", FormatZero: "-",
}

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
				formatAddress(res),
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
		"Service",
		"Usage Operation",
		"Hourly Before",
		"Hourly After",
		"Monthly Delta",
	})

	formattedHourlyBefore := money.FormatMoney(pricing.hourlyTotalBefore)
	formattedHourlyAfter := money.FormatMoney(pricing.hourlyTotalAfter)
	formattedMonthlyDelta := money.FormatMoney(pricing.monthlyTotalDelta)
	if pricing.monthlyTotalDelta > 0 {
		formattedMonthlyDelta = "+" + formattedMonthlyDelta
	}

	table.SetFooter([]string{
		"",
		"",
		"Total",
		formattedHourlyBefore,
		formattedHourlyAfter,
		formattedMonthlyDelta,
	})
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	table.SetBorder(false)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetAutoWrapText(false)
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

	formattedBefore := money.FormatMoney(hourlyBefore)
	formattedAfter := money.FormatMoney(hourlyAfter)
	formattedMonthlyDelta := money.FormatMoney(monthlyDelta)
	if monthlyDelta > 0 {
		formattedMonthlyDelta = "+" + formattedMonthlyDelta
	}

	pricing.tableData = append(pricing.tableData, []string{
		formatAddress(res),
		price.ServiceCode,
		price.UsageOperation,
		formattedBefore,
		formattedAfter,
		formattedMonthlyDelta,
	})
}

func formatAddress(res Resource) string {
	actionIcon := ""
	switch res.Action {
	case "create":
		actionIcon = "+"
	case "delete":
		actionIcon = "-"
	case "update":
		actionIcon = "~"
	default:
	}

	return actionIcon + " " + res.Address
}
