package tesla

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestFanSpeedConstants(t *testing.T) {
	// Test fan speed constants
	if FanSpeedOff != 0 {
		t.Errorf("Expected FanSpeedOff to be 0, got %d", FanSpeedOff)
	}
	
	if FanSpeed1 != 1 {
		t.Errorf("Expected FanSpeed1 to be 1, got %d", FanSpeed1)
	}
	
	if FanSpeed10 != 10 {
		t.Errorf("Expected FanSpeed10 to be 10, got %d", FanSpeed10)
	}
	
	if FanSpeedAuto != 11 {
		t.Errorf("Expected FanSpeedAuto to be 11, got %d", FanSpeedAuto)
	}
}

func TestAirflowPatternConstants(t *testing.T) {
	// Test airflow pattern constants
	if AirflowFace != 0 {
		t.Errorf("Expected AirflowFace to be 0, got %d", AirflowFace)
	}
	
	if AirflowFeet != 1 {
		t.Errorf("Expected AirflowFeet to be 1, got %d", AirflowFeet)
	}
	
	if AirflowDefrost != 2 {
		t.Errorf("Expected AirflowDefrost to be 2, got %d", AirflowDefrost)
	}
	
	if AirflowAuto != 7 {
		t.Errorf("Expected AirflowAuto to be 7, got %d", AirflowAuto)
	}
}

func TestDefrosterModeConstants(t *testing.T) {
	// Test defroster mode constants
	if DefrosterOff != 0 {
		t.Errorf("Expected DefrosterOff to be 0, got %d", DefrosterOff)
	}
	
	if DefrosterNormal != 1 {
		t.Errorf("Expected DefrosterNormal to be 1, got %d", DefrosterNormal)
	}
	
	if DefrosterMax != 2 {
		t.Errorf("Expected DefrosterMax to be 2, got %d", DefrosterMax)
	}
}

func TestHVACStateStructure(t *testing.T) {
	state := &HVACState{
		IsOn:                true,
		DriverTempCelsius:   22.0,
		PassengerTempCelsius: 24.0,
		InsideTempCelsius:   20.0,
		OutsideTempCelsius:  15.0,
		FanStatus:           3,
		IsFrontDefrosterOn:  false,
		IsRearDefrosterOn:   false,
		IsAutoConditioning:  true,
		MinTempCelsius:      15.0,
		MaxTempCelsius:      30.0,
		LeftTempDirection:   1,
		RightTempDirection:  1,
		IsPreconditioning:   false,
		BioweaponModeOn:     false,
	}
	
	// Test basic fields
	if !state.IsOn {
		t.Error("Expected HVAC to be on")
	}
	
	if state.DriverTempCelsius != 22.0 {
		t.Errorf("Expected driver temp 22.0, got %.1f", state.DriverTempCelsius)
	}
	
	if state.PassengerTempCelsius != 24.0 {
		t.Errorf("Expected passenger temp 24.0, got %.1f", state.PassengerTempCelsius)
	}
	
	if state.InsideTempCelsius != 20.0 {
		t.Errorf("Expected inside temp 20.0, got %.1f", state.InsideTempCelsius)
	}
	
	if state.OutsideTempCelsius != 15.0 {
		t.Errorf("Expected outside temp 15.0, got %.1f", state.OutsideTempCelsius)
	}
	
	if state.FanStatus != 3 {
		t.Errorf("Expected fan status 3, got %d", state.FanStatus)
	}
	
	if state.IsFrontDefrosterOn {
		t.Error("Expected front defroster to be off")
	}
	
	if state.IsRearDefrosterOn {
		t.Error("Expected rear defroster to be off")
	}
	
	if !state.IsAutoConditioning {
		t.Error("Expected auto conditioning to be on")
	}
	
	if state.MinTempCelsius != 15.0 {
		t.Errorf("Expected min temp 15.0, got %.1f", state.MinTempCelsius)
	}
	
	if state.MaxTempCelsius != 30.0 {
		t.Errorf("Expected max temp 30.0, got %.1f", state.MaxTempCelsius)
	}
	
	if state.LeftTempDirection != 1 {
		t.Errorf("Expected left temp direction 1, got %d", state.LeftTempDirection)
	}
	
	if state.RightTempDirection != 1 {
		t.Errorf("Expected right temp direction 1, got %d", state.RightTempDirection)
	}
	
	if state.IsPreconditioning {
		t.Error("Expected preconditioning to be off")
	}
	
	if state.BioweaponModeOn {
		t.Error("Expected bioweapon mode to be off")
	}
}

