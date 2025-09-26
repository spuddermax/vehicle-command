package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/teslamotors/vehicle-command/internal/tesla"
)

func main() {
	var (
		configPath = flag.String("config", tesla.GetDefaultConfigPath(), "Path to configuration file")
		action     = flag.String("action", "show", "Action to perform: show, create, validate, set-vin, set-key")
		vin        = flag.String("vin", "", "Vehicle VIN (for set-vin action)")
		keyFile    = flag.String("key-file", "", "Private key file path (for set-key action)")
		tokenFile  = flag.String("token-file", "", "OAuth token file path (for set-token action)")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	switch *action {
	case "show":
		showConfig(*configPath)
	case "create":
		createConfig(*configPath)
	case "validate":
		validateConfig(*configPath)
	case "set-vin":
		if *vin == "" {
			fmt.Fprintf(os.Stderr, "Error: VIN is required for set-vin action\n")
			os.Exit(1)
		}
		setVIN(*configPath, *vin)
	case "set-key":
		if *keyFile == "" {
			fmt.Fprintf(os.Stderr, "Error: key-file is required for set-key action\n")
			os.Exit(1)
		}
		setKeyFile(*configPath, *keyFile)
	case "set-token":
		if *tokenFile == "" {
			fmt.Fprintf(os.Stderr, "Error: token-file is required for set-token action\n")
			os.Exit(1)
		}
		setTokenFile(*configPath, *tokenFile)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown action '%s'\n", *action)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Tesla HVAC Configuration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tesla-config [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -config string")
	fmt.Println("        Path to configuration file (default: ~/.config/tesla-hvac/config.json)")
	fmt.Println("  -action string")
	fmt.Println("        Action to perform: show, create, validate, set-vin, set-key, set-token (default: show)")
	fmt.Println("  -vin string")
	fmt.Println("        Vehicle VIN (for set-vin action)")
	fmt.Println("  -key-file string")
	fmt.Println("        Private key file path (for set-key action)")
	fmt.Println("  -token-file string")
	fmt.Println("        OAuth token file path (for set-token action)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  show      - Display current configuration")
	fmt.Println("  create    - Create a new configuration file with defaults")
	fmt.Println("  validate  - Validate the configuration file")
	fmt.Println("  set-vin   - Set the vehicle VIN")
	fmt.Println("  set-key   - Set the private key file path")
	fmt.Println("  set-token - Set the OAuth token file path")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  tesla-config -action create")
	fmt.Println("  tesla-config -action set-vin -vin 5YJ3E1EA4KF123456")
	fmt.Println("  tesla-config -action set-key -key-file ~/.tesla/private_key.pem")
	fmt.Println("  tesla-config -action validate")
}

func showConfig(configPath string) {
	config, err := tesla.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Configuration loaded from: %s\n", configPath)
	fmt.Println()
	fmt.Printf("Tesla Configuration:\n")
	fmt.Printf("  VIN: %s\n", config.Tesla.VIN)
	fmt.Printf("  Private Key File: %s\n", config.Tesla.PrivateKeyFile)
	fmt.Printf("  OAuth Token File: %s\n", config.Tesla.OAuthTokenFile)
	fmt.Printf("  Connection Timeout: %v\n", config.Tesla.ConnectionTimeout)
	fmt.Printf("  Scan Timeout: %v\n", config.Tesla.ScanTimeout)
	fmt.Printf("  Max Concurrent Requests: %d\n", config.Tesla.MaxConcurrentRequests)
	fmt.Printf("  Request Timeout: %v\n", config.Tesla.RequestTimeout)
	fmt.Printf("  Scan Retries: %d\n", config.Tesla.ScanRetries)
	fmt.Printf("  Scan Delay: %v\n", config.Tesla.ScanDelay)
	fmt.Println()
	fmt.Printf("Client Configuration:\n")
	fmt.Printf("  Client Name: %s\n", config.Client.ClientName)
	fmt.Printf("  Client Version: %s\n", config.Client.ClientVersion)
	fmt.Printf("  Keep Alive Interval: %v\n", config.Client.KeepAliveInterval)
	fmt.Printf("  Health Check Interval: %v\n", config.Client.HealthCheckInterval)
	fmt.Printf("  Enable Auto Reconnect: %t\n", config.Client.EnableAutoReconnect)
	fmt.Printf("  Enable Health Checks: %t\n", config.Client.EnableHealthChecks)
	fmt.Printf("  Enable Metrics: %t\n", config.Client.EnableMetrics)
	fmt.Println()
	fmt.Printf("Retry Configuration:\n")
	fmt.Printf("  Max Retries: %d\n", config.Retry.MaxRetries)
	fmt.Printf("  Initial Delay: %v\n", config.Retry.InitialDelay)
	fmt.Printf("  Max Delay: %v\n", config.Retry.MaxDelay)
	fmt.Printf("  Backoff Factor: %.2f\n", config.Retry.BackoffFactor)
	fmt.Printf("  Jitter: %t\n", config.Retry.Jitter)
	fmt.Println()
	fmt.Printf("Circuit Breaker Configuration:\n")
	fmt.Printf("  Max Failures: %d\n", config.CircuitBreaker.MaxFailures)
	fmt.Printf("  Reset Timeout: %v\n", config.CircuitBreaker.ResetTimeout)
	fmt.Printf("  Half Open Max Calls: %d\n", config.CircuitBreaker.HalfOpenMaxCalls)
	fmt.Println()
	fmt.Printf("Logging Configuration:\n")
	fmt.Printf("  Level: %s\n", config.Logging.Level)
	fmt.Printf("  Format: %s\n", config.Logging.Format)
	fmt.Printf("  Output: %s\n", config.Logging.Output)
	fmt.Printf("  File Path: %s\n", config.Logging.FilePath)
	fmt.Printf("  Max Size: %d MB\n", config.Logging.MaxSize)
	fmt.Printf("  Max Backups: %d\n", config.Logging.MaxBackups)
	fmt.Printf("  Max Age: %d days\n", config.Logging.MaxAge)
}

func createConfig(configPath string) {
	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists: %s\n", configPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Configuration creation cancelled.")
			return
		}
	}

	// Create default config
	config := tesla.DefaultConfig()
	config.ConfigPath = configPath

	// Load environment variables
	config.LoadFromEnv()

	// Save config
	if err := config.Save(); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("Configuration created successfully: %s\n", configPath)
	fmt.Println("You can now edit the configuration file or use the set-* actions to configure specific values.")
}

func validateConfig(configPath string) {
	config, err := tesla.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := config.Validate(); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration is valid.")
}

func setVIN(configPath string, vin string) {
	config, err := tesla.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	config.Tesla.VIN = vin

	if err := config.Save(); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("VIN set to: %s\n", vin)
}

func setKeyFile(configPath string, keyFile string) {
	config, err := tesla.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Check if key file exists
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		fmt.Printf("Warning: Key file does not exist: %s\n", keyFile)
	}

	config.Tesla.PrivateKeyFile = keyFile

	if err := config.Save(); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("Private key file set to: %s\n", keyFile)
}

func setTokenFile(configPath string, tokenFile string) {
	config, err := tesla.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Check if token file exists
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		fmt.Printf("Warning: Token file does not exist: %s\n", tokenFile)
	}

	config.Tesla.OAuthTokenFile = tokenFile

	if err := config.Save(); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("OAuth token file set to: %s\n", tokenFile)
}
