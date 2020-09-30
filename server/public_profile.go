package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// publicProfile will return the public profile for the corresponding paymail address
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-paymail/pull/7/files
func publicProfile(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")

	// Parse, sanitize and basic validation
	_, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, req, "invalid-parameter", "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if domain != paymailDomain {
		ErrorResponse(w, req, "unknown-domain", "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// todo: lookup the paymail address in a data-store, database, etc - get the Name & Avatar (return 404 if not found)

	// todo: add caching for fast responses since the Name & Avatar don't change often, use dependency keys for cache busting

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, &paymail.PublicProfile{
		Avatar: "insert-avatar-url-here", // todo: insert the image url for the avatar
		Name:   "insert-name-here",       // todo: insert the name of the user
	})
}
