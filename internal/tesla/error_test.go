package tesla

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

func TestErrorWrapping(t *testing.T) {
	// Test that errors are properly wrapped
	originalErr := errors.New("original error")
	wrappedErr := ErrNotConnected
	
	// Test error messages
	if wrappedErr.Error() != "not connected to vehicle" {
		t.Errorf("Expected 'not connected to vehicle', got '%s'", wrappedErr.Error())
	}
	
	// Test error types
	if !errors.Is(wrappedErr, ErrNotConnected) {
		t.Error("Expected errors.Is to return true for ErrNotConnected")
	}
	
	if errors.Is(wrappedErr, originalErr) {
		t.Error("Expected errors.Is to return false for different error")
	}
}

func TestErrorTypes(t *testing.T) {
	// Test all error types are defined and have correct messages
	errorTests := []struct {
		err     error
		message string
	}{
		{ErrNotConnected, "not connected to vehicle"},
		{ErrConnectionLost, "connection to vehicle lost"},
		{ErrOperationTimeout, "operation timed out"},
		{ErrRetryExhausted, "retry attempts exhausted"},
		{ErrCircuitOpen, "circuit breaker is open"},
	}
	
	for _, test := range errorTests {
		if test.err == nil {
			t.Errorf("Expected error to not be nil")
			continue
		}
		
		if test.err.Error() != test.message {
			t.Errorf("Expected error message '%s', got '%s'", test.message, test.err.Error())
		}
	}
}

func TestErrorHandlingInRetryLogic(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test that retry logic properly handles different error types
	ctx := context.Background()
	
	// Test with ErrNotConnected
	err := client.retryWithBackoff(ctx, "test_not_connected", func() error {
		return ErrNotConnected
	})
	
	if err == nil {
		t.Error("Expected error after retries")
	}
	
	// Test with custom error
	customErr := errors.New("custom error")
	err = client.retryWithBackoff(ctx, "test_custom_error", func() error {
		return customErr
	})
	
	if err == nil {
		t.Error("Expected error after retries")
	}
	
	// Test with wrapped error
	wrappedErr := errors.New("wrapped: " + ErrNotConnected.Error())
	err = client.retryWithBackoff(ctx, "test_wrapped_error", func() error {
		return wrappedErr
	})
	
	if err == nil {
		t.Error("Expected error after retries")
	}
}

func TestErrorHandlingInCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     100 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Test that circuit breaker handles different error types
	errorTests := []error{
		ErrNotConnected,
		ErrConnectionLost,
		ErrOperationTimeout,
		errors.New("custom error"),
	}
	
	for _, testErr := range errorTests {
		err := cb.Call(func() error {
			return testErr
		})
		
		if err != testErr {
			t.Errorf("Expected error %v, got %v", testErr, err)
		}
	}
}

func TestErrorHandlingInHVACMethods(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	// Test that all HVAC methods return ErrNotConnected when not connected
	hvacTests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "GetHVACState",
			fn: func() error {
				_, err := client.GetHVACState(ctx)
				return err
			},
		},
		{
			name: "SetTemperature",
			fn: func() error {
				return client.SetTemperature(ctx, 22.0, 22.0)
			},
		},
		{
			name: "SetClimateOn",
			fn: func() error {
				return client.SetClimateOn(ctx)
			},
		},
		{
			name: "SetClimateOff",
			fn: func() error {
				return client.SetClimateOff(ctx)
			},
		},
		{
			name: "SetFanSpeed",
			fn: func() error {
				return client.SetFanSpeed(ctx, FanSpeed5)
			},
		},
		{
			name: "GetFanSpeed",
			fn: func() error {
				_, err := client.GetFanSpeed(ctx)
				return err
			},
		},
		{
			name: "SetAirflowPattern",
			fn: func() error {
				return client.SetAirflowPattern(ctx, AirflowFace)
			},
		},
		{
			name: "GetAirflowPattern",
			fn: func() error {
				_, err := client.GetAirflowPattern(ctx)
				return err
			},
		},
		{
			name: "SetDefroster",
			fn: func() error {
				return client.SetDefroster(ctx, true, false)
			},
		},
		{
			name: "SetAutoMode",
			fn: func() error {
				return client.SetAutoMode(ctx, true)
			},
		},
		{
			name: "GetAutoMode",
			fn: func() error {
				_, err := client.GetAutoMode(ctx)
				return err
			},
		},
	}
	
	for _, test := range hvacTests {
		t.Run(test.name, func(t *testing.T) {
			err := test.fn()
			if err != ErrNotConnected {
				t.Errorf("Expected ErrNotConnected, got %v", err)
			}
		})
	}
}

