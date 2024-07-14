package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewStore(client *mongo.Client, dbName string) Store {
	database := client.Database(dbName)
	return &MongoStore{
		client: client,
		db:     database,
	}
}
