package main

import (
	"context"

	"github.com/tonicpow/go-paymail"
	"github.com/tonicpow/go-paymail/server"
)

// Example demo implementation of a service provider
type demoServiceProvider struct {
	// Extend your dependencies or custom values
}

// GetPaymailByAlias is a demo implementation of this interface
func (d *demoServiceProvider) GetPaymailByAlias(_ context.Context, alias, domain string) (*server.PaymailAddress, error) {

	// Get the data from the demo database
	return DemoGetPaymailByAlias(alias, domain)
}

// CreateAddressResolutionResponse is a demo implementation of this interface
func (d *demoServiceProvider) CreateAddressResolutionResponse(ctx context.Context, alias, domain string,
	senderValidation bool) (*paymail.ResolutionInformation, error) {

	// Generate a new destination / output for the basic address resolution
	return DemoCreateAddressResolutionResponse(ctx, alias, domain, senderValidation)
}

// CreateP2PDestinationResponse is a demo implementation of this interface
func (d *demoServiceProvider) CreateP2PDestinationResponse(ctx context.Context, alias, domain string,
	satoshis uint64) (*paymail.PaymentDestinationInformation, error) {

	// Generate a new destination for the p2p request
	return DemoCreateP2PDestinationResponse(ctx, alias, domain, satoshis)
}

// RecordTransaction is a demo implementation of this interface
func (d *demoServiceProvider) RecordTransaction(ctx context.Context,
	p2pTx *paymail.P2PTransaction) (*paymail.P2PTransactionResponse, error) {

	// Record the tx into your datastore layer
	return DemoRecordTransaction(ctx, p2pTx)
}
