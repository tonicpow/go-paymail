package paymail

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitcoinschema/go-bitcoin"
)

// TestSenderRequest_Sign will test the method Sign()
func TestSenderRequest_Sign(t *testing.T) {

	// Test private key
	key, err := bitcoin.CreatePrivateKeyString()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
		Purpose:      "testing",
	}

	var signature string

	// Test signing (invalid key)
	if signature, err = senderRequest.Sign(""); err == nil {
		t.Fatalf("error was expected")
	} else if len(signature) > 0 {
		t.Fatalf("signature should be empty")
	}

	// Test signing (invalid key)
	if signature, err = senderRequest.Sign("0"); err == nil {
		t.Fatalf("error was expected")
	} else if len(signature) > 0 {
		t.Fatalf("signature should be empty")
	}

	// Test signing (invalid dt)
	senderRequest.Dt = ""
	if signature, err = senderRequest.Sign(key); err == nil {
		t.Fatalf("error was expected")
	} else if len(signature) > 0 {
		t.Fatalf("signature should be empty")
	}

	// Test signing (invalid senderHandle)
	senderRequest.Dt = time.Now().UTC().Format(time.RFC3339)
	senderRequest.SenderHandle = ""
	if signature, err = senderRequest.Sign(key); err == nil {
		t.Fatalf("error was expected")
	} else if len(signature) > 0 {
		t.Fatalf("signature should be empty")
	}

	// Test signing (valid)
	senderRequest.SenderHandle = "mrz@moneybutton.com"
	if signature, err = senderRequest.Sign(key); err != nil {
		t.Fatalf("error occurred in sign: %s", err.Error())
	} else if len(signature) == 0 {
		t.Fatalf("signature was expected but empty")
	}

	// Get address for verification
	var address string
	if address, err = bitcoin.GetAddressFromPrivateKey(key, false); err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Verify the signature
	if err = senderRequest.Verify(address, signature); err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}

// ExampleSenderRequest_Sign example using Sign()
//
// See more examples in /examples/
func ExampleSenderRequest_Sign() {

	// Test private key
	key := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
		Purpose:      "testing",
	}

	// Sign the sender request
	signature, err := senderRequest.Sign(key)
	if err != nil {
		fmt.Printf("error occurred in sign: %s", err.Error())
		return
	} else if len(signature) == 0 {
		fmt.Printf("signature was empty")
		return
	}

	// Cannot display signature as it changes because of the "dt" field
	fmt.Printf("signature created!")
	// Output:signature created!
}

// BenchmarkSenderRequest_Sign benchmarks the method Sign()
func BenchmarkSenderRequest_Sign(b *testing.B) {

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
		Purpose:      "testing",
	}

	for i := 0; i < b.N; i++ {
		_, _ = senderRequest.Sign("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	}
}

// TestSenderRequest_Verify will test the method Verify()
func TestSenderRequest_Verify(t *testing.T) {

	// Test private key
	key, err := bitcoin.CreatePrivateKeyString()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
		Purpose:      "testing",
	}

	// Sign
	var signature string
	if signature, err = senderRequest.Sign(key); err != nil {
		t.Fatalf("error occurred in sign: %s", err.Error())
	} else if len(signature) == 0 {
		t.Fatalf("signature was expected but empty")
	}

	// Get address from private key
	var address string
	if address, err = bitcoin.GetAddressFromPrivateKey(key, false); err != nil {
		t.Fatalf("error occurred in AddressFromPrivateKey: %s", err.Error())
	}

	// Try verifying (valid)
	if err = senderRequest.Verify(address, signature); err != nil {
		t.Fatalf("error occurred in Verify: %s", err.Error())
	}

	// Try verifying (invalid)
	if err = senderRequest.Verify("", signature); err == nil {
		t.Fatalf("error should have occurred")
	}

	// Try verifying (invalid)
	if err = senderRequest.Verify(address, ""); err == nil {
		t.Fatalf("error should have occurred")
	}

	// Try verifying (invalid)
	if err = senderRequest.Verify(address, "0"); err == nil {
		t.Fatalf("error should have occurred")
	}
}

// ExampleSenderRequest_Verify example using Verify()
//
// See more examples in /examples/
func ExampleSenderRequest_Verify() {

	// Example sender request
	senderRequest := &SenderRequest{
		Dt:           "2020-10-02T16:43:39Z",
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
		Purpose:      "testing",
	}

	// Try verifying (valid) (using an address and a signature - previously generated for example)
	if err := senderRequest.Verify(
		"1MRXps9AaAhHiZwpAvVqaX9J8UAjFhbgGw",
		"IDQWyhXrMrV0++c8lCzp6opWdDkbwEgNjHIOH+TRn9K6fJkOexICLiD9XzLajSlFezHJWgJTicCRv641zhk6rOY=",
	); err != nil {
		fmt.Printf("error occurred in Verify: %s", err.Error())
		return
	}

	fmt.Printf("signature verified!")
	// Output:signature verified!
}

// BenchmarkSenderRequest_Verify benchmarks the method Verify()
func BenchmarkSenderRequest_Verify(b *testing.B) {

	// Example sender request
	senderRequest := &SenderRequest{
		Dt:           "2020-10-02T16:43:39Z",
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
		Purpose:      "testing",
	}

	for i := 0; i < b.N; i++ {
		_ = senderRequest.Verify(
			"1MRXps9AaAhHiZwpAvVqaX9J8UAjFhbgGw",
			"IDQWyhXrMrV0++c8lCzp6opWdDkbwEgNjHIOH+TRn9K6fJkOexICLiD9XzLajSlFezHJWgJTicCRv641zhk6rOY=",
		)
	}
}
