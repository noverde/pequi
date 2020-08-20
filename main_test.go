package main

import "testing"

func TestNormalizeAddress(t *testing.T) {
	exp := "http://localhost/"
	addr := normalizeURLString("")
	if addr != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, addr)
	}

	exp = "http://localhost/"
	addr = normalizeURLString("localhost")
	if addr != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, addr)
	}

	exp = "http://mydomain.com:8080/"
	addr = normalizeURLString("mydomain.com:8080")
	if addr != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, addr)
	}

	exp = "http://mydomain.com/"
	addr = normalizeURLString("http://mydomain.com")
	if addr != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, addr)
	}

	exp = "https://mydomain.com/"
	addr = normalizeURLString("https://mydomain.com")
	if addr != exp {
		t.Errorf("Expected error \"%v\" and got \"%v\".", exp, addr)
	}
}
