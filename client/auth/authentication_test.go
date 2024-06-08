package auth

import (
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
	pkce string
}

func (m MockPkceGenerator) generate() string {
	return m.pkce
}

func TestAuthenticator(t *testing.T) {
	t.Run("it should open a browser window to the correct URL",
		func(t *testing.T) {
			// Given a pkce generator
			pkceGenerator := MockPkceGenerator{"pkce"}

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
			authenticator.Authenticate()

			// Then it should match the expected command
			if r := recover(); r != nil {
				t.Error() // Message will be printed by the panic
			}
		},
	)
}
