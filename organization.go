package peeringdb

import (
	"encoding/json"
	"time"
)

// organizationResource is the top-level structure when parsing the JSON output
// from the API. This structure is not used if the Organization JSON object is
// included as a field in another JSON object. This structure is used only if
// the proper namespace is queried.
type organizationResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []Organization `json:"data"`
}

// Organization is a structure representing an Organization. An organization
// can be seen as an enterprise linked to networks, facilities and internet
// exchange points.
type Organization struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	AKA                 string    `json:"aka"`
	NameLong            string    `json:"name_long"`
	Website             string    `json:"website"`
	Notes               string    `json:"notes"`
	Require2FA          bool      `json:"require_2fa"`
	NetworkSet          []int     `json:"net_set"`
	FacilitySet         []int     `json:"fac_set"`
	InternetExchangeSet []int     `json:"ix_set"`
	CarrierSet          []int     `json:"carrier_set"`
	CampusSet           []int     `json:"campus_set"`
	Address1            string    `json:"address1"`
	Address2            string    `json:"address2"`
	City                string    `json:"city"`
	Country             string    `json:"country"`
	State               string    `json:"state"`
	Zipcode             string    `json:"zipcode"`
	Floor               string    `json:"floor"`
	Suite               string    `json:"suite"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
	Status              string    `json:"status"`
	SocialMedia         []SocialMedia `json:"social_media"`
}

// getOrganizationResource returns a pointer to an organizationResource
// structure corresponding to the API JSON response. An error can be returned
// if something went wrong.
func (api *API) getOrganizationResource(search map[string]interface{}) (*organizationResource, error) {
	// Get the OrganizationResource from the API
	response, err := api.lookup(organizationNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &organizationResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetOrganization returns a pointer to a slice of Organization structures that
// the PeeringDB API can provide matching the given search parameters map. If
// an error occurs, the returned error will be non-nil. The returned value can
// be nil if no object could be found.
func (api *API) GetOrganization(search map[string]interface{}) (*[]Organization, error) {
	// Ask for the all Organization objects
	organizationResource, err := api.getOrganizationResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all Organization objects, will be nil if slice is empty
	return &organizationResource.Data, nil
}

// GetAllOrganizations returns a pointer to a slice of Organization structures
// that the PeeringDB API can provide. If an error occurs, the returned error
// will be non-nil. The can be nil if no object could be found.
func (api *API) GetAllOrganizations() (*[]Organization, error) {
	// Return all Organization objects
	return api.GetOrganization(nil)
}

// GetOrganizationByID returns a pointer to a Organization structure that
// matches the given ID. If the ID is lesser than 0, it will return nil. The
// returned error will be non-nil if an issue as occurred while trying to query
// the API. If for some reasons the API returns more than one object for the
// given ID (but it must not) only the first will be used for the returned
// value.
func (api *API) GetOrganizationByID(id int) (*Organization, error) {
	// No point of looking for the organization with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the Organization given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	organizations, err := api.GetOrganization(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No Organization matching the ID
	if len(*organizations) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*organizations)[0], nil
}
