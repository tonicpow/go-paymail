package paymail

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// newTestClient will return a client for testing purposes
func newTestClient() (*Client, error) {
	// Create a Resty Client
	client := resty.New()

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())

	// Set options
	options, err := DefaultClientOptions()
	if err != nil {
		return nil, err
	}

	// Set test options
	options.RequestTracing = true
	options.DNSTimeout = 15

	// Create a new client
	var newClient *Client
	newClient, err = NewClient(options, client, nil)
	if err != nil {
		return nil, err
	}

	// Set the customer resolver with known defaults
	r := newCustomResolver(
		newClient.Resolver,
		map[string][]string{
			"moneybutton.com": {"44.225.125.175", "35.165.117.200", "54.190.182.236"},
			"test.com":        {"44.225.125.175", "35.165.117.200", "54.190.182.236"},
		},
		map[string][]*net.SRV{
			DefaultServiceName + DefaultProtocol + "moneybutton.com": {{Target: "www.moneybutton.com", Port: 443, Priority: 10, Weight: 10}},
			"invalid" + DefaultProtocol + "moneybutton.com":          {{Target: "www.moneybutton.com", Port: 443, Priority: 10, Weight: 10}},
			DefaultServiceName + DefaultProtocol + "relayx.io":       {{Target: "relayx.io", Port: 443, Priority: 10, Weight: 10}},
		},
		map[string][]net.IPAddr{
			"example.com": {net.IPAddr{IP: net.ParseIP("8.8.8.8"), Zone: "eth0"}},
		},
	)

	// Set the custom resolver
	newClient.Resolver = r

	// Return the mocking client
	return newClient, nil
}

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("default client", func(t *testing.T) {
		client, err := NewClient(nil, nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, defaultDNSTimeout, client.Options.DNSTimeout)
		assert.Equal(t, defaultDNSPort, client.Options.DNSPort)
		assert.Equal(t, defaultUserAgent, client.Options.UserAgent)
		assert.Equal(t, defaultNameServerNetwork, client.Options.NameServerNetwork)
		assert.Equal(t, defaultNameServer, client.Options.NameServer)
		assert.Equal(t, defaultSSLTimeout, client.Options.SSLTimeout)
		assert.Equal(t, defaultSSLDeadline, client.Options.SSLDeadline)
		assert.Equal(t, defaultGetTimeout, client.Options.GetTimeout)
		assert.Equal(t, defaultRetryCount, client.Options.RetryCount)
		assert.Equal(t, false, client.Options.RequestTracing)
		assert.Equal(t, defaultPostTimeout, client.Options.PostTimeout)
		assert.NotEqual(t, 0, len(client.Options.BRFCSpecs))
		assert.Greater(t, len(client.Options.BRFCSpecs), 6)
	})

	t.Run("custom http client", func(t *testing.T) {
		customHTTPClient := resty.New()
		customHTTPClient.SetTimeout(defaultGetTimeout * time.Second)
		client, err := NewClient(nil, customHTTPClient, nil)
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})

	t.Run("custom options", func(t *testing.T) {
		options, err := DefaultClientOptions()
		assert.NotNil(t, options)
		assert.NoError(t, err)
		options.UserAgent = "custom user agent"

		var client *Client
		client, err = NewClient(options, nil, nil)
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})

	t.Run("custom resolver", func(t *testing.T) {
		r := newCustomResolver(nil, nil, nil, nil)
		client, err := NewClient(nil, nil, r)
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})

	t.Run("no brfcs", func(t *testing.T) {
		options, err := DefaultClientOptions()
		assert.NoError(t, err)
		assert.NotNil(t, options)

		// Remove the specs (empty)
		options.BRFCSpecs = nil

		var client *Client
		client, err = NewClient(options, nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

// ExampleNewClient example using NewClient()
//
// See more examples in /examples/
func ExampleNewClient() {
	client, err := NewClient(nil, nil, nil)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}
	fmt.Printf("loaded client: %s", client.Options.UserAgent)
	// Output:loaded client: go-paymail: v0.1.0
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(nil, nil, nil)
	}
}

// TestDefaultClientOptions will test the method DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	options, err := DefaultClientOptions()
	assert.NoError(t, err)
	assert.NotNil(t, options)

	assert.Equal(t, defaultDNSTimeout, options.DNSTimeout)
	assert.Equal(t, defaultDNSPort, options.DNSPort)
	assert.Equal(t, defaultUserAgent, options.UserAgent)
	assert.Equal(t, defaultNameServerNetwork, options.NameServerNetwork)
	assert.Equal(t, defaultNameServer, options.NameServer)
	assert.Equal(t, defaultSSLTimeout, options.SSLTimeout)
	assert.Equal(t, defaultSSLDeadline, options.SSLDeadline)
	assert.Equal(t, defaultGetTimeout, options.GetTimeout)
	assert.Equal(t, defaultRetryCount, options.RetryCount)
	assert.Equal(t, false, options.RequestTracing)
	assert.Equal(t, defaultPostTimeout, options.PostTimeout)
	assert.NotEqual(t, 0, len(options.BRFCSpecs))
	assert.Greater(t, len(options.BRFCSpecs), 6)
}

// ExampleDefaultClientOptions example using DefaultClientOptions()
//
// See more examples in /examples/
func ExampleDefaultClientOptions() {
	options, err := DefaultClientOptions()
	if err != nil {
		fmt.Printf("error loading options: %s", err.Error())
		return
	}
	fmt.Printf("loaded options: %s", options.UserAgent)
	// Output:loaded options: go-paymail: v0.1.0
}

// BenchmarkDefaultClientOptions benchmarks the method DefaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DefaultClientOptions()
	}
}
