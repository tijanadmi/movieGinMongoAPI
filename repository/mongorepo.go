package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)



// MongoClient combines all collection clients
type MongoClient struct {
	Hall        HallClient
	Movie       MovieClient
	Repertoire  RepertoireClient
	Reservation ReservationClient
	Users       UsersClient
	Client      *mongo.Client // polje za čuvanje klijenta MongoDB
}

// NewMongoClient initializes the MongoDB clients and sets up their collections
func NewMongoClient(client *mongo.Client) *MongoClient {
	return &MongoClient{
		Hall: HallClient{
			Col: getCollection(client, "halls"),
		},
		Movie: MovieClient{
			Col: getCollection(client, "movies"),
		},
		Repertoire: RepertoireClient{
			Col: getCollection(client, "repertoires"),
		},
		Reservation: ReservationClient{
			Col: getCollection(client, "reservations"),
		},
		Users: UsersClient{
			Col: getCollection(client, "users"),
		},
		Client: client, // Čuvanje klijenta MongoDB za upravljanje sesijom
	}
}
