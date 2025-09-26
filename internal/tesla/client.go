package tesla

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/teslamotors/vehicle-command/internal/authentication"
	"github.com/teslamotors/vehicle-command/pkg/connector/ble"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/vehicle"
	universal "github.com/teslamotors/vehicle-command/pkg/protocol/protobuf/universalmessage"
)

// Error types for better error handling
var (
	ErrNotConnected     = errors.New("not connected to vehicle")
	ErrConnectionLost   = errors.New("connection to vehicle lost")
	ErrOperationTimeout = errors.New("operation timed out")
	ErrRetryExhausted   = errors.New("retry attempts exhausted")
	ErrCircuitOpen      = errors.New("circuit breaker is open")
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries      int           `json:"max_retries"`
	InitialDelay    time.Duration `json:"initial_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	BackoffFactor   float64       `json:"backoff_factor"`
	Jitter          bool          `json:"jitter"`
}

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	MaxFailures     int           `json:"max_failures"`
	ResetTimeout    time.Duration `json:"reset_timeout"`
	HalfOpenMaxCalls int          `json:"half_open_max_calls"`
}

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	CircuitClosed CircuitBreakerState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config        CircuitBreakerConfig
	state         CircuitBreakerState
	failureCount  int
	lastFailTime  time.Time
	successCount  int
	mutex         sync.RWMutex
}

// Client represents a Tesla vehicle client for HVAC operations
type Client struct {
	vehicle         *vehicle.Vehicle
	vin             string
	conn            *ble.Connection
	logger          *log.Logger
	retryConfig     RetryConfig
	circuitBreaker  *CircuitBreaker
	lastHealthCheck time.Time
	healthMutex     sync.RWMutex
}

// HVACState represents the current state of the vehicle's HVAC system
type HVACState struct {
	IsOn                bool    `json:"is_on"`
	DriverTempCelsius   float32 `json:"driver_temp_celsius"`
	PassengerTempCelsius float32 `json:"passenger_temp_celsius"`
	InsideTempCelsius   float32 `json:"inside_temp_celsius"`
	OutsideTempCelsius  float32 `json:"outside_temp_celsius"`
	FanStatus           int32   `json:"fan_status"`
	IsFrontDefrosterOn  bool    `json:"is_front_defroster_on"`
	IsRearDefrosterOn   bool    `json:"is_rear_defroster_on"`
	IsAutoConditioning  bool    `json:"is_auto_conditioning"`
	MinTempCelsius      float32 `json:"min_temp_celsius"`
	MaxTempCelsius      float32 `json:"max_temp_celsius"`
	LeftTempDirection   int32   `json:"left_temp_direction"`
	RightTempDirection  int32   `json:"right_temp_direction"`
	IsPreconditioning   bool    `json:"is_preconditioning"`
	BioweaponModeOn     bool    `json:"bioweapon_mode_on"`
}

// FanSpeed represents the fan speed levels
type FanSpeed int32

const (
	FanSpeedOff FanSpeed = iota
	FanSpeed1
	FanSpeed2
	FanSpeed3
	FanSpeed4
	FanSpeed5
	FanSpeed6
	FanSpeed7
	FanSpeed8
	FanSpeed9
	FanSpeed10
	FanSpeedAuto
)

// AirflowPattern represents the airflow direction patterns
type AirflowPattern int32

const (
	AirflowFace AirflowPattern = iota
	AirflowFeet
	AirflowDefrost
	AirflowFaceFeet
	AirflowFeetDefrost
	AirflowFaceDefrost
	AirflowFaceFeetDefrost
	AirflowAuto
)

// DefrosterMode represents the defroster mode
type DefrosterMode int32

const (
	DefrosterOff DefrosterMode = iota
	DefrosterNormal
	DefrosterMax
)

// HVACSettings represents HVAC control settings
type HVACSettings struct {
	DriverTempCelsius   float32 `json:"driver_temp_celsius"`
	PassengerTempCelsius float32 `json:"passenger_temp_celsius"`
	IsAutoMode          bool    `json:"is_auto_mode"`
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// Call executes a function with circuit breaker protection
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if circuit is open
	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailTime) > cb.config.ResetTimeout {
			cb.state = CircuitHalfOpen
			cb.successCount = 0
		} else {
			return ErrCircuitOpen
		}
	}

	// Check if circuit is half-open and we've exceeded max calls
	if cb.state == CircuitHalfOpen && cb.successCount >= cb.config.HalfOpenMaxCalls {
		return ErrCircuitOpen
	}

	// Execute the function
	err := fn()
	if err != nil {
		cb.failureCount++
		cb.lastFailTime = time.Now()
		
		if cb.failureCount >= cb.config.MaxFailures {
			cb.state = CircuitOpen
		}
		return err
	}

	// Success - reset failure count and update state
	cb.failureCount = 0
	cb.successCount++
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}

	return nil
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// NewClient creates a new Tesla client with default retry and circuit breaker configuration
func NewClient(vin string, logger *log.Logger) *Client {
	return &Client{
		vin:    vin,
		logger: logger,
		retryConfig: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			Jitter:        true,
		},
		circuitBreaker: NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures:     5,
			ResetTimeout:     60 * time.Second,
			HalfOpenMaxCalls: 3,
		}),
	}
}

// NewClientWithConfig creates a new Tesla client with custom configuration
func NewClientWithConfig(vin string, logger *log.Logger, retryConfig RetryConfig, circuitConfig CircuitBreakerConfig) *Client {
	return &Client{
		vin:    vin,
		logger: logger,
		retryConfig: retryConfig,
		circuitBreaker: NewCircuitBreaker(circuitConfig),
	}
}

// NewClientFromConfig creates a new Tesla client from a configuration
func NewClientFromConfig(config *Config, logger *log.Logger) *Client {
	return &Client{
		vin:    config.Tesla.VIN,
		logger: logger,
		retryConfig: config.Retry,
		circuitBreaker: NewCircuitBreaker(config.CircuitBreaker),
	}
}

// NewClientWithConfigManager creates a new Tesla client with configuration management
func NewClientWithConfigManager(configManager *ConfigManager, logger *log.Logger) *Client {
	config := configManager.GetConfig()
	client := NewClientFromConfig(config, logger)

	// Register callback to update client configuration when config changes
	configManager.RegisterCallback(func(oldConfig, newConfig *Config) error {
		client.retryConfig = newConfig.Retry
		client.circuitBreaker = NewCircuitBreaker(newConfig.CircuitBreaker)
		client.vin = newConfig.Tesla.VIN
		return nil
	})

	return client
}

// retryWithBackoff executes a function with exponential backoff retry logic
func (c *Client) retryWithBackoff(ctx context.Context, operation string, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute the function with circuit breaker protection
		err := c.circuitBreaker.Call(fn)
		if err == nil {
			if attempt > 0 {
				c.logger.Printf("Operation '%s' succeeded on attempt %d", operation, attempt+1)
			}
			return nil
		}

		lastErr = err
		
		// Don't retry on the last attempt
		if attempt == c.retryConfig.MaxRetries {
			break
		}

		// Calculate delay with exponential backoff
		delay := c.calculateDelay(attempt)
		
		c.logger.Printf("Operation '%s' failed on attempt %d: %v. Retrying in %v", 
			operation, attempt+1, err, delay)

		// Wait with context cancellation support
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("%w: %s failed after %d attempts: %v", 
		ErrRetryExhausted, operation, c.retryConfig.MaxRetries+1, lastErr)
}

// calculateDelay calculates the delay for the given attempt using exponential backoff
func (c *Client) calculateDelay(attempt int) time.Duration {
	delay := float64(c.retryConfig.InitialDelay) * math.Pow(c.retryConfig.BackoffFactor, float64(attempt))
	
	// Cap at max delay
	if delay > float64(c.retryConfig.MaxDelay) {
		delay = float64(c.retryConfig.MaxDelay)
	}
	
	// Add jitter if enabled
	if c.retryConfig.Jitter {
		// Add up to 25% jitter
		jitter := delay * 0.25 * (0.5 - math.Mod(float64(time.Now().UnixNano()), 1.0))
		delay += jitter
	}
	
	return time.Duration(delay)
}

// withTimeout wraps a context with a timeout
func (c *Client) withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// checkConnectionHealth verifies if the connection is still healthy
func (c *Client) checkConnectionHealth(ctx context.Context) error {
	c.healthMutex.Lock()
	defer c.healthMutex.Unlock()
	
	// Don't check too frequently
	if time.Since(c.lastHealthCheck) < 5*time.Second {
		return nil
	}
	
	if c.vehicle == nil || c.conn == nil {
		return ErrNotConnected
	}
	
	// Try to get a simple state to verify connection
	_, err := c.vehicle.GetState(ctx, vehicle.StateCategoryClimate)
	if err != nil {
		c.logger.Printf("Health check failed: %v", err)
		return fmt.Errorf("%w: %v", ErrConnectionLost, err)
	}
	
	c.lastHealthCheck = time.Now()
	return nil
}

// ensureConnection ensures the connection is healthy, reconnecting if necessary
func (c *Client) ensureConnection(ctx context.Context, privateKeyFile string) error {
	// Check if we have a connection
	if c.vehicle == nil || c.conn == nil {
		return ErrNotConnected
	}
	
	// Check connection health
	if err := c.checkConnectionHealth(ctx); err != nil {
		c.logger.Printf("Connection health check failed, attempting to reconnect: %v", err)
		
		// Close existing connection
		c.Disconnect()
		
		// Attempt to reconnect
		return c.Connect(ctx, privateKeyFile)
	}
	
	return nil
}

// Connect establishes a BLE connection to the Tesla vehicle with retry logic
func (c *Client) Connect(ctx context.Context, privateKeyFile string) error {
	// Add timeout to connection process
	connectCtx, cancel := c.withTimeout(ctx, 60*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(connectCtx, "connect", func() error {
		return c.connectInternal(connectCtx, privateKeyFile)
	})
}

// ConnectWithConfig establishes a BLE connection using configuration settings
func (c *Client) ConnectWithConfig(ctx context.Context, config *Config) error {
	// Use connection timeout from config
	connectCtx, cancel := c.withTimeout(ctx, config.Tesla.ConnectionTimeout)
	defer cancel()
	
	return c.retryWithBackoff(connectCtx, "connect", func() error {
		return c.connectInternalWithConfig(connectCtx, config)
	})
}

// connectInternalWithConfig performs the actual connection logic using config
func (c *Client) connectInternalWithConfig(ctx context.Context, config *Config) error {
	c.logger.Printf("Scanning for vehicle VIN: %s", config.Tesla.VIN)
	
	// Use scan timeout from config
	scanCtx, cancel := c.withTimeout(ctx, config.Tesla.ScanTimeout)
	defer cancel()
	
	// Scan for the vehicle with retries
	var scan *ble.ScanResult
	var err error
	
	for attempt := 0; attempt < config.Tesla.ScanRetries; attempt++ {
		scan, err = ble.ScanVehicleBeacon(scanCtx, config.Tesla.VIN)
		if err == nil {
			break
		}
		
		if attempt < config.Tesla.ScanRetries-1 {
			c.logger.Printf("Scan attempt %d failed: %v. Retrying in %v", 
				attempt+1, err, config.Tesla.ScanDelay)
			time.Sleep(config.Tesla.ScanDelay)
		}
	}
	
	if err != nil {
		return fmt.Errorf("failed to scan for vehicle after %d attempts: %w", 
			config.Tesla.ScanRetries, err)
	}
	
	c.logger.Printf("Found vehicle: %s (%s) %ddBm", scan.LocalName, scan.Address, scan.RSSI)
	
	// Create BLE connection
	conn, err := ble.NewConnectionFromScanResult(ctx, config.Tesla.VIN, scan)
	if err != nil {
		return fmt.Errorf("failed to create BLE connection: %w", err)
	}
	c.conn = conn
	
	// Load private key if provided
	var privateKey authentication.ECDHPrivateKey
	if config.Tesla.PrivateKeyFile != "" {
		privateKey, err = protocol.LoadPrivateKey(config.Tesla.PrivateKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
	}
	
	// Create vehicle instance
	car, err := vehicle.NewVehicle(conn, privateKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create vehicle instance: %w", err)
	}
	c.vehicle = car
	
	// Connect to vehicle
	if err := car.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to vehicle: %w", err)
	}
	
	// Start session for authenticated commands
	if err := car.StartSession(ctx, []universal.Domain{universal.Domain_DOMAIN_VEHICLE_SECURITY, universal.Domain_DOMAIN_INFOTAINMENT}); err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	
	c.logger.Println("Successfully connected to Tesla vehicle")
	return nil
}

// connectInternal performs the actual connection logic
func (c *Client) connectInternal(ctx context.Context, privateKeyFile string) error {
	c.logger.Printf("Scanning for vehicle VIN: %s", c.vin)
	
	// Scan for the vehicle
	scan, err := ble.ScanVehicleBeacon(ctx, c.vin)
	if err != nil {
		return fmt.Errorf("failed to scan for vehicle: %w", err)
	}
	
	c.logger.Printf("Found vehicle: %s (%s) %ddBm", scan.LocalName, scan.Address, scan.RSSI)
	
	// Create BLE connection
	conn, err := ble.NewConnectionFromScanResult(ctx, c.vin, scan)
	if err != nil {
		return fmt.Errorf("failed to create BLE connection: %w", err)
	}
	c.conn = conn
	
	// Load private key if provided
	var privateKey authentication.ECDHPrivateKey
	if privateKeyFile != "" {
		privateKey, err = protocol.LoadPrivateKey(privateKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
	}
	
	// Create vehicle instance
	car, err := vehicle.NewVehicle(conn, privateKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create vehicle instance: %w", err)
	}
	c.vehicle = car
	
	// Connect to vehicle
	if err := car.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to vehicle: %w", err)
	}
	
	// Start session for authenticated commands
	if err := car.StartSession(ctx, []universal.Domain{universal.Domain_DOMAIN_VEHICLE_SECURITY, universal.Domain_DOMAIN_INFOTAINMENT}); err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	
	c.logger.Println("Successfully connected to Tesla vehicle")
	return nil
}

// Disconnect closes the connection to the vehicle
func (c *Client) Disconnect() {
	if c.vehicle != nil {
		c.vehicle.Disconnect()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.logger.Println("Disconnected from Tesla vehicle")
}

// GetHVACState retrieves the current HVAC state from the vehicle with retry logic
func (c *Client) GetHVACState(ctx context.Context) (*HVACState, error) {
	// Add timeout to state retrieval
	stateCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	var result *HVACState
	err := c.retryWithBackoff(stateCtx, "get_hvac_state", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		// Request climate state
		state, err := c.vehicle.GetState(stateCtx, vehicle.StateCategoryClimate)
		if err != nil {
			return fmt.Errorf("failed to get climate state: %w", err)
		}
		
		// Parse climate state from the actual response
		climateState := state.GetClimateState()
		if climateState == nil {
			return fmt.Errorf("no climate state data received")
		}
		
		result = &HVACState{
			IsOn:                climateState.GetIsClimateOn(),
			DriverTempCelsius:   climateState.GetDriverTempSetting(),
			PassengerTempCelsius: climateState.GetPassengerTempSetting(),
			InsideTempCelsius:   climateState.GetInsideTempCelsius(),
			OutsideTempCelsius:  climateState.GetOutsideTempCelsius(),
			FanStatus:           climateState.GetFanStatus(),
			IsFrontDefrosterOn:  climateState.GetIsFrontDefrosterOn(),
			IsRearDefrosterOn:   climateState.GetIsRearDefrosterOn(),
			IsAutoConditioning:  climateState.GetIsAutoConditioningOn(),
			MinTempCelsius:      climateState.GetMinAvailTempCelsius(),
			MaxTempCelsius:      climateState.GetMaxAvailTempCelsius(),
			LeftTempDirection:   climateState.GetLeftTempDirection(),
			RightTempDirection:  climateState.GetRightTempDirection(),
			IsPreconditioning:   climateState.GetIsPreconditioning(),
			BioweaponModeOn:     climateState.GetBioweaponModeOn(),
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

// SetTemperature sets the driver and passenger temperature with retry logic
func (c *Client) SetTemperature(ctx context.Context, driverTemp, passengerTemp float32) error {
	// Add timeout to temperature setting
	tempCtx, cancel := c.withTimeout(ctx, 15*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(tempCtx, "set_temperature", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Printf("Setting temperature - Driver: %.1f°C, Passenger: %.1f°C", driverTemp, passengerTemp)
		
		return c.vehicle.ChangeClimateTemp(tempCtx, driverTemp, passengerTemp)
	})
}

// SetClimateOn turns the climate system on with retry logic
func (c *Client) SetClimateOn(ctx context.Context) error {
	// Add timeout to climate control
	climateCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(climateCtx, "set_climate_on", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Println("Turning climate system on")
		return c.vehicle.ClimateOn(climateCtx)
	})
}

// SetClimateOff turns the climate system off with retry logic
func (c *Client) SetClimateOff(ctx context.Context) error {
	// Add timeout to climate control
	climateCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(climateCtx, "set_climate_off", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Println("Turning climate system off")
		return c.vehicle.ClimateOff(climateCtx)
	})
}

// SetFanSpeed sets the fan speed level
func (c *Client) SetFanSpeed(ctx context.Context, speed FanSpeed) error {
	if c.vehicle == nil {
		return fmt.Errorf("not connected to vehicle")
	}
	
	c.logger.Printf("Setting fan speed to: %d", speed)
	
	// Convert FanSpeed to int32 for the vehicle command
	speedInt := int32(speed)
	if speed == FanSpeedAuto {
		speedInt = -1 // Use -1 for auto mode
	}
	
	// Note: The Tesla library doesn't have direct fan speed control
	// This would need to be implemented using low-level commands
	// For now, we'll log the request and return an error
	return fmt.Errorf("fan speed control not yet implemented - would set to %d", speedInt)
}

// GetFanSpeed returns the current fan speed level with retry logic
func (c *Client) GetFanSpeed(ctx context.Context) (FanSpeed, error) {
	// Add timeout to fan speed retrieval
	fanCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	var result FanSpeed
	err := c.retryWithBackoff(fanCtx, "get_fan_speed", func() error {
		state, err := c.GetHVACState(fanCtx)
		if err != nil {
			return fmt.Errorf("failed to get HVAC state: %w", err)
		}
		
		// Convert fan status to FanSpeed enum
		fanStatus := state.FanStatus
		if fanStatus == -1 {
			result = FanSpeedAuto
		} else if fanStatus == 0 {
			result = FanSpeedOff
		} else if fanStatus > 0 && fanStatus <= 10 {
			result = FanSpeed(fanStatus)
		} else {
			return fmt.Errorf("unknown fan status: %d", fanStatus)
		}
		
		return nil
	})
	
	if err != nil {
		return FanSpeedOff, err
	}
	
	return result, nil
}

// SetAirflowPattern sets the airflow direction pattern
func (c *Client) SetAirflowPattern(ctx context.Context, pattern AirflowPattern) error {
	if c.vehicle == nil {
		return fmt.Errorf("not connected to vehicle")
	}
	
	c.logger.Printf("Setting airflow pattern to: %d", pattern)
	
	// Note: The Tesla library doesn't have direct airflow pattern control
	// This would need to be implemented using low-level commands
	// For now, we'll log the request and return an error
	return fmt.Errorf("airflow pattern control not yet implemented - would set to %d", pattern)
}

// GetAirflowPattern returns the current airflow pattern with retry logic
func (c *Client) GetAirflowPattern(ctx context.Context) (AirflowPattern, error) {
	// Add timeout to airflow pattern retrieval
	airflowCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	var result AirflowPattern
	err := c.retryWithBackoff(airflowCtx, "get_airflow_pattern", func() error {
		state, err := c.GetHVACState(airflowCtx)
		if err != nil {
			return fmt.Errorf("failed to get HVAC state: %w", err)
		}
		
		// Determine airflow pattern based on defroster and temperature direction settings
		if state.IsFrontDefrosterOn {
			if state.IsRearDefrosterOn {
				result = AirflowFaceFeetDefrost
			} else {
				result = AirflowDefrost
			}
		} else {
			// Use temperature direction to determine pattern
			leftDir := state.LeftTempDirection
			rightDir := state.RightTempDirection
			
			// This is a simplified mapping - actual implementation would need more logic
			if leftDir == 1 && rightDir == 1 {
				result = AirflowFace
			} else if leftDir == 2 && rightDir == 2 {
				result = AirflowFeet
			} else if leftDir == 3 && rightDir == 3 {
				result = AirflowFaceFeet
			} else {
				result = AirflowAuto
			}
		}
		
		return nil
	})
	
	if err != nil {
		return AirflowAuto, err
	}
	
	return result, nil
}

// SetDefroster sets the front and rear defroster state
func (c *Client) SetDefroster(ctx context.Context, front, rear bool) error {
	if c.vehicle == nil {
		return fmt.Errorf("not connected to vehicle")
	}
	
	c.logger.Printf("Setting defroster - Front: %v, Rear: %v", front, rear)
	
	// Note: The Tesla library doesn't have direct defroster control methods
	// This would need to be implemented using the low-level Send method
	// For now, we'll just log the request
	return fmt.Errorf("defroster control not yet implemented")
}

// SetAutoMode sets the auto conditioning mode
func (c *Client) SetAutoMode(ctx context.Context, enabled bool) error {
	if c.vehicle == nil {
		return fmt.Errorf("not connected to vehicle")
	}
	
	c.logger.Printf("Setting auto mode to: %v", enabled)
	
	// Use the existing ClimateOn/ClimateOff methods for now
	// A proper auto mode toggle would need to be implemented using low-level commands
	if enabled {
		return c.vehicle.ClimateOn(ctx)
	} else {
		return c.vehicle.ClimateOff(ctx)
	}
}

// GetAutoMode returns the current auto conditioning mode with retry logic
func (c *Client) GetAutoMode(ctx context.Context) (bool, error) {
	// Add timeout to auto mode retrieval
	autoCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	var result bool
	err := c.retryWithBackoff(autoCtx, "get_auto_mode", func() error {
		state, err := c.GetHVACState(autoCtx)
		if err != nil {
			return fmt.Errorf("failed to get HVAC state: %w", err)
		}
		
		result = state.IsAutoConditioning
		return nil
	})
	
	if err != nil {
		return false, err
	}
	
	return result, nil
}

// SetSeatHeater sets the seat heater level for the specified seat with retry logic
func (c *Client) SetSeatHeater(ctx context.Context, seat vehicle.SeatPosition, level vehicle.Level) error {
	// Add timeout to seat heater control
	heaterCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(heaterCtx, "set_seat_heater", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Printf("Setting seat heater - Seat: %v, Level: %v", seat, level)
		
		// Use the existing SetSeatHeater method from the vehicle library
		levels := map[vehicle.SeatPosition]vehicle.Level{seat: level}
		return c.vehicle.SetSeatHeater(heaterCtx, levels)
	})
}

// SetSeatCooler sets the seat cooler level for the specified seat with retry logic
func (c *Client) SetSeatCooler(ctx context.Context, seat vehicle.SeatPosition, level vehicle.Level) error {
	// Add timeout to seat cooler control
	coolerCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(coolerCtx, "set_seat_cooler", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Printf("Setting seat cooler - Seat: %v, Level: %v", seat, level)
		
		// Use the existing SetSeatCooler method from the vehicle library
		return c.vehicle.SetSeatCooler(coolerCtx, level, seat)
	})
}

// SetSteeringWheelHeater sets the steering wheel heater state with retry logic
func (c *Client) SetSteeringWheelHeater(ctx context.Context, enabled bool) error {
	// Add timeout to steering wheel heater control
	steeringCtx, cancel := c.withTimeout(ctx, 10*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(steeringCtx, "set_steering_wheel_heater", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Printf("Setting steering wheel heater to: %v", enabled)
		
		// Use the existing SetSteeringWheelHeater method from the vehicle library
		return c.vehicle.SetSteeringWheelHeater(steeringCtx, enabled)
	})
}

// SetPreconditioningMax sets the preconditioning max mode with retry logic
func (c *Client) SetPreconditioningMax(ctx context.Context, enabled bool, manualOverride bool) error {
	// Add timeout to preconditioning control
	precondCtx, cancel := c.withTimeout(ctx, 15*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(precondCtx, "set_preconditioning_max", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Printf("Setting preconditioning max - Enabled: %v, Manual Override: %v", enabled, manualOverride)
		
		// Use the existing SetPreconditioningMax method from the vehicle library
		return c.vehicle.SetPreconditioningMax(precondCtx, enabled, manualOverride)
	})
}

// SetBioweaponDefenseMode sets the bioweapon defense mode with retry logic
func (c *Client) SetBioweaponDefenseMode(ctx context.Context, enabled bool, manualOverride bool) error {
	// Add timeout to bioweapon defense control
	bioCtx, cancel := c.withTimeout(ctx, 15*time.Second)
	defer cancel()
	
	return c.retryWithBackoff(bioCtx, "set_bioweapon_defense_mode", func() error {
		if c.vehicle == nil {
			return ErrNotConnected
		}
		
		c.logger.Printf("Setting bioweapon defense mode - Enabled: %v, Manual Override: %v", enabled, manualOverride)
		
		// Use the existing SetBioweaponDefenseMode method from the vehicle library
		return c.vehicle.SetBioweaponDefenseMode(bioCtx, enabled, manualOverride)
	})
}

// IsConnected returns true if the client is connected to a vehicle
func (c *Client) IsConnected() bool {
	return c.vehicle != nil && c.conn != nil
}

// GetVIN returns the vehicle identification number
func (c *Client) GetVIN() string {
	return c.vin
}

// SetTimeout sets the timeout for vehicle operations
func (c *Client) SetTimeout(timeout time.Duration) {
	// This would be implemented by setting context timeouts
	// For now, it's a placeholder
}
