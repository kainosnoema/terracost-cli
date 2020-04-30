package prices

import (
	"github.com/kainosnoema/terracost/cli/terraform"
)

// Difficult to compile an accurate list as there's no complete documentation.
// Example: https://docs.aws.amazon.com/AmazonS3/latest/dev/aws-usage-report-understand.html

var regionMap = map[string]string{
	"af-south-1":     "AFS1",
	"ap-east-1":      "APE1",
	"ap-northeast-1": "APN1",
	"ap-northeast-2": "APN2",
	"ap-northeast-3": "APN3",
	"ap-south-1":     "APS3", // confusing but accurate
	"ap-southeast-1": "APS1",
	"ap-southeast-2": "APS2",
	"ca-central-1":   "CAN1",
	"cn-north-1":     "CNN1", // or CN1 or CN
	"cn-northwest-1": "CNW1",
	"eu-central-1":   "EUC1",
	"eu-north-1":     "EUN1",
	"eu-south-1":     "EUS1",
	"eu-west-1":      "EU", // or EUW1
	"eu-west-2":      "EUW2",
	"eu-west-3":      "EUW3",
	"me-south-1":     "MES1",
	"sa-east-1":      "SAE1",
	"us-east-1":      "USE1",
	"us-east-2":      "USE2",
	"us-gov-west-1":  "UGW1", // or GOVW1 or GOV
	"us-gov-east-1":  "UGE1", // or GOVE1
	"us-west-1":      "USW1",
	"us-west-2":      "USW2",
}

type ChangesQueries struct {
	Before []PriceQuery
	After  []PriceQuery
}

func ResourceChangesQueries(region string, res terraform.ResourceChangeJSON) ChangesQueries {
	switch res.Type {
	case "aws_instance":
		return AWSInstance(region, res.Change)
	case "aws_db_instance":
		return AWSDBInstance(region, res.Change)
	default:
		return ChangesQueries{}
	}
}
