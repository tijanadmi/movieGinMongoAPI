package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
)

func createRandomUser(t *testing.T) *models.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := models.User{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		Roles:    []string{"user", "admin"},
	}

	user, err := testStore.Users.InsertUser(context.Background(), &arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.NotZero(t, user.DateOfCreation)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.Users.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Roles, user2.Roles)
	require.WithinDuration(t, user1.DateOfCreation, user2.DateOfCreation, time.Second)
}
