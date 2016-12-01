package peeringdb

import "fmt"

func Example() {
	api := NewAPI()

	// Look for the organization given a name
	search := make(map[string]interface{})
	search["name"] = "LUXNETWORK S.A."

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
	// Output: LUXNETWORK S.A.
}

func ExampleAPI_GetASN() {
	api := NewAPI()
	as29467 := api.GetASN(29467)

	fmt.Printf("Name:      %s\n", as29467.Name)
	fmt.Printf("AS number: %d\n", as29467.ASN)
	// Output:
	// Name:      LUXNETWORK S.A.
	// AS number: 29467
}
