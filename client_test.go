package paymail

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
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

	// Return the mocking client with default options
	return NewClient(options, client)
}

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil, nil)
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	if client == nil {
		t.Fatal("failed to load client")
	}

	if client.Options.DNSTimeout != defaultDNSTimeout {
		t.Fatal("defaultDNSTimeout does not match default")
	}

	if client.Options.DNSPort != defaultDNSPort {
		t.Fatal("defaultDNSPort does not match default")
	}

	if client.Options.UserAgent != defaultUserAgent {
		t.Fatal("defaultUserAgent does not match default")
	}

	if client.Options.NameServerNetwork != defaultNameServerNetwork {
		t.Fatal("defaultNameServerNetwork does not match default")
	}

	if client.Options.NameServer != defaultNameServer {
		t.Fatal("defaultNameServer does not match default")
	}

	if client.Options.SSLTimeout != defaultSSLTimeout {
		t.Fatal("defaultSSLTimeout does not match default")
	}

	if client.Options.SSLDeadline != defaultSSLDeadline {
		t.Fatal("defaultSSLDeadline does not match default")
	}

	if client.Options.GetTimeout != defaultGetTimeout {
		t.Fatal("defaultGetTimeout does not match default")
	}

	if client.Options.RetryCount != defaultRetryCount {
		t.Fatal("defaultRetryCount does not match default")
	}

	if client.Options.RequestTracing {
		t.Fatal("RequestTracing should be false by default")
	}

	if client.Options.PostTimeout != defaultPostTimeout {
		t.Fatal("defaultPostTimeout does not match default")
	}

	if client.Resolver.Dial == nil {
		t.Fatal("client.Resolver.Dial was nil")
	}

	if len(client.Options.BRFCSpecs) == 0 {
		t.Fatal("client.Options.BRFCSpecs was empty")
	}

	if len(client.Options.BRFCSpecs) < 6 {
		t.Fatal("client.Options.BRFCSpecs was less than 10 (missing default specs)")
	}
}

// TestDefaultClientOptions will test the method DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	options, err := DefaultClientOptions()
	if err != nil {
		t.Fatalf("error loading options: %s", err.Error())
	}

	if options.DNSTimeout != defaultDNSTimeout {
		t.Fatal("defaultDNSTimeout does not match default")
	}

	if options.DNSPort != defaultDNSPort {
		t.Fatal("defaultDNSPort does not match default")
	}

	if options.UserAgent != defaultUserAgent {
		t.Fatal("defaultUserAgent does not match default")
	}

	if options.NameServerNetwork != defaultNameServerNetwork {
		t.Fatal("defaultNameServerNetwork does not match default")
	}

	if options.NameServer != defaultNameServer {
		t.Fatal("defaultNameServer does not match default")
	}

	if options.SSLTimeout != defaultSSLTimeout {
		t.Fatal("defaultSSLTimeout does not match default")
	}

	if options.SSLDeadline != defaultSSLDeadline {
		t.Fatal("defaultSSLDeadline does not match default")
	}

	if options.GetTimeout != defaultGetTimeout {
		t.Fatal("defaultGetTimeout does not match default")
	}

	if options.RetryCount != defaultRetryCount {
		t.Fatal("defaultRetryCount does not match default")
	}

	if options.PostTimeout != defaultPostTimeout {
		t.Fatal("defaultPostTimeout does not match default")
	}

	if options.RequestTracing {
		t.Fatal("RequestTracing should be false by default")
	}

	if len(options.BRFCSpecs) == 0 {
		t.Fatal("options.BRFCSpecs was empty")
	}

	if len(options.BRFCSpecs) < 6 {
		t.Fatal("options.BRFCSpecs was less than 10 (missing default specs)")
	}

}
