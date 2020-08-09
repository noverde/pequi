package main

import (
	"errors"
	"log"
	"os"

	"github.com/teris-io/shortid"
)

const storeMaxRetries = 10

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

func storePut(url string) (string, error) {
	// If fails or duplicate, try to generate N (MAX_RETRIES) times
	for count := 0; count < storeMaxRetries; count++ {
		slug, err := shortid.Generate()
		if err != nil {
			log.Printf("Shortid generation failed (%d)", count)
			continue
		}
		if err = firestorePut(slug, url); err == nil {
			return slug, nil
		}
		log.Printf("Slug store failed (%d)", count)
	}
	return "", errors.New("Error trying to generate unique slug")
}
