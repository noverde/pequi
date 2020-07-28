package main

import "github.com/teris-io/shortid"

func storeInit() {
	firestoreInit()
}

func storeClose() {
	firestoreClose()
}

func storeGet(hash string) (string, error) {
	return firestoreGet(hash)
}

// TODO: Implement Later
func storePut(data *payload) (string, error) {
	return shortid.Generate()
}
