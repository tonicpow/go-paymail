package paymail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
Example:
{
	"hex": "01000000012adda020db81f2155ebba69e7c841275517ebf91674268c32ff2f5c7e2853b2c010000006b483045022100872051ef0b6c47714130c12a067db4f38b988bfc22fe270731c2146f5229386b02207abf68bbf092ec03e2c616defcc4c868ad1fc3cdbffb34bcedfab391a1274f3e412102affe8c91d0a61235a3d07b1903476a2e2f7a90451b2ed592fea9937696a07077ffffffff02ed1a0000000000001976a91491b3753cf827f139d2dc654ce36f05331138ddb588acc9670300000000001976a914da036233873cc6489ff65a0185e207d243b5154888ac00000000",
	"metadata": {
		"note": "Human readable information related to the tx."
		"pubkey": "<somepubkey>",
		"sender": "someone@example.tld",
		"signature": "signature(txid)",
	},
	"reference": "someRefId"
}
*/

// P2PRawTransaction is the request body for the P2P transaction request
type P2PRawTransaction struct {
	Hex       string    `json:"hex"`       // The raw transaction, encoded as a hexadecimal string
	MetaData  *MetaData `json:"metadata"`  // An object containing data associated with the transaction
	Reference string    `json:"reference"` // Reference for the payment (from previous P2P Destination request)
}

// MetaData is an object containing data associated with the transaction
type MetaData struct {
	Note      string `json:"note,omitempty"`      // A human readable bit of information about the payment
	PubKey    string `json:"pubkey,omitempty"`    // Public key to validate the signature (if signature is given)
	Sender    string `json:"sender,omitempty"`    // The paymail of the person that originated the transaction
	Signature string `json:"signature,omitempty"` // A signature of the tx id made by the sender
}

// P2PTransaction is the response to the request
type P2PTransaction struct {
	StandardResponse
	Note string `json:"note"` // Some human readable note
	TxID string `json:"txid"` // The txid of the broadcasted tx
}

// SendP2PTransaction will submit a transaction hex string (tx_hex) to a paymail provider
//
// Specs: https://docs.moneybutton.com/docs/paymail-06-p2p-transactions.html
func (c *Client) SendP2PTransaction(p2pURL, alias, domain string, transaction *P2PRawTransaction) (response *P2PTransaction, err error) {

	// Require a valid url
	if len(p2pURL) == 0 || !strings.Contains(p2pURL, "https://") {
		err = fmt.Errorf("invalid url: %s", p2pURL)
		return
	} else if len(alias) == 0 {
		err = fmt.Errorf("missing alias")
		return
	} else if len(domain) == 0 {
		err = fmt.Errorf("missing domain")
		return
	}

	// Basic requirements for request
	if transaction == nil {
		err = fmt.Errorf("transaction cannot be nil")
		return
	} else if len(transaction.Hex) == 0 {
		err = fmt.Errorf("hex is required")
		return
	} else if len(transaction.Reference) == 0 {
		err = fmt.Errorf("reference is required")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/api/rawtx/{alias}@{domain.tld}
	// https://<host-discovery-target>/api/receive-transaction/{alias}@{domain.tld}
	reqURL := strings.Replace(strings.Replace(p2pURL, "{alias}", alias, -1), "{domain.tld}", domain, -1)

	// Fire the POST request
	var resp StandardResponse
	if resp, err = c.postRequest(reqURL, transaction); err != nil {
		return
	}

	// Start the response
	response = &P2PTransaction{StandardResponse: resp}

	// Test the status code
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {

		// Paymail address not found?
		if response.StatusCode == http.StatusNotFound {
			err = fmt.Errorf("paymail address not found")
		} else {
			serverError := &ServerError{}
			if err = json.Unmarshal(resp.Body, serverError); err != nil {
				return
			}
			err = fmt.Errorf("bad response from paymail provider: code %d, message: %s", response.StatusCode, serverError.Message)
		}

		return
	}

	// Decode the body of the response
	if err = json.Unmarshal(resp.Body, &response); err != nil {
		return
	}

	// Check for a reference number
	if len(response.TxID) == 0 {
		err = fmt.Errorf("missing a returned txid")
		return
	}

	return
}
