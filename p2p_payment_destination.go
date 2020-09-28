package paymail

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bitcoinsv/bsvd/chaincfg"
	"github.com/bitcoinsv/bsvd/txscript"
	"github.com/bitcoinsv/bsvutil"
	"github.com/go-resty/resty/v2"
)

/*
Example:
{
  "satoshis": 1000100
}
*/

// PaymentRequest is the request body for the P2P payment request
type PaymentRequest struct {
	Satoshis uint64 `json:"satoshis"` // The amount, in Satoshis, that the sender intends to transfer to the receiver
}

// PaymentDestination is the response from the GetP2PPaymentDestination() request
//
// The reference is unique for the payment destination request
type PaymentDestination struct {
	StandardResponse
	Outputs   []*output `json:"outputs"`   // A list of outputs
	Reference string    `json:"reference"` // A reference for the payment, created by the receiver of the transaction
}

// output is returned inside the payment destination response
//
// There can be several outputs in one response based on the amount of satoshis being transferred and
// the rules in place by the Paymail provider
type output struct {
	Address  string `json:"address,omitempty"`  // Hex encoded locking script
	Satoshis uint64 `json:"satoshis,omitempty"` // Number of satoshis for that output
	Script   string `json:"script"`             // Hex encoded locking script
}

// GetP2PPaymentDestination will return list of outputs for the P2P transactions to use
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func (c *Client) GetP2PPaymentDestination(p2pURL, alias, domain string, paymentRequest *PaymentRequest) (response *PaymentDestination, err error) {

	// Require a valid url
	if len(p2pURL) == 0 || !strings.Contains(p2pURL, "https://") {
		err = fmt.Errorf("invalid url: %s", p2pURL)
		return
	}

	// Basic requirements for resolution request
	if paymentRequest == nil {
		err = fmt.Errorf("paymentRequest cannot be nil")
		return
	} else if paymentRequest.Satoshis == 0 {
		err = fmt.Errorf("satoshis is required")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/api/rawtx/{alias}@{domain.tld}
	// https://<host-discovery-target>/api/p2p-payment-destination/{alias}@{domain.tld}
	reqURL := strings.Replace(strings.Replace(p2pURL, "{alias}", alias, -1), "{domain.tld}", domain, -1)

	// Set POST defaults
	c.Resty.SetTimeout(time.Duration(c.Options.PostTimeout) * time.Second)

	// Set the user agent
	req := c.Resty.R().SetBody(paymentRequest).SetHeader("User-Agent", c.Options.UserAgent)

	// Enable tracing
	if c.Options.RequestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Post(reqURL); err != nil {
		return
	}

	// Start response
	response = new(PaymentDestination)

	// Tracing enabled?
	if c.Options.RequestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Test the status code
	response.StatusCode = resp.StatusCode()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {

		// Paymail address not found?
		if response.StatusCode == http.StatusNotFound {
			err = fmt.Errorf("paymail address not found")
		} else {
			je := &JSONError{}
			if err = json.Unmarshal(resp.Body(), je); err != nil {
				return
			}
			err = fmt.Errorf("bad response from paymail provider: code %d, message: %s", response.StatusCode, je.Message)
		}

		return
	}

	// Decode the body of the response
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return
	}

	// Check for a reference number
	if len(response.Reference) == 0 {
		err = fmt.Errorf("missing a returned reference value")
		return
	}

	// No outputs?
	if len(response.Outputs) == 0 {
		err = fmt.Errorf("missing a returned output")
		return
	}

	// Loop all outputs
	for index, out := range response.Outputs {

		// No script returned
		if len(out.Script) == 0 {
			continue
		}

		// Decode the hex string into bytes
		var script []byte
		if script, err = hex.DecodeString(out.Script); err != nil {
			return
		}

		// Extract the components from the script
		var addresses []bsvutil.Address
		if _, addresses, _, err = txscript.ExtractPkScriptAddrs(script, &chaincfg.MainNetParams); err != nil {
			return
		}

		// Missing an address?
		if len(addresses) == 0 {
			err = fmt.Errorf("invalid output script, missing an address")
			return
		}

		// Extract the address from the pubkey hash
		var address *bsvutil.LegacyAddressPubKeyHash
		if address, err = bsvutil.NewLegacyAddressPubKeyHash(addresses[0].ScriptAddress(), &chaincfg.MainNetParams); err != nil {
			return
		}

		// Use the encoded version of the address
		response.Outputs[index].Address = address.EncodeAddress()
	}

	return
}
