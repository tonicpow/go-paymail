package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Load the client
	client, err := paymail.NewClient(nil, nil)
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the capabilities
	var capabilities *paymail.Capabilities
	capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort)
	if err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities:", capabilities)

	// Check if capabilities exist
	found := capabilities.Has(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	log.Println("capabilities found: ", found)
}
