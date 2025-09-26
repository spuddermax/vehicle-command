package tesla

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the complete configuration for the Tesla HVAC client
type Config struct {
	// Tesla API Configuration
	Tesla TeslaConfig `json:"tesla"`

	// Client Configuration
	Client ClientConfig `json:"client"`

	// Retry Configuration
	Retry RetryConfig `json:"retry"`

	// Circuit Breaker Configuration
	CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`

	// Logging Configuration
	Logging LoggingConfig `json:"logging"`

	// File paths
	ConfigPath string `json:"-"` // Path to config file (not serialized)
}

// TeslaConfig holds Tesla-specific API configuration
type TeslaConfig struct {
	// Vehicle Information
	VIN string `json:"vin"`

	// Authentication
	PrivateKeyFile string `json:"private_key_file"`
	OAuthTokenFile string `json:"oauth_token_file"`

	// Connection Settings
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	ScanTimeout       time.Duration `json:"scan_timeout"`

	// API Settings
	MaxConcurrentRequests int `json:"max_concurrent_requests"`
	RequestTimeout        time.Duration `json:"request_timeout"`

	// Vehicle Discovery
	ScanRetries int `json:"scan_retries"`
	ScanDelay   time.Duration `json:"scan_delay"`
}

// ClientConfig holds client-specific configuration
type ClientConfig struct {
	// Client Identification
	ClientName    string `json:"client_name"`
	ClientVersion string `json:"client_version"`

	// Connection Management
	KeepAliveInterval time.Duration `json:"keep_alive_interval"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`

	// Feature Flags
	EnableAutoReconnect bool `json:"enable_auto_reconnect"`
	EnableHealthChecks  bool `json:"enable_health_checks"`
	EnableMetrics       bool `json:"enable_metrics"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`      // debug, info, warn, error
	Format     string `json:"format"`     // json, text
	Output     string `json:"output"`     // stdout, stderr, file
	FilePath   string `json:"file_path"`  // Path to log file (if output is file)
	MaxSize    int    `json:"max_size"`   // Max log file size in MB
	MaxBackups int    `json:"max_backups"` // Max number of backup files
	MaxAge     int    `json:"max_age"`    // Max age in days
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Tesla: TeslaConfig{
			ConnectionTimeout:     60 * time.Second,
			ScanTimeout:          30 * time.Second,
			MaxConcurrentRequests: 5,
			RequestTimeout:       10 * time.Second,
			ScanRetries:          3,
			ScanDelay:            2 * time.Second,
		},
		Client: ClientConfig{
			ClientName:           "tesla-hvac-client",
			ClientVersion:        "1.0.0",
			KeepAliveInterval:    30 * time.Second,
			HealthCheckInterval:  60 * time.Second,
			EnableAutoReconnect:  true,
			EnableHealthChecks:   true,
			EnableMetrics:        false,
		},
		Retry: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			Jitter:        true,
		},
		CircuitBreaker: CircuitBreakerConfig{
			MaxFailures:      5,
			ResetTimeout:     60 * time.Second,
			HalfOpenMaxCalls: 3,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
		},
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if file doesn't exist
		config := DefaultConfig()
		config.ConfigPath = configPath
		
		// Save default config to file
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config file: %w", err)
		}
		
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.ConfigPath = configPath
	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	if c.ConfigPath == "" {
		return fmt.Errorf("config path not set")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(c.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(c.ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate Tesla config
	if c.Tesla.VIN == "" {
		return fmt.Errorf("tesla.vin is required")
	}

	if c.Tesla.ConnectionTimeout <= 0 {
		return fmt.Errorf("tesla.connection_timeout must be positive")
	}

	if c.Tesla.ScanTimeout <= 0 {
		return fmt.Errorf("tesla.scan_timeout must be positive")
	}

	if c.Tesla.MaxConcurrentRequests <= 0 {
		return fmt.Errorf("tesla.max_concurrent_requests must be positive")
	}

	if c.Tesla.RequestTimeout <= 0 {
		return fmt.Errorf("tesla.request_timeout must be positive")
	}

	// Validate retry config
	if c.Retry.MaxRetries < 0 {
		return fmt.Errorf("retry.max_retries must be non-negative")
	}

	if c.Retry.InitialDelay <= 0 {
		return fmt.Errorf("retry.initial_delay must be positive")
	}

	if c.Retry.MaxDelay <= 0 {
		return fmt.Errorf("retry.max_delay must be positive")
	}

	if c.Retry.BackoffFactor <= 0 {
		return fmt.Errorf("retry.backoff_factor must be positive")
	}

	// Validate circuit breaker config
	if c.CircuitBreaker.MaxFailures <= 0 {
		return fmt.Errorf("circuit_breaker.max_failures must be positive")
	}

	if c.CircuitBreaker.ResetTimeout <= 0 {
		return fmt.Errorf("circuit_breaker.reset_timeout must be positive")
	}

	if c.CircuitBreaker.HalfOpenMaxCalls <= 0 {
		return fmt.Errorf("circuit_breaker.half_open_max_calls must be positive")
	}

	// Validate logging config
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("logging.level must be one of: debug, info, warn, error")
	}

	validFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validFormats[c.Logging.Format] {
		return fmt.Errorf("logging.format must be one of: json, text")
	}

	validOutputs := map[string]bool{
		"stdout": true, "stderr": true, "file": true,
	}
	if !validOutputs[c.Logging.Output] {
		return fmt.Errorf("logging.output must be one of: stdout, stderr, file")
	}

	if c.Logging.Output == "file" && c.Logging.FilePath == "" {
		return fmt.Errorf("logging.file_path is required when output is file")
	}

	return nil
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Tesla configuration
	if vin := os.Getenv("TESLA_VIN"); vin != "" {
		c.Tesla.VIN = vin
	}
	if keyFile := os.Getenv("TESLA_PRIVATE_KEY_FILE"); keyFile != "" {
		c.Tesla.PrivateKeyFile = keyFile
	}
	if tokenFile := os.Getenv("TESLA_OAUTH_TOKEN_FILE"); tokenFile != "" {
		c.Tesla.OAuthTokenFile = tokenFile
	}

	// Client configuration
	if name := os.Getenv("TESLA_CLIENT_NAME"); name != "" {
		c.Client.ClientName = name
	}
	if version := os.Getenv("TESLA_CLIENT_VERSION"); version != "" {
		c.Client.ClientVersion = version
	}

	// Logging configuration
	if level := os.Getenv("TESLA_LOG_LEVEL"); level != "" {
		c.Logging.Level = level
	}
	if format := os.Getenv("TESLA_LOG_FORMAT"); format != "" {
		c.Logging.Format = format
	}
	if output := os.Getenv("TESLA_LOG_OUTPUT"); output != "" {
		c.Logging.Output = output
	}
	if filePath := os.Getenv("TESLA_LOG_FILE"); filePath != "" {
		c.Logging.FilePath = filePath
	}
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "./tesla-hvac-config.json"
	}
	return filepath.Join(homeDir, ".config", "tesla-hvac", "config.json")
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	return GetConfigPath()
}
