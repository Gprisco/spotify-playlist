package auth

import (
	"errors"
	"sync"
)

var lock = &sync.Mutex{}

type Store struct {
	code  string
	token string
}

var storeSingleton *Store

func GetCredentialStore() *Store {
	if storeSingleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if storeSingleton == nil {
			storeSingleton = &Store{}
		}
	}

	return storeSingleton
}

func (s *Store) SetCode(code string) error {
	if len(s.code) > 0 {
		return errors.New("An authorization code was already present")
	}

	s.code = code
	return nil
}

func (s *Store) GetCode() string {
	return s.code
}

func (s *Store) SetToken(token string) error {
	if len(s.token) > 0 {
		return errors.New("An authorization code was already present")
	}

	s.token = token
	return nil
}

func (s *Store) GetToken() string {
	return s.token
}
