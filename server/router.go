package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// Handlers is used to isolate loading the routes (used for testing)
func Handlers() *httprouter.Router {

	// Create a new router
	r := apirouter.New()

	// Turned off all CORs - should be accessed outside of a browser
	r.CrossOriginEnabled = false
	r.CrossOriginAllowCredentials = false
	r.CrossOriginAllowOriginAll = false

	// Register all actions
	registerRoutes(r)

	// Return the router
	return r.HTTPRouter
}

// registerRoutes register all the package specific routes
func registerRoutes(router *apirouter.Router) {

	// Set the main index page (navigating to slash)
	router.HTTPRouter.GET("/", router.Request(index))
	// router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders) // Disabled for security

	// Set the health request (used for load balancers)
	router.HTTPRouter.GET("/health", router.RequestNoLogging(health))
	router.HTTPRouter.OPTIONS("/health", router.SetCrossOriginHeaders)
	router.HTTPRouter.HEAD("/health", router.SetCrossOriginHeaders)

	// Set the 404 handler (any request not detected)
	router.HTTPRouter.NotFound = http.HandlerFunc(notFound)

	// Set the method not allowed
	router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)

	// Set the capabilities (service discovery)
	router.HTTPRouter.GET("/.well-known/"+paymail.DefaultServiceName, router.Request(showCapabilities))

	// Set the PKI request (public key information)
	router.HTTPRouter.GET("/v1/"+paymail.DefaultServiceName+"/id/:paymailAddress", router.Request(showPKI))
}
