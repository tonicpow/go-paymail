package server

import (
	"time"

	"github.com/tonicpow/go-paymail"
)

// ConfigOps allow functional options to be supplied
// that overwrite default options.
type ConfigOps func(c *Configuration)

// defaultConfigOptions will return a Configuration struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultConfigOptions() *Configuration {
	return &Configuration{
		APIVersion:              DefaultAPIVersion,
		BasicRoutes:             &basicRoutes{},
		BSVAliasVersion:         paymail.DefaultBsvAliasVersion,
		Capabilities:            genericCapabilities(paymail.DefaultBsvAliasVersion, DefaultSenderValidation),
		Port:                    DefaultServerPort,
		Prefix:                  DefaultPrefix,
		SenderValidationEnabled: DefaultSenderValidation,
		ServiceName:             paymail.DefaultServiceName,
		Timeout:                 DefaultTimeout,
	}
}

// WithGenericCapabilities will load the generic Paymail capabilities
func WithGenericCapabilities() ConfigOps {
	return func(c *Configuration) {
		c.Capabilities = genericCapabilities(c.BSVAliasVersion, c.SenderValidationEnabled)
	}
}

// WithCapabilities will modify the capabilities
func WithCapabilities(capabilities *Capabilities) ConfigOps {
	return func(c *Configuration) {
		if capabilities != nil {
			// todo: validate that these are valid capabilities (string->url path)
			c.Capabilities = capabilities
		}
	}
}

// WithBasicRoutes will turn on all the basic routes
func WithBasicRoutes() ConfigOps {
	return func(c *Configuration) {
		c.BasicRoutes = &basicRoutes{
			Add404Route:    true,
			AddHealthRoute: true,
			AddIndexRoute:  true,
			AddNotAllowed:  true,
		}
	}
}

// WithTimeout will set a custom timeout
func WithTimeout(timeout time.Duration) ConfigOps {
	return func(c *Configuration) {
		if timeout > 0 {
			c.Timeout = timeout
		}
	}
}

// WithServiceName will set a custom service name
func WithServiceName(serviceName string) ConfigOps {
	return func(c *Configuration) {
		if len(serviceName) > 0 {
			c.ServiceName = serviceName
		}
	}
}

// WithSenderValidation will enable sender validation
func WithSenderValidation() ConfigOps {
	return func(c *Configuration) {
		c.SenderValidationEnabled = true
	}
}

// WithDomain will add the domain if not found
func WithDomain(domain string) ConfigOps {
	return func(c *Configuration) {
		if len(domain) > 0 {
			// todo: attempt to add, but cannot return the error
			_ = c.AddDomain(domain)
		}
	}
}

// WithPort will overwrite the default port
func WithPort(port int) ConfigOps {
	return func(c *Configuration) {
		if port > 0 {
			c.Port = port
		}
	}
}
