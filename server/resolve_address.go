package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
	"github.com/tonicpow/go-paymail"
)

/*
Incoming Data Object Example:
{
    "senderName": "MrZ",
    "senderHandle": "mrz@moneybutton.com",
    "dt": "2020-04-09T16:08:06.419Z",
    "amount": 551,
    "purpose": "message to receiver"
}
*/

// resolveAddress will return the payment destination (bitcoin address) for the corresponding paymail address
//
// Specs: http://bsvalias.org/04-01-basic-address-resolution.html
func resolveAddress(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")
	amount := params.GetUint64("amount")
	dateTime := params.GetString("dt")
	purpose := params.GetString("purpose")
	senderHandle := params.GetString("senderHandle")
	senderName := params.GetString("senderName")

	// Parse, sanitize and basic validation
	_, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, req, "invalid-parameter", "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if domain != paymailDomain {
		ErrorResponse(w, req, "unknown-domain", "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Check required fields
	if len(senderHandle) == 0 {
		ErrorResponse(w, req, "missing-parameter", "missing required field: senderHandle", http.StatusBadRequest)
		return
	} else if len(dateTime) == 0 {
		ErrorResponse(w, req, "missing-parameter", "missing required field: dt", http.StatusBadRequest)
		return
	}

	// Validate the date/time
	if err := paymail.ValidateTimestamp(dateTime); err != nil {
		ErrorResponse(w, req, "invalid-parameter", "invalid format: dt", http.StatusBadRequest)
		return
	}

	// todo: validate senderHandle (must be a valid paymail address)

	// todo: use the additional fields?
	logger.NoFileData(logger.DEBUG, fmt.Sprintf("amount: %d, purpose: %s, senderName: %s", amount, purpose, senderName))

	// todo: lookup the paymail address in a data-store, database, etc - get the PubKey (return 404 if not found)

	// todo: add caching for fast responses since the pubkey will not change

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, &paymail.Resolution{
		Address:   address,
		Output:    "insert-output-script-here", // todo: insert the actual script
		Signature: "insert-signature-here",     // todo: insert the real signature if required
	})
}
