package main

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	firestoreProjectID  string
	firestoreCollection string
	firestoreContext    context.Context
	firestoreClient     *firestore.Client
)

func firestoreInit() {
	var err error

	log.Printf("Initializing Firestore storage engine")

	firestoreProjectID = "*detect-project-id*"
	if value, ok := os.LookupEnv("FIRESTORE_PROJECT"); ok {
		firestoreProjectID = value
	}
	log.Printf("Firestore project: %s", firestoreProjectID)

	firestoreCollection = "urls"
	if value, ok := os.LookupEnv("FIRESTORE_COLLECTION"); ok {
		firestoreCollection = value
	}
	log.Printf("Firestore Collection: %s", firestoreCollection)

	firestoreContext = context.Background()
	firestoreClient, err = firestore.NewClient(firestoreContext, firestoreProjectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
}

func firestoreClose() {
	firestoreClient.Close()
}

func firestoreAuth(token string) bool {
	return false
}

func firestoreGet(hash string) (string, error) {
	doc, err := firestoreClient.Collection(firestoreCollection).Doc(hash).Get(firestoreContext)
	if err == nil {
		if url, ok := doc.Data()["url"]; ok {
			return url.(string), nil
		}
	}

	return "", err
}

func firestorePut(slug string, url string) error {
	// Reference collection item and run a transaction to prevent duplication
	ref := firestoreClient.Collection(firestoreCollection).Doc(slug)
	err := firestoreClient.RunTransaction(firestoreContext,
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
				"url":        url,
				"created_at": firestore.ServerTimestamp,
			})
		})

	return err
}
