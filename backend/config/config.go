/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 *
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/r4j3sh-com/soxyCheckerGui/backend/checker"
)

// Config represents the application configuration
type Config struct {
	// LastProxyType is the last used proxy type
	LastProxyType checker.ProxyType `json:"lastProxyType"`

	// LastEndpoint is the last used endpoint for checking proxies
	LastEndpoint string `json:"lastEndpoint"`

	// LastThreadCount is the last used thread count
	LastThreadCount int `json:"lastThreadCount"`

	// LastUpstreamProxy is the last used upstream proxy
	LastUpstreamProxy string `json:"lastUpstreamProxy"`

	// LastUpstreamProxyType is the last used upstream proxy type
	LastUpstreamProxyType checker.ProxyType `json:"lastUpstreamProxyType"`

	// DefaultEndpoints is a list of predefined endpoints for checking proxies
	DefaultEndpoints []string `json:"defaultEndpoints"`

	// MaxThreads is the maximum allowed thread count
	MaxThreads int `json:"maxThreads"`

	// Theme is the UI theme (light or dark)
	Theme string `json:"theme"`

	// EnableGeolocation enables geolocation for proxies
	EnableGeolocation bool `json:"enableGeolocation"`

	// ExportFormat is the default format for exporting proxies
	ExportFormat string `json:"exportFormat"`

	// AutoSaveResults enables automatic saving of results
	AutoSaveResults bool `json:"autoSaveResults"`

	// AutoSavePath is the path for automatically saved results
	AutoSavePath string `json:"autoSavePath"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		LastProxyType:         checker.HTTP,
		LastEndpoint:          "https://api.ipify.org",
		LastThreadCount:       20,
		LastUpstreamProxy:     "",
		LastUpstreamProxyType: checker.HTTP,
		DefaultEndpoints: []string{
			"https://api.ipify.org",
			"https://ifconfig.me/ip",
			"https://icanhazip.com",
			"https://ipinfo.io/ip",
			"https://checkip.amazonaws.com",
		},
		MaxThreads:        100,
		Theme:             "system",
		EnableGeolocation: true,
		ExportFormat:      "plain", // plain, with-type, json
		AutoSaveResults:   false,
		AutoSavePath:      "",
	}
}

var (
	instance *ConfigManager
	once     sync.Once
)

// ConfigManager handles loading and saving of application configuration
type ConfigManager struct {
	config     *Config
	configPath string
	mutex      sync.RWMutex
}

// GetInstance returns the singleton instance of ConfigManager
func GetInstance() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{
			config: DefaultConfig(),
		}
		instance.configPath = getConfigPath()

		// Update the Load call to handle the error (around line 104)
		if err := instance.Load(); err != nil {
			// Handle the error appropriately - you might want to log it or return it
			// For now, we'll just ignore it to maintain existing behavior
			_ = err
		}
	})
	return instance
}

// Load loads the configuration from disk
func (cm *ConfigManager) Load() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		dir := filepath.Dir(cm.configPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Save default config
		return cm.save()
	}

	// Read config file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	if err := json.Unmarshal(data, &cm.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// Save saves the configuration to disk
func (cm *ConfigManager) Save() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	return cm.save()
}

// save is an internal method to save the config (must be called with mutex locked)
func (cm *ConfigManager) save() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns a copy of the current configuration
func (cm *ConfigManager) GetConfig() Config {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Return a copy to avoid race conditions
	return *cm.config
}

// UpdateConfig updates the configuration
func (cm *ConfigManager) UpdateConfig(updater func(*Config)) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Apply updates
	updater(cm.config)

	// Save changes
	return cm.save()
}

// UpdateLastProxyType updates the last used proxy type
func (cm *ConfigManager) UpdateLastProxyType(proxyType checker.ProxyType) error {
	return cm.UpdateConfig(func(c *Config) {
		c.LastProxyType = proxyType
	})
}

// UpdateLastEndpoint updates the last used endpoint
func (cm *ConfigManager) UpdateLastEndpoint(endpoint string) error {
	return cm.UpdateConfig(func(c *Config) {
		c.LastEndpoint = endpoint
	})
}

// UpdateLastThreadCount updates the last used thread count
func (cm *ConfigManager) UpdateLastThreadCount(threadCount int) error {
	return cm.UpdateConfig(func(c *Config) {
		c.LastThreadCount = threadCount
	})
}

// UpdateLastUpstreamProxy updates the last used upstream proxy
func (cm *ConfigManager) UpdateLastUpstreamProxy(proxy string, proxyType checker.ProxyType) error {
	return cm.UpdateConfig(func(c *Config) {
		c.LastUpstreamProxy = proxy
		c.LastUpstreamProxyType = proxyType
	})
}

// UpdateTheme updates the UI theme
func (cm *ConfigManager) UpdateTheme(theme string) error {
	return cm.UpdateConfig(func(c *Config) {
		c.Theme = theme
	})
}

// UpdateGeolocation updates the geolocation setting
func (cm *ConfigManager) UpdateGeolocation(enable bool) error {
	return cm.UpdateConfig(func(c *Config) {
		c.EnableGeolocation = enable
	})
}

// UpdateExportFormat updates the export format
func (cm *ConfigManager) UpdateExportFormat(format string) error {
	return cm.UpdateConfig(func(c *Config) {
		c.ExportFormat = format
	})
}

// UpdateAutoSave updates the auto-save settings
func (cm *ConfigManager) UpdateAutoSave(enable bool, path string) error {
	return cm.UpdateConfig(func(c *Config) {
		c.AutoSaveResults = enable
		c.AutoSavePath = path
	})
}

// getConfigPath returns the path to the config file based on the OS
func getConfigPath() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		// On Windows, use %APPDATA%
		configDir = filepath.Join(os.Getenv("APPDATA"), "SoxyCheckerGui")
	case "darwin":
		// On macOS, use ~/Library/Application Support
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support", "SoxyCheckerGui")
	default:
		// On Linux/Unix, use ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		configDir = filepath.Join(homeDir, ".config", "SoxyCheckerGui")
	}

	return filepath.Join(configDir, "config.json")
}
