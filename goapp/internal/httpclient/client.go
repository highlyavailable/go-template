package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"goapp/internal/logging"
)

// Config holds HTTP client configuration
type Config struct {
	// Timeouts
	Timeout     time.Duration `envconfig:"TIMEOUT" default:"30s"`
	DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" default:"10s"`
	TLSTimeout  time.Duration `envconfig:"TLS_TIMEOUT" default:"10s"`

	// Connection pooling
	MaxIdleConns        int           `envconfig:"MAX_IDLE_CONNS" default:"100"`
	MaxIdleConnsPerHost int           `envconfig:"MAX_IDLE_CONNS_PER_HOST" default:"10"`
	IdleConnTimeout     time.Duration `envconfig:"IDLE_CONN_TIMEOUT" default:"90s"`

	// TLS
	InsecureSkipVerify bool   `envconfig:"INSECURE_SKIP_VERIFY" default:"false"`
	CertFile           string `envconfig:"CERT_FILE"`
	KeyFile            string `envconfig:"KEY_FILE"`

	// Proxy
	ProxyURL string `envconfig:"PROXY_URL"`

	// Retry
	MaxRetries   int           `envconfig:"MAX_RETRIES" default:"3"`
	RetryWaitMin time.Duration `envconfig:"RETRY_WAIT_MIN" default:"1s"`
	RetryWaitMax time.Duration `envconfig:"RETRY_WAIT_MAX" default:"30s"`

	// Headers
	UserAgent string            `envconfig:"USER_AGENT" default:"goapp/1.0"`
	Headers   map[string]string `envconfig:"HEADERS"`
}

// Client wraps http.Client with enterprise features
type Client struct {
	client *http.Client
	config Config
	logger logging.Logger
}

// New creates a new enterprise HTTP client
func New(cfg Config, logger logging.Logger) (*Client, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   cfg.DialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   cfg.TLSTimeout,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		DisableCompression:    false,
		ResponseHeaderTimeout: cfg.Timeout,
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		MinVersion:         tls.VersionTLS12, // Enterprise security: minimum TLS 1.2
	}

	// Load client certificates if provided
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	transport.TLSClientConfig = tlsConfig

	// Configure proxy if provided
	if cfg.ProxyURL != "" {
		proxyFunc, err := createProxyFunc(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to configure proxy: %w", err)
		}
		transport.Proxy = proxyFunc
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Log the redirect
			logger.Debugf("Redirecting to %s", req.URL.String())

			// Limit redirects for security
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	return &Client{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// Do executes an HTTP request with retry logic and logging
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Add default headers
	c.addDefaultHeaders(req)

	var resp *http.Response
	var err error

	// Retry logic
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			waitTime := c.calculateBackoff(attempt)
			c.logger.Infof("Retrying request after %v (attempt %d/%d)", waitTime, attempt, c.config.MaxRetries)
			time.Sleep(waitTime)
		}

		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			// Success or client error (don't retry)
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt < c.config.MaxRetries {
			c.logger.Warnf("Request failed, will retry: %v", err)
		}
	}

	if err != nil {
		c.logger.Errorf("Request failed after %d attempts: %v", c.config.MaxRetries+1, err)
		return nil, err
	}

	c.logger.Debugf("Request completed: %s %s -> %d", req.Method, req.URL.String(), resp.StatusCode)
	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, url, contentType string, body interface{}) (*http.Response, error) {
	req, err := createRequestWithBody(ctx, http.MethodPost, url, contentType, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	return c.Do(req)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, url, contentType string, body interface{}) (*http.Response, error) {
	req, err := createRequestWithBody(ctx, http.MethodPut, url, contentType, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	return c.Do(req)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}
	return c.Do(req)
}

// HealthCheck performs a health check against a URL
func (c *Client) HealthCheck(ctx context.Context, url string) error {
	resp, err := c.Get(ctx, url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("health check failed: status %d", resp.StatusCode)
	}

	return nil
}

// Close cleans up the client resources
func (c *Client) Close() {
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

// addDefaultHeaders adds default headers to the request
func (c *Client) addDefaultHeaders(req *http.Request) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}

	// Add custom headers
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}
}

// calculateBackoff calculates exponential backoff with jitter
func (c *Client) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: min * (2 ^ attempt)
	backoff := c.config.RetryWaitMin * time.Duration(1<<uint(attempt-1))
	if backoff > c.config.RetryWaitMax {
		backoff = c.config.RetryWaitMax
	}
	return backoff
}

// createProxyFunc creates a proxy function from URL string
func createProxyFunc(proxyURL string) (func(*http.Request) (*url.URL, error), error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	return http.ProxyURL(u), nil
}
