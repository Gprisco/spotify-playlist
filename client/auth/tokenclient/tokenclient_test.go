package tokenclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

// To mock the http client, we mock the underlying roundtripper
// mockRoundTripper implements http.RoundTripper for testing purposes
type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestSpotifyTokenClient(t *testing.T) {
	t.Run("it should build a valid request and return the access token",
		func(t *testing.T) {
			// Given a mock round tripper and some assertions on the http request
			mockResponse := `{"access_token": "expected-access-token"}`
			mockRoundTripper := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					// Check if the request has the expected properties
					if req.Method != "POST" {
						t.Errorf("Expected POST method, got %s", req.Method)
					}

					if req.URL.String() != "https://accounts.spotify.com/api/token" {
						t.Errorf("Expected URL 'https://accounts.spotify.com/api/token', got %s", req.URL.String())
					}

					err := req.ParseForm()
					if err != nil {
						t.Fatal("Error during form parsing")
					}
					assertPostFormParam(t, req.PostForm, "grant_type", "authorization_code")
					assertPostFormParam(t, req.PostForm, "code", "expected-code")
					assertPostFormParam(t, req.PostForm, "redirect_uri", "expected-redirect-uri")
					assertPostFormParam(t, req.PostForm, "client_id", "expected-client-id")
					assertPostFormParam(t, req.PostForm, "code_verifier", "expected-code-verifier")

					// Create a mock response with no errors
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
						Header:     make(http.Header),
					}, nil
				},
			}

			// And the subject under test using the above mock
			tokenClient := SpotifyTokenClient{
				client:      &http.Client{Transport: mockRoundTripper},
				clientId:    "expected-client-id",
				redirectUri: "expected-redirect-uri",
			}

			// When calling GetToken
			accessToken, err := tokenClient.GetToken("expected-code", "expected-code-verifier")

			if err != nil {
				t.Errorf("GetToken returned an error: %s", err.Error())
			}

			if accessToken != "expected-access-token" {
				t.Errorf("Expected access token 'expected-access-token', but got: %s", accessToken)
			}
		},
	)

	t.Run("it should return the error from http client if any",
		func(t *testing.T) {
			// Given a failing round tripper
			mockRoundTripper := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("mock http error")
				},
			}

			// And a spotify token client using it
			tokenClient := SpotifyTokenClient{
				client:      &http.Client{Transport: mockRoundTripper},
				clientId:    "expected-client-id",
				redirectUri: "expected-redirect-uri",
			}

			// When calling GetToken
			_, err := tokenClient.GetToken("a code", "a verifier")

			const expectedError = `Post "https://accounts.spotify.com/api/token": mock http error`
			if err.Error() != expectedError {
				t.Errorf("Expected '%s', but got '%s'", expectedError, err.Error())
			}
		},
	)

	t.Run("it should return an error if the status code is not 200",
		func(t *testing.T) {
			// Given a round tripper returning a non-OK response
			mockRoundTripper := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusInternalServerError}, nil
				},
			}

			// And a spotify token client using it
			tokenClient := SpotifyTokenClient{
				client:      &http.Client{Transport: mockRoundTripper},
				clientId:    "expected-client-id",
				redirectUri: "expected-redirect-uri",
			}

			// When calling GetToken
			_, err := tokenClient.GetToken("a code", "a verifier")

			// Then an error should be returned
			expectedError := fmt.Sprintf("received non-OK response: %d", http.StatusInternalServerError)
			if err.Error() != expectedError {
				t.Errorf("Expected '%s', but got '%s'", expectedError, err.Error())
			}
		},
	)

	t.Run("it should return an error if the access token is not present in the response",
		func(t *testing.T) {
			// Given a round tripper returning an empty body
			mockRoundTripper := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("{}")),
					}, nil
				},
			}

			// And a spotify token client using it
			tokenClient := SpotifyTokenClient{
				client:      &http.Client{Transport: mockRoundTripper},
				clientId:    "expected-client-id",
				redirectUri: "expected-redirect-uri",
			}

			// When calling GetToken
			_, err := tokenClient.GetToken("a code", "a verifier")

			// Then an error should be returned
			expectedError := "access_token not found or is not a string"
			if err.Error() != expectedError {
				t.Errorf("Expected '%s', but got '%s'", expectedError, err.Error())
			}
		},
	)

	t.Run("it should return an error if the access token is not present in the response",
		func(t *testing.T) {
			// Given a round tripper returning a non-valid json
			mockRoundTripper := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("not valid json"))}, nil
				},
			}

			// And a spotify token client using it
			tokenClient := SpotifyTokenClient{
				client:      &http.Client{Transport: mockRoundTripper},
				clientId:    "expected-client-id",
				redirectUri: "expected-redirect-uri",
			}

			// When calling GetToken
			_, err := tokenClient.GetToken("a code", "a verifier")

			// Then an error should be returned
			expectedError := "failed to unmarshal JSON: invalid character 'o' in literal null (expecting 'u')"
			if err.Error() != expectedError {
				t.Errorf("Expected error to contain '%s', but got '%s'", expectedError, err.Error())
			}
		},
	)
}

func assertPostFormParam(t *testing.T, body url.Values, key string, expected string) {
	if body.Get(key) != expected {
		t.Errorf("Expected form to contain '%s': '%s', but got '%s'", key, expected, body.Get(key))
	}
}
