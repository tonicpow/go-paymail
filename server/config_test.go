package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-paymail"
)

// Mock implementation of a service provider
type mockServiceProvider struct {
	// Extend your dependencies or custom values
}

// GetPaymailByAlias is a demo implementation of this interface
func (m *mockServiceProvider) GetPaymailByAlias(_ context.Context, alias, domain string,
	_ *RequestMetadata) (*paymail.AddressInformation, error) {

	// Get the data from the demo database
	return nil, nil
}

// CreateAddressResolutionResponse is a demo implementation of this interface
func (m *mockServiceProvider) CreateAddressResolutionResponse(ctx context.Context, alias, domain string,
	senderValidation bool, _ *RequestMetadata) (*paymail.ResolutionInformation, error) {

	// Generate a new destination / output for the basic address resolution
	return nil, nil
}

// CreateP2PDestinationResponse is a demo implementation of this interface
func (m *mockServiceProvider) CreateP2PDestinationResponse(ctx context.Context, alias, domain string,
	satoshis uint64, _ *RequestMetadata) (*paymail.PaymentDestinationInformation, error) {

	// Generate a new destination for the p2p request
	return nil, nil
}

// RecordTransaction is a demo implementation of this interface
func (m *mockServiceProvider) RecordTransaction(ctx context.Context,
	p2pTx *paymail.P2PTransaction, _ *RequestMetadata) (*paymail.P2PTransactionInformation, error) {

	// Record the tx into your datastore layer
	return nil, nil
}

// TestConfiguration_Validate will test the method Validate()
func TestConfiguration_Validate(t *testing.T) {
	// todo: finish test!
}

// TestConfiguration_IsAllowedDomain will test the method IsAllowedDomain()
func TestConfiguration_IsAllowedDomain(t *testing.T) {
	// todo: finish test!
}

// TestConfiguration_AddDomain will test the method AddDomain()
func TestConfiguration_AddDomain(t *testing.T) {
	t.Parallel()

	t.Run("no domain", func(t *testing.T) {
		testDomain := "test.com"
		c, err := NewConfig(new(mockServiceProvider), WithDomain(testDomain), WithGenericCapabilities())
		require.NoError(t, err)
		require.NotNil(t, c)

		err = c.AddDomain("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDomainMissing)
	})

	t.Run("sanitized domain", func(t *testing.T) {
		testDomain := "WWW.TEST.COM"
		addDomain := "testER.com"
		c, err := NewConfig(new(mockServiceProvider), WithDomain(testDomain), WithGenericCapabilities())
		require.NoError(t, err)
		require.NotNil(t, c)

		err = c.AddDomain(addDomain)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(c.PaymailDomains))
		assert.Equal(t, "test.com", c.PaymailDomains[0].Name)
		assert.Equal(t, "tester.com", c.PaymailDomains[1].Name)
	})
}

// TestConfiguration_EnrichCapabilities will test the method EnrichCapabilities()
func TestConfiguration_EnrichCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("basic enrich", func(t *testing.T) {
		testDomain := "test.com"
		c, err := NewConfig(new(mockServiceProvider), WithDomain(testDomain), WithGenericCapabilities())
		require.NoError(t, err)
		require.NotNil(t, c)

		capabilities := c.EnrichCapabilities(testDomain)
		assert.Equal(t, 5, len(capabilities.Capabilities))
		assert.Equal(t, paymail.DefaultBsvAliasVersion, c.Capabilities.BsvAlias)
		assert.Equal(t, "https://"+testDomain+"/v1/bsvalias/address/{alias}@{domain.tld}", capabilities.Capabilities[paymail.BRFCPaymentDestination])
		assert.Equal(t, "https://"+testDomain+"/v1/bsvalias/id/{alias}@{domain.tld}", capabilities.Capabilities[paymail.BRFCPki])
		assert.Equal(t, "https://"+testDomain+"/v1/bsvalias/public-profile/{alias}@{domain.tld}", capabilities.Capabilities[paymail.BRFCPublicProfile])
		assert.Equal(t, "https://"+testDomain+"/v1/bsvalias/verify-pubkey/{alias}@{domain.tld}/{pubkey}", capabilities.Capabilities[paymail.BRFCVerifyPublicKeyOwner])
		assert.Equal(t, false, capabilities.Capabilities[paymail.BRFCSenderValidation])
	})

	t.Run("multiple times", func(t *testing.T) {
		testDomain := "test.com"
		c, err := NewConfig(new(mockServiceProvider), WithDomain("test.com"), WithGenericCapabilities())
		require.NoError(t, err)
		require.NotNil(t, c)

		capabilities := c.EnrichCapabilities(testDomain)
		assert.Equal(t, 5, len(capabilities.Capabilities))

		capabilities = c.EnrichCapabilities(testDomain)
		assert.Equal(t, 5, len(capabilities.Capabilities))
	})
}

// TestGenerateServiceURL will test the method GenerateServiceURL()
func TestGenerateServiceURL(t *testing.T) {
	t.Parallel()

	t.Run("valid values", func(t *testing.T) {
		u := GenerateServiceURL("https://", "test.com", "v1", "bsvalias")
		assert.Equal(t, "https://test.com/v1/bsvalias", u)
	})

	t.Run("all invalid values", func(t *testing.T) {
		u := GenerateServiceURL("", "", "", "")
		assert.Equal(t, "", u)
	})

	t.Run("missing prefix", func(t *testing.T) {
		u := GenerateServiceURL("", "test.com", "v1", "")
		assert.Equal(t, "", u)
	})

	t.Run("missing domain", func(t *testing.T) {
		u := GenerateServiceURL("https://", "", "v1", "")
		assert.Equal(t, "", u)
	})

	t.Run("no api version", func(t *testing.T) {
		u := GenerateServiceURL("https://", "test", "", "bsvalias")
		assert.Equal(t, "https://test/bsvalias", u)
	})

	t.Run("no service name", func(t *testing.T) {
		u := GenerateServiceURL("https://", "test", "v1", "")
		assert.Equal(t, "https://test/v1", u)
	})
}

// TestNewConfig will test the method NewConfig()
func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("no values and no provider", func(t *testing.T) {
		c, err := NewConfig(nil)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrServiceProviderNil)
		assert.Nil(t, c)
	})

	t.Run("missing domain", func(t *testing.T) {
		c, err := NewConfig(new(mockServiceProvider))
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDomainMissing)
		assert.Nil(t, c)
	})

	t.Run("valid client - minimum options", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, 5, len(c.Capabilities.Capabilities))
		assert.Equal(t, "test.com", c.PaymailDomains[0].Name)
	})

	t.Run("custom port", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithPort(12345),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, 12345, c.Port)
	})

	t.Run("custom timeout", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithTimeout(10*time.Second),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, 10*time.Second, c.Timeout)
	})

	t.Run("custom service name", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithServiceName("custom"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, "custom", c.ServiceName)
	})

	t.Run("sender validation enabled", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithSenderValidation(),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, true, c.SenderValidationEnabled)
	})

	t.Run("with basic routes", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithBasicRoutes(),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.NotNil(t, c.BasicRoutes)
		assert.Equal(t, true, c.BasicRoutes.Add404Route)
		assert.Equal(t, true, c.BasicRoutes.AddIndexRoute)
		assert.Equal(t, true, c.BasicRoutes.AddHealthRoute)
		assert.Equal(t, true, c.BasicRoutes.AddNotAllowed)
	})
}
