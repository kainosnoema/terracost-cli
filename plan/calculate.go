package plan

import (
	"github.com/kainosnoema/terracost/cli/prices"
	"github.com/kainosnoema/terracost/cli/terraform"
)

// Resource maps a Terraform resource to AWS pricing
type Resource struct {
	Address string
	Action  string
	Before  prices.ByID
	After   prices.ByID
}

// Calculate takes a TF plan, fetches AWS prices, and returns priced Resources
func Calculate(tfPlan *terraform.PlanJSON) ([]Resource, error) {
	resources := []Resource{}
	priceLookup := prices.NewLookup()

	for _, res := range tfPlan.ResourceChanges {
		resource := Resource{
			Address: res.Address,
			Action:  res.Change.Actions[0],
			Before:  prices.ByID{},
			After:   prices.ByID{},
		}

		changesPriceIDs := prices.ResourceChangesPriceIDs(tfPlan.Region(), res)
		for _, beforePriceID := range changesPriceIDs.Before {
			resource.Before[beforePriceID] = priceLookup.Add(beforePriceID)
		}
		for _, afterPriceID := range changesPriceIDs.After {
			resource.After[afterPriceID] = priceLookup.Add(afterPriceID)
		}

		resources = append(resources, resource)
	}

	err := priceLookup.Perform()
	if err != nil {
		return nil, err
	}

	return resources, nil
}
