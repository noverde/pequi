package store

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
	set(string, string) error
}

var drivers = make(map[string]storeDriver)

// Store is the database storage driver.
type Store struct {
	driver     storeDriver
	driverName string
	AuthToken  string
}

// Register makes a store driver available by the provided name.
func Register(name string, driver storeDriver) {
	drivers[name] = driver
}

// New creates a new instance of Store object
func New(driverName string) (*Store, error) {
	if driverName == "" {
		return nil, errors.New("Driver name cannot be empty. Use .Default method instead")
	}

	driver, exists := drivers[driverName]
	if !exists {
		return nil, errors.New("Driver does not exists")
	}

	s := &Store{
		driver:     driver,
		driverName: driverName,
		AuthToken:  os.Getenv("AUTHORIZATION_TOKEN"),
	}
	s.driver.init()

	return s, nil
}

// Default ...
func Default() (*Store, error) {
	return New("memory")
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

// Set adds URL to datastore
func (s *Store) Set(url string) (string, error) {
	// If fails or duplicate, try to generate N (MAX_RETRIES) times
	for count := 0; count < storeMaxRetries; count++ {
		slug, err := shortid.Generate()
		if err != nil {
			log.Printf("Shortid generation failed (%d)", count)
			continue
		}
		if err = s.driver.set(slug, url); err == nil {
			return slug, nil
		}
		log.Printf("Slug store failed (%d)", count)
	}
	return "", errors.New("Error trying to generate unique slug")
}
