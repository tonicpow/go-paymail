package paymail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

/*
Default:
{
    "avatar": "https://<domain><image>",
    "name": "<name>"
}
*/

// PublicProfile is the result returned from GetPublicProfile()
type PublicProfile struct {
	StandardResponse
	Avatar string `json:"avatar"` // A URL that returns a 180x180 image. It can accept an optional parameter `s` to return an image of width and height `s`. The image should be JPEG, PNG, or GIF.
	Name   string `json:"name"`   // A string up to 100 characters long. (name or nickname)
}

// GetPublicProfile will return a valid public profile
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-paymail/pull/7/files
func (c *Client) GetPublicProfile(publicProfileURL, alias, domain string) (response *PublicProfile, err error) {

	// Require a valid url
	if len(publicProfileURL) == 0 || !strings.Contains(publicProfileURL, "https://") {
		err = fmt.Errorf("invalid url: %s", publicProfileURL)
		return
	}

	// Basic requirements for resolution request
	if len(alias) == 0 {
		err = fmt.Errorf("missing alias")
		return
	} else if len(domain) == 0 {
		err = fmt.Errorf("missing domain")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/public-profile/{alias}@{domain.tld}
	reqURL := strings.Replace(strings.Replace(publicProfileURL, "{alias}", alias, -1), "{domain.tld}", domain, -1)

	// Set the user agent
	req := c.Resty.R().SetHeader("User-Agent", c.Options.UserAgent)

	// Enable tracing
	if c.Options.RequestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Get(reqURL); err != nil {
		return
	}

	// New struct
	response = new(PublicProfile)

	// Tracing enabled?
	if c.Options.RequestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Test the status code (200 or 304 is valid)
	response.StatusCode = resp.StatusCode()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {
		je := &JSONError{}
		if err = json.Unmarshal(resp.Body(), je); err != nil {
			return
		}
		err = fmt.Errorf("bad response from paymail provider: code %d, message: %s", response.StatusCode, je.Message)
		return
	}

	// Decode the body of the response
	err = json.Unmarshal(resp.Body(), &response)

	return
}
