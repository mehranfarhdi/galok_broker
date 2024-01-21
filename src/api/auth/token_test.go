package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreateTokenAndExtractTokenID(t *testing.T) {
	// Set your API_SECRET for testing (replace with your actual secret)
	os.Setenv("API_SECRET", "your_secret_key")

	// Create a token for testing
	userID := uint32(123)

	username := "mohammad"

	email := "mohammad@gmail.com"

	isAdmin := true

	token, err := CreateToken(userID, username, email, isAdmin)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	// Create a mock HTTP request with the token in the Authorization header
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Extract the user ID from the token in the request
	extractedUserID, err := ExtractTokenID(req)
	if err != nil {
		t.Fatalf("Error extracting user ID: %v", err)
	}

	// Verify that the extracted user ID matches the original user ID
	if extractedUserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, extractedUserID)
	}
}

func TestTokenValid(t *testing.T) {
	// Set your API_SECRET for testing (replace with your actual secret)
	os.Setenv("API_SECRET", "your_secret_key")

	// Create a token for testing
	userID := uint32(123)

	username := "mohammad"

	email := "mohammad@gmail.com"

	isAdmin := true

	token, err := CreateToken(userID, username, email, isAdmin)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	// Create a mock HTTP request with the token in the Authorization header
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Check if the token in the request is valid
	err = TokenValid(req)
	if err != nil {
		t.Errorf("Token validation failed: %v", err)
	}
}

func TestExtractToken(t *testing.T) {
	// Create a mock HTTP request with the token in the URL query
	req := httptest.NewRequest("GET", "/test?token=testToken", nil)

	// Extract the token from the request
	extractedToken := ExtractToken(req)

	// Verify that the extracted token matches the expected value
	expectedToken := "testToken"
	if extractedToken != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, extractedToken)
	}
}
