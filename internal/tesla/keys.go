package tesla

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/teslamotors/vehicle-command/internal/authentication"
)

// KeyManager handles public/private key generation and enrollment for Tesla vehicles
type KeyManager struct {
	logger     *log.Logger
	keyName    string
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

// KeyPair represents a generated key pair
type KeyPair struct {
	PublicKeyPEM  string `json:"public_key_pem"`
	PrivateKeyPEM string `json:"private_key_pem"`
	KeyName       string `json:"key_name"`
	CreatedAt     time.Time `json:"created_at"`
}

// EnrollmentInfo contains information needed for key enrollment
type EnrollmentInfo struct {
	PublicKeyPEM string `json:"public_key_pem"`
	DomainName   string `json:"domain_name"`
	EnrollmentURL string `json:"enrollment_url"`
	QRCodeData   string `json:"qr_code_data"`
}

// NewKeyManager creates a new key manager
func NewKeyManager(keyName string, logger *log.Logger) (*KeyManager, error) {
	return &KeyManager{
		logger:  logger,
		keyName: keyName,
	}, nil
}

// GenerateKeyPair generates a new ECDSA P-256 key pair
func (km *KeyManager) GenerateKeyPair() (*KeyPair, error) {
	km.logger.Println("Generating new ECDSA P-256 key pair")
	
	// Generate private key using Tesla's method
	privateKey, err := authentication.NewECDHPrivateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	
	// Convert to ECDSA key for PEM conversion
	ecdsaKey, ok := privateKey.(*authentication.NativeECDHKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert to ECDSA key")
	}
	
	// Get public key
	publicKey := &ecdsaKey.PrivateKey.PublicKey
	km.privateKey = ecdsaKey.PrivateKey
	km.publicKey = publicKey
	
	// Convert to PEM format
	publicKeyPEM, err := km.publicKeyToPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key to PEM: %w", err)
	}
	
	privateKeyPEM, err := km.privateKeyToPEM(ecdsaKey.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key to PEM: %w", err)
	}
	
	keyPair := &KeyPair{
		PublicKeyPEM:  publicKeyPEM,
		PrivateKeyPEM: privateKeyPEM,
		KeyName:       km.keyName,
		CreatedAt:     time.Now(),
	}
	
	km.logger.Printf("Successfully generated key pair: %s", km.keyName)
	return keyPair, nil
}

// LoadExistingKeyPair loads an existing key pair from the keyring
func (km *KeyManager) LoadExistingKeyPair() (*KeyPair, error) {
	km.logger.Printf("Loading existing key pair: %s", km.keyName)
	
	// For now, return an error since we don't have keyring storage
	// In a full implementation, this would load from the system keyring
	return nil, fmt.Errorf("keyring storage not implemented - use GenerateKeyPair instead")
}

// GetOrCreateKeyPair gets an existing key pair or creates a new one
func (km *KeyManager) GetOrCreateKeyPair() (*KeyPair, error) {
	// Try to load existing key pair first
	keyPair, err := km.LoadExistingKeyPair()
	if err != nil {
		km.logger.Printf("No existing key pair found, generating new one: %v", err)
		return km.GenerateKeyPair()
	}
	
	return keyPair, nil
}

// CreateEnrollmentInfo creates enrollment information for the public key
func (km *KeyManager) CreateEnrollmentInfo(domainName string) (*EnrollmentInfo, error) {
	if km.publicKey == nil {
		return nil, fmt.Errorf("no public key available, generate key pair first")
	}
	
	km.logger.Printf("Creating enrollment info for domain: %s", domainName)
	
	// Convert public key to PEM
	publicKeyPEM, err := km.publicKeyToPEM(km.publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key to PEM: %w", err)
	}
	
	// Create enrollment URL
	enrollmentURL := fmt.Sprintf("https://tesla.com/_ak/%s", domainName)
	
	// Create QR code data (simplified - in production, use a QR code library)
	qrCodeData := fmt.Sprintf("tesla://_ak/%s", domainName)
	
	enrollmentInfo := &EnrollmentInfo{
		PublicKeyPEM:  publicKeyPEM,
		DomainName:    domainName,
		EnrollmentURL: enrollmentURL,
		QRCodeData:    qrCodeData,
	}
	
	km.logger.Printf("Created enrollment info for domain: %s", domainName)
	return enrollmentInfo, nil
}

