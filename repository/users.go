package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tijanadmi/moveginmongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// type Users interface {
// 	GetUserByUsername(ctx context.Context, username string, password string) (*models.User, error)
// }

// BooksClient is the client responsible for querying mongodb
type UsersClient struct {
	Col *mongo.Collection
}

// Get returns a user by username
func (c *UsersClient) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var dbUser models.User

	res := c.Col.FindOne(ctx, bson.M{"username": username})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return &dbUser, res.Err()
		}
		log.Print(fmt.Errorf("error when finding the dbUser [%s]: %q", username, res.Err()))
		return &dbUser, res.Err()
	}

	if err := res.Decode(&dbUser); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", username, err))
		return &dbUser, err
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
	// 	return dbUser, err
	// }

	return &dbUser, nil
}

// AddHall adds a new hall to the MongoDB collection
func (c *UsersClient) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.DateOfCreation = time.Now()
	user.DateOfLastUpdate = time.Now()

	result, err := c.Col.InsertOne(ctx, user)

	if err != nil {
		log.Print(fmt.Errorf("could not add new user: %w", err))
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// getUserByUsername traži korisnika po korisničkom imenu
// func (c *UsersClient) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
//     var user models.User
//     filter := bson.M{"username": username}
//     err := c.Col.FindOne(ctx, filter).Decode(&user)
//     if err == mongo.ErrNoDocuments {
//         return nil, nil // Vraća nil ako korisnik nije pronađen
//     }
//     return &user, err
// }
