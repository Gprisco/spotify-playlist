package auth

import (
	"fmt"
	"net/http"
)

type CommandExecutor interface {
	executeCommand(string) error
}

type PkceGenerator interface {
	generate() string
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
	q.Add("code_challenge", a.pkceGenerator.generate())

	request.URL.RawQuery = q.Encode()

	a.commandExecutor.executeCommand(fmt.Sprintf("open %s", request.URL.String()))
}
