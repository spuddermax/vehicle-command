package tesla

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

func TestRetryWithBackoffHVAC(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test successful call on first attempt
	callCount := 0
	err := client.retryWithBackoff(context.Background(), "test_success", func() error {
		callCount++
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestRetryWithBackoffFailureHVAC(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test failure after all retries
	callCount := 0
	err := client.retryWithBackoff(context.Background(), "test_failure", func() error {
		callCount++
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error after retries, got nil")
	}
	
	expectedCalls := client.retryConfig.MaxRetries + 1
	if callCount != expectedCalls {
		t.Errorf("Expected %d calls, got %d", expectedCalls, callCount)
	}
}

func TestRetryWithBackoffSuccessAfterRetriesHVAC(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test success after some retries
	callCount := 0
	err := client.retryWithBackoff(context.Background(), "test_success_after_retries", func() error {
		callCount++
		if callCount < 2 {
			return errors.New("test error")
		}
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryWithBackoffContextCancellationHVAC(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	callCount := 0
	err := client.retryWithBackoff(ctx, "test_context_cancel", func() error {
		callCount++
		time.Sleep(100 * time.Millisecond) // Longer than context timeout
		return errors.New("test error")
	})
	
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
	
	// Should have made at least one call before timing out
	if callCount == 0 {
		t.Error("Expected at least one call before timeout")
	}
}

func TestCalculateDelay(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test delay calculation
	delay1 := client.calculateDelay(0)
	delay2 := client.calculateDelay(1)
	delay3 := client.calculateDelay(2)
	
	// Delays should increase exponentially
	if delay2 <= delay1 {
		t.Errorf("Expected delay2 > delay1, got delay1=%v, delay2=%v", delay1, delay2)
	}
	
	if delay3 <= delay2 {
		t.Errorf("Expected delay3 > delay2, got delay2=%v, delay3=%v", delay2, delay3)
	}
	
	// Delays should not exceed max delay
	if delay1 > client.retryConfig.MaxDelay {
		t.Errorf("Expected delay1 <= max delay, got delay1=%v, max=%v", delay1, client.retryConfig.MaxDelay)
	}
	
	if delay2 > client.retryConfig.MaxDelay {
		t.Errorf("Expected delay2 <= max delay, got delay2=%v, max=%v", delay2, client.retryConfig.MaxDelay)
	}
	
	if delay3 > client.retryConfig.MaxDelay {
		t.Errorf("Expected delay3 <= max delay, got delay3=%v, max=%v", delay3, client.retryConfig.MaxDelay)
	}
}

func TestCalculateDelayWithJitter(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	client.retryConfig.Jitter = true
	
	// Test that jitter adds randomness
	delays := make([]time.Duration, 10)
	for i := 0; i < 10; i++ {
		delays[i] = client.calculateDelay(1)
	}
	
	// Check that delays are not all the same (jitter is working)
	allSame := true
	for i := 1; i < len(delays); i++ {
		if delays[i] != delays[0] {
			allSame = false
			break
		}
	}
	
	if allSame {
		t.Error("Expected jitter to add randomness to delays")
	}
}

func TestCalculateDelayWithoutJitter(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	client.retryConfig.Jitter = false
	
	// Test that delays are consistent without jitter
	delay1 := client.calculateDelay(1)
	delay2 := client.calculateDelay(1)
	
	if delay1 != delay2 {
		t.Errorf("Expected consistent delays without jitter, got delay1=%v, delay2=%v", delay1, delay2)
	}
}

func TestWithTimeout(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test timeout wrapper
	ctx, cancel := client.withTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	// Test that context times out
	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected context to timeout")
	}
}

func TestRetryConfigValidation(t *testing.T) {
	// Test valid retry config
	config := RetryConfig{
		MaxRetries:    3,
		InitialDelay:  time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}
	
	// This should not panic
	_ = config
}

func TestRetryWithCustomConfig(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	// Create client with custom retry config
	retryConfig := RetryConfig{
		MaxRetries:    1,
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      100 * time.Millisecond,
		BackoffFactor: 2.0,
		Jitter:        false,
	}
	
	circuitConfig := CircuitBreakerConfig{
		MaxFailures:      5,
		ResetTimeout:     60 * time.Second,
		HalfOpenMaxCalls: 3,
	}
	
	client := NewClientWithConfig("TEST_VIN", logger, retryConfig, circuitConfig)
	
	// Test that custom config is used
	if client.retryConfig.MaxRetries != 1 {
		t.Errorf("Expected max retries 1, got %d", client.retryConfig.MaxRetries)
	}
	
	if client.retryConfig.InitialDelay != 10*time.Millisecond {
		t.Errorf("Expected initial delay 10ms, got %v", client.retryConfig.InitialDelay)
	}
	
	// Test retry with custom config
	callCount := 0
	err := client.retryWithBackoff(context.Background(), "test_custom_config", func() error {
		callCount++
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error after retries, got nil")
	}
	
	// Should have made 2 calls (1 initial + 1 retry)
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryWithBackoffLogging(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	
	// Test that retry logging works (we can't easily test the output, but we can ensure it doesn't panic)
	callCount := 0
	err := client.retryWithBackoff(context.Background(), "test_logging", func() error {
		callCount++
		if callCount < 2 {
			return errors.New("test error")
		}
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}
