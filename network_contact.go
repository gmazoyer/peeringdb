package peeringdb

import (
	"encoding/json"
	"time"
)

// networkContactResource is the top-level structure when parsing the JSON
// output from the API. This structure is not used if the NetworkContact JSON
// object is included as a field in another JSON object. This structure is used
// only if the proper namespace is queried.
type networkContactResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []NetworkContact `json:"data"`
}

// NetworkContact represents a contact for a network.
type NetworkContact struct {
	ID        int       `json:"id"`
	NetworkID int       `json:"net_id"`
	Network   Network   `json:"net"`
	Role      string    `json:"role"`
	Visible   string    `json:"visible"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	URL       string    `json:"url"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Status    string    `json:"status"`
}

// getNetworkContactResource returns a pointer to an networkContactResource
// structure corresponding to the API JSON response. An error can be returned
// if something went wrong.
func (api *API) getNetworkContactResource(search map[string]interface{}) (*networkContactResource, error) {
	// Get the NetworkContactResource from the API
	response, err := api.lookup(networkContactNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &networkContactResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetNetworkContact returns a pointer to a slice of NetworkContact structures
// that the PeeringDB API can provide matching the given search parameters map.
// If an error occurs, the returned error will be non-nil. The returned value
// can be nil if no object could be found.
func (api *API) GetNetworkContact(search map[string]interface{}) (*[]NetworkContact, error) {
	// Ask for the all NetworkContact objects
	networkContactResource, err := api.getNetworkContactResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all NetworkContact objects, will be nil if slice is empty
	return &networkContactResource.Data, nil
}

// GetAllNetworkContacts returns a pointer to a slice of NetworkContact
// structures that the PeeringDB API can provide. If an error occurs, the
// returned error will be non-nil. The can be nil if no object could be found.
func (api *API) GetAllNetworkContacts() (*[]NetworkContact, error) {
	// Return all NetworkContact objects
	return api.GetNetworkContact(nil)
}

// GetNetworkContactByID returns a pointer to a NetworkContact structure that
// matches the given ID. If the ID is lesser than 0, it will return nil. The
// returned error will be non-nil if an issue as occurred while trying to query
// the API. If for some reasons the API returns more than one object for the
// given ID (but it must not) only the first will be used for the returned
// value.
func (api *API) GetNetworkContactByID(id int) (*NetworkContact, error) {
	// No point of looking for the network contact with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the NetworkContact given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	networkContacts, err := api.GetNetworkContact(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No NetworkContact matching the ID
	if len(*networkContacts) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*networkContacts)[0], nil
}
