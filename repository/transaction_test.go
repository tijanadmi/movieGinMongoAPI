package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddReservation(t *testing.T) {
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
	require.Equal(t, arg.MovieID, reservation.MovieID)
	require.Equal(t, arg.Date, reservation.Date)
	require.Equal(t, arg.Time, reservation.Time)
	require.Equal(t, arg.Hall, reservation.Hall)
	require.Equal(t, arg.ReservSeats, reservation.ReservSeats)

	require.NotZero(t, reservation.ID)
	require.NotZero(t, reservation.CreationDate)
}

// func TestCancelReservation(t *testing.T) {
// 	user := createRandomUser(t)
// 	movie := createRandomMovie(t)
// 	repertoire := createRandomRepertoireForMovie(t, movie.ID)
// 	arg := AddReservationParams{
// 		Username:    user.Username,
// 		MovieID:     movie.ID.Hex(),
// 		Date:        repertoire.Date,
// 		Time:        repertoire.Time,
// 		Hall:        repertoire.Hall,
// 		ReservSeats: []string{"A1", "A2"},
// 	}

// 	reservation,err := testStore.AddReservation(context.Background(), arg)
// 	require.NoError(t, err)

// 	err = testStore.DeleteReservation(context.Background(), reservation1.ID.Hex())
// 	require.NoError(t, err)

// 	Reservation2, err := testStore.GetReservationById(context.Background(), reservation1.ID.Hex())
// 	fmt.Println(Reservation2, err)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrReservationNotFound.Error())
// 	require.Empty(t, Reservation2)
// }