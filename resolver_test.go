package paymail

import (
	"context"
	"fmt"
	"net"
)

// resolver for mocking requests
type resolver struct {
	hosts        map[string][]string
	ipAddresses  map[string][]net.IPAddr
	liveResolver ResolverInterface
	srvRecords   map[string][]*net.SRV
}

// newCustomResolver will return a custom resolver with specific records hard coded ,
func newCustomResolver(liveResolver ResolverInterface, hosts map[string][]string,
	srvRecords map[string][]*net.SRV, ipAddresses map[string][]net.IPAddr) ResolverInterface {
	return &resolver{
		hosts:        hosts,
		ipAddresses:  ipAddresses,
		liveResolver: liveResolver,
		srvRecords:   srvRecords,
	}
}

// LookupHost will lookup a host
func (r *resolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	records := r.hosts[host]
	if len(records) != 0 {
		return records, nil
	}
	return r.liveResolver.LookupHost(ctx, host)
}

// LookupIPAddr will lookup an ip address
func (r *resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	records := r.ipAddresses[host]
	if len(records) != 0 {
		return records, nil
	}
	return r.liveResolver.LookupIPAddr(ctx, host)
}

// LookupSRV will lookup an SRV record
func (r *resolver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	records := r.srvRecords[service+proto+name]
	if len(records) != 0 {
		if service == "invalid" { // Returns an invalid cname
			return fmt.Sprintf("_%s._%s", service, proto), records, nil
		}
		return fmt.Sprintf("_%s._%s.%s.", service, proto, name), records, nil
	}
	return r.liveResolver.LookupSRV(ctx, service, proto, name)
}
