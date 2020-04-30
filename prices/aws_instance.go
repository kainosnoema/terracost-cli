package prices

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kainosnoema/terracost/cli/terraform"
)

func AWSInstance(region string, changes terraform.ChangeJSON) ChangesPriceIDs {
	changesPriceIDs := ChangesPriceIDs{}

	if changes.Before != nil {
		changesPriceIDs.Before = []PriceID{ec2Instance(region, changes.Before)}
	}

	if changes.After != nil {
		changesPriceIDs.After = []PriceID{ec2Instance(region, changes.After)}
	}

	return changesPriceIDs
}

func ec2Instance(region string, changeAttrs map[string]interface{}) PriceID {
	ec2UsageOperation := fmt.Sprintf("%s-BoxUsage:%s:%s",
		regionMap[region],
		changeAttrs["instance_type"].(string),
		imageUsageOperation(region, changeAttrs["ami"].(string)),
	)

	return PriceID{
		ServiceCode:    "AmazonEC2",
		UsageOperation: ec2UsageOperation,
	}
}

// TODO: make a single API call for all AMIs
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
