package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-paymail"
)

// testConfig loads a basic test configuration
func testConfig(t *testing.T, domain string) *Configuration {
	c, err := NewConfig(
		new(mockServiceProvider),
		WithDomain(domain),
		WithGenericCapabilities(),
	)
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

// TestConfiguration_Validate will test the method Validate()
func TestConfiguration_Validate(t *testing.T) {
	t.Parallel()

	t.Run("missing domain", func(t *testing.T) {
		c := &Configuration{}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDomainMissing)
	})

	t.Run("missing port", func(t *testing.T) {
		c := &Configuration{
			PaymailDomains: []*Domain{{Name: "test.com"}},
		}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPortMissing)
	})

	t.Run("missing service name", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			PaymailDomains: []*Domain{{Name: "test.com"}},
		}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrServiceNameMissing)
	})

	t.Run("invalid service name", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			ServiceName:    "$*%*",
			PaymailDomains: []*Domain{{Name: "test.com"}},
		}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrServiceNameMissing)
	})

	t.Run("missing capabilities", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			ServiceName:    "test",
			PaymailDomains: []*Domain{{Name: "test.com"}},
		}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrCapabilitiesMissing)
	})

	t.Run("invalid capabilities", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			ServiceName:    "test",
			PaymailDomains: []*Domain{{Name: "test.com"}},
			Capabilities: &paymail.CapabilitiesPayload{
				BsvAlias: "",
			},
		}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrBsvAliasMissing)
	})

	t.Run("zero capabilities", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			ServiceName:    "test",
			PaymailDomains: []*Domain{{Name: "test.com"}},
			Capabilities: &paymail.CapabilitiesPayload{
				BsvAlias: "test",
			},
		}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrCapabilitiesMissing)
	})

	t.Run("basic valid configuration", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			ServiceName:    "test",
			PaymailDomains: []*Domain{{Name: "test.com"}},
			Capabilities:   GenericCapabilities("test", false),
		}
		err := c.Validate()
		require.NoError(t, err)
	})

	t.Run("configuration with domain validation disabled", func(t *testing.T) {
		c := &Configuration{
			Port:           12345,
			ServiceName:    "test",
			PaymailDomains: []*Domain{},
			Capabilities:   GenericCapabilities("test", false),
		}
		assert.False(t, c.PaymailDomainsValidationDisabled)
		err := c.Validate()
		assert.ErrorIs(t, err, ErrDomainMissing)

		c.PaymailDomainsValidationDisabled = true
		err = c.Validate()
		assert.NoError(t, err)
	})
}

// TestConfiguration_IsAllowedDomain will test the method IsAllowedDomain()
func TestConfiguration_IsAllowedDomain(t *testing.T) {
	t.Parallel()

	t.Run("empty domain", func(t *testing.T) {
		c := testConfig(t, "test.com")
		require.NotNil(t, c)

		success := c.IsAllowedDomain("")
		assert.Equal(t, false, success)
	})

	t.Run("domain found", func(t *testing.T) {
		c := testConfig(t, "test.com")
		require.NotNil(t, c)

		success := c.IsAllowedDomain("test.com")
		assert.Equal(t, true, success)
	})

	t.Run("sanitized domain found", func(t *testing.T) {
		c := testConfig(t, "test.com")
		require.NotNil(t, c)

		success := c.IsAllowedDomain("WWW.test.COM")
		assert.Equal(t, true, success)
	})

	t.Run("both domains are sanitized", func(t *testing.T) {
		c := testConfig(t, "WwW.Test.Com")
		require.NotNil(t, c)

		success := c.IsAllowedDomain("WWW.test.COM")
		assert.Equal(t, true, success)
	})

	t.Run("domain validation on", func(t *testing.T) {
		c := testConfig(t, "WwW.Test.Com")
		c.PaymailDomainsValidationDisabled = false
		require.NotNil(t, c)

		assert.Equal(t, true, c.IsAllowedDomain("test.com"))
		assert.Equal(t, false, c.IsAllowedDomain("test2.com"))
	})

	t.Run("domain validation off", func(t *testing.T) {
		c := testConfig(t, "WwW.Test.Com")
		c.PaymailDomainsValidationDisabled = true
		require.NotNil(t, c)

		assert.Equal(t, true, c.IsAllowedDomain("test.com"))
		assert.Equal(t, true, c.IsAllowedDomain("test2.com"))
	})
}

