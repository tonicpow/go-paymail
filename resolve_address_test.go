package paymail

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

// TestClient_ResolveAddress will test the method ResolveAddress()
func TestClient_ResolveAddress(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err != nil {
		t.Fatalf("error occurred in ResolveAddress: %s", err.Error())
	} else if resolution == nil {
		t.Fatalf("resolution was nil")
	} else if resolution.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusOK)
	}

	// Check if response has output
	if len(resolution.Output) == 0 {
		t.Fatalf("missing output script value")
	}

	// Check that we got an address
	if len(resolution.Address) == 0 {
		t.Fatalf("missing address value")
	}
}

// ExampleClient_ResolveAddress example using ResolveAddress()
//
// See more examples in /examples/
func ExampleClient_ResolveAddress() {
	// Load the client
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://www.moneybutton.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err != nil {
		fmt.Printf("error occurred in ResolveAddress: %s", err.Error())
		return
	}
	if len(resolution.Address) > 0 {
		fmt.Printf("address found!")
		// fmt.Printf("address found: %s", resolution.Address) // Disabled since the address changes often
	}
	// Output:address found!
}

// BenchmarkClient_ResolveAddress benchmarks the method ResolveAddress()
func BenchmarkClient_ResolveAddress(b *testing.B) {
	client, _ := NewClient(nil, nil)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	for i := 0; i < b.N; i++ {
		_, _ = client.ResolveAddress("https://www.moneybutton.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	}
}

// TestClient_ResolveAddressStatusNotModified will test the method ResolveAddress()
func TestClient_ResolveAddressStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err != nil {
		t.Fatalf("error occurred in ResolveAddress: %s", err.Error())
	} else if resolution == nil {
		t.Fatalf("resolution was nil")
	} else if resolution.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusNotModified)
	}

	// Check if response has output
	if len(resolution.Output) == 0 {
		t.Fatalf("missing output script value")
	}

	// Check that we got an address
	if len(resolution.Address) == 0 {
		t.Fatalf("missing address value")
	}
}

// TestClient_ResolveAddressInvalidURL will test the method ResolveAddress()
func TestClient_ResolveAddressInvalidURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("invalid-domain", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should have been nil")
	}
}

// TestClient_ResolveAddressSenderRequestNil will test the method ResolveAddress()
func TestClient_ResolveAddressSenderRequestNil(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", nil)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should have been nil")
	}
}

// TestClient_ResolveAddressSenderRequestDt will test the method ResolveAddress()
func TestClient_ResolveAddressSenderRequestDt(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           "",
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should have been nil")
	}
}

// TestClient_ResolveAddressSenderRequestHandle will test the method ResolveAddress()
func TestClient_ResolveAddressSenderRequestHandle(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should have been nil")
	}
}

// TestClient_ResolveAddressMissingAlias will test the method ResolveAddress()
func TestClient_ResolveAddressMissingAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should have been nil")
	}
}

// TestClient_ResolveAddressMissingDomain will test the method ResolveAddress()
func TestClient_ResolveAddressMissingDomain(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"output": "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should have been nil")
	}
}

// TestClient_ResolveAddressBadRequest will test the method ResolveAddress()
func TestClient_ResolveAddressBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_ResolveAddressHTTPError will test the method ResolveAddress()
func TestClient_ResolveAddressHTTPError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewErrorResponder(fmt.Errorf("error in request")),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution != nil {
		t.Fatalf("resolution should be nil")
	}
}

// TestClient_ResolveAddressBadError will test the method ResolveAddress()
func TestClient_ResolveAddressBadError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": request failed}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_ResolveAddressPaymailNotFound will test the method ResolveAddress()
func TestClient_ResolveAddressPaymailNotFound(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotFound,
			`{"message": "not found"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusNotFound {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusNotFound)
	}
}

// TestClient_ResolveAddressBadJSON will test the method ResolveAddress()
func TestClient_ResolveAddressBadJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"output: 76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusOK)
	}
}

// TestClient_ResolveAddressMissingOutput will test the method ResolveAddress()
func TestClient_ResolveAddressMissingOutput(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"output": ""}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusOK)
	}
}

// TestClient_ResolveAddressInvalidOutput will test the method ResolveAddress()
func TestClient_ResolveAddressInvalidOutput(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"output": "12345678"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusOK)
	}
}

// TestClient_ResolveAddressInvalidHex will test the method ResolveAddress()
func TestClient_ResolveAddressInvalidHex(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"output": "7e00bb007d4960727eb11d92a052502c"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusOK)
	}
}

// TestClient_ResolveAddressInvalidHexLength will test the method ResolveAddress()
func TestClient_ResolveAddressInvalidHexLength(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/address/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"output": "0"}`,
		),
	)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339), // UTC is assumed
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress("https://test.com/api/v1/bsvalias/address/{alias}@{domain.tld}", "mrz", "moneybutton.com", senderRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if resolution == nil {
		t.Fatalf("resolution should have not been nil")
	} else if resolution.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", resolution.StatusCode, http.StatusOK)
	}
}
