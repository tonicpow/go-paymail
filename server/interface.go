package server

// PaymailAddress is an internal struct for paymail addresses
type PaymailAddress struct {
	Alias       string `json:"alias"`        // Alias or handle of the paymail
	Avatar      string `json:"avatar"`       // This is the url of the user (public profile)
	ID          uint64 `json:"id"`           // Unique identifier
	LastAddress string `json:"last_address"` // This is used as a temp address for now (should be via xPub)
	Name        string `json:"name"`         // This is the name of the user (public profile)
	PrivateKey  string `json:"private_key"`  // PrivateKey hex encoded
	PubKey      string `json:"pubkey"`       // PublicKey hex encoded
}

// PaymailServerInterface the paymail server interface that needs to be implemented
type PaymailServerInterface interface {
	GetPaymailByAlias(alias string) *PaymailAddress
}
