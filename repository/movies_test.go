package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
)

func createRandomMovie(t *testing.T) *models.Movie {

	arg := models.Movie{
		Title:     util.RandomString(50),
		Duration:  int32(util.RandomInt(100, 250)),
		Genre:     util.RandomString(200),
		Directors: util.RandomString(200),
		Actors:    util.RandomString(200),
		Screening: time.Now(),
		Plot:      util.RandomString(200),
		Poster:    util.RandomString(200),
	}

	movie, err := testStore.Movie.AddMovie(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, movie)

	require.Equal(t, arg.Title, movie.Title)
	require.Equal(t, arg.Duration, movie.Duration)
	require.Equal(t, arg.Genre, movie.Genre)
	require.Equal(t, arg.Directors, movie.Directors)
	require.Equal(t, arg.Actors, movie.Actors)
	require.Equal(t, arg.Screening, movie.Screening)
	require.Equal(t, arg.Plot, movie.Plot)
	require.Equal(t, arg.Poster, movie.Poster)

	require.NotZero(t, movie.ID)
	require.NotZero(t, movie.CreatedAt)

	return movie
}

func TestCreateMovie(t *testing.T) {
	createRandomMovie(t)
}

func TestListMovie(t *testing.T) {

	movies, err := testStore.Movie.ListMovies(context.Background())
	require.NoError(t, err)

	createRandomMovie(t)

	movies1, err := testStore.Movie.ListMovies(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, movies1)

	require.Equal(t, len(movies)+1, len(movies1))
}

func TestGetMovie(t *testing.T) {
	movie1 := createRandomMovie(t)
	movie2, err := testStore.Movie.GetMovie(context.Background(), movie1.ID.Hex())
	require.NoError(t, err)
	require.NotEmpty(t, movie2)

	require.Equal(t, movie1.ID, movie2.ID)
	require.Equal(t, movie1.Title, movie2.Title)
	require.Equal(t, movie1.Duration, movie2.Duration)
	require.Equal(t, movie1.Genre, movie2.Genre)
	require.Equal(t, movie1.Directors, movie2.Directors)
	require.Equal(t, movie1.Actors, movie2.Actors)
	require.WithinDuration(t, movie1.Screening, movie2.Screening, time.Second)
	require.Equal(t, movie1.Plot, movie2.Plot)
	require.Equal(t, movie1.Poster, movie2.Poster)
	require.WithinDuration(t, movie1.CreatedAt, movie2.CreatedAt, time.Second)

}

func TestUpdateMovie(t *testing.T) {
	movie1 := createRandomMovie(t)

	arg := models.Movie{
		Title:     util.RandomString(50),
		Duration:  movie1.Duration,
		Genre:     movie1.Genre,
		Directors: movie1.Directors,
		Actors:    movie1.Actors,
		Screening: movie1.Screening,
		Plot:      movie1.Plot,
		Poster:    movie1.Poster,
	}

	movie2, err := testStore.Movie.UpdateMovie(context.Background(), movie1.ID.Hex(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, movie2)

	require.Equal(t, movie1.ID, movie2.ID)
	require.Equal(t, arg.Title, movie2.Title)
	require.Equal(t, movie1.Duration, movie2.Duration)
	require.Equal(t, movie1.Genre, movie2.Genre)
	require.Equal(t, movie1.Directors, movie2.Directors)
	require.Equal(t, movie1.Actors, movie2.Actors)
	require.WithinDuration(t, movie1.Screening, movie2.Screening, time.Second)
	require.Equal(t, movie1.Plot, movie2.Plot)
	require.Equal(t, movie1.Poster, movie2.Poster)
	//require.WithinDuration(t, movie1.CreatedAt, movie2.CreatedAt, time.Second)
}

func TestDeleteMovie(t *testing.T) {
	movie1 := createRandomMovie(t)
	err := testStore.Movie.DeleteMovie(context.Background(), movie1.ID.Hex())
	require.NoError(t, err)

	movie2, err := testStore.Movie.GetMovie(context.Background(), movie1.ID.Hex())
	fmt.Println(movie2, err)
	require.Error(t, err)
	require.EqualError(t, err, ErrMovieNotFound.Error())
	require.Empty(t, movie2)
}
