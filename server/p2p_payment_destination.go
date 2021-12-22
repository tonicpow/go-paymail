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
  "satoshis": 1000100,
}
*/

// p2pDestination will return an output script(s) for a destination (used with SendP2PTransaction)
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func (c *Configuration) p2pDestination(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(w, req, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: params.GetUint64("satoshis"),
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		ErrorResponse(w, req, ErrorMissingSatoshis, "missing parameter: satoshis", http.StatusBadRequest)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(req, alias, domain, "")
	md.PaymentDestination = paymentRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(req.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(w, req, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(w, req, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Create the response
	var response *paymail.PaymentDestinationInformation
	if response, err = c.actions.CreateP2PDestinationResponse(
		req.Context(), alias, domain, paymentRequest.Satoshis, md,
	); err != nil {
		ErrorResponse(w, req, ErrorScript, "error creating output script(s): "+err.Error(), http.StatusExpectationFailed)
		return
	}

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, response)
}
