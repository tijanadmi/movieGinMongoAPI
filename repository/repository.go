package repository

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func getCollection(client *mongo.Client, col string) *mongo.Collection {

	return client.Database("userDB").Collection(col)
}

func setupIndexes(ctx context.Context, collection *mongo.Collection, key string) {
	idxOpt := &options.IndexOptions{}
	idxOpt.SetUnique(true)
	mod := mongo.IndexModel{
		Keys: bson.M{
			key: 1, // index in ascending order
		},
		Options: idxOpt,
	}

	ind, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal(fmt.Errorf("Indexes().CreateOne() ERROR: %w", err))
	} else {
		// BooksHandler call returns string of the index name
		log.Printf("CreateOne() index: %s\n", ind)
	}
}
