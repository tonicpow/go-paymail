package paymail

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

/*
Alternative checks:

https://dnsviz.net/d/domain.com/dnssec/
https://dnssec-analyzer.verisignlabs.com/domain.com
*/

// DNSCheckResult struct is returned for the DNS check
type DNSCheckResult struct {
	Answer       answer    `json:"answer"`
	CheckTime    time.Time `json:"check_time"`
	DNSSEC       bool      `json:"dnssec"`
	Domain       string    `json:"domain,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	NSEC         nsec      `json:"nsec"`
}

// nsec struct for NSEC type
type nsec struct {
	NSEC       *dns.NSEC       `json:"nsec,omitempty"`
	NSEC3      *dns.NSEC3      `json:"nsec_3,omitempty"`
	NSEC3PARAM *dns.NSEC3PARAM `json:"nsec_3_param,omitempty"`
	Type       string          `json:"type,omitempty"`
}

// answer struct the answer of the DNS question
type answer struct {
	CalculatedDS      []*domainDS     `json:"calculate_ds,omitempty"`
	DNSKEYRecordCount int             `json:"dnskey_record_count,omitempty"`
	DNSKEYRecords     []*domainDNSKEY `json:"dnskey_records,omitempty"`
	DSRecordCount     int             `json:"ds_record_count,omitempty"`
	DSRecords         []*domainDS     `json:"ds_records,omitempty"`
	Matching          matching        `json:"matching,omitempty"`
}

// matching struct for information
type matching struct {
	DNSKEY []*domainDNSKEY `json:"dnskey,omitempty"`
	DS     []*domainDS     `json:"ds,omitempty"`
}

// domainDS struct
type domainDS struct {
	Algorithm  uint8  `json:"algorithm,omitempty"`
	Digest     string `json:"digest,omitempty"`
	DigestType uint8  `json:"digest_type,omitempty"`
	KeyTag     uint16 `json:"key_tag,omitempty"`
}

// domainDNSKEY struct
type domainDNSKEY struct {
	Algorithm    uint8     `json:"algorithm,omitempty"`
	CalculatedDS *domainDS `json:"calculate_ds,omitempty"`
	Flags        uint16    `json:"flags,omitempty"`
	Protocol     uint8     `json:"protocol,omitempty"`
	PublicKey    string    `json:"public_key,omitempty"`
}

// Domains that DO NOT work properly for DNSSEC validation
var (

	// todo: find a way to make these work
	// https://network-tools.com/nslookup/ for a heroku app produces 0 results
	domainsWithIssues = []string{
		"herokuapp.com", // CNAME on heroku is a pointer, and thus there is no NS returned
	}
)

// CheckDNSSEC will check the DNSSEC for a given domain
//
// Paymail providers should have DNSSEC enabled for their domain
func (c *Client) CheckDNSSEC(domain string) (result *DNSCheckResult) {

	// Start the new result
	result = new(DNSCheckResult)
	result.CheckTime = time.Now()

	var err error

	// Valid domain name (ASCII or IDN)
	if domain, err = idna.ToASCII(domain); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in ToASCII: %s", err.Error())
		return
	}

	// Validate domain
	if domain, err = publicsuffix.EffectiveTLDPlusOne(domain); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in EffectiveTLDPlusOne: %s", err.Error())
		return
	}

	// Set the valid domain now
	result.Domain = domain

	// Check known domain issues
	for _, d := range domainsWithIssues {
		if strings.Contains(result.Domain, d) {
			result.ErrorMessage = fmt.Sprintf("%s cannot be validated due to a known issue with %s", result.Domain, d)
			return
		}
	}

	// Set the TLD
	tld, _ := publicsuffix.PublicSuffix(domain)

	// Set the registry name server
	var registryNameserver string
	if registryNameserver, err = resolveOneNS(tld, c.Options.NameServer, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveOneNS: %s", err.Error())
		return
	}

	// Set the domain name server
	var domainNameserver string
	if domainNameserver, err = resolveOneNS(domain, c.Options.NameServer, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveOneNS: %s", err.Error())
		return
	}

	// Domain name servers at registrar Host
	var domainDsRecord []*domainDS
	if domainDsRecord, err = resolveDomainDS(domain, registryNameserver, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveDomainDS: %s", err.Error())
		return
	}

	// Set the records and count
	result.Answer.DSRecords = domainDsRecord
	result.Answer.DSRecordCount = cap(domainDsRecord)

	// Resolve domain DNSKey
	var dnsKey []*domainDNSKEY
	if dnsKey, err = resolveDomainDNSKEY(domain, domainNameserver, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveDomainDNSKEY: %s", err.Error())
		return
	}

	// Set the DNSKEY records
	result.Answer.DNSKEYRecords = dnsKey
	result.Answer.DNSKEYRecordCount = cap(result.Answer.DNSKEYRecords)

	// Check the digest type
	var digest uint8
	if cap(result.Answer.DSRecords) != 0 {
		digest = result.Answer.DSRecords[0].DigestType
	}

	// Check the DS record
	if result.Answer.DSRecordCount > 0 && result.Answer.DNSKEYRecordCount > 0 {
		var calculatedDS []*domainDS
		if calculatedDS, err = calculateDSRecord(domain, domainNameserver, c.Options.DNSPort, digest); err != nil {
			result.ErrorMessage = fmt.Sprintf("failed in calculateDSRecord: %s", err.Error())
			return
		}
		result.Answer.CalculatedDS = calculatedDS
	}

	// Resolve the domain NSEC
	var nsec *dns.NSEC
	if nsec, err = resolveDomainNSEC(domain, c.Options.NameServer, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveDomainNSEC: %s", err.Error())
		return
	} else if nsec != nil {
		result.NSEC.Type = "nsec"
		result.NSEC.NSEC = nsec
	}

	// Resolve the domain NSEC3
	var nsec3 *dns.NSEC3
	if nsec3, err = resolveDomainNSEC3(domain, c.Options.NameServer, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveDomainNSEC3: %s", err.Error())
		return
	} else if nsec3 != nil {
		result.NSEC.Type = "nsec3"
		result.NSEC.NSEC3 = nsec3
	}

	// Resolve the domain NSEC3PARAM
	var nsec3param *dns.NSEC3PARAM
	if nsec3param, err = resolveDomainNSEC3PARAM(domain, c.Options.NameServer, c.Options.DNSPort); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed in resolveDomainNSEC3PARAM: %s", err.Error())
		return
	} else if nsec3param != nil {
		result.NSEC.Type = "nsec3param"
		result.NSEC.NSEC3PARAM = nsec3param
	}

	// Check the keys and set the DNSSEC flag
	if result.Answer.DSRecordCount > 0 && result.Answer.DNSKEYRecordCount > 0 {
		var filtered []*domainDS
		var dnsKeys []*domainDNSKEY
		for _, e := range result.Answer.DSRecords {
			for i, f := range result.Answer.CalculatedDS {
				if f.Digest == e.Digest {
					filtered = append(filtered, f)
					dnsKeys = append(dnsKeys, result.Answer.DNSKEYRecords[i])
				}
			}
		}
		result.Answer.Matching.DS = filtered
		result.Answer.Matching.DNSKEY = dnsKeys
		result.DNSSEC = true
	} else {
		result.DNSSEC = false
	}

	// Complete
	return
}

/*
Source: https://github.com/binaryfigments/dnscheck
License: https://github.com/binaryfigments/dnscheck/blob/master/LICENSE
*/

// resolveOneNS will resolve one name server
func resolveOneNS(domain, nameServer, dnsPort string) (string, error) {
	var answer []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	m.MsgHdr.RecursionDesired = true
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return "", err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.NS); ok {
			answer = append(answer, a.Ns)
		}
	}
	if len(answer) < 1 || answer == nil {
		return "", err
	}
	return answer[0], nil
}

// resolveDomainNSEC will resolve a domain NSEC
func resolveDomainNSEC(domain, nameServer, dnsPort string) (*dns.NSEC, error) {
	var answer *dns.NSEC
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNSEC)
	m.MsgHdr.RecursionDesired = true
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return nil, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.NSEC); ok {
			answer = a
			return answer, nil
		}
	}
	return nil, nil
}

// resolveDomainNSEC3 will resolve a domain NSEC3
func resolveDomainNSEC3(domain, nameServer, dnsPort string) (*dns.NSEC3, error) {
	var answer *dns.NSEC3
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNSEC3)
	m.MsgHdr.RecursionDesired = true
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return nil, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.NSEC3); ok {
			answer = a
			return answer, nil
		}
	}
	return nil, nil
}

// resolveDomainNSEC3PARAM will resolve a domain NSEC3PARAM
func resolveDomainNSEC3PARAM(domain, nameServer, dnsPort string) (*dns.NSEC3PARAM, error) {
	var answer *dns.NSEC3PARAM
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNSEC3PARAM)
	m.MsgHdr.RecursionDesired = true
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return nil, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.NSEC3PARAM); ok {
			answer = a
			return answer, nil
		}
	}
	return nil, nil
}

// resolveDomainDS will resolve a domain DS
func resolveDomainDS(domain, nameServer, dnsPort string) ([]*domainDS, error) {
	var ds []*domainDS
	m := new(dns.Msg)
	m.MsgHdr.RecursionDesired = true
	m.SetQuestion(dns.Fqdn(domain), dns.TypeDS)
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return ds, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.DS); ok {
			readKey := &domainDS{
				Algorithm:  a.Algorithm,
				Digest:     a.Digest,
				DigestType: a.DigestType,
				KeyTag:     a.KeyTag,
			}
			ds = append(ds, readKey)
		}
	}
	return ds, nil
}

// resolveDomainDNSKEY will resolve a domain DNSKEY
func resolveDomainDNSKEY(domain, nameServer, dnsPort string) ([]*domainDNSKEY, error) {
	var dnskey []*domainDNSKEY

	m := new(dns.Msg)
	m.MsgHdr.RecursionDesired = true
	m.SetQuestion(dns.Fqdn(domain), dns.TypeDNSKEY)
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return dnskey, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.DNSKEY); ok {
			readKey := &domainDNSKEY{
				Algorithm: a.Algorithm,
				Flags:     a.Flags,
				Protocol:  a.Protocol,
				PublicKey: a.PublicKey,
			}
			dnskey = append(dnskey, readKey)
		}
	}
	return dnskey, err
}

// calculateDSRecord function for generating DS records from the DNSKEY
// Input: domain, digest and name server from the host
// Output: one of more structs with DS information
func calculateDSRecord(domain, nameServer, dnsPort string, digest uint8) ([]*domainDS, error) {
	var calculatedDS []*domainDS

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeDNSKEY)
	m.SetEdns0(4096, true)
	m.MsgHdr.RecursionDesired = true
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameServer+":"+dnsPort)
	if err != nil {
		return calculatedDS, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.DNSKEY); ok {
			calculatedKey := &domainDS{
				Algorithm:  a.ToDS(digest).Algorithm,
				Digest:     a.ToDS(digest).Digest,
				DigestType: a.ToDS(digest).DigestType,
				KeyTag:     a.ToDS(digest).KeyTag,
			}
			calculatedDS = append(calculatedDS, calculatedKey)
		}
	}
	return calculatedDS, nil
}
