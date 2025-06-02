package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// createRequestWithBody creates an HTTP request with a body
func createRequestWithBody(ctx context.Context, method, url, contentType string, body interface{}) (*http.Request, error) {
	var reader io.Reader
	
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			reader = v
		case string:
			reader = strings.NewReader(v)
		case []byte:
			reader = bytes.NewReader(v)
		default:
			// Assume JSON serializable
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body to JSON: %w", err)
			}
			reader = bytes.NewReader(jsonData)
			if contentType == "" {
				contentType = "application/json"
			}
		}
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, err
	}
	
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	
	return req, nil
}

// DecodeJSON decodes JSON response body into target
func DecodeJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}
	
	return json.NewDecoder(resp.Body).Decode(target)
}

// ReadBody reads the response body as bytes
func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ReadBodyAsString reads the response body as string
func ReadBodyAsString(resp *http.Response) (string, error) {
	body, err := ReadBody(resp)
	if err != nil {
		return "", err
	}
	return string(body), nil
}