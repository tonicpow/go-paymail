package paymail

import (
	"context"
	"net"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the paymail client/configuration
type Client struct {
	options *clientOptions `json:"options"` // Options are all the default settings / configuration
}

// ClientOptions holds all the configuration for client requests and default resources
type clientOptions struct {
	brfcSpecs         []*BRFCSpec   `json:"brfc_specs"`          // List of BRFC specifications
	dnsPort           string        `json:"dns_port"`            // Default DNS port for SRV checks
	dnsTimeout        time.Duration `json:"dns_timeout"`         // Default timeout in seconds for DNS fetching
	httpTimeout       time.Duration `json:"get_timeout"`         // Default timeout in seconds for GET requests
	nameServer        string        `json:"name_server"`         // Default name server for DNS checks
	nameServerNetwork string        `json:"name_server_network"` // Default name server network
	requestTracing    bool          `json:"request_tracing"`     // If enabled, it will trace the request timing
	retryCount        int           `json:"retry_count"`         // Default retry count for HTTP requests
	sslDeadline       time.Duration `json:"ssl_deadline"`        // Default timeout in seconds for SSL deadline
	sslTimeout        time.Duration `json:"ssl_timeout"`         // Default timeout in seconds for SSL timeout
	userAgent         string        `json:"user_agent"`          // User agent for all outgoing requests
	resolver          DNSResolver   `json:"-"`
	httpClient        *resty.Client `json:"-"`
}

type ClientOps func(c *clientOptions)

func WithDNSPort(port string) ClientOps {
	return func(c *clientOptions) {
		c.dnsPort = port
	}
}

func WithDNSTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.dnsTimeout = timeout
	}
}
func WithBRFCSpecs(specs []*BRFCSpec) ClientOps {
	return func(c *clientOptions) {
		c.brfcSpecs = specs
	}
}

func WithHttpTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.httpTimeout = timeout
	}
}

func WithNameServer(ip string) ClientOps {
	return func(c *clientOptions) {
		c.nameServer = ip
	}
}
func WithNameServerNetwork(network string) ClientOps {
	return func(c *clientOptions) {
		c.nameServerNetwork = network
	}
}

func WithRequestTracing() ClientOps {
	return func(c *clientOptions) {
		c.requestTracing = true
	}
}
func WithRetryCount(retries int) ClientOps {
	return func(c *clientOptions) {
		c.retryCount = retries
	}
}

func WithSSLTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.sslTimeout = timeout
	}
}

func WithSSLDeadline(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.sslDeadline = timeout
	}
}

func WithUserAgent(userAgent string) ClientOps {
	return func(c *clientOptions) {
		c.userAgent = userAgent
	}
}

func WithCustomResolver(resolver DNSResolver) ClientOps {
	return func(c *clientOptions) {
		c.resolver = resolver
	}
}

func WithCustomHttpClient(client *resty.Client) ClientOps {
	return func(c *clientOptions) {
		c.httpClient = client
	}
}

// defaultClientOptions will return an Options struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() (opts *clientOptions, err error) {
	// Set the default options
	opts = &clientOptions{
		dnsPort:           defaultDNSPort,
		dnsTimeout:        defaultDNSTimeout,
		httpTimeout:       defaultHttpTimeout,
		nameServer:        defaultNameServer,
		nameServerNetwork: defaultNameServerNetwork,
		requestTracing:    false,
		retryCount:        defaultRetryCount,
		sslDeadline:       defaultSSLDeadline,
		sslTimeout:        defaultSSLTimeout,
		userAgent:         defaultUserAgent,
	}
	// Load the default BRFC specs
	err = opts.LoadBRFCs("")
	return
}

// DNSResolver is a custom resolver interface for testing
type DNSResolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
}

// NewClient creates a new client for all outgoing paymail requests
//
// If no options are given, it will use the DefaultClientOptions()
// If no client is supplied it will use a default Resty HTTP client
func NewClient(opts ...ClientOps) (*Client, error) {
	defaults, err := defaultClientOptions()
	if err != nil {
		return nil, err
	}
	// Create a new client
	client := &Client{
		options: defaults,
	}
	// overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}
	// default brfcs
	if len(client.options.brfcSpecs) == 0 {
		// Check for specs (if not set, use the defaults)
		if err := client.options.LoadBRFCs(""); err != nil {
			return nil, err
		}
	}
	// Set the resolver
	if client.options.resolver == nil {
		r := client.defaultResolver()
		client.options.resolver = &r
	}
	// Set the Resty HTTP client
	if client.options.httpClient == nil {
		client.options.httpClient = resty.New()
		// Set defaults (for GET requests)
		client.options.httpClient.SetTimeout(time.Duration(client.options.httpTimeout) * time.Second)
		client.options.httpClient.SetRetryCount(client.options.retryCount)
	}
	return client, nil
}

// getRequest is a standard GET request for all outgoing HTTP requests
func (c *Client) getRequest(requestURL string) (response StandardResponse, err error) {
	// Set the user agent
	req := c.Options.httpClient.R().SetHeader("User-Agent", c.Options.userAgent)
	// Enable tracing
	if c.Options.requestTracing {
		req.EnableTrace()
	}
	// Fire the request
	var resp *resty.Response
	if resp, err = req.Get(requestURL); err != nil {
		return
	}
	// Tracing enabled?
	if c.Options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}
	// Set the status code
	response.StatusCode = resp.StatusCode()
	// Set the body
	response.Body = resp.Body()
	return
}

// postRequest is a standard PORT request for all outgoing HTTP requests
func (c *Client) postRequest(requestURL string, data interface{}) (response StandardResponse, err error) {
	// Set the user agent
	req := c.Options.httpClient.R().SetBody(data).SetHeader("User-Agent", c.Options.userAgent)
	// Enable tracing
	if c.Options.requestTracing {
		req.EnableTrace()
	}
	// Fire the request
	var resp *resty.Response
	if resp, err = req.Post(requestURL); err != nil {
		return
	}
	// Tracing enabled?
	if c.Options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}
	// Set the status code
	response.StatusCode = resp.StatusCode()
	// Set the body
	response.Body = resp.Body()
	return
}
