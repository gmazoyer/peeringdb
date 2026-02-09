package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// Network represents an Autonomous System identified by an AS number and other
// details. It belongs to an Organization, contains one or more
// NetworkContact, and is part of several Facility and InternetExchangeLAN.
type Network struct {
	ID                                int           `json:"id"`
	OrganizationID                    int           `json:"org_id"`
	Organization                      Organization  `json:"org,omitempty"`
	Name                              string        `json:"name"`
	AKA                               string        `json:"aka"`
	NameLong                          string        `json:"name_long"`
	Website                           string        `json:"website"`
	ASN                               int           `json:"asn"`
	LookingGlass                      string        `json:"looking_glass"`
	RouteServer                       string        `json:"route_server"`
	IRRASSet                          string        `json:"irr_as_set"`
	InfoType                          string        `json:"info_type"`
	InfoTypes                         []string      `json:"info_types"`
	InfoPrefixes4                     int           `json:"info_prefixes4"`
	InfoPrefixes6                     int           `json:"info_prefixes6"`
	InfoTraffic                       string        `json:"info_traffic"`
	InfoRatio                         string        `json:"info_ratio"`
	InfoScope                         string        `json:"info_scope"`
	InfoUnicast                       bool          `json:"info_unicast"`
	InfoMulticast                     bool          `json:"info_multicast"`
	InfoIPv6                          bool          `json:"info_ipv6"`
	InfoNeverViaRouteServers          bool          `json:"info_never_via_route_servers"`
	InternetExchangeCount             int           `json:"ix_count"`
	FacilityCount                     int           `json:"fac_count"`
	Notes                             string        `json:"notes"`
	NetworkInternetExchangeLANUpdated time.Time     `json:"netixlan_updated"`
	NetworkFacilityUpdated            time.Time     `json:"netfac_updated"`
	NetworkContactUpdated             time.Time     `json:"poc_updated"`
	PolicyURL                         string        `json:"policy_url"`
	PolicyGeneral                     string        `json:"policy_general"`
	PolicyLocations                   string        `json:"policy_locations"`
	PolicyRatio                       bool          `json:"policy_ratio"`
	PolicyContracts                   string        `json:"policy_contracts"`
	NetworkFacilitySet                []int         `json:"netfac_set"`
	NetworkInternetExchangeLANSet     []int         `json:"netixlan_set"`
	NetworkContactSet                 []int         `json:"poc_set"`
	AllowIXPUpdate                    bool          `json:"allow_ixp_update"`
	StatusDashboard                   string        `json:"status_dashboard"`
	RIRStatus                         string        `json:"rir_status"`
	RIRStatusUpdated                  time.Time     `json:"rir_status_updated"`
	Created                           time.Time     `json:"created"`
	Updated                           time.Time     `json:"updated"`
	Status                            string        `json:"status"`
	SocialMedia                       []SocialMedia `json:"social_media"`
}

// GetNetwork returns a slice of Network structures matching the given search
// parameters.
func (api *API) GetNetwork(ctx context.Context, search url.Values) ([]Network, error) {
	return fetch[Network](ctx, api, networkNamespace, search)
}

// GetAllNetworks returns all Network structures available from the API.
func (api *API) GetAllNetworks(ctx context.Context) ([]Network, error) {
	return api.GetNetwork(ctx, nil)
}

// GetNetworkByID returns the Network matching the given ID, or nil if not
// found.
func (api *API) GetNetworkByID(ctx context.Context, id int) (*Network, error) {
	return fetchByID[Network](ctx, api, networkNamespace, id)
}

// NetworkFacility links a Network with a Facility, indicating where a network
// is located. It can be used to search common facilities between several
// networks to know where they can interconnect directly.
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

// GetNetworkFacility returns a slice of NetworkFacility structures matching
// the given search parameters.
func (api *API) GetNetworkFacility(ctx context.Context, search url.Values) ([]NetworkFacility, error) {
	return fetch[NetworkFacility](ctx, api, networkFacilityNamespace, search)
}

// GetAllNetworkFacilities returns all NetworkFacility structures available
// from the API.
func (api *API) GetAllNetworkFacilities(ctx context.Context) ([]NetworkFacility, error) {
	return api.GetNetworkFacility(ctx, nil)
}

// GetNetworkFacilityByID returns the NetworkFacility matching the given ID,
// or nil if not found.
func (api *API) GetNetworkFacilityByID(ctx context.Context, id int) (*NetworkFacility, error) {
	return fetchByID[NetworkFacility](ctx, api, networkFacilityNamespace, id)
}

// NetworkInternetExchangeLAN represents a network's connection to an Internet
// exchange LAN. It can be used to find common IX LANs between several
// networks.
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
	BFDSupport             bool                `json:"bfd_support"`
	Operational            bool                `json:"operational"`
	NetworkSideID          int                 `json:"net_side_id"`
	InternetExchangeSideID int                 `json:"ix_side_id"`
	Created                time.Time           `json:"created"`
	Updated                time.Time           `json:"updated"`
	Status                 string              `json:"status"`
}

// GetNetworkInternetExchangeLAN returns a slice of
// NetworkInternetExchangeLAN structures matching the given search parameters.
func (api *API) GetNetworkInternetExchangeLAN(ctx context.Context, search url.Values) ([]NetworkInternetExchangeLAN, error) {
	return fetch[NetworkInternetExchangeLAN](ctx, api, networkInternetExchangeLANNamespace, search)
}

// GetAllNetworkInternetExchangeLANs returns all NetworkInternetExchangeLAN
// structures available from the API.
func (api *API) GetAllNetworkInternetExchangeLANs(ctx context.Context) ([]NetworkInternetExchangeLAN, error) {
	return api.GetNetworkInternetExchangeLAN(ctx, nil)
}

// GetNetworkInternetExchangeLANByID returns the NetworkInternetExchangeLAN
// matching the given ID, or nil if not found.
func (api *API) GetNetworkInternetExchangeLANByID(ctx context.Context, id int) (*NetworkInternetExchangeLAN, error) {
	return fetchByID[NetworkInternetExchangeLAN](ctx, api, networkInternetExchangeLANNamespace, id)
}
