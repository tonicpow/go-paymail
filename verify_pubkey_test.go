package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// TestClient_VerifyPubKey will test the method VerifyPubKey()
func TestClient_VerifyPubKey(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err != nil {
		t.Fatalf("error occurred in VerifyPubKey: %s", err.Error())
	} else if verification == nil {
		t.Fatalf("verification was nil")
	} else if verification.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", verification.BsvAlias, DefaultBsvAliasVersion)
	} else if verification.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusOK)
	}

	// Check the handle
	if verification.Handle != "mrz@moneybutton.com" {
		t.Fatalf("handle %s did not match: %s", verification.Handle, "mrz@moneybutton.com")
	}

	// Check the PubKey
	if verification.PubKey != "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10" {
		t.Fatalf("pubkey %s did not match: %s", verification.Handle, "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	}

	// Check the match
	if !verification.Match {
		t.Fatalf("match was false, should be true")
	}

}

// ExampleClient_VerifyPubKey example using VerifyPubKey()
//
// See more examples in /examples/
func ExampleClient_VerifyPubKey() {
	// Load the client
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Verify PubKey
	var verification *Verification
	verification, err = client.VerifyPubKey("https://www.moneybutton.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err != nil {
		fmt.Printf("error getting verification: " + err.Error())
		return
	}
	fmt.Printf("verified %s handle with pubkey: %s", verification.Handle, verification.PubKey)
	// Output:verified mrz@moneybutton.com handle with pubkey: 02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10
}

// BenchmarkClient_VerifyPubKey benchmarks the method VerifyPubKey()
func BenchmarkClient_VerifyPubKey(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.VerifyPubKey("https://www.moneybutton.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	}
}

// TestClient_VerifyPubKeyStatusNotModified will test the method VerifyPubKey()
func TestClient_VerifyPubKeyStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err != nil {
		t.Fatalf("error occurred in VerifyPubKey: %s", err.Error())
	} else if verification == nil {
		t.Fatalf("verification was nil")
	} else if verification.BsvAlias != DefaultBsvAliasVersion {
		t.Fatalf("BsvAlias was: %s and not: %s", verification.BsvAlias, DefaultBsvAliasVersion)
	} else if verification.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusNotModified)
	}
}

// TestClient_VerifyPubKeyMissingURL will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("invalid-url", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification != nil {
		t.Fatalf("verification should be nil")
	}
}

// TestClient_VerifyPubKeyMissingAlias will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification != nil {
		t.Fatalf("verification should be nil")
	}
}

// TestClient_VerifyPubKeyMissingDomain will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingDomain(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification != nil {
		t.Fatalf("verification should be nil")
	}
}

// TestClient_VerifyPubKeyMissingPubKey will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingPubKey(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification != nil {
		t.Fatalf("verification should be nil")
	}
}

// TestClient_VerifyPubKeyBadRequest will test the method VerifyPubKey()
func TestClient_VerifyPubKeyBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_VerifyPubKeyHTTPError will test the method VerifyPubKey()
func TestClient_VerifyPubKeyHTTPError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewErrorResponder(fmt.Errorf("error in request")),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification != nil {
		t.Fatalf("verification should be nil")
	}
}

// TestClient_VerifyPubKeyBadError will test the method VerifyPubKey()
func TestClient_VerifyPubKeyBadError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": request failed}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_VerifyPubKeyBadJSON will test the method VerifyPubKey()
func TestClient_VerifyPubKeyBadJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": 1.0,handle: mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusOK)
	}
}

// TestClient_VerifyPubKeyMissingHandle will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingHandle(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusOK)
	}
}

// TestClient_VerifyPubKeyMissingBsvAlias will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingBsvAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "","handle": "mrz@moneybutton.com","pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusOK)
	}
}

// TestClient_VerifyPubKeyMissingPubKeyResponse will test the method VerifyPubKey()
func TestClient_VerifyPubKeyMissingPubKeyResponse(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusOK)
	}
}

// TestClient_VerifyPubKeyInvalidPubKey will test the method VerifyPubKey()
func TestClient_VerifyPubKeyInvalidPubKey(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, "https://test.com/api/v1/bsvalias/verifypubkey/mrz@moneybutton.com/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"bsvalias": "1.0","handle": "mrz@moneybutton.com","pubkey": "12345678","match": true}`,
		),
	)

	// Fire the request
	var verification *Verification
	verification, err = client.VerifyPubKey("https://test.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}", "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if verification == nil {
		t.Fatalf("verification should not be nil")
	} else if verification.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", verification.StatusCode, http.StatusOK)
	}
}
