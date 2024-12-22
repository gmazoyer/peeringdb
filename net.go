package peeringdb

import (
	"encoding/json"
	"time"
)

// networkResource is the top-level structure when parsing the JSON output from
// the API. This structure is not used if the Network JSON object is included
// as a field in another JSON object. This structure is used only if the proper
// namespace is queried.
type networkResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []Network `json:"data"`
}

// Network is a structure representing a network. Basically, a network is an
// Autonomous System identified by an AS number and other details. It belongs
// to an Organization, contains one or more NetworkContact, and is part of
// several Facility and InternetExchangeLAN.
type Network struct {
	ID                                int               `json:"id"`
	OrganizationID                    int               `json:"org_id"`
	Organization                      Organization      `json:"org,omitempty"`
	Name                              string            `json:"name"`
	AKA                               string            `json:"aka"`
	NameLong                          string            `json:"name_long"`
	Website                           string            `json:"website"`
	SocialMedia                       []SocialMediaItem `json:"social_media"`
	ASN                               int               `json:"asn"`
	LookingGlass                      string            `json:"looking_glass"`
	RouteServer                       string            `json:"route_server"`
	IRRASSet                          string            `json:"irr_as_set"`
	InfoType                          string            `json:"info_type"`
	InfoTypes                         []string          `json:"info_types"`
	InfoPrefixes4                     int               `json:"info_prefixes4"`
	InfoPrefixes6                     int               `json:"info_prefixes6"`
	InfoTraffic                       string            `json:"info_traffic"`
	InfoRatio                         string            `json:"info_ratio"`
	InfoScope                         string            `json:"info_scope"`
	InfoUnicast                       bool              `json:"info_unicast"`
	InfoMulticast                     bool              `json:"info_multicast"`
	InfoIPv6                          bool              `json:"info_ipv6"`
	InfoNeverViaRouteServers          bool              `json:"info_never_via_route_servers"`
	IXCount                           int               `json:"ix_count"`
	FacCount                          int               `json:"fac_count"`
	Notes                             string            `json:"notes"`
	NetworkInternetExchangeLANUpdated time.Time         `json:"netixlan_updated"`
	NetworkFacilityUpdated            time.Time         `json:"netfac_updated"`
	NetworkContactUpdated             time.Time         `json:"poc_updated"`
	PolicyURL                         string            `json:"policy_url"`
	PolicyGeneral                     string            `json:"policy_general"`
	PolicyLocations                   string            `json:"policy_locations"`
	PolicyRatio                       bool              `json:"policy_ratio"`
	PolicyContracts                   string            `json:"policy_contracts"`
	NetworkFacilitySet                []int             `json:"netfac_set"`
	NetworkInternetExchangeLANSet     []int             `json:"netixlan_set"`
	NetworkContactSet                 []int             `json:"poc_set"`
	AllowIXPUpdate                    bool              `json:"allow_ixp_update"`
	StatusDashboard                   string            `json:"status_dashboard"`
	RIRStatus                         string            `json:"rir_status"`
	RIRStatusUpdated                  time.Time         `json:"rir_status_updated"`
	Created                           time.Time         `json:"created"`
	Updated                           time.Time         `json:"updated"`
	Status                            string            `json:"status"`
}

