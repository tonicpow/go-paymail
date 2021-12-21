package server

import (
	"time"
)

// Server default values
const (
	DefaultServerPort       = 3000             // Port for the server
	DefaultPaymailDomain    = "localhost"      // Domain (if not set)
	DefaultAPIVersion       = "v1"             // Version of API
	DefaultPrefix           = "https://"       // Paymail specs require SSL
	DefaultSenderValidation = false            // If true, it requires extra sender validation
	DefaultTimeout          = 15 * time.Second // Default timeouts
)

// basicRoutes is the configuration for basic server routes
type basicRoutes struct {
	Add404Route    bool `json:"add_404_route,omitempty"`
	AddHealthRoute bool `json:"add_health_route,omitempty"`
	AddIndexRoute  bool `json:"add_index_route,omitempty"`
	AddNotAllowed  bool `json:"add_not_allowed,omitempty"`
}

// PaymailAddress is an internal struct for paymail addresses and their corresponding information
type PaymailAddress struct {
	Alias       string `json:"alias"`        // Alias or handle of the paymail
	Avatar      string `json:"avatar"`       // This is the url of the user (public profile)
	Domain      string `json:"domain"`       // Domain of the paymail
	ID          uint64 `json:"id"`           // Unique identifier
	LastAddress string `json:"last_address"` // This is used as a temp address for now (should be via xPub)
	Name        string `json:"name"`         // This is the name of the user (public profile)
	PrivateKey  string `json:"private_key"`  // PrivateKey hex encoded
	PubKey      string `json:"pubkey"`       // PublicKey hex encoded
}

// AddressResolutionOutput is an internal struct for the old address resolution
type AddressResolutionOutput struct {
	LastAddress  string `json:"last_address"`  // This is used as a temp address for now (should be via xPub)
	OutputScript string `json:"output_script"` // This is the output script
	Signature    string `json:"signature"`     // This the signature if it was required
}
