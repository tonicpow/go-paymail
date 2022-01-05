package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_SendP2PTransaction will test the method SendP2PTransaction()
func TestClient_SendP2PTransaction(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`,
		),
	)

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, http.StatusOK, transaction.StatusCode)
	assert.NotEqual(t, 0, len(transaction.TxID))
	assert.NotEqual(t, 0, len(transaction.Note))
}

// ExampleClient_SendP2PTransaction example using SendP2PTransaction()
//
// See more examples in /examples/
func ExampleClient_SendP2PTransaction() {
	// Load the client (using a TestClient for this example since a live transaction is not possible)
	client := newTestClient(nil)

	// Create mock response (Using a mocked response since a live transaction is not possible)
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusOK, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567"}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	if err != nil {
		fmt.Printf("error occurred in SendP2PTransaction: %s", err.Error())
		return
	}
	fmt.Printf("transaction was successful: %s", transaction.TxID)
	// Output:transaction was successful: f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1
}

// BenchmarkClient_SendP2PTransaction benchmarks the method SendP2PTransaction()
func BenchmarkClient_SendP2PTransaction(b *testing.B) {
	client := newTestClient(nil)

	// Create response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusOK, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	transaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: testAlias + "@" + testDomain},
		Reference: "1234567",
	}

	for i := 0; i < b.N; i++ {
		_, _ = client.SendP2PTransaction(
			testServerURL+"receive-transaction/{alias}@{domain.tld}",
			testAlias, testDomain, transaction)
	}
}

// TestClient_SendP2PTransactionStatusNotModified will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusNotModified(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567"}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}",
		testAlias, testDomain, rawTransaction,
	)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, http.StatusNotModified, transaction.StatusCode)
}

// TestClient_SendP2PTransactionStatusMissingURL will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusMissingURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		"invalid-url", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusMissingAlias will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusMissingAlias(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", "", testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusMissingDomain will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusMissingDomain(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, "", rawTransaction,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusNilTransaction will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusNilTransaction(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, nil,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusMissingHex will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusMissingHex(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusMissingReference will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusMissingReference(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusHTTPError will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusHTTPError(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewErrorResponder(fmt.Errorf("error in request")))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.Nil(t, transaction)
}

// TestClient_SendP2PTransactionStatusBadRequest will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusBadRequest(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusBadRequest, `{"message": "request failed"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.NotNil(t, transaction)
}

// TestClient_SendP2PTransactionStatusPaymailNotFound will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusPaymailNotFound(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotFound, `{"message": "not found"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.NotNil(t, transaction)
	assert.Equal(t, http.StatusNotFound, transaction.StatusCode)
}

// TestClient_SendP2PTransactionStatusBadErrorJSON will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusBadErrorJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusBadRequest, `{"message: request failed"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.NotNil(t, transaction)
	assert.Equal(t, http.StatusBadRequest, transaction.StatusCode)
}

// TestClient_SendP2PTransactionStatusBadJSON will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusBadJSON(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note:test note",txid":"f3ddfabf7a7a84cfa20016e61df24dff32953d4023a3002cb5a98d6da4ef9bf1"}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.NotNil(t, transaction)
}

// TestClient_SendP2PTransactionStatusMissingTxID will test the method SendP2PTransaction()
func TestClient_SendP2PTransactionStatusMissingTxID(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	// Create a client with options
	client := newTestClient(t)

	// Create mock response
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"receive-transaction/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(http.StatusNotModified, `{"note":"test note","txid":""}`))

	// Raw TX
	rawTransaction := &P2PTransaction{
		Hex:       "some-raw-hex",
		MetaData:  &P2PMetaData{Note: "test note", Sender: "someone@" + testDomain},
		Reference: "1234567",
	}

	// Fire the request
	transaction, err := client.SendP2PTransaction(
		testServerURL+"receive-transaction/{alias}@{domain.tld}", testAlias, testDomain, rawTransaction,
	)
	require.Error(t, err)
	require.NotNil(t, transaction)
	assert.Equal(t, 0, len(transaction.TxID))
}