func TestHVACSettingsStructure(t *testing.T) {
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

func TestSetFanSpeedNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetFanSpeed(ctx, FanSpeed5)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestGetFanSpeedNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	_, err := client.GetFanSpeed(ctx)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetAirflowPatternNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetAirflowPattern(ctx, AirflowFace)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestGetAirflowPatternNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	_, err := client.GetAirflowPattern(ctx)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetDefrosterNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetDefroster(ctx, true, false)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetAutoModeNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetAutoMode(ctx, true)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestGetAutoModeNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	_, err := client.GetAutoMode(ctx)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetSeatHeaterNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetSeatHeater(ctx, 0, 0) // Using 0 as placeholder for vehicle.SeatPosition
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetSeatCoolerNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetSeatCooler(ctx, 0, 0) // Using 0 as placeholder for vehicle.SeatPosition
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetSteeringWheelHeaterNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetSteeringWheelHeater(ctx, true)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetPreconditioningMaxNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetPreconditioningMax(ctx, true, false)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestSetBioweaponDefenseModeNotConnected(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	client := NewClient("TEST_VIN", logger)
	ctx := context.Background()
	
	err := client.SetBioweaponDefenseMode(ctx, true, false)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}

func TestClientFromConfig(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	config := DefaultConfig()
	config.Tesla.VIN = "CONFIG_VIN_123"
	
	client := NewClientFromConfig(config, logger)
	
	if client == nil {
		t.Fatal("NewClientFromConfig returned nil")
	}
	
	if client.GetVIN() != "CONFIG_VIN_123" {
		t.Errorf("Expected VIN 'CONFIG_VIN_123', got '%s'", client.GetVIN())
	}
	
	// Test that retry config is set
	if client.retryConfig.MaxRetries != config.Retry.MaxRetries {
		t.Errorf("Expected max retries %d, got %d", config.Retry.MaxRetries, client.retryConfig.MaxRetries)
	}
}

func TestClientWithConfigManager(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/test-config.json"
	
	// Create config manager
	configManager, err := NewConfigManager(configPath)
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}
	defer configManager.Close()
	
	// Set VIN in config
	err = configManager.UpdateConfig(func(c *Config) {
		c.Tesla.VIN = "MANAGER_VIN_123"
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}
	
	// Create client with config manager
	client := NewClientWithConfigManager(configManager, logger)
	
	if client == nil {
		t.Fatal("NewClientWithConfigManager returned nil")
	}
	
	if client.GetVIN() != "MANAGER_VIN_123" {
		t.Errorf("Expected VIN 'MANAGER_VIN_123', got '%s'", client.GetVIN())
	}
}

func TestHVACErrorTypes(t *testing.T) {
	// Test error types are defined
	if ErrNotConnected == nil {
		t.Error("ErrNotConnected should not be nil")
	}
	
	if ErrConnectionLost == nil {
		t.Error("ErrConnectionLost should not be nil")
	}
	
	if ErrOperationTimeout == nil {
		t.Error("ErrOperationTimeout should not be nil")
	}
	
	if ErrRetryExhausted == nil {
		t.Error("ErrRetryExhausted should not be nil")
	}
	
	if ErrCircuitOpen == nil {
		t.Error("ErrCircuitOpen should not be nil")
	}
	
	// Test error messages
	if ErrNotConnected.Error() != "not connected to vehicle" {
		t.Errorf("Expected 'not connected to vehicle', got '%s'", ErrNotConnected.Error())
	}
	
	if ErrConnectionLost.Error() != "connection to vehicle lost" {
		t.Errorf("Expected 'connection to vehicle lost', got '%s'", ErrConnectionLost.Error())
	}
	
	if ErrOperationTimeout.Error() != "operation timed out" {
		t.Errorf("Expected 'operation timed out', got '%s'", ErrOperationTimeout.Error())
	}
	
	if ErrRetryExhausted.Error() != "retry attempts exhausted" {
		t.Errorf("Expected 'retry attempts exhausted', got '%s'", ErrRetryExhausted.Error())
	}
	
	if ErrCircuitOpen.Error() != "circuit breaker is open" {
		t.Errorf("Expected 'circuit breaker is open', got '%s'", ErrCircuitOpen.Error())
	}
}
