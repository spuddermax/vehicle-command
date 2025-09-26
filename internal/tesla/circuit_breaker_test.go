package tesla

import (
	"errors"
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      3,
		ResetTimeout:     5 * time.Second,
		HalfOpenMaxCalls: 2,
	}
	
	cb := NewCircuitBreaker(config)
	
	if cb == nil {
		t.Fatal("NewCircuitBreaker returned nil")
	}
	
	if cb.config.MaxFailures != 3 {
		t.Errorf("Expected max failures 3, got %d", cb.config.MaxFailures)
	}
	
	if cb.config.ResetTimeout != 5*time.Second {
		t.Errorf("Expected reset timeout 5s, got %v", cb.config.ResetTimeout)
	}
	
	if cb.state != CircuitClosed {
		t.Errorf("Expected initial state CircuitClosed, got %v", cb.state)
	}
}

func TestCircuitBreakerClosedState(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      2,
		ResetTimeout:     100 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Test successful call
	callCount := 0
	err := cb.Call(func() error {
		callCount++
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
	
	if cb.state != CircuitClosed {
		t.Errorf("Expected state CircuitClosed, got %v", cb.state)
	}
	
	if cb.failureCount != 0 {
		t.Errorf("Expected failure count 0, got %d", cb.failureCount)
	}
}

func TestCircuitBreakerFailureCounting(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      2,
		ResetTimeout:     100 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Test first failure
	err := cb.Call(func() error {
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	if cb.state != CircuitClosed {
		t.Errorf("Expected state CircuitClosed, got %v", cb.state)
	}
	
	if cb.failureCount != 1 {
		t.Errorf("Expected failure count 1, got %d", cb.failureCount)
	}
	
	// Test second failure (should open circuit)
	err = cb.Call(func() error {
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	if cb.state != CircuitOpen {
		t.Errorf("Expected state CircuitOpen, got %v", cb.state)
	}
	
	if cb.failureCount != 2 {
		t.Errorf("Expected failure count 2, got %d", cb.failureCount)
	}
}

func TestCircuitBreakerOpenState(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Open the circuit
	cb.Call(func() error {
		return errors.New("test error")
	})
	
	if cb.state != CircuitOpen {
		t.Errorf("Expected state CircuitOpen, got %v", cb.state)
	}
	
	// Test that calls fail immediately in open state
	err := cb.Call(func() error {
		return nil // This should not be called
	})
	
	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreakerHalfOpenState(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Open the circuit
	cb.Call(func() error {
		return errors.New("test error")
	})
	
	if cb.state != CircuitOpen {
		t.Errorf("Expected state CircuitOpen, got %v", cb.state)
	}
	
	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)
	
	// Test successful call in half-open state
	callCount := 0
	err := cb.Call(func() error {
		callCount++
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if cb.state != CircuitClosed {
		t.Errorf("Expected state CircuitClosed, got %v", cb.state)
	}
	
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestCircuitBreakerHalfOpenFailure(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Open the circuit
	cb.Call(func() error {
		return errors.New("test error")
	})
	
	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)
	
	// Test failure in half-open state
	err := cb.Call(func() error {
		return errors.New("test error")
	})
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	if cb.state != CircuitOpen {
		t.Errorf("Expected state CircuitOpen, got %v", cb.state)
	}
}

func TestCircuitBreakerHalfOpenMaxCalls(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Open the circuit
	cb.Call(func() error {
		return errors.New("test error")
	})
	
	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)
	
	// First call in half-open state should succeed
	callCount := 0
	err := cb.Call(func() error {
		callCount++
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if cb.state != CircuitClosed {
		t.Errorf("Expected state CircuitClosed, got %v", cb.state)
	}
	
	// Second call should also succeed (circuit is now closed)
	err = cb.Call(func() error {
		callCount++
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected successful call, got error: %v", err)
	}
	
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestCircuitBreakerGetState(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Test initial state
	if cb.GetState() != CircuitClosed {
		t.Errorf("Expected initial state CircuitClosed, got %v", cb.GetState())
	}
	
	// Open the circuit
	cb.Call(func() error {
		return errors.New("test error")
	})
	
	// Test open state
	if cb.GetState() != CircuitOpen {
		t.Errorf("Expected state CircuitOpen, got %v", cb.GetState())
	}
	
	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)
	
	// Test half-open state (before any calls)
	// Note: The circuit breaker transitions to half-open when we try to call it
	// after the reset timeout, not just by waiting
	state := cb.GetState()
	if state != CircuitHalfOpen && state != CircuitOpen {
		t.Errorf("Expected state CircuitHalfOpen or CircuitOpen, got %v", state)
	}
}

func TestCircuitBreakerConcurrentAccess(t *testing.T) {
	config := CircuitBreakerConfig{
		MaxFailures:      10,
		ResetTimeout:     100 * time.Millisecond,
		HalfOpenMaxCalls: 5,
	}
	
	cb := NewCircuitBreaker(config)
	
	// Test concurrent calls
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			err := cb.Call(func() error {
				return errors.New("test error")
			})
			if err == nil {
				t.Error("Expected error, got nil")
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Circuit should be open
	if cb.GetState() != CircuitOpen {
		t.Errorf("Expected state CircuitOpen, got %v", cb.GetState())
	}
}
