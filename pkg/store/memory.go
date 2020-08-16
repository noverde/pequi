package store

import (
	"errors"
	"log"
	"sync"
	"time"
)

type memoryDriver struct {
	mutex sync.RWMutex
	data  map[string]memoryData
}

type memoryData struct {
	url       string
	createdAt time.Time
}

func init() {
	Register("memory", &memoryDriver{})
}

func (d *memoryDriver) init() {
	log.Printf("Initializing memory storage engine")
	d.data = make(map[string]memoryData)
}

func (d *memoryDriver) close() {
	d.data = nil
}

func (d *memoryDriver) auth(token string) bool {
	return false
}

func (d *memoryDriver) get(slug string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	value, exists := d.data[slug]
	if exists {
		return value.url, nil
	}

	return "", errors.New("Item does not exists")
}

func (d *memoryDriver) set(slug string, url string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.data[slug]; exists {
		return errors.New("Item already exists")
	}

	d.data[slug] = memoryData{
		url:       url,
		createdAt: time.Now(),
	}

	return nil
}
