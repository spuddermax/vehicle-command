package tesla

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestNewKeyManager(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	// Test creating key manager
	manager, err := NewKeyManager("test-key", logger)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}
	
	if manager == nil {
		t.Fatal("Key manager is nil")
	}
	
	// Keyring is not used in this simplified implementation
	
	if manager.logger == nil {
		t.Error("Logger is nil")
	}
	
	if manager.keyName != "test-key" {
		t.Errorf("Expected key name 'test-key', got '%s'", manager.keyName)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewKeyManager("test-generate", logger)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}
	
	// Generate key pair
	keyPair, err := manager.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	
	if keyPair == nil {
		t.Fatal("Key pair is nil")
	}
	
	if keyPair.KeyName != "test-generate" {
		t.Errorf("Expected key name 'test-generate', got '%s'", keyPair.KeyName)
	}
	
	if keyPair.PublicKeyPEM == "" {
		t.Error("Public key PEM should not be empty")
	}
	
	if keyPair.PrivateKeyPEM == "" {
		t.Error("Private key PEM should not be empty")
	}
	
	if keyPair.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	
	// Verify we can get the keys
	if manager.GetPrivateKey() == nil {
		t.Error("Private key should be available")
	}
	
	if manager.GetPublicKey() == nil {
		t.Error("Public key should be available")
	}
}

func TestKeyPair(t *testing.T) {
	keyPair := &KeyPair{
		PublicKeyPEM:  "test-public-key",
		PrivateKeyPEM: "test-private-key",
		KeyName:       "test-key",
		CreatedAt:     time.Now(),
	}
	
	if keyPair.PublicKeyPEM != "test-public-key" {
		t.Errorf("Expected public key 'test-public-key', got '%s'", keyPair.PublicKeyPEM)
	}
	
	if keyPair.PrivateKeyPEM != "test-private-key" {
		t.Errorf("Expected private key 'test-private-key', got '%s'", keyPair.PrivateKeyPEM)
	}
	
	if keyPair.KeyName != "test-key" {
		t.Errorf("Expected key name 'test-key', got '%s'", keyPair.KeyName)
	}
}

func TestEnrollmentInfo(t *testing.T) {
	enrollmentInfo := &EnrollmentInfo{
		PublicKeyPEM:  "test-public-key",
		DomainName:    "example.com",
		EnrollmentURL: "https://tesla.com/_ak/example.com",
		QRCodeData:    "tesla://_ak/example.com",
	}
	
	if enrollmentInfo.PublicKeyPEM != "test-public-key" {
		t.Errorf("Expected public key 'test-public-key', got '%s'", enrollmentInfo.PublicKeyPEM)
	}
	
	if enrollmentInfo.DomainName != "example.com" {
		t.Errorf("Expected domain 'example.com', got '%s'", enrollmentInfo.DomainName)
	}
	
	if enrollmentInfo.EnrollmentURL != "https://tesla.com/_ak/example.com" {
		t.Errorf("Expected enrollment URL 'https://tesla.com/_ak/example.com', got '%s'", enrollmentInfo.EnrollmentURL)
	}
	
	if enrollmentInfo.QRCodeData != "tesla://_ak/example.com" {
		t.Errorf("Expected QR code data 'tesla://_ak/example.com', got '%s'", enrollmentInfo.QRCodeData)
	}
}

func TestCreateEnrollmentInfo(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewKeyManager("test-enrollment", logger)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}
	
	// Generate key pair first
	_, err = manager.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	
	// Create enrollment info
	enrollmentInfo, err := manager.CreateEnrollmentInfo("example.com")
	if err != nil {
		t.Fatalf("Failed to create enrollment info: %v", err)
	}
	
	if enrollmentInfo == nil {
		t.Fatal("Enrollment info is nil")
	}
	
	if enrollmentInfo.DomainName != "example.com" {
		t.Errorf("Expected domain 'example.com', got '%s'", enrollmentInfo.DomainName)
	}
	
	if enrollmentInfo.EnrollmentURL != "https://tesla.com/_ak/example.com" {
		t.Errorf("Expected enrollment URL 'https://tesla.com/_ak/example.com', got '%s'", enrollmentInfo.EnrollmentURL)
	}
	
	if enrollmentInfo.PublicKeyPEM == "" {
		t.Error("Public key PEM should not be empty")
	}
}

func TestCreateEnrollmentInfoWithoutKey(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewKeyManager("test-no-key", logger)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}
	
	// Try to create enrollment info without generating a key pair
	_, err = manager.CreateEnrollmentInfo("example.com")
	if err == nil {
		t.Error("Expected error when no public key is available")
	}
}

func TestGetKeyInfo(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewKeyManager("test-info", logger)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}
	
	// Get key info before generating keys
	info := manager.GetKeyInfo()
	if info["key_name"] != "test-info" {
		t.Errorf("Expected key name 'test-info', got '%s'", info["key_name"])
	}
	
	if info["has_private"] != false {
		t.Error("Expected has_private to be false")
	}
	
	if info["has_public"] != false {
		t.Error("Expected has_public to be false")
	}
	
	// Generate key pair
	_, err = manager.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	
	// Get key info after generating keys
	info = manager.GetKeyInfo()
	if info["has_private"] != true {
		t.Error("Expected has_private to be true")
	}
	
	if info["has_public"] != true {
		t.Error("Expected has_public to be true")
	}
	
	if info["curve"] != "P-256" {
		t.Errorf("Expected curve 'P-256', got '%s'", info["curve"])
	}
	
	if info["key_size"] != 256 {
		t.Errorf("Expected key size 256, got %v", info["key_size"])
	}
}

func TestCreateEnrollmentInstructions(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewKeyManager("test-instructions", logger)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}
	
	// Generate key pair first
	_, err = manager.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	
	// Create enrollment instructions
	instructions, err := manager.CreateEnrollmentInstructions("example.com")
	if err != nil {
		t.Fatalf("Failed to create enrollment instructions: %v", err)
	}
	
	if instructions == "" {
		t.Error("Instructions should not be empty")
	}
	
	// Check that instructions contain expected content
	if !contains(instructions, "Tesla Vehicle Key Enrollment Instructions") {
		t.Error("Instructions should contain title")
	}
	
	if !contains(instructions, "example.com") {
		t.Error("Instructions should contain domain name")
	}
	
	if !contains(instructions, "https://tesla.com/_ak/example.com") {
		t.Error("Instructions should contain enrollment URL")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && contains(s[1:], substr)
}
