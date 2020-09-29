package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// showPKI will return the public key information for the corresponding paymail address
func showPKI(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")

	// Parse, sanitize and basic validation
	_, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest,
			apirouter.ErrorFromRequest(req, "invalid paymail address: "+incomingPaymail, "invalid paymail: "+incomingPaymail, http.StatusBadRequest, http.StatusBadRequest, incomingPaymail))
		return
	} else if domain != paymailDomain {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest,
			apirouter.ErrorFromRequest(req, "domain unknown: "+domain, "domain unknown: "+domain, http.StatusBadRequest, http.StatusBadRequest, domain))
		return
	}

	// todo: lookup the paymail address in a data-store, database, etc - get the PubKey

	// todo: add caching for fast responses since the pubkey will not change

	// Return the PKI response
	apirouter.ReturnResponse(w, req, http.StatusOK, &paymail.PKI{
		BsvAlias: paymail.DefaultBsvAliasVersion,
		Handle:   address,
		PubKey:   "insert-pubkey-here", // todo: insert the pubkey
	})
}
