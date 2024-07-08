package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createRandomRepertoire(t *testing.T) *models.Repertoire {
	movie := createRandomMovie(t)
	hall := createRandomHall(t)

	now := time.Now()
	formattedDate := now.Format("2006-01-02")
	dateValue, _ := util.ParseDate(formattedDate)

	arg := models.Repertoire{
		MovieID:         movie.ID,
		DateSt:          formattedDate,
		Date:            dateValue,
		Time:            "10:00",
		Hall:            hall.Name,
		NumOfTickets:    50,
		NumOfResTickets: 0,
	}

	repertoire, err := testStore.Repertoire.AddRepertoire(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, repertoire)

	require.Equal(t, arg.MovieID, repertoire.MovieID)
	require.Equal(t, arg.DateSt, repertoire.DateSt)
	require.Equal(t, arg.Date, repertoire.Date)
	require.Equal(t, arg.Time, repertoire.Time)
	require.Equal(t, arg.Hall, repertoire.Hall)
	require.Equal(t, arg.NumOfTickets, repertoire.NumOfTickets)
	require.Equal(t, arg.NumOfResTickets, repertoire.NumOfResTickets)

	require.NotZero(t, repertoire.ID)
	require.NotZero(t, repertoire.CreatedAt)

	return repertoire
}

func createRandomRepertoireForMovie(t *testing.T, movieId primitive.ObjectID) *models.Repertoire {

	hall := createRandomHall(t)

	now := time.Now()
	formattedDate := now.Format("2006-01-02")
	dateValue, _ := util.ParseDate(formattedDate)

	arg := models.Repertoire{
		MovieID:         movieId,
		DateSt:          formattedDate,
		Date:            dateValue,
		Time:            "10:00",
		Hall:            hall.Name,
		NumOfTickets:    50,
		NumOfResTickets: 0,
	}

	repertoire, err := testStore.Repertoire.AddRepertoire(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, repertoire)

	require.Equal(t, arg.MovieID, repertoire.MovieID)
	require.Equal(t, arg.DateSt, repertoire.DateSt)
	require.Equal(t, arg.Date, repertoire.Date)
	require.Equal(t, arg.Time, repertoire.Time)
	require.Equal(t, arg.Hall, repertoire.Hall)
	require.Equal(t, arg.NumOfTickets, repertoire.NumOfTickets)
	require.Equal(t, arg.NumOfResTickets, repertoire.NumOfResTickets)

	require.NotZero(t, repertoire.ID)
	require.NotZero(t, repertoire.CreatedAt)

	return repertoire
}

func TestCreateRepertoire(t *testing.T) {
	createRandomRepertoire(t)
}

func TestListRepertoire(t *testing.T) {

	repertoires, err := testStore.Repertoire.ListRepertoires(context.Background())
	require.NoError(t, err)

	createRandomRepertoire(t)

	repertoires1, err := testStore.Repertoire.ListRepertoires(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, repertoires1)

	require.Equal(t, len(repertoires)+1, len(repertoires1))
}

func TestGetRepertoire(t *testing.T) {
	repertoire1 := createRandomRepertoire(t)
	repertoire2, err := testStore.Repertoire.GetRepertoire(context.Background(), repertoire1.ID.Hex())
	require.NoError(t, err)
	require.NotEmpty(t, repertoire2)

	require.Equal(t, repertoire1.ID, repertoire2.ID)
	require.Equal(t, repertoire1.MovieID, repertoire2.MovieID)
	require.Equal(t, repertoire1.DateSt, repertoire2.DateSt)
	require.Equal(t, repertoire1.Date, repertoire2.Date)
	require.Equal(t, repertoire1.Time, repertoire2.Time)
	require.Equal(t, repertoire1.Hall, repertoire2.Hall)
	require.Equal(t, repertoire1.NumOfTickets, repertoire2.NumOfTickets)
	require.Equal(t, repertoire1.NumOfResTickets, repertoire2.NumOfResTickets)
	require.WithinDuration(t, repertoire1.CreatedAt, repertoire2.CreatedAt, time.Second)

}

