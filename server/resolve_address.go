package server

import (
	"net/http"

	"github.com/bitcoinschema/go-bitcoin"
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
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(w, req, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if domain != paymailDomain {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Check for required fields
	if len(senderRequest.SenderHandle) == 0 {
		ErrorResponse(w, req, ErrorInvalidSenderHandle, "senderHandle is empty", http.StatusBadRequest)
		return
	} else if len(senderRequest.Dt) == 0 {
		ErrorResponse(w, req, ErrorInvalidDt, "dt is empty", http.StatusBadRequest)
		return
	}

	// Validate the timestamp
	if err := paymail.ValidateTimestamp(senderRequest.Dt); err != nil {
		ErrorResponse(w, req, ErrorInvalidDt, "invalid dt: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation on sender handle
	if err := paymail.ValidatePaymail(senderRequest.SenderHandle); err != nil {
		ErrorResponse(w, req, ErrorInvalidSenderHandle, "invalid senderHandle: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Only validate signatures if sender validation is enabled (skip if disabled)
	if senderValidationEnabled {
		if len(senderRequest.Signature) > 0 {

			// todo: get pubkey of senderHandle - PKI, convert pubkey into address for signing

			// todo: validate the signature against the message components (add real keyAddress)
			if err := senderRequest.Verify("", senderRequest.Signature); err != nil {
				ErrorResponse(w, req, ErrorInvalidSignature, "invalid signature: "+err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			ErrorResponse(w, req, ErrorInvalidSignature, "missing required signature", http.StatusBadRequest)
			return
		}
	}

	// todo: lookup the paymail address in a data-store, database, etc - get the PubKey (return 404 if not found)

	// todo: add caching for fast responses since the pubkey will not change

	// Find in mock database
	foundPaymail := getPaymailByAlias(alias)
	if foundPaymail == nil {
		ErrorResponse(w, req, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Start the response
	response := &paymail.Resolution{}

	// Generate the script
	var err error
	response.Output, err = bitcoin.ScriptFromAddress(foundPaymail.LastAddress)
	if err != nil {
		ErrorResponse(w, req, ErrorScript, "error generating script: "+err.Error(), http.StatusNotFound)
		return
	}

	// Create a signature of output if senderValidation is enabled
	if senderValidationEnabled {
		response.Signature, err = bitcoin.SignMessage(foundPaymail.PrivateKey, response.Output)
		if err != nil {
			ErrorResponse(w, req, ErrorInvalidSignature, "invalid signature: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, response)
}
