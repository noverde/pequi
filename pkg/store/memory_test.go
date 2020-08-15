package store

import "testing"

func TestMemoryStoreAuth(t *testing.T) {
	store := memoryDriver{}
	store.init()
	defer store.close()

	expect := false
	result := store.auth("")
	if result != expect {
		t.Errorf("Expected %v and got %v.", expect, result)
	}
}

func TestMemoryStorePutGet(t *testing.T) {
	store := memoryDriver{}
	store.init()
	defer store.close()

	slug := "qwert"
	url := "https://github.com/noverde/pequi"
	err := store.put(slug, url)
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}

	val, err := store.get(slug)
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	if val != url {
		t.Errorf("Expected \"%v\" and got \"%v\".", val, url)
	}
}

func TestMemoryStoreItemDoesNotExists(t *testing.T) {
	store := memoryDriver{}
	store.init()
	defer store.close()

	exp := "Item does not exists"
	_, err := store.get("qwert")
	if err.Error() != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, err)
	}
}

func TestMemoryStoreItemAlreadyExists(t *testing.T) {
	store := memoryDriver{}
	store.init()
	defer store.close()

	exp := "Item already exists"
	url := "https://github.com/noverde/pequi"
	slug := "qwert"

	err := store.put(slug, url)
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	err = store.put(slug, url)
	if err.Error() != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, err)
	}
}
