/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 *
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

package backend

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/r4j3sh-com/soxyCheckerGui/backend/checker"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	manager    *checker.Manager
	resultsMux sync.Mutex
	results    []ProxyResult
}

// ProxyResult represents the result of a proxy check
type ProxyResult struct {
	Proxy      string  `json:"proxy"`
	Type       string  `json:"type"`
	Status     string  `json:"status"`
	Latency    float64 `json:"latency,omitempty"`
	OutgoingIP string  `json:"outgoingIp,omitempty"`
	Geo        string  `json:"geo,omitempty"`
	Error      string  `json:"error,omitempty"`
}

// Stats represents the statistics of proxy checks
type Stats struct {
	Total           int            `json:"Total"`
	Live            int            `json:"Live"`
	Dead            int            `json:"Dead"`
	Errors          int            `json:"Errors"`
	Pending         int            `json:"Pending"`
	SuccessRate     float64        `json:"SuccessRate"`
	AverageSpeed    int64          `json:"AverageSpeed"`
	ChecksPerSecond float64        `json:"ChecksPerSecond"`
	StartTime       time.Time      `json:"StartTime"`
	TypeCounts      map[string]int `json:"TypeCounts"`
}

// CheckParams represents the parameters for a proxy check
type CheckParams struct {
	ProxyList     []string `json:"ProxyList"`
	ProxyType     string   `json:"ProxyType"`
	Endpoint      string   `json:"Endpoint"`
	Threads       int      `json:"Threads"`
	UpstreamProxy string   `json:"UpstreamProxy,omitempty"`
	UpstreamType  string   `json:"UpstreamType,omitempty"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		manager: checker.NewManager(),
		results: make([]ProxyResult, 0),
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// StartCheck starts checking proxies with the given parameters
func (a *App) StartCheck(params CheckParams) string {
	// Log the start of the check
	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Starting check with %d proxies, type: %s, threads: %d",
		len(params.ProxyList), params.ProxyType, params.Threads))

	// Clear previous results
	a.resultsMux.Lock()
	a.results = make([]ProxyResult, 0, len(params.ProxyList))
	a.resultsMux.Unlock()

	// Update initial stats
	stats := Stats{
		Total:      len(params.ProxyList),
		Pending:    len(params.ProxyList),
		Live:       0,
		Dead:       0,
		Errors:     0,
		TypeCounts: make(map[string]int),
	}
	runtime.EventsEmit(a.ctx, "stats-update", stats)

	// Convert parameters to checker.ProxyCheckRequest
	checkRequest := checker.ProxyCheckRequest{
		ProxyList:     params.ProxyList,
		ProxyType:     checker.ProxyType(params.ProxyType),
		Endpoint:      params.Endpoint,
		Threads:       params.Threads,
		UpstreamProxy: params.UpstreamProxy,
		UpstreamType:  checker.ProxyType(params.UpstreamType),
	}

	// Start the check in the manager
	go a.manager.Start(checkRequest,
		// Log callback
		func(msg string) {
			runtime.EventsEmit(a.ctx, "log", msg)
		},
		// Update callback
		func() {
			a.updateResults()
			a.updateStats()
		})

	// Emit check status
	runtime.EventsEmit(a.ctx, "check-status", "running")

	return "Check started"
}

// PauseCheck pauses the current check

func (a *App) PauseCheck() string {
	fmt.Println("PauseCheck called")
	runtime.EventsEmit(a.ctx, "log", "Pausing check...")

	if a.manager == nil || !a.manager.IsRunning() {
		runtime.EventsEmit(a.ctx, "log", "No check in progress to pause")
		return "No check in progress"
	}

	/* if a.manager != nil && a.manager.IsRunning() && !a.manager.IsPaused() {
		// Use ForcePause instead of Pause for immediate effect
		a.manager.ForcePause()
		runtime.EventsEmit(a.ctx, "check-status", "paused")
		runtime.EventsEmit(a.ctx, "log", "Check paused")
	} */

	if a.manager.IsPaused() {
		runtime.EventsEmit(a.ctx, "log", "Check is already paused")
		return "Check already paused"
	}

	if a.manager.Pause() {
		// Start a goroutine to track pause progress
		go func() {
			// Wait a moment for worker count to be properly set
			time.Sleep(200 * time.Millisecond)

			totalWorkers := a.manager.GetWorkerCount()
			if totalWorkers <= 0 {
				// If no workers reported, use thread count from stats
				stats := a.manager.GetStats()
				totalWorkers = stats.ThreadCount
			}

			// Ensure we have at least one worker to avoid division by zero
			if totalWorkers <= 0 {
				totalWorkers = 1 // Prevent division by zero
			}

			runtime.EventsEmit(a.ctx, "check-status", "pausing")
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Pausing %d workers...", totalWorkers))

			// Set a timeout for the pause operation
			timeoutChan := time.After(5 * time.Second)

			// Poll until all workers are paused or timeout occurs
			maxAttempts := 300 // 30 seconds max (100ms * 300)
			for i := 0; i < maxAttempts; i++ {
				select {
				case <-timeoutChan:
					// Timeout reached, force transition to paused state
					runtime.EventsEmit(a.ctx, "check-status", "paused")
					runtime.EventsEmit(a.ctx, "log", "Pause timeout reached, forcing paused state")
					return
				default:
					pausedWorkers := a.manager.GetPausedWorkerCount()

					// Emit progress event
					runtime.EventsEmit(a.ctx, "pause-progress", map[string]interface{}{
						"paused":  pausedWorkers,
						"total":   totalWorkers,
						"percent": float64(pausedWorkers) / float64(totalWorkers) * 100,
					})

					// Check if all workers are paused
					if pausedWorkers >= totalWorkers && totalWorkers > 0 {
						runtime.EventsEmit(a.ctx, "check-status", "paused")
						runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Check paused - all %d workers stopped", pausedWorkers))
						return
					}

					// Wait before checking again
					time.Sleep(100 * time.Millisecond)
				}
			}

			// If we get here, we've exceeded maxAttempts without all workers pausing
			runtime.EventsEmit(a.ctx, "check-status", "paused")
			runtime.EventsEmit(a.ctx, "log", "Maximum pause attempts reached, forcing paused state")
		}()

		return "Check pausing"
	}

	return "Failed to pause check"
}

// ResumeCheck resumes the current paused check
func (a *App) ResumeCheck() string {
	fmt.Println("ResumeCheck called")
	runtime.EventsEmit(a.ctx, "log", "Resuming check...")

	if a.manager == nil || !a.manager.IsRunning() {
		runtime.EventsEmit(a.ctx, "log", "No check in progress to resume")
		return "No check in progress"
	}

	if !a.manager.IsPaused() {
		runtime.EventsEmit(a.ctx, "log", "Check is not paused")
		return "Check not paused"
	}

	if a.manager.Resume() {
		runtime.EventsEmit(a.ctx, "check-status", "running")
		runtime.EventsEmit(a.ctx, "log", "Check resumed")
		return "Check resumed"
	}

	return "Failed to resume check"
}

// StopCheck stops the current check gracefully
func (a *App) StopCheck() string {
	fmt.Println("StopCheck called")
	runtime.EventsEmit(a.ctx, "log", "Stopping check gracefully...")
	if a.manager != nil {
		a.manager.Stop(true)

	}
	runtime.EventsEmit(a.ctx, "check-status", "stopped")
	return "Check stopped"
}

// ForceStopCheck forces the current check to stop immediately
/* func (a *App) ForceStopCheck() string {
	fmt.Println("ForceStopCheck called")
	runtime.EventsEmit(a.ctx, "log", "Force stopping check...")
	if a.manager != nil {
		a.manager.Stop(true)
	}
	runtime.EventsEmit(a.ctx, "check-status", "stopped")
	return "Check force stopped"
} */

// ClearResults clears all results and resets the manager
func (a *App) ClearResults() string {
	fmt.Println("ClearResults called")

	// Clear the app's results
	a.resultsMux.Lock()
	a.results = []ProxyResult{}
	a.resultsMux.Unlock()

	// If there's a manager, try to clear its results too
	if a.manager != nil {
		// Check if the manager is running
		if !a.manager.IsRunning() || a.manager.IsPaused() {
			// If the manager has a ClearResults method, call it
			// Otherwise, create a new manager instance
			if clearMethod, ok := interface{}(a.manager).(interface{ ClearResults() }); ok {
				clearMethod.ClearResults()
			} else {
				// Create a new manager instance to effectively clear all results
				a.manager = checker.NewManager()
			}
		} else {
			runtime.EventsEmit(a.ctx, "log", "Cannot clear results while check is running. Stop or pause first.")
		}
	}

	// Emit events to update the UI
	runtime.EventsEmit(a.ctx, "results-update", []ProxyResult{})
	runtime.EventsEmit(a.ctx, "stats-update", Stats{
		Total:      0,
		Pending:    0,
		Live:       0,
		Dead:       0,
		Errors:     0,
		TypeCounts: make(map[string]int),
	})

	return "Results cleared"
}

// GetWorkingProxies returns a list of working proxies
func (a *App) GetWorkingProxies() []string {
	// First check if we have results in the App struct
	a.resultsMux.Lock()
	appResults := a.results
	a.resultsMux.Unlock()

	workingProxies := []string{}

	// Check results from the App struct
	for _, result := range appResults {
		status := strings.ToLower(result.Status)
		// Check if the proxy is live/working - check for multiple possible status values
		if status == "live" || status == "working" || status == "success" {
			workingProxies = append(workingProxies, result.Proxy)
		}
	}

	// If we found working proxies, return them
	if len(workingProxies) > 0 {
		//fmt.Printf("Found %d working proxies in App results\n", len(workingProxies))
		return workingProxies
	}

	// If no working proxies found in App results, check the manager's results
	if a.manager != nil {
		// Get results from the manager
		managerResults := a.manager.GetResults()
		fmt.Printf("Manager has %d total results\n", len(managerResults))

		// Check results from the manager
		for _, result := range managerResults {
			// Check if the proxy is live/working - check for multiple possible status values
			if result.Status == "live" || result.Status == "working" || result.Status == "success" {
				workingProxies = append(workingProxies, result.Proxy)
			}
		}
	}

	fmt.Printf("Total working proxies found: %d\n", len(workingProxies))
	return workingProxies
}

// updateResults gets the latest results from the manager and updates the app's results
func (a *App) updateResults() {
	managerResults := a.manager.GetResults()

	a.resultsMux.Lock()
	defer a.resultsMux.Unlock()

	// Convert checker.ProxyResult to app.ProxyResult
	a.results = make([]ProxyResult, len(managerResults))
	for i, r := range managerResults {
		a.results[i] = ProxyResult{
			Proxy:      r.Proxy,
			Type:       string(r.Type),
			Status:     string(r.Status),
			Latency:    float64(r.Latency),
			OutgoingIP: r.OutgoingIP,
			Geo:        r.Country,
			Error:      r.Error,
		}
	}

	// Emit results update
	runtime.EventsEmit(a.ctx, "results-update", a.results)
}

// updateStats updates and emits the current stats
func (a *App) updateStats() {
	managerStats := a.manager.GetStats()

	// Convert checker.Stats to app.Stats
	stats := Stats{
		Total:           managerStats.Total,
		Live:            managerStats.Live,
		Dead:            managerStats.Dead,
		Pending:         managerStats.Pending,
		Errors:          managerStats.Errors,
		SuccessRate:     managerStats.SuccessRate,
		AverageSpeed:    managerStats.AverageSpeed,
		ChecksPerSecond: managerStats.ChecksPerSecond,
		StartTime:       managerStats.StartTime,
		TypeCounts:      make(map[string]int),
	}

	// Convert type counts
	for t, count := range managerStats.TypeCounts {
		stats.TypeCounts[string(t)] = count
	}

	runtime.EventsEmit(a.ctx, "stats-update", stats)
}
