package paymail

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tonicpow/go-paymail/interfaces"
)

// Client is the Paymail client configuration and options
type Client struct {
	httpClient *resty.Client          // HTTP client for GET/POST requests
	options    *clientOptions         // Options are all the default settings / configuration
	resolver   interfaces.DNSResolver // Resolver for DNS look ups
}

// ClientOptions holds all the configuration for client requests and default resources
type clientOptions struct {
	brfcSpecs         []*BRFCSpec   // List of BRFC specifications
	dnsPort           string        // Default DNS port for SRV checks
	dnsTimeout        time.Duration // Default timeout in seconds for DNS fetching
	httpTimeout       time.Duration // Default timeout in seconds for GET requests
	nameServer        string        // Default name server for DNS checks
	nameServerNetwork string        // Default name server network
	requestTracing    bool          // If enabled, it will trace the request timing
	retryCount        int           // Default retry count for HTTP requests
	sslDeadline       time.Duration // Default timeout in seconds for SSL deadline
	sslTimeout        time.Duration // Default timeout in seconds for SSL timeout
	userAgent         string        // User agent for all outgoing requests
}

// ClientOps allow functional options to be supplied
// that overwrite default go-paymail client options.
type ClientOps func(c *clientOptions)

// NewClient creates a new client for all paymail requests
//
// If no options are given, it will use the defaultClientOptions()
// If no client is supplied it will use a default Resty HTTP client
func NewClient(opts ...ClientOps) (*Client, error) {

	// Start with the defaults
	defaults, err := defaultClientOptions()
	if err != nil {
		return nil, err
	}

	// Create a new client
	client := &Client{
		options: defaults,
	}

	// Overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}

	// Check for specs (if not set, use the defaults)
	if len(client.options.brfcSpecs) == 0 {
		if err = client.options.LoadBRFCs(""); err != nil {
			return nil, err
		}
	}

	// Set the resolver
	if client.resolver == nil {
		r := client.defaultResolver()
		client.resolver = &r
	}

	// Set the Resty HTTP client
	if client.httpClient == nil {
		client.httpClient = resty.New()

		// Set defaults (for GET requests)
		client.httpClient.SetTimeout(client.options.httpTimeout)
		client.httpClient.SetRetryCount(client.options.retryCount)
	}
	return client, nil
}

// GetBRFCs will return the list of specs
func (c *Client) GetBRFCs() []*BRFCSpec {
	return c.options.brfcSpecs
}

// GetUserAgent will return the user agent string of the client
func (c *Client) GetUserAgent() string {
	return c.options.userAgent
}

// GetResolver will return the internal resolver from the client
func (c *Client) GetResolver() interfaces.DNSResolver {
	return c.resolver
}

// getRequest is a standard GET request for all outgoing HTTP requests
func (c *Client) getRequest(requestURL string) (response StandardResponse, err error) {

	// Set the user agent
	req := c.httpClient.R().SetHeader("User-Agent", c.options.userAgent)

	// Enable tracing
	if c.options.requestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Get(requestURL); err != nil {
		return
	}

	// Tracing enabled?
	if c.options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code
	response.StatusCode = resp.StatusCode()

	// Set the body
	response.Body = resp.Body()
	return
}

// postRequest is a standard POST request for all outgoing HTTP requests
func (c *Client) postRequest(requestURL string, data interface{}) (response StandardResponse, err error) {

	// Set the user agent
	req := c.httpClient.R().SetBody(data).SetHeader("User-Agent", c.options.userAgent)

	// Enable tracing
	if c.options.requestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Post(requestURL); err != nil {
		return
	}

	// Tracing enabled?
	if c.options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code
	response.StatusCode = resp.StatusCode()

	// Set the body
	response.Body = resp.Body()
	return
}