// TestConfiguration_AddDomain will test the method AddDomain()
func TestConfiguration_AddDomain(t *testing.T) {
	t.Parallel()

	t.Run("no domain", func(t *testing.T) {
		testDomain := "test.com"
		c := testConfig(t, testDomain)
		require.NotNil(t, c)

		err := c.AddDomain("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDomainMissing)
	})

	t.Run("sanitized domain", func(t *testing.T) {
		testDomain := "WWW.TEST.COM"
		addDomain := "testER.com"
		c := testConfig(t, testDomain)
		require.NotNil(t, c)

		err := c.AddDomain(addDomain)
		require.NoError(t, err)

		assert.Equal(t, 2, len(c.PaymailDomains))
		assert.Equal(t, "test.com", c.PaymailDomains[0].Name)
		assert.Equal(t, "tester.com", c.PaymailDomains[1].Name)
	})

	t.Run("domain already exists", func(t *testing.T) {
		testDomain := "test.com"
		addDomain := "test.com"
		c := testConfig(t, testDomain)
		require.NotNil(t, c)

		err := c.AddDomain(addDomain)
		require.NoError(t, err)

		assert.Equal(t, 1, len(c.PaymailDomains))
		assert.Equal(t, "test.com", c.PaymailDomains[0].Name)
	})
}

// TestConfiguration_EnrichCapabilities will test the method EnrichCapabilities()
func TestConfiguration_EnrichCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("basic enrich", func(t *testing.T) {
		testDomain := "test.com"
		c := testConfig(t, testDomain)
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
		c := testConfig(t, testDomain)
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
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrServiceProviderNil)
		assert.Nil(t, c)
	})

	t.Run("missing domain", func(t *testing.T) {
		c, err := NewConfig(new(mockServiceProvider))
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDomainMissing)
		assert.Nil(t, c)
	})

	t.Run("valid client - minimum options", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, 5, len(c.Capabilities.Capabilities))
		assert.Equal(t, "test.com", c.PaymailDomains[0].Name)
	})

	t.Run("custom port", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithPort(12345),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, 12345, c.Port)
	})

	t.Run("custom timeout", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithTimeout(10*time.Second),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, 10*time.Second, c.Timeout)
	})

	t.Run("custom service name", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithServiceName("custom"),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, "custom", c.ServiceName)
	})

	t.Run("sender validation enabled", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithSenderValidation(),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, true, c.SenderValidationEnabled)
	})

	t.Run("with custom capabilities", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithCapabilities(GenericCapabilities("test", false)),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, 5, len(c.Capabilities.Capabilities))
		assert.Equal(t, "test", c.Capabilities.BsvAlias)
	})

	t.Run("with basic routes", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithBasicRoutes(),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.BasicRoutes)
		assert.Equal(t, true, c.BasicRoutes.Add404Route)
		assert.Equal(t, true, c.BasicRoutes.AddIndexRoute)
		assert.Equal(t, true, c.BasicRoutes.AddHealthRoute)
		assert.Equal(t, true, c.BasicRoutes.AddNotAllowed)
	})

	t.Run("domain validation disabled", func(t *testing.T) {
		c, err := NewConfig(
			new(mockServiceProvider),
			WithDomain("test.com"),
			WithPort(12345),
			WithDomainValidationDisabled(),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, 12345, c.Port)
		assert.Equal(t, true, c.PaymailDomainsValidationDisabled)
	})
}
