package peeringdb

import (
	"encoding/json"
	"time"
)

// campusResource is the top-level structure when parsing the JSON output
// from the API. This structure is not used if the Campus JSON object is
// included as a field in another JSON object. This structure is used only if
// the proper namespace is queried.
type campusResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []Campus `json:"data"`
}

// Campus is the representation of a site where facilities are.
type Campus struct {
	ID               int          `json:"id"`
	OrganizationID   int          `json:"org_id"`
	OrganizationName string       `json:"org_name"`
	Organization     Organization `json:"organization,omitempty"`
	Name             string       `json:"name"`
	NameLong         string       `json:"name_long"`
	AKA              string       `json:"aka"`
	Website          string       `json:"website"`
	Notes            string       `json:"notes"`
	Created          time.Time    `json:"created"`
	Updated          time.Time    `json:"updated"`
	Status           string       `json:"status"`
	City             string       `json:"city"`
	Country          string       `json:"country"`
	State            string       `json:"state"`
	Zipcode          string       `json:"zipcode"`
	FacilitySet      []int        `json:"fac_set"`
	SocialMedia      []struct {
		Service    string `json:"service"`
		Identifier string `json:"identifier"`
	} `json:"social_media"`
}

// getCampusResource returns a pointer to a campusResource structure
// corresponding to the API JSON response. An error can be returned if
// something went wrong.
func (api *API) getCampusResource(search map[string]interface{}) (*campusResource, error) {
	// Get the CampusResource from the API
	response, err := api.lookup(campusNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &campusResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetCampus returns a pointer to a slice of Campus structures that the
// PeeringDB API can provide matching the given search parameters map. If an
// error occurs, the returned error will be non-nil. The returned value can be
// nil if no object could be found.
func (api *API) GetCampus(search map[string]interface{}) (*[]Campus, error) {
	// Ask for the all Campus objects
	campusResource, err := api.getCampusResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all Campus objects, will be nil if slice is empty
	return &campusResource.Data, nil
}

// GetAllCampuses returns a pointer to a slice of Campus structures that the
// PeeringDB API can provide. If an error occurs, the returned error will be
// non-nil. The can be nil if no object could be found.
func (api *API) GetAllCampuses() (*[]Campus, error) {
	// Return all Campus objects
	return api.GetCampus(nil)
}

// GetCampusByID returns a pointer to a Campus structure that matches the
// given ID. If the ID is lesser than 0, it will return nil. The returned
// error will be non-nil if an issue as occurred while trying to query the
// API. If for some reasons the API returns more than one object for the
// given ID (but it must not) only the first will be used for the returned
// value.
func (api *API) GetCampusByID(id int) (*Campus, error) {
	// No point of looking for the campus with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the Campus given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	campuses, err := api.GetCampus(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No Campus matching the ID
	if len(*campuses) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*campuses)[0], nil
}
