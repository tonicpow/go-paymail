package interfaces

import (
	"context"
	"net"
)

// DNSResolver is a custom resolver interface for testing
type DNSResolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
}
