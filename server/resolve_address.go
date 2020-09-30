package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

/*
Incoming Data Object Example:
{
    "senderName": "MrZ",
    "senderHandle": "mrz@moneybutton.com",
    "dt": "2020-04-09T16:08:06.419Z",
    "amount": 551,
    "purpose": "message to receiver",
	"signature": "SIGNATURE-IF-REQUIRED-IN-CONFIG"
}
*/

// resolveAddress will return the payment destination (bitcoin address) for the corresponding paymail address
//
// Specs: http://bsvalias.org/04-01-basic-address-resolution.html
func resolveAddress(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")

	// Start the SenderRequest
	senderRequest := &paymail.SenderRequest{
		Amount:       params.GetUint64("amount"),
		Dt:           params.GetString("dt"),
		Purpose:      params.GetString("purpose"),
		SenderHandle: params.GetString("senderHandle"),
		SenderName:   params.GetString("senderName"),
		Signature:    params.GetString("signature"),
	}

	// Parse, sanitize and basic validation
	_, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, req, "invalid-parameter", "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if domain != paymailDomain {
		ErrorResponse(w, req, "unknown-domain", "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Check for required fields
	if len(senderRequest.SenderHandle) == 0 {
		ErrorResponse(w, req, "missing-sender-handle", "missing required field: senderHandle", http.StatusBadRequest)
		return
	} else if len(senderRequest.Dt) == 0 {
		ErrorResponse(w, req, "missing-dt", "missing required field: dt", http.StatusBadRequest)
		return
	}

	// Validate the timestamp
	// todo: check that it's within X amount of time in the past and not in X amount into the future
	if err := paymail.ValidateTimestamp(senderRequest.Dt); err != nil {
		ErrorResponse(w, req, "invalid-dt", "invalid dt format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation on sender handle
	// todo: check that the sender handle is a real paymail address?
	if err := paymail.ValidatePaymail(senderRequest.SenderHandle); err != nil {
		ErrorResponse(w, req, "invalid-sender-handle", "invalid senderHandle: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Only validate signatures if sender validation is enabled
	if senderValidationEnabled {
		if len(senderRequest.Signature) > 0 {
			// todo: validate the signature against the message components
			if err := senderRequest.Verify("", senderRequest.Signature); err != nil {
				ErrorResponse(w, req, "invalid-signature", "invalid signature: "+err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			ErrorResponse(w, req, "missing-signature", "missing required signature", http.StatusBadRequest)
			return
		}
	}

	// todo: lookup the paymail address in a data-store, database, etc - get the PubKey (return 404 if not found)

	// todo: add caching for fast responses since the pubkey will not change

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, &paymail.Resolution{
		Address:   address,
		Output:    "insert-output-script-here", // todo: insert the actual script
		Signature: "insert-signature-here",     // todo: insert the real signature if required (is this the previous signature from SenderRequest?)
	})
}
