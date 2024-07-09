package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStore(client *mongo.Client, dbName string) *MongoStore {
	database := client.Database(dbName)
	return &MongoStore{
		client: client,
		db:     database,
	}
}

