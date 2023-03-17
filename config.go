package main

import (
	"strings"

	auth "github.com/korylprince/go-ad-auth/v3"
)

// Config configure chronicle-ui
type Config struct {
	LDAPServer   string `required:"true"`
	LDAPPort     int    `required:"true" default:"389"`
	LDAPBaseDN   string `required:"true"`
	LDAPGroup    string `required:"true"`
	LDAPSecurity string `required:"true" default:"none"`

	SQLDSN string `required:"true"`

	ProxyHeaders bool   `default:"false"`
	ListenAddr   string `default:":80"`
}

// AuthConfig returns the auth.Config from Config
func (c *Config) AuthConfig() *auth.Config {
	config := &auth.Config{
		Server: c.LDAPServer,
		Port:   c.LDAPPort,
		BaseDN: c.LDAPBaseDN,
	}
	switch strings.ToLower(c.LDAPSecurity) {
	case "tls":
		config.Security = auth.SecurityTLS
	case "starttls":
		config.Security = auth.SecurityStartTLS
	default:
		config.Security = auth.SecurityNone
	}
	return config
}
