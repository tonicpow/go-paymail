package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_GetPKI will test the method GetPKI()
func TestClient_GetPKI(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPKI(http.StatusOK)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.NoError(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, DefaultBsvAliasVersion, pki.BsvAlias)
		assert.Equal(t, http.StatusOK, pki.StatusCode)
		assert.Equal(t, testAlias+"@"+testDomain, pki.Handle)
		assert.Equal(t, "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10", pki.PubKey)
	})

	t.Run("successful response - status not modified", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPKI(http.StatusNotModified)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.NoError(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, DefaultBsvAliasVersion, pki.BsvAlias)
		assert.Equal(t, http.StatusNotModified, pki.StatusCode)
		assert.Equal(t, testAlias+"@"+testDomain, pki.Handle)
		assert.Equal(t, "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10", pki.PubKey)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusBadRequest, pki.StatusCode)
	})

	t.Run("bad error", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusBadRequest, pki.StatusCode)
	})

	t.Run("invalid alias", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"`+DefaultServiceName+`": "","handle": "`+testAlias+`@`+testDomain+`","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusOK, pki.StatusCode)
		assert.Equal(t, "", pki.BsvAlias)
	})

	t.Run("invalid json", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`","handle": 1,pubkey: 02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusOK, pki.StatusCode)
		assert.Equal(t, "", pki.BsvAlias)
		assert.Equal(t, "", pki.PubKey)
	})

	t.Run("returned incorrect handle", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`","handle": "invalid@`+testDomain+`","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusOK, pki.StatusCode)
		assert.NotEqual(t, testAlias+"@"+testDomain, pki.Handle)
	})

	t.Run("missing pubkey", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`","handle": "`+testAlias+`@`+testDomain+`","pubkey": ""}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusOK, pki.StatusCode)
		assert.Equal(t, testAlias+"@"+testDomain, pki.Handle)
		assert.Equal(t, "", pki.PubKey)
	})

	t.Run("invalid pubkey length", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`",
"handle": "`+testAlias+`@`+testDomain+`","pubkey": "wrong-length"}`,
			),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, pki)
		assert.Equal(t, http.StatusOK, pki.StatusCode)
		assert.Equal(t, testAlias+"@"+testDomain, pki.Handle)
		assert.NotEqual(t, PubKeyLength, len(pki.PubKey))
	})

	t.Run("invalid url", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPKI(http.StatusOK)

		pki, err := client.GetPKI("invalid-url", testAlias, testDomain)
		require.Error(t, err)
		require.Nil(t, pki)
	})

	t.Run("missing alias", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPKI(http.StatusOK)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", "", testDomain)
		require.Error(t, err)
		require.Nil(t, pki)
	})

	t.Run("missing domain", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPKI(http.StatusOK)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, "")
		require.Error(t, err)
		require.Nil(t, pki)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.Nil(t, pki)
	})
}

// mockGetPKI is used for mocking the response
func mockGetPKI(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, testServerURL+"id/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			statusCode,
			`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`",
"handle": "`+testAlias+`@`+testDomain+`",
"pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
		),
	)
}

// ExampleClient_GetPKI example using GetPKI()
//
// See more examples in /examples/
func ExampleClient_GetPKI() {
	// Load the client
	client := newTestClient(nil)

	mockGetPKI(http.StatusOK)

	// Get the pki
	pki, err := client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
	if err != nil {
		fmt.Printf("error getting pki: " + err.Error())
		return
	}
	fmt.Printf("found %s handle with pubkey: %s", pki.Handle, pki.PubKey)
	// Output:found mrz@test.com handle with pubkey: 02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10
}

// BenchmarkClient_GetPKI benchmarks the method GetPKI()
func BenchmarkClient_GetPKI(b *testing.B) {
	client := newTestClient(nil)
	mockGetPKI(http.StatusOK)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPKI(testServerURL+"id/{alias}@{domain.tld}", testAlias, testDomain)
	}
}
