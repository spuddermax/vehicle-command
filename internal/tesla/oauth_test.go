package tesla

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewOAuthManager(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	
	// Test creating OAuth manager
	manager, err := NewOAuthManager(logger)
	if err != nil {
		t.Fatalf("Failed to create OAuth manager: %v", err)
	}
	
	if manager == nil {
		t.Fatal("OAuth manager is nil")
	}
	
	if manager.keyring == nil {
		t.Error("Keyring is nil")
	}
	
	if manager.logger == nil {
		t.Error("Logger is nil")
	}
}

func TestOAuthToken(t *testing.T) {
	token := &OAuthToken{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		TokenType:    "Bearer",
		Scope:        "openid email offline_access",
	}
	
	if token.AccessToken != "test_access_token" {
		t.Errorf("Expected access token 'test_access_token', got '%s'", token.AccessToken)
	}
	
	if token.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", token.TokenType)
	}
	
	if token.ExpiresAt.Before(time.Now()) {
		t.Error("Token should not be expired")
	}
}

func TestIsTokenValid(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewOAuthManager(logger)
	if err != nil {
		t.Fatalf("Failed to create OAuth manager: %v", err)
	}
	
	// Test valid token
	validToken := &OAuthToken{
		AccessToken: "test_token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	
	if !manager.IsTokenValid(validToken) {
		t.Error("Valid token should be considered valid")
	}
	
	// Test expired token
	expiredToken := &OAuthToken{
		AccessToken: "test_token",
		ExpiresAt:   time.Now().Add(-1 * time.Hour),
	}
	
	if manager.IsTokenValid(expiredToken) {
		t.Error("Expired token should not be considered valid")
	}
	
	// Test nil token
	if manager.IsTokenValid(nil) {
		t.Error("Nil token should not be considered valid")
	}
}

func TestCreateDefaultToken(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewOAuthManager(logger)
	if err != nil {
		t.Fatalf("Failed to create OAuth manager: %v", err)
	}
	
	token, err := manager.CreateDefaultToken()
	if err != nil {
		t.Fatalf("Failed to create default token: %v", err)
	}
	
	if token == nil {
		t.Fatal("Default token is nil")
	}
	
	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}
	
	if token.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", token.TokenType)
	}
	
	if !manager.IsTokenValid(token) {
		t.Error("Default token should be valid")
	}
}

func TestGetEnvironmentToken(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewOAuthManager(logger)
	if err != nil {
		t.Fatalf("Failed to create OAuth manager: %v", err)
	}
	
	// Test with environment variables set
	os.Setenv("TESLA_ACCESS_TOKEN", "env_access_token")
	os.Setenv("TESLA_REFRESH_TOKEN", "env_refresh_token")
	os.Setenv("TESLA_TOKEN_EXPIRES_AT", time.Now().Add(24*time.Hour).Format(time.RFC3339))
	os.Setenv("TESLA_TOKEN_SCOPE", "openid email offline_access")
	
	token, err := manager.GetEnvironmentToken()
	if err != nil {
		t.Fatalf("Failed to get environment token: %v", err)
	}
	
	if token.AccessToken != "env_access_token" {
		t.Errorf("Expected access token 'env_access_token', got '%s'", token.AccessToken)
	}
	
	if token.RefreshToken != "env_refresh_token" {
		t.Errorf("Expected refresh token 'env_refresh_token', got '%s'", token.RefreshToken)
	}
	
	// Test without environment variables
	os.Unsetenv("TESLA_ACCESS_TOKEN")
	
	_, err = manager.GetEnvironmentToken()
	if err == nil {
		t.Error("Expected error when TESLA_ACCESS_TOKEN is not set")
	}
	
	// Clean up
	os.Unsetenv("TESLA_REFRESH_TOKEN")
	os.Unsetenv("TESLA_TOKEN_EXPIRES_AT")
	os.Unsetenv("TESLA_TOKEN_SCOPE")
}

func TestValidateTokenWithTesla(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewOAuthManager(logger)
	if err != nil {
		t.Fatalf("Failed to create OAuth manager: %v", err)
	}
	
	ctx := context.Background()
	
	// Test valid token
	validToken := &OAuthToken{
		AccessToken: "valid_token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	
	err = manager.ValidateTokenWithTesla(ctx, validToken)
	if err != nil {
		t.Errorf("Valid token should pass validation: %v", err)
	}
	
	// Test nil token
	err = manager.ValidateTokenWithTesla(ctx, nil)
	if err == nil {
		t.Error("Nil token should fail validation")
	}
	
	// Test empty token
	emptyToken := &OAuthToken{
		AccessToken: "",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	
	err = manager.ValidateTokenWithTesla(ctx, emptyToken)
	if err == nil {
		t.Error("Empty token should fail validation")
	}
}

func TestGetTokenForVehicle(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", log.LstdFlags)
	manager, err := NewOAuthManager(logger)
	if err != nil {
		t.Fatalf("Failed to create OAuth manager: %v", err)
	}
	
	// This test will fail because we haven't stored a token yet
	// In a real test, we would store a token first
	_, err = manager.GetTokenForVehicle("TEST_VIN_123")
	if err == nil {
		t.Error("Expected error when no token is stored")
	}
}
