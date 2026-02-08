package peeringdb

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
)

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
	// ErrBuildingRequest is the error that will be returned if the HTTP
	// request to call the API cannot be built as expected.
	ErrBuildingRequest = errors.New("error while building the request to send to the peeringdb api")
	// ErrQueryingAPI is the error that will be returned if there is an issue
	// while making the request to the API.
	ErrQueryingAPI = errors.New("error while querying peeringdb api")
	// ErrRateLimitExceeded is the error that will be returned if the API rate
	// limit is exceeded.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// SocialMedia represents a social media link for a PeeringDB entity.
type SocialMedia struct {
	Service    string `json:"service"`
	Identifier string `json:"identifier"`
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

// formatSearchParameters is used to format parameters for a request. When
// building the search string the keys will be used in the alphabetic order.
func formatSearchParameters(parameters map[string]interface{}) string {
	// Nothing in slice, just return empty string
	if parameters == nil {
		return ""
	}

	var search string
	var keys []string

	// Get all map keys
	for i := range parameters {
		keys = append(keys, i)
	}

	// Sort the keys slice
	sort.Strings(keys)

	// For each element, append it to the request separated by a & symbol.
	for _, key := range keys {
		search = search + "&" + key + "=" + url.QueryEscape(fmt.Sprintf("%v", parameters[key]))
	}

	return search
}

// formatURL is used to format a URL to make a request on PeeringDB API.
func formatURL(base, namespace string, search map[string]interface{}) string {
	return fmt.Sprintf("%s%s?depth=1%s", base, namespace,
		formatSearchParameters(search))
}

// lookup is used to query the PeeringDB API given a namespace to use and data
// to format the request. It returns an HTTP response that the caller must
// decode with a JSON decoder.
func (api *API) lookup(namespace string, search map[string]interface{}) (*http.Response, error) {
	url := formatURL(api.url, namespace, search)

	// Prepare the GET request to the API, no need to set a body since
	// everything is in the URL
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuildingRequest, err)
	}

	if api.apiKey != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Api-Key %s", api.apiKey))
	}

	response, err := api.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryingAPI, err)
	}

	// Special handling for PeeringDB rate limit
	if response.StatusCode == http.StatusTooManyRequests {
		response.Body.Close()
		return nil, ErrRateLimitExceeded
	}
	// Generic handling for non-OK responses
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		return nil, fmt.Errorf("%s: %s", response.Status, body)
	}

	return response, nil
}

// GetASN is a simplified function to get PeeringDB details about a given AS
// number. It basically gets the Net object matching the AS number. If the AS
// number cannot be found, nil is returned.
func (api *API) GetASN(asn int) (*Network, error) {
	search := make(map[string]interface{})
	search["asn"] = asn

	// Actually fetch the Network from PeeringDB
	network, err := api.GetNetwork(search)

	// Error, so nil pointer returned
	if err != nil {
		return nil, err
	}

	if len(*network) == 0 {
		return nil, fmt.Errorf("no network found for ASN %d", asn)
	}
	return &(*network)[0], nil
}
