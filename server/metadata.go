package server

import (
	"net/http"

	apirouter "github.com/mrz1836/go-api-router"
)

// CreateMetadata will create the base metadata using the request
func CreateMetadata(req *http.Request, alias, domain, optionalNote string) *RequestMetadata {
	return &RequestMetadata{
		Alias:      alias,
		Domain:     domain,
		IPAddress:  apirouter.GetClientIPAddress(req),
		Note:       optionalNote,
		RequestURI: req.RequestURI,
		UserAgent:  req.UserAgent(),
	}
}
