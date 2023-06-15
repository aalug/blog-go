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

func TestQueries_RemoveTagsFromPost(t *testing.T) {
	post := createRandomPost(t)
	tag := createRandomTag(t)

	params := AddTagToPostParams{
		PostID: post.ID,
		TagID:  tag.ID,
	}
	err := testQueries.AddTagToPost(context.Background(), params)
	require.NoError(t, err)

	deleteParams := DeleteTagsFromPostParams{
		PostID: int32(post.ID),
		TagIds: []int32{tag.ID},
	}

	err = testQueries.DeleteTagsFromPost(context.Background(), deleteParams)
	require.NoError(t, err)
}
