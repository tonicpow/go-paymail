package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_GetPublicProfile will test the method GetPublicProfile()
func TestClient_GetPublicProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPublicProfile(http.StatusOK)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, http.StatusOK, profile.StatusCode)
		assert.Equal(t, testName, profile.Name)
		assert.Equal(t, testAvatar, profile.Avatar)
	})

	t.Run("successful response - status not modified", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPublicProfile(http.StatusNotModified)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, http.StatusNotModified, profile.StatusCode)
		assert.Equal(t, testName, profile.Name)
		assert.Equal(t, testAvatar, profile.Avatar)
	})

	t.Run("missing url", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPublicProfile(http.StatusOK)

		profile, err := client.GetPublicProfile("invalid-url", testAlias, testDomain)
		require.Error(t, err)
		require.Nil(t, profile)
	})

	t.Run("missing alias", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPublicProfile(http.StatusOK)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", "", testDomain)
		require.Error(t, err)
		require.Nil(t, profile)
	})

	t.Run("missing domain", func(t *testing.T) {
		client := newTestClient(t)

		mockGetPublicProfile(http.StatusOK)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, "")
		require.Error(t, err)
		require.Nil(t, profile)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, http.StatusBadRequest, profile.StatusCode)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.Nil(t, profile)
	})

	t.Run("error occurred", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, http.StatusBadRequest, profile.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		client := newTestClient(t)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"name": MrZ,avatar: https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
			),
		)

		profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		require.Error(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, http.StatusOK, profile.StatusCode)
		assert.Equal(t, "", profile.Name)
		assert.Equal(t, "", profile.Avatar)
	})
}

// mockGetPublicProfile is used for mocking the response
func mockGetPublicProfile(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			statusCode, `{"name": "`+testName+`","avatar": "`+testAvatar+`"}`,
		),
	)
}

// ExampleClient_GetPublicProfile example using GetPublicProfile()
//
// See more examples in /examples/
func ExampleClient_GetPublicProfile() {
	// Load the client
	client := newTestClient(nil)

	mockGetPublicProfile(http.StatusOK)

	// Get profile
	profile, err := client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
	if err != nil {
		fmt.Printf("error getting profile: " + err.Error())
		return
	}
	fmt.Printf("found profile for: %s", profile.Name)
	// Output:found profile for: MrZ
}

// BenchmarkClient_GetPublicProfile benchmarks the method GetPublicProfile()
func BenchmarkClient_GetPublicProfile(b *testing.B) {
	client := newTestClient(nil)
	mockGetPublicProfile(http.StatusOK)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
	}
}
