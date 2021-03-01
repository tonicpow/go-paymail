package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Create a client with options
	client, err := paymail.NewClient(nil, nil, nil)
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Load additional specification(s)
	additionalSpec := `[{"author": "andy (nChain)","id": "57dd1f54fc67","title": "BRFC Specifications","url": "http://bsvalias.org/01-02-brfc-id-assignment.html","version": "1"}]`
	if err = client.Options.LoadBRFCs(additionalSpec); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	log.Printf("total specifications loaded: %d", len(client.Options.brfcSpecs))
}