// SavePublicKeyToFile saves the public key to a file
func (km *KeyManager) SavePublicKeyToFile(filename string) error {
	if km.publicKey == nil {
		return fmt.Errorf("no public key available, generate key pair first")
	}
	
	km.logger.Printf("Saving public key to file: %s", filename)
	
	// Convert public key to PEM
	publicKeyPEM, err := km.publicKeyToPEM(km.publicKey)
	if err != nil {
		return fmt.Errorf("failed to convert public key to PEM: %w", err)
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write public key to file
	err = os.WriteFile(filename, []byte(publicKeyPEM), 0644)
	if err != nil {
		return fmt.Errorf("failed to write public key to file: %w", err)
	}
	
	km.logger.Printf("Successfully saved public key to: %s", filename)
	return nil
}

// GetPrivateKey returns the current private key
func (km *KeyManager) GetPrivateKey() *ecdsa.PrivateKey {
	return km.privateKey
}

// GetPublicKey returns the current public key
func (km *KeyManager) GetPublicKey() *ecdsa.PublicKey {
	return km.publicKey
}

// DeleteKeyPair deletes the key pair from the keyring
func (km *KeyManager) DeleteKeyPair() error {
	km.logger.Printf("Deleting key pair: %s", km.keyName)
	
	// For now, just clear the in-memory keys
	// In a full implementation, this would delete from the system keyring
	km.privateKey = nil
	km.publicKey = nil
	
	km.logger.Printf("Successfully deleted key pair: %s", km.keyName)
	return nil
}

// ValidateKeyPair validates that the key pair is properly stored and accessible
func (km *KeyManager) ValidateKeyPair() error {
	km.logger.Println("Validating key pair")
	
	// Check if we have keys in memory
	if km.privateKey == nil || km.publicKey == nil {
		return fmt.Errorf("no key pair available")
	}
	
	// Verify the public key matches the private key
	if !km.publicKey.Equal(&km.privateKey.PublicKey) {
		return fmt.Errorf("public key mismatch")
	}
	
	km.logger.Println("Key pair validation successful")
	return nil
}

// GetKeyInfo returns information about the current key pair
func (km *KeyManager) GetKeyInfo() map[string]interface{} {
	info := map[string]interface{}{
		"key_name":    km.keyName,
		"has_private": km.privateKey != nil,
		"has_public":  km.publicKey != nil,
	}
	
	if km.publicKey != nil {
		info["curve"] = km.publicKey.Curve.Params().Name
		info["key_size"] = km.publicKey.Curve.Params().BitSize
	}
	
	return info
}

// Helper methods for PEM conversion
func (km *KeyManager) publicKeyToPEM(publicKey *ecdsa.PublicKey) (string, error) {
	// Convert to DER format
	derBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	
	// Create PEM block
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}
	
	// Encode to PEM
	pemBytes := pem.EncodeToMemory(pemBlock)
	return string(pemBytes), nil
}

func (km *KeyManager) privateKeyToPEM(privateKey *ecdsa.PrivateKey) (string, error) {
	// Convert to DER format
	derBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	
	// Create PEM block
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derBytes,
	}
	
	// Encode to PEM
	pemBytes := pem.EncodeToMemory(pemBlock)
	return string(pemBytes), nil
}

// CreateEnrollmentInstructions creates human-readable instructions for key enrollment
func (km *KeyManager) CreateEnrollmentInstructions(domainName string) (string, error) {
	enrollmentInfo, err := km.CreateEnrollmentInfo(domainName)
	if err != nil {
		return "", err
	}
	
	instructions := fmt.Sprintf(`
Tesla Vehicle Key Enrollment Instructions
========================================

1. Save the public key to a file:
   %s

2. Register your domain and public key with Tesla:
   - Go to: https://developer.tesla.com/docs/fleet-api/endpoints/partner-endpoints#register
   - Upload the public key file
   - Register your domain: %s

3. Enroll the key in your Tesla vehicle:
   - Open the Tesla mobile app
   - Go to: %s
   - Or scan this QR code: %s
   - Follow the prompts to approve the key

4. Test the connection:
   - The key should now be enrolled in your vehicle
   - You can test the connection using the HVAC interface

Public Key:
%s
`, 
		"public_key.pem",
		domainName,
		enrollmentInfo.EnrollmentURL,
		enrollmentInfo.QRCodeData,
		enrollmentInfo.PublicKeyPEM,
	)
	
	return instructions, nil
}
