/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 *
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

package checker

import (
	"sync"
	"time"
)

// Stats represents statistics about proxy checks
type Stats struct {
	// Total is the total number of proxies loaded
	Total int `json:"total"`

	// Live is the number of working proxies
	Live int `json:"live"`

	// Dead is the number of non-working proxies
	Dead int `json:"dead"`

	// Errors is the number of proxies that resulted in errors
	Errors int `json:"errors"`

	// Pending is the number of proxies waiting to be checked
	Pending int `json:"pending"`

	// Checking is the number of proxies currently being checked
	Checking int `json:"checking"`

	// TypeCounts is a map of proxy types to their counts
	TypeCounts map[ProxyType]int `json:"typeCounts"`

	// SuccessRate is the percentage of successful checks (live proxies)
	SuccessRate float64 `json:"successRate"`

	// AverageSpeed is the average check speed in milliseconds
	AverageSpeed int64 `json:"averageSpeed"`

	// Number of threads used for checking
	ThreadCount int `json:"threadCount"`

	// ChecksPerSecond is the number of checks completed per second
	ChecksPerSecond float64 `json:"checksPerSecond"`

	// StartTime is when the check started
	StartTime time.Time `json:"startTime"`

	// ElapsedTime is the duration since the check started
	ElapsedTime time.Duration `json:"elapsedTime"`

	// EstimatedTimeRemaining is the estimated time to complete all checks
	EstimatedTimeRemaining time.Duration `json:"estimatedTimeRemaining"`
}

// StatsTracker keeps track of proxy check statistics
type StatsTracker struct {
	stats      Stats
	mutex      sync.RWMutex
	startTime  time.Time
	totalTime  int64
	totalCount int
}

// NewStatsTracker creates a new StatsTracker
func NewStatsTracker() *StatsTracker {
	return &StatsTracker{
		stats: Stats{
			TypeCounts: make(map[ProxyType]int),
			StartTime:  time.Now(),
		},
		startTime: time.Now(),
	}
}

// Reset resets the statistics
func (st *StatsTracker) Reset(totalProxies int) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.stats = Stats{
		Total:      totalProxies,
		Pending:    totalProxies,
		TypeCounts: make(map[ProxyType]int),
		StartTime:  time.Now(),
	}

	st.startTime = time.Now()
	st.totalTime = 0
	st.totalCount = 0
}

// UpdateWithResult updates statistics based on a proxy check result
func (st *StatsTracker) UpdateWithResult(result *ProxyResult) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	// Update type counts
	if result.Type != "" {
		st.stats.TypeCounts[result.Type] = st.stats.TypeCounts[result.Type] + 1
	}

	// Update status counts
	switch result.Status {
	case StatusLive:
		st.stats.Live++
		st.stats.Pending--

		// Update speed statistics
		if result.Latency > 0 {
			st.totalTime += result.Latency
			st.totalCount++
			st.stats.AverageSpeed = st.totalTime / int64(st.totalCount)
		}

	case StatusDead:
		st.stats.Dead++
		st.stats.Pending--

	case StatusError:
		st.stats.Errors++
		st.stats.Pending--

	case StatusChecking:
		st.stats.Checking++
		st.stats.Pending--

	case StatusPending:
		// No change needed for pending status
	}

	// Calculate success rate
	completedChecks := st.stats.Live + st.stats.Dead + st.stats.Errors
	if completedChecks > 0 {
		st.stats.SuccessRate = float64(st.stats.Live) / float64(completedChecks) * 100
	}

	// Calculate elapsed time and checks per second
	st.stats.ElapsedTime = time.Since(st.startTime)
	if st.stats.ElapsedTime.Seconds() > 0 {
		st.stats.ChecksPerSecond = float64(completedChecks) / st.stats.ElapsedTime.Seconds()
	}

	// Estimate time remaining
	if st.stats.ChecksPerSecond > 0 && st.stats.Pending > 0 {
		remainingSeconds := float64(st.stats.Pending) / st.stats.ChecksPerSecond
		st.stats.EstimatedTimeRemaining = time.Duration(remainingSeconds * float64(time.Second))
	}
}

// MarkCheckingAsDead marks all checking proxies as dead
// Used when force stopping a check
func (st *StatsTracker) MarkCheckingAsDead() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.stats.Dead += st.stats.Checking
	st.stats.Checking = 0

	// Recalculate success rate
	completedChecks := st.stats.Live + st.stats.Dead + st.stats.Errors
	if completedChecks > 0 {
		st.stats.SuccessRate = float64(st.stats.Live) / float64(completedChecks) * 100
	}
}

// GetStats returns a copy of the current statistics
func (st *StatsTracker) GetStats() Stats {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	// Create a copy of the stats to avoid race conditions
	statsCopy := Stats{
		Total:                  st.stats.Total,
		Live:                   st.stats.Live,
		Dead:                   st.stats.Dead,
		Errors:                 st.stats.Errors,
		Pending:                st.stats.Pending,
		Checking:               st.stats.Checking,
		SuccessRate:            st.stats.SuccessRate,
		AverageSpeed:           st.stats.AverageSpeed,
		ChecksPerSecond:        st.stats.ChecksPerSecond,
		StartTime:              st.stats.StartTime,
		ElapsedTime:            st.stats.ElapsedTime,
		EstimatedTimeRemaining: st.stats.EstimatedTimeRemaining,
		TypeCounts:             make(map[ProxyType]int),
	}

	// Copy the type counts map
	for k, v := range st.stats.TypeCounts {
		statsCopy.TypeCounts[k] = v
	}

	return statsCopy
}

// UpdateElapsedTime updates the elapsed time and estimated time remaining
// This should be called periodically to keep time estimates accurate
func (st *StatsTracker) UpdateElapsedTime() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.stats.ElapsedTime = time.Since(st.startTime)

	completedChecks := st.stats.Live + st.stats.Dead + st.stats.Errors
	if st.stats.ElapsedTime.Seconds() > 0 {
		st.stats.ChecksPerSecond = float64(completedChecks) / st.stats.ElapsedTime.Seconds()
	}

	if st.stats.ChecksPerSecond > 0 && st.stats.Pending > 0 {
		remainingSeconds := float64(st.stats.Pending) / st.stats.ChecksPerSecond
		st.stats.EstimatedTimeRemaining = time.Duration(remainingSeconds * float64(time.Second))
	}
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	} else if d < time.Hour {
		minutes := d / time.Minute
		seconds := (d % time.Minute) / time.Second
		return minutes.String() + "m " + seconds.String() + "s"
	} else {
		hours := d / time.Hour
		minutes := (d % time.Hour) / time.Minute
		return hours.String() + "h " + minutes.String() + "m"
	}
}
