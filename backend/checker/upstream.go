/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 *
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

package checker

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

// ProxyType represents the type of proxy (duplicate with manager.go)
/* type ProxyType string

const (
	HTTP    ProxyType = "http"
	HTTPS   ProxyType = "https"
	SOCKS4  ProxyType = "socks4"
	SOCKS5  ProxyType = "socks5"
	UNKNOWN ProxyType = "unknown"
) */

// UpstreamProxy represents a proxy that will be used to route all proxy checks
type UpstreamProxy struct {
	Address string
	Type    ProxyType
	Timeout time.Duration
}

// NewUpstreamProxy creates a new upstream proxy configuration
func NewUpstreamProxy(address string, proxyType ProxyType, timeout time.Duration) *UpstreamProxy {
	return &UpstreamProxy{
		Address: address,
		Type:    proxyType,
		Timeout: timeout,
	}
}

// CreateDialer creates a dialer that routes connections through the upstream proxy
func (up *UpstreamProxy) CreateDialer() (proxy.Dialer, error) {
	if up.Address == "" {
		// If no upstream proxy is specified, return a direct dialer
		return &net.Dialer{Timeout: up.Timeout}, nil
	}

	return createUpstreamDialer(up.Address, up.Type, up.Timeout)
}

// CreateHTTPTransport creates an HTTP transport that routes connections through the upstream proxy
func (up *UpstreamProxy) CreateHTTPTransport() (*http.Transport, error) {
	if up.Address == "" {
		// If no upstream proxy is specified, return a direct transport
		return &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   up.Timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   up.Timeout,
			ResponseHeaderTimeout: up.Timeout,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
		}, nil
	}

	// Create a dialer that uses the upstream proxy
	upstreamDialer, err := up.CreateDialer()
	if err != nil {
		return nil, fmt.Errorf("failed to create upstream dialer: %w", err)
	}

	// Create a transport that uses the upstream dialer
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return upstreamDialer.Dial(network, addr)
		},
		TLSHandshakeTimeout:   up.Timeout,
		ResponseHeaderTimeout: up.Timeout,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
	}

	return transport, nil
}

// CreateHTTPClient creates an HTTP client that routes connections through the upstream proxy
func (up *UpstreamProxy) CreateHTTPClient() (*http.Client, error) {
	transport, err := up.CreateHTTPTransport()
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: transport,
		Timeout:   up.Timeout,
	}, nil
}

// TestUpstreamConnection tests if the upstream proxy is working
func (up *UpstreamProxy) TestUpstreamConnection(endpoint string) (string, error) {
	if up.Address == "" {
		return "", fmt.Errorf("no upstream proxy specified")
	}

	// Create a client that uses the upstream proxy
	client, err := up.CreateHTTPClient()
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Make a request to the endpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add common headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upstream proxy connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// The response should contain the outgoing IP
	outgoingIP := string(body)
	if outgoingIP == "" {
		return "", ErrEmptyResponse
	}

	return outgoingIP, nil
}

// GetProxyTypeFromString converts a string to a ProxyType
func GetProxyTypeFromString(proxyType string) ProxyType {
	switch proxyType {
	case "http":
		return HTTP
	case "https":
		return HTTPS
	case "socks4":
		return SOCKS4
	case "socks5":
		return SOCKS5
	default:
		return UNKNOWN
	}
}

// String returns the string representation of a ProxyType
func (pt ProxyType) String() string {
	return string(pt)
}

// IsValid checks if the ProxyType is valid
func (pt ProxyType) IsValid() bool {
	return pt == HTTP || pt == HTTPS || pt == SOCKS4 || pt == SOCKS5
}
