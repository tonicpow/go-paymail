package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

// TestClient_GetCapabilities will test the method GetCapabilities()
func TestClient_GetCapabilities(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","capabilities": {"6745385c3fc0": false,"pki": "https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}","paymentDestination": "https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}"}}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err != nil {
		t.Fatalf("error occurred in GetCapabilities: %s", err.Error())
	} else if capabilities == nil {
		t.Fatalf("capabilities was nil")
	} else if capabilities.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", capabilities.BsvAlias, DefaultBsvAliasVersion)
	} else if capabilities.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusOK)
	}

	// Check if response has pki
	if !capabilities.Has(BRFCPki, "") {
		t.Fatalf("%s was not found in the GetCapabilities response", BRFCPki)
	}
}

// ExampleClient_GetCapabilities example using GetCapabilities()
//
// See more examples in /examples/
func ExampleClient_GetCapabilities() {
	// Load the client
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Get the capabilities
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("moneybutton.com", DefaultPort)
	if err != nil {
		fmt.Printf("error getting capabilities: " + err.Error())
		return
	}
	fmt.Printf("found %d capabilities", len(capabilities.Capabilities))
	// Output:found 7 capabilities
}

// BenchmarkClient_GetCapabilities benchmarks the method GetCapabilities()
func BenchmarkClient_GetCapabilities(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetCapabilities("moneybutton.com", DefaultPort)
	}
}

// TestClient_GetCapabilitiesStatusNotModified will test the method GetCapabilities()
func TestClient_GetCapabilitiesStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"bsvalias": "1.0","capabilities": {"6745385c3fc0": false,"pki": "https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}","paymentDestination": "https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}"}}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err != nil {
		t.Fatalf("error occurred in GetCapabilities: %s", err.Error())
	} else if capabilities == nil {
		t.Fatalf("capabilities was nil")
	} else if capabilities.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", capabilities.BsvAlias, DefaultBsvAliasVersion)
	} else if capabilities.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusNotModified)
	}
}

// TestClient_GetCapabilitiesBadRequest will test the method GetCapabilities()
func TestClient_GetCapabilitiesBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities != nil && len(capabilities.Capabilities) > 0 {
		t.Fatalf("capabilities should be empty: %v", capabilities.Capabilities)
	} else if capabilities != nil && capabilities.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetCapabilitiesMissingTarget will test the method GetCapabilities()
func TestClient_GetCapabilitiesMissingTarget(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("", DefaultPort)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities != nil && len(capabilities.Capabilities) > 0 {
		t.Fatalf("capabilities should be empty: %v", capabilities.Capabilities)
	}
}

// TestClient_GetCapabilitiesMissingPort will test the method GetCapabilities()
func TestClient_GetCapabilitiesMissingPort(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", 0)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities != nil && len(capabilities.Capabilities) > 0 {
		t.Fatalf("capabilities should be empty: %v", capabilities.Capabilities)
	}
}

// TestClient_GetCapabilitiesHTTPError will test the method GetCapabilities()
func TestClient_GetCapabilitiesHTTPError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewErrorResponder(fmt.Errorf("error in request")),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities != nil && len(capabilities.Capabilities) > 0 {
		t.Fatalf("capabilities should be empty: %v", capabilities.Capabilities)
	} else if capabilities != nil && capabilities.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetCapabilitiesBadError will test the method GetCapabilities()
func TestClient_GetCapabilitiesBadError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": request failed}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities != nil && len(capabilities.Capabilities) > 0 {
		t.Fatalf("capabilities should be empty: %v", capabilities.Capabilities)
	} else if capabilities != nil && capabilities.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetCapabilitiesInvalidQuotes will test the method GetCapabilities()
func TestClient_GetCapabilitiesInvalidQuotes(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{“bsvalias“: “1.0“,“capabilities“: {“6745385c3fc0“: false,“pki“: “https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}“,“paymentDestination“: “https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}“}}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err != nil {
		t.Fatalf("error occurred in GetCapabilities: %s", err.Error())
	} else if capabilities == nil {
		t.Fatalf("capabilities was nil")
	} else if capabilities.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", capabilities.BsvAlias, DefaultBsvAliasVersion)
	} else if capabilities.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusOK)
	}
}

// TestClient_GetCapabilitiesInvalidAlias will test the method GetCapabilities()
func TestClient_GetCapabilitiesInvalidAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"bsvalias": "","capabilities": {"6745385c3fc0": false,"pki": "https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}","paymentDestination": "https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}"}}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities == nil {
		t.Fatalf("capabilities should not be nil")
	} else if capabilities.BsvAlias == DefaultBsvAliasVersion {
		t.Fatalf("capabilities version should not match: %v", DefaultBsvAliasVersion)
	} else if capabilities.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusNotModified)
	}
}

