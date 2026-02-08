package peeringdb

import "testing"

func TestFormatSearchParameters(t *testing.T) {
	var searchMap map[string]interface{}
	var expected string
	var searchParameters string

	// Test for nil map
	expected = ""
	searchParameters = formatSearchParameters(nil)
	if searchParameters != expected {
		t.Errorf("formatSearchParameters, want '%s' got '%s'", expected,
			searchParameters)
	}

	// Test for empty map
	searchMap = make(map[string]interface{})
	expected = ""
	searchParameters = formatSearchParameters(searchMap)
	if searchParameters != expected {
		t.Errorf("formatSearchParameters, want '%s' got '%s'", expected,
			searchParameters)
	}

	// Test one value
	searchMap = make(map[string]interface{})
	searchMap["id"] = 10
	expected = "&id=10"
	searchParameters = formatSearchParameters(searchMap)
	if searchParameters != expected {
		t.Errorf("formatSearchParameters, want '%s' got '%s'", expected,
			searchParameters)
	}

	// Test two values
	searchMap = make(map[string]interface{})
	searchMap["id"] = 10
	searchMap["asn"] = 65536
	expected = "&asn=65536&id=10"
	searchParameters = formatSearchParameters(searchMap)
	if searchParameters != expected {
		t.Errorf("formatSearchParameters, want '%s' got '%s'", expected,
			searchParameters)
	}
}

func TestFormatURL(t *testing.T) {
	var expected string
	var url string

	base := "https://www.peeringdb.com/api/"
	searchMap := make(map[string]interface{})
	searchMap["id"] = 10

	// Test fac namespace with search parameter
	expected = "https://www.peeringdb.com/api/fac?depth=1&id=10"
	url = formatURL(base, facilityNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ix namespace with search parameter
	expected = "https://www.peeringdb.com/api/ix?depth=1&id=10"
	url = formatURL(base, internetExchangeNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ixfac namespace with search parameter
	expected = "https://www.peeringdb.com/api/ixfac?depth=1&id=10"
	url = formatURL(base, internetExchangeFacilityNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ixlan namespace with search parameter
	expected = "https://www.peeringdb.com/api/ixlan?depth=1&id=10"
	url = formatURL(base, internetExchangeLANNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ixpfx namespace with search parameter
	expected = "https://www.peeringdb.com/api/ixpfx?depth=1&id=10"
	url = formatURL(base, internetExchangePrefixNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test net namespace with search parameter
	expected = "https://www.peeringdb.com/api/net?depth=1&id=10"
	url = formatURL(base, networkNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test netfac namespace with search parameter
	expected = "https://www.peeringdb.com/api/netfac?depth=1&id=10"
	url = formatURL(base, networkFacilityNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test netixlan namespace with search parameter
	expected = "https://www.peeringdb.com/api/netixlan?depth=1&id=10"
	url = formatURL(base, networkInternetExchangeLANNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test org namespace with search parameter
	expected = "https://www.peeringdb.com/api/org?depth=1&id=10"
	url = formatURL(base, organizationNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test poc namespace with search parameter
	expected = "https://www.peeringdb.com/api/poc?depth=1&id=10"
	url = formatURL(base, networkContactNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}
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

func TestGetASN(t *testing.T) {
	api := NewAPI()
	expectedASN := 201281
	net, err := api.GetASN(expectedASN)

	if err != nil {
		t.Fail()
		return
	}

	if net.ASN != expectedASN {
		t.Errorf("GetASN, want ASN '%d' got '%d'", expectedASN, net.ASN)
	}
}
