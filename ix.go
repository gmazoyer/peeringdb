package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// InternetExchange represents an Internet exchange point. It is directly
// linked to the Organization that manages the IX.
type InternetExchange struct {
	ID                     int           `json:"id"`
	OrganizationID         int           `json:"org_id"`
	Organization           Organization  `json:"org,omitempty"`
	Name                   string        `json:"name"`
	AKA                    string        `json:"aka"`
	NameLong               string        `json:"name_long"`
	City                   string        `json:"city"`
	Country                string        `json:"country"`
	RegionContinent        string        `json:"region_continent"`
	Media                  string        `json:"media"`
	Notes                  string        `json:"notes"`
	ProtoUnicast           bool          `json:"proto_unicast"`
	ProtoMulticast         bool          `json:"proto_multicast"`
	ProtoIPv6              bool          `json:"proto_ipv6"`
	Website                string        `json:"website"`
	URLStats               string        `json:"url_stats"`
	TechEmail              string        `json:"tech_email"`
	TechPhone              string        `json:"tech_phone"`
	PolicyEmail            string        `json:"policy_email"`
	PolicyPhone            string        `json:"policy_phone"`
	SalesPhone             string        `json:"sales_phone"`
	SalesEmail             string        `json:"sales_email"`
	FacilitySet            []int         `json:"fac_set"`
	InternetExchangeLANSet []int         `json:"ixlan_set"`
	NetworkCount           int           `json:"net_count"`
	FacilityCount          int           `json:"fac_count"`
	IxfNetCount            int           `json:"ixf_net_count"`
	IxfLastImport          time.Time     `json:"ixf_last_import"`
	IxfImportRequest       time.Time     `json:"ixf_import_request"`
	IxfImportRequestStatus string        `json:"ixf_import_request_status"`
	ServiceLevel           string        `json:"service_level"`
	Terms                  string        `json:"terms"`
	StatusDashboard        string        `json:"status_dashboard"`
	Created                time.Time     `json:"created"`
	Updated                time.Time     `json:"updated"`
	Status                 string        `json:"status"`
	SocialMedia            []SocialMedia `json:"social_media"`
}

// GetInternetExchange returns a slice of InternetExchange structures matching
// the given search parameters.
func (api *API) GetInternetExchange(ctx context.Context, search url.Values) ([]InternetExchange, error) {
	return fetch[InternetExchange](ctx, api, internetExchangeNamespace, search)
}

// GetAllInternetExchanges returns all InternetExchange structures available
// from the API.
func (api *API) GetAllInternetExchanges(ctx context.Context) ([]InternetExchange, error) {
	return api.GetInternetExchange(ctx, nil)
}

// GetInternetExchangeByID returns the InternetExchange matching the given ID,
// or nil if not found.
func (api *API) GetInternetExchangeByID(ctx context.Context, id int) (*InternetExchange, error) {
	return fetchByID[InternetExchange](ctx, api, internetExchangeNamespace, id)
}

// InternetExchangeLAN represents one of the networks (LANs) of an Internet
// exchange point, with details like MTU, VLAN support, etc.
type InternetExchangeLAN struct {
	ID                         int              `json:"id"`
	InternetExchangeID         int              `json:"ix_id"`
	InternetExchange           InternetExchange `json:"ix,omitempty"`
	Name                       string           `json:"name"`
	Description                string           `json:"descr"`
	MTU                        int              `json:"mtu"`
	Dot1QSupport               bool             `json:"dot1q_support"`
	RouteServerASN             int              `json:"rs_asn"`
	ARPSponge                  string           `json:"arp_sponge"`
	NetworkSet                 []int            `json:"net_set"`
	InternetExchangePrefixSet  []int            `json:"ixpfx_set"`
	IXFIXPMemberListURL        string           `json:"ixf_ixp_member_list_url"`
	IXFIXPMemberListURLVisible string           `json:"ixf_ixp_member_list_url_visible"`
	IXFIXPImportEnabled        bool             `json:"ixf_ixp_import_enabled"`
	Created                    time.Time        `json:"created"`
	Updated                    time.Time        `json:"updated"`
	Status                     string           `json:"status"`
}

