package tesla

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN_123", logger)
	
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	
	if client.GetVIN() != "TEST_VIN_123" {
		t.Errorf("Expected VIN 'TEST_VIN_123', got '%s'", client.GetVIN())
	}
	
	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}
}

func TestClientMethodsWithoutConnection(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN_123", logger)
	ctx := context.Background()
	
	// Test methods that should fail when not connected
	_, err := client.GetHVACState(ctx)
	if err == nil {
		t.Error("Expected GetHVACState to fail when not connected")
	}
	
	err = client.SetTemperature(ctx, 22.0, 22.0)
	if err == nil {
		t.Error("Expected SetTemperature to fail when not connected")
	}
	
	err = client.SetClimateOn(ctx)
	if err == nil {
		t.Error("Expected SetClimateOn to fail when not connected")
	}
	
	err = client.SetClimateOff(ctx)
	if err == nil {
		t.Error("Expected SetClimateOff to fail when not connected")
	}
}

func TestHVACState(t *testing.T) {
	state := &HVACState{
		IsOn:                true,
		DriverTempCelsius:   22.0,
		PassengerTempCelsius: 22.0,
		InsideTempCelsius:   20.0,
		OutsideTempCelsius:  15.0,
		FanStatus:           2,
		IsFrontDefrosterOn:  false,
		IsRearDefrosterOn:   false,
		IsAutoConditioning:  true,
	}
	
	if !state.IsOn {
		t.Error("Expected HVAC to be on")
	}
	
	if state.DriverTempCelsius != 22.0 {
		t.Errorf("Expected driver temp 22.0, got %.1f", state.DriverTempCelsius)
	}
	
	if state.FanStatus != 2 {
		t.Errorf("Expected fan status 2, got %d", state.FanStatus)
	}
}

func TestHVACSettings(t *testing.T) {
	settings := &HVACSettings{
		DriverTempCelsius:   22.0,
		PassengerTempCelsius: 24.0,
		IsAutoMode:          true,
	}
	
	if settings.DriverTempCelsius != 22.0 {
		t.Errorf("Expected driver temp 22.0, got %.1f", settings.DriverTempCelsius)
	}
	
	if settings.PassengerTempCelsius != 24.0 {
		t.Errorf("Expected passenger temp 24.0, got %.1f", settings.PassengerTempCelsius)
	}
	
	if !settings.IsAutoMode {
		t.Error("Expected auto mode to be enabled")
	}
}

func TestSetTimeout(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN_123", logger)
	
	// Test that SetTimeout doesn't panic
	client.SetTimeout(30 * time.Second)
}
