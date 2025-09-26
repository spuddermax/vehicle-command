package tesla

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/fsnotify.v1"
)

// ConfigManager manages configuration loading, saving, and hot-reloading
type ConfigManager struct {
	config     *Config
	configPath string
	watcher    *fsnotify.Watcher
	mu         sync.RWMutex
	callbacks  []ConfigChangeCallback
	stopCh     chan struct{}
}

// ConfigChangeCallback is called when configuration changes
type ConfigChangeCallback func(oldConfig, newConfig *Config) error

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) (*ConfigManager, error) {
	// Load initial configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load initial config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Load environment variables
	config.LoadFromEnv()

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	// Watch the config file directory
	configDir := getConfigDir(configPath)
	if err := watcher.Add(configDir); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch config directory: %w", err)
	}

	cm := &ConfigManager{
		config:     config,
		configPath: configPath,
		watcher:    watcher,
		callbacks:  make([]ConfigChangeCallback, 0),
		stopCh:     make(chan struct{}),
	}

	// Start watching for changes
	go cm.watchForChanges()

	return cm, nil
}

// GetConfig returns the current configuration (thread-safe)
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// UpdateConfig updates the configuration and saves it to file
func (cm *ConfigManager) UpdateConfig(updater func(*Config)) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Create a copy of the current config
	oldConfig := *cm.config

	// Apply updates
	updater(cm.config)

	// Validate the updated configuration
	if err := cm.config.Validate(); err != nil {
		// Restore old config if validation fails
		cm.config = &oldConfig
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Save to file
	if err := cm.config.Save(); err != nil {
		// Restore old config if save fails
		cm.config = &oldConfig
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Notify callbacks
	for _, callback := range cm.callbacks {
		if err := callback(&oldConfig, cm.config); err != nil {
			// Log error but don't fail the update
			fmt.Printf("Warning: config change callback failed: %v\n", err)
		}
	}

	return nil
}

// RegisterCallback registers a callback for configuration changes
func (cm *ConfigManager) RegisterCallback(callback ConfigChangeCallback) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.callbacks = append(cm.callbacks, callback)
}

// Reload reloads the configuration from file
func (cm *ConfigManager) Reload() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Load new configuration
	newConfig, err := LoadConfig(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	// Validate new configuration
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid reloaded configuration: %w", err)
	}

	// Load environment variables
	newConfig.LoadFromEnv()

	// Store old config for callbacks
	oldConfig := cm.config

	// Update configuration
	cm.config = newConfig

	// Notify callbacks
	for _, callback := range cm.callbacks {
		if err := callback(oldConfig, cm.config); err != nil {
			// Log error but don't fail the reload
			fmt.Printf("Warning: config change callback failed: %v\n", err)
		}
	}

	return nil
}

// watchForChanges watches for configuration file changes
func (cm *ConfigManager) watchForChanges() {
	for {
		select {
		case event, ok := <-cm.watcher.Events:
			if !ok {
				return
			}

			// Check if the config file was modified
			if event.Op&fsnotify.Write == fsnotify.Write && event.Name == cm.configPath {
				// Debounce rapid changes
				time.Sleep(100 * time.Millisecond)

				// Reload configuration
				if err := cm.Reload(); err != nil {
					fmt.Printf("Failed to reload configuration: %v\n", err)
				}
			}

		case err, ok := <-cm.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Configuration watcher error: %v\n", err)

		case <-cm.stopCh:
			return
		}
	}
}

// Close stops the configuration manager and cleans up resources
func (cm *ConfigManager) Close() error {
	close(cm.stopCh)
	return cm.watcher.Close()
}

// getConfigDir returns the directory containing the config file
func getConfigDir(configPath string) string {
	// Extract directory from config path
	// For now, just return the parent directory
	// In a real implementation, you might want to use filepath.Dir()
	return "."
}

// CreateDefaultConfig creates a default configuration file at the specified path
func CreateDefaultConfig(configPath string) error {
	config := DefaultConfig()
	config.ConfigPath = configPath
	return config.Save()
}

// ConfigBuilder provides a fluent interface for building configurations
type ConfigBuilder struct {
	config *Config
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithVIN sets the vehicle VIN
func (cb *ConfigBuilder) WithVIN(vin string) *ConfigBuilder {
	cb.config.Tesla.VIN = vin
	return cb
}

// WithPrivateKeyFile sets the private key file path
func (cb *ConfigBuilder) WithPrivateKeyFile(keyFile string) *ConfigBuilder {
	cb.config.Tesla.PrivateKeyFile = keyFile
	return cb
}

// WithOAuthTokenFile sets the OAuth token file path
func (cb *ConfigBuilder) WithOAuthTokenFile(tokenFile string) *ConfigBuilder {
	cb.config.Tesla.OAuthTokenFile = tokenFile
	return cb
}

// WithRetryConfig sets the retry configuration
func (cb *ConfigBuilder) WithRetryConfig(retryConfig RetryConfig) *ConfigBuilder {
	cb.config.Retry = retryConfig
	return cb
}

// WithCircuitBreakerConfig sets the circuit breaker configuration
func (cb *ConfigBuilder) WithCircuitBreakerConfig(circuitConfig CircuitBreakerConfig) *ConfigBuilder {
	cb.config.CircuitBreaker = circuitConfig
	return cb
}

// WithLoggingConfig sets the logging configuration
func (cb *ConfigBuilder) WithLoggingConfig(loggingConfig LoggingConfig) *ConfigBuilder {
	cb.config.Logging = loggingConfig
	return cb
}

// WithClientConfig sets the client configuration
func (cb *ConfigBuilder) WithClientConfig(clientConfig ClientConfig) *ConfigBuilder {
	cb.config.Client = clientConfig
	return cb
}

// Build returns the built configuration
func (cb *ConfigBuilder) Build() *Config {
	return cb.config
}

// BuildAndSave builds the configuration and saves it to file
func (cb *ConfigBuilder) BuildAndSave(configPath string) (*Config, error) {
	config := cb.Build()
	config.ConfigPath = configPath

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	if err := config.Save(); err != nil {
		return nil, fmt.Errorf("failed to save configuration: %w", err)
	}

	return config, nil
}
