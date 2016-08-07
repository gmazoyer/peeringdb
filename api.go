package peeringdb

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
)

const (
	baseAPI                             = "https://peeringdb.com/api/"
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
	ErrBuildingURL = errors.New("Error while building the URL to call the PeeringDB API.")
	// ErrBuildingRequest is the error that will be returned if the HTTP
	// request to call the API cannot be built as expected.
	ErrBuildingRequest = errors.New("Error while building the request to send to the PeeringDB API.")
	// ErrQueryingAPI is the error that will be returned if there is an issue
	// while making the request to the API.
	ErrQueryingAPI = errors.New("Error while querying PeeringDB API.")
)

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
		search = fmt.Sprintf("%s&%s=%v", search, key, parameters[key])
	}

	return search
}

// formatURL is used to format The URL to call the PeeringDB API.
func formatURL(namespace string, search map[string]interface{}) string {
	return fmt.Sprintf("%s%s?depth=1%s", baseAPI, namespace,
		formatSearchParameters(search))
}

// lookup is used to query the PeeringDB API given a namespace to use and data
// to format the request. It returns an HTTP response that the caller must
// decode with a JSON decoder. If auth is provided, non-nil slice with two
// values, the first one being the username and the second one being the
// password, an authentication is made.
func lookup(namespace string, auth []string, search map[string]interface{}) (*http.Response, error) {
	url := formatURL(namespace, search)
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
	if (auth != nil) && (len(auth) == 2) {
		request.SetBasicAuth(auth[0], auth[1])
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
func GetASN(asn int) *Network {
	search := make(map[string]interface{})
	search["asn"] = asn

	// Actually fetch the Network from PeeringDB
	network, err := GetNetwork(search)

	// Error, so nil pointer returned
	if err != nil {
		return nil
	}

	return &(*network)[0]
}
