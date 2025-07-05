/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 *
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

package checker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ProxyType represents the type of proxy
type ProxyType string

const (
	Auto    ProxyType = "auto"
	HTTP    ProxyType = "http"
	HTTPS   ProxyType = "https"
	SOCKS4  ProxyType = "socks4"
	SOCKS5  ProxyType = "socks5"
	UNKNOWN ProxyType = "unknown"
)

// ProxyCheckRequest represents a request to check proxies
type ProxyCheckRequest struct {
	ProxyList     []string  // List of proxies to check (ip:port format)
	ProxyType     ProxyType // Type of proxies to check
	Endpoint      string    // Endpoint to check against
	Threads       int       // Number of threads to use
	UpstreamProxy string    // Optional upstream proxy (ip:port format)
	UpstreamType  ProxyType // Type of upstream proxy
}

// ProxyResult represents the result of a proxy check (result.go)
/* type ProxyResult struct {
	Proxy      string    // Proxy address (ip:port)
	Type       ProxyType // Type of proxy
	Status     string    // Status (LIVE, DEAD)
	Latency    int64     // Latency in milliseconds
	OutgoingIP string    // Outgoing IP address
	Geo        string    // Geo location (if available)
	Error      string    // Error message (if any)
}
*/
// Stats represents the statistics of proxy checks
/* type Stats struct {
	Total        int               // Total number of proxies
	Live         int               // Number of live proxies
	Dead         int               // Number of dead proxies
	Errors       int               // Number of errors
	TypeCounts   map[ProxyType]int // Count of each proxy type
	AverageSpeed int64             // Average speed in milliseconds
} */

// Manager handles proxy checking operations
type Manager struct {
	mutex             sync.Mutex
	workingMutex      sync.Mutex
	running           bool
	paused            bool
	results           []ProxyResult
	working           []string
	stats             Stats
	stopChan          chan struct{}
	pauseChan         chan struct{}
	resumeChan        chan struct{}
	workerCount       int
	pausedWorkerCount int32
}

// NewManager creates a new proxy checker manager
/* func NewManager() *Manager {
	return &Manager{
		results: make([]ProxyResult, 0),
		working: make([]string, 0),
		stats: Stats{
			TypeCounts: make(map[ProxyType]int),
		},
	}
} */

// GetWorkerCount returns the total number of workers
func (m *Manager) GetWorkerCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.workerCount
}

// GetPausedWorkerCount returns the number of workers that have been paused
func (m *Manager) GetPausedWorkerCount() int {
	return int(atomic.LoadInt32(&m.pausedWorkerCount))
}

// IncrementPausedWorkerCount increments the paused worker count
func (m *Manager) IncrementPausedWorkerCount() {
	atomic.AddInt32(&m.pausedWorkerCount, 1)
}

// ResetPausedWorkerCount resets the paused worker count
func (m *Manager) ResetPausedWorkerCount() {
	atomic.StoreInt32(&m.pausedWorkerCount, 0)
}

// NewManager creates a new proxy checker manager
func NewManager() *Manager {
	return &Manager{
		stopChan:   make(chan struct{}),
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
		stats: Stats{
			TypeCounts: make(map[ProxyType]int),
		},
		results: make([]ProxyResult, 0),
		mutex:   sync.Mutex{},
	}
}

