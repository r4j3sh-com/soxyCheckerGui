/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

package checker

import (
	"time"
)

// ProxyStatus represents the status of a proxy check
type ProxyStatus string

const (
	// StatusPending indicates the proxy check has not started yet
	StatusPending ProxyStatus = "pending"

	// StatusChecking indicates the proxy check is in progress
	StatusChecking ProxyStatus = "checking"

	// StatusLive indicates the proxy is working
	StatusLive ProxyStatus = "live"

	// StatusDead indicates the proxy is not working
	StatusDead ProxyStatus = "dead"

	// StatusError indicates an error occurred during the proxy check
	StatusError ProxyStatus = "error"
)

// ProxyResult represents the result of a proxy check
type ProxyResult struct {
	// Proxy is the proxy address in format ip:port
	Proxy string `json:"proxy"`

	// Type is the detected or specified proxy type
	Type ProxyType `json:"type"`

	// Status is the current status of the proxy
	Status ProxyStatus `json:"status"`

	// Latency is the time it took to check the proxy in milliseconds
	Latency int64 `json:"latency"`

	// OutgoingIP is the IP address seen by the endpoint when using this proxy
	OutgoingIP string `json:"outgoingIp"`

	// Country is the country of the proxy (if geolocation is enabled)
	Country string `json:"country"`

	// CountryCode is the ISO country code of the proxy (if geolocation is enabled)
	CountryCode string `json:"countryCode"`

	// Error is the error message if the proxy check failed
	Error string `json:"error"`

	// Timestamp is when the check was completed
	Timestamp time.Time `json:"timestamp"`

	// Anonymous indicates if the proxy is anonymous (doesn't reveal your IP)
	Anonymous bool `json:"anonymous"`

	// SupportsHTTPS indicates if the proxy supports HTTPS connections
	SupportsHTTPS bool `json:"supportsHttps"`
}

// NewPendingResult creates a new ProxyResult with status pending
func NewPendingResult(proxy string, proxyType ProxyType) *ProxyResult {
	return &ProxyResult{
		Proxy:     proxy,
		Type:      proxyType,
		Status:    StatusPending,
		Timestamp: time.Now(),
	}
}

// SetChecking updates the result status to checking
func (r *ProxyResult) SetChecking() {
	r.Status = StatusChecking
	r.Timestamp = time.Now()
}

// SetLive updates the result to indicate a successful check
func (r *ProxyResult) SetLive(latency int64, outgoingIP string) {
	r.Status = StatusLive
	r.Latency = latency
	r.OutgoingIP = outgoingIP
	r.Error = ""
	r.Timestamp = time.Now()
}

// SetDead updates the result to indicate a failed check
func (r *ProxyResult) SetDead(err string) {
	r.Status = StatusDead
	r.Error = err
	r.Timestamp = time.Now()
}

// SetError updates the result to indicate an error during check
func (r *ProxyResult) SetError(err string) {
	r.Status = StatusError
	r.Error = err
	r.Timestamp = time.Now()
}

// SetType updates the proxy type
func (r *ProxyResult) SetType(proxyType ProxyType) {
	r.Type = proxyType
	r.Timestamp = time.Now()
}

// SetGeoInfo updates the geolocation information
func (r *ProxyResult) SetGeoInfo(country string, countryCode string) {
	r.Country = country
	r.CountryCode = countryCode
}

// SetAnonymous updates the anonymity status
func (r *ProxyResult) SetAnonymous(anonymous bool) {
	r.Anonymous = anonymous
}

// SetSupportsHTTPS updates whether the proxy supports HTTPS
func (r *ProxyResult) SetSupportsHTTPS(supportsHTTPS bool) {
	r.SupportsHTTPS = supportsHTTPS
}

// Clone creates a copy of the ProxyResult
func (r *ProxyResult) Clone() *ProxyResult {
	return &ProxyResult{
		Proxy:         r.Proxy,
		Type:          r.Type,
		Status:        r.Status,
		Latency:       r.Latency,
		OutgoingIP:    r.OutgoingIP,
		Country:       r.Country,
		CountryCode:   r.CountryCode,
		Error:         r.Error,
		Timestamp:     r.Timestamp,
		Anonymous:     r.Anonymous,
		SupportsHTTPS: r.SupportsHTTPS,
	}
}

// ProxyResultList is a list of ProxyResult objects
type ProxyResultList []*ProxyResult

// Clone creates a deep copy of the ProxyResultList
func (l ProxyResultList) Clone() ProxyResultList {
	if l == nil {
		return nil
	}

	result := make(ProxyResultList, len(l))
	for i, r := range l {
		result[i] = r.Clone()
	}

	return result
}

// FilterByStatus returns a new list containing only results with the specified status
func (l ProxyResultList) FilterByStatus(status ProxyStatus) ProxyResultList {
	var result ProxyResultList

	for _, r := range l {
		if r.Status == status {
			result = append(result, r)
		}
	}

	return result
}

// FilterByType returns a new list containing only results with the specified type
func (l ProxyResultList) FilterByType(proxyType ProxyType) ProxyResultList {
	var result ProxyResultList

	for _, r := range l {
		if r.Type == proxyType {
			result = append(result, r)
		}
	}

	return result
}

// GetLiveProxies returns a list of working proxy addresses (ip:port format)
func (l ProxyResultList) GetLiveProxies() []string {
	var result []string

	for _, r := range l {
		if r.Status == StatusLive {
			result = append(result, r.Proxy)
		}
	}

	return result
}

// GetLiveProxiesWithType returns a list of working proxy addresses with their types
// Format: "type://ip:port"
func (l ProxyResultList) GetLiveProxiesWithType() []string {
	var result []string

	for _, r := range l {
		if r.Status == StatusLive {
			result = append(result, string(r.Type)+"://"+r.Proxy)
		}
	}

	return result
}
