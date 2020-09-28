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
	p2pDestinationURL := capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")
	p2pSendURL := capabilities.GetString(paymail.BRFCP2PTransactions, "")

	// Create the basic paymentRequest to achieve a payment destination (how many sats are you planning to send?)
	paymentRequest := &paymail.PaymentRequest{Satoshis: 1000}

	// Get the p2p destination
	var destination *paymail.PaymentDestination
	destination, err = client.GetP2PPaymentDestination(p2pDestinationURL, "mrz", "moneybutton.com", paymentRequest)
	if err != nil {
		log.Fatal("error getting destination: " + err.Error())
	}
	log.Printf("destination returned reference: %s and outputs: %d", destination.Reference, len(destination.Outputs))

	// Create a new transaction
	rawTransaction := &paymail.P2PRawTransaction{
		Hex: "replace-with-raw-transaction-hex", // todo: replace with a real transaction
		MetaData: &paymail.MetaData{
			Note:   "Thanks for dinner!",
			Sender: "mrz@moneybutton.com",
		},
		Reference: destination.Reference,
	}

	// Send the p2p destination
	var transaction *paymail.P2PTransaction
	transaction, err = client.SendP2PTransaction(p2pSendURL, "satchmo", "moneybutton.com", rawTransaction)
	if err != nil {
		log.Fatal("error sending transaction: " + err.Error())
	}
	log.Printf("transaction sent: %s", transaction.TxID)
}
