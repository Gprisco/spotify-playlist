package callback

import (
	"net/http/httptest"
	"testing"
)

func TestCallback(t *testing.T) {
	t.Run("it should send the code query param back", func(t *testing.T) {
		// Given a channel for catching the result
		channel := make(chan *CallbackResult)

		// and a callback server
		server := &CallbackContext{channel: channel}

		// and a request
		req := httptest.NewRequest("GET", "/callback?code=123code&error=mock-error", nil)

		// and a response recorder
		response := httptest.NewRecorder()

		// When
		go server.ServeHTTP(response, req)

		// Then the code should be as expected
		result := <-channel
		if result.Code != "123code" {
			t.Errorf("Expected code to be 123code, got %s", result.Code)
		}

		// and the error to be as expected
		if result.Err != "mock-error" {
			t.Errorf("Expected error to be mock-error, got %s", result.Err)
		}

		// and the response status code should be 200 OK
		if response.Code != 200 {
			t.Errorf("Expected status code to be 200, got %d", response.Code)
		}
	})
}
