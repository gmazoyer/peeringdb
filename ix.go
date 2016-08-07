package peeringdb

import (
	"encoding/json"
	"time"
)

// InternetExchangeResource is the top-level structure when parsing the JSON
// output from the API. This structure is not used if the InternetExchange JSON
// object is included as a field in another JSON object. This structure is used
// only if the proper namespace is queried.
type InternetExchangeResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []InternetExchange `json:"data"`
}

// InternetExchange is a structure representing an Internet exchange point. It
// is directly linked to the Organization that manage the IX.
type InternetExchange struct {
	ID                     int          `json:"id"`
	OrganizationID         int          `json:"org_id"`
	Organization           Organization `json:"org,omitempty"`
	Name                   string       `json:"name"`
	NameLong               string       `json:"name_long"`
	City                   string       `json:"city"`
	Country                string       `json:"country"`
	RegionContinent        string       `json:"region_continent"`
	Media                  string       `json:"media"`
	Notes                  string       `json:"notes"`
	ProtoUnicast           bool         `json:"proto_unicast"`
	ProtoMulticast         bool         `json:"proto_multicast"`
	ProtoIPv6              bool         `json:"proto_ipv6"`
	Website                string       `json:"website"`
	URLStats               string       `json:"url_stats"`
	TechEmail              string       `json:"tech_email"`
	TechPhone              string       `json:"tech_phone"`
	PolicyEmail            string       `json:"policy_email"`
	PolicyPhone            string       `json:"policy_phone"`
	FacilitySet            []int        `json:"fac_set"`
	InternetExchangeLANSet []int        `json:"ixlan_set"`
	Created                time.Time    `json:"created"`
	Updated                time.Time    `json:"updated"`
	Status                 string       `json:"status"`
}

