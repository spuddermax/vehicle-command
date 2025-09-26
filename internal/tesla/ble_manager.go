package tesla

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/teslamotors/vehicle-command/internal/authentication"
	"github.com/teslamotors/vehicle-command/pkg/connector/ble"
	"github.com/teslamotors/vehicle-command/pkg/vehicle"
	universal "github.com/teslamotors/vehicle-command/pkg/protocol/protobuf/universalmessage"
)

// BLEManager handles Bluetooth Low Energy connections to Tesla vehicles
type BLEManager struct {
	logger        *log.Logger
	adapterID     string
	scanTimeout   time.Duration
	connTimeout   time.Duration
	sessionTimeout time.Duration
	retryInterval  time.Duration
	maxRetries    int
}

// ConnectionState represents the current state of the BLE connection
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateScanning
	StateConnecting
	StateConnected
	StateSessionActive
	StateError
)

func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateScanning:
		return "scanning"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateSessionActive:
		return "session_active"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// BLEConnection represents an active BLE connection to a Tesla vehicle
type BLEConnection struct {
	vin            string
	conn           *ble.Connection
	vehicle        *vehicle.Vehicle
	state          ConnectionState
	lastError      error
	connectedAt    time.Time
	sessionActive  bool
	logger         *log.Logger
	sessionTimeout time.Duration
	mutex          sync.RWMutex
}

// ScanResult represents a discovered Tesla vehicle during scanning
type ScanResult struct {
	VIN         string
	LocalName   string
	Address     string
	RSSI        int16
	DiscoveredAt time.Time
}

// NewBLEManager creates a new BLE manager
func NewBLEManager(logger *log.Logger) *BLEManager {
	return &BLEManager{
		logger:         logger,
		scanTimeout:    30 * time.Second,
		connTimeout:    10 * time.Second,
		sessionTimeout: 5 * time.Second,
		retryInterval:  2 * time.Second,
		maxRetries:     3,
	}
}

// SetAdapterID sets the Bluetooth adapter ID to use
func (bm *BLEManager) SetAdapterID(adapterID string) {
	bm.adapterID = adapterID
}

// SetTimeouts sets various timeout values
func (bm *BLEManager) SetTimeouts(scanTimeout, connTimeout, sessionTimeout time.Duration) {
	bm.scanTimeout = scanTimeout
	bm.connTimeout = connTimeout
	bm.sessionTimeout = sessionTimeout
}

// SetRetryConfig sets retry configuration
func (bm *BLEManager) SetRetryConfig(retryInterval time.Duration, maxRetries int) {
	bm.retryInterval = retryInterval
	bm.maxRetries = maxRetries
}

// InitializeAdapter initializes the Bluetooth adapter
func (bm *BLEManager) InitializeAdapter() error {
	bm.logger.Println("Initializing Bluetooth adapter")
	
	err := ble.InitAdapterWithID(bm.adapterID)
	if err != nil {
		if ble.IsAdapterError(err) {
			bm.logger.Printf("Bluetooth adapter error: %s", ble.AdapterErrorHelpMessage(err))
		} else {
			bm.logger.Printf("Failed to initialize Bluetooth adapter: %s", err)
		}
		return fmt.Errorf("failed to initialize Bluetooth adapter: %w", err)
	}
	
	bm.logger.Println("Bluetooth adapter initialized successfully")
	return nil
}

// ScanForVehicle scans for a Tesla vehicle with the specified VIN
func (bm *BLEManager) ScanForVehicle(ctx context.Context, vin string) (*ScanResult, error) {
	bm.logger.Printf("Scanning for vehicle VIN: %s", vin)
	
	// Create context with timeout
	scanCtx, cancel := context.WithTimeout(ctx, bm.scanTimeout)
	defer cancel()
	
	// Scan for the vehicle
	scan, err := ble.ScanVehicleBeacon(scanCtx, vin)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for vehicle: %w", err)
	}
	
	result := &ScanResult{
		VIN:          vin,
		LocalName:    scan.LocalName,
		Address:      scan.Address,
		RSSI:         scan.RSSI,
		DiscoveredAt: time.Now(),
	}
	
	bm.logger.Printf("Found vehicle: %s (%s) %ddBm", result.LocalName, result.Address, result.RSSI)
	return result, nil
}

// ConnectToVehicle connects to a Tesla vehicle using BLE
func (bm *BLEManager) ConnectToVehicle(ctx context.Context, vin string, privateKey authentication.ECDHPrivateKey) (*BLEConnection, error) {
	bm.logger.Printf("Connecting to vehicle VIN: %s", vin)
	
	// Create context with timeout
	connCtx, cancel := context.WithTimeout(ctx, bm.connTimeout)
	defer cancel()
	
	// Scan for the vehicle first
	scanResult, err := bm.ScanForVehicle(connCtx, vin)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for vehicle: %w", err)
	}
	
	// Create BLE connection
	conn, err := ble.NewConnectionFromScanResult(connCtx, vin, &ble.ScanResult{
		LocalName: scanResult.LocalName,
		Address:   scanResult.Address,
		RSSI:      scanResult.RSSI,
		Connectable: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create BLE connection: %w", err)
	}
	
	// Create vehicle instance
	car, err := vehicle.NewVehicle(conn, privateKey, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create vehicle instance: %w", err)
	}
	
	// Connect to vehicle
	if err := car.Connect(connCtx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to connect to vehicle: %w", err)
	}
	
	// Create connection object
	bleConn := &BLEConnection{
		vin:            vin,
		conn:           conn,
		vehicle:        car,
		state:          StateConnected,
		connectedAt:    time.Now(),
		logger:         bm.logger,
		sessionTimeout: bm.sessionTimeout,
	}
	
	bm.logger.Printf("Successfully connected to vehicle: %s", vin)
	return bleConn, nil
}

