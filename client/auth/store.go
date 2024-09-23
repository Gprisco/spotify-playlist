package auth

import "sync"

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
