package configs

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client instance.
var _DBClient, _DBClientDisconnect = ConnectDB()

func ConnectDB() (*mongo.Client, func()) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(Env.MongodbURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	return client, func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func GetDB() *mongo.Database {
	return _DBClient.Database(Env.DatabaseName)
}

// getting database collections.
func GetCollection(collectionName string) *mongo.Collection {
	db := GetDB()
	collection := db.Collection(collectionName)
	return collection
}

func CloseDB() {
	if _DBClientDisconnect != nil {
		_DBClientDisconnect()
	}
}