func TestErrorHandlingInSeatMethods(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	// Test seat heater/cooler methods
	seatTests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "SetSeatHeater",
			fn: func() error {
				return client.SetSeatHeater(ctx, 0, 0) // Using 0 as placeholder
			},
		},
		{
			name: "SetSeatCooler",
			fn: func() error {
				return client.SetSeatCooler(ctx, 0, 0) // Using 0 as placeholder
			},
		},
		{
			name: "SetSteeringWheelHeater",
			fn: func() error {
				return client.SetSteeringWheelHeater(ctx, true)
			},
		},
		{
			name: "SetPreconditioningMax",
			fn: func() error {
				return client.SetPreconditioningMax(ctx, true, false)
			},
		},
		{
			name: "SetBioweaponDefenseMode",
			fn: func() error {
				return client.SetBioweaponDefenseMode(ctx, true, false)
			},
		},
	}
	
	for _, test := range seatTests {
		t.Run(test.name, func(t *testing.T) {
			err := test.fn()
			if err != ErrNotConnected {
				t.Errorf("Expected ErrNotConnected, got %v", err)
			}
		})
	}
}

func TestErrorHandlingInConnectionMethods(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	// Test connection health check
	err := client.checkConnectionHealth(ctx)
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
	
	// Test ensure connection
	err = client.ensureConnection(ctx, "")
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestErrorHandlingWithContextCancellation(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	err := client.retryWithBackoff(ctx, "test_cancelled", func() error {
		return errors.New("test error")
	})
	
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestErrorHandlingWithContextTimeout(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	// Wait for timeout
	time.Sleep(10 * time.Millisecond)
	
	err := client.retryWithBackoff(ctx, "test_timeout", func() error {
		return errors.New("test error")
	})
	
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestErrorHandlingInCircuitBreakerStates(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Test closed state
	err := cb.Call(func() error {
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error in closed state")
	}
	
	// Test open state
	err = cb.Call(func() error {
		return nil // This should not be called
	})
	
	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
	
	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)
	
	// Test half-open state
	err = cb.Call(func() error {
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected success in half-open state, got %v", err)
	}
}

func TestErrorHandlingInConfiguration(t *testing.T) {
	// Test configuration validation errors
	config := DefaultConfig()
	config.Tesla.VIN = "" // Invalid VIN
	
	err := config.Validate()
	if err == nil {
		t.Error("Expected validation error for empty VIN")
	}
	
	// Test configuration loading errors
	_, err = LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent config file")
	}
}

func TestErrorHandlingInConfigManager(t *testing.T) {
	// Test config manager with invalid path
	_, err := NewConfigManager("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent config path")
	}
}

func TestErrorHandlingInClientCreation(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	// Test client creation with empty VIN
	client := NewClient("", logger)
	if client == nil {
		t.Fatal("NewClient should not return nil even with empty VIN")
	}
	
	// Test client creation with nil logger
	client = NewClient("TEST_VIN", nil)
	if client == nil {
		t.Fatal("NewClient should not return nil even with nil logger")
	}
}

func TestErrorHandlingInRetryExhaustion(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Set max retries to 0 to test immediate exhaustion
	client.retryConfig.MaxRetries = 0
	
	ctx := context.Background()
	err := client.retryWithBackoff(ctx, "test_exhaustion", func() error {
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error after retry exhaustion")
	}
	
	// Check that error contains retry exhaustion message
	if err.Error() != "retry attempts exhausted: test_exhaustion failed after 1 attempts: test error" {
		t.Errorf("Expected retry exhaustion error, got: %v", err)
	}
}
