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
	"net"
	"net/http"
	"net/url"
	"time"

	socks "golang.org/x/net/proxy"
)

// DetectProxyType attempts to automatically detect the type of proxy
// It tries each protocol in order: SOCKS5, SOCKS4, HTTPS, HTTP
func DetectProxyType(proxy string, timeout time.Duration) (ProxyType, error) {
	// Try each protocol in sequence
	protocols := []struct {
		checkFunc func(string, time.Duration) bool
		proxyType ProxyType
	}{
		{checkSOCKS5Quick, SOCKS5},
		{checkSOCKS4Quick, SOCKS4},
		{checkHTTPSQuick, HTTPS},
		{checkHTTPQuick, HTTP},
	}

	for _, protocol := range protocols {
		if protocol.checkFunc(proxy, timeout) {
			return protocol.proxyType, nil
		}
	}

	return "", fmt.Errorf("could not detect proxy type")
}

// Quick check functions for auto-detection

// checkHTTPQuick performs a quick check to see if a proxy supports HTTP
func checkHTTPQuick(proxy string, timeout time.Duration) bool {
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		return false
	}

	// Create a transport with the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
		IdleConnTimeout:     timeout,
	}

	// Create a client with the transport
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	// Try to connect to a known endpoint
	req, err := http.NewRequest("HEAD", "http://www.google.com", nil)
	if err != nil {
		return false
	}

	// Set a short timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// If we got a response, the proxy is working
	return resp.StatusCode >= 200 && resp.StatusCode < 500
}

// checkHTTPSQuick performs a quick check to see if a proxy supports HTTPS
func checkHTTPSQuick(proxy string, timeout time.Duration) bool {
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		return false
	}

	// Create a transport with the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
		IdleConnTimeout:     timeout,
	}

	// Create a client with the transport
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	// Try to connect to a known HTTPS endpoint
	req, err := http.NewRequest("HEAD", "https://www.google.com", nil)
	if err != nil {
		return false
	}

	// Set a short timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// If we got a response, the proxy is working
	return resp.StatusCode >= 200 && resp.StatusCode < 500
}

// checkSOCKS4Quick performs a quick check to see if a proxy supports SOCKS4
func checkSOCKS4Quick(proxy string, timeout time.Duration) bool {
	// Parse the proxy address (host, port, err)
	_, _, err := net.SplitHostPort(proxy)
	if err != nil {
		return false
	}

	// Create a SOCKS4 dialer
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// Try to connect to the proxy
	conn, err := dialer.Dial("tcp", proxy)
	if err != nil {
		return false
	}
	defer conn.Close()

	// Prepare SOCKS4 handshake
	// SOCKS4 request format:
	// VN(1) | CD(1) | DSTPORT(2) | DSTIP(4) | USERID(variable) | NULL(1)
	// VN = 4 (SOCKS version)
	// CD = 1 (connect command)
	request := []byte{
		4,     // SOCKS version
		1,     // CONNECT command
		0, 80, // Port 80
		8, 8, 8, 8, // IP (8.8.8.8)
		0, // User ID (empty)
	}

	// Set a deadline for the connection
	conn.SetDeadline(time.Now().Add(timeout))

	// Send the request
	_, err = conn.Write(request)
	if err != nil {
		return false
	}

	// Read the response
	response := make([]byte, 8)
	_, err = conn.Read(response)
	if err != nil {
		return false
	}

	// Check if the response indicates success
	// SOCKS4 response format:
	// VN(1) | CD(1) | DSTPORT(2) | DSTIP(4)
	// CD = 90 (request granted)
	return response[1] == 90
}

// checkSOCKS5Quick performs a quick check to see if a proxy supports SOCKS5
func checkSOCKS5Quick(proxy string, timeout time.Duration) bool {
	// Create a SOCKS5 dialer
	dialer, err := socks.SOCKS5("tcp", proxy, nil, &net.Dialer{
		Timeout: timeout,
	})
	if err != nil {
		return false
	}

	// Try to connect to a known endpoint
	conn, err := dialer.Dial("tcp", "www.google.com:80")
	if err != nil {
		return false
	}
	defer conn.Close()

	// If we got here, the proxy is working
	return true
}
