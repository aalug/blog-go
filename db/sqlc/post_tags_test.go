package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestQueries_AddTagToPost tests the add tag to post function
func TestQueries_AddTagToPost(t *testing.T) {
	post := createRandomPost(t)
	tag := createRandomTag(t)

	params := AddTagToPostParams{
		PostID: post.ID,
		TagID:  tag.ID,
	}
	err := testQueries.AddTagToPost(context.Background(), params)
	require.NoError(t, err)
}

func TestQueries_RemoveTagFromPost(t *testing.T) {
	post := createRandomPost(t)
	tag := createRandomTag(t)

	params := AddTagToPostParams{
		PostID: post.ID,
		TagID:  tag.ID,
	}
	err := testQueries.AddTagToPost(context.Background(), params)
	require.NoError(t, err)

	params2 := RemoveTagFromPostParams{
		PostID: post.ID,
		TagID:  tag.ID,
	}

	err = testQueries.RemoveTagFromPost(context.Background(), params2)
	require.NoError(t, err)
}
