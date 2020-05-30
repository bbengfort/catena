/*
Package config implements the catena server configuration. The Config object uses tags
to handle reasonable defaults and reading configurations from the environment. It can
also be loaded from a YAML or JSON file or modified from a urifave CLI context. File
configurations can be looked up with a search path that checks standard linux
configuration locations. This ensures full flexibility in configuration for development
and production environments, while also being overkill for the job.
*/
package config

import (
	"fmt"
	"time"
)

// New creates a new configuration object with specified defaults and any values loaded
// from the environment. If the config object is improperly defined or it cannot parse
// the values from the environment it will return an error.
func New() (c Config, err error) {
	c = Config{}

	// Load the defaults from the struct tags
	if err = defaults(&c); err != nil {
		return Config{}, err
	}

	// Load values from the environment
	if err = environs(&c); err != nil {
		return Config{}, err
	}

	return c, nil
}

// Config defines the required configuration for the Catena server.
type Config struct {
	Domain string `default:"localhost" env:"CATENA_DOMAIN"`
	Addr   string `default:"127.0.0.1" env:"CATENA_BIND_ADDR"`
	Port   uint16 `default:"8888" env:"CATENA_PORT"`
	NoTLS  bool   `env:"CATENA_NO_TLS"`
	DBURL  string `env:"DATABASE_URL"`
	Routes struct {
		RedirectTrailingSlash  bool `default:"true"`
		RedirectFixedPath      bool `default:"true"`
		HandleMethodNotAllowed bool `default:"true"`
	}
	ReadTimeout  time.Duration `default:"10s" env:"CATENA_READ_TIMEOUT"`
	WriteTimeout time.Duration `default:"20s" env:"CATENA_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `default:"5m" env:"CATENA_IDLE_TIMEOUT"`
}

// Endpoint returns the human readable endpoint using either the domain or the bind addr
// with the correct protocol and port if required.
func (c Config) Endpoint() string {
	host := c.Domain
	if host == "" || (host == "localhost" && c.Addr != "127.0.0.1") {
		if c.Addr != "" {
			host = c.Addr
		}
	}

	if c.NoTLS {
		if c.Port != 80 {
			return fmt.Sprintf("http://%s:%d", host, c.Port)
		}
		return "http://" + host
	}

	if c.Port != 443 {
		return fmt.Sprintf("https://%s:%d", host, c.Port)
	}
	return "https://" + host
}

// BindAddr returns the bind address and port for the http server to listen on.
func (c Config) BindAddr() string {
	if c.Port > 0 {
		return fmt.Sprintf("%s:%d", c.Addr, c.Port)
	}
	if c.NoTLS {
		return c.Addr + ":80"
	}
	return c.Addr + ":443"
}
