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
	// This is required first to get the corresponding P2P PaymentResolution endpoint url
	var capabilities *paymail.Capabilities
	capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort)
	if err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities:", capabilities)

	// Extract the URL from the capabilities response
	p2pURL := capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")

	// Create the basic paymentRequest to achieve a payment destination (how many sats are you planning to send?)
	paymentRequest := &paymail.PaymentRequest{Satoshis: 1000}

	// Get the p2p destination
	var destination *paymail.PaymentDestination
	destination, err = client.GetP2PPaymentDestination(p2pURL, "mrz", "moneybutton.com", paymentRequest)
	if err != nil {
		log.Fatal("error getting destination: " + err.Error())
	}
	log.Printf("destination returned reference: %s and outputs: %d", destination.Reference, len(destination.Outputs))
}
