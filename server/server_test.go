package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateServer will test the method CreateServer()
func TestCreateServer(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := &Configuration{
			Port:    12345,
			Timeout: 10 * time.Second,
		}
		s := CreateServer(config)
		require.NotNil(t, s)
		assert.IsType(t, &http.Server{}, s)
		assert.Equal(t, fmt.Sprintf(":%d", config.Port), s.Addr)
		assert.Equal(t, config.Timeout, s.WriteTimeout)
		assert.Equal(t, config.Timeout, s.ReadTimeout)
	})
}

// TestStart will test the method Start()
func TestStart(t *testing.T) {
	t.Run("run server", func(t *testing.T) {
		/*
			// todo: run in a non-blocking way to test
				config := &Configuration{
					Port:    12345,
					Timeout: 10 * time.Second,
				}
				s := CreateServer(config)
				StartServer(s)
		*/
	})
}

// Test_removePort will test the method removePort()
func Test_removePort(t *testing.T) {
	testDomain := "domain.com"

	t.Run("valid removal", func(t *testing.T) {
		host := testDomain + ":1234"
		rp := removePort(host)
		assert.Equal(t, rp, testDomain)
	})

	t.Run("valid removal (no port)", func(t *testing.T) {
		host := testDomain + ":"
		rp := removePort(host)
		assert.Equal(t, rp, testDomain)
	})

	t.Run("no port", func(t *testing.T) {
		rp := removePort(testDomain)
		assert.Equal(t, rp, testDomain)
	})
}

// Test_getHost will test the method getHost()
func Test_getHost(t *testing.T) {
	testDomain := "domain.com"

	t.Run("valid host with port", func(t *testing.T) {
		req, err := http.NewRequestWithContext(
			context.Background(), http.MethodGet,
			"http://"+testDomain+":1234", nil,
		)
		require.NoError(t, err)
		require.NotNil(t, req)

		host := getHost(req)
		assert.Equal(t, testDomain, host)
	})

	t.Run("valid host with no port", func(t *testing.T) {
		req, err := http.NewRequestWithContext(
			context.Background(), http.MethodGet,
			"http://"+testDomain+"/", nil,
		)
		require.NoError(t, err)
		require.NotNil(t, req)

		host := getHost(req)
		assert.Equal(t, testDomain, host)
	})
}
