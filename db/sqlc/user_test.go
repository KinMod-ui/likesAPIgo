package db

import (
	"context"
	"likesApi/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	username := util.RandomUsername()

	user, err := testQueries.CreateUser(context.Background(), username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, username)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {

	createRandomUser(t)
}

func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, user := range users {
		require.NotEmpty(t, user)
	}

}
