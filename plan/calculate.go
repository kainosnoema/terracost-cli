package plan

import (
	"fmt"

	"github.com/kainosnoema/terracost/cli/mapping"
	"github.com/kainosnoema/terracost/cli/prices"
)

// Resource maps a Terraform resource to AWS pricing
type Resource struct {
	Type   string
	Name   string
	Prices []prices.Price
}

// Calculate takes a TF plan, fetches AWS prices, and returns priced Resources
func Calculate(tfPlan *PlanJSON) ([]Resource, error) {
	region := tfPlan.Region()
	resources := []Resource{}
	priceQueries := map[prices.PriceQuery]prices.Price{}

	for _, res := range tfPlan.ResourceChanges {
		var serviceCode string
		var usageOperation string

		switch res.Type {
		case "aws_instance":
			if res.Change.Before == nil { // creating
				serviceCode = "AmazonEC2"
				usageOperation = mapping.EC2Instance(region, res.Change.After)
			} else if res.Change.After == nil { // deleting

			} else { // updating

			}
		case "aws_db_instance":
			if res.Change.Before == nil { // creating
				serviceCode = "AmazonRDS"
				usageOperation = mapping.RDSInstance(region, res.Change.After)
			} else if res.Change.After == nil { // deleting

			} else { // updating

			}
		default:
			return nil, fmt.Errorf("resource type not supported: %s", res.Type)
		}

		priceQuery := prices.PriceQuery{
			ServiceCode:    serviceCode,
			UsageOperation: usageOperation,
		}
		priceQueries[priceQuery] = prices.Price{
			ServiceCode:    priceQuery.ServiceCode,
			UsageOperation: priceQuery.UsageOperation,
		}
		resources = append(resources, Resource{
			Type:   res.Type,
			Name:   res.Name,
			Prices: []prices.Price{priceQueries[priceQuery]},
		})
	}

	var lookup []prices.PriceQuery
	for q := range priceQueries {
		lookup = append(lookup, q)
	}

	priceRes, err := prices.Lookup(lookup)
	if err != nil {
		return nil, err
	}

	for _, price := range priceRes {
		priceQueries[prices.PriceQuery{
			ServiceCode:    price.ServiceCode,
			UsageOperation: price.UsageOperation,
		}] = price
	}

	for k, resource := range resources {
		resource.Prices = []prices.Price{priceQueries[prices.PriceQuery{
			ServiceCode:    resource.Prices[0].ServiceCode,
			UsageOperation: resource.Prices[0].UsageOperation,
		}]}
		resources[k] = resource
	}

	return resources, nil
}
