package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func addLike(t *testing.T, user User, content Content) Like {
	if (user == User{}) {
		user = createRandomUser(t)
	}

	if (content == Content{}) {
		content = createRandomContent(t, User{})
	}
	arg := CreateLikeParams{
		UserID:    user.ID,
		ContentID: content.ID,
	}

	likeAdd, err := testQueries.CreateLike(
		context.Background(), arg,
	)
	require.NoError(t, err)
	require.NotEmpty(t, likeAdd)
	require.Equal(t, likeAdd.UserID, user.ID)
	require.Equal(t, likeAdd.ContentID, content.ID)
	require.NotEmpty(t, likeAdd.UpdateAt)
	require.NotEmpty(t, likeAdd.Liked)

	return likeAdd

}

func TestCreateLike(t *testing.T) {
	addLike(t, User{}, Content{})
}

func TestGetLike(t *testing.T) {
	likeAdd1 := addLike(t, User{}, Content{})

	args := GetLikeParams{
		UserID:    likeAdd1.UserID,
		ContentID: likeAdd1.ContentID,
	}

	likeAdd2, err := testQueries.GetLike(
		context.Background(), args,
	)
	require.NoError(t, err)
	require.NotEmpty(t, likeAdd2)

	require.Equal(t, likeAdd1.Liked, likeAdd2)

}

func TestTotalLikes(t *testing.T) {

	user := createRandomUser(t)
	content := createRandomContent(t, user)
	addLike(t, user, content)
	for i := 0; i < 15; i++ {
		newUser := createRandomUser(t)
		addLike(t, newUser, content)
	}

	likes, err := testQueries.TotalLikesForContent(
		context.Background(), content.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, likes)
	require.Equal(t, likes, int64(16))
}
