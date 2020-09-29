package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// index basic request to /
func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	apirouter.ReturnResponse(w, req, http.StatusOK, map[string]interface{}{"message": "Welcome to the Paymail Server ✌(◕‿-)✌"})
}

// health is a basic request to return a health response
func health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// notFound handles all 404 requests
func notFound(w http.ResponseWriter, req *http.Request) {
	apirouter.ReturnResponse(w, req, http.StatusNotFound,
		apirouter.ErrorFromRequest(req, fmt.Sprintf("%d occurred: %s", http.StatusNotFound, req.RequestURI), "request not found", http.StatusNotFound, http.StatusNotFound, req.RequestURI))
}

// methodNotAllowed handles all 405 requests
func methodNotAllowed(w http.ResponseWriter, req *http.Request) {
	apirouter.ReturnResponse(w, req, http.StatusMethodNotAllowed,
		apirouter.ErrorFromRequest(req, fmt.Sprintf("%d occurred: %s method: %s", http.StatusMethodNotAllowed, req.RequestURI, req.Method), "method not allowed", http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, req.Method))
}
