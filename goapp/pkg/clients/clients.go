package clients

import (
	"crypto/tls"
	"fmt"
	"goapp/pkg/logging"
	"net"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type SecureClientConfig struct {
	CertFile       string // Path to the client certificate
	KeyFile        string // Path to the client key
	ProxyURL       string // URL of the proxy server
	URLForConnTest string // URL to test the client connection
}

// NewInsecureClient creates an HTTP client that does not verify TLS certificates
func NewInsecureClient(urlForRequest string) (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	err := testClient(client, urlForRequest)
	if err != nil {
		logging.Logger.Error("Error testing client", zap.Error(err))
		fmt.Println("Error testing client:", err)
		return nil, err
	}

	return client, nil
}

// NewSecureClient creates an HTTP client that uses proxies and verifies TLS certificates
func NewSecureClient(clientConfig SecureClientConfig) (*http.Client, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientConfig.CertFile, clientConfig.KeyFile)
	if err != nil {
		return nil, err
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Parse proxy URL
	proxy, err := url.Parse(clientConfig.ProxyURL)
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

	err = testClient(client, clientConfig.URLForConnTest)
	if err != nil {
		logging.Logger.Error("Error testing client", zap.Error(err))
		fmt.Println("Error testing client:", err)
		return nil, err
	}

	return client, nil
}

func testClient(c *http.Client, urlForRequest string) error {
	if urlForRequest == "" {
		urlForRequest = "https://www.google.com"
	}

	// assert valid URL
	_, err := url.ParseRequestURI(urlForRequest)
	if err != nil {
		logging.Logger.Error("Error parsing URL", zap.Error(err))
		return err
	}

	// Create a new insecure client
	resp, err := c.Get(urlForRequest)
	if err != nil {
		logging.Logger.Error("Error making request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response status code: %d", resp.StatusCode)
	}

	logging.Logger.Info("Client test successful")
	return nil
}
