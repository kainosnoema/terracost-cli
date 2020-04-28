package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/machinebox/graphql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	planTypes "github.com/kainosnoema/terracost/cli/plan"
)

var regionMap = map[string]string{
	"us-east-1": "USE1",
	"us-east-2": "USE2",
	"us-west-1": "USW1",
	"us-west-2": "USW2",
}

var query = `query ($lookup: [PriceQuery!]) {
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

type price struct {
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

type priceQuery struct {
	ServiceCode    string
	UsageOperation string
}

type resource struct {
	Type           string
	Name           string
	ServiceCode    string
	UsageOperation string
}

type response struct {
	Prices []price
}

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan [plan-file]",
	Short: "Calculate costs for a Terraform plan file",
	Long:  "",
	Run:   plan,
}

func plan(cmd *cobra.Command, args []string) {
	fmt.Println("Planning...")

	planFile, err := ioutil.TempFile("", "tc-plan")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(planFile.Name())

	planCmd := exec.Command("terraform", "plan", "-out", planFile.Name())
	// planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr
	err = planCmd.Run()
	if err != nil {
		log.Fatalf("planCmd.Run() failed with %s\n", err)
	}

	fmt.Println("Calculating...")

	showCmd := exec.Command("terraform", "show", "-json", planFile.Name())
	out, err := showCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("showCmd.Run() failed with %s\n", err)
	}

	var tfPlan planTypes.PlanJSON
	err = json.Unmarshal(out, &tfPlan)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	resources := []resource{}
	prices := map[priceQuery]price{}
	region := tfPlan.Configuration.ProviderConfig.AWS.Expressions.Region.ConstantValue
	for _, res := range tfPlan.ResourceChanges {
		switch res.Type {
		case "aws_instance":
			if res.Change.Before == nil { // creating
				ami := res.Change.After["ami"].(string)
				instanceType := res.Change.After["instance_type"].(string)
				usageOperation := fmt.Sprintf(
					"%s-BoxUsage:%s:%s",
					regionMap[region],
					instanceType,
					imageUsageOperation(region, ami),
				)
				resources = append(resources, resource{
					Type:           res.Type,
					Name:           res.Name,
					ServiceCode:    "AmazonEC2",
					UsageOperation: usageOperation,
				})
				prices[priceQuery{
					ServiceCode:    "AmazonEC2",
					UsageOperation: usageOperation,
				}] = price{}
			} else if res.Change.After == nil { // deleting

			} else { // updating

			}
		default:
			fmt.Println("resource type not recognized: ", res.Type)
		}
	}

	var lookup []priceQuery
	for q := range prices {
		lookup = append(lookup, q)
	}

	req := graphql.NewRequest(query)
	req.Var("lookup", lookup)

	var res response
	client := graphql.NewClient("http://localhost:3000/api/graphql")
	if err := client.Run(context.Background(), req, &res); err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	for _, price := range res.Prices {
		prices[priceQuery{
			ServiceCode:    price.ServiceCode,
			UsageOperation: price.UsageOperation,
		}] = price
	}

	data := [][]string{}
	hourlyTotal := 0.0
	monthlyTotal := 0.0
	for _, resource := range resources {
		price := prices[priceQuery{
			ServiceCode:    resource.ServiceCode,
			UsageOperation: resource.UsageOperation,
		}]
		hourlyCost, _ := strconv.ParseFloat(price.Dimensions[0].PricePerUnit, 32)
		monthlyCost := hourlyCost * 730
		hourlyTotal += hourlyCost
		monthlyTotal += monthlyCost
		data = append(data, []string{
			resource.Type,
			resource.Name,
			resource.ServiceCode,
			resource.UsageOperation,
			"$" + strconv.FormatFloat(hourlyCost, 'f', 3, 32),
			"$" + strconv.FormatFloat(monthlyCost, 'f', 3, 32),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Resource",
		"Name",
		"Service",
		"Usage Operation",
		"Hourly Cost",
		"Monthly Cost",
	})
	table.SetFooter([]string{
		"",
		"",
		"",
		"Total",
		"$" + strconv.FormatFloat(hourlyTotal, 'f', 3, 32),
		"$" + strconv.FormatFloat(monthlyTotal, 'f', 3, 32),
	})
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
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