// Start begins checking proxies with the given request
func (m *Manager) Start(req ProxyCheckRequest, logCb func(string), updateCb func()) {
	m.mutex.Lock()
	if m.running {
		m.mutex.Unlock()
		logCb("Check already in progress")
		return
	}

	// Reset state
	m.running = true
	m.paused = false
	m.results = []ProxyResult{}
	m.working = []string{}
	m.stats = Stats{
		Total:       len(req.ProxyList),
		Pending:     len(req.ProxyList),
		TypeCounts:  make(map[ProxyType]int),
		ThreadCount: req.Threads,
	}
	m.workerCount = req.Threads
	m.stopChan = make(chan struct{})
	m.pauseChan = make(chan struct{})
	m.resumeChan = make(chan struct{})
	m.ResetPausedWorkerCount()
	m.mutex.Unlock()
	logThgreadCount := fmt.Sprintf("Total worker threads: %d", req.Threads)

	logCb(logThgreadCount)
	logCb("Starting proxy check with " + string(req.ProxyType) + " type")

	// Create work queue
	jobs := make(chan string, len(req.ProxyList))
	for _, proxy := range req.ProxyList {
		jobs <- proxy
	}
	close(jobs)

	// Create wait group for workers
	var wg sync.WaitGroup
	wg.Add(req.Threads)

	// Track total latency for average calculation
	var totalLatency int64
	var liveCount int
	var latencyMutex sync.Mutex

	// Start worker goroutines
	for i := 0; i < req.Threads; i++ {
		go func(id int) {
			defer wg.Done()

			for proxy := range jobs {
				select {
				case <-m.stopChan:
					return
				case <-m.pauseChan:
					logCb(fmt.Sprintf("Worker %d paused", id))
					select {
					case <-m.resumeChan:
						logCb(fmt.Sprintf("Worker %d resumed", id))
					case <-m.stopChan:
						return
					}
				default:
					// Check proxy
					logCb("Checking proxy: " + proxy)

					// Determine proxy type
					proxyType := req.ProxyType
					defaultTimeout := 10 * time.Second
					if proxyType == Auto {
						// Auto-detect proxy type
						detectedType, err := DetectProxyType(proxy, defaultTimeout)
						if err != nil {
							logCb("Auto-detection failed for " + proxy + ": " + err.Error())
							proxyType = HTTP
						} else {
							proxyType = detectedType
							logCb("Auto-detected " + proxy + " as " + string(proxyType))
						}
					}

					// Perform the check
					start := time.Now()
					result := ProxyResult{
						Proxy: proxy,
						Type:  proxyType,
					}

					// Check the proxy based on its type
					var err error
					var outgoingIP string

					switch proxyType {
					case HTTP:
						outgoingIP, err = CheckHTTP(proxy, req.Endpoint, defaultTimeout, req.UpstreamProxy, req.UpstreamType)
					case HTTPS:
						outgoingIP, err = CheckHTTPS(proxy, req.Endpoint, defaultTimeout, req.UpstreamProxy, req.UpstreamType)
					case SOCKS4:
						outgoingIP, err = CheckSOCKS4(proxy, req.Endpoint, defaultTimeout, req.UpstreamProxy, req.UpstreamType)
					case SOCKS5:
						outgoingIP, err = CheckSOCKS5(proxy, req.Endpoint, defaultTimeout, req.UpstreamProxy, req.UpstreamType)
					default:
						err = fmt.Errorf("unsupported proxy type: %s", proxyType)
					}

					// Calculate latency
					result.Latency = time.Since(start).Milliseconds()

					// Set result status based on check outcome
					if err != nil {
						result.Status = "DEAD"
						result.Error = err.Error()
					} else {
						result.Status = "LIVE"
						result.OutgoingIP = outgoingIP

						// Update latency stats
						latencyMutex.Lock()
						totalLatency += result.Latency
						liveCount++
						latencyMutex.Unlock()
					}

					// Update results and stats
					m.mutex.Lock()
					m.results = append(m.results, result)

					// Update stats
					if result.Status == "LIVE" {
						m.stats.Live++
						m.workingMutex.Lock()
						m.working = append(m.working, proxy)
						m.workingMutex.Unlock()
					} else if result.Status == "DEAD" {
						m.stats.Dead++
					} else {
						m.stats.Errors++
					}

					m.stats.TypeCounts[proxyType]++

					// Calculate average speed
					if liveCount > 0 {
						m.stats.AverageSpeed = totalLatency / int64(liveCount)
					}

					m.mutex.Unlock()

					// Notify UI
					updateCb()
				}
			}
		}(i)
	}

	// Wait for completion in a separate goroutine
	go func() {
		wg.Wait()
		m.mutex.Lock()
		m.running = false
		m.paused = false
		m.mutex.Unlock()
		logCb("Proxy check completed")
		updateCb()
	}()
}

// Stop stops the current check operation
func (m *Manager) Stop(force bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running {
		return
	}

	// Check if stopChan is already closed
	select {
	case <-m.stopChan:
		// Channel is already closed, create a new one for future use
		m.stopChan = make(chan struct{})
	default:
		// Channel is still open, close it to signal workers to stop
		close(m.stopChan)
	}

	m.running = false

	// For graceful stop, the running flag will be set to false when all workers finish
}

