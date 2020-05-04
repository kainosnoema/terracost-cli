// Package apiserver is a minimal Terracost API server used for testing
// the prices package interaction with the API's GraphQL interface
package apiserver

import (
	"encoding/json"
	"net/http"
)

// Handler is an implementation of net/http.Handler that provides a stub
// Terracost API server implementation with a "/api/graphql" endpoint:
var Handler http.Handler

type handler struct{}

type dimension struct {
	BeginRange   string
	EndRange     string
	PricePerUnit string
	Unit         string
	RateCode     string
	Description  string
}

type price struct {
	ServiceCode    string
	UsageOperation string
	Dimensions     []dimension
	Updated        string
}

type response struct {
	Prices []price
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/api/graphql":
		h.serveAPIGraphQL(w, req)
	default:
		w.WriteHeader(404)
	}
}

func (h handler) serveAPIGraphQL(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type graphQLRes struct {
		Data response `json:"data"`
	}
	json.NewEncoder(w).Encode(graphQLRes{response{
		[]price{{
			ServiceCode:    "AmazonEC2",
			UsageOperation: "USW2-BoxUsage:t2.micro:RunInstances",
			Dimensions:     []dimension{{"0", "Inf", "0.1", "Hours", "", "Description"}},
			Updated:        "",
		}},
	}})
}

func init() {
	Handler = handler{}
}
