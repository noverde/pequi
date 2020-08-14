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
