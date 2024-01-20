package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetMiddlewareJSON(t *testing.T) {
	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/test", nil)

	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Define a test handler using SetMiddlewareJSON
	handler := SetMiddlewareJSON(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the Content-Type header is set to application/json
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", contentType)
		}
	})

	// Call the test handler with the test request and response recorder
	handler(w, req)
}

//func TestSetMiddlewareAuthentication(t *testing.T) {
//	// Create a test HTTP request
//	req := httptest.NewRequest("GET", "/test", nil)
//
//	// Create a test HTTP response recorder
//	w := httptest.NewRecorder()
//
//	// Mock the TokenValid function to always return nil (no error)
//	mockTokenValid := func(r *http.Request) error {
//		return nil
//	}
//
//	// Define a test handler using SetMiddlewareAuthentication
//	handler := SetMiddlewareAuthentication(func(w http.ResponseWriter, r *http.Request) {
//		// This function will be called if authentication is successful
//	})
//
//	// Replace the original TokenValid function with the mock function
//	originalTokenValid := auth.TokenValid
//	auth.TokenValid = mockTokenValid
//	defer func() {
//		auth.TokenValid = originalTokenValid
//	}()
//
//	// Call the test handler with the test request and response recorder
//	handler(w, req)
//
//	// Verify that the status code is not Unauthorized (401)
//	if w.Code == http.StatusUnauthorized {
//		t.Errorf("Expected status code other than Unauthorized (401)")
//	}
//}
