package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mrz1836/go-logger"
)

// CreateServer will create a basic Paymail Server
func CreateServer(c *Configuration) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", c.Port), // Address to run the server on
		Handler:           Handlers(c),                // Load all the routes
		ReadHeaderTimeout: c.Timeout,                  // Basic default timeout for header read requests
		ReadTimeout:       c.Timeout,                  // Basic default timeout for read requests
		WriteTimeout:      c.Timeout,                  // Basic default timeout for write requests
	}
}

// StartServer will run the Paymail server
func StartServer(srv *http.Server) {
	logger.Data(2, logger.DEBUG, "starting go paymail server...", logger.MakeParameter("address", srv.Addr))
	logger.Fatalln(srv.ListenAndServe())
}

// getHost tries its best to return the request host
func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		return removePort(r.Host)
	}
	if len(r.URL.Host) == 0 {
		return removePort(r.Host)
	}
	return r.URL.Host
}

// removePort will attempt to remove the port if found
func removePort(host string) string {
	// Slice off any port information.
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	return host
}