// Pause pauses the current check operation
func (m *Manager) Pause() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running || m.paused {
		return false
	}

	m.paused = true
	m.ResetPausedWorkerCount()
	close(m.pauseChan)
	return true
}

// SetWorkerCount sets the worker count
func (m *Manager) SetWorkerCount(count int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.workerCount = count
}

// Resume resumes the current check operation
func (m *Manager) Resume() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running || !m.paused {
		return false
	}

	m.paused = false
	m.pauseChan = make(chan struct{})
	close(m.resumeChan)
	m.resumeChan = make(chan struct{})
	return true
}

// IsPaused returns whether the check operation is paused
func (m *Manager) IsPaused() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.paused
}

// ForceStop immediately terminates all proxy checking operations
func (m *Manager) ForceStop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running {
		return
	}

	// Close the stop channel to signal all workers to stop
	close(m.stopChan)

	// Reset channels
	m.stopChan = make(chan struct{})
	m.pauseChan = make(chan struct{})
	m.resumeChan = make(chan struct{})

	// Reset state
	m.running = false
	m.paused = false
	atomic.StoreInt32(&m.pausedWorkerCount, 0)
}

// ForcePause immediately pauses all proxy checking operations
func (m *Manager) ForcePause() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running || m.paused {
		return false
	}

	// Set paused state immediately
	m.paused = true

	// Close the pause channel to signal all workers to pause
	close(m.pauseChan)

	// Reset the pause channel for future use
	m.pauseChan = make(chan struct{})

	// Reset the paused worker count
	atomic.StoreInt32(&m.pausedWorkerCount, int32(m.workerCount))

	return true
}

// GetResults returns the current results
func (m *Manager) GetResults() []ProxyResult {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Return a copy to avoid race conditions
	results := make([]ProxyResult, len(m.results))
	copy(results, m.results)
	return results
}

// ClearResults clears all results and resets the statistics
func (m *Manager) ClearResults() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Only allow clearing if not currently running
	if m.running && !m.paused {
		return
	}

	// Clear results and working proxies
	m.results = []ProxyResult{}
	m.working = []string{}

	// Reset statistics
	m.stats = Stats{
		TypeCounts: make(map[ProxyType]int),
	}
}

// GetWorkingProxies returns the list of working proxies
/* func (m *Manager) GetWorkingProxies() []string {
	m.workingMutex.Lock()
	defer m.workingMutex.Unlock()

	// Return a copy to avoid race conditions
	working := make([]string, len(m.working))
	copy(working, m.working)
	return working
} */

// GetStats returns the current statistics
func (m *Manager) GetStats() Stats {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Return a copy to avoid race conditions
	stats := Stats{
		Total:        m.stats.Total,
		Pending:      m.stats.Pending,
		Live:         m.stats.Live,
		Dead:         m.stats.Dead,
		Errors:       m.stats.Errors,
		AverageSpeed: m.stats.AverageSpeed,
		TypeCounts:   make(map[ProxyType]int),
	}

	for k, v := range m.stats.TypeCounts {
		stats.TypeCounts[k] = v
	}

	// Recalculate pending count to ensure accuracy
	stats.Pending = stats.Total - stats.Live - stats.Dead - stats.Errors

	return stats
}

// IsRunning returns whether a check is currently running
func (m *Manager) IsRunning() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.running
}

// DetectProxyType attempts to detect the type of a proxy
/* func DetectProxyType(proxyAddr string, timeout time.Duration) (ProxyType, error) {
	// Try SOCKS5 first
	if _, err := CheckSOCKS5(proxyAddr, "https://api.ipify.org", timeout, "", Auto); err == nil {
		return SOCKS5, nil
	}

	// Try SOCKS4
	if _, err := CheckSOCKS4(proxyAddr, "https://api.ipify.org", timeout, "", Auto); err == nil {
		return SOCKS4, nil
	}

	// Try HTTPS
	if _, err := CheckHTTPS(proxyAddr, "https://api.ipify.org", timeout, "", Auto); err == nil {
		return HTTPS, nil
	}

	// Try HTTP
	if _, err := CheckHTTP(proxyAddr, "https://api.ipify.org", timeout, "", Auto); err == nil {
		return HTTP, nil
	}

	return "", fmt.Errorf("could not detect proxy type")
} */
