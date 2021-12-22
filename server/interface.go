package server

import (
	"context"

	"github.com/tonicpow/go-paymail"
)

// PaymailServiceProvider the paymail server interface that needs to be implemented
type PaymailServiceProvider interface {
	CreateAddressResolutionResponse(ctx context.Context, alias, domain string, senderValidation bool) (*paymail.ResolutionInformation, error)
	CreateP2PDestinationResponse(ctx context.Context, alias, domain string, satoshis uint64) (*paymail.PaymentDestinationInformation, error)
	GetPaymailByAlias(ctx context.Context, alias, domain string) (*PaymailAddress, error)
	RecordTransaction(ctx context.Context, p2pTx *paymail.P2PTransaction) (*paymail.P2PTransactionResponse, error)
}
