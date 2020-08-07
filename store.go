package main

import (
	"os"

	"github.com/teris-io/shortid"
)

var storeAuthToken = ""

func storeInit() {
	if value, ok := os.LookupEnv("AUTHORIZATION_TOKEN"); ok {
		storeAuthToken = value
	}

	firestoreInit()
}

func storeClose() {
	firestoreClose()
}

func storeAuthorization(token string) bool {
	if storeAuthToken != "" && storeAuthToken == token {
		return true
	}

	return firestoreAuthorization(token)
}

func storeGet(hash string) (string, error) {
	return firestoreGet(hash)
}

// TODO: Implement Later
func storePut(data *payload) (string, error) {
	return shortid.Generate()
}
