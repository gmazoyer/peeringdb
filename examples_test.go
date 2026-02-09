package peeringdb

import (
	"context"
	"fmt"
	"net/url"
)

func Example() {
	api := NewAPI()
	ctx := context.Background()

	// Look for the organization given a name
	search := url.Values{}
	search.Set("name", "Guillaume Mazoyer")

	// Get the organization
	organizations, err := api.GetOrganization(ctx, search)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(organizations) < 1 {
		fmt.Printf("No organization found with name '%s'\n", search.Get("name"))
		return
	}

	if len(organizations) > 1 {
		fmt.Printf("More than one organizations found with name '%s'\n",
			search.Get("name"))
		return
	}

	org := organizations[0]

	// Find if there are networks linked to the organization
	for _, networkID := range org.NetworkSet {
		network, err := api.GetNetworkByID(ctx, networkID)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(network.Name)
		}
	}
}

func ExampleAPI_GetASN() {
	api := NewAPI()
	as201281, err := api.GetASN(context.Background(), 201281)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Name: %s\n", as201281.Name)
	fmt.Printf("ASN:  %d\n", as201281.ASN)
}
