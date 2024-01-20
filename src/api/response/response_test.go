package responses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Sample data for testing
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}

	// Call the JSON function to write the JSON response
	JSON(w, http.StatusOK, data)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse the response body
	var responseBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error decoding JSON: %v", err)
	}

	// Check the content of the response body
	expectedMessage := "Hello, World!"
	if responseBody["message"] != expectedMessage {
		t.Errorf("Expected message %s, got %s", expectedMessage, responseBody["message"])
	}
}

func TestERROR(t *testing.T) {
	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Sample error for testing
	err := someFunctionThatReturnsAnError()

	// Call the ERROR function to write the error response
	ERROR(w, http.StatusBadRequest, err)

	// Check the response status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse the response body
	var responseBody map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error decoding JSON: %v", err)
	}

	// Check the content of the response body
	expectedError := err.Error()
	if responseBody["error"] != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, responseBody["error"])
	}
}

func someFunctionThatReturnsAnError() error {
	// This function is just for testing, you can replace it with your actual error-producing logic
	return nil
}
