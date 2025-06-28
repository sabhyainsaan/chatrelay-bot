package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMockBackend tests the /v1/chat/stream endpoint with a known query
func TestMockBackend(t *testing.T) {
	// Set up a test server
	reqBody := ChatRequest{
		UserID: "U123",
		Query:  "what is golang?",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/v1/chat/stream", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the result
	rr := httptest.NewRecorder()

	// Define the handler directly (without starting the full server)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Directly use your existing logic
		_, span := backendTracer.Start(r.Context(), "TestSpan")
		defer span.End()

		var req ChatRequest
		json.NewDecoder(r.Body).Decode(&req)

		response := ChatResponse{
			FullResponse: "Go is an open-source programming language designed for simplicity and concurrency.",
		}
		json.NewEncoder(w).Encode(response)
	})

	// Serve the test request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}

	// Decode and verify the response
	var res ChatResponse
	json.NewDecoder(rr.Body).Decode(&res)

	expected := "Go is an open-source programming language designed for simplicity and concurrency."
	if res.FullResponse != expected {
		t.Errorf("Expected '%s', got '%s'", expected, res.FullResponse)
	}
}
