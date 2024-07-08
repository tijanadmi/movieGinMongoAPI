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

// RepertoireClient is the client responsible for querying mongodb
type ReservationClient struct {
	Col *mongo.Collection
}

var (
	ErrReservationNotFound = errors.New("reservation not found")
)

// AddReservation adds a new reservation to the MongoDB collection
func (c *ReservationClient) InsertReservation(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error) {
	reservation.ID = primitive.NewObjectID()
	reservation.CreationDate = time.Now()
	result, err := c.Col.InsertOne(ctx, reservation)
	if err != nil {
		log.Print(fmt.Errorf("could not add new reservation: %w", err))
		return nil, err
	}
	reservation.ID = result.InsertedID.(primitive.ObjectID)
	return reservation, nil
}

// GetReservationById returns a reservations based on its ID
func (c *ReservationClient) GetReservationById(ctx context.Context, id string) (*models.Reservation, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var reservation models.Reservation
	result := c.Col.FindOne(ctx, bson.M{"_id": objID})
	err = result.Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	return &reservation, nil
}

// GetReservation returns a all reservation based on username
func (c *ReservationClient) GetAllReservationsForUser(ctx context.Context, username string) ([]models.Reservation, error) {
	reservations := make([]models.Reservation, 0)

	cur, err := c.Col.Find(ctx, bson.M{"username": username})
	if err != nil {
		log.Print(fmt.Errorf("could not get all reservations [%s]: %w", username, err))
		return nil, err
	}

	if err = cur.All(ctx, &reservations); err != nil {
		log.Print(fmt.Errorf("could marshall the repertoires results: %w", err))
		return nil, err
	}

	return reservations, nil
}

// DeleteReservation deletes a reservation based on its ID
func (c *ReservationClient) DeleteReservation(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := c.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the repertoire with id [%s]: %w", id, err))
		return err
	}
	if res.DeletedCount == 0 {
		return ErrReservationNotFound
	}

	return nil
}
