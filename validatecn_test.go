package traefik_commonname_validator_plugin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
)

// certificates by example/docker-compose.yml
const (
	signingCA = `-----BEGIN CERTIFICATE-----
MIIB4jCCAWigAwIBAgIIPcdDGEihqqgwCgYIKoZIzj0EAwMwEzERMA8GA1UEAxMI
Z29vZC1vbmUwIBcNMjQxMTI3MTQ0ODA1WhgPMjEyNDExMjcxNDQ4MDVaMBMxETAP
BgNVBAMTCGdvb2Qtb25lMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEswhpvhWehT5n
OMqltZllVHz07KNdEcn5+PcWYjwQF/zWomCQE5E8i5+gTNgLDiXyy/WmfjnmWwER
pwihmg9gZyaooRnIB0NV8qtjcJDAFjC2yO4fotE+ABksV3am3KpBo4GGMIGDMA4G
A1UdDwEB/wQEAwIChDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYD
VR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUyIHhTZUqC3CJxxLXMJd9qyAZqdgw
HwYDVR0jBBgwFoAUyIHhTZUqC3CJxxLXMJd9qyAZqdgwCgYIKoZIzj0EAwMDaAAw
ZQIwKkA86fB+s+pQk3vHTZCpfaV0Vxtv/cnXWzx9gWWZ6R5xuhhvhqmie8q+gmyf
0hObAjEA+lTAjTYncNri9jdab7NehSJfozs0Hd2Eubn/NDI7+TDBpFG2TrRywgnU
gUWxwSNQ
-----END CERTIFICATE-----`
	authClientCrt = `-----BEGIN CERTIFICATE-----
MIIB0zCCAVqgAwIBAgIIUtj0W9Ynb+EwCgYIKoZIzj0EAwMwEzERMA8GA1UEAxMI
Z29vZC1vbmUwHhcNMjQxMTI3MTQ0ODA1WhcNMjYxMjI3MTQ0ODA1WjAWMRQwEgYD
VQQDEwthdXRoLWNsaWVudDB2MBAGByqGSM49AgEGBSuBBAAiA2IABEE2qaeXUact
1fjVILwjoN6AYdA4r9yh7dyRzqcseIlKkWpxZsi+HrQ1qu9zlAIUQG7k3B5TuPOW
D6Md+gRY7hAjtwMSkUFFTAMrz1wmJ7HdbBC0T0iiFzJTMiclpnI7cKN4MHYwDgYD
VR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNV
HRMBAf8EAjAAMB8GA1UdIwQYMBaAFMiB4U2VKgtwiccS1zCXfasgGanYMBYGA1Ud
EQQPMA2CC2F1dGgtY2xpZW50MAoGCCqGSM49BAMDA2cAMGQCMGn7nUdtJOISt+7+
JaSbAtDXMbPXNn9nyPx0iFDAVugSx7SCzb+8aHYEt8DRk6hkVwIwJwIhNiTMqsp/
PJRTbVciRRBDI0rlXk1w6ld9ajiq6BGxY8tVPzqKxfKlL1bfpJKF
-----END CERTIFICATE-----`
	authClient2Crt = `-----BEGIN CERTIFICATE-----
MIIB1jCCAVygAwIBAgIIEKdMGC8j5L4wCgYIKoZIzj0EAwMwEzERMA8GA1UEAxMI
Z29vZC1vbmUwHhcNMjQxMTI3MTQ0ODA1WhcNMjYxMjI3MTQ0ODA1WjAXMRUwEwYD
VQQDEwxhdXRoMi1jbGllbnQwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAAQP6Dlr5mB4
rnpFVHe0kGvZ5lYq5Kvft9NnnujWYgDqTPzzviOZemvXZShQQ996ndpDuLbXrAp9
fzw/77CWyPgEaLZ/GWCy5Pu84PCRM1O7U8WT9Na8sxKvOzC1f7bbS66jeTB3MA4G
A1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYD
VR0TAQH/BAIwADAfBgNVHSMEGDAWgBTIgeFNlSoLcInHEtcwl32rIBmp2DAXBgNV
HREEEDAOggxhdXRoMi1jbGllbnQwCgYIKoZIzj0EAwMDaAAwZQIwImC5CAxSFWr8
fhT54brSgEbr8lFxJGeo4OzWzB12lP1kvb+PmT54QaHmOFeTSqOfAjEAoWCfhkWh
NWvmzZuu6yvspxgA1F7S0wvFT9cqDRyZsWW/b0xvl1HOH0/EJ0wc0WkR
-----END CERTIFICATE-----`
)

func TestValidateCN(t *testing.T) {
	cfg := CreateConfig()
	assertEqual(t, len(cfg.Allowed), 1)
	cfg.Allowed = []string{"auth-client"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "test-plugin")
	assertNoError(t, err)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	assertNoError(t, err)
	req.TLS = buildTLSWith([]string{authClientCrt})
	handler.ServeHTTP(recorder, req)
	assertEqual(t, recorder.Result().StatusCode, 200)

	recorder = httptest.NewRecorder()
	req.TLS = buildTLSWith([]string{authClient2Crt})
	handler.ServeHTTP(recorder, req)
	assertEqual(t, recorder.Result().StatusCode, 403)
}

func TestValidateCNBothAllowed(t *testing.T) {
	cfg := CreateConfig()
	assertEqual(t, len(cfg.Allowed), 1)
	cfg.Allowed = []string{"auth-client", "auth2-client"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "test-plugin")
	assertNoError(t, err)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	assertNoError(t, err)
	req.TLS = buildTLSWith([]string{authClientCrt, signingCA})
	handler.ServeHTTP(recorder, req)
	assertEqual(t, recorder.Result().StatusCode, 200)

	recorder = httptest.NewRecorder()
	req.TLS = buildTLSWith([]string{authClient2Crt, signingCA})
	handler.ServeHTTP(recorder, req)
	assertEqual(t, recorder.Result().StatusCode, 200)
}

type Equatable interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~string | ~bool
}

func assertEqual[TT Equatable](t *testing.T, expected, actual TT) {
	t.Helper()
	if expected != actual {
		t.Errorf("got status code %v, want %v", expected, actual)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// buildTLSWith copied from https://github.com/traefik/traefik/blob/master/pkg/middlewares/passtlsclientcert/pass_tls_client_cert_test.go
func buildTLSWith(certContents []string) *tls.ConnectionState {
	var peerCertificates []*x509.Certificate

	for _, certContent := range certContents {
		peerCertificates = append(peerCertificates, getCertificate(certContent))
	}

	return &tls.ConnectionState{PeerCertificates: peerCertificates}
}

// getCertificate copied from https://github.com/traefik/traefik/blob/master/pkg/middlewares/passtlsclientcert/pass_tls_client_cert_test.go
func getCertificate(certContent string) *x509.Certificate {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(signingCA))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(certContent))
	if block == nil {
		panic("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	return cert
}
