package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
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

	firestoreProjectID = ""
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

func firestoreGet(hash string) (string, error) {
	doc, err := firestoreClient.Collection(firestoreCollection).Doc(hash).Get(firestoreContext)
	if err == nil {
		if url, ok := doc.Data()["url"]; ok {
			return url.(string), nil
		}
	}

	return "", err
}

// TODO: Implement Later
func firestorePut() {

}
