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

var drivers = make(map[string]storeDriver)

// Store is the database storage driver.
type Store struct {
	driver    storeDriver
	AuthToken string
}

// Register makes a store driver available by the provided name.
func Register(name string, driver storeDriver) {
	drivers[name] = driver
}

// NewStore creates a new instance of Store object
func NewStore() *Store {
	s := new(Store)
	s.AuthToken = os.Getenv("AUTHORIZATION_TOKEN")

	// Get storage driver or use firestore as default.
	name := os.Getenv("STORAGE_DRIVER")
	if name == "" {
		name = "firestore"
	}
	log.Printf("Storage driver: %s", name)

	// Initialize storage driver.
	s.driver = drivers[name]
	s.driver.init()

	return s
}

// Close cleanup driver connection
func (s *Store) Close() {
	s.driver.close()
}

// Auth check for valid authentication token
func (s *Store) Auth(token string) bool {
	if s.AuthToken != "" && s.AuthToken == token {
		return true
	}

	return s.driver.auth(token)
}

// Get item from datastore
func (s *Store) Get(hash string) (string, error) {
	return s.driver.get(hash)
}

// Put adds URL to datastore
func (s *Store) Put(url string) (string, error) {
	// If fails or duplicate, try to generate N (MAX_RETRIES) times
	for count := 0; count < storeMaxRetries; count++ {
		slug, err := shortid.Generate()
		if err != nil {
			log.Printf("Shortid generation failed (%d)", count)
			continue
		}
		if err = s.driver.put(slug, url); err == nil {
			return slug, nil
		}
		log.Printf("Slug store failed (%d)", count)
	}
	return "", errors.New("Error trying to generate unique slug")
}
