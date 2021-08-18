package mongo

import (
	"context"
	"wallet-engine/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollIndexes binds a slice of mongo.IndexModel to a collection specified as `string`
type CollIndexes map[string][]mongo.IndexModel

// EnforceValidation creates the collection if it doesn't already exist with the validation rule specified.
// If the collection already exists it returns nil
func EnforceValidation(db *mongo.Database) error {
	colls, err := db.ListCollectionNames(context.TODO(), bson.M{"name": config.WalletCollectioName})
	if err != nil {
		return err
	}

	if len(colls) < 1 {
		return db.CreateCollection(context.TODO(), config.WalletCollectioName, &options.CreateCollectionOptions{
			Validator: bson.M{
				"$jsonSchema": bson.M{
					"properties": bson.M{
						"balance": bson.M{
							"minimum":     0,
							"description": "balance cannot be a negative value",
						},
					},
				},
			},
		})
	}

	return nil
}

// CreateIndexes creates the indexes for a collection if it doesn't already exist.
func CreateIndexes(db *mongo.Database, indexes CollIndexes) error {
	for coll, indexes := range indexes {

		_, err := db.Collection(coll).Indexes().CreateMany(context.TODO(), indexes)
		if err != nil {
			return err
		}

	}

	return nil
}
