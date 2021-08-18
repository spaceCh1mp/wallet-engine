package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect takes a url string and returns the specified database name
func Connect(url, dbName string) *mongo.Database {

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalf("could not fetch client from url [ %s ]: %v", url, err)
	}

	if err := client.Connect(context.Background()); err != nil {
		log.Fatalf("could not initialize client: %v", err)
	}

	return client.Database(dbName)
}
