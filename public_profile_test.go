package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// TestClient_GetPublicProfile will test the method GetPublicProfile()
func TestClient_GetPublicProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"name": "MrZ","avatar": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://test.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err != nil {
		t.Fatalf("error occurred in GetPublicProfile: %s", err.Error())
	} else if profile == nil {
		t.Fatalf("profile was nil")
	} else if profile.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", profile.StatusCode, http.StatusOK)
	}

	// Check the name
	if len(profile.Name) == 0 {
		t.Fatalf("name was empty")
	}

	// Check the avatar
	if len(profile.Avatar) == 0 {
		t.Fatalf("avatar was empty")
	}
}

// ExampleClient_GetPublicProfile example using GetPublicProfile()
func ExampleClient_GetPublicProfile() {
	// Load the client
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://www.moneybutton.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err != nil {
		fmt.Printf("error getting profile: " + err.Error())
		return
	}
	fmt.Printf("found profile for: %s", profile.Name)
	// Output:found profile for: MrZ
}

// BenchmarkClient_GetPublicProfile benchmarks the method GetPublicProfile()
func BenchmarkClient_GetPublicProfile(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPublicProfile("https://www.moneybutton.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	}
}

// TestClient_GetPublicProfileStatusNotModified will test the method GetPublicProfile()
func TestClient_GetPublicProfileStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"name": "MrZ","avatar": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://test.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err != nil {
		t.Fatalf("error occurred in GetPublicProfile: %s", err.Error())
	} else if profile == nil {
		t.Fatalf("profile was nil")
	} else if profile.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", profile.StatusCode, http.StatusNotModified)
	}
}

// TestClient_GetPublicProfileMissingURL will test the method GetPublicProfile()
func TestClient_GetPublicProfileMissingURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"name": "MrZ","avatar": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("invalid-url", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if profile != nil {
		t.Fatalf("profile should be nil")
	}
}

// TestClient_GetPublicProfileMissingAlias will test the method GetPublicProfile()
func TestClient_GetPublicProfileMissingAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"name": "MrZ","avatar": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://test.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if profile != nil {
		t.Fatalf("profile should be nil")
	}
}

// TestClient_GetPublicProfileMissingDomain will test the method GetPublicProfile()
func TestClient_GetPublicProfileMissingDomain(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"name": "MrZ","avatar": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://test.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if profile != nil {
		t.Fatalf("profile should be nil")
	}
}

// TestClient_GetPublicProfileBadRequest will test the method GetPublicProfile()
func TestClient_GetPublicProfileBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://test.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if profile != nil && profile.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", profile.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetPublicProfileBadJSON will test the method GetPublicProfile()
func TestClient_GetPublicProfileBadJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/public-profile/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"name": MrZ,avatar: https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
		),
	)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile("https://test.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if profile != nil && profile.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", profile.StatusCode, http.StatusBadRequest)
	}
}
