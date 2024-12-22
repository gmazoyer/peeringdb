package peeringdb

import (
	"encoding/json"
	"time"
)

// facilityResource is the top-level structure when parsing the JSON output
// from the API. This structure is not used if the Facility JSON object is
// included as a field in another JSON object. This structure is used only if
// the proper namespace is queried.
type facilityResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []Facility `json:"data"`
}

// Facility is the representation of a location where network operators and
// Internet exchange points are located. Most of the time you know a facility
// as a datacenter.
type Facility struct {
	ID                        int               `json:"id"`
	OrganizationID            int               `json:"org_id"`
	OrganizationName          string            `json:"org_name"`
	Organization              Organization      `json:"organization,omitempty"`
	Name                      string            `json:"name"`
	AKA                       string            `json:"aka"`
	NameLong                  string            `json:"name_long"`
	Website                   string            `json:"website"`
	SocialMedia               []SocialMediaItem `json:"social_media"`
	CLLI                      string            `json:"clli"`
	Rencode                   string            `json:"rencode"`
	Npanxx                    string            `json:"npanxx"`
	Notes                     string            `json:"notes"`
	NetCount                  int               `json:"net_count"`
	IXCount                   int               `json:"ix_count"`
	SalesEmail                string            `json:"sales_email"`
	SalesPhone                string            `json:"sales_phone"`
	TechEmail                 string            `json:"tech_email"`
	TechPhone                 string            `json:"tech_phone"`
	AvailableVoltageServices  []string          `json:"available_voltage_services"`
	DiverseServingSubstations bool              `json:"diverse_serving_substations"`
	Property                  string            `json:"property"`
	RegionContinent           string            `json:"region_continent"`
	StatusDashboard           string            `json:"status_dashboard"`
	Created                   time.Time         `json:"created"`
	Updated                   time.Time         `json:"updated"`
	Status                    string            `json:"status"`
	Address1                  string            `json:"address1"`
	Address2                  string            `json:"address2"`
	City                      string            `json:"city"`
	Country                   string            `json:"country"`
	State                     string            `json:"state"`
	Zipcode                   string            `json:"zipcode"`
	Floor                     string            `json:"floor"`
	Suite                     string            `json:"suite"`
	Latitude                  float64           `json:"latitude"`
	Longitude                 float64           `json:"longitude"`
}

// getFacilityResource returns a pointer to a facilityResource structure
// corresponding to the API JSON response. An error can be returned if
// something went wrong.
func (api *API) getFacilityResource(search map[string]interface{}) (*facilityResource, error) {
	// Get the FacilityResource from the API
	response, err := api.lookup(facilityNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &facilityResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetFacility returns a pointer to a slice of Facility structures that the
// PeeringDB API can provide matching the given search parameters map. If an
// error occurs, the returned error will be non-nil. The returned value can be
// nil if no object could be found.
func (api *API) GetFacility(search map[string]interface{}) (*[]Facility, error) {
	// Ask for the all Facility objects
	facilyResource, err := api.getFacilityResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all Facility objects, will be nil if slice is empty
	return &facilyResource.Data, nil
}

// GetAllFacilities returns a pointer to a slice of Facility structures that
// the PeeringDB API can provide. If an error occurs, the returned error will
// be non-nil. The can be nil if no object could be found.
func (api *API) GetAllFacilities() (*[]Facility, error) {
	// Return all Facility objects
	return api.GetFacility(nil)
}

// GetFacilityByID returns a pointer to a Facility structure that matches the
// given ID. If the ID is lesser than 0, it will return nil. The returned error
// will be non-nil if an issue as occurred while trying to query the API. If for
// some reasons the API returns more than one object for the given ID (but it
// must not) only the first will be used for the returned value.
func (api *API) GetFacilityByID(id int) (*Facility, error) {
	// No point of looking for the facility with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the Facility given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	facilities, err := api.GetFacility(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No Facility matching the ID
	if len(*facilities) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*facilities)[0], nil
}
