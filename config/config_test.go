package config_test

import (
	"os"
	"testing"
	"time"

	. "github.com/bbengfort/catena/config"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaults(t *testing.T) {
	c, err := New()
	require.NoError(t, err, "could not configure config defaults")

	// Top Level Defaults
	require.Equal(t, "localhost", c.Domain)
	require.Equal(t, "127.0.0.1", c.Addr)
	require.Equal(t, uint16(8888), c.Port)
	require.False(t, c.NoTLS)
	require.Empty(t, c.DBURL)

	// Routes Defaults
	require.True(t, c.Routes.RedirectTrailingSlash)
	require.True(t, c.Routes.RedirectFixedPath)
	require.True(t, c.Routes.HandleMethodNotAllowed)

	// Timeouts Defaults
	require.Equal(t, 10*time.Second, c.ReadTimeout)
	require.Equal(t, 20*time.Second, c.WriteTimeout)
	require.Equal(t, 5*time.Minute, c.IdleTimeout)
}

func TestConfigEnviron(t *testing.T) {
	// Set the environment
	envvars := map[string]string{
		"CATENA_DOMAIN":        "catena.dev",
		"CATENA_BIND_ADDR":     "0.0.0.0",
		"CATENA_PORT":          "443",
		"CATENA_NO_TLS":        "true",
		"DATABASE_URL":         "postgres://user@localhost:5432/db",
		"CATENA_READ_TIMEOUT":  "1m",
		"CATENA_WRITE_TIMEOUT": "500ms",
		"CATENA_IDLE_TIMEOUT":  "3h",
	}

	for key, val := range envvars {
		require.NoError(t, os.Setenv(key, val))
	}

	// Ensure environemnt is unset at end of the test
	defer func() {
		for key := range envvars {
			require.NoError(t, os.Unsetenv(key))
		}
	}()

	c, err := New()
	require.NoError(t, err, "could not configure config defaults")

	// Top Level Defaults
	require.Equal(t, "catena.dev", c.Domain)
	require.Equal(t, "0.0.0.0", c.Addr)
	require.Equal(t, uint16(443), c.Port)
	require.True(t, c.NoTLS)
	require.Equal(t, "postgres://user@localhost:5432/db", c.DBURL)

	// Routes Defaults
	require.Equal(t, true, c.Routes.RedirectTrailingSlash)
	require.Equal(t, true, c.Routes.RedirectFixedPath)
	require.True(t, c.Routes.HandleMethodNotAllowed)

	// Timeouts Defaults
	require.Equal(t, 1*time.Minute, c.ReadTimeout)
	require.Equal(t, 500*time.Millisecond, c.WriteTimeout)
	require.Equal(t, 180*time.Minute, c.IdleTimeout)
}
