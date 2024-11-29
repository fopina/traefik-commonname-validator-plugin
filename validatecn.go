// Package traefik_commonname_validator_plugin a plugin to only allow some client certificate Subject CNs.
package traefik_commonname_validator_plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Allowed []string `json:"allowed,omitempty"`
	Debug   bool     `json:"debug,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Allowed: make([]string, 1),
		Debug:   false,
	}
}

// ValidateCN a ValidateCN plugin.
type ValidateCN struct {
	next    http.Handler
	allowed []string
	debug   bool
	name    string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Debug {
		log.Printf("Configuration: %v", config)
	}
	if len(config.Allowed) == 0 {
		return nil, fmt.Errorf("allowed cannot be empty")
	}

	return &ValidateCN{
		allowed: config.Allowed,
		debug:   config.Debug,
		next:    next,
		name:    name,
	}, nil
}

func (p *ValidateCN) getCertInfo(tls *tls.ConnectionState) string {
	if tls == nil {
		return ""
	}
	// we only care about the first cert of the chain, the leaf
	// rest of the chain has (or should have) been validated by Traefik mTLS
	if len(tls.PeerCertificates) == 0 {
		return ""
	}
	return tls.PeerCertificates[0].Subject.CommonName
}

func (p *ValidateCN) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	actualCN := p.getCertInfo(req.TLS)
	for _, allowedCN := range p.allowed {
		if actualCN == allowedCN {
			p.next.ServeHTTP(rw, req)
			return
		}
	}

	if p.debug {
		log.Printf("REJECTED: %s not part of %s", actualCN, p.allowed)
	}
	rw.WriteHeader(http.StatusForbidden)
	fmt.Fprintln(rw, "Forbidden")
}
