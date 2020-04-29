package plan

import (
	"github.com/kainosnoema/terracost/cli/prices"
	"github.com/kainosnoema/terracost/cli/terraform"
)

// Resource maps a Terraform resource to AWS pricing
type Resource struct {
	Type   string
	Name   string
	Prices map[prices.PriceQuery]prices.Price
}

// Calculate takes a TF plan, fetches AWS prices, and returns priced Resources
func Calculate(tfPlan *terraform.PlanJSON) ([]Resource, error) {
	region := tfPlan.Region()
	resources := []Resource{}
	allQueries := map[prices.PriceQuery]prices.Price{}

	for _, res := range tfPlan.ResourceChanges {
		resource := Resource{
			Type:   res.Type,
			Name:   res.Name,
			Prices: map[prices.PriceQuery]prices.Price{},
		}

		for _, priceQuery := range prices.Resource(region, res) {
			emptyPrice := prices.Price{
				ServiceCode:    priceQuery.ServiceCode,
				UsageOperation: priceQuery.UsageOperation,
			}
			resource.Prices[priceQuery] = emptyPrice
			allQueries[priceQuery] = emptyPrice
		}

		resources = append(resources, resource)
	}

	var lookup []prices.PriceQuery
	for q := range allQueries {
		lookup = append(lookup, q)
	}

	priceRes, err := prices.Lookup(lookup)
	if err != nil {
		return nil, err
	}

	for _, price := range priceRes {
		allQueries[prices.PriceQuery{
			ServiceCode:    price.ServiceCode,
			UsageOperation: price.UsageOperation,
		}] = price
	}

	for k, resource := range resources {
		resPrices := resource.Prices
		for priceQuery := range resPrices {
			resPrices[priceQuery] = allQueries[prices.PriceQuery{
				ServiceCode:    priceQuery.ServiceCode,
				UsageOperation: priceQuery.UsageOperation,
			}]
		}
		resource.Prices = resPrices
		resources[k] = resource
	}

	return resources, nil
}