// getNetworkResource returns a pointer to an networkResource structure
// corresponding to the API JSON response. An error can be returned if
// something went wrong.
func (api *API) getNetworkResource(search map[string]interface{}) (*networkResource, error) {
	// Get the NetworkResource from the API
	response, err := api.lookup(networkNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &networkResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetNetwork returns a pointer to a slice of Network structures that the
// PeeringDB API can provide matching the given search parameters map. If an
// error occurs, the returned error will be non-nil. The returned value can be
// nil if no object could be found.
func (api *API) GetNetwork(search map[string]interface{}) (*[]Network, error) {
	// Ask for the all Network objects
	networkResource, err := api.getNetworkResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all Network objects, will be nil if slice is empty
	return &networkResource.Data, nil
}

// GetAllNetworks returns a pointer to a slice of Network structures that the
// PeeringDB API can provide. If an error occurs, the returned error will be
// non-nil. The can be nil if no object could be found.
func (api *API) GetAllNetworks() (*[]Network, error) {
	// Return all Network objects
	return api.GetNetwork(nil)
}

// GetNetworkByID returns a pointer to a Network structure that matches the
// given ID. If the ID is lesser than 0, it will return nil. The returned error
// will be non-nil if an issue as occurred while trying to query the API. If for
// some reasons the API returns more than one object for the given ID (but it
// must not) only the first will be used for the returned value.
func (api *API) GetNetworkByID(id int) (*Network, error) {
	// No point of looking for the network with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the Network given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	networks, err := api.GetNetwork(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No Network matching the ID
	if len(*networks) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*networks)[0], nil
}

// networkFacilityResource is the top-level structure when parsing the JSON
// output from the API. This structure is not used if the NetFacility JSON
// object is included as a field in another JSON object. This structure is used
// only if the proper namespace is queried.
type networkFacilityResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []NetworkFacility `json:"data"`
}

