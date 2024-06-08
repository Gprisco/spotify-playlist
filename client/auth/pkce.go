package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type PkceGenerator interface {
	GenerateCodeVerifier() (string, error)
	GenerateCodeChallenge(verifier string) string
}

type RandomPkceGenerator struct{}

func (p *RandomPkceGenerator) GenerateCodeVerifier() (string, error) {
	const codeVerifierLength = 32
	verifier := make([]byte, codeVerifierLength)
	_, err := rand.Read(verifier)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(verifier), nil
}

func (p *RandomPkceGenerator) GenerateCodeChallenge(verifier string) string {
	sha := sha256.New()
	sha.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sha.Sum(nil))
}
