package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tijanadmi/moveginmongo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createRandomReservation(t *testing.T) *models.Reservation {
	user := createRandomUser(t)
	movie := createRandomMovie(t)
	repertoire := createRandomRepertoireForMovie(t, movie.ID)
	arg := models.Reservation{
		Username:      user.Username,
		UserID:        user.ID,
		MovieID:       movie.ID,
		RepertoiresID: repertoire.ID,
		MovieTitle:    movie.Title,
		Date:          repertoire.Date,
		Time:          repertoire.Time,
		Hall:          repertoire.Hall,
		ReservSeats:   []string{"A1", "A2"},
	}

	reservation, err := testStore.Reservation.InsertReservation(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, reservation)

	require.Equal(t, arg.Username, reservation.Username)
	require.Equal(t, arg.UserID, reservation.UserID)
	require.Equal(t, arg.MovieID, reservation.MovieID)
	require.Equal(t, arg.RepertoiresID, reservation.RepertoiresID)
	require.Equal(t, arg.MovieTitle, reservation.MovieTitle)
	require.Equal(t, arg.Date, reservation.Date)
	require.Equal(t, arg.Time, reservation.Time)
	require.Equal(t, arg.Hall, reservation.Hall)
	require.Equal(t, arg.ReservSeats, reservation.ReservSeats)

	require.NotZero(t, reservation.ID)
	require.NotZero(t, reservation.CreationDate)

	return reservation
}

func TestCreateReservation(t *testing.T) {
	createRandomReservation(t)
}

func createRandomReservationForUser(t *testing.T, userID primitive.ObjectID, username string) *models.Reservation {

	movie := createRandomMovie(t)
	repertoire := createRandomRepertoireForMovie(t, movie.ID)
	arg := models.Reservation{
		Username:      username,
		UserID:        userID,
		MovieID:       movie.ID,
		RepertoiresID: repertoire.ID,
		MovieTitle:    movie.Title,
		Date:          repertoire.Date,
		Time:          repertoire.Time,
		Hall:          repertoire.Hall,
		ReservSeats:   []string{"A1", "A2"},
	}

	reservation, err := testStore.Reservation.InsertReservation(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, reservation)

	require.Equal(t, arg.Username, reservation.Username)
	require.Equal(t, arg.UserID, reservation.UserID)
	require.Equal(t, arg.MovieID, reservation.MovieID)
	require.Equal(t, arg.RepertoiresID, reservation.RepertoiresID)
	require.Equal(t, arg.MovieTitle, reservation.MovieTitle)
	require.Equal(t, arg.Date, reservation.Date)
	require.Equal(t, arg.Time, reservation.Time)
	require.Equal(t, arg.Hall, reservation.Hall)
	require.Equal(t, arg.ReservSeats, reservation.ReservSeats)

	require.NotZero(t, reservation.ID)
	require.NotZero(t, reservation.CreationDate)

	return reservation
}

func TestCreateReservationForUser(t *testing.T) {
	user := createRandomUser(t)
	createRandomReservationForUser(t, user.ID, user.Username)
}

func TestGetAllReservationsForUser(t *testing.T) {
	user := createRandomUser(t)
	reservations, err := testStore.Reservation.GetAllReservationsForUser(context.Background(), user.Username)
	require.NoError(t, err)

	createRandomReservationForUser(t, user.ID, user.Username)

	reservations1, err := testStore.Reservation.GetAllReservationsForUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, reservations1)

	require.Equal(t, len(reservations)+1, len(reservations1))
}

func TestGetReservation(t *testing.T) {
	reservation1 := createRandomReservation(t)
	reservation2, err := testStore.Reservation.GetReservationById(context.Background(), reservation1.ID.Hex())
	require.NoError(t, err)
	require.NotEmpty(t, reservation2)

	require.Equal(t, reservation1.ID, reservation2.ID)
	require.Equal(t, reservation1.Username, reservation2.Username)
	require.Equal(t, reservation1.UserID, reservation2.UserID)
	require.Equal(t, reservation1.MovieID, reservation2.MovieID)
	require.Equal(t, reservation1.RepertoiresID, reservation2.RepertoiresID)
	require.Equal(t, reservation1.MovieTitle, reservation2.MovieTitle)
	require.Equal(t, reservation1.Date, reservation2.Date)
	require.Equal(t, reservation1.Time, reservation2.Time)
	require.Equal(t, reservation1.Hall, reservation2.Hall)
	require.Equal(t, reservation1.ReservSeats, reservation2.ReservSeats)
	require.WithinDuration(t, reservation1.CreationDate, reservation2.CreationDate, time.Second)

}

func TestDeleteReservation(t *testing.T) {
	reservation1 := createRandomReservation(t)
	err := testStore.Reservation.DeleteReservation(context.Background(), reservation1.ID.Hex())
	require.NoError(t, err)

	Reservation2, err := testStore.Reservation.GetReservationById(context.Background(), reservation1.ID.Hex())
	fmt.Println(Reservation2, err)
	require.Error(t, err)
	require.EqualError(t, err, ErrReservationNotFound.Error())
	require.Empty(t, Reservation2)
}
