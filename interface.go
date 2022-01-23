package paymail

import (
	"context"
	"net"

	"github.com/go-resty/resty/v2"
	"github.com/tonicpow/go-paymail/interfaces"
)

// ClientInterface is the Paymail client interface
type ClientInterface interface {
	CheckDNSSEC(domain string) (result *DNSCheckResult)
	CheckSSL(host string) (valid bool, err error)
	GetBRFCs() []*BRFCSpec
	GetCapabilities(target string, port int) (response *CapabilitiesResponse, err error)
	GetOptions() *ClientOptions
	GetP2PPaymentDestination(p2pURL, alias, domain string, paymentRequest *PaymentRequest) (response *PaymentDestinationResponse, err error)
	GetPKI(pkiURL, alias, domain string) (response *PKIResponse, err error)
	GetPublicProfile(publicProfileURL, alias, domain string) (response *PublicProfileResponse, err error)
	GetResolver() interfaces.DNSResolver
	GetSRVRecord(service, protocol, domainName string) (srv *net.SRV, err error)
	GetUserAgent() string
	ResolveAddress(resolutionURL, alias, domain string, senderRequest *SenderRequest) (response *ResolutionResponse, err error)
	SendP2PTransaction(p2pURL, alias, domain string, transaction *P2PTransaction) (response *P2PTransactionResponse, err error)
	ValidateSRVRecord(ctx context.Context, srv *net.SRV, port, priority, weight uint16) error
	VerifyPubKey(verifyURL, alias, domain, pubKey string) (response *VerificationResponse, err error)
	WithCustomHTTPClient(client *resty.Client) ClientInterface
	WithCustomResolver(resolver interfaces.DNSResolver) ClientInterface
}
