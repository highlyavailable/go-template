package clients

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// NewInsecureClient creates an HTTP client that does not verify TLS certificates
func NewInsecureClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// NewSecureClient creates an HTTP client that uses proxies and verifies TLS certificates
func NewSecureClient(proxyURL string, certFile string, keyFile string) (*http.Client, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Parse proxy URL
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: tlsConfig,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}
