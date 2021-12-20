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

	configuration.RegisterBasicRoutes(r)
	configuration.RegisterRoutes(r)

	// Return the router
	return r.HTTPRouter
}

// RegisterBasicRoutes register the basic routes to the http router
func (config *Configuration) RegisterBasicRoutes(r *apirouter.Router) {
	// Register basic server routes
	config.registerBasicRoutes(r)
}

// RegisterRoutes register all the available paymail routes to the http router
func (config *Configuration) RegisterRoutes(r *apirouter.Router) {
	// Register paymail routes
	config.registerPaymailRoutes(r)
}

// registerBasicRoutes will register basic server related routes
func (config *Configuration) registerBasicRoutes(router *apirouter.Router) {

	if config.BasicRoutes.AddIndexRoute {
		// Set the main index page (navigating to slash)
		router.HTTPRouter.GET("/", router.Request(index))
		// router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders) // Disabled for security
	}

	if config.BasicRoutes.AddHealthRoute {
		// Set the health request (used for load balancers)
		router.HTTPRouter.GET("/health", router.RequestNoLogging(health))
		router.HTTPRouter.OPTIONS("/health", router.SetCrossOriginHeaders)
		router.HTTPRouter.HEAD("/health", router.SetCrossOriginHeaders)
	}

	if config.BasicRoutes.Add404Route {
		// Set the 404 handler (any request not detected)
		router.HTTPRouter.NotFound = http.HandlerFunc(notFound)
	}

	if config.BasicRoutes.AddNotAllowed {
		// Set the method not allowed
		router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)
	}
}

// registerPaymailRoutes will register all paymail related routes
func (config *Configuration) registerPaymailRoutes(router *apirouter.Router) {

	// Capabilities (service discovery)
	router.HTTPRouter.GET(
		"/.well-known/"+config.ServiceName,
		router.Request(config.showCapabilities),
	)

	// PKI request (public key information)
	router.HTTPRouter.GET(
		"/"+paymailAPIVersion+"/"+config.ServiceName+"/id/:PaymailAddress",
		router.Request(config.showPKI),
	)

	// Verify PubKey request (public key verification to paymail address)
	router.HTTPRouter.GET(
		"/"+paymailAPIVersion+"/"+config.ServiceName+"/verify-pubkey/:PaymailAddress/:pubKey",
		router.Request(config.verifyPubKey),
	)

	// Payment Destination request (address resolution)
	router.HTTPRouter.POST(
		"/"+paymailAPIVersion+"/"+config.ServiceName+"/address/:PaymailAddress",
		router.Request(config.resolveAddress),
	)

	// Public Profile request (returns Name & Avatar)
	router.HTTPRouter.GET(
		"/"+paymailAPIVersion+"/"+config.ServiceName+"/public-profile/:PaymailAddress",
		router.Request(config.publicProfile),
	)

	// P2P Destination request (returns output & reference)
	router.HTTPRouter.POST(
		"/"+paymailAPIVersion+"/"+config.ServiceName+"/p2p-payment-destination/:PaymailAddress",
		router.Request(config.p2pDestination),
	)

	// P2P Receive Tx request (receives the P2P transaction, broadcasts, returns tx_id)
	router.HTTPRouter.POST(
		"/"+paymailAPIVersion+"/"+config.ServiceName+"/receive-transaction/:PaymailAddress",
		router.Request(config.p2pReceiveTx),
	)
}
