package main

import "github.com/tonicpow/go-paymail/server"

func main() {

	// initialize the mock database - only for testing
	server.InitMockDatabase()

	// Default server runs on port 3000 and timeout requests after 15 seconds
	config := server.NewConfiguration("test.com", new(serverInterface))
	config.Port = 3003
	config.Timeout = 10
	config.BasicRoutes.Add404Route = true
	config.BasicRoutes.AddHealthRoute = true
	config.BasicRoutes.AddIndexRoute = true
	config.BasicRoutes.AddNotAllowed = true

	server.Start(config)
}

// Example mock implementation
type serverInterface struct{}

func (s *serverInterface) GetPaymailByAlias(alias string) *server.PaymailAddress {
	return server.MockGetPaymailByAlias(alias)
}