func TestGetRepertoireByMovieDateTimeHall(t *testing.T) {
	repertoire1 := createRandomRepertoire(t)
	repertoire2, err := testStore.Repertoire.GetRepertoireByMovieDateTimeHall(context.Background(), repertoire1.MovieID.Hex(), repertoire1.Date, repertoire1.Time, repertoire1.Hall)
	require.NoError(t, err)
	require.NotEmpty(t, repertoire2)

	require.Equal(t, repertoire1.ID, repertoire2.ID)
	require.Equal(t, repertoire1.MovieID, repertoire2.MovieID)
	require.Equal(t, repertoire1.DateSt, repertoire2.DateSt)
	require.Equal(t, repertoire1.Date, repertoire2.Date)
	require.Equal(t, repertoire1.Time, repertoire2.Time)
	require.Equal(t, repertoire1.Hall, repertoire2.Hall)
	require.Equal(t, repertoire1.NumOfTickets, repertoire2.NumOfTickets)
	require.Equal(t, repertoire1.NumOfResTickets, repertoire2.NumOfResTickets)
	require.WithinDuration(t, repertoire1.CreatedAt, repertoire2.CreatedAt, time.Second)

}

func TestGetAllRepertoireForMovie(t *testing.T) {
	movie := createRandomMovie(t)
	repertoire1 := createRandomRepertoireForMovie(t, movie.ID)
	repertoires, err := testStore.Repertoire.GetAllRepertoireForMovie(context.Background(), movie.ID.Hex(), repertoire1.Date, repertoire1.Date)
	require.NoError(t, err)

	repertoire11 := createRandomRepertoireForMovie(t, movie.ID)
	var repertoires2 []models.Repertoire
	if repertoire1.Date.Before(repertoire11.Date) {
		repertoires2, err = testStore.Repertoire.GetAllRepertoireForMovie(context.Background(), movie.ID.Hex(), repertoire1.Date, repertoire11.Date)
	} else {
		repertoires2, err = testStore.Repertoire.GetAllRepertoireForMovie(context.Background(), movie.ID.Hex(), repertoire11.Date, repertoire1.Date)
	}
	repertoires2, err = testStore.Repertoire.GetAllRepertoireForMovie(context.Background(), movie.ID.Hex(), repertoire1.Date, repertoire1.Date)
	require.NoError(t, err)
	require.NotEmpty(t, repertoires2)

	require.Equal(t, len(repertoires)+1, len(repertoires2))

}

func TestUpdateRepertoire(t *testing.T) {
	repertoire1 := createRandomRepertoire(t)

	arg := models.Repertoire{
		MovieID:         repertoire1.MovieID,
		DateSt:          repertoire1.DateSt,
		Date:            repertoire1.Date,
		Time:            "12:00",
		Hall:            repertoire1.Hall,
		NumOfTickets:    50,
		NumOfResTickets: 2,
		ReservSeats:     []string{"A1", "A2"},
	}

	repertoire2, err := testStore.Repertoire.UpdateRepertoire(context.Background(), repertoire1.ID.Hex(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, repertoire2)

	require.Equal(t, repertoire1.ID, repertoire2.ID)
	require.Equal(t, repertoire1.MovieID, repertoire2.MovieID)
	require.Equal(t, repertoire1.DateSt, repertoire2.DateSt)
	require.Equal(t, repertoire1.Date, repertoire2.Date)
	require.Equal(t, arg.Time, repertoire2.Time)
	require.Equal(t, repertoire1.Hall, repertoire2.Hall)
	require.Equal(t, arg.NumOfTickets, repertoire2.NumOfTickets)
	require.Equal(t, arg.NumOfResTickets, repertoire2.NumOfResTickets)
	require.Equal(t, arg.NumOfResTickets, repertoire2.NumOfResTickets)
	require.NotEmpty(t, repertoire2.ReservSeats)

}

func TestDeleteRepertoire(t *testing.T) {
	repertoire1 := createRandomRepertoire(t)
	err := testStore.Repertoire.DeleteRepertoire(context.Background(), repertoire1.ID.Hex())
	require.NoError(t, err)

	repertoire2, err := testStore.Repertoire.GetRepertoire(context.Background(), repertoire1.ID.Hex())
	fmt.Println(repertoire2, err)
	require.Error(t, err)
	require.EqualError(t, err, ErrRepertoireNotFound.Error())
	require.Empty(t, repertoire2)
}

func TestDeleteRepertoireForMovie(t *testing.T) {
	repertoire1 := createRandomRepertoire(t)
	err := testStore.Repertoire.DeleteRepertoireForMovie(context.Background(), repertoire1.MovieID.Hex())
	require.NoError(t, err)

	repertoire2, err := testStore.Repertoire.GetRepertoire(context.Background(), repertoire1.ID.Hex())
	fmt.Println(repertoire2, err)
	require.Error(t, err)
	require.EqualError(t, err, ErrRepertoireNotFound.Error())
	require.Empty(t, repertoire2)
}
