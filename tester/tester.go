// Package tester is the testing package
package tester

import (
	"context"
	"fmt"
	"net"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/tonicpow/go-paymail/interfaces"
)

// Resolver for mocking requests
type Resolver struct {
	hosts        map[string][]string
	ipAddresses  map[string][]net.IPAddr
	liveResolver interfaces.DNSResolver
	srvRecords   map[string][]*net.SRV
}

// NewCustomResolver will return a custom resolver with specific records hard coded ,
func NewCustomResolver(liveResolver interfaces.DNSResolver, hosts map[string][]string,
	srvRecords map[string][]*net.SRV, ipAddresses map[string][]net.IPAddr) interfaces.DNSResolver {
	return &Resolver{
		hosts:        hosts,
		ipAddresses:  ipAddresses,
		liveResolver: liveResolver,
		srvRecords:   srvRecords,
	}
}

// LookupHost will lookup a host
func (r *Resolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	records, ok := r.hosts[host]
	if ok {
		return records, nil
	}
	return r.liveResolver.LookupHost(ctx, host)
}

// LookupIPAddr will look up an ip address
func (r *Resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	records, ok := r.ipAddresses[host]
	if ok {
		return records, nil
	}
	return r.liveResolver.LookupIPAddr(ctx, host)
}

// LookupSRV will look up an SRV record
func (r *Resolver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	records, ok := r.srvRecords[service+proto+name]
	if ok {
		if service == "invalid" { // Returns an invalid cname
			return fmt.Sprintf("_%s._%s", service, proto), records, nil
		}
		return fmt.Sprintf("_%s._%s.%s.", service, proto, name), records, nil
	}
	return r.liveResolver.LookupSRV(ctx, service, proto, name)
}

// MockResty will return a mocked Resty client
func MockResty() *resty.Client {

	// Create a Resty Client
	client := resty.New()

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())

	return client
}
