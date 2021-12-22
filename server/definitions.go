package server

import (
	"time"

	"github.com/tonicpow/go-paymail"
)

// Server default values
const (
	DefaultAPIVersion       = "v1"             // Version of API
	DefaultPrefix           = "https://"       // Paymail specs require SSL
	DefaultSenderValidation = false            // If true, it requires extra sender validation
	DefaultServerPort       = 3000             // Port for the server
	DefaultTimeout          = 15 * time.Second // Default timeouts
)

// basicRoutes is the configuration for basic server routes
type basicRoutes struct {
	Add404Route    bool `json:"add_404_route,omitempty"`
	AddHealthRoute bool `json:"add_health_route,omitempty"`
	AddIndexRoute  bool `json:"add_index_route,omitempty"`
	AddNotAllowed  bool `json:"add_not_allowed,omitempty"`
}

// RequestMetadata is the struct with extra metadata
type RequestMetadata struct {
	Alias              string                  `json:"alias,omitempty"`               // Alias of the paymail
	Domain             string                  `json:"domain,omitempty"`              // Domain of the request
	IPAddress          string                  `json:"ip_address,omitempty"`          // IP address of the requesting user
	Note               string                  `json:"note,omitempty"`                // Generic note field used for extra information
	PaymentDestination *paymail.PaymentRequest `json:"payment_destination,omitempty"` // Information from the P2P Payment Destination request
	RequestURI         string                  `json:"request_uri,omitempty"`         // Full requesting URL path
	ResolveAddress     *paymail.SenderRequest  `json:"resolve_address,omitempty"`     // Information from the Resolve Address request
	UserAgent          string                  `json:"user_agent,omitempty"`          // User agent of the requesting user
}
