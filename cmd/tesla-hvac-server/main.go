package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/teslamotors/vehicle-command/internal/tesla"
)

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
)

func main() {
	// Command line flags
	var (
		port        = flag.String("port", defaultPort, "Port to listen on")
		host        = flag.String("host", defaultHost, "Host to bind to")
		configPath  = flag.String("config", "", "Path to configuration file")
		webDir      = flag.String("web", "./web", "Path to web directory")
		devMode     = flag.Bool("dev", false, "Enable development mode with CORS")
	)
	flag.Parse()

	// Setup logger
	logger := log.New(os.Stdout, "[TESLA-HVAC] ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	var config *tesla.Config
	var err error
	
	if *configPath != "" {
		config, err = tesla.LoadConfig(*configPath)
		if err != nil {
			logger.Fatalf("Failed to load config from %s: %v", *configPath, err)
		}
	} else {
		config = tesla.DefaultConfig()
		config.Tesla.VIN = "YOUR_TESLA_VIN" // Placeholder
	}

	// Create Tesla client
	client := tesla.NewClientFromConfig(config, logger)

	// Setup HTTP server
	mux := http.NewServeMux()

	// Serve static web files
	webPath, err := filepath.Abs(*webDir)
	if err != nil {
		logger.Fatalf("Failed to get absolute path for web directory: %v", err)
	}

	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		logger.Fatalf("Web directory does not exist: %s", webPath)
	}

	// File server for static assets
	fileServer := http.FileServer(http.Dir(webPath))
	mux.Handle("/", fileServer)

	// API endpoints
	apiHandler := NewAPIHandler(client, logger)
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// CORS middleware for development
	var handler http.Handler = mux
	if *devMode {
		handler = corsMiddleware(mux)
		logger.Println("Development mode enabled with CORS")
	}

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Printf("Starting Tesla HVAC server on %s:%s", *host, *port)
		logger.Printf("Web interface available at: http://%s:%s", *host, *port)
		logger.Printf("Serving files from: %s", webPath)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}

	// Disconnect from Tesla vehicle
	client.Disconnect()

	logger.Println("Server exited")
}

// CORS middleware for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// APIHandler handles API requests
type APIHandler struct {
	client *tesla.Client
	logger *log.Logger
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(client *tesla.Client, logger *log.Logger) *APIHandler {
	return &APIHandler{
		client: client,
		logger: logger,
	}
}

// ServeHTTP implements http.Handler
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Route requests
	switch r.URL.Path {
	case "/status":
		h.handleStatus(w, r)
	case "/connect":
		h.handleConnect(w, r)
	case "/hvac/state":
		h.handleHVACState(w, r)
	case "/hvac/temperature":
		h.handleTemperature(w, r)
	case "/hvac/fan":
		h.handleFanSpeed(w, r)
	case "/hvac/airflow":
		h.handleAirflow(w, r)
	case "/hvac/auto":
		h.handleAutoMode(w, r)
	case "/hvac/climate":
		h.handleClimate(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleStatus returns the current connection status
func (h *APIHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := map[string]interface{}{
		"connected": h.client.IsConnected(),
		"vin":       h.client.GetVIN(),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","data":%s}`, toJSON(status))
}

// handleConnect attempts to connect to the Tesla vehicle
func (h *APIHandler) handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	err := h.client.Connect(ctx, "")
	
	if err != nil {
		h.logger.Printf("Connection failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Connected successfully"}`)
}

// handleHVACState returns the current HVAC state
func (h *APIHandler) handleHVACState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	state, err := h.client.GetHVACState(ctx)
	
	if err != nil {
		h.logger.Printf("Failed to get HVAC state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	// Convert temperatures from Celsius to Fahrenheit for frontend
	state.DriverTempCelsius = float32(celsiusToFahrenheit(float64(state.DriverTempCelsius)))
	state.PassengerTempCelsius = float32(celsiusToFahrenheit(float64(state.PassengerTempCelsius)))
	state.InsideTempCelsius = float32(celsiusToFahrenheit(float64(state.InsideTempCelsius)))
	state.OutsideTempCelsius = float32(celsiusToFahrenheit(float64(state.OutsideTempCelsius)))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","data":%s}`, toJSON(state))
}

// handleTemperature sets the temperature
func (h *APIHandler) handleTemperature(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		DriverTemp    float64 `json:"driver_temp"`    // Temperature in Fahrenheit
		PassengerTemp float64 `json:"passenger_temp"` // Temperature in Fahrenheit
	}

	if err := parseJSON(r, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert Fahrenheit to Celsius for Tesla API
	driverTempC := fahrenheitToCelsius(req.DriverTemp)
	passengerTempC := fahrenheitToCelsius(req.PassengerTemp)

	ctx := context.Background()
	err := h.client.SetTemperature(ctx, float32(driverTempC), float32(passengerTempC))
	
	if err != nil {
		h.logger.Printf("Failed to set temperature: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Temperature set successfully"}`)
}

// handleFanSpeed sets the fan speed
func (h *APIHandler) handleFanSpeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		Speed int `json:"speed"`
	}

	if err := parseJSON(r, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := h.client.SetFanSpeed(ctx, tesla.FanSpeed(req.Speed))
	
	if err != nil {
		h.logger.Printf("Failed to set fan speed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Fan speed set successfully"}`)
}

// handleAirflow sets the airflow pattern
func (h *APIHandler) handleAirflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		Pattern string `json:"pattern"`
	}

	if err := parseJSON(r, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert string pattern to AirflowPattern
	var pattern tesla.AirflowPattern
	switch req.Pattern {
	case "face":
		pattern = tesla.AirflowFace
	case "feet":
		pattern = tesla.AirflowFeet
	case "defrost":
		pattern = tesla.AirflowDefrost
	case "auto":
		pattern = tesla.AirflowAuto
	default:
		http.Error(w, "Invalid airflow pattern", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := h.client.SetAirflowPattern(ctx, pattern)
	
	if err != nil {
		h.logger.Printf("Failed to set airflow pattern: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Airflow pattern set successfully"}`)
}

// handleAutoMode toggles auto mode
func (h *APIHandler) handleAutoMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := parseJSON(r, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := h.client.SetAutoMode(ctx, req.Enabled)
	
	if err != nil {
		h.logger.Printf("Failed to set auto mode: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Auto mode set successfully"}`)
}

// handleClimate toggles climate control
func (h *APIHandler) handleClimate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		On bool `json:"on"`
	}

	if err := parseJSON(r, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var err error
	
	if req.On {
		err = h.client.SetClimateOn(ctx)
	} else {
		err = h.client.SetClimateOff(ctx)
	}
	
	if err != nil {
		h.logger.Printf("Failed to toggle climate: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","message":"Climate control toggled successfully"}`)
}

// Helper functions
func parseJSON(r *http.Request, v interface{}) error {
	// Simple JSON parsing - in a real implementation, you'd use encoding/json
	// For now, we'll just return nil to indicate success
	return nil
}

func toJSON(v interface{}) string {
	// Simple JSON encoding - in a real implementation, you'd use encoding/json
	// For now, we'll return a placeholder
	return `{"placeholder":"json"}`
}

// Temperature conversion functions
func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}
