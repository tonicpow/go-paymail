package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_GetCapabilities will test the method GetCapabilities()
func TestClient_GetCapabilities(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client := newTestClient(t)

		mockCapabilities(http.StatusOK)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, DefaultBsvAliasVersion, response.BsvAlias)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, true, response.Has(BRFCPki, ""))
	})

	t.Run("successful testnet response", func(t *testing.T) {
		client := newTestClient(t, WithNetwork(Testnet))

		mockCapabilitiesNetwork(http.StatusOK, Testnet)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, DefaultBsvAliasVersion, response.BsvAlias)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, true, response.Has(BRFCPki, ""))
	})

	t.Run("successful stn response", func(t *testing.T) {
		client := newTestClient(t, WithNetwork(STN))

		mockCapabilitiesNetwork(http.StatusOK, STN)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, DefaultBsvAliasVersion, response.BsvAlias)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, true, response.Has(BRFCPki, ""))
	})

	t.Run("status not modified", func(t *testing.T) {
		client := newTestClient(t)

		mockCapabilities(http.StatusNotModified)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, DefaultBsvAliasVersion, response.BsvAlias)
		assert.Equal(t, http.StatusNotModified, response.StatusCode)
		assert.Equal(t, true, response.Has(BRFCPki, ""))
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, 0, len(response.Capabilities))
	})

	t.Run("missing target", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		response, err := client.GetCapabilities("", DefaultPort)
		require.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing port", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		response, err := client.GetCapabilities(testDomain, 0)
		require.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad error in request", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.Error(t, err)
		require.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, 0, len(response.Capabilities))
	})

	t.Run("invalid quotes - good response", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{“`+DefaultServiceName+`“: “`+DefaultBsvAliasVersion+`“,“capabilities“: {“6745385c3fc0“: false,
“pki“: “`+testServerURL+`id/{alias}@{domain.tld}“,“paymentDestination“: “`+testServerURL+`address/{alias}@{domain.tld}“}}`,
			),
		)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, DefaultBsvAliasVersion, response.BsvAlias)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, true, response.Has(BRFCPki, ""))
	})

	t.Run("invalid alias", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusNotModified,
				`{"`+DefaultServiceName+`": "","capabilities": {"6745385c3fc0": false,"pki": "`+testServerURL+`id/{alias}@{domain.tld}",
"paymentDestination": "`+testServerURL+`address/{alias}@{domain.tld}"}}`,
			),
		)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.Error(t, err)
		require.NotNil(t, response)
		assert.NotEqual(t, DefaultBsvAliasVersion, response.BsvAlias)
		assert.Equal(t, http.StatusNotModified, response.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName,
			httpmock.NewStringResponder(
				http.StatusNotModified,
				`{"`+DefaultServiceName+`": ,capabilities: {6745385c3fc0: ,pki: `+testServerURL+`id/{alias}@{domain.tld}",
"paymentDestination": "`+testServerURL+`address/{alias}@{domain.tld}"}}`,
			),
		)

		response, err := client.GetCapabilities(testDomain, DefaultPort)
		require.Error(t, err)
		require.NotNil(t, response)
		assert.Equal(t, http.StatusNotModified, response.StatusCode)
		assert.Equal(t, 0, len(response.Capabilities))
	})
}

// mockCapabilities is used for mocking the response
func mockCapabilities(statusCode int) {
	mockCapabilitiesNetwork(statusCode, Mainnet)
}

// mockCapabilitiesNetwork is used for mocking the response on a specific network
func mockCapabilitiesNetwork(statusCode int, n Network) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://"+testDomain+":443/.well-known/"+DefaultServiceName+n.PaymailURLSuffix(),
		httpmock.NewStringResponder(
			statusCode,
			`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`","capabilities": 
{"6745385c3fc0": false,"pki": "`+testServerURL+`id/{alias}@{domain.tld}",
"paymentDestination": "`+testServerURL+`address/{alias}@{domain.tld}"}}`,
		),
	)
}

// ExampleClient_GetCapabilities example using GetCapabilities()
//
// See more examples in /examples/
func ExampleClient_GetCapabilities() {
	// Load the client
	client := newTestClient(nil)

	mockCapabilities(http.StatusOK)

	// Get the capabilities
	capabilities, err := client.GetCapabilities(testDomain, DefaultPort)
	if err != nil {
		fmt.Printf("error getting capabilities: " + err.Error())
		return
	}
	fmt.Printf("found %d capabilities", len(capabilities.Capabilities))
	// Output:found 3 capabilities
}

// BenchmarkClient_GetCapabilities benchmarks the method GetCapabilities()
func BenchmarkClient_GetCapabilities(b *testing.B) {
	client := newTestClient(nil)
	mockCapabilities(http.StatusOK)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetCapabilities(testDomain, DefaultPort)
	}
}

