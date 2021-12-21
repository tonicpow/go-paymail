package server

import (
	"context"

	"github.com/tonicpow/go-paymail"
)

// PaymailServiceProvider the paymail server interface that needs to be implemented
type PaymailServiceProvider interface {
	GetPaymailByAlias(ctx context.Context, alias, domain string) (*PaymailAddress, error)
	CreateAddressResolutionResponse(ctx context.Context, alias, domain string, senderValidation bool) (*paymail.ResolutionInformation, error)
	CreateP2PDestinationResponse(ctx context.Context, alias, domain string, satoshis uint64) (*paymail.PaymentDestinationInformation, error)
}
