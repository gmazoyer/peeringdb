package peeringdb

import "fmt"

func Example() {
	api := NewAPI()

	// Look for the organization given a name
	search := make(map[string]interface{})
	search["name"] = "Guillaume Mazoyer"

	// Get the organization, pointer to slice returned
	organizations, err := api.GetOrganization(search)

	// If an error as occurred, print it
	if err != nil {
		fmt.Println(err)
		return
	}

	// No organization found
	if len(*organizations) < 1 {
		fmt.Printf("No organization found with name '%s'\n", search["name"])
		return
	}

	// Several organizations found
	if len(*organizations) > 1 {
		fmt.Printf("More than one organizations found with name '%s'\n",
			search["name"])
		return
	}

	// Get the first found organization
	org := (*organizations)[0]

	// Find if there are networks linked to the organization
	if len(org.NetworkSet) > 0 {
		// For each network
		for _, networkID := range org.NetworkSet {
			// Get the details and print it
			network, err := api.GetNetworkByID(networkID)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Print(network.Name)
			}
		}
	}
	// Output: Guillaume Mazoyer
}

func ExampleAPI_GetASN() {
	api := NewAPI()
	as201281, err := api.GetASN(201281)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Name: %s\n", as201281.Name)
	fmt.Printf("ASN:  %d\n", as201281.ASN)
	// Output:
	// Name: Guillaume Mazoyer
	// ASN:  201281
}
