package tesla

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestIsConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Initially not connected
	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}
	
	// Test with mock connection (we can't easily test real connection in unit tests)
	// This would require mocking the vehicle and BLE connection
}

func TestGetVIN(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN_123", logger)
	
	if client.GetVIN() != "TEST_VIN_123" {
		t.Errorf("Expected VIN 'TEST_VIN_123', got '%s'", client.GetVIN())
	}
}

func TestSetTimeoutConnection(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test that SetTimeout doesn't panic
	client.SetTimeout(30 * time.Second)
}

func TestDisconnect(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test that Disconnect doesn't panic when not connected
	client.Disconnect()
	
	// Test that Disconnect doesn't panic when called multiple times
	client.Disconnect()
}

func TestConnectWithInvalidVIN(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("", logger) // Empty VIN
	
	ctx := context.Background()
	err := client.Connect(ctx, "")
	
	// This should fail because we can't scan for an empty VIN
	if err == nil {
		t.Error("Expected error with empty VIN")
	}
}

func TestConnectWithContextCancellation(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	// Wait for context to timeout
	time.Sleep(10 * time.Millisecond)
	
	err := client.Connect(ctx, "")
	
	// Should fail due to context cancellation
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestConnectWithConfig(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	config := DefaultConfig()
	config.Tesla.VIN = "CONFIG_VIN_123"
	config.Tesla.ConnectionTimeout = 1 * time.Millisecond
	config.Tesla.ScanTimeout = 1 * time.Millisecond
	config.Tesla.ScanRetries = 1
	
	ctx := context.Background()
	err := client.ConnectWithConfig(ctx, config)
	
	// This should fail because we can't actually connect to a real vehicle in unit tests
	if err == nil {
		t.Error("Expected error when connecting with config")
	}
}

func TestCheckConnectionHealth(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	ctx := context.Background()
	err := client.checkConnectionHealth(ctx)
	
	// Should fail when not connected
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestEnsureConnection(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	ctx := context.Background()
	err := client.ensureConnection(ctx, "")
	
	// Should fail when not connected
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestClientWithCustomRetryConfig(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	retryConfig := RetryConfig{
		MaxRetries:    5,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      1 * time.Second,
		BackoffFactor: 1.5,
		Jitter:        false,
	}
	
	circuitConfig := CircuitBreakerConfig{
		MaxFailures:      3,
		ResetTimeout:     30 * time.Second,
		HalfOpenMaxCalls: 2,
	}
	
	client := NewClientWithConfig("TEST_VIN", logger, retryConfig, circuitConfig)
	
	if client == nil {
		t.Fatal("NewClientWithConfig returned nil")
	}
	
	if client.GetVIN() != "TEST_VIN" {
		t.Errorf("Expected VIN 'TEST_VIN', got '%s'", client.GetVIN())
	}
	
	// Test that custom retry config is used
	if client.retryConfig.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", client.retryConfig.MaxRetries)
	}
	
	if client.retryConfig.InitialDelay != 100*time.Millisecond {
		t.Errorf("Expected initial delay 100ms, got %v", client.retryConfig.InitialDelay)
	}
	
	if client.retryConfig.MaxDelay != 1*time.Second {
		t.Errorf("Expected max delay 1s, got %v", client.retryConfig.MaxDelay)
	}
	
	if client.retryConfig.BackoffFactor != 1.5 {
		t.Errorf("Expected backoff factor 1.5, got %.2f", client.retryConfig.BackoffFactor)
	}
	
	if client.retryConfig.Jitter {
		t.Error("Expected jitter to be false")
	}
	
	// Test that custom circuit breaker config is used
	if client.circuitBreaker.config.MaxFailures != 3 {
		t.Errorf("Expected max failures 3, got %d", client.circuitBreaker.config.MaxFailures)
	}
	
	if client.circuitBreaker.config.ResetTimeout != 30*time.Second {
		t.Errorf("Expected reset timeout 30s, got %v", client.circuitBreaker.config.ResetTimeout)
	}
	
	if client.circuitBreaker.config.HalfOpenMaxCalls != 2 {
		t.Errorf("Expected half open max calls 2, got %d", client.circuitBreaker.config.HalfOpenMaxCalls)
	}
}

func TestClientConnectionState(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test initial state
	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}
	
	// Test VIN
	if client.GetVIN() != "TEST_VIN" {
		t.Errorf("Expected VIN 'TEST_VIN', got '%s'", client.GetVIN())
	}
	
	// Test that we can call methods without panicking
	ctx := context.Background()
	
	// These should all fail with ErrNotConnected or retry exhaustion
	_, err := client.GetHVACState(ctx)
	if err != ErrNotConnected && err.Error() != "retry attempts exhausted: get_hvac_state failed after 4 attempts: not connected to vehicle" {
		t.Errorf("Expected ErrNotConnected or retry exhaustion, got %v", err)
	}
	
	err = client.SetTemperature(ctx, 22.0, 22.0)
	if err != ErrNotConnected && err.Error() != "retry attempts exhausted: set_temperature failed after 4 attempts: circuit breaker is open" {
		t.Errorf("Expected ErrNotConnected or retry exhaustion, got %v", err)
	}
	
	err = client.SetClimateOn(ctx)
	if err != ErrNotConnected && err.Error() != "retry attempts exhausted: set_climate_on failed after 4 attempts: circuit breaker is open" {
		t.Errorf("Expected ErrNotConnected or retry exhaustion, got %v", err)
	}
	
	err = client.SetClimateOff(ctx)
	if err != ErrNotConnected && err.Error() != "retry attempts exhausted: set_climate_off failed after 4 attempts: circuit breaker is open" {
		t.Errorf("Expected ErrNotConnected or retry exhaustion, got %v", err)
	}
}

func TestClientWithNilLogger(t *testing.T) {
	// Test that client can be created with nil logger
	client := NewClient("TEST_VIN", nil)
	
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	
	if client.GetVIN() != "TEST_VIN" {
		t.Errorf("Expected VIN 'TEST_VIN', got '%s'", client.GetVIN())
	}
	
	// Test that methods don't panic with nil logger
	ctx := context.Background()
	_, err := client.GetHVACState(ctx)
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestClientRetryConfigDefaults(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test default retry config
	if client.retryConfig.MaxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", client.retryConfig.MaxRetries)
	}
	
	if client.retryConfig.InitialDelay != time.Second {
		t.Errorf("Expected initial delay 1s, got %v", client.retryConfig.InitialDelay)
	}
	
	if client.retryConfig.MaxDelay != 30*time.Second {
		t.Errorf("Expected max delay 30s, got %v", client.retryConfig.MaxDelay)
	}
	
	if client.retryConfig.BackoffFactor != 2.0 {
		t.Errorf("Expected backoff factor 2.0, got %.2f", client.retryConfig.BackoffFactor)
	}
	
	if !client.retryConfig.Jitter {
		t.Error("Expected jitter to be true")
	}
}

func TestClientCircuitBreakerDefaults(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test default circuit breaker config
	if client.circuitBreaker.config.MaxFailures != 5 {
		t.Errorf("Expected max failures 5, got %d", client.circuitBreaker.config.MaxFailures)
	}
	
	if client.circuitBreaker.config.ResetTimeout != 60*time.Second {
		t.Errorf("Expected reset timeout 60s, got %v", client.circuitBreaker.config.ResetTimeout)
	}
	
	if client.circuitBreaker.config.HalfOpenMaxCalls != 3 {
		t.Errorf("Expected half open max calls 3, got %d", client.circuitBreaker.config.HalfOpenMaxCalls)
	}
	
	// Test initial circuit breaker state
	if client.circuitBreaker.GetState() != CircuitClosed {
		t.Errorf("Expected initial state CircuitClosed, got %v", client.circuitBreaker.GetState())
	}
}
