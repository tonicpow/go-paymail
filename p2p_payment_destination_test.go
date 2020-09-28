package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// TestClient_GetP2PPaymentDestination will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestination(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err != nil {
		t.Fatalf("error occurred in GetP2PPaymentDestination: %s", err.Error())
	} else if destination == nil {
		t.Fatalf("destination was nil")
	} else if destination.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusOK)
	}

	// Check if response has output
	if len(destination.Outputs) == 0 {
		t.Fatalf("missing output(s) script value")
	}

	// Check that we have satoshis
	if destination.Outputs[0].Satoshis != 100 {
		t.Fatalf("missing satoshis, %d != %d", destination.Outputs[0].Satoshis, 100)
	}

	// Check that we got an address
	if len(destination.Reference) == 0 {
		t.Fatalf("missing reference value")
	}
}

// ExampleClient_GetP2PPaymentDestination example using GetP2PPaymentDestination()
func ExampleClient_GetP2PPaymentDestination() {
	// Load the client
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://www.moneybutton.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err != nil {
		fmt.Printf("error occurred in GetP2PPaymentDestination: %s", err.Error())
		return
	}
	if len(destination.Outputs) > 0 {
		fmt.Printf("payment destination found!")
		// fmt.Printf("reference found: %s", resolution.Reference) // Disabled since the reference changes often
	}
	// Output:payment destination found!
}

// BenchmarkClient_GetP2PPaymentDestination benchmarks the method GetP2PPaymentDestination()
func BenchmarkClient_GetP2PPaymentDestination(b *testing.B) {
	client, _ := NewClient(nil, nil)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	for i := 0; i < b.N; i++ {
		_, _ = client.GetP2PPaymentDestination("https://www.moneybutton.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	}
}

// TestClient_GetP2PPaymentDestinationStatusNotModified will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotModified,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err != nil {
		t.Fatalf("error occurred in GetP2PPaymentDestination: %s", err.Error())
	} else if destination == nil {
		t.Fatalf("destination was nil")
	} else if destination.StatusCode != http.StatusNotModified {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusNotModified)
	}
}

// TestClient_GetP2PPaymentDestinationBadURL will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationBadURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("invalid-url", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination != nil {
		t.Fatalf("destination should be nil")
	}
}

// TestClient_GetP2PPaymentDestinationPaymentNil will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationPaymentNil(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", nil)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination != nil {
		t.Fatalf("destination should be nil")
	}
}

// TestClient_GetP2PPaymentDestinationMissingSatoshis will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationMissingSatoshis(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	// Set the payment request
	paymentRequest := &PaymentRequest{Satoshis: 0}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination != nil {
		t.Fatalf("destination should be nil")
	}
}

// TestClient_GetP2PPaymentDestinationBadRequest will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": "request failed"}`,
		),
	)

	// Set the payment request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetP2PPaymentDestinationHTTPError will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationHTTPError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewErrorResponder(fmt.Errorf("error in request")),
	)

	// Set the payment request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination != nil {
		t.Fatalf("destination should be nil")
	}
}

// TestClient_GetP2PPaymentDestinationAddressNotFound will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationAddressNotFound(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusNotFound,
			`{"message": "not found"}`,
		),
	)

	// Set the payment request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusNotFound {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusNotFound)
	}
}

// TestClient_GetP2PPaymentDestinationBadErrorJSON will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationBadErrorJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{"message": request failed}`,
		),
	)

	// Set the payment request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusBadRequest)
	}
}

// TestClient_GetP2PPaymentDestinationBadJSON will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationBadJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{script: 76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusOK)
	}
}

// TestClient_GetP2PPaymentDestinationMissingReference will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationMissingReference(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": ""}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusOK)
	}
}

// TestClient_GetP2PPaymentDestinationMissingOutputs will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationMissingOutputs(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [],"reference": "12345678"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusOK)
	}
}

// TestClient_GetP2PPaymentDestinationInvalidScript will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationInvalidScript(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "12345678","satoshis": 100}],"reference": "12345678"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusOK)
	}
}

// TestClient_GetP2PPaymentDestinationInvalidHex will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestinationInvalidHex(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, "https://test.com/api/v1/bsvalias/p2p-payment-destination/mrz@moneybutton.com",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"outputs": [{"script": "0","satoshis": 100}],"reference": "12345678"}`,
		),
	)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination("https://test.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}", "mrz", "moneybutton.com", paymentRequest)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if destination == nil {
		t.Fatalf("destination should not be nil")
	} else if destination.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode was: %d and not: %d", destination.StatusCode, http.StatusOK)
	}
}
