package auth

import (
	"encoding/base64"
	"testing"
)

var pkceGenerator RandomPkceGenerator = RandomPkceGenerator{}

func TestPKCE_GenerateCodeVerifier(t *testing.T) {
	// When generating a verifier
	verifier, _ := pkceGenerator.GenerateCodeVerifier()

	// Then ensure the verifier is 43-128 characters long as per PKCE specification.
	verifierLength := len(verifier)
	if verifierLength < 43 || verifierLength > 128 {
		t.Fatalf("Expected code verifier length between 43 and 128, got %d", verifierLength)
	}

	// Ensure the verifier is valid base64url string.
	if _, err := base64.RawURLEncoding.DecodeString(verifier); err != nil {
		t.Fatalf("Expected valid base64url encoding, got %v", err)
	}
}

func TestPKCE_GenerateCodeChallenge(t *testing.T) {
	// Given a verifier
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expectedChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	// When generating the challenge
	challenge := pkceGenerator.GenerateCodeChallenge(verifier)

	if challenge != expectedChallenge {
		t.Fatalf("Expected code challenge %s, got %s", expectedChallenge, challenge)
	}
}

func BenchmarkPKCE_GenerateCodeVerifier(b *testing.B) {
	pkce := &RandomPkceGenerator{}
	for i := 0; i < b.N; i++ {
		_, _ = pkce.GenerateCodeVerifier()
	}
}

func BenchmarkPKCE_GenerateCodeChallenge(b *testing.B) {
	pkce := &RandomPkceGenerator{}
	verifier, _ := pkce.GenerateCodeVerifier()
	for i := 0; i < b.N; i++ {
		_ = pkce.GenerateCodeChallenge(verifier)
	}
}
