package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Tracer for the mock backend service
var backendTracer = otel.Tracer("chatrelay/backend")

type ChatRequest struct {
	UserID string `json:"user_id"`
	Query  string `json:"query"`
}

type ChatResponse struct {
	FullResponse string `json:"full_response"`
}

func StartMockServer() {
	http.HandleFunc("/v1/chat/stream", func(w http.ResponseWriter, r *http.Request) {
		// Start tracing span for incoming backend request
		_, span := backendTracer.Start(r.Context(), "MockBackendResponse")
		defer span.End()

		var req ChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Add request data to trace
		span.SetAttributes(
			attribute.String("user_id", req.UserID),
			attribute.String("query", req.Query),
		)

		// Simulate a delay like real backend
		time.Sleep(2 * time.Second)

		response := ChatResponse{
			FullResponse: "Goroutines are lightweight threads managed by the Go runtime.",
		}

		// Add response to trace
		span.SetAttributes(attribute.String("response", response.FullResponse))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Mock backend started on http://localhost:8080")
	go http.ListenAndServe(":8080", nil)
}
