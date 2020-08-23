package store

import (
	"os"
	"testing"
)

func TestFirestoreAuth(t *testing.T) {
	os.Setenv("FIRESTORE_PROJECT", "test")
	os.Setenv("FIRESTORE_COLLECTION", "test")

	store := firestoreDriver{}
	store.init()
	defer store.close()

	expect := false
	result := store.auth("")
	if result != expect {
		t.Errorf("Expected %v and got %v.", expect, result)
	}
}
