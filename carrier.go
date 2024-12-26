package peeringdb

import (
	"encoding/json"
	"time"
)

// carrierResource is the top-level structure when parsing the JSON output
// from the API. This structure is not used if the Carrier JSON object is
// included as a field in another JSON object. This structure is used only if
// the proper namespace is queried.
type carrierResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []Carrier `json:"data"`
}

// Carrier is the representation of a network able to provider transport from
// one facility to another.
type Carrier struct {
	ID               int          `json:"id"`
	OrganizationID   int          `json:"org_id"`
	OrganizationName string       `json:"org_name"`
	Organization     Organization `json:"organization,omitempty"`
	Name             string       `json:"name"`
	AKA              string       `json:"aka"`
	NameLong         string       `json:"name_long"`
	Website          string       `json:"website"`
	Notes            string       `json:"notes"`
	Created          time.Time    `json:"created"`
	Updated          time.Time    `json:"updated"`
	Status           string       `json:"status"`
	SocialMedia      []struct {
		Service    string `json:"service"`
		Identifier string `json:"identifier"`
	} `json:"social_media"`
}

// getCarrierResource returns a pointer to a carrierResource structure
// corresponding to the API JSON response. An error can be returned if
// something went wrong.
func (api *API) getCarrierResource(search map[string]interface{}) (*carrierResource, error) {
	// Get the CarrierResource from the API
	response, err := api.lookup(carrierNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &carrierResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetCarrier returns a pointer to a slice of Carrier structures that the
// PeeringDB API can provide matching the given search parameters map. If an
// error occurs, the returned error will be non-nil. The returned value can be
// nil if no object could be found.
func (api *API) GetCarrier(search map[string]interface{}) (*[]Carrier, error) {
	// Ask for the all Carrier objects
	carrierResource, err := api.getCarrierResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all Carrier objects, will be nil if slice is empty
	return &carrierResource.Data, nil
}

// GetAllCarriers returns a pointer to a slice of Carrier structures that the
// PeeringDB API can provide. If an error occurs, the returned error will be
// non-nil. The can be nil if no object could be found.
func (api *API) GetAllCarriers() (*[]Carrier, error) {
	// Return all Carrier objects
	return api.GetCarrier(nil)
}

// GetCarrierByID returns a pointer to a Carrier structure that matches the
// given ID. If the ID is lesser than 0, it will return nil. The returned
// error will be non-nil if an issue as occurred while trying to query the
// API. If for some reasons the API returns more than one object for the
// given ID (but it must not) only the first will be used for the returned
// value.
func (api *API) GetCarrierByID(id int) (*Carrier, error) {
	// No point of looking for the carrier with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the Carrier given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	carriers, err := api.GetCarrier(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No Carrier matching the ID
	if len(*carriers) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*carriers)[0], nil
}

// carrierFacilityResource is the top-level structure when parsing the JSON
// output from the API. This structure is not used if the CarrierFacility JSON
// object is included as a field in another JSON object. This structure is
// used only if the proper namespace is queried.
type carrierFacilityResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []CarrierFacility `json:"data"`
}

// CarrierFacility is a structure used to link an Carrier structure with a
// Facility structure. It helps to know in which facilities a carrier
// operates.
type CarrierFacility struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CarrierID  int       `json:"carrier_id"`
	Carrier    Carrier   `json:"carrier,omitempty"`
	FacilityID int       `json:"fac_id"`
	Facility   Facility  `json:"fac,omitempty"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	Status     string    `json:"status"`
}

// getCarrierFacilityResource returns a pointer to an carrierFacilityResource
// structure corresponding to the API JSON response. An error can be returned
// if something went wrong.
func (api *API) getCarrierFacilityResource(search map[string]interface{}) (*carrierFacilityResource, error) {
	// Get the CarrierFacilityResource from the API
	response, err := api.lookup(carrierFacilityNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &carrierFacilityResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetCarrierFacility returns a pointer to a slice of CarrierFacility
// structures that the PeeringDB API can provide matching the given search
// parameters map. If an error occurs, the returned error will be non-nil. The
// returned value can be nil if no object could be found.
func (api *API) GetCarrierFacility(search map[string]interface{}) (*[]CarrierFacility, error) {
	// Ask for the all CarrierFacility objects
	carrierFacilityResource, err := api.getCarrierFacilityResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all InternetExchangeFacility objects, will be nil if slice is
	// empty
	return &carrierFacilityResource.Data, nil
}

// GetAllCarrierFacilities returns a pointer to a slice of CarrierFacility
// structures that the PeeringDB API can provide. If an error occurs, the
// returned error will be non-nil. The can be nil if no object could be found.
func (api *API) GetAllCarrierFacilities() (*[]CarrierFacility, error) {
	// Return all CarrierFacility objects
	return api.GetCarrierFacility(nil)
}

// GetCarrierFacilityByID returns a pointer to a CarrierFacility structure
// that matches the given ID. If the ID is lesser than 0, it will return nil.
// The returned error will be non-nil if an issue as occurred while trying to
// query the API. If for some reasons the API returns more than one object for
// the given ID (but it must not) only the first will be used for the returned
// value.
func (api *API) GetCarrierFacilityByID(id int) (*CarrierFacility, error) {
	// No point of looking for the carrier facility with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the CarrierFacility given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	carrierFacilities, err := api.GetCarrierFacility(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No CarrierFacility matching the ID
	if len(*carrierFacilities) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*carrierFacilities)[0], nil
}
