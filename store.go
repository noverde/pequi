package main

import (
	"os"

	"github.com/teris-io/shortid"
)

var storeAuthToken string

func storeInit() {
	storeAuthToken = os.Getenv("AUTHORIZATION_TOKEN")
	firestoreInit()
}

func storeClose() {
	firestoreClose()
}

func storeAuth(token string) bool {
	if storeAuthToken != "" && storeAuthToken == token {
		return true
	}

	return firestoreAuth(token)
}

func storeGet(hash string) (string, error) {
	return firestoreGet(hash)
}

// TODO: Implement Later
func storePut(data *payload) (string, error) {
	return shortid.Generate()
}
