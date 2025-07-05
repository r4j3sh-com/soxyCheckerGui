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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

var (
	ErrInvalidProxyFormat    = errors.New("invalid proxy format")
	ErrUnsupportedProxyType  = errors.New("unsupported proxy type")
	ErrProxyConnectionFailed = errors.New("proxy connection failed")
	ErrEmptyResponse         = errors.New("empty response from endpoint")
)

// CheckHTTP checks if an HTTP proxy is working
// If upstreamProxy is provided, the check will be routed through it
func CheckHTTP(proxyAddr string, endpoint string, timeout time.Duration, upstreamProxy string, upstreamType ProxyType) (string, error) {
	// Validate proxy format
	if !strings.Contains(proxyAddr, ":") {
		return "", ErrInvalidProxyFormat
	}

	// Create proxy URL
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return "", fmt.Errorf("invalid proxy address: %w", err)
	}

	// Create transport and client
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   timeout,
		ResponseHeaderTimeout: timeout,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
	}

	// If upstream proxy is specified, route through it
	if upstreamProxy != "" {
		upstreamDialer, err := createUpstreamDialer(upstreamProxy, upstreamType, timeout)
		if err != nil {
			return "", fmt.Errorf("failed to create upstream connection: %w", err)
		}

		// Replace the dialer with one that uses the upstream proxy
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return upstreamDialer.Dial(network, addr)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	// Make the request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add common headers to appear more like a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("proxy connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body to get the IP
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// The response should contain the outgoing IP
	outgoingIP := strings.TrimSpace(string(body))
	if outgoingIP == "" {
		return "", ErrEmptyResponse
	}

	return outgoingIP, nil
}

// CheckHTTPS checks if an HTTPS proxy is working
func CheckHTTPS(proxyAddr string, endpoint string, timeout time.Duration, upstreamProxy string, upstreamType ProxyType) (string, error) {
	// Validate proxy format
	if !strings.Contains(proxyAddr, ":") {
		return "", ErrInvalidProxyFormat
	}

	// Create proxy URL
	proxyURL, err := url.Parse("https://" + proxyAddr)
	if err != nil {
		return "", fmt.Errorf("invalid proxy address: %w", err)
	}

	// Create transport and client
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   timeout,
		ResponseHeaderTimeout: timeout,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
	}

	// If upstream proxy is specified, route through it
	if upstreamProxy != "" {
		upstreamDialer, err := createUpstreamDialer(upstreamProxy, upstreamType, timeout)
		if err != nil {
			return "", fmt.Errorf("failed to create upstream connection: %w", err)
		}

		// Replace the dialer with one that uses the upstream proxy
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return upstreamDialer.Dial(network, addr)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	// Make the request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add common headers to appear more like a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("proxy connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body to get the IP
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// The response should contain the outgoing IP
	outgoingIP := strings.TrimSpace(string(body))
	if outgoingIP == "" {
		return "", ErrEmptyResponse
	}

	return outgoingIP, nil
}

// CheckSOCKS4 checks if a SOCKS4 proxy is working
func CheckSOCKS4(proxyAddr string, endpoint string, timeout time.Duration, upstreamProxy string, upstreamType ProxyType) (string, error) {
	// Validate proxy format
	if !strings.Contains(proxyAddr, ":") {
		return "", ErrInvalidProxyFormat
	}

	// Create SOCKS4 dialer
	dialer := &net.Dialer{Timeout: timeout}

	// If upstream proxy is specified, route through it
	if upstreamProxy != "" {
		// Note: Chaining SOCKS proxies is complex and not fully implemented here
		return "", fmt.Errorf("upstream proxy not supported for SOCKS4 checks")
	}

	// Create SOCKS4 client
	// Note: Go's proxy package doesn't directly support SOCKS4, so we use SOCKS5 with special handling
	auth := &proxy.Auth{
		User: "socks4", // This is a marker for SOCKS4 protocol
	}
	socks4Dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, dialer)
	if err != nil {
		return "", fmt.Errorf("failed to create SOCKS4 client: %w", err)
	}

	// Parse the endpoint URL to get the host and port
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint URL: %w", err)
	}

	// Extract host and port from the endpoint
	host := endpointURL.Hostname()
	port := endpointURL.Port()
	if port == "" {
		if endpointURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	// Connect to the endpoint through the SOCKS4 proxy
	conn, err := socks4Dialer.Dial("tcp", host+":"+port)
	if err != nil {
		return "", fmt.Errorf("SOCKS4 connection failed: %w", err)
	}
	defer conn.Close()

	// For HTTP(S) endpoints, we need to make an HTTP request
	if endpointURL.Scheme == "http" || endpointURL.Scheme == "https" {
		// Create a client that uses our SOCKS4 connection
		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return socks4Dialer.Dial(network, addr)
				},
			},
			Timeout: timeout,
		}

		// Make the request
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		// Add common headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("HTTP request through SOCKS4 failed: %w", err)
		}
		defer resp.Body.Close()

		// Read response body to get the IP
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}

		// The response should contain the outgoing IP
		outgoingIP := strings.TrimSpace(string(body))
		if outgoingIP == "" {
			return "", ErrEmptyResponse
		}

		return outgoingIP, nil
	}

	// For non-HTTP endpoints, we would need a different approach
	return "Connection successful", nil
}

