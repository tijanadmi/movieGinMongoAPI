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

// MovieModel sa CRUD operacijama
// MovieClient is the client responsible for querying MongoDB
type MovieClient struct {
	Col *mongo.Collection
}

var (
	ErrMovieNotFound = errors.New("movie not found")
)

func (c *MovieClient) InitMovies(ctx context.Context) {
	setupIndexes(ctx, c.Col, "title")
}

// AddMovie adds a new movie to the MongoDB collection
func (c *MovieClient) AddMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	movie.ID = primitive.NewObjectID()
	movie.CreatedAt = time.Now()
	result, err := c.Col.InsertOne(ctx, movie)
	if err != nil {
		log.Print(fmt.Errorf("could not add new movie: %w", err))
		return nil, err
	}
	movie.ID = result.InsertedID.(primitive.ObjectID)

	return movie, nil
}

// ListMovies returns all movies from the MongoDB collection
func (c *MovieClient) ListMovies(ctx context.Context) ([]models.Movie, error) {
	movies := make([]models.Movie, 0)
	cur, err := c.Col.Find(ctx, bson.M{})
	if err != nil {
		log.Print(fmt.Errorf("could not get all movies: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &movies); err != nil {
		log.Print(fmt.Errorf("could not marshall the movies results: %w", err))
		return nil, err
	}

	return movies, nil
}

// GetMovie returns a movie by ID from the MongoDB collection
func (c *MovieClient) GetMovie(ctx context.Context, id string) (*models.Movie, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var movie models.Movie
	result := c.Col.FindOne(ctx, bson.M{"_id": objID})
	err = result.Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}

	return &movie, nil
}

// UpdateMovie updates a movie by ID in the MongoDB collection
func (c *MovieClient) UpdateMovie(ctx context.Context, id string, movie models.Movie) (*models.Movie, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, err := c.Col.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{
		{"$set", bson.D{
			{"title", movie.Title},
			{"duration", movie.Duration},
			{"genre", movie.Genre},
			{"directors", movie.Directors},
			{"actors", movie.Actors},
			{"screening", movie.Screening},
			{"plot", movie.Plot},
			{"poster", movie.Poster},
			{"screenings", movie.Screenings},
		}},
	})
	if err != nil {
		log.Print(fmt.Errorf("could not update movie with id [%s]: %w", id, err))
		return &models.Movie{}, err
	}

	if res.MatchedCount == 0 {
		return &models.Movie{}, ErrHallNotFound
	}
	movie.ID = objID

	return &movie, nil
}

// DeleteMovie deletes a movie by ID from the MongoDB collection
func (c *MovieClient) DeleteMovie(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := c.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error deleting the movie with id [%s]: %w", id, err))
		return err
	}

	if res.DeletedCount == 0 {
		return ErrMovieNotFound
	}

	return nil
}

// GetHall returns a hall by ID from the MongoDB collection
func (c *MovieClient) SearchMovies(ctx context.Context, movieId string) ([]models.Movie, error) {
	movies := make([]models.Movie, 0)
	fmt.Println(movieId)
	// Provera inicijalizacije kolekcije
	if c.Col == nil {
		log.Print(fmt.Errorf("collection is not initialized:"))
		return nil, fmt.Errorf("collection is not initialized")
	}

	// Dinamičko kreiranje match stage-a

	var matchStage bson.D
	if movieId != "0" {
		objectId, err := primitive.ObjectIDFromHex(movieId)
		if err != nil {
			log.Print(fmt.Errorf("invalid movie ID: %w", err))
			return nil, err
		}
		matchStage = bson.D{{"$match", bson.D{{"_id", objectId}}}}
		//matchStage = bson.D{{"$match", bson.D{{"_id", movieId}}}}
	} else {
		matchStage = bson.D{{"$match", bson.D{}}}
	}

	pipeline := mongo.Pipeline{
		matchStage,
		{
			{"$lookup", bson.D{
				{"from", "repertoires"},
				{"localField", "_id"},
				{"foreignField", "movieId"},
				{"as", "screenings"},
			}},
		},
		{
			{"$project", bson.D{
				{"_id", 1},
				{"title", 1},
				{"duration", 1},
				{"genre", 1},
				{"directors", 1},
				{"actors", 1},
				{"screening", 1},
				{"plot", 1},
				{"poster", 1},
				{"screenings.date", 1},
				{"screenings.time", 1},
				{"screenings.hall", 1},
			}},
		},
	}

	// Izvršavanje agregacije
	cursor, err := c.Col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(fmt.Errorf("could not aggregate movies: %w", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	// Parsiranje rezultata
	if err := cursor.All(ctx, &movies); err != nil {
		log.Print(fmt.Errorf("could not unmarshal the movies results: %w", err))
		return nil, err
	}

	return movies, nil

}
