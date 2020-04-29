package mapping

import (
	"github.com/kainosnoema/terracost/cli/prices"
	"github.com/kainosnoema/terracost/cli/terraform"
)

func Resource(region string, res terraform.ResourceChangeJSON) []prices.PriceQuery {
	priceQueries := []prices.PriceQuery{}

	switch res.Type {
	case "aws_instance":
		if res.Change.Before == nil { // creating
			priceQueries = append(priceQueries, EC2Instance(region, res.Change.After)...)
		} else if res.Change.After == nil { // deleting
		} else { // updating
		}
	case "aws_db_instance":
		if res.Change.Before == nil { // creating
			priceQueries = append(priceQueries, RDSInstance(region, res.Change.After)...)
		} else if res.Change.After == nil { // deleting
		} else { // updating
		}
	default:
	}

	return priceQueries
}
