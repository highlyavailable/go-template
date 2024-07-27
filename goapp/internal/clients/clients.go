package clients

import (
	"crypto/tls"
	"fmt"
	"goapp/internal/logging"
	"net"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// NewInsecureClient creates an HTTP client that does not verify TLS certificates
func NewInsecureClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	err := testClient(client)
	if err != nil {
		logging.Error("Error testing client", zap.Error(err))
		fmt.Println("Error testing client:", err)
		return nil, err
	}

	return client, nil
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

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	err = testClient(client)
	if err != nil {
		logging.Error("Error testing client", zap.Error(err))
		fmt.Println("Error testing client:", err)
		return nil, err
	}

	return client, nil
}

func testClient(c *http.Client) error {
	// Create a new insecure client
	resp, err := c.Get("https://www.google.com")
	if err != nil {
		logging.Error("Error making request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response status code: %d", resp.StatusCode)
	}

	logging.Info("Client test successful")
	return nil
}
