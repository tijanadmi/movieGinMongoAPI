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

func createRandomHall(t *testing.T) *models.Hall {

	arg := models.Hall{
		Name: util.RandomHall(),
		Rows: []string{
			"A",
			"B",
			"C",
			"D",
			"E"},
		Cols: []int{
			1,
			2,
			3,
			4,
			5},
	}

	hall, err := testStore.Hall.InsertHall(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, hall)

	require.Equal(t, arg.Name, hall.Name)
	require.Equal(t, arg.Rows, hall.Rows)
	require.Equal(t, arg.Cols, hall.Cols)

	require.NotZero(t, hall.ID)
	require.NotZero(t, hall.CreatedAt)

	return hall
}

func TestCreateHall(t *testing.T) {
	createRandomHall(t)
}

func TestGetHall(t *testing.T) {
	hall1 := createRandomHall(t)
	hall2, err := testStore.Hall.GetHall(context.Background(), hall1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, hall2)

	for _, hall := range hall2 {
		require.Equal(t, hall1.ID, hall.ID)
		require.Equal(t, hall1.Name, hall.Name)
		require.Equal(t, hall1.Rows, hall.Rows)
		require.Equal(t, hall1.Cols, hall.Cols)
		require.WithinDuration(t, hall1.CreatedAt, hall.CreatedAt, time.Second)
	}
}

func TestListHall(t *testing.T) {

	halls, err := testStore.Hall.ListHalls(context.Background())
	require.NoError(t, err)

	createRandomHall(t)

	halls1, err := testStore.Hall.ListHalls(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, halls1)

	require.Equal(t, len(halls)+1, len(halls1))
}

func TestUpdateHall(t *testing.T) {
	hall1 := createRandomHall(t)

	arg := models.Hall{
		Name: hall1.Name,
		Rows: []string{
			"A",
			"B",
			"C",
			"D",
			"E",
			"F"},
		Cols: hall1.Cols,
	}

	hall2, err := testStore.Hall.UpdateHall(context.Background(), hall1.ID.Hex(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, hall2)

	require.Equal(t, hall1.ID, hall2.ID)
	require.Equal(t, hall1.Name, hall2.Name)
	require.Equal(t, arg.Rows, hall2.Rows)
	require.Equal(t, hall1.Cols, hall2.Cols)
}

func TestDeleteHall(t *testing.T) {
	hall1 := createRandomHall(t)
	err := testStore.Hall.DeleteHall(context.Background(), hall1.ID.Hex())
	require.NoError(t, err)

	hall2, err := testStore.Hall.GetHallById(context.Background(), hall1.ID.Hex())
	fmt.Println(hall2, err)
	require.Error(t, err)
	require.EqualError(t, err, ErrHallNotFound.Error())
	require.Empty(t, hall2)
}
