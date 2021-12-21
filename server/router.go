package server

import (
	"net/http"

	apirouter "github.com/mrz1836/go-api-router"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
)

// Handlers are used to isolate loading the routes (used for testing)
func Handlers(configuration *Configuration) *nrhttprouter.Router {

	// Create a new router
	r := apirouter.New()

	// Turned off all CORs - should be accessed outside a browser
	r.CrossOriginEnabled = false
	r.CrossOriginAllowCredentials = false
	r.CrossOriginAllowOriginAll = false

	// Register the routes
	configuration.RegisterBasicRoutes(r)
	configuration.RegisterRoutes(r)

	// Return the router
	return r.HTTPRouter
}

// RegisterBasicRoutes register the basic routes to the http router
func (c *Configuration) RegisterBasicRoutes(r *apirouter.Router) {
	c.registerBasicRoutes(r)
}

// RegisterRoutes register all the available paymail routes to the http router
func (c *Configuration) RegisterRoutes(r *apirouter.Router) {
	c.registerPaymailRoutes(r)
}

// registerBasicRoutes will register basic server related routes
func (c *Configuration) registerBasicRoutes(router *apirouter.Router) {

	// Set the main index page (navigating to slash)
	if c.BasicRoutes.AddIndexRoute {
		router.HTTPRouter.GET("/", router.Request(index))
		// router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders) // Disabled for security
	}

	// Set the health request (used for load balancers)
	if c.BasicRoutes.AddHealthRoute {
		router.HTTPRouter.GET("/health", router.RequestNoLogging(health))
		router.HTTPRouter.OPTIONS("/health", router.SetCrossOriginHeaders)
		router.HTTPRouter.HEAD("/health", router.SetCrossOriginHeaders)
	}

	// Set the 404 handler (any request not detected)
	if c.BasicRoutes.Add404Route {
		router.HTTPRouter.NotFound = http.HandlerFunc(notFound)
	}

	// Set the method not allowed
	if c.BasicRoutes.AddNotAllowed {
		router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)
	}
}

// registerPaymailRoutes will register all paymail related routes
func (c *Configuration) registerPaymailRoutes(router *apirouter.Router) {

	// Capabilities (service discovery)
	router.HTTPRouter.GET(
		"/.well-known/"+c.ServiceName,
		router.Request(c.showCapabilities),
	)

	// PKI request (public key information)
	router.HTTPRouter.GET(
		"/"+c.APIVersion+"/"+c.ServiceName+"/id/:PaymailAddress",
		router.Request(c.showPKI),
	)

	// Verify PubKey request (public key verification to paymail address)
	router.HTTPRouter.GET(
		"/"+c.APIVersion+"/"+c.ServiceName+"/verify-pubkey/:PaymailAddress/:pubKey",
		router.Request(c.verifyPubKey),
	)

	// Payment Destination request (address resolution)
	router.HTTPRouter.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/address/:PaymailAddress",
		router.Request(c.resolveAddress),
	)

	// Public Profile request (returns Name & Avatar)
	router.HTTPRouter.GET(
		"/"+c.APIVersion+"/"+c.ServiceName+"/public-profile/:PaymailAddress",
		router.Request(c.publicProfile),
	)

	// P2P Destination request (returns output & reference)
	router.HTTPRouter.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/p2p-payment-destination/:PaymailAddress",
		router.Request(c.p2pDestination),
	)

	// P2P Receive Tx request (receives the P2P transaction, broadcasts, returns tx_id)
	router.HTTPRouter.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/receive-transaction/:PaymailAddress",
		router.Request(c.p2pReceiveTx),
	)
}
