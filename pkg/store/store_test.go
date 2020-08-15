package store

import "testing"

func TestStoreEmptyDriver(t *testing.T) {
	exp := "Driver name cannot be empty. Use .Default method instead"
	_, err := New("")

	if err.Error() != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, err)
	}
}

func TestStoreInvalidDriver(t *testing.T) {
	exp := "Driver does not exists"
	_, err := New("nonono")

	if err.Error() != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, err)
	}
}

func TestStoreDefaultDriver(t *testing.T) {
	exp := "memory"
	store, err := Default()
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	if store.driverName != exp {
		t.Errorf("Expected driver \"%v\" and got \"%v\".", exp, store.driverName)
	}
	defer store.Close()
}

func TestStoreAuthEmpty(t *testing.T) {
	store, err := Default()
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	defer store.Close()

	expect := false
	result := store.Auth("")
	if result != expect {
		t.Errorf("Expected %v and got %v.", expect, result)
	}
}

func TestStoreAuth(t *testing.T) {
	store, err := Default()
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	defer store.Close()

	token := "qwert"
	store.AuthToken = token
	expect := true
	result := store.Auth(token)
	if result != expect {
		t.Errorf("Expected %v and got %v.", expect, result)
	}
}

func TestStorePutGet(t *testing.T) {
	store, err := Default()
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	defer store.Close()

	url := "https://github.com/noverde/pequi"
	res, err := store.Put(url)
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}

	val, err := store.Get(res)
	if err != nil {
		t.Errorf("Error \"%v\".", err)
	}
	if val != url {
		t.Errorf("Expected \"%v\" and got \"%v\".", val, url)
	}

}
