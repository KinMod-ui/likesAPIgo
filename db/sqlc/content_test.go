package db

import (
	"context"
	"likesApi/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomContent(t *testing.T, user User) Content {
	title := util.RandomData()
	if (user == User{}) {
		user = createRandomUser(t)
	}

	arg := CreateContentParams{
		Title:  title,
		UserID: user.ID,
	}

	content, err := testQueries.CreateContent(
		context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, content)

	require.Equal(t, content.Title, title)
	require.Equal(t, content.UserID, user.ID)
	require.NotZero(t, content.CreatedAt)
	require.NotZero(t, content.ID)
	return content
}

func TestCreateContent(t *testing.T) {
	createRandomContent(t, User{})
}

func TestGetContent(t *testing.T) {
	c1 := createRandomContent(t, User{})
	c2, err := testQueries.GetContent(
		context.Background(), c1.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, c2)
	require.Equal(t, c1, c2)
}

func TestListContent(t *testing.T) {

	user := createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomContent(t, user)
	}

	arg := ListContentOfUserParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 5,
	}

	contents, err := testQueries.ListContentOfUser(
		context.Background(), arg,
	)
	require.NoError(t, err)
	require.Len(t, contents, 5)

	for _, cts := range contents {
		require.NotEmpty(t, cts)
	}
}
