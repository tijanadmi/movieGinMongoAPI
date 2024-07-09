package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tijanadmi/movieginmongoapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRepertoireNotFound = errors.New("repertoire not found")
)

// AddRepertoire adds a new repertoire to the MongoDB collection
func (r *MongoStore) AddRepertoire(ctx context.Context, repertoire *models.Repertoire) (*models.Repertoire, error) {
	repertoire.ID = primitive.NewObjectID()
	repertoire.CreatedAt = time.Now()
	// Provera da li je numOfResTickets postavljen, ako nije postavi na 0
	if repertoire.NumOfResTickets == 0 {
		repertoire.NumOfResTickets = 0
	}
	fmt.Println("Repository NumOfResTickets", repertoire.NumOfResTickets)
	result, err := r.db.Collection("repertoires").InsertOne(ctx, repertoire)
	if err != nil {
		log.Print(fmt.Errorf("could not add new repertoire: %w", err))
		return nil, err
	}

	repertoire.ID = result.InsertedID.(primitive.ObjectID)
	return repertoire, nil
}

// ListRepertoires returns all repertoires from the MongoDB collection
func (r *MongoStore) ListRepertoires(ctx context.Context) ([]models.Repertoire, error) {
	repertoires := make([]models.Repertoire, 0)
	cur, err := r.db.Collection("repertoires").Find(ctx, bson.M{})
	if err != nil {
		log.Print(fmt.Errorf("could not get all repertoires: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &repertoires); err != nil {
		log.Print(fmt.Errorf("could marshall the repertoires results: %w", err))
		return nil, err
	}

	return repertoires, nil
}

// GetRepertoire returns a repertoire based on its ID
func (r *MongoStore) GetRepertoire(ctx context.Context, id string) (*models.Repertoire, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var repertoire models.Repertoire
	result := r.db.Collection("repertoires").FindOne(ctx, bson.M{"_id": objID})
	err = result.Decode(&repertoire)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrRepertoireNotFound
		}
		return nil, err
	}

	return &repertoire, nil
}

// GetRepertoire returns a repertoire based on its ID
func (r *MongoStore) GetRepertoireByMovieDateTimeHall(ctx context.Context, movieId string, dateValue time.Time, timeValue string, hallValue string) (models.Repertoire, error) {
	var repertoire models.Repertoire
	movieID, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{
		"movieId": movieID,
		"date":    dateValue,
		"time":    timeValue,
		"hall":    hallValue,
	}
	res := r.db.Collection("repertoires").FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return repertoire, nil
		}
		log.Print(fmt.Errorf("error when finding the repertoire [%s]: %q", movieID, res.Err()))
		return repertoire, res.Err()
	}

	if err := res.Decode(&repertoire); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", movieID, err))
		return repertoire, err
	}
	return repertoire, nil
}

// GetRepertoire returns a repertoires based on its movieId
func (r *MongoStore) GetAllRepertoireForMovie(ctx context.Context, movieId string, startDate time.Time, endDate time.Time) ([]models.Repertoire, error) {
	repertoires := make([]models.Repertoire, 0)

	movieID, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{
		"movieId": movieID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	cur, err := r.db.Collection("repertoires").Find(ctx, filter)
	if err != nil {
		log.Print(fmt.Errorf("could not get repertoires for period  %s - %s [%s]: %w", startDate, endDate, movieId, err))
		return nil, err
	}

	if err = cur.All(ctx, &repertoires); err != nil {
		log.Print(fmt.Errorf("could marshall the repertoires results: %w", err))
		return nil, err
	}

	return repertoires, nil
}

// UpdateRepertoire updates a repertoire based on its ID
func (r *MongoStore) UpdateRepertoire(ctx context.Context, id string, repertoire models.Repertoire) (*models.Repertoire, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Repertoire{}, err
	}
	res, err := r.db.Collection("repertoires").UpdateOne(ctx, bson.M{"_id": objID}, bson.D{
		{"$set", bson.D{
			{"movieId", repertoire.MovieID},
			{"date", repertoire.Date},
			{"time", repertoire.Time},
			{"hall", repertoire.Hall},
			{"numOfTickets", repertoire.NumOfTickets},
			{"numOfResTickets", repertoire.NumOfResTickets},
			{"reservSeats", repertoire.ReservSeats},
		}},
	})
	if err != nil {
		log.Print(fmt.Errorf("could not update repertoire with id [%s]: %w", id, err))
		return &models.Repertoire{}, err
	}

	if res.MatchedCount == 0 {
		return &models.Repertoire{}, ErrRepertoireNotFound
	}
	repertoire.ID = objID

	return &repertoire, nil
}

// DeleteRepertoire deletes a repertoire based on its ID
func (r *MongoStore) DeleteRepertoire(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.db.Collection("repertoires").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the repertoire with id [%s]: %w", id, err))
		return err
	}

	if res.DeletedCount == 0 {
		return ErrRepertoireNotFound
	}

	return nil
}

// DeleteRepertoire deletes a repertoire based on its ID
func (r *MongoStore) DeleteRepertoireForMovie(ctx context.Context, movieId string) error {
	movieID, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		return err
	}

	res, err := r.db.Collection("repertoires").DeleteMany(ctx, bson.M{"movieId": movieID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the repertoire with id [%s]: %w", movieId, err))
		return err
	}

	if res.DeletedCount == 0 {
		return ErrRepertoireNotFound
	}

	return nil
}
