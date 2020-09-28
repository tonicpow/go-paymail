package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// TestClient_GetPKI will test the method GetPKI()
func TestClient_GetPKI(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err != nil {
		t.Fatalf("error occurred in GetPKI: %s", err.Error())
	} else if pki == nil {
		t.Fatalf("pki was nil")
	} else if pki.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", pki.BsvAlias, DefaultBsvAliasVersion)
	} else if pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}

	// Check if response has pki
	if pki.Handle != "mrz@moneybutton.com" {
		t.Fatalf("handle %s did not match: %s", pki.Handle, "mrz@moneybutton.com")
	}
}

// ExampleClient_GetPKI example using GetPKI()
func ExampleClient_GetPKI() {
	// Load the client
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Get the pki
	var pki *PKI
	pki, err = client.GetPKI("https://www.moneybutton.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err != nil {
		fmt.Printf("error getting pki: " + err.Error())
		return
	}
	fmt.Printf("found %s handle with pubkey: %s", pki.Handle, pki.PubKey)
	// Output:found mrz@moneybutton.com handle with pubkey: 02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10
}

// BenchmarkClient_GetPKI benchmarks the method GetPKI()
func BenchmarkClient_GetPKI(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPKI("https://www.moneybutton.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	}
}

// TestClient_GetPKIStatusNotModified will test the method GetPKI()
func TestClient_GetPKIStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err != nil {
		t.Fatalf("error occurred in GetPKI: %s", err.Error())
	} else if pki == nil {
		t.Fatalf("pki was nil")
	} else if pki.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", pki.BsvAlias, DefaultBsvAliasVersion)
	} else if pki.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}

// TestClient_GetPKIBadRequest will test the method GetPKI()
func TestClient_GetPKIBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetPKIBadError will test the method GetPKI()
func TestClient_GetPKIBadError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": request failed}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetPKIInvalidAlias will test the method GetPKI()
func TestClient_GetPKIInvalidAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}

// TestClient_GetPKIInvalidJSON will test the method GetPKI()
func TestClient_GetPKIInvalidJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": 1,pubkey: 02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}

// TestClient_GetPKIInvalidHandle will test the method GetPKI()
func TestClient_GetPKIInvalidHandle(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "invalid@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	} else if pki != nil && pki.Handle == "mrz@moneybutton.com" {
		t.Fatalf("Handle was: %s and not: %s", pki.Handle, "mrz@moneybutton.com")
	}
}

// TestClient_GetPKIMissingPubKey will test the method GetPKI()
func TestClient_GetPKIMissingPubKey(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": ""}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}

// TestClient_GetPKIInvalidPubKeyLength will test the method GetPKI()
func TestClient_GetPKIInvalidPubKeyLength(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "wrong-length"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}

// TestClient_GetPKIInvalidURL will test the method GetPKI()
func TestClient_GetPKIInvalidURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"}`,
		),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("invalid-url", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}

// TestClient_GetPKIHTTPError will test the method GetPKI()
func TestClient_GetPKIHTTPError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create valid response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/id/mrz@moneybutton.com",
		httpmock.NewErrorResponder(fmt.Errorf("error in request")),
	)

	// Fire the request
	var pki *PKI
	pki, err = client.GetPKI("https://test.com/api/v1/bsvalias/id/{alias}@{domain.tld}", "mrz", "moneybutton.com")
	if err == nil {
		t.Fatalf("error should have occurred in GetPKI")
	} else if pki != nil && pki.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", pki.StatusCode, http.StatusOK)
	}
}
