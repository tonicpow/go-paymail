package server

import (
	"context"

	"github.com/tonicpow/go-paymail"
)

// Mock implementation of a service provider
type mockServiceProvider struct{}

// GetPaymailByAlias is a demo implementation of this interface
func (m *mockServiceProvider) GetPaymailByAlias(_ context.Context, _, _ string,
	_ *RequestMetadata) (*paymail.AddressInformation, error) {

	// Get the data from the demo database
	return nil, nil
}

// CreateAddressResolutionResponse is a demo implementation of this interface
func (m *mockServiceProvider) CreateAddressResolutionResponse(_ context.Context, _, _ string,
	_ bool, _ *RequestMetadata) (*paymail.ResolutionInformation, error) {

	// Generate a new destination / output for the basic address resolution
	return nil, nil
}

// CreateP2PDestinationResponse is a demo implementation of this interface
func (m *mockServiceProvider) CreateP2PDestinationResponse(_ context.Context, _, _ string,
	_ uint64, _ *RequestMetadata) (*paymail.PaymentDestinationInformation, error) {

	// Generate a new destination for the p2p request
	return nil, nil
}

// RecordTransaction is a demo implementation of this interface
func (m *mockServiceProvider) RecordTransaction(_ context.Context,
	_ *paymail.P2PTransaction, _ *RequestMetadata) (*paymail.P2PTransactionInformation, error) {

	// Record the tx into your datastore layer
	return nil, nil
}
