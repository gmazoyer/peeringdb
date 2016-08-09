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

	base := "https://peeringdb.com/api/"
	searchMap := make(map[string]interface{})
	searchMap["id"] = 10

	// Test fac namespace with search parameter
	expected = "https://peeringdb.com/api/fac?depth=1&id=10"
	url = formatURL(base, facilityNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ix namespace with search parameter
	expected = "https://peeringdb.com/api/ix?depth=1&id=10"
	url = formatURL(base, internetExchangeNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ixfac namespace with search parameter
	expected = "https://peeringdb.com/api/ixfac?depth=1&id=10"
	url = formatURL(base, internetExchangeFacilityNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ixlan namespace with search parameter
	expected = "https://peeringdb.com/api/ixlan?depth=1&id=10"
	url = formatURL(base, internetExchangeLANNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test ixpfx namespace with search parameter
	expected = "https://peeringdb.com/api/ixpfx?depth=1&id=10"
	url = formatURL(base, internetExchangePrefixNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test net namespace with search parameter
	expected = "https://peeringdb.com/api/net?depth=1&id=10"
	url = formatURL(base, networkNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test netfac namespace with search parameter
	expected = "https://peeringdb.com/api/netfac?depth=1&id=10"
	url = formatURL(base, networkFacilityNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test netixlan namespace with search parameter
	expected = "https://peeringdb.com/api/netixlan?depth=1&id=10"
	url = formatURL(base, networkInternetExchangeLANNamepsace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test org namespace with search parameter
	expected = "https://peeringdb.com/api/org?depth=1&id=10"
	url = formatURL(base, organizationNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}

	// Test poc namespace with search parameter
	expected = "https://peeringdb.com/api/poc?depth=1&id=10"
	url = formatURL(base, networkContactNamespace, searchMap)
	if url != expected {
		t.Errorf("formatURL, want '%s' got '%s'", expected, url)
	}
}

func TestNewAPI(t *testing.T) {
	var expectedURL string
	var api *API

	// Test to use the public PeeringDB API
	api = NewAPI()
	expectedURL = "https://peeringdb.com/api/"
	if api.url != expectedURL {
		t.Errorf("formatURL, want '%s' got '%s'", expectedURL, api.url)
	}
}

func TestNewAPIWithAuth(t *testing.T) {
	var expectedURL, expectedLogin, expectedPassword string
	var api *API

	// Test to use the public PeeringDB API with authentication
	api = NewAPIWithAuth("test", "123")
	expectedURL = "https://peeringdb.com/api/"
	expectedLogin = "test"
	expectedPassword = "123"
	if api.url != expectedURL {
		t.Errorf("formatURL, want '%s' got '%s'", expectedURL, api.url)
	}
	if api.authLogin != expectedLogin {
		t.Errorf("formatURL, want '%s' got '%s'", expectedLogin, api.authLogin)
	}
	if api.authPassword != expectedPassword {
		t.Errorf("formatURL, want '%s' got '%s'", expectedPassword, api.authPassword)
	}
}

func TestNewAPIFromURL(t *testing.T) {
	var expectedURL string
	var api *API

	// Test to see if an empty string parameter will force to use the public
	// PeeringDB API.
	api = NewAPIFromURL("")
	expectedURL = "https://peeringdb.com/api/"
	if api.url != expectedURL {
		t.Errorf("formatURL, want '%s' got '%s'", expectedURL, api.url)
	}

	// Test with
	api = NewAPIFromURL("http://localhost/api/")
	expectedURL = "http://localhost/api/"
	if api.url != expectedURL {
		t.Errorf("formatURL, want '%s' got '%s'", expectedURL, api.url)
	}
}

func TestNewAPIFromURLWithAuth(t *testing.T) {
	var expectedURL, expectedLogin, expectedPassword string
	var api *API

	// Test to see if an empty string parameter will force to use the public
	// PeeringDB API.
	api = NewAPIFromURLWithAuth("", "test", "123")
	expectedURL = "https://peeringdb.com/api/"
	expectedLogin = "test"
	expectedPassword = "123"
	if api.url != expectedURL {
		t.Errorf("formatURL, want '%s' got '%s'", expectedURL, api.url)
	}
	if api.authLogin != expectedLogin {
		t.Errorf("formatURL, want '%s' got '%s'", expectedLogin, api.authLogin)
	}
	if api.authPassword != expectedPassword {
		t.Errorf("formatURL, want '%s' got '%s'", expectedPassword, api.authPassword)
	}

	// Test with
	api = NewAPIFromURLWithAuth("http://localhost/api/", "test", "123")
	expectedURL = "http://localhost/api/"
	expectedLogin = "test"
	expectedPassword = "123"
	if api.url != expectedURL {
		t.Errorf("formatURL, want '%s' got '%s'", expectedURL, api.url)
	}
	if api.authLogin != expectedLogin {
		t.Errorf("formatURL, want '%s' got '%s'", expectedLogin, api.authLogin)
	}
	if api.authPassword != expectedPassword {
		t.Errorf("formatURL, want '%s' got '%s'", expectedPassword, api.authPassword)
	}
}

func TestGetASN(t *testing.T) {
	api := NewAPI()
	expectedASN := 29467
	net := api.GetASN(expectedASN)

	if net.ASN != expectedASN {
		t.Errorf("GetASN, want ASN '%d' got '%d'", expectedASN, net.ASN)
	}
}