// GetInternetExchangeLAN returns a slice of InternetExchangeLAN structures
// matching the given search parameters.
func (api *API) GetInternetExchangeLAN(ctx context.Context, search url.Values) ([]InternetExchangeLAN, error) {
	return fetch[InternetExchangeLAN](ctx, api, internetExchangeLANNamespace, search)
}

// GetAllInternetExchangeLANs returns all InternetExchangeLAN structures
// available from the API.
func (api *API) GetAllInternetExchangeLANs(ctx context.Context) ([]InternetExchangeLAN, error) {
	return api.GetInternetExchangeLAN(ctx, nil)
}

// GetInternetExchangeLANByID returns the InternetExchangeLAN matching the
// given ID, or nil if not found.
func (api *API) GetInternetExchangeLANByID(ctx context.Context, id int) (*InternetExchangeLAN, error) {
	return fetchByID[InternetExchangeLAN](ctx, api, internetExchangeLANNamespace, id)
}

// InternetExchangePrefix represents a prefix used by an Internet exchange
// point, directly linked to an InternetExchangeLAN.
type InternetExchangePrefix struct {
	ID                    int                 `json:"id"`
	InternetExchangeLANID int                 `json:"ixlan_id"`
	InternetExchangeLAN   InternetExchangeLAN `json:"ixlan,omitempty"`
	Protocol              string              `json:"protocol"`
	Prefix                string              `json:"prefix"`
	InDFZ                 bool                `json:"in_dfz"`
	Created               time.Time           `json:"created"`
	Updated               time.Time           `json:"updated"`
	Status                string              `json:"status"`
}

// GetInternetExchangePrefix returns a slice of InternetExchangePrefix
// structures matching the given search parameters.
func (api *API) GetInternetExchangePrefix(ctx context.Context, search url.Values) ([]InternetExchangePrefix, error) {
	return fetch[InternetExchangePrefix](ctx, api, internetExchangePrefixNamespace, search)
}

// GetAllInternetExchangePrefixes returns all InternetExchangePrefix structures
// available from the API.
func (api *API) GetAllInternetExchangePrefixes(ctx context.Context) ([]InternetExchangePrefix, error) {
	return api.GetInternetExchangePrefix(ctx, nil)
}

// GetInternetExchangePrefixByID returns the InternetExchangePrefix matching
// the given ID, or nil if not found.
func (api *API) GetInternetExchangePrefixByID(ctx context.Context, id int) (*InternetExchangePrefix, error) {
	return fetchByID[InternetExchangePrefix](ctx, api, internetExchangePrefixNamespace, id)
}

// InternetExchangeFacility links an InternetExchange with a Facility,
// indicating where an IX can be found or what IXs are in a given facility.
type InternetExchangeFacility struct {
	ID                 int              `json:"id"`
	Name               string           `json:"name"`
	City               string           `json:"city"`
	Country            string           `json:"country"`
	InternetExchangeID int              `json:"ix_id"`
	InternetExchange   InternetExchange `json:"ix,omitempty"`
	FacilityID         int              `json:"fac_id"`
	Facility           Facility         `json:"fac,omitempty"`
	Created            time.Time        `json:"created"`
	Updated            time.Time        `json:"updated"`
	Status             string           `json:"status"`
}

// GetInternetExchangeFacility returns a slice of InternetExchangeFacility
// structures matching the given search parameters.
func (api *API) GetInternetExchangeFacility(ctx context.Context, search url.Values) ([]InternetExchangeFacility, error) {
	return fetch[InternetExchangeFacility](ctx, api, internetExchangeFacilityNamespace, search)
}

// GetAllInternetExchangeFacilities returns all InternetExchangeFacility
// structures available from the API.
func (api *API) GetAllInternetExchangeFacilities(ctx context.Context) ([]InternetExchangeFacility, error) {
	return api.GetInternetExchangeFacility(ctx, nil)
}

// GetInternetExchangeFacilityByID returns the InternetExchangeFacility
// matching the given ID, or nil if not found.
func (api *API) GetInternetExchangeFacilityByID(ctx context.Context, id int) (*InternetExchangeFacility, error) {
	return fetchByID[InternetExchangeFacility](ctx, api, internetExchangeFacilityNamespace, id)
}
