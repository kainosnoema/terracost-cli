package prices

import (
	"context"

	"github.com/machinebox/graphql"
)

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

// Price represents AWS pricing information
type Price struct {
	PriceID
	Dimensions []Dimension
	Updated    string
}

// PriceID is a composite key of AWS service code and usage operation
type PriceID struct {
	ServiceCode    string
	UsageOperation string
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

// ByID maps a price query to the price result
type ByID map[PriceID]*Price

// Lookup gathers a unique set of PriceIDs and looks them up using the Terracost API
type Lookup interface {
	Add(id PriceID) *Price
	Perform() error
}

// NewLookup returns a new Lookup
func NewLookup() Lookup {
	return &lookup{prices: ByID{}}
}

type lookup struct {
	prices ByID
}

// Add tales a PriceID and returns a price pointer that will be populated when looked up
func (pl *lookup) Add(id PriceID) *Price {
	pendingPrice := pl.prices[id]
	if pendingPrice == nil {
		pendingPrice = &Price{PriceID: id}
		pl.prices[id] = pendingPrice
	}
	return pendingPrice
}

// Perform hits the Terracost API with the list of PriceIDs and populates the prices
func (pl *lookup) Perform() error {
	if len(pl.prices) == 0 {
		return nil
	}

	lookupIDs := []PriceID{}
	for q := range pl.prices {
		lookupIDs = append(lookupIDs, q)
	}

	req := graphql.NewRequest(priceQueryGql)
	req.Var("lookup", lookupIDs)

	var response struct{ Prices []Price }
	client := graphql.NewClient(apiURL)
	if err := client.Run(context.Background(), req, &response); err != nil {
		return err
	}

	for _, p := range response.Prices {
		price := pl.prices[PriceID{
			ServiceCode:    p.ServiceCode,
			UsageOperation: p.UsageOperation,
		}]
		price.Dimensions = p.Dimensions
		price.Updated = p.Updated
	}

	return nil
}
