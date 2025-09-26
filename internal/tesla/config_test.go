package tesla

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	
	// Test Tesla config defaults
	if config.Tesla.ConnectionTimeout != 60*time.Second {
		t.Errorf("Expected connection timeout 60s, got %v", config.Tesla.ConnectionTimeout)
	}
	
	if config.Tesla.ScanTimeout != 30*time.Second {
		t.Errorf("Expected scan timeout 30s, got %v", config.Tesla.ScanTimeout)
	}
	
	if config.Tesla.MaxConcurrentRequests != 5 {
		t.Errorf("Expected max concurrent requests 5, got %d", config.Tesla.MaxConcurrentRequests)
	}
	
	// Test retry config defaults
	if config.Retry.MaxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", config.Retry.MaxRetries)
	}
	
	if config.Retry.InitialDelay != time.Second {
		t.Errorf("Expected initial delay 1s, got %v", config.Retry.InitialDelay)
	}
	
	if config.Retry.BackoffFactor != 2.0 {
		t.Errorf("Expected backoff factor 2.0, got %.2f", config.Retry.BackoffFactor)
	}
	
	// Test circuit breaker defaults
	if config.CircuitBreaker.MaxFailures != 5 {
		t.Errorf("Expected max failures 5, got %d", config.CircuitBreaker.MaxFailures)
	}
	
	if config.CircuitBreaker.ResetTimeout != 60*time.Second {
		t.Errorf("Expected reset timeout 60s, got %v", config.CircuitBreaker.ResetTimeout)
	}
	
	// Test logging defaults
	if config.Logging.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", config.Logging.Level)
	}
	
	if config.Logging.Format != "text" {
		t.Errorf("Expected log format 'text', got '%s'", config.Logging.Format)
	}
}

func TestConfigValidation(t *testing.T) {
	config := DefaultConfig()
	config.Tesla.VIN = "TEST_VIN_123" // Set a VIN for validation
	
	// Test valid config
	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid: %v", err)
	}
	
	// Test missing VIN
	config.Tesla.VIN = ""
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for missing VIN")
	}
	
	// Test invalid connection timeout
	config = DefaultConfig()
	config.Tesla.ConnectionTimeout = -1 * time.Second
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for negative connection timeout")
	}
	
	// Test invalid retry config
	config = DefaultConfig()
	config.Retry.MaxRetries = -1
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for negative max retries")
	}
	
	// Test invalid circuit breaker config
	config = DefaultConfig()
	config.CircuitBreaker.MaxFailures = 0
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for zero max failures")
	}
	
	// Test invalid logging level
	config = DefaultConfig()
	config.Logging.Level = "invalid"
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for invalid log level")
	}
	
	// Test invalid logging format
	config = DefaultConfig()
	config.Logging.Format = "invalid"
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for invalid log format")
	}
	
	// Test invalid logging output
	config = DefaultConfig()
	config.Logging.Output = "invalid"
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for invalid log output")
	}
	
	// Test file output without file path
	config = DefaultConfig()
	config.Logging.Output = "file"
	config.Logging.FilePath = ""
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for file output without file path")
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.json")
	
	// Create config
	config := DefaultConfig()
	config.Tesla.VIN = "TEST_VIN_123"
	config.Tesla.PrivateKeyFile = "/test/private_key.pem"
	config.ConfigPath = configPath
	
	// Save config
	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
	
	// Load config
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify loaded config
	if loadedConfig.Tesla.VIN != "TEST_VIN_123" {
		t.Errorf("Expected VIN 'TEST_VIN_123', got '%s'", loadedConfig.Tesla.VIN)
	}
	
	if loadedConfig.Tesla.PrivateKeyFile != "/test/private_key.pem" {
		t.Errorf("Expected private key file '/test/private_key.pem', got '%s'", loadedConfig.Tesla.PrivateKeyFile)
	}
	
	if loadedConfig.ConfigPath != configPath {
		t.Errorf("Expected config path '%s', got '%s'", configPath, loadedConfig.ConfigPath)
	}
}

func TestConfigLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("TESLA_VIN", "ENV_VIN_123")
	os.Setenv("TESLA_PRIVATE_KEY_FILE", "/env/private_key.pem")
	os.Setenv("TESLA_CLIENT_NAME", "env-client")
	os.Setenv("TESLA_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("TESLA_VIN")
		os.Unsetenv("TESLA_PRIVATE_KEY_FILE")
		os.Unsetenv("TESLA_CLIENT_NAME")
		os.Unsetenv("TESLA_LOG_LEVEL")
	}()
	
	config := DefaultConfig()
	config.LoadFromEnv()
	
	if config.Tesla.VIN != "ENV_VIN_123" {
		t.Errorf("Expected VIN 'ENV_VIN_123', got '%s'", config.Tesla.VIN)
	}
	
	if config.Tesla.PrivateKeyFile != "/env/private_key.pem" {
		t.Errorf("Expected private key file '/env/private_key.pem', got '%s'", config.Tesla.PrivateKeyFile)
	}
	
	if config.Client.ClientName != "env-client" {
		t.Errorf("Expected client name 'env-client', got '%s'", config.Client.ClientName)
	}
	
	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.Logging.Level)
	}
}

func TestConfigBuilder(t *testing.T) {
	config := NewConfigBuilder().
		WithVIN("BUILDER_VIN_123").
		WithPrivateKeyFile("/builder/private_key.pem").
		WithRetryConfig(RetryConfig{
			MaxRetries:    5,
			InitialDelay:  2 * time.Second,
			MaxDelay:      60 * time.Second,
			BackoffFactor: 3.0,
			Jitter:        false,
		}).
		WithCircuitBreakerConfig(CircuitBreakerConfig{
			MaxFailures:      10,
			ResetTimeout:     120 * time.Second,
			HalfOpenMaxCalls: 5,
		}).
		Build()
	
	if config.Tesla.VIN != "BUILDER_VIN_123" {
		t.Errorf("Expected VIN 'BUILDER_VIN_123', got '%s'", config.Tesla.VIN)
	}
	
	if config.Tesla.PrivateKeyFile != "/builder/private_key.pem" {
		t.Errorf("Expected private key file '/builder/private_key.pem', got '%s'", config.Tesla.PrivateKeyFile)
	}
	
	if config.Retry.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", config.Retry.MaxRetries)
	}
	
	if config.Retry.InitialDelay != 2*time.Second {
		t.Errorf("Expected initial delay 2s, got %v", config.Retry.InitialDelay)
	}
	
	if config.CircuitBreaker.MaxFailures != 10 {
		t.Errorf("Expected max failures 10, got %d", config.CircuitBreaker.MaxFailures)
	}
}

func TestConfigBuilderSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "builder-config.json")
	
	config, err := NewConfigBuilder().
		WithVIN("SAVE_VIN_123").
		WithPrivateKeyFile("/save/private_key.pem").
		BuildAndSave(configPath)
	
	if err != nil {
		t.Fatalf("Failed to build and save config: %v", err)
	}
	
	if config.Tesla.VIN != "SAVE_VIN_123" {
		t.Errorf("Expected VIN 'SAVE_VIN_123', got '%s'", config.Tesla.VIN)
	}
	
	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	if path == "" {
		t.Error("Expected non-empty config path")
	}
	
	// Should contain tesla-hvac
	if !contains(path, "tesla-hvac") {
		t.Errorf("Expected config path to contain 'tesla-hvac', got '%s'", path)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "create-config.json")
	
	err := CreateDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to create default config: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
	
	// Load and verify config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load created config: %v", err)
	}
	
	if config.ConfigPath != configPath {
		t.Errorf("Expected config path '%s', got '%s'", configPath, config.ConfigPath)
	}
}

