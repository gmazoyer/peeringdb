package peeringdb

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

const (
	baseAPI                             = "https://www.peeringdb.com/api/"
	facilityNamespace                   = "fac"
	internetExchangeNamespace           = "ix"
	internetExchangeFacilityNamespace   = "ixfac"
	internetExchangeLANNamespace        = "ixlan"
	internetExchangePrefixNamespace     = "ixpfx"
	networkNamespace                    = "net"
	networkFacilityNamespace            = "netfac"
	networkInternetExchangeLANNamepsace = "netixlan"
	organizationNamespace               = "org"
	networkContactNamespace             = "poc"
)

var (
	// ErrBuildingURL is the error that will be returned if the URL to call the
	// API cannot be built as expected.
	ErrBuildingURL = errors.New("error while building the URL to call the peeringdb api")
	// ErrBuildingRequest is the error that will be returned if the HTTP
	// request to call the API cannot be built as expected.
	ErrBuildingRequest = errors.New("error while building the request to send to the peeringdb api")
	// ErrQueryingAPI is the error that will be returned if there is an issue
	// while making the request to the API.
	ErrQueryingAPI = errors.New("error while querying peeringdb api")
)

// API is the structure used to interact with the PeeringDB API. This is the
// main structure of this package. All functions to make API calls are
// associated to this structure.
type API struct {
	url      string
	login    string
	password string
	apiKey   string
}

// NewAPI returns a pointer to a new API structure. It uses the publicly known
// PeeringDB API endpoint.
func NewAPI() *API {
	return &API{
		url:      baseAPI,
		login:    "",
		password: "",
	}
}

// NewAPIWithAuth returns a pointer to a new API structure. The API will point
// to the publicly known PeeringDB API endpoint and will use the provided login
// and password to attempt an authentication while making API calls.
func NewAPIWithAuth(login, password string) *API {
	return &API{
		url:      baseAPI,
		login:    login,
		password: password,
	}
}

// NewAPIWithAuth returns a pointer to a new API structure. The API will point
// to the publicly known PeeringDB API endpoint and will use the provided login
// and password to attempt an authentication while making API calls.
func NewAPIWithAPIKey(apiKey string) *API {
	return &API{
		url:    baseAPI,
		apiKey: apiKey,
	}
}

// NewAPIFromURL returns a pointer to a new API structure from a given URL. If
// the given URL is empty it will use the default PeeringDB API URL.
func NewAPIFromURL(url string) *API {
	if url == "" {
		return NewAPI()
	}

	return &API{
		url:      url,
		login:    "",
		password: "",
	}
}

// NewAPIFromURLWithAuth returns a pointer to a new API structure from a given
// URL. If the given URL is empty it will use the default PeeringDB API URL. It
// will use the provided login and password to attempt an authentication while
// making API calls.
func NewAPIFromURLWithAuth(url, login, password string) *API {
	if url == "" {
		return NewAPIWithAuth(login, password)
	}

	return &API{
		url:      url,
		login:    login,
		password: password,
	}
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

// formatURL is used to format The URL to call the PeeringDB API.
func formatURL(base, namespace string, search map[string]interface{}) string {
	return fmt.Sprintf("%s%s?depth=1%s", base, namespace,
		formatSearchParameters(search))
}

// lookup is used to query the PeeringDB API given a namespace to use and data
// to format the request. It returns an HTTP response that the caller must
// decode with a JSON decoder.
func (api *API) lookup(namespace string, search map[string]interface{}) (*http.Response, error) {
	url := formatURL(api.url, namespace, search)
	if url == "" {
		return nil, ErrBuildingURL
	}

	// Prepare the GET request to the API, no need to set a body since
	// everything is in the URL
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, ErrBuildingRequest
	}

	// If auth credentials are provided, use them
	if (api.login != "") && (api.password != "") {
		request.SetBasicAuth(api.login, api.password)
	}

	if api.apiKey != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Api-Key %s", api.apiKey))
	}

	// Send the request to the API using a simple HTTP client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, ErrQueryingAPI
	}

	return response, nil
}

// GetASN is a simplified function to get PeeringDB details about a given AS
// number. It basically gets the Net object matching the AS number. If the AS
// number cannot be found, nil is returned.
func (api *API) GetASN(asn int) *Network {
	search := make(map[string]interface{})
	search["asn"] = asn

	// Actually fetch the Network from PeeringDB
	network, err := api.GetNetwork(search)

	// Error, so nil pointer returned
	if err != nil {
		return nil
	}

	if len(*network) == 0 {
		return nil
	}
	return &(*network)[0]
}
