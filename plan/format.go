package plan

import (
	"io"
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
	"github.com/mitchellh/colorstring"
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

	table.SetFooter([]string{
		"",
		"",
		"Total",
		money.FormatMoney(pricing.hourlyTotalBefore),
		money.FormatMoney(pricing.hourlyTotalAfter),
		money.FormatMoney(pricing.monthlyTotalDelta),
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
	price := res.After[priceID]
	beforePrice := findBeforePrice(res.Before, priceID)

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

	pricing.tableData = append(pricing.tableData, []string{
		formatAddress(res),
		price.ServiceCode,
		formatDescription(beforePrice, price),
		money.FormatMoney(hourlyBefore),
		money.FormatMoney(hourlyAfter),
		money.FormatMoney(monthlyDelta),
	})
}

func formatAddress(res Resource) string {
	actionIcon := ""
	switch res.Action {
	case "create":
		actionIcon = "[green]+[reset]"
	case "delete":
		actionIcon = "[red]-[reset]"
	case "update":
		actionIcon = "[yellow]~[reset]"
	default:
	}

	return colorstring.Color(actionIcon + " " + res.Address)
}

func formatDescription(beforePrice, price *prices.Price) string {
	if beforePrice != nil && price != nil {
		if beforePrice.UsageOperation == price.UsageOperation {
			return price.UsageOperation
		}
		return beforePrice.UsageOperation + " -> " + price.UsageOperation
	} else if beforePrice != nil {
		return beforePrice.UsageOperation
	}

	return price.UsageOperation
}

func findBeforePrice(beforePrices prices.ByID, priceID prices.PriceID) *prices.Price {
	if beforePrice, ok := beforePrices[priceID]; ok {
		return beforePrice
	}
	usagePrefix := strings.SplitN(priceID.UsageOperation, ":", 2)[0]
	for id, price := range beforePrices {
		if strings.HasPrefix(id.UsageOperation, usagePrefix) {
			return price
		}
	}
	return nil
}
