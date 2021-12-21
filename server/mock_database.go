package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/mrz1836/go-logger"
)

// paymailAddressTable is the mocked data for the example server (table: paymail_address)
var paymailAddressTable []*PaymailAddress

// Create the list of mock aliases to create on load
var mockAliases = []struct {
	alias  string
	avatar string
	id     uint64
	name   string
}{
	{"mrz", "https://github.com/mrz1836.png", 1, "MrZ"},
	{"satchmo", "https://github.com/rohenaz.png", 2, "Satchmo"},
}

// InitMockDatabase run on load
func InitMockDatabase() {

	// Generate a paymail addresses
	for _, mock := range mockAliases {
		if err := generatePaymail(mock.alias, mock.id); err != nil {
			log.Fatalf("failed to create paymail address in mock database for alias: %s id: %d", mock.alias, mock.id)
		}
	}

	// Log the paymail addresses in database
	logger.Data(2, logger.DEBUG, fmt.Sprintf("found %d paymails in the mock database", len(mockAliases)))
}

// generatePaymail will make a new row in the mock database
// creates a private key and pubkey
func generatePaymail(alias string, id uint64) error {

	// Start a row
	row := &PaymailAddress{ID: id, Alias: alias}

	var err error

	// Generate new private key
	if row.PrivateKey, err = bitcoin.CreatePrivateKeyString(); err != nil {
		return err
	}

	// Get address
	if row.LastAddress, err = bitcoin.GetAddressFromPrivateKeyString(row.PrivateKey, true); err != nil {
		return err
	}

	// Derive a pubkey from private key
	if row.PubKey, err = bitcoin.PubKeyFromPrivateKeyString(row.PrivateKey, true); err != nil {
		return err
	}

	// Load some mock paymail addresses
	paymailAddressTable = append(paymailAddressTable, row)

	return nil
}

// MockGetPaymailByAlias will find a paymail address given an alias
func MockGetPaymailByAlias(alias string) (*PaymailAddress, error) {
	alias = strings.ToLower(alias)
	for i, row := range paymailAddressTable {
		if alias == row.Alias {
			return paymailAddressTable[i], nil
		}
	}
	return nil, nil
}
