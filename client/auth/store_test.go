package auth

import (
	"testing"
)

func TestSingleton(t *testing.T) {
	t.Run("it should create one and only one instance",
		func(t *testing.T) {
			// given 2 instances of a store
			store1 := GetCredentialStore()
			store2 := GetCredentialStore()

			// expect them to be the same
			if store1 != store2 {
				t.Error("The instances are different")
			}
		})
}

func TestStore(t *testing.T) {
	t.Run("it should store the code only once", func(t *testing.T) {
		// Given a store
		store := GetCredentialStore()

		// When setting a code
		store.SetCode("test")

		// and trying to overwrite it
		err := store.SetCode("not_allowed")

		// Then it should have the correct code
		if store.GetCode() != "test" {
			t.Errorf("Expected %v, but got %v", "test", store.GetCode())
		}

		// and err should not be null
		if err == nil {
			t.Error("it should not be allowed to overwrite the code")
		}
	})

	t.Run("it should store the token only once", func(t *testing.T) {
		// Given a store
		store := GetCredentialStore()

		// When setting a token
		store.SetToken("test_token")

		// and trying to overwrite it
		err := store.SetToken("not_allowed_token")

		// Then it should have the correct code
		if store.GetToken() != "test_token" {
			t.Errorf("Expected %v, but got %v", "test_token", store.GetToken())
		}

		// and err should not be null
		if err == nil {
			t.Error("it should not be allowed to overwrite the token")
		}
	})
}
