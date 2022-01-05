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
