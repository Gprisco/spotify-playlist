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
		panic(fmt.Sprintf("Expected\t%s\ngot\t\t%s", m.expectedCommand, command))
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

			// Then it should match the expected command
			defer func() {
				if r := recover(); r != nil {
					t.Error() // Message will be printed by the panic
				}
			}()

			// When starting the authentication flow
			authenticator.Authenticate()
		},
	)

	t.Run("it should panic when pkce generator returns an error",
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

			// Then it should recover from a panic
			defer func() {
				if r := recover(); r == nil {
					t.Error()
				}
			}()

			// When starting the authentication flow
			authenticator.Authenticate()
		},
	)

	t.Run("it should panic when the command executor returns an error",
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

			// Then it should recover from a panic
			defer func() {
				if r := recover(); r == nil {
					t.Error()
				}
			}()

			// When starting the authentication flow
			authenticator.Authenticate()
		},
	)
}
