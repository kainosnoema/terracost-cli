package plan

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"

	"github.com/kainosnoema/terracost/cli/plan/mappings"
)

// Resource maps a Terraform resource to AWS pricing
type Resource struct {
	Type  string
	Name  string
	Price Price
}

// Price represents AWS pricing information
type Price struct {
	ServiceCode    string
	UsageOperation string
	Dimensions     []dimension
	Updated        string
}

type dimension struct {
	BeginRange   string
	EndRange     string
	PricePerUnit string
	Unit         string
	RateCode     string
	Description  string
}

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

type priceQuery struct {
	ServiceCode    string
	UsageOperation string
}

type priceResponse struct {
	Prices []Price
}

// Calculate takes a TF plan, fetches AWS prices, and returns priced Resources
func Calculate(tfPlan *PlanJSON) ([]Resource, error) {
	region := tfPlan.Configuration.ProviderConfig.AWS.Expressions.Region.ConstantValue
	resources := []Resource{}
	prices := map[priceQuery]Price{}

	for _, res := range tfPlan.ResourceChanges {
		var serviceCode string
		var usageOperation string

		switch res.Type {
		case "aws_instance":
			if res.Change.Before == nil { // creating
				serviceCode = "AmazonEC2"
				usageOperation = mappings.EC2Instance(region, res.Change.After)
			} else if res.Change.After == nil { // deleting

			} else { // updating

			}
		case "aws_db_instance":
			if res.Change.Before == nil { // creating
				serviceCode = "AmazonRDS"
				usageOperation = mappings.RDSInstance(region, res.Change.After)
			} else if res.Change.After == nil { // deleting

			} else { // updating

			}
		default:
			return nil, fmt.Errorf("resource type not supported: %s", res.Type)
		}

		priceQuery := priceQuery{
			ServiceCode:    serviceCode,
			UsageOperation: usageOperation,
		}
		prices[priceQuery] = Price{
			ServiceCode:    priceQuery.ServiceCode,
			UsageOperation: priceQuery.UsageOperation,
		}
		resources = append(resources, Resource{
			Type:  res.Type,
			Name:  res.Name,
			Price: prices[priceQuery],
		})
	}

	var lookup []priceQuery
	for q := range prices {
		lookup = append(lookup, q)
	}

	req := graphql.NewRequest(priceQueryGql)
	req.Var("lookup", lookup)

	var res priceResponse
	client := graphql.NewClient("http://localhost:3000/api/graphql")
	if err := client.Run(context.Background(), req, &res); err != nil {
		return nil, err
	}

	for _, price := range res.Prices {
		prices[priceQuery{
			ServiceCode:    price.ServiceCode,
			UsageOperation: price.UsageOperation,
		}] = price
	}

	for k, resource := range resources {
		resource.Price = prices[priceQuery{
			ServiceCode:    resource.Price.ServiceCode,
			UsageOperation: resource.Price.UsageOperation,
		}]
		resources[k] = resource
	}

	return resources, nil
}
