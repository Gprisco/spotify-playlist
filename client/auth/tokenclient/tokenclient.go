package tokenclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpotifyTokenClient struct {
	client      *http.Client
	clientId    string
	redirectUri string
}

type TokenClient interface {
	GetToken(
		code string,
		codeVerifier string,
	) (string, error)
}

func (s *SpotifyTokenClient) GetToken(
	code string,
	codeVerifier string,
) (string, error) {
	const endpoint = "https://accounts.spotify.com/api/token"

	reqBody := url.Values{}
	reqBody.Add("grant_type", "authorization_code")
	reqBody.Add("code", code)
	reqBody.Add("redirect_uri", s.redirectUri)
	reqBody.Add("client_id", s.clientId)
	reqBody.Add("code_verifier", codeVerifier)

	resp, err := s.client.PostForm(endpoint, reqBody)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Extract the access_token field
	accessToken, ok := responseMap["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found or is not a string")
	}

	return accessToken, nil
}