// TestClient_GetCapabilitiesInvalidJSON will test the method GetCapabilities()
func TestClient_GetCapabilitiesInvalidJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com:443/.well-known/bsvalias",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"bsvalias": ,capabilities: {6745385c3fc0: ,pki: https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}","paymentDestination": "https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}"}}`,
		),
	)

	// Fire the request
	var capabilities *Capabilities
	capabilities, err = client.GetCapabilities("test.com", DefaultPort)
	if err == nil {
		t.Fatalf("error should have occurred in GetCapabilities")
	} else if capabilities == nil {
		t.Fatalf("capabilities should not be nil")
	} else if len(capabilities.Capabilities) > 0 {
		t.Fatalf("capabilities count should be zero")
	} else if capabilities.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", capabilities.StatusCode, http.StatusNotModified)
	}
}

// TestCapabilities_Has will test the method Has()
func TestCapabilities_Has(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		capabilities  *Capabilities
		brfcID        string
		alternateID   string
		expectedFound bool
	}{
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "6745385c3fc0", "alternate_id", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "6745385c3fc0", "", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "alternate_id", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "6745385c3fc0", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "wrong", "wrong", false},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "wrong", "6745385c3fc0", true},
	}

	// Test all
	for _, test := range tests {
		if output := test.capabilities.Has(test.brfcID, test.alternateID); output != test.expectedFound {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%v] expected, received: [%v]", t.Name(), test.brfcID, test.alternateID, test.expectedFound, output)
		}
	}
}

// ExampleCapabilities_Has example using Has()
//
// See more examples in /examples/
func ExampleCapabilities_Has() {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}

	found := capabilities.Has("6745385c3fc0", "alternate_id")
	fmt.Printf("found brfc: %v", found)
	// Output:found brfc: true
}

// BenchmarkCapabilities_Has benchmarks the method Has()
func BenchmarkCapabilities_Has(b *testing.B) {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}

	for i := 0; i < b.N; i++ {
		_ = capabilities.Has("6745385c3fc0", "alternate_id")
	}
}

// TestCapabilities_GetBool will test the method GetBool()
func TestCapabilities_GetBool(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		capabilities  *Capabilities
		brfcID        string
		alternateID   string
		expectedValue bool
	}{
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "6745385c3fc0", "alternate_id", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "6745385c3fc0", "", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "alternate_id", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "6745385c3fc0", "6745385c3fc0", true},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "wrong", "wrong", false},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "wrong", "6745385c3fc0", true},
	}

	// Test all
	for _, test := range tests {
		if output := test.capabilities.GetBool(test.brfcID, test.alternateID); output != test.expectedValue {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%v] expected, received: [%v]", t.Name(), test.brfcID, test.alternateID, test.expectedValue, output)
		}
	}
}

// ExampleCapabilities_GetBool example using GetBool()
//
// See more examples in /examples/
func ExampleCapabilities_GetBool() {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}

	found := capabilities.GetBool("6745385c3fc0", "alternate_id")
	fmt.Printf("found value: %v", found)
	// Output:found value: true
}

// BenchmarkCapabilities_GetBool benchmarks the method GetBool()
func BenchmarkCapabilities_GetBool(b *testing.B) {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities:     map[string]interface{}{"6745385c3fc0": true, "alternate_id": true, "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}

	for i := 0; i < b.N; i++ {
		_ = capabilities.GetBool("6745385c3fc0", "alternate_id")
	}
}

// TestCapabilities_GetString will test the method GetString()
func TestCapabilities_GetString(t *testing.T) {

	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		capabilities  *Capabilities
		brfcID        string
		alternateID   string
		expectedValue string
	}{
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "pki", "0c4339ef99c2", "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "0c4339ef99c2", "pki", "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "0c4339ef99c2", "0c4339ef99c2", "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "pki", "", "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "wrong", "wrong", ""},
		{&Capabilities{
			StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
			BsvAlias:         DefaultServiceName,
			Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
		}, "wrong", "pki", "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}

	// Test all
	for _, test := range tests {
		if output := test.capabilities.GetString(test.brfcID, test.alternateID); output != test.expectedValue {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.brfcID, test.alternateID, test.expectedValue, output)
		}
	}
}

// ExampleCapabilities_GetString example using GetString()
//
// See more examples in /examples/
func ExampleCapabilities_GetString() {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}

	found := capabilities.GetString("pki", "0c4339ef99c2")
	fmt.Printf("found value: %v", found)
	// Output:found value: https://domain.com/bsvalias/id/{alias}@{domain.tld}
}

// BenchmarkCapabilities_GetString benchmarks the method GetString()
func BenchmarkCapabilities_GetString(b *testing.B) {
	capabilities := &Capabilities{
		StandardResponse: StandardResponse{StatusCode: http.StatusOK, Tracing: resty.TraceInfo{TotalTime: 200}},
		BsvAlias:         DefaultServiceName,
		Capabilities:     map[string]interface{}{"6745385c3fc0": false, "pki": "https://domain.com/bsvalias/id/{alias}@{domain.tld}", "0c4339ef99c2": "https://domain.com/bsvalias/id/{alias}@{domain.tld}"},
	}
	for i := 0; i < b.N; i++ {
		_ = capabilities.GetString("pki", "0c4339ef99c2")
	}
}
