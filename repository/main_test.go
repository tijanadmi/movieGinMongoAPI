package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tijanadmi/moveginmongo/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testStore *MongoClient

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	clientOptions := options.Client().ApplyURI(config.MongoURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("cannot connect to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("cannot ping MongoDB:", err)
	}

	testStore = NewMongoClient(client)
	os.Exit(m.Run())
}