// getInternetExchangeResource returns a pointer to an InternetExchangeResource
// structure corresponding to the API JSON response. An error can be returned
// if something went wrong.
func getInternetExchangeResource(search map[string]interface{}) (*InternetExchangeResource, error) {
	// Get the InternetExchangeResource from the API
	response, err := lookup(internetExchangeNamespace, nil, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &InternetExchangeResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetInternetExchange returns a pointer to a slice of InternetExchange
// structures that the PeeringDB API can provide matching the given search
// parameters map. If an error occurs, the returned error will be non-nil. The
// returned value can be nil if no object could be found.
func GetInternetExchange(search map[string]interface{}) (*[]InternetExchange, error) {
	// Ask for the all InternetExchange objects
	internetExchangeResource, err := getInternetExchangeResource(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// Return all InternetExchange objects, will be nil if slice is empty
	return &internetExchangeResource.Data, nil
}

// GetAllInternetExchanges returns a pointer to a slice of InternetExchange
// structures that the PeeringDB API can provide. If an error occurs, the
// returned error will be non-nil. The can be nil if no object could be found.
func GetAllInternetExchanges() (*[]InternetExchange, error) {
	// Return all InternetExchange objects
	return GetInternetExchange(nil)
}

// GetInternetExchangeByID returns a pointer to a InternetExchange structure
// that matches the given ID. If the ID is lesser than 0, it will return nil.
// The returned error will be non-nil if an issue as occured while trying to
// query the API. If for some reasons the API returns more than one object for
// the given ID (but it must not) only the first will be used for the returned
// value.
func GetInternetExchangeByID(id int) (*InternetExchange, error) {
	// No point of looking for the Internet exchange with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the InternetExchange given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	internetExchanges, err := GetInternetExchange(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// No InternetExchange matching the ID
	if len(*internetExchanges) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*internetExchanges)[0], nil
}

// InternetExchangeLANResource is the top-level structure when parsing the JSON
// output from the API. This structure is not used if the InternetExchangeLAN
// JSON object is included as a field in another JSON object. This structure is
// used only if the proper namespace is queried.
type InternetExchangeLANResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []InternetExchangeLAN `json:"data"`
}

// InternetExchangeLAN is a structure representing the one of the network (LAN)
// of an Internet exchange points. It contains details about the LAN like the
// MTU, VLAN support, etc.
type InternetExchangeLAN struct {
	ID                        int              `json:"id"`
	InternetExchangeID        int              `json:"ix_id"`
	InternetExchange          InternetExchange `json:"ix,omitempty"`
	Name                      string           `json:"name"`
	Description               string           `json:"descr"`
	MTU                       int              `json:"mtu"`
	Dot1QSupport              bool             `json:"dot1q_support"`
	RouteServerASN            int              `json:"rs_asn"`
	ARPSponge                 string           `json:"arp_sponge"`
	NetworkSet                []int            `json:"net_set"`
	InternetExchangePrefixSet []int            `json:"ixpfx_set"`
	Status                    string           `json:"status"`
	Updated                   time.Time        `json:"updated"`
	Created                   time.Time        `json:"created"`
}

// getInternetExchangeLANResource returns a pointer to an
// InternetExchangeLANResource structure corresponding to the API JSON
// response. An error can be returned if  something went wrong.
func getInternetExchangeLANResource(search map[string]interface{}) (*InternetExchangeLANResource, error) {
	// Get the InternetExchangeLANResource from the API
	response, err := lookup(internetExchangeLANNamespace, nil, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &InternetExchangeLANResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetInternetExchangeLAN returns a pointer to a slice of InternetExchangeLAN
// structures that the PeeringDB API can provide matching the given search
// parameters map. If an error occurs, the returned error will be non-nil. The
// returned value can be nil if no object could be found.
func GetInternetExchangeLAN(search map[string]interface{}) (*[]InternetExchangeLAN, error) {
	// Ask for the all InternetExchangeLAN objects
	internetExchangeLANResource, err := getInternetExchangeLANResource(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// Return all InternetExchangeLAN objects, will be nil if slice is empty
	return &internetExchangeLANResource.Data, nil
}

// GetAllInternetExchangeLANs returns a pointer to a slice of
// InternetExchangeLAN structures that the PeeringDB API can provide. If an
// error occurs, the returned error will be non-nil. The can be nil if no
// object could be found.
func GetAllInternetExchangeLANs() (*[]InternetExchangeLAN, error) {
	// Return all InternetExchangeLAN objects
	return GetInternetExchangeLAN(nil)
}

// GetInternetExchangeLANByID returns a pointer to a InternetExchangeLAN
// structure that matches the given ID. If the ID is lesser than 0, it will
// return nil. The returned error will be non-nil if an issue as occured while
// trying to query the API. If for some reasons the API returns more than one
// object for the given ID (but it must not) only the first will be used for
// the returned value.
func GetInternetExchangeLANByID(id int) (*InternetExchangeLAN, error) {
	// No point of looking for the Internet exchange LAN with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the InternetExchangeLAN given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	ixLANs, err := GetInternetExchangeLAN(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// No InternetExchangeLAN matching the ID
	if len(*ixLANs) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*ixLANs)[0], nil
}

// InternetExchangePrefixResource is the top-level structure when parsing the
// JSON output from the API. This structure is not used if the
// InternetExchangePrefix JSON object is included as a field in another JSON
// object. This structure is used only if the proper namespace is queried.
type InternetExchangePrefixResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []InternetExchangePrefix `json:"data"`
}

// InternetExchangePrefix is a structure representing the prefix used by an
// Internet exchange point. It is directly linked to an InternetExchangeLAN.
type InternetExchangePrefix struct {
	ID                    int                 `json:"id"`
	InternetExchangeLANID int                 `json:"ixlan_id"`
	InternetExchangeLAN   InternetExchangeLAN `json:"ixlan,omitempty"`
	Protocol              string              `json:"protocol"`
	Prefix                string              `json:"prefix"`
	Created               time.Time           `json:"created"`
	Updated               time.Time           `json:"updated"`
	Status                string              `json:"status"`
}

// getInternetExchangePrefixResource returns a pointer to an
// InternetExchangePrefixResource structure corresponding to the API JSON
// response. An error can be returned if something went wrong.
func getInternetExchangePrefixResource(search map[string]interface{}) (*InternetExchangePrefixResource, error) {
	// Get the InternetExchangePrefixResource from the API
	response, err := lookup(internetExchangePrefixNamespace, nil, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &InternetExchangePrefixResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetInternetExchangePrefix returns a pointer to a slice of
// InternetExchangePrefix structures that the PeeringDB API can provide
// matching the given search parameters map. If an error occurs, the returned
// error will be non-nil. The returned value can be nil if no object could be
// found.
func GetInternetExchangePrefix(search map[string]interface{}) (*[]InternetExchangePrefix, error) {
	// Ask for the all InternetExchangePrefix objects
	internetExchangePrefixResource, err := getInternetExchangePrefixResource(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// Return all InternetExchangePrefix objects, will be nil if slice is empty
	return &internetExchangePrefixResource.Data, nil
}

// GetAllInternetExchangePrefixes returns a pointer to a slice of
// InternetExchangePrefix structures that the PeeringDB API can provide. If an
// error occurs, the returned error will be non-nil. The can be nil if no
// object could be found.
func GetAllInternetExchangePrefixes() (*[]InternetExchangePrefix, error) {
	// Return all InternetExchangePrefix objects
	return GetInternetExchangePrefix(nil)
}

// GetInternetExchangePrefixByID returns a pointer to a InternetExchangePrefix
// structure that matches the given ID. If the ID is lesser than 0, it will
// return nil. The returned error will be non-nil if an issue as occured while
// trying to query the API. If for some reasons the API returns more than one
// object for the given ID (but it must not) only the first will be used for
// the returned value.
func GetInternetExchangePrefixByID(id int) (*InternetExchangePrefix, error) {
	// No point of looking for the Internet exchange prefix with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the InternetExchangePrefix given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	ixPrefixes, err := GetInternetExchangePrefix(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// No InternetExchangePrefix matching the ID
	if len(*ixPrefixes) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*ixPrefixes)[0], nil
}

// InternetExchangeFacilityResource is the top-level structure when parsing the
// JSON output from the API. This structure is not used if the
// InternetExchangeFacility JSON object is included as a field in another JSON
// object. This structure is used only if the proper namespace is queried.
type InternetExchangeFacilityResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []InternetExchangeFacility `json:"data"`
}

// InternetExchangeFacility is a structure used to link an InternetExchange
// structure with a Facility structure. It helps to know where an Internet
// exchange points can be found, or what Internet exchange points can be found
// in a given facility.
type InternetExchangeFacility struct {
	ID                 int              `json:"id"`
	InternetExchangeID int              `json:"ix_id"`
	InternetExchange   InternetExchange `json:"ix,omitempty"`
	FacilityID         int              `json:"fac_id"`
	Facility           Facility         `json:"fac,omitempty"`
	Created            time.Time        `json:"created"`
	Updated            time.Time        `json:"updated"`
	Status             string           `json:"status"`
}

// getInternetExchangeFacilityResource returns a pointer to an
// InternetExchangeFacilityResource structure corresponding to the API JSON
// response. An error can be returned if something went wrong.
func getInternetExchangeFacilityResource(search map[string]interface{}) (*InternetExchangeFacilityResource, error) {
	// Get the InternetExchangeFacilityResource from the API
	response, err := lookup(internetExchangeFacilityNamespace, nil, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &InternetExchangeFacilityResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetInternetExchangeFacility returns a pointer to a slice of
// InternetExchangeFacility structures that the PeeringDB API can provide
// matching the given search parameters map. If an error occurs, the returned
// error will be non-nil. The returned value can be nil if no object could be
// found.
func GetInternetExchangeFacility(search map[string]interface{}) (*[]InternetExchangeFacility, error) {
	// Ask for the all InternetExchangeFacility objects
	internetExchangeFacilityResource, err := getInternetExchangeFacilityResource(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// Return all InternetExchangeFacility objects, will be nil if slice is
	// empty
	return &internetExchangeFacilityResource.Data, nil
}

// GetAllInternetExchangeFacilities returns a pointer to a slice of
// InternetExchangeFacility structures that the PeeringDB API can provide. If
// an error occurs, the returned error will be non-nil. The can be nil if no
// object could be found.
func GetAllInternetExchangeFacilities() (*[]InternetExchangeFacility, error) {
	// Return all InternetExchangeFacility objects
	return GetInternetExchangeFacility(nil)
}

// GetInternetExchangeFacilityByID returns a pointer to a
// InternetExchangeFacility structure that matches the given ID. If the ID is
// lesser than 0, it will return nil. The returned error will be non-nil if an
// issue as occured while trying to query the API. If for some reasons the API
// returns more than one object for the given ID (but it must not) only the
// first will be used for the returned value.
func GetInternetExchangeFacilityByID(id int) (*InternetExchangeFacility, error) {
	// No point of looking for the Internet exchange facility with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the InternetExchangeFacility given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	ixFacilities, err := GetInternetExchangeFacility(search)

	// Error as occured while querying the API
	if err != nil {
		return nil, err
	}

	// No InternetExchangeFacility matching the ID
	if len(*ixFacilities) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*ixFacilities)[0], nil
}
