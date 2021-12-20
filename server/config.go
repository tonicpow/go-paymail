package server

import "github.com/tonicpow/go-paymail"

// Basic Configuration for the server
const (
	paymailAPIVersion = "v1" // Version of API
)

type basicRoutes struct {
	Add404Route    bool `json:"add_404_route,omitempty"`
	AddHealthRoute bool `json:"add_health_route,omitempty"`
	AddIndexRoute  bool `json:"add_index_route,omitempty"`
	AddNotAllowed  bool `json:"add_not_allowed,omitempty"`
}

// Configuration paymail server configuration object
type Configuration struct {
	BasicRoutes             *basicRoutes           `json:"basic_routes,omitempty"`
	Capabilities            *Capabilities          `json:"capabilities,omitempty"`
	actions                 PaymailServerInterface `json:"interface,omitempty"`
	PaymailDomain           string                 `json:"paymail_domain,omitempty"`
	Port                    int                    `json:"port,omitempty"`
	SenderValidationEnabled bool                   `json:"sender_validation_enabled,omitempty"`
	ServiceName             string                 `json:"service_name,omitempty"`
	ServiceURL              string                 `json:"service_url,omitempty"`
	Timeout                 int                    `json:"timeout,omitempty"`
}

// NewConfiguration create a new Configuration for the paymail server
func NewConfiguration(paymailDomain string, serverInterface PaymailServerInterface) *Configuration {
	config := &Configuration{
		BasicRoutes:             &basicRoutes{},
		actions:                 serverInterface,
		PaymailDomain:           paymailDomain,
		Port:                    3000,
		SenderValidationEnabled: false,
		ServiceName:             paymail.DefaultServiceName,
		ServiceURL:              "https://" + paymailDomain + "/" + paymailAPIVersion + "/" + paymail.DefaultServiceName + "/",
		Timeout:                 15,
	}

	config.Capabilities = createCapabilities(config)

	return config
}
