package prices

import (
	"context"

	"github.com/machinebox/graphql"
)

// Price represents AWS pricing information
type Price struct {
	ServiceCode    string
	UsageOperation string
	Dimensions     []Dimension
	Updated        string
}

// Dimension is the price for a specific dimension/range of the product
type Dimension struct {
	BeginRange   string
	EndRange     string
	PricePerUnit string
	Unit         string
	RateCode     string
	Description  string
}

// PriceQuery is used to look up pricing from the Terracost API
type PriceQuery struct {
	ServiceCode    string
	UsageOperation string
}

var apiURL = "http://localhost:3000/api/graphql"
var priceQueryGql = `query ($lookup: [PriceQuery!]) {
	Prices(lookup: $lookup) {
		ServiceCode
		UsageOperation
		Dimensions {
			BeginRange
			EndRange
			Unit
			PricePerUnit
			Description
		}
	}
}`

// Lookup hits the Terracost API with a list of queries and returns prices
func Lookup(queries []PriceQuery) ([]Price, error) {
	req := graphql.NewRequest(priceQueryGql)
	req.Var("lookup", queries)

	var response struct{ Prices []Price }
	client := graphql.NewClient(apiURL)
	if err := client.Run(context.Background(), req, &response); err != nil {
		return nil, err
	}

	return response.Prices, nil
}