// NetworkFacility is a structure used to link a Network with a Facility. It
// helps to know where a network is located (it can be in several facilities).
// For example, it can be used to search common facilities between several
// networks to know where they can interconnect themselves directly.
type NetworkFacility struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	NetworkID  int       `json:"net_id"`
	Network    Network   `json:"net,omitempty"`
	FacilityID int       `json:"fac_id"`
	Facility   Facility  `json:"fac,omitempty"`
	LocalASN   int       `json:"local_asn"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	Status     string    `json:"status"`
}

// getNetworkFacilityResource returns a pointer to an networkFacilityResource
// structure corresponding to the API JSON response. An error can be returned
// if something went wrong.
func (api *API) getNetworkFacilityResource(search map[string]interface{}) (*networkFacilityResource, error) {
	// Get the NetworkFacilityResource from the API
	response, err := api.lookup(networkFacilityNamespace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &networkFacilityResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetNetworkFacility returns a pointer to a slice of NetworkFacility
// structures that the PeeringDB API can provide matching the given search
// parameters map. If an error occurs, the returned error will be non-nil. The
// returned value can be nil if no object could be found.
func (api *API) GetNetworkFacility(search map[string]interface{}) (*[]NetworkFacility, error) {
	// Ask for the all NetworkFacility objects
	networkFacilityResource, err := api.getNetworkFacilityResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all NetworkFacility objects, will be nil if slice is empty
	return &networkFacilityResource.Data, nil
}

// GetAllNetworkFacilities returns a pointer to a slice of NetworkFacility
// structures that the PeeringDB API can provide. If an error occurs, the
// returned error will be non-nil. The can be nil if no object could be found.
func (api *API) GetAllNetworkFacilities() (*[]NetworkFacility, error) {
	// Return all NetFacility objects
	return api.GetNetworkFacility(nil)
}

// GetNetworkFacilityByID returns a pointer to a NetworkFacility structure that
// matches the given ID. If the ID is lesser than 0, it will return nil. The
// returned error will be non-nil if an issue as occurred while trying to query
// the API. If for some reasons the API returns more than one object for the
// given ID (but it must not) only the first will be used for the returned
// value.
func (api *API) GetNetworkFacilityByID(id int) (*NetworkFacility, error) {
	// No point of looking for the network facility with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the NetworkFacility given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	networkFacilities, err := api.GetNetworkFacility(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No NetworkFacility matching the ID
	if len(*networkFacilities) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*networkFacilities)[0], nil
}

// networkInternetExchangeLANResource is the top-level structure when parsing
// the JSON output from the API. This structure is not used if the
// NetworkInternetExchangeLAN JSON object is included as a field in another
// JSON object. This structure is used only if the proper namespace is queried.
type networkInternetExchangeLANResource struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []NetworkInternetExchangeLAN `json:"data"`
}

// NetworkInternetExchangeLAN is a structure allowing to know to which
// InternetExchangeLAN a network is connected. It can be used, for example, to
// know what are the common Internet exchange LANs between several networks.
type NetworkInternetExchangeLAN struct {
	ID                     int                 `json:"id"`
	NetworkID              int                 `json:"net_id"`
	Network                Network             `json:"net,omitempty"`
	InternetExchangeID     int                 `json:"ix_id"`
	InternetExchange       InternetExchange    `json:"ix,omitempty"`
	Name                   string              `json:"name"`
	InternetExchangeLANID  int                 `json:"ixlan_id"`
	InternetExchangeLAN    InternetExchangeLAN `json:"ixlan,omitempty"`
	Notes                  string              `json:"notes"`
	Speed                  int                 `json:"speed"`
	ASN                    int                 `json:"asn"`
	IPAddr4                string              `json:"ipaddr4"`
	IPAddr6                string              `json:"ipaddr6"`
	IsRSPeer               bool                `json:"is_rs_peer"`
	Operational            bool                `json:"operational"`
	NetworkSideID          int                 `json:"net_side_id"`
	InternetExchangeSideID int                 `json:"ix_side_id"`
	Created                time.Time           `json:"created"`
	Updated                time.Time           `json:"updated"`
	Status                 string              `json:"status"`
}

// getNetworkInternetExchangeLANResource returns a pointer to an
// networkInternetExchangeLANResource structure corresponding to the API JSON
// response. An error can be returned if something went wrong.
func (api *API) getNetworkInternetExchangeLANResource(search map[string]interface{}) (*networkInternetExchangeLANResource, error) {
	// Get the NetworkInternetExchangeLANResource from the API
	response, err := api.lookup(networkInternetExchangeLANNamepsace, search)
	if err != nil {
		return nil, err
	}

	// Ask for cleanup once we are done
	defer response.Body.Close()

	// Decode what the API has given to us
	resource := &networkInternetExchangeLANResource{}
	err = json.NewDecoder(response.Body).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetNetworkInternetExchangeLAN returns a pointer to a slice of
// NetworkInternetExchangeLAN structures that the PeeringDB API can provide
// matching the given search parameters map. If an error occurs, the returned
// error will be non-nil. The returned value can be nil if no object could be
// found.
func (api *API) GetNetworkInternetExchangeLAN(search map[string]interface{}) (*[]NetworkInternetExchangeLAN, error) {
	// Ask for the all NetInternetExchangeLAN objects
	networkInternetExchangeLANResource, err := api.getNetworkInternetExchangeLANResource(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// Return all NetInternetExchangeLAN objects, will be nil if slice is empty
	return &networkInternetExchangeLANResource.Data, nil
}

// GetAllNetworkInternetExchangeLANs returns a pointer to a slice of
// NetworkInternetExchangeLAN structures that the PeeringDB API can provide. If
// an error occurs, the returned error will be non-nil. The can be nil if no
// object could be found.
func (api *API) GetAllNetworkInternetExchangeLANs() (*[]NetworkInternetExchangeLAN, error) {
	// Return all NetworkInternetExchangeLAN objects
	return api.GetNetworkInternetExchangeLAN(nil)
}

// GetNetworkInternetExchangeLANByID returns a pointer to a
// NetworkInternetExchangeLAN structure that matches the given ID. If the ID is
// lesser than 0, it will return nil. The returned error will be non-nil if an
// issue as occurred while trying to query the API. If for some reasons the API
// returns more than one object for the given ID (but it must not) only the
// first will be used for the returned value.
func (api *API) GetNetworkInternetExchangeLANByID(id int) (*NetworkInternetExchangeLAN, error) {
	// No point of looking for the Internet exchange LAN with an ID < 0
	if id < 0 {
		return nil, nil
	}

	// Ask for the NetworkInternetExchangeLAN given it ID
	search := make(map[string]interface{})
	search["id"] = id

	// Actually ask for it
	networkInternetExchangeLANs, err := api.GetNetworkInternetExchangeLAN(search)

	// Error as occurred while querying the API
	if err != nil {
		return nil, err
	}

	// No NetworkInternetExchangeLAN matching the ID
	if len(*networkInternetExchangeLANs) < 1 {
		return nil, nil
	}

	// Only return the first match, they must be only one match (ID being
	// unique)
	return &(*networkInternetExchangeLANs)[0], nil
}