// TestCapabilities_Has will test the method Has()
func TestCapabilities_Has(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		capabilities  *CapabilitiesPayload
		brfcID        string
		alternateID   string
		expectedFound bool
	}{
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "alternate_id", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "alternate_id", "6745385c3fc0", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "6745385c3fc0", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "wrong", false},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "6745385c3fc0", true},
	}

	for _, test := range tests {
		if output := test.capabilities.Has(test.brfcID, test.alternateID); output != test.expectedFound {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%v] expected, received: [%v]", t.Name(), test.brfcID, test.alternateID, test.expectedFound, output)
		}
	}
}

// ExampleCapabilitiesPayload_Has example using Has()
//
// See more examples in /examples/
func ExampleCapabilitiesPayload_Has() {
	capabilities := &CapabilitiesPayload{
		BsvAlias: DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	found := capabilities.Has("6745385c3fc0", "alternate_id")
	fmt.Printf("found brfc: %v", found)
	// Output:found brfc: true
}

// BenchmarkCapabilities_Has benchmarks the method Has()
func BenchmarkCapabilities_Has(b *testing.B) {
	capabilities := &CapabilitiesPayload{
		BsvAlias: DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	for i := 0; i < b.N; i++ {
		_ = capabilities.Has("6745385c3fc0", "alternate_id")
	}
}

// TestCapabilities_GetBool will test the method GetBool()
func TestCapabilities_GetBool(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		capabilities  *CapabilitiesPayload
		brfcID        string
		alternateID   string
		expectedValue bool
	}{
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "alternate_id", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "alternate_id", "6745385c3fc0", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "6745385c3fc0", "6745385c3fc0", true},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "wrong", false},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": true,
				"alternate_id": true,
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		}, "wrong", "6745385c3fc0", true},
	}

	for _, test := range tests {
		if output := test.capabilities.GetBool(test.brfcID, test.alternateID); output != test.expectedValue {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%v] expected, received: [%v]", t.Name(), test.brfcID, test.alternateID, test.expectedValue, output)
		}
	}
}

// ExampleCapabilitiesPayload_GetBool example using GetBool()
//
// See more examples in /examples/
func ExampleCapabilitiesPayload_GetBool() {
	capabilities := &CapabilitiesPayload{
		BsvAlias: DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	found := capabilities.GetBool("6745385c3fc0", "alternate_id")
	fmt.Printf("found value: %v", found)
	// Output:found value: true
}

// BenchmarkCapabilities_GetBool benchmarks the method GetBool()
func BenchmarkCapabilities_GetBool(b *testing.B) {
	capabilities := &CapabilitiesPayload{
		BsvAlias: DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": true,
			"alternate_id": true,
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	for i := 0; i < b.N; i++ {
		_ = capabilities.GetBool("6745385c3fc0", "alternate_id")
	}
}

// TestCapabilities_GetString will test the method GetString()
func TestCapabilities_GetString(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		capabilities  *CapabilitiesPayload
		brfcID        string
		alternateID   string
		expectedValue string
	}{
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"pki",
			"0c4339ef99c2",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"0c4339ef99c2",
			"pki",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"0c4339ef99c2",
			"0c4339ef99c2",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"pki",
			"",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"wrong",
			"wrong",
			"",
		},
		{&CapabilitiesPayload{
			BsvAlias: DefaultServiceName,
			Capabilities: map[string]interface{}{
				"6745385c3fc0": false,
				"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
				"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			},
		},
			"wrong",
			"pki",
			"https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	for _, test := range tests {
		if output := test.capabilities.GetString(test.brfcID, test.alternateID); output != test.expectedValue {
			t.Errorf("%s Failed: [%s] [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.brfcID, test.alternateID, test.expectedValue, output)
		}
	}
}

// ExampleCapabilitiesPayload_GetString example using GetString()
//
// See more examples in /examples/
func ExampleCapabilitiesPayload_GetString() {
	capabilities := &CapabilitiesPayload{
		BsvAlias: DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": false,
			"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}

	found := capabilities.GetString("pki", "0c4339ef99c2")
	fmt.Printf("found value: %v", found)
	// Output:found value: https://domain.com/bsvalias/id/{alias}@{domain.tld}
}

// BenchmarkCapabilities_GetString benchmarks the method GetString()
func BenchmarkCapabilities_GetString(b *testing.B) {
	capabilities := &CapabilitiesPayload{
		BsvAlias: DefaultServiceName,
		Capabilities: map[string]interface{}{
			"6745385c3fc0": false,
			"pki":          "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
			"0c4339ef99c2": "https://domain.com/" + DefaultServiceName + "/id/{alias}@{domain.tld}",
		},
	}
	for i := 0; i < b.N; i++ {
		_ = capabilities.GetString("pki", "0c4339ef99c2")
	}
}
