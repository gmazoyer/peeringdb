package peeringdb

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// newTestServer creates a test HTTP server that responds with the given
// status code and body for any request.
func newTestServer(t *testing.T, statusCode int, body interface{}) (*httptest.Server, *API) {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(body)
	}))

	api := NewAPI(WithURL(ts.URL + "/"))
	t.Cleanup(ts.Close)

	return ts, api
}

func TestNewAPI(t *testing.T) {
	tests := []struct {
		name       string
		opts       []Option
		wantURL    string
		wantAPIKey string
	}{
		{
			name:    "default",
			wantURL: baseAPI,
		},
		{
			name:       "with API key",
			opts:       []Option{WithAPIKey("test123")},
			wantURL:    baseAPI,
			wantAPIKey: "test123",
		},
		{
			name:    "with custom URL",
			opts:    []Option{WithURL("http://localhost/api/")},
			wantURL: "http://localhost/api/",
		},
		{
			name:       "with URL and API key",
			opts:       []Option{WithURL("http://localhost/api/"), WithAPIKey("test123")},
			wantURL:    "http://localhost/api/",
			wantAPIKey: "test123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(tt.opts...)
			if api.url != tt.wantURL {
				t.Errorf("url: want %q, got %q", tt.wantURL, api.url)
			}
			if api.apiKey != tt.wantAPIKey {
				t.Errorf("apiKey: want %q, got %q", tt.wantAPIKey, api.apiKey)
			}
			if api.client == nil {
				t.Error("client should not be nil")
			}
		})
	}
}

func TestFetchSuccess(t *testing.T) {
	resp := resource[Organization]{
		Data: []Organization{{ID: 1, Name: "Test Org"}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	orgs, err := api.GetOrganization(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0].Name != "Test Org" {
		t.Errorf("expected name %q, got %q", "Test Org", orgs[0].Name)
	}
}

func TestFetchWithSearchParams(t *testing.T) {
	var receivedURL string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resource[Network]{})
	}))
	defer ts.Close()

	api := NewAPI(WithURL(ts.URL + "/"))
	search := url.Values{}
	search.Set("asn", "65536")

	_, err := api.GetNetwork(context.Background(), search)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedURL != "/net?depth=1&asn=65536" {
		t.Errorf("unexpected URL: %s", receivedURL)
	}
}

func TestFetchWithAPIKey(t *testing.T) {
	var receivedAuth string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resource[Network]{})
	}))
	defer ts.Close()

	api := NewAPI(WithURL(ts.URL+"/"), WithAPIKey("secret-key"))
	_, err := api.GetNetwork(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuth != "Api-Key secret-key" {
		t.Errorf("expected auth header %q, got %q", "Api-Key secret-key", receivedAuth)
	}
}

func TestFetchRateLimit(t *testing.T) {
	_, api := newTestServer(t, http.StatusTooManyRequests, nil)

	_, err := api.GetNetwork(context.Background(), nil)
	if !errors.Is(err, ErrRateLimitExceeded) {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}
}

