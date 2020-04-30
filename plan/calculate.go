package plan

import (
	"github.com/kainosnoema/terracost/cli/prices"
	"github.com/kainosnoema/terracost/cli/terraform"
)

// Resource maps a Terraform resource to AWS pricing
type Resource struct {
	Address string
	Action  string
	Before  Prices
	After   Prices
}

// Prices maps a price query to the price result
type Prices map[prices.PriceQuery]*prices.Price

// Calculate takes a TF plan, fetches AWS prices, and returns priced Resources
func Calculate(tfPlan *terraform.PlanJSON) ([]Resource, error) {
	region := tfPlan.Region()
	resources := []Resource{}
	pendingPrices := Prices{}

	for _, res := range tfPlan.ResourceChanges {
		resource := Resource{
			Address: res.Address,
			Action:  res.Change.Actions[0],
			Before:  Prices{},
			After:   Prices{},
		}

		changesQueries := prices.ResourceChangesQueries(region, res)

		for _, beforeQuery := range changesQueries.Before {
			pendingPrice := pendingPrices[beforeQuery]
			if pendingPrice == nil {
				pendingPrice = &prices.Price{
					ServiceCode:    beforeQuery.ServiceCode,
					UsageOperation: beforeQuery.UsageOperation,
				}
				pendingPrices[beforeQuery] = pendingPrice
			}
			resource.Before[beforeQuery] = pendingPrice
		}

		for _, afterQuery := range changesQueries.After {
			pendingPrice := pendingPrices[afterQuery]
			if pendingPrice == nil {
				pendingPrice = &prices.Price{
					ServiceCode:    afterQuery.ServiceCode,
					UsageOperation: afterQuery.UsageOperation,
				}
				pendingPrices[afterQuery] = pendingPrice
			}
			resource.After[afterQuery] = pendingPrice
		}

		resources = append(resources, resource)
	}

	var queries []prices.PriceQuery
	for q := range pendingPrices {
		queries = append(queries, q)
	}

	priceRes, err := prices.LookupQueries(queries)
	if err != nil {
		return nil, err
	}

	for _, price := range priceRes {
		pendingPrice := pendingPrices[prices.PriceQuery{
			ServiceCode:    price.ServiceCode,
			UsageOperation: price.UsageOperation,
		}]
		pendingPrice.Dimensions = price.Dimensions
		pendingPrice.Updated = price.Updated
	}

	return resources, nil
}
