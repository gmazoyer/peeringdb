package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// Campus is the representation of a site where facilities are.
type Campus struct {
	ID               int           `json:"id"`
	OrganizationID   int           `json:"org_id"`
	OrganizationName string        `json:"org_name"`
	Organization     Organization  `json:"organization,omitempty"`
	Name             string        `json:"name"`
	NameLong         string        `json:"name_long"`
	AKA              string        `json:"aka"`
	Website          string        `json:"website"`
	Notes            string        `json:"notes"`
	Created          time.Time     `json:"created"`
	Updated          time.Time     `json:"updated"`
	Status           string        `json:"status"`
	City             string        `json:"city"`
	Country          string        `json:"country"`
	State            string        `json:"state"`
	Zipcode          string        `json:"zipcode"`
	FacilitySet      []int         `json:"fac_set"`
	SocialMedia      []SocialMedia `json:"social_media"`
}

// GetCampus returns a slice of Campus structures matching the given search
// parameters.
func (api *API) GetCampus(ctx context.Context, search url.Values) ([]Campus, error) {
	return fetch[Campus](ctx, api, campusNamespace, search)
}

// GetAllCampuses returns all Campus structures available from the API.
func (api *API) GetAllCampuses(ctx context.Context) ([]Campus, error) {
	return api.GetCampus(ctx, nil)
}

// GetCampusByID returns the Campus matching the given ID, or nil if not
// found.
func (api *API) GetCampusByID(ctx context.Context, id int) (*Campus, error) {
	return fetchByID[Campus](ctx, api, campusNamespace, id)
}
