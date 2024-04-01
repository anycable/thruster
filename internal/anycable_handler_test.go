package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnyCableHandler_ServeHTTP(t *testing.T) {
	websocketHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("WebSocket handler response"))
	})

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Next handler response"))
	})

	handler := &AnyCableHandler{
		handler: websocketHandler,
		next:    nextHandler,
	}

	t.Run("Redirect to WebSocketHandler for /cable", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/cable", nil)
		handler.ServeHTTP(recorder, request)

		response := recorder.Result()
		body := recorder.Body.String()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		if body != "WebSocket handler response" {
			t.Errorf("Expected body to be 'WebSocket handler response', got '%s'", body)
		}
	})

	t.Run("Pass other requests to the next handler", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/not_cable", nil)
		handler.ServeHTTP(recorder, request)

		response := recorder.Result()
		body := recorder.Body.String()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		if body != "Next handler response" {
			t.Errorf("Expected body to be 'Next handler response', got '%s'", body)
		}
	})
}
