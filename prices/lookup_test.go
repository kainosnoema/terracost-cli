package prices

import (
	"net/http/httptest"
	"testing"

	"github.com/kainosnoema/terracost-cli/prices/testdata/apiserver"
)

func TestLookup(t *testing.T) {
	s := httptest.NewServer(apiserver.Handler)
	defer s.Close()

	apiURL = s.URL + "/api/graphql"
	t.Run("basic lookup", func(t *testing.T) {
		priceID := PriceID{"AmazonEC2", "USW2-BoxUsage:t2.micro:RunInstances"}

		lookup := NewLookup()
		price := lookup.Add(priceID)
		err := lookup.Perform()
		if err != nil {
			t.Fatal(err)
		}
		if len(price.Dimensions) == 0 {
			t.Fatalf("failed to lookup price for %v", priceID)
		}
	})
}
