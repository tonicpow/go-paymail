package server

import (
	"errors"
	"net/http"

	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// Error codes for server response errors
const (
	ErrorFindingPaymail      = "error-finding-paymail"
	ErrorInvalidDt           = "invalid-dt"
	ErrorInvalidParameter    = "invalid-parameter"
	ErrorInvalidPubKey       = "invalid-pubkey"
	ErrorInvalidSenderHandle = "invalid-sender-handle"
	ErrorInvalidSignature    = "invalid-signature"
	ErrorMethodNotFound      = "method-405"
	ErrorMissingHex          = "missing-hex"
	ErrorMissingReference    = "missing-reference"
	ErrorMissingSatoshis     = "missing-satoshis"
	ErrorPaymailNotFound     = "not-found"
	ErrorRecordingTx         = "error-recording-tx"
	ErrorRequestNotFound     = "request-404"
	ErrorScript              = "script-error"
	ErrorUnknownDomain       = "unknown-domain"
)

var (
	// ErrDomainMissing is the error for missing domain
	ErrDomainMissing = errors.New("domain is missing")

	// ErrServiceProviderNil is the error for having a nil service provider
	ErrServiceProviderNil = errors.New("service provider is nil")

	// ErrPortMissing is when the port is not found
	ErrPortMissing = errors.New("missing a port")

	// ErrServiceNameMissing is when the service name is not found
	ErrServiceNameMissing = errors.New("missing service name")

	// ErrCapabilitiesMissing is when the capabilities struct is nil or not set
	ErrCapabilitiesMissing = errors.New("missing capabilities struct")

	// ErrBsvAliasMissing is when the bsv alias version is missing
	ErrBsvAliasMissing = errors.New("missing bsv alias version")
)

// ErrorResponse is a standard way to return errors to the client
//
// Specs: http://bsvalias.org/99-01-recommendations.html
func ErrorResponse(w http.ResponseWriter, req *http.Request, code, message string, statusCode int) {
	apirouter.ReturnResponse(w, req, statusCode, &paymail.ServerError{Code: code, Message: message})
}
