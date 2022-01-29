package paymail

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-paymail/tester"
)

// newTestClient will return a client for testing purposes
func newTestClient(t *testing.T, opts ...ClientOps) ClientInterface {

	// Create a Resty Client
	httpClient := tester.MockResty()
	if t != nil {
		require.NotNil(t, httpClient)
	}

	// Create a new client
	client, err := NewClient(append([]ClientOps{WithRequestTracing(), WithDNSTimeout(15 * time.Second)}, opts...)...)
	if t != nil {
		require.NotNil(t, client)
		require.NoError(t, err)
	}

	_ = client.WithCustomHTTPClient(httpClient)

	// Set the customer resolver with known defaults
	r := tester.NewCustomResolver(
		client.GetResolver(),
		map[string][]string{
			testDomain:      {"44.225.125.175", "35.165.117.200", "54.190.182.236"},
			"norecords.com": {},
		},
		map[string][]*net.SRV{
			DefaultServiceName + DefaultProtocol + testDomain:      {{Target: "www." + testDomain, Port: 443, Priority: 10, Weight: 10}},
			"invalid" + DefaultProtocol + testDomain:               {{Target: "www." + testDomain, Port: 443, Priority: 10, Weight: 10}},
			DefaultServiceName + DefaultProtocol + "relayx.io":     {{Target: "relayx.io", Port: 443, Priority: 10, Weight: 10}},
			DefaultServiceName + DefaultProtocol + "norecords.com": {},
		},
		map[string][]net.IPAddr{
			"example.com": {net.IPAddr{IP: net.ParseIP("8.8.8.8"), Zone: "eth0"}},
		},
	)

	// Set the custom resolver
	client.WithCustomResolver(r)
	return client
}

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("default client", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, defaultDNSTimeout, client.GetOptions().dnsTimeout)
		assert.Equal(t, defaultDNSPort, client.GetOptions().dnsPort)
		assert.Equal(t, defaultUserAgent, client.GetOptions().userAgent)
		assert.Equal(t, defaultNameServerNetwork, client.GetOptions().nameServerNetwork)
		assert.Equal(t, defaultNameServer, client.GetOptions().nameServer)
		assert.Equal(t, defaultSSLTimeout, client.GetOptions().sslTimeout)
		assert.Equal(t, defaultSSLDeadline, client.GetOptions().sslDeadline)
		assert.Equal(t, defaultHTTPTimeout, client.GetOptions().httpTimeout)
		assert.Equal(t, defaultRetryCount, client.GetOptions().retryCount)
		assert.Equal(t, false, client.GetOptions().requestTracing)
		assert.NotEqual(t, 0, len(client.GetOptions().brfcSpecs))
		assert.Greater(t, len(client.GetBRFCs()), 6)
	})

	t.Run("custom http client", func(t *testing.T) {
		customHTTPClient := resty.New()
		customHTTPClient.SetTimeout(defaultHTTPTimeout)
		client, err := NewClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
		client.WithCustomHTTPClient(customHTTPClient)
	})

	t.Run("custom dns port", func(t *testing.T) {
		client, err := NewClient(WithDNSPort("54"))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "54", client.GetOptions().dnsPort)
	})

	t.Run("custom http timeout", func(t *testing.T) {
		client, err := NewClient(WithHTTPTimeout(10 * time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 10*time.Second, client.GetOptions().httpTimeout)
	})

	t.Run("custom name server", func(t *testing.T) {
		client, err := NewClient(WithNameServer("9.9.9.9"))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "9.9.9.9", client.GetOptions().nameServer)
	})

	t.Run("custom name server network", func(t *testing.T) {
		client, err := NewClient(WithNameServerNetwork("tcp"))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "tcp", client.GetOptions().nameServerNetwork)
	})

	t.Run("custom retry count", func(t *testing.T) {
		client, err := NewClient(WithRetryCount(3))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 3, client.GetOptions().retryCount)
	})

	t.Run("custom ssl timeout", func(t *testing.T) {
		client, err := NewClient(WithSSLTimeout(7 * time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 7*time.Second, client.GetOptions().sslTimeout)
	})

	t.Run("custom ssl deadline", func(t *testing.T) {
		client, err := NewClient(WithSSLDeadline(7 * time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 7*time.Second, client.GetOptions().sslDeadline)
	})

	t.Run("custom options", func(t *testing.T) {
		client, err := NewClient(WithUserAgent("custom user agent"))
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})

	t.Run("custom resolver", func(t *testing.T) {
		r := tester.NewCustomResolver(nil, nil, nil, nil)
		client, err := NewClient()
		assert.NotNil(t, client)
		assert.NoError(t, err)
		client.WithCustomResolver(r)
	})

	t.Run("no brfcs", func(t *testing.T) {
		var client ClientInterface
		client, err := NewClient(WithBRFCSpecs(nil))
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

// TestClient_GetBRFCs will test the method GetBRFCs()
func TestClient_GetBRFCs(t *testing.T) {
	t.Parallel()

	t.Run("get brfcs", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
		brfcs := client.GetBRFCs()
		assert.Equal(t, 23, len(brfcs))
		assert.Equal(t, "b2aa66e26b43", brfcs[0].ID)
	})
}

// TestClient_GetUserAgent will test the method GetUserAgent()
func TestClient_GetUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("get user agent", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
		userAgent := client.GetUserAgent()
		assert.Equal(t, defaultUserAgent, userAgent)
	})
}

// TestClient_GetResolver will test the method GetResolver()
func TestClient_GetResolver(t *testing.T) {
	t.Parallel()

	t.Run("get resolver", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
		r := client.GetResolver()
		assert.NotNil(t, r)
	})
}

// ExampleNewClient example using NewClient()
//
// See more examples in /examples/
func ExampleNewClient() {
	client, err := NewClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}
	fmt.Printf("loaded client: %s", client.GetOptions().userAgent)
	// Output:loaded client: go-paymail: v0.6.0
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(nil)
	}
}

// TestDefaultClientOptions will test the method defaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	options, err := defaultClientOptions()
	assert.NoError(t, err)
	assert.NotNil(t, options)

	assert.Equal(t, defaultDNSTimeout, options.dnsTimeout)
	assert.Equal(t, defaultDNSPort, options.dnsPort)
	assert.Equal(t, defaultUserAgent, options.userAgent)
	assert.Equal(t, defaultNameServerNetwork, options.nameServerNetwork)
	assert.Equal(t, defaultNameServer, options.nameServer)
	assert.Equal(t, defaultSSLTimeout, options.sslTimeout)
	assert.Equal(t, defaultSSLDeadline, options.sslDeadline)
	assert.Equal(t, defaultHTTPTimeout, options.httpTimeout)
	assert.Equal(t, defaultRetryCount, options.retryCount)
	assert.Equal(t, false, options.requestTracing)
	assert.NotEqual(t, 0, len(options.brfcSpecs))
	assert.Greater(t, len(options.brfcSpecs), 6)
}

// BenchmarkDefaultClientOptions benchmarks the method defaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = defaultClientOptions()
	}
}
