package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// Capabilities is the standard response for returning the Paymail capabilities
type Capabilities struct {
	BsvAlias     string                 `json:"bsvalias"`     // Version of the bsvalias
	Capabilities map[string]interface{} `json:"capabilities"` // List of the capabilities
}

// genericCapabilities will make generic capabilities
func genericCapabilities(bsvAliasVersion string, senderValidation bool) *Capabilities {
	return &Capabilities{
		BsvAlias: bsvAliasVersion,
		Capabilities: map[string]interface{}{
			paymail.BRFCPaymentDestination:   "/address/{alias}@{domain.tld}",
			paymail.BRFCPki:                  "/id/{alias}@{domain.tld}",
			paymail.BRFCPublicProfile:        "/public-profile/{alias}@{domain.tld}",
			paymail.BRFCSenderValidation:     senderValidation,
			paymail.BRFCVerifyPublicKeyOwner: "/verify-pubkey/{alias}@{domain.tld}/{pubkey}",
		},
	}
}

// showCapabilities will return the service discovery results for the server
// and list all active capabilities of the Paymail server
//
// Specs: http://bsvalias.org/02-02-capability-discovery.html
func (c *Configuration) showCapabilities(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Check the domain (allowed, and used for capabilities response)
	// todo: bake this into middleware? This is protecting the "req" domain name (like CORs)
	domain := getHost(req)
	if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Set the service URL
	apirouter.ReturnResponse(w, req, http.StatusOK, c.EnrichCapabilities(domain))
}
