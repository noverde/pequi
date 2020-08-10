package main

import (
	"errors"
	"log"
	"os"

	"github.com/teris-io/shortid"
)

const storeMaxRetries = 5

type storeDriver interface {
	init()
	close()
	auth(string) bool
	get(string) (string, error)
	put(string, string) error
}

var (
	driver         storeDriver
	storeAuthToken string
)

func storeInit() {
	storeAuthToken = os.Getenv("AUTHORIZATION_TOKEN")

	driver = &firestoreDriver{}
	driver.init()
}

func storeClose() {
	driver.close()
}

func storeAuth(token string) bool {
	if storeAuthToken != "" && storeAuthToken == token {
		return true
	}

	return driver.auth(token)
}

func storeGet(hash string) (string, error) {
	return driver.get(hash)
}

func storePut(url string) (string, error) {
	// If fails or duplicate, try to generate N (MAX_RETRIES) times
	for count := 0; count < storeMaxRetries; count++ {
		slug, err := shortid.Generate()
		if err != nil {
			log.Printf("Shortid generation failed (%d)", count)
			continue
		}
		if err = driver.put(slug, url); err == nil {
			return slug, nil
		}
		log.Printf("Slug store failed (%d)", count)
	}
	return "", errors.New("Error trying to generate unique slug")
}
