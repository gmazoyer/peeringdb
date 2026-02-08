package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// Organization represents an enterprise linked to networks, facilities,
// and internet exchange points.
type Organization struct {
	ID                  int           `json:"id"`
	Name                string        `json:"name"`
	AKA                 string        `json:"aka"`
	NameLong            string        `json:"name_long"`
	Website             string        `json:"website"`
	Notes               string        `json:"notes"`
	Require2FA          bool          `json:"require_2fa"`
	NetworkSet          []int         `json:"net_set"`
	FacilitySet         []int         `json:"fac_set"`
	InternetExchangeSet []int         `json:"ix_set"`
	CarrierSet          []int         `json:"carrier_set"`
	CampusSet           []int         `json:"campus_set"`
	Address1            string        `json:"address1"`
	Address2            string        `json:"address2"`
	City                string        `json:"city"`
	Country             string        `json:"country"`
	State               string        `json:"state"`
	Zipcode             string        `json:"zipcode"`
	Floor               string        `json:"floor"`
	Suite               string        `json:"suite"`
	Latitude            float64       `json:"latitude"`
	Longitude           float64       `json:"longitude"`
	Created             time.Time     `json:"created"`
	Updated             time.Time     `json:"updated"`
	Status              string        `json:"status"`
	SocialMedia         []SocialMedia `json:"social_media"`
}

// GetOrganization returns a slice of Organization structures matching the
// given search parameters.
func (api *API) GetOrganization(ctx context.Context, search url.Values) ([]Organization, error) {
	return fetch[Organization](ctx, api, organizationNamespace, search)
}

// GetAllOrganizations returns all Organization structures available from the
// API.
func (api *API) GetAllOrganizations(ctx context.Context) ([]Organization, error) {
	return api.GetOrganization(ctx, nil)
}

// GetOrganizationByID returns the Organization matching the given ID, or nil
// if not found.
func (api *API) GetOrganizationByID(ctx context.Context, id int) (*Organization, error) {
	return fetchByID[Organization](ctx, api, organizationNamespace, id)
}
