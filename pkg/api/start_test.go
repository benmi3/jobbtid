/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETStart(t *testing.T) {
	t.Run("No Query returns Bad Reqest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/start", nil)
		response := httptest.NewRecorder()

		start(response, request)

		got := response.Body.String()
		want := `{"statusCode":400,"message":"Bad Request"}`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if response.Code != http.StatusBadRequest {
			t.Errorf("got %d, want %d", response.Code, http.StatusBadRequest)
		}
	})
}

func TestPOSTStart(t *testing.T) {
	t.Run("NoBody return Bad Request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/start", nil)
		response := httptest.NewRecorder()

		start(response, request)

		got := response.Body.String()
		want := `{"statusCode":400,"message":"Bad Request"}`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
		if response.Code != http.StatusBadRequest {
			t.Errorf("got %d, want %d", response.Code, http.StatusBadRequest)
		}
	})
}
