package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tijanadmi/movieginmongoapi/models"
)

func CreateRandomAddReservation(t *testing.T) *models.Reservation {
	user := createRandomUser(t)
	movie := createRandomMovie(t)
	repertoire := createRandomRepertoireForMovie(t, movie.ID)
	arg := AddReservationParams{
		Username:    user.Username,
		MovieID:     movie.ID.Hex(),
		Date:        repertoire.Date,
		Time:        repertoire.Time,
		Hall:        repertoire.Hall,
		ReservSeats: []string{"A1", "A2"},
	}

	reservation, err := testStore.AddReservation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, reservation)

	require.Equal(t, arg.Username, reservation.Username)
	require.Equal(t, arg.MovieID, reservation.MovieID.Hex())
	require.Equal(t, arg.Date, reservation.Date)
	require.Equal(t, arg.Time, reservation.Time)
	require.Equal(t, arg.Hall, reservation.Hall)
	require.Equal(t, arg.ReservSeats, reservation.ReservSeats)

	require.NotZero(t, reservation.ID)
	require.NotZero(t, reservation.CreationDate)

	return reservation
}

func TestAddReservation(t *testing.T) {
	CreateRandomAddReservation(t)
}

func TestCancelReservation(t *testing.T) {
	reservation := CreateRandomAddReservation(t)

	err := testStore.DeleteReservation(context.Background(), reservation.ID.Hex())
	require.NoError(t, err)

	reservation1, err := testStore.GetReservationById(context.Background(), reservation.ID.Hex())
	//fmt.Println(reservation1, err)
	require.Error(t, err)
	require.EqualError(t, err, ErrReservationNotFound.Error())
	require.Empty(t, reservation1)
}
