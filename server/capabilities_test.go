package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-paymail"
)

// TestGenericCapabilities will test the method GenericCapabilities()
func TestGenericCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("valid values", func(t *testing.T) {
		c := GenericCapabilities("test", true)
		require.NotNil(t, c)
		assert.Equal(t, "test", c.BsvAlias)
		assert.Equal(t, 5, len(c.Capabilities))
	})

	t.Run("no alias version", func(t *testing.T) {
		c := GenericCapabilities("", true)
		require.NotNil(t, c)
		assert.Equal(t, "", c.BsvAlias)
	})

	t.Run("sender validation", func(t *testing.T) {
		c := GenericCapabilities("", true)
		require.NotNil(t, c)
		assert.Equal(t, true, c.Capabilities[paymail.BRFCSenderValidation])
	})
}

// TestP2PCapabilities will test the method P2PCapabilities()
func TestP2PCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("valid values", func(t *testing.T) {
		c := P2PCapabilities("test", true)
		require.NotNil(t, c)
		assert.Equal(t, "test", c.BsvAlias)
		assert.Equal(t, 7, len(c.Capabilities))
	})

	t.Run("no alias version", func(t *testing.T) {
		c := P2PCapabilities("", true)
		require.NotNil(t, c)
		assert.Equal(t, "", c.BsvAlias)
	})

	t.Run("sender validation", func(t *testing.T) {
		c := P2PCapabilities("", true)
		require.NotNil(t, c)
		assert.Equal(t, true, c.Capabilities[paymail.BRFCSenderValidation])
	})

	t.Run("has p2p routes", func(t *testing.T) {
		c := P2PCapabilities("", true)
		require.NotNil(t, c)
		assert.NotEmpty(t, c.Capabilities[paymail.BRFCP2PTransactions])
		assert.NotEmpty(t, c.Capabilities[paymail.BRFCP2PPaymentDestination])
	})
}
