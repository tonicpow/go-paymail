package paymail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNetwork_String will test the method String()
func TestNetwork_String(t *testing.T) {
	t.Run("valid networks", func(t *testing.T) {
		assert.Equal(t, "mainnet", Mainnet.String())
		assert.Equal(t, "testnet", Testnet.String())
		assert.Equal(t, "STN", STN.String())
	})
	t.Run("invalid network", func(t *testing.T) {
		b := Network(8)
		assert.Equal(t, "not recognized", b.String())
	})
}

// TestNetwork_PaymailURLSuffix will test the method URLSuffix()
func TestNetwork_PaymailURLSuffix(t *testing.T) {
	t.Run("valid networks", func(t *testing.T) {
		assert.Equal(t, "", Mainnet.URLSuffix())
		assert.Equal(t, "-testnet", Testnet.URLSuffix())
		assert.Equal(t, "-stn", STN.URLSuffix())
	})
	t.Run("invalid network", func(t *testing.T) {
		b := new(Network)
		assert.Equal(t, "", b.URLSuffix())
	})
}
