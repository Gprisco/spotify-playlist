package callback

import (
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	t.Run("it should return the callback result",
		func(t *testing.T) {
			// Given a timeout
			timeout := 1 * time.Second

			// prepare to send an http GET request to /callback including code and error query param
			go sendCallback(t, "expectedCode", "expectedError")

			// When handling the callback
			result := HandleCallback(timeout)

			// Then the result should be returned
			if result == nil {
				t.Errorf("The callback result was not returned")
			}
		},
	)
}

// NOTE: this code assumes that the callback server starts in 5 seconds
func sendCallback(t *testing.T, expectedCode string, expectedError string) {
	for range time.Tick(time.Second * 5) {
    client := &http.Client{}

    req, err := http.NewRequest(
      http.MethodGet,
      "http://localhost:8080/callback?code="+expectedCode+"&error="+expectedError,
      nil,
    )

    if err != nil {
      t.Errorf("Error creating the request: %s", err.Error())
    }

    _, err = client.Do(req)
    if err != nil {
      t.Errorf("Error sending the request: %s", err.Error())
    }

    break
  }
}
