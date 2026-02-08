package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// Carrier is the representation of a network able to provide transport from
// one facility to another.
type Carrier struct {
	ID               int           `json:"id"`
	OrganizationID   int           `json:"org_id"`
	OrganizationName string        `json:"org_name"`
	Organization     Organization  `json:"organization,omitempty"`
	Name             string        `json:"name"`
	AKA              string        `json:"aka"`
	NameLong         string        `json:"name_long"`
	Website          string        `json:"website"`
	Notes            string        `json:"notes"`
	Created          time.Time     `json:"created"`
	Updated          time.Time     `json:"updated"`
	Status           string        `json:"status"`
	SocialMedia      []SocialMedia `json:"social_media"`
}

// GetCarrier returns a slice of Carrier structures matching the given search
// parameters.
func (api *API) GetCarrier(ctx context.Context, search url.Values) ([]Carrier, error) {
	return fetch[Carrier](ctx, api, carrierNamespace, search)
}

// GetAllCarriers returns all Carrier structures available from the API.
func (api *API) GetAllCarriers(ctx context.Context) ([]Carrier, error) {
	return api.GetCarrier(ctx, nil)
}

// GetCarrierByID returns the Carrier matching the given ID, or nil if not
// found.
func (api *API) GetCarrierByID(ctx context.Context, id int) (*Carrier, error) {
	return fetchByID[Carrier](ctx, api, carrierNamespace, id)
}

// CarrierFacility links a Carrier with a Facility, indicating in which
// facilities a carrier operates.
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

// GetCarrierFacility returns a slice of CarrierFacility structures matching
// the given search parameters.
func (api *API) GetCarrierFacility(ctx context.Context, search url.Values) ([]CarrierFacility, error) {
	return fetch[CarrierFacility](ctx, api, carrierFacilityNamespace, search)
}

// GetAllCarrierFacilities returns all CarrierFacility structures available
// from the API.
func (api *API) GetAllCarrierFacilities(ctx context.Context) ([]CarrierFacility, error) {
	return api.GetCarrierFacility(ctx, nil)
}

// GetCarrierFacilityByID returns the CarrierFacility matching the given ID,
// or nil if not found.
func (api *API) GetCarrierFacilityByID(ctx context.Context, id int) (*CarrierFacility, error) {
	return fetchByID[CarrierFacility](ctx, api, carrierFacilityNamespace, id)
}