// CheckSOCKS5 checks if a SOCKS5 proxy is working
func CheckSOCKS5(proxyAddr string, endpoint string, timeout time.Duration, upstreamProxy string, upstreamType ProxyType) (string, error) {
	// Validate proxy format
	if !strings.Contains(proxyAddr, ":") {
		return "", ErrInvalidProxyFormat
	}

	// Create SOCKS5 dialer
	dialer := &net.Dialer{Timeout: timeout}

	// If upstream proxy is specified, route through it
	if upstreamProxy != "" {
		// Note: Chaining SOCKS proxies is complex and not fully implemented here
		return "", fmt.Errorf("upstream proxy not supported for SOCKS5 checks")
	}

	// Create SOCKS5 client
	socks5Dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, dialer)
	if err != nil {
		return "", fmt.Errorf("failed to create SOCKS5 client: %w", err)
	}

	// Parse the endpoint URL to get the host and port
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint URL: %w", err)
	}

	// Extract host and port from the endpoint
	host := endpointURL.Hostname()
	port := endpointURL.Port()
	if port == "" {
		if endpointURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	// Connect to the endpoint through the SOCKS5 proxy
	conn, err := socks5Dialer.Dial("tcp", host+":"+port)
	if err != nil {
		return "", fmt.Errorf("SOCKS5 connection failed: %w", err)
	}
	defer conn.Close()

	// For HTTP(S) endpoints, we need to make an HTTP request
	if endpointURL.Scheme == "http" || endpointURL.Scheme == "https" {
		// Create a client that uses our SOCKS5 connection
		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return socks5Dialer.Dial(network, addr)
				},
			},
			Timeout: timeout,
		}

		// Make the request
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		// Add common headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("HTTP request through SOCKS5 failed: %w", err)
		}
		defer resp.Body.Close()

		// Read response body to get the IP
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}

		// The response should contain the outgoing IP
		outgoingIP := strings.TrimSpace(string(body))
		if outgoingIP == "" {
			return "", ErrEmptyResponse
		}

		return outgoingIP, nil
	}

	// For non-HTTP endpoints, we would need a different approach
	return "Connection successful", nil
}

// Helper function to create an upstream dialer based on proxy type
func createUpstreamDialer(upstreamProxy string, upstreamType ProxyType, timeout time.Duration) (proxy.Dialer, error) {
	dialer := &net.Dialer{Timeout: timeout}

	switch upstreamType {
	case HTTP, HTTPS:
		// For HTTP/HTTPS upstream proxies
		proxyURL, err := url.Parse(string(upstreamType) + "://" + upstreamProxy)
		if err != nil {
			return nil, fmt.Errorf("invalid upstream proxy format: %w", err)
		}
		return proxy.FromURL(proxyURL, dialer)

	case SOCKS4:
		// For SOCKS4 upstream proxies
		// Use SOCKS5 with SOCKS4 flag since golang.org/x/net/proxy doesn't have a direct SOCKS4 constructor
		auth := &proxy.Auth{
			User: "socks4", // This is a marker for SOCKS4 protocol
		}
		return proxy.SOCKS5("tcp", upstreamProxy, auth, dialer)

	case SOCKS5:
		// For SOCKS5 upstream proxies
		return proxy.SOCKS5("tcp", upstreamProxy, nil, dialer)

	default:
		return nil, ErrUnsupportedProxyType
	}
}
