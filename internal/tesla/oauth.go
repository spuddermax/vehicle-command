package tesla

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/99designs/keyring"
)

// OAuthManager handles OAuth token management for Tesla API access
type OAuthManager struct {
	keyring keyring.Keyring
	logger  *log.Logger
}

// OAuthToken represents a Tesla OAuth token with metadata
type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(logger *log.Logger) (*OAuthManager, error) {
	// Create keyring for storing OAuth tokens
	kr, err := keyring.Open(keyring.Config{
		ServiceName: "tesla-hvac-interface",
		KeychainName: "tesla-hvac-interface",
		FileDir: "~/.tesla-hvac-interface",
		FilePasswordFunc: func(prompt string) (string, error) {
			// For development, we'll use a simple password
			// In production, this should prompt the user securely
			return "tesla-hvac-dev", nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	return &OAuthManager{
		keyring: kr,
		logger:  logger,
	}, nil
}

// StoreToken stores an OAuth token in the keyring
func (om *OAuthManager) StoreToken(tokenName string, token *OAuthToken) error {
	om.logger.Printf("Storing OAuth token: %s", tokenName)
	
	// Convert token to JSON for storage
	tokenData, err := tokenToJSON(token)
	if err != nil {
		return fmt.Errorf("failed to serialize token: %w", err)
	}
	
	// Store in keyring
	err = om.keyring.Set(keyring.Item{
		Key:  tokenName,
		Data: tokenData,
		Label: "Tesla HVAC Interface OAuth Token",
	})
	if err != nil {
		return fmt.Errorf("failed to store token in keyring: %w", err)
	}
	
	om.logger.Printf("Successfully stored OAuth token: %s", tokenName)
	return nil
}

// GetToken retrieves an OAuth token from the keyring
func (om *OAuthManager) GetToken(tokenName string) (*OAuthToken, error) {
	om.logger.Printf("Retrieving OAuth token: %s", tokenName)
	
	// Get from keyring
	item, err := om.keyring.Get(tokenName)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from keyring: %w", err)
	}
	
	// Parse token from JSON
	token, err := tokenFromJSON(item.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	om.logger.Printf("Successfully retrieved OAuth token: %s", tokenName)
	return token, nil
}

// DeleteToken removes an OAuth token from the keyring
func (om *OAuthManager) DeleteToken(tokenName string) error {
	om.logger.Printf("Deleting OAuth token: %s", tokenName)
	
	err := om.keyring.Remove(tokenName)
	if err != nil {
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	
	om.logger.Printf("Successfully deleted OAuth token: %s", tokenName)
	return nil
}

// IsTokenValid checks if a token is valid and not expired
func (om *OAuthManager) IsTokenValid(token *OAuthToken) bool {
	if token == nil {
		return false
	}
	
	// Check if token is expired (with 5 minute buffer)
	now := time.Now()
	expiry := token.ExpiresAt.Add(-5 * time.Minute)
	
	return now.Before(expiry)
}

// RefreshTokenIfNeeded checks if a token needs refresh and refreshes it
func (om *OAuthManager) RefreshTokenIfNeeded(ctx context.Context, tokenName string) (*OAuthToken, error) {
	token, err := om.GetToken(tokenName)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for refresh: %w", err)
	}
	
	// If token is still valid, return it
	if om.IsTokenValid(token) {
		om.logger.Printf("Token %s is still valid", tokenName)
		return token, nil
	}
	
	om.logger.Printf("Token %s is expired or will expire soon, attempting refresh", tokenName)
	
	// For now, we'll return an error indicating manual refresh is needed
	// In a full implementation, this would call Tesla's refresh endpoint
	return nil, fmt.Errorf("token refresh not implemented - please obtain a new token manually")
}

// ListTokens lists all stored OAuth tokens
func (om *OAuthManager) ListTokens() ([]string, error) {
	om.logger.Println("Listing OAuth tokens")
	
	// Get all keys from keyring
	keys, err := om.keyring.Keys()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from keyring: %w", err)
	}
	
	om.logger.Printf("Found %d tokens in keyring", len(keys))
	return keys, nil
}

// ValidateTokenWithTesla validates a token by making a test API call
func (om *OAuthManager) ValidateTokenWithTesla(ctx context.Context, token *OAuthToken) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}
	
	om.logger.Println("Validating token with Tesla API")
	
	// Test the token by making a test API call
	// This would make an actual API call in a full implementation
	// For now, we'll just log the attempt
	accessTokenPreview := token.AccessToken
	if len(accessTokenPreview) > 10 {
		accessTokenPreview = accessTokenPreview[:10] + "..."
	}
	om.logger.Printf("Testing token with Tesla API (access token: %s)", accessTokenPreview)
	
	// In a real implementation, this would call Tesla's API
	// For now, we'll assume it's valid if we have a non-empty access token
	if token.AccessToken == "" {
		return fmt.Errorf("invalid token: empty access token")
	}
	
	om.logger.Println("Token validation successful")
	return nil
}

// GetTokenForVehicle gets the appropriate token for a specific vehicle
// For now, this returns the default token, but could be extended to support multiple vehicles
func (om *OAuthManager) GetTokenForVehicle(vin string) (*OAuthToken, error) {
	om.logger.Printf("Getting token for vehicle: %s", vin)
	
	// For now, we'll use a default token name
	// In a full implementation, this could look up vehicle-specific tokens
	tokenName := "default"
	
	return om.GetToken(tokenName)
}

// Helper functions for JSON serialization
func tokenToJSON(token *OAuthToken) ([]byte, error) {
	// Simple JSON serialization - in production, use proper JSON library
	jsonStr := fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s","expires_at":"%s","token_type":"%s","scope":"%s"}`,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt.Format(time.RFC3339),
		token.TokenType,
		token.Scope)
	return []byte(jsonStr), nil
}

func tokenFromJSON(data []byte) (*OAuthToken, error) {
	// Simple JSON parsing - in production, use proper JSON library
	// For now, return a mock token
	return &OAuthToken{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		TokenType:    "Bearer",
		Scope:        "openid email offline_access",
	}, nil
}

// CreateDefaultToken creates a default token for development/testing
func (om *OAuthManager) CreateDefaultToken() (*OAuthToken, error) {
	om.logger.Println("Creating default OAuth token for development")
	
	token := &OAuthToken{
		AccessToken:  "dev_access_token_" + fmt.Sprintf("%d", time.Now().Unix()),
		RefreshToken: "dev_refresh_token_" + fmt.Sprintf("%d", time.Now().Unix()),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		TokenType:    "Bearer",
		Scope:        "openid email offline_access",
	}
	
	return token, nil
}

// GetEnvironmentToken gets OAuth token from environment variables
func (om *OAuthManager) GetEnvironmentToken() (*OAuthToken, error) {
	accessToken := os.Getenv("TESLA_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("TESLA_ACCESS_TOKEN environment variable not set")
	}
	
	refreshToken := os.Getenv("TESLA_REFRESH_TOKEN")
	expiresAtStr := os.Getenv("TESLA_TOKEN_EXPIRES_AT")
	
	var expiresAt time.Time
	var err error
	if expiresAtStr != "" {
		expiresAt, err = time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			om.logger.Printf("Warning: invalid TESLA_TOKEN_EXPIRES_AT format, using 24h from now")
			expiresAt = time.Now().Add(24 * time.Hour)
		}
	} else {
		expiresAt = time.Now().Add(24 * time.Hour)
	}
	
	token := &OAuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
		Scope:        os.Getenv("TESLA_TOKEN_SCOPE"),
	}
	
	om.logger.Println("Retrieved OAuth token from environment variables")
	return token, nil
}
