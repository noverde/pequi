package store

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type firestoreDriver struct {
	projectID  string
	collection string
	context    context.Context
	client     *firestore.Client
}

func init() {
	Register("firestore", &firestoreDriver{})
}

func (d *firestoreDriver) init() {
	var err error

	log.Printf("Initializing Firestore storage engine")

	d.projectID = "*detect-project-id*"
	if value, ok := os.LookupEnv("FIRESTORE_PROJECT"); ok {
		d.projectID = value
	}
	log.Printf("Firestore project: %s", d.projectID)

	d.collection = "urls"
	if value, ok := os.LookupEnv("FIRESTORE_COLLECTION"); ok {
		d.collection = value
	}
	log.Printf("Firestore Collection: %s", d.collection)

	d.context = context.Background()
	d.client, err = firestore.NewClient(d.context, d.projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
}

func (d *firestoreDriver) close() {
	d.client.Close()
}

func (d *firestoreDriver) auth(token string) bool {
	return false
}

func (d *firestoreDriver) get(slug string) (string, error) {
	doc, err := d.client.Collection(d.collection).Doc(slug).Get(d.context)
	if err == nil {
		if url, ok := doc.Data()["long_url"]; ok {
			return url.(string), nil
		}
	}

	return "", err
}

func (d *firestoreDriver) set(slug string, url string) error {
	// Reference collection item and run a transaction to prevent duplication
	ref := d.client.Collection(d.collection).Doc(slug)
	err := d.client.RunTransaction(d.context,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Try to get item
			_, err := tx.Get(ref)
			if err == nil {
				return errors.New("Item already exists")
			}

			// We can proceed only if error is notfound
			if status.Code(err) != codes.NotFound {
				return err
			}

			// Save data
			return tx.Set(ref, map[string]interface{}{
				"long_url":   url,
				"created_at": firestore.ServerTimestamp,
			})
		})

	return err
}