func TestFetchServerError(t *testing.T) {
	_, api := newTestServer(t, http.StatusInternalServerError, "something went wrong")

	_, err := api.GetNetwork(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestFetchInvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not json"))
	}))
	defer ts.Close()

	api := NewAPI(WithURL(ts.URL + "/"))
	_, err := api.GetNetwork(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFetchConnectionError(t *testing.T) {
	api := NewAPI(WithURL("http://127.0.0.1:1/"))
	_, err := api.GetNetwork(context.Background(), nil)
	if !errors.Is(err, ErrQueryingAPI) {
		t.Errorf("expected ErrQueryingAPI, got %v", err)
	}
}

func TestFetchContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block forever; the cancelled context should abort the request
		select {}
	}))
	defer ts.Close()

	api := NewAPI(WithURL(ts.URL + "/"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := api.GetNetwork(ctx, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestFetchByIDSuccess(t *testing.T) {
	resp := resource[Network]{
		Data: []Network{{ID: 42, Name: "Test Net", ASN: 65536}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	net, err := api.GetNetworkByID(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if net == nil {
		t.Fatal("expected non-nil network")
	}
	if net.ID != 42 {
		t.Errorf("expected ID 42, got %d", net.ID)
	}
}

func TestFetchByIDNotFound(t *testing.T) {
	resp := resource[Network]{Data: []Network{}}
	_, api := newTestServer(t, http.StatusOK, resp)

	net, err := api.GetNetworkByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if net != nil {
		t.Errorf("expected nil for not found, got %+v", net)
	}
}

func TestFetchByIDInvalid(t *testing.T) {
	api := NewAPI()

	_, err := api.GetNetworkByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("expected ErrInvalidID for id=0, got %v", err)
	}

	_, err = api.GetNetworkByID(context.Background(), -1)
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("expected ErrInvalidID for id=-1, got %v", err)
	}
}

func TestGetASNSuccess(t *testing.T) {
	resp := resource[Network]{
		Data: []Network{{ID: 1, Name: "Test AS", ASN: 201281}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	net, err := api.GetASN(context.Background(), 201281)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if net.ASN != 201281 {
		t.Errorf("expected ASN 201281, got %d", net.ASN)
	}
}

func TestGetASNNotFound(t *testing.T) {
	resp := resource[Network]{Data: []Network{}}
	_, api := newTestServer(t, http.StatusOK, resp)

	_, err := api.GetASN(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error for ASN not found")
	}
}

func TestGetOrganization(t *testing.T) {
	resp := resource[Organization]{
		Data: []Organization{{ID: 1, Name: "Test Org"}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	orgs, err := api.GetOrganization(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 || orgs[0].Name != "Test Org" {
		t.Errorf("unexpected result: %+v", orgs)
	}
}

func TestGetInternetExchange(t *testing.T) {
	resp := resource[InternetExchange]{
		Data: []InternetExchange{{ID: 1, Name: "Test IX"}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	ixs, err := api.GetInternetExchange(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ixs) != 1 || ixs[0].Name != "Test IX" {
		t.Errorf("unexpected result: %+v", ixs)
	}
}

func TestGetCarrier(t *testing.T) {
	resp := resource[Carrier]{
		Data: []Carrier{{ID: 1, Name: "Test Carrier"}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	carriers, err := api.GetCarrier(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(carriers) != 1 || carriers[0].Name != "Test Carrier" {
		t.Errorf("unexpected result: %+v", carriers)
	}
}

func TestGetCampus(t *testing.T) {
	resp := resource[Campus]{
		Data: []Campus{{ID: 1, Name: "Test Campus"}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	campuses, err := api.GetCampus(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(campuses) != 1 || campuses[0].Name != "Test Campus" {
		t.Errorf("unexpected result: %+v", campuses)
	}
}

func TestGetNetworkContact(t *testing.T) {
	resp := resource[NetworkContact]{
		Data: []NetworkContact{{ID: 1, Name: "Admin Contact"}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	contacts, err := api.GetNetworkContact(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 1 || contacts[0].Name != "Admin Contact" {
		t.Errorf("unexpected result: %+v", contacts)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	api := NewAPI(WithHTTPClient(custom))
	if api.client != custom {
		t.Error("expected custom HTTP client to be set")
	}
}

func TestGetAllOrganizations(t *testing.T) {
	resp := resource[Organization]{Data: []Organization{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	orgs, err := api.GetAllOrganizations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 {
		t.Errorf("expected 1 org, got %d", len(orgs))
	}
}

func TestGetOrganizationByID(t *testing.T) {
	resp := resource[Organization]{Data: []Organization{{ID: 5, Name: "Org 5"}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	org, err := api.GetOrganizationByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org == nil || org.ID != 5 {
		t.Errorf("unexpected result: %+v", org)
	}
}

func TestGetAllNetworks(t *testing.T) {
	resp := resource[Network]{Data: []Network{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nets, err := api.GetAllNetworks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nets) != 1 {
		t.Errorf("expected 1 network, got %d", len(nets))
	}
}

func TestGetAllNetworkFacilities(t *testing.T) {
	resp := resource[NetworkFacility]{Data: []NetworkFacility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nfs, err := api.GetAllNetworkFacilities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nfs) != 1 {
		t.Errorf("expected 1 netfac, got %d", len(nfs))
	}
}

func TestGetNetworkFacility(t *testing.T) {
	resp := resource[NetworkFacility]{Data: []NetworkFacility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nfs, err := api.GetNetworkFacility(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nfs) != 1 {
		t.Errorf("expected 1 netfac, got %d", len(nfs))
	}
}

func TestGetNetworkFacilityByID(t *testing.T) {
	resp := resource[NetworkFacility]{Data: []NetworkFacility{{ID: 3}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nf, err := api.GetNetworkFacilityByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nf == nil || nf.ID != 3 {
		t.Errorf("unexpected result: %+v", nf)
	}
}

func TestGetNetworkInternetExchangeLAN(t *testing.T) {
	resp := resource[NetworkInternetExchangeLAN]{Data: []NetworkInternetExchangeLAN{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nixlans, err := api.GetNetworkInternetExchangeLAN(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nixlans) != 1 {
		t.Errorf("expected 1, got %d", len(nixlans))
	}
}

func TestGetAllNetworkInternetExchangeLANs(t *testing.T) {
	resp := resource[NetworkInternetExchangeLAN]{Data: []NetworkInternetExchangeLAN{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nixlans, err := api.GetAllNetworkInternetExchangeLANs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nixlans) != 1 {
		t.Errorf("expected 1, got %d", len(nixlans))
	}
}

func TestGetNetworkInternetExchangeLANByID(t *testing.T) {
	resp := resource[NetworkInternetExchangeLAN]{Data: []NetworkInternetExchangeLAN{{ID: 7}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	nixlan, err := api.GetNetworkInternetExchangeLANByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nixlan == nil || nixlan.ID != 7 {
		t.Errorf("unexpected result: %+v", nixlan)
	}
}

func TestGetAllInternetExchanges(t *testing.T) {
	resp := resource[InternetExchange]{Data: []InternetExchange{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	ixs, err := api.GetAllInternetExchanges(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ixs) != 1 {
		t.Errorf("expected 1, got %d", len(ixs))
	}
}

func TestGetInternetExchangeByID(t *testing.T) {
	resp := resource[InternetExchange]{Data: []InternetExchange{{ID: 2}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	ix, err := api.GetInternetExchangeByID(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ix == nil || ix.ID != 2 {
		t.Errorf("unexpected result: %+v", ix)
	}
}

func TestGetInternetExchangeLAN(t *testing.T) {
	resp := resource[InternetExchangeLAN]{Data: []InternetExchangeLAN{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	lans, err := api.GetInternetExchangeLAN(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lans) != 1 {
		t.Errorf("expected 1, got %d", len(lans))
	}
}

func TestGetAllInternetExchangeLANs(t *testing.T) {
	resp := resource[InternetExchangeLAN]{Data: []InternetExchangeLAN{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	lans, err := api.GetAllInternetExchangeLANs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lans) != 1 {
		t.Errorf("expected 1, got %d", len(lans))
	}
}

func TestGetInternetExchangeLANByID(t *testing.T) {
	resp := resource[InternetExchangeLAN]{Data: []InternetExchangeLAN{{ID: 4}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	lan, err := api.GetInternetExchangeLANByID(context.Background(), 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lan == nil || lan.ID != 4 {
		t.Errorf("unexpected result: %+v", lan)
	}
}

func TestGetInternetExchangePrefix(t *testing.T) {
	resp := resource[InternetExchangePrefix]{Data: []InternetExchangePrefix{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	pfxs, err := api.GetInternetExchangePrefix(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pfxs) != 1 {
		t.Errorf("expected 1, got %d", len(pfxs))
	}
}

func TestGetAllInternetExchangePrefixes(t *testing.T) {
	resp := resource[InternetExchangePrefix]{Data: []InternetExchangePrefix{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	pfxs, err := api.GetAllInternetExchangePrefixes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pfxs) != 1 {
		t.Errorf("expected 1, got %d", len(pfxs))
	}
}

func TestGetInternetExchangePrefixByID(t *testing.T) {
	resp := resource[InternetExchangePrefix]{Data: []InternetExchangePrefix{{ID: 6}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	pfx, err := api.GetInternetExchangePrefixByID(context.Background(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pfx == nil || pfx.ID != 6 {
		t.Errorf("unexpected result: %+v", pfx)
	}
}

func TestGetInternetExchangeFacility(t *testing.T) {
	resp := resource[InternetExchangeFacility]{Data: []InternetExchangeFacility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	ixfacs, err := api.GetInternetExchangeFacility(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ixfacs) != 1 {
		t.Errorf("expected 1, got %d", len(ixfacs))
	}
}

func TestGetAllInternetExchangeFacilities(t *testing.T) {
	resp := resource[InternetExchangeFacility]{Data: []InternetExchangeFacility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	ixfacs, err := api.GetAllInternetExchangeFacilities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ixfacs) != 1 {
		t.Errorf("expected 1, got %d", len(ixfacs))
	}
}

func TestGetInternetExchangeFacilityByID(t *testing.T) {
	resp := resource[InternetExchangeFacility]{Data: []InternetExchangeFacility{{ID: 8}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	ixfac, err := api.GetInternetExchangeFacilityByID(context.Background(), 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ixfac == nil || ixfac.ID != 8 {
		t.Errorf("unexpected result: %+v", ixfac)
	}
}

func TestGetFacility(t *testing.T) {
	resp := resource[Facility]{Data: []Facility{{ID: 1, Name: "Test DC"}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	facs, err := api.GetFacility(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(facs) != 1 || facs[0].Name != "Test DC" {
		t.Errorf("unexpected result: %+v", facs)
	}
}

func TestGetAllFacilities(t *testing.T) {
	resp := resource[Facility]{Data: []Facility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	facs, err := api.GetAllFacilities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(facs) != 1 {
		t.Errorf("expected 1, got %d", len(facs))
	}
}

func TestGetFacilityByID(t *testing.T) {
	resp := resource[Facility]{Data: []Facility{{ID: 9}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	fac, err := api.GetFacilityByID(context.Background(), 9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fac == nil || fac.ID != 9 {
		t.Errorf("unexpected result: %+v", fac)
	}
}

func TestGetAllCarriers(t *testing.T) {
	resp := resource[Carrier]{Data: []Carrier{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	carriers, err := api.GetAllCarriers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(carriers) != 1 {
		t.Errorf("expected 1, got %d", len(carriers))
	}
}

func TestGetCarrierByID(t *testing.T) {
	resp := resource[Carrier]{Data: []Carrier{{ID: 3}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	c, err := api.GetCarrierByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil || c.ID != 3 {
		t.Errorf("unexpected result: %+v", c)
	}
}

func TestGetCarrierFacility(t *testing.T) {
	resp := resource[CarrierFacility]{Data: []CarrierFacility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	cfs, err := api.GetCarrierFacility(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfs) != 1 {
		t.Errorf("expected 1, got %d", len(cfs))
	}
}

func TestGetAllCarrierFacilities(t *testing.T) {
	resp := resource[CarrierFacility]{Data: []CarrierFacility{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	cfs, err := api.GetAllCarrierFacilities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfs) != 1 {
		t.Errorf("expected 1, got %d", len(cfs))
	}
}

func TestGetCarrierFacilityByID(t *testing.T) {
	resp := resource[CarrierFacility]{Data: []CarrierFacility{{ID: 5}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	cf, err := api.GetCarrierFacilityByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cf == nil || cf.ID != 5 {
		t.Errorf("unexpected result: %+v", cf)
	}
}

func TestGetAllCampuses(t *testing.T) {
	resp := resource[Campus]{Data: []Campus{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	campuses, err := api.GetAllCampuses(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(campuses) != 1 {
		t.Errorf("expected 1, got %d", len(campuses))
	}
}

func TestGetCampusByID(t *testing.T) {
	resp := resource[Campus]{Data: []Campus{{ID: 2}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	c, err := api.GetCampusByID(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil || c.ID != 2 {
		t.Errorf("unexpected result: %+v", c)
	}
}

func TestGetAllNetworkContacts(t *testing.T) {
	resp := resource[NetworkContact]{Data: []NetworkContact{{ID: 1}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	contacts, err := api.GetAllNetworkContacts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 1 {
		t.Errorf("expected 1, got %d", len(contacts))
	}
}

func TestGetNetworkContactByID(t *testing.T) {
	resp := resource[NetworkContact]{Data: []NetworkContact{{ID: 4}}}
	_, api := newTestServer(t, http.StatusOK, resp)
	c, err := api.GetNetworkContactByID(context.Background(), 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil || c.ID != 4 {
		t.Errorf("unexpected result: %+v", c)
	}
}

func TestSocialMediaDeserialization(t *testing.T) {
	resp := resource[Organization]{
		Data: []Organization{{
			ID:   1,
			Name: "Social Org",
			SocialMedia: []SocialMedia{
				{Service: "website", Identifier: "https://example.com"},
				{Service: "twitter", Identifier: "@example"},
			},
		}},
	}
	_, api := newTestServer(t, http.StatusOK, resp)

	orgs, err := api.GetOrganization(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs[0].SocialMedia) != 2 {
		t.Fatalf("expected 2 social media entries, got %d", len(orgs[0].SocialMedia))
	}
	if orgs[0].SocialMedia[0].Service != "website" {
		t.Errorf("expected service %q, got %q", "website", orgs[0].SocialMedia[0].Service)
	}
}
