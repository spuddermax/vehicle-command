package tesla

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewBLEManager(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	manager := NewBLEManager(logger)
	if manager == nil {
		t.Fatal("BLE manager is nil")
	}
	
	if manager.logger == nil {
		t.Error("Logger is nil")
	}
	
	if manager.scanTimeout != 30*time.Second {
		t.Errorf("Expected scan timeout 30s, got %v", manager.scanTimeout)
	}
	
	if manager.connTimeout != 10*time.Second {
		t.Errorf("Expected connection timeout 10s, got %v", manager.connTimeout)
	}
	
	if manager.maxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", manager.maxRetries)
	}
}

func TestBLEManagerConfiguration(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager := NewBLEManager(logger)
	
	// Test adapter ID setting
	manager.SetAdapterID("hci0")
	if manager.adapterID != "hci0" {
		t.Errorf("Expected adapter ID 'hci0', got '%s'", manager.adapterID)
	}
	
	// Test timeout setting
	manager.SetTimeouts(60*time.Second, 20*time.Second, 10*time.Second)
	if manager.scanTimeout != 60*time.Second {
		t.Errorf("Expected scan timeout 60s, got %v", manager.scanTimeout)
	}
	
	if manager.connTimeout != 20*time.Second {
		t.Errorf("Expected connection timeout 20s, got %v", manager.connTimeout)
	}
	
	if manager.sessionTimeout != 10*time.Second {
		t.Errorf("Expected session timeout 10s, got %v", manager.sessionTimeout)
	}
	
	// Test retry configuration
	manager.SetRetryConfig(5*time.Second, 5)
	if manager.retryInterval != 5*time.Second {
		t.Errorf("Expected retry interval 5s, got %v", manager.retryInterval)
	}
	
	if manager.maxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", manager.maxRetries)
	}
}

func TestConnectionState(t *testing.T) {
	// Test state string conversion
	states := []struct {
		state    ConnectionState
		expected string
	}{
		{StateDisconnected, "disconnected"},
		{StateScanning, "scanning"},
		{StateConnecting, "connecting"},
		{StateConnected, "connected"},
		{StateSessionActive, "session_active"},
		{StateError, "error"},
	}
	
	for _, s := range states {
		if s.state.String() != s.expected {
			t.Errorf("Expected state '%s', got '%s'", s.expected, s.state.String())
		}
	}
}

func TestScanResult(t *testing.T) {
	scanResult := &ScanResult{
		VIN:          "TEST_VIN_123",
		LocalName:    "Tesla Model Y",
		Address:      "AA:BB:CC:DD:EE:FF",
		RSSI:         -45,
		DiscoveredAt: time.Now(),
	}
	
	if scanResult.VIN != "TEST_VIN_123" {
		t.Errorf("Expected VIN 'TEST_VIN_123', got '%s'", scanResult.VIN)
	}
	
	if scanResult.LocalName != "Tesla Model Y" {
		t.Errorf("Expected local name 'Tesla Model Y', got '%s'", scanResult.LocalName)
	}
	
	if scanResult.Address != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("Expected address 'AA:BB:CC:DD:EE:FF', got '%s'", scanResult.Address)
	}
	
	if scanResult.RSSI != -45 {
		t.Errorf("Expected RSSI -45, got %d", scanResult.RSSI)
	}
	
	if scanResult.DiscoveredAt.IsZero() {
		t.Error("DiscoveredAt should not be zero")
	}
}

func TestBLEConnection(t *testing.T) {
	// Create a mock BLE connection
	conn := &BLEConnection{
		vin:         "TEST_VIN_123",
		state:       StateDisconnected,
		connectedAt: time.Now(),
	}
	
	if conn.GetVIN() != "TEST_VIN_123" {
		t.Errorf("Expected VIN 'TEST_VIN_123', got '%s'", conn.GetVIN())
	}
	
	if conn.GetState() != StateDisconnected {
		t.Errorf("Expected state disconnected, got %s", conn.GetState())
	}
	
	if conn.IsConnected() {
		t.Error("Expected connection to be disconnected")
	}
	
	if conn.IsSessionActive() {
		t.Error("Expected session to be inactive")
	}
	
	if conn.GetVehicle() != nil {
		t.Error("Expected vehicle to be nil")
	}
	
	if conn.GetLastError() != nil {
		t.Error("Expected no last error")
	}
	
	// Test connection duration
	duration := conn.GetConnectionDuration()
	if duration < 0 {
		t.Error("Connection duration should be non-negative")
	}
}

func TestBLEConnectionStateTransitions(t *testing.T) {
	conn := &BLEConnection{
		vin:         "TEST_VIN_123",
		state:       StateDisconnected,
		connectedAt: time.Now(),
	}
	
	// Test initial state
	if conn.GetState() != StateDisconnected {
		t.Errorf("Expected initial state disconnected, got %s", conn.GetState())
	}
	
	// Test state transitions (simulated)
	conn.mutex.Lock()
	conn.state = StateScanning
	conn.mutex.Unlock()
	
	if conn.GetState() != StateScanning {
		t.Errorf("Expected state scanning, got %s", conn.GetState())
	}
	
	conn.mutex.Lock()
	conn.state = StateConnecting
	conn.mutex.Unlock()
	
	if conn.GetState() != StateConnecting {
		t.Errorf("Expected state connecting, got %s", conn.GetState())
	}
	
	conn.mutex.Lock()
	conn.state = StateConnected
	conn.mutex.Unlock()
	
	if !conn.IsConnected() {
		t.Error("Expected connection to be active")
	}
	
	conn.mutex.Lock()
	conn.state = StateSessionActive
	conn.sessionActive = true
	conn.mutex.Unlock()
	
	if !conn.IsSessionActive() {
		t.Error("Expected session to be active")
	}
}

func TestBLEConnectionInfo(t *testing.T) {
	conn := &BLEConnection{
		vin:         "TEST_VIN_123",
		state:       StateConnected,
		connectedAt: time.Now().Add(-5 * time.Minute),
	}
	
	info := conn.GetConnectionInfo()
	
	if info["vin"] != "TEST_VIN_123" {
		t.Errorf("Expected VIN 'TEST_VIN_123', got '%s'", info["vin"])
	}
	
	if info["state"] != "connected" {
		t.Errorf("Expected state 'connected', got '%s'", info["state"])
	}
	
	if info["session_active"] != false {
		t.Error("Expected session to be inactive")
	}
	
	// Test with error
	conn.mutex.Lock()
	conn.lastError = fmt.Errorf("test error")
	conn.mutex.Unlock()
	
	info = conn.GetConnectionInfo()
	if info["last_error"] != "test error" {
		t.Errorf("Expected last error 'test error', got '%s'", info["last_error"])
	}
}

func TestRetryWithBackoff(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager := NewBLEManager(logger)
	
	// Test successful operation
	attempts := 0
	err := manager.RetryWithBackoff(context.Background(), func() error {
		attempts++
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
	
	// Test failing operation
	attempts = 0
	err = manager.RetryWithBackoff(context.Background(), func() error {
		attempts++
		return fmt.Errorf("test error")
	})
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryWithBackoffContextCancellation(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager := NewBLEManager(logger)
	
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	attempts := 0
	err := manager.RetryWithBackoff(ctx, func() error {
		attempts++
		return fmt.Errorf("test error")
	})
	
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
	
	// The first attempt will still run before checking the context
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}
