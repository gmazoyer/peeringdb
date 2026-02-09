package peeringdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const Version = "0.1.0"

const (
	baseAPI                             = "https://www.peeringdb.com/api/"
	facilityNamespace                   = "fac"
	carrierNamespace                    = "carrier"
	carrierFacilityNamespace            = "carrierfac"
	campusNamespace                     = "campus"
	internetExchangeNamespace           = "ix"
	internetExchangeFacilityNamespace   = "ixfac"
	internetExchangeLANNamespace        = "ixlan"
	internetExchangePrefixNamespace     = "ixpfx"
	networkNamespace                    = "net"
	networkFacilityNamespace            = "netfac"
	networkInternetExchangeLANNamespace = "netixlan"
	organizationNamespace               = "org"
	networkContactNamespace             = "poc"
)

var (
	// ErrBuildingRequest is returned if the HTTP request to call the API
	// cannot be built as expected.
	ErrBuildingRequest = errors.New("error building request for peeringdb api")
	// ErrQueryingAPI is returned if there is an issue while making the
	// request to the API.
	ErrQueryingAPI = errors.New("error querying peeringdb api")
	// ErrRateLimitExceeded is returned if the API rate limit is exceeded.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	// ErrInvalidID is returned when a resource ID is invalid (e.g. <= 0).
	ErrInvalidID = errors.New("invalid resource ID")
)

// SocialMedia represents a social media link for a PeeringDB entity.
type SocialMedia struct {
	Service    string `json:"service"`
	Identifier string `json:"identifier"`
}

// resource is a generic top-level structure for parsing JSON responses from
// the PeeringDB API.
type resource[T any] struct {
	Meta struct {
		Generated float64 `json:"generated,omitempty"`
	} `json:"meta"`
	Data []T `json:"data"`
}

// API is the structure used to interact with the PeeringDB API.
type API struct {
	url    string
	apiKey string
	client *http.Client
}

// Option configures an API instance.
type Option func(*API)

// WithURL sets a custom API endpoint URL.
func WithURL(u string) Option {
	return func(api *API) {
		api.url = u
	}
}

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) Option {
	return func(api *API) {
		api.apiKey = apiKey
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(api *API) {
		api.client = client
	}
}

// NewAPI returns a new API instance configured with the given options.
func NewAPI(opts ...Option) *API {
	api := &API{
		url:    baseAPI,
		client: &http.Client{},
	}
	for _, opt := range opts {
		opt(api)
	}
	return api
}

// fetch queries the PeeringDB API for a given namespace and returns a slice
// of results. It uses generics to avoid per-entity boilerplate.
func fetch[T any](ctx context.Context, api *API, namespace string, search url.Values) ([]T, error) {
	u := api.url + namespace + "?depth=1"
	if len(search) > 0 {
		u += "&" + search.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuildingRequest, err)
	}

	request.Header.Set("User-Agent", "go-peeringdb/"+Version)
	if api.apiKey != "" {
		request.Header.Set("Authorization", "Api-Key "+api.apiKey)
	}

	response, err := api.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryingAPI, err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimitExceeded
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("%s: %s", response.Status, body)
	}

	var r resource[T]
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r.Data, nil
}

// fetchByID queries the PeeringDB API for a single resource by ID.
func fetchByID[T any](ctx context.Context, api *API, namespace string, id int) (*T, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	search := url.Values{}
	search.Set("id", strconv.Itoa(id))

	results, err := fetch[T](ctx, api, namespace, search)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

// GetASN returns the Network matching the given AS number. If the AS number
// cannot be found, an error is returned.
func (api *API) GetASN(ctx context.Context, asn int) (*Network, error) {
	search := url.Values{}
	search.Set("asn", strconv.Itoa(asn))

	networks, err := api.GetNetwork(ctx, search)
	if err != nil {
		return nil, err
	}

	if len(networks) == 0 {
		return nil, fmt.Errorf("no network found for ASN %d", asn)
	}

	return &networks[0], nil
}
