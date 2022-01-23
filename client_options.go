package paymail

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tonicpow/go-paymail/interfaces"
)

// defaultClientOptions will return an ClientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() (opts *ClientOptions, err error) {
	// Set the default options
	opts = &ClientOptions{
		dnsPort:           defaultDNSPort,
		dnsTimeout:        defaultDNSTimeout,
		httpTimeout:       defaultHTTPTimeout,
		nameServer:        defaultNameServer,
		nameServerNetwork: defaultNameServerNetwork,
		requestTracing:    false,
		retryCount:        defaultRetryCount,
		sslDeadline:       defaultSSLDeadline,
		sslTimeout:        defaultSSLTimeout,
		userAgent:         defaultUserAgent,
	}

	// Load the default BRFC specs
	opts.brfcSpecs, err = LoadBRFCs("")
	return
}

// WithDNSPort can be supplied with a custom dns port to perform SRV checks on.
// Default is 53.
func WithDNSPort(port string) ClientOps {
	return func(c *ClientOptions) {
		c.dnsPort = port
	}
}

// WithDNSTimeout can be supplied to overwrite the default dns srv check timeout.
// The default is 5 seconds.
func WithDNSTimeout(timeout time.Duration) ClientOps {
	return func(c *ClientOptions) {
		c.dnsTimeout = timeout
	}
}

// WithBRFCSpecs allows custom specs to be supplied to extend or replace the defaults.
func WithBRFCSpecs(specs []*BRFCSpec) ClientOps {
	return func(c *ClientOptions) {
		c.brfcSpecs = specs
	}
}

// WithHTTPTimeout can be supplied to adjust the default http client timeouts.
// The http client is used when querying paymail services for capabilities
// Default timeout is 20 seconds.
func WithHTTPTimeout(timeout time.Duration) ClientOps {
	return func(c *ClientOptions) {
		c.httpTimeout = timeout
	}
}

// WithNameServer can be supplied to overwrite the default name server used to resolve srv requests.
// default is 8.8.8.8.
func WithNameServer(ip string) ClientOps {
	return func(c *ClientOptions) {
		c.nameServer = ip
	}
}

// WithNameServerNetwork can overwrite the default network protocol to use.
// The default is udp.
func WithNameServerNetwork(network string) ClientOps {
	return func(c *ClientOptions) {
		c.nameServerNetwork = network
	}
}

// WithRequestTracing will enable tracing.
// Tracing is disabled by default.
func WithRequestTracing() ClientOps {
	return func(c *ClientOptions) {
		c.requestTracing = true
	}
}

// WithRetryCount will overwrite the default retry count for http requests.
// Default retries is 2.
func WithRetryCount(retries int) ClientOps {
	return func(c *ClientOptions) {
		c.retryCount = retries
	}
}

// WithSSLTimeout will overwrite the default ssl timeout.
// Default timeout is 10 seconds.
func WithSSLTimeout(timeout time.Duration) ClientOps {
	return func(c *ClientOptions) {
		c.sslTimeout = timeout
	}
}

// WithSSLDeadline will overwrite the default ssl deadline.
// Default is 10 seconds.
func WithSSLDeadline(timeout time.Duration) ClientOps {
	return func(c *ClientOptions) {
		c.sslDeadline = timeout
	}
}

// WithUserAgent will overwrite the default useragent.
// Default is go-paymail + version.
func WithUserAgent(userAgent string) ClientOps {
	return func(c *ClientOptions) {
		c.userAgent = userAgent
	}
}

// WithCustomResolver will allow you to supply a custom  dns resolver,
// useful for testing etc.
func (c *Client) WithCustomResolver(resolver interfaces.DNSResolver) ClientInterface {
	c.resolver = resolver
	return c
}

// WithCustomHTTPClient will overwrite the default client with a custom client.
func (c *Client) WithCustomHTTPClient(client *resty.Client) ClientInterface {
	c.httpClient = client
	return c
}
