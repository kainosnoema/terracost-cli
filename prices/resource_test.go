package prices

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kainosnoema/terracost-cli/terraform"
)

type stubEC2API struct {
	ec2iface.EC2API
}

func (m *stubEC2API) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	output := &ec2.DescribeImagesOutput{
		Images: []*ec2.Image{
			{
				UsageOperation: aws.String("RunInstances"),
			},
		},
	}
	return output, nil
}

type testCase struct {
	changeJSON     terraform.ResourceChangeJSON
	changePriceIDs ChangesPriceIDs
}

func TestResource(t *testing.T) {
	initEC2Api = func(region string) ec2iface.EC2API { return &stubEC2API{} }

	awsInstanceTestCases := []testCase{
		{
			terraform.ResourceChangeJSON{
				Address: "aws_instance.server[0]",
				Type:    "aws_instance",
				Change: terraform.ChangeJSON{
					Actions: []string{"create"},
					Before:  nil,
					After: map[string]interface{}{
						"ami":           "ami-abcd1234",
						"instance_type": "t2.micro",
					},
				},
			},
			ChangesPriceIDs{
				Before: nil,
				After:  []PriceID{{"AmazonEC2", "USW2-BoxUsage:t2.micro:RunInstances"}},
			},
		},
		{
			terraform.ResourceChangeJSON{
				Address: "aws_instance.server[0]",
				Type:    "aws_instance",
				Change: terraform.ChangeJSON{
					Actions: []string{"delete"},
					Before: map[string]interface{}{
						"ami":           "ami-abcd1234",
						"instance_type": "t2.micro",
					},
					After: nil,
				},
			},
			ChangesPriceIDs{
				Before: []PriceID{{"AmazonEC2", "USW2-BoxUsage:t2.micro:RunInstances"}},
				After:  nil,
			},
		},
		{
			terraform.ResourceChangeJSON{
				Address: "aws_instance.server[0]",
				Type:    "aws_instance",
				Change: terraform.ChangeJSON{
					Actions: []string{"update"},
					Before: map[string]interface{}{
						"ami":           "ami-abcd1234",
						"instance_type": "t2.micro",
					},
					After: map[string]interface{}{
						"ami":           "ami-abcd1234",
						"instance_type": "t2.small",
					},
				},
			},
			ChangesPriceIDs{
				Before: []PriceID{{"AmazonEC2", "USW2-BoxUsage:t2.micro:RunInstances"}},
				After:  []PriceID{{"AmazonEC2", "USW2-BoxUsage:t2.small:RunInstances"}},
			},
		},
	}

	for _, testCase := range awsInstanceTestCases {
		t.Run(testCase.changeJSON.Address, func(t *testing.T) {
			changePriceIDs := ResourceChangesPriceIDs("us-west-2", testCase.changeJSON)
			if !reflect.DeepEqual(changePriceIDs, testCase.changePriceIDs) {
				t.Errorf("incorrect price IDs: %v, expected: %v", changePriceIDs, testCase.changePriceIDs)
			}
		})
	}

}
