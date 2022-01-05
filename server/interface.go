package server

import (
	"context"

	"github.com/tonicpow/go-paymail"
)

// PaymailServiceProvider the paymail server interface that needs to be implemented
type PaymailServiceProvider interface {
	CreateAddressResolutionResponse(
		ctx context.Context,
		alias, domain string,
		senderValidation bool,
		metaData *RequestMetadata,
	) (*paymail.ResolutionPayload, error)

	CreateP2PDestinationResponse(
		ctx context.Context,
		alias, domain string,
		satoshis uint64,
		metaData *RequestMetadata,
	) (*paymail.PaymentDestinationPayload, error)

	GetPaymailByAlias(
		ctx context.Context,
		alias, domain string,
		metaData *RequestMetadata,
	) (*paymail.AddressInformation, error)

	RecordTransaction(
		ctx context.Context,
		p2pTx *paymail.P2PTransaction,
		metaData *RequestMetadata,
	) (*paymail.P2PTransactionPayload, error)
}
