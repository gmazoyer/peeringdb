package peeringdb

import (
	"context"
	"net/url"
	"time"
)

// NetworkContact represents a contact for a network.
type NetworkContact struct {
	ID        int       `json:"id"`
	NetworkID int       `json:"net_id"`
	Network   Network   `json:"net"`
	Role      string    `json:"role"`
	Visible   string    `json:"visible"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	URL       string    `json:"url"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Status    string    `json:"status"`
}

// GetNetworkContact returns a slice of NetworkContact structures matching the
// given search parameters.
func (api *API) GetNetworkContact(ctx context.Context, search url.Values) ([]NetworkContact, error) {
	return fetch[NetworkContact](ctx, api, networkContactNamespace, search)
}

// GetAllNetworkContacts returns all NetworkContact structures available from
// the API.
func (api *API) GetAllNetworkContacts(ctx context.Context) ([]NetworkContact, error) {
	return api.GetNetworkContact(ctx, nil)
}

// GetNetworkContactByID returns the NetworkContact matching the given ID, or
// nil if not found.
func (api *API) GetNetworkContactByID(ctx context.Context, id int) (*NetworkContact, error) {
	return fetchByID[NetworkContact](ctx, api, networkContactNamespace, id)
}
