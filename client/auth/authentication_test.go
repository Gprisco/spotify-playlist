package auth

import (
	"errors"
	"fmt"
	"testing"
)

// Mock Command Executor
type MockCommandExecutor struct {
	expectedCommand string
	errorReturned   error
}

func (m MockCommandExecutor) executeCommand(command string) error {
	if command != m.expectedCommand {
		return errors.New(fmt.Sprintf("Expected\t%s\ngot\t\t%s", m.expectedCommand, command))
	}

	return m.errorReturned
}

// Mock PKCE Generator
type MockPkceGenerator struct {
	pkce     string
	verifier string

	verifierError error
}

func (m MockPkceGenerator) GenerateCodeVerifier() (string, error) {
	return m.verifier, m.verifierError
}

func (m MockPkceGenerator) GenerateCodeChallenge(verifier string) string {
	if m.verifier != verifier {
		panic("It should use the generated verifier")
	}

	return m.pkce
}

func TestAuthenticator(t *testing.T) {
	t.Run("it should open a browser window to the correct URL",
		func(t *testing.T) {
			// Given a pkce generator
			pkceGenerator := MockPkceGenerator{"pkce", "verifier", nil}

			// and a command executor expecting the correct command
			successfulCommandExecutor := MockCommandExecutor{
				"open " +
					"https://accounts.spotify.com/authorize?" +
					"client_id=clientId&" +
					"code_challenge=pkce&" +
					"code_challenge_method=S256&" +
					"redirect_uri=redirectUrl&" +
					"response_type=code&" +
					"scope=user-read-private",
				nil,
			}

			// and an authenticator using it
			authenticator := NewAuthenticator(
				"clientId",
				"redirectUrl",
				successfulCommandExecutor,
				pkceGenerator,
			)

			// When starting the authentication flow
			err := authenticator.Authenticate()

			if err != nil {
				t.Errorf("The authentication went wrong: %s", err.Error())
			}
		},
	)

	t.Run("it should return an error when pkce generator returns an error",
		func(t *testing.T) {
			// Given a pkce generator which returns an error
			pkceGenerator := MockPkceGenerator{
				"ignored",
				"ignored",
				errors.New("Error generating the verifier"),
			}

			// and a command executor
			successfulCommandExecutor := MockCommandExecutor{
				"ignored",
				nil,
			}

			// and an authenticator using it
			authenticator := NewAuthenticator(
				"clientId",
				"redirectUrl",
				successfulCommandExecutor,
				pkceGenerator,
			)

			// When starting the authentication flow
			err := authenticator.Authenticate()

			if err == nil {
				t.Errorf("The authentication did not return an error as expected")
			}
		},
	)

	t.Run("it should return an error when the command executor returns an error",
		func(t *testing.T) {
			// Given a pkce generator which returns no error
			pkceGenerator := MockPkceGenerator{
				"ignored",
				"ignored",
				nil,
			}

			// and a command executor which returns an error
			successfulCommandExecutor := MockCommandExecutor{
				"ignored",
				errors.New("Command execution error"),
			}

			// and an authenticator using it
			authenticator := NewAuthenticator(
				"clientId",
				"redirectUrl",
				successfulCommandExecutor,
				pkceGenerator,
			)

			// When starting the authentication flow
			err := authenticator.Authenticate()

			if err == nil {
				t.Errorf("The authentication did not return an error as expected")
			}
		},
	)
}
