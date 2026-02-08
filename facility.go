package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// Facility is the representation of a location where network operators and
// Internet exchange points are located. Most of the time you know a facility
// as a datacenter.
type Facility struct {
	ID                        int           `json:"id"`
	OrganizationID            int           `json:"org_id"`
	OrganizationName          string        `json:"org_name"`
	Organization              Organization  `json:"organization,omitempty"`
	CampusID                  int           `json:"campus_id"`
	Campus                    Campus        `json:"campus,omitempty"`
	Name                      string        `json:"name"`
	AKA                       string        `json:"aka"`
	NameLong                  string        `json:"name_long"`
	Website                   string        `json:"website"`
	CLLI                      string        `json:"clli"`
	Rencode                   string        `json:"rencode"`
	Npanxx                    string        `json:"npanxx"`
	Notes                     string        `json:"notes"`
	NetCount                  int           `json:"net_count"`
	IXCount                   int           `json:"ix_count"`
	CarrierCount              int           `json:"carrier_count"`
	SalesEmail                string        `json:"sales_email"`
	SalesPhone                string        `json:"sales_phone"`
	TechEmail                 string        `json:"tech_email"`
	TechPhone                 string        `json:"tech_phone"`
	AvailableVoltageServices  []string      `json:"available_voltage_services"`
	DiverseServingSubstations bool          `json:"diverse_serving_substations"`
	Property                  string        `json:"property"`
	RegionContinent           string        `json:"region_continent"`
	StatusDashboard           string        `json:"status_dashboard"`
	Created                   time.Time     `json:"created"`
	Updated                   time.Time     `json:"updated"`
	Status                    string        `json:"status"`
	Address1                  string        `json:"address1"`
	Address2                  string        `json:"address2"`
	City                      string        `json:"city"`
	Country                   string        `json:"country"`
	State                     string        `json:"state"`
	Zipcode                   string        `json:"zipcode"`
	Floor                     string        `json:"floor"`
	Suite                     string        `json:"suite"`
	Latitude                  float64       `json:"latitude"`
	Longitude                 float64       `json:"longitude"`
	SocialMedia               []SocialMedia `json:"social_media"`
}

// GetFacility returns a slice of Facility structures matching the given
// search parameters.
func (api *API) GetFacility(ctx context.Context, search url.Values) ([]Facility, error) {
	return fetch[Facility](ctx, api, facilityNamespace, search)
}

// GetAllFacilities returns all Facility structures available from the API.
func (api *API) GetAllFacilities(ctx context.Context) ([]Facility, error) {
	return api.GetFacility(ctx, nil)
}

// GetFacilityByID returns the Facility matching the given ID, or nil if not
// found.
func (api *API) GetFacilityByID(ctx context.Context, id int) (*Facility, error) {
	return fetchByID[Facility](ctx, api, facilityNamespace, id)
}
