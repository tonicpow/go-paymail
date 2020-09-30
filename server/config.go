package server

import "github.com/tonicpow/go-paymail"

// Basic configuration for the server
const (
	serviceURL              = "https://" + paymailDomain + "/v1/" + paymail.DefaultServiceName + "/" // This is appended to all URLs
	paymailDomain           = "test.com"                                                             // This is the primary domain for the paymail service
	senderValidationEnabled = false                                                                  // Turn on if all address resolution requests need a valid signature
)
