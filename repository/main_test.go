package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tijanadmi/movieginmongoapi/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testStore Store

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

	// create a repository

	testStore = NewStore(client, config.Database)
	os.Exit(m.Run())
}
