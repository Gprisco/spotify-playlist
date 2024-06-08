package auth

import (
	"fmt"
	"net/http"
)

type CommandExecutor interface {
	executeCommand(string) error
}

type Authenticator struct {
	clientId        string
	redirectUrl     string
	commandExecutor CommandExecutor
	pkceGenerator   PkceGenerator
}

func NewAuthenticator(
	clientId string,
	redirectUrl string,
	commandExecutor CommandExecutor,
	pkceGenerator PkceGenerator,
) *Authenticator {
	return &Authenticator{
		clientId,
		redirectUrl,
		commandExecutor,
		pkceGenerator,
	}
}

func (a *Authenticator) Authenticate() {
	request, err := http.NewRequest(
		http.MethodGet,
		"https://accounts.spotify.com/authorize",
		nil,
	)

	if err != nil {
		panic(err)
	}

	q := request.URL.Query()
	q.Add("client_id", a.clientId)
	q.Add("redirect_uri", a.redirectUrl)
	q.Add("response_type", "code")
	q.Add("scope", "user-read-private")
	q.Add("code_challenge_method", "S256")

	verifier, err := a.pkceGenerator.GenerateCodeVerifier()

	if err != nil {
		panic(err)
	}

	q.Add("code_challenge", a.pkceGenerator.GenerateCodeChallenge(verifier))

	request.URL.RawQuery = q.Encode()

	err = a.commandExecutor.executeCommand(fmt.Sprintf("open %s", request.URL.String()))

	if err != nil {
		panic(err)
	}
}
