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
