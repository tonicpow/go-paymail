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
	// This is required first to get the corresponding PKI endpoint url
	var capabilities *paymail.Capabilities
	capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort)
	if err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities:", capabilities)

	// Extract the PKI URL from the capabilities response
	pkiURL := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)

	// Get the actual PKI
	var pki *paymail.PKI
	pki, err = client.GetPKI(pkiURL, "mrz", "moneybutton.com")
	if err != nil {
		log.Fatal("error getting pki: " + err.Error())
	}
	log.Println("found pki:", pki)
}
