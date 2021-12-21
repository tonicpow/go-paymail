package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/tonicpow/go-paymail"
	"github.com/tonicpow/go-paymail/server"
)

// paymailAddressTable is the demo data for the example server (table: paymail_address)
var demoPaymailAddressTable []*server.PaymailAddress

// Create the list of demo aliases to create on load
var demoAliases = []struct {
	alias  string
	domain string
	avatar string
	id     uint64
	name   string
}{
	{"mrz", "test.com", "https://github.com/mrz1836.png", 1, "MrZ"},
	{"mrz", "another.com", "https://github.com/mrz1836.png", 4, "MrZ"},
	{"satchmo", "test.com", "https://github.com/rohenaz.png", 2, "Satchmo"},
	{"siggi", "test.com", "https://github.com/icellan.png", 3, "Siggi"},
}

// InitDemoDatabase creates demo data for the database based on the given aliases
func InitDemoDatabase() error {

	// Generate paymail address records
	for _, demo := range demoAliases {
		if err := generateDemoPaymail(
			demo.alias,
			demo.domain,
			demo.avatar,
			demo.name,
			demo.id,
		); err != nil {
			return fmt.Errorf("failed to create paymail address in demo database for alias: %s id: %d", demo.alias, demo.id)
		}
	}

	return nil
}

// generateDemoPaymail will make a new row in the demo database
//
// NOTE: creates a private key and pubkey
func generateDemoPaymail(alias, domain, avatar, name string, id uint64) (err error) {

	// Start a row
	row := &server.PaymailAddress{
		Alias:  alias,
		Avatar: avatar,
		Domain: domain,
		ID:     id,
		Name:   name,
	}

	// Generate new private key
	if row.PrivateKey, err = bitcoin.CreatePrivateKeyString(); err != nil {
		return
	}

	// Get address
	if row.LastAddress, err = bitcoin.GetAddressFromPrivateKeyString(
		row.PrivateKey, true,
	); err != nil {
		return
	}

	// Derive a pubkey from private key
	if row.PubKey, err = bitcoin.PubKeyFromPrivateKeyString(
		row.PrivateKey, true,
	); err != nil {
		return
	}

	// Add to the table
	demoPaymailAddressTable = append(demoPaymailAddressTable, row)

	return
}

// DemoGetPaymailByAlias will find a paymail address given an alias
func DemoGetPaymailByAlias(alias, domain string) (*server.PaymailAddress, error) {
	for i, row := range demoPaymailAddressTable {
		if strings.EqualFold(alias, row.Alias) && strings.EqualFold(domain, row.Domain) {
			return demoPaymailAddressTable[i], nil
		}
	}
	return nil, nil
}

// DemoCreateAddressResolutionResponse will create a new destination for the address resolution
func DemoCreateAddressResolutionResponse(_ context.Context, alias, domain string,
	senderValidation bool) (*paymail.ResolutionInformation, error) {

	// Get the paymail record
	p, err := DemoGetPaymailByAlias(alias, domain)
	if err != nil {
		return nil, err
	}

	// Start the response
	response := &paymail.ResolutionInformation{}

	// Generate the script
	if response.Output, err = bitcoin.ScriptFromAddress(
		p.LastAddress,
	); err != nil {
		return nil, errors.New("error generating script: " + err.Error())
	}

	// Create a signature of output if senderValidation is enabled
	if senderValidation {
		if response.Signature, err = bitcoin.SignMessage(
			p.PrivateKey, response.Output, false,
		); err != nil {
			return nil, errors.New("invalid signature: " + err.Error())
		}
	}

	return response, nil
}

// DemoCreateP2PDestinationResponse will create a basic resolution response for the demo
func DemoCreateP2PDestinationResponse(ctx context.Context, alias, domain string,
	satoshis uint64) (*paymail.PaymentDestinationInformation, error) {

	// Get the paymail record
	p, err := DemoGetPaymailByAlias(alias, domain)
	if err != nil {
		return nil, err
	}

	// Start the output
	output := &paymail.PaymentOutput{
		Satoshis: satoshis,
	}

	// Generate the script
	if output.Script, err = bitcoin.ScriptFromAddress(
		p.LastAddress,
	); err != nil {
		return nil, err
	}

	// Create the response
	return &paymail.PaymentDestinationInformation{
		Outputs:   []*paymail.PaymentOutput{output},
		Reference: "1234567890", // todo: this should be unique per request
	}, nil
}