// StartSession starts an authenticated session with the vehicle
func (bc *BLEConnection) StartSession(ctx context.Context) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	
	bc.logger.Printf("Starting session with vehicle: %s", bc.vin)
	
	// Create context with timeout
	sessionCtx, cancel := context.WithTimeout(ctx, bc.sessionTimeout)
	defer cancel()
	
	// Start session with all supported domains
	domains := []universal.Domain{
		universal.Domain_DOMAIN_VEHICLE_SECURITY,
		universal.Domain_DOMAIN_INFOTAINMENT,
	}
	
	err := bc.vehicle.StartSession(sessionCtx, domains)
	if err != nil {
		bc.state = StateError
		bc.lastError = err
		return fmt.Errorf("failed to start session: %w", err)
	}
	
	bc.state = StateSessionActive
	bc.sessionActive = true
	bc.lastError = nil
	
	bc.logger.Printf("Session started successfully with vehicle: %s", bc.vin)
	return nil
}

// Disconnect closes the connection to the vehicle
func (bc *BLEConnection) Disconnect() {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	
	bc.logger.Printf("Disconnecting from vehicle: %s", bc.vin)
	
	if bc.vehicle != nil {
		bc.vehicle.Disconnect()
	}
	
	if bc.conn != nil {
		bc.conn.Close()
	}
	
	bc.state = StateDisconnected
	bc.sessionActive = false
	bc.vehicle = nil
	bc.conn = nil
	
	bc.logger.Printf("Disconnected from vehicle: %s", bc.vin)
}

// GetState returns the current connection state
func (bc *BLEConnection) GetState() ConnectionState {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.state
}

// IsConnected returns true if the connection is active
func (bc *BLEConnection) IsConnected() bool {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.state == StateConnected || bc.state == StateSessionActive
}

// IsSessionActive returns true if the session is active
func (bc *BLEConnection) IsSessionActive() bool {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.state == StateSessionActive && bc.sessionActive
}

// GetVehicle returns the vehicle instance
func (bc *BLEConnection) GetVehicle() *vehicle.Vehicle {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.vehicle
}

// GetVIN returns the vehicle identification number
func (bc *BLEConnection) GetVIN() string {
	return bc.vin
}

// GetLastError returns the last error that occurred
func (bc *BLEConnection) GetLastError() error {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.lastError
}

// GetConnectionDuration returns how long the connection has been active
func (bc *BLEConnection) GetConnectionDuration() time.Duration {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return time.Since(bc.connectedAt)
}

// Reconnect attempts to reconnect to the vehicle
func (bc *BLEConnection) Reconnect(ctx context.Context, privateKey authentication.ECDHPrivateKey) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	
	bc.logger.Printf("Reconnecting to vehicle: %s", bc.vin)
	
	// Disconnect first
	if bc.vehicle != nil {
		bc.vehicle.Disconnect()
	}
	if bc.conn != nil {
		bc.conn.Close()
	}
	
	// Reset state
	bc.state = StateDisconnected
	bc.sessionActive = false
	bc.vehicle = nil
	bc.conn = nil
	
	// Create new connection
	conn, err := ble.NewConnectionFromScanResult(ctx, bc.vin, &ble.ScanResult{
		LocalName: "Tesla", // This would need to be stored from previous scan
		Address:   "",      // This would need to be stored from previous scan
		RSSI:      0,       // This would need to be stored from previous scan
		Connectable: true,
	})
	if err != nil {
		bc.state = StateError
		bc.lastError = err
		return fmt.Errorf("failed to create new connection: %w", err)
	}
	
	// Create vehicle instance
	car, err := vehicle.NewVehicle(conn, privateKey, nil)
	if err != nil {
		conn.Close()
		bc.state = StateError
		bc.lastError = err
		return fmt.Errorf("failed to create vehicle instance: %w", err)
	}
	
	// Connect to vehicle
	if err := car.Connect(ctx); err != nil {
		conn.Close()
		bc.state = StateError
		bc.lastError = err
		return fmt.Errorf("failed to connect to vehicle: %w", err)
	}
	
	// Update connection
	bc.conn = conn
	bc.vehicle = car
	bc.state = StateConnected
	bc.connectedAt = time.Now()
	bc.lastError = nil
	
	bc.logger.Printf("Reconnected to vehicle: %s", bc.vin)
	return nil
}

// RetryWithBackoff retries an operation with exponential backoff
func (bm *BLEManager) RetryWithBackoff(ctx context.Context, operation func() error) error {
	var lastErr error
	
	for attempt := 0; attempt < bm.maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(bm.retryInterval * time.Duration(attempt)):
				// Continue with retry
			}
		}
		
		bm.logger.Printf("Attempt %d/%d", attempt+1, bm.maxRetries)
		
		err := operation()
		if err == nil {
			return nil
		}
		
		lastErr = err
		bm.logger.Printf("Attempt %d failed: %v", attempt+1, err)
	}
	
	return fmt.Errorf("operation failed after %d attempts: %w", bm.maxRetries, lastErr)
}

// GetConnectionInfo returns information about the connection
func (bc *BLEConnection) GetConnectionInfo() map[string]interface{} {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	
	info := map[string]interface{}{
		"vin":            bc.vin,
		"state":          bc.state.String(),
		"connected_at":   bc.connectedAt,
		"duration":       time.Since(bc.connectedAt).String(),
		"session_active": bc.sessionActive,
	}
	
	if bc.lastError != nil {
		info["last_error"] = bc.lastError.Error()
	}
	
	return info
}
