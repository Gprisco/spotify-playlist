package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"prisco.dev/spotify-playlist/client/auth/callback"
)

type CommandExecutor interface {
	executeCommand(string) error
}

type Authenticator struct {
	clientId        string
	redirectUrl     string
	commandExecutor CommandExecutor
	pkceGenerator   PkceGenerator
	callbackHandler callback.CallbackHandler
	credentialStore *Store
}

func NewAuthenticator(
	clientId string,
	redirectUrl string,
	commandExecutor CommandExecutor,
	pkceGenerator PkceGenerator,
	callbackHandler callback.CallbackHandler,
	credentialsStore *Store,
) *Authenticator {
	return &Authenticator{
		clientId,
		redirectUrl,
		commandExecutor,
		pkceGenerator,
		callbackHandler,
		credentialsStore,
	}
}

// Authenticate() starts the OAuth2 authentication flow using PKCE method
func (a *Authenticator) Authenticate() error {
	request, err := a.buildRequest()

	if err != nil {
		return err
	}

	// Open a browser window using the `open` command
	err = a.commandExecutor.executeCommand(fmt.Sprintf("open %s", request.URL.String()))

	if err != nil {
		return errors.New(fmt.Sprintf(
			"Error opening browser window: %s",
			err.Error(),
		))
	}

	// Handle the OAuth 2 callback
	callback := a.callbackHandler(30 * time.Second)

	if callback.Err != "" {
		return errors.New(callback.Err)
	}

	a.credentialStore.Code = callback.Code

	return nil
}

func (a *Authenticator) buildRequest() (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		"https://accounts.spotify.com/authorize",
		nil,
	)

	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Error in creating the http request %s",
			err.Error(),
		))
	}

	q := request.URL.Query()
	q.Add("client_id", a.clientId)
	q.Add("redirect_uri", a.redirectUrl)
	q.Add("response_type", "code")
	q.Add("scope", "user-read-private")
	q.Add("code_challenge_method", "S256")

	// Generate a code verifier using the provided generator
	verifier, err := a.pkceGenerator.GenerateCodeVerifier()

	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Error generating the code verifier: %s",
			err.Error(),
		))
	}

	q.Add("code_challenge", a.pkceGenerator.GenerateCodeChallenge(verifier))

	request.URL.RawQuery = q.Encode()

	return request, nil
}
