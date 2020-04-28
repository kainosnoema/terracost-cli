package plan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/machinebox/graphql"
)

var regionMap = map[string]string{
	"us-east-1": "USE1",
	"us-east-2": "USE2",
	"us-west-1": "USW1",
	"us-west-2": "USW2",
}

// Resource maps a Terraform resource to an AWS service, usage, and pricing
type Resource struct {
	Type           string
	Name           string
	ServiceCode    string
	UsageOperation string
	Price          Price
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
		switch res.Type {
		case "aws_instance":
			if res.Change.Before == nil { // creating
				usageOperation := fmt.Sprintf(
					"%s-BoxUsage:%s:%s",
					regionMap[region],
					res.Change.After["instance_type"].(string),
					imageUsageOperation(region, res.Change.After["ami"].(string)),
				)
				resources = append(resources, Resource{
					Type:           res.Type,
					Name:           res.Name,
					ServiceCode:    "AmazonEC2",
					UsageOperation: usageOperation,
				})
				prices[priceQuery{
					ServiceCode:    "AmazonEC2",
					UsageOperation: usageOperation,
				}] = Price{}
			} else if res.Change.After == nil { // deleting

			} else { // updating

			}
		default:
			return nil, fmt.Errorf("resource type not recognized: %s", res.Type)
		}
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
			ServiceCode:    resource.ServiceCode,
			UsageOperation: resource.UsageOperation,
		}]
		resources[k] = resource
	}

	return resources, nil
}

func imageUsageOperation(region, ami string) string {
	svc := ec2.New(session.New(&aws.Config{Region: &region}))
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{
			aws.String(ami),
		},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		fmt.Println(err.Error())
		return "RunInstances"
	}

	for _, img := range result.Images {
		return aws.StringValue(img.UsageOperation)
	}

	return "RunInstances"
}
