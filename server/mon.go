package server

import (
	"wallet-engine/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	trueBool = true
)

var indexes = map[string][]mongo.IndexModel{
	config.WalletCollectioName: {
		{Keys: map[string]interface{}{"userID": 1}, Options: options.Index().SetUnique(true)},
	},
}
