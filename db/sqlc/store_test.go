package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSQLStore_AddTagsToPost tests the AddTagsToPost method
func TestSQLStore_AddTagsToPost(t *testing.T) {
	post := createRandomPost(t)
	tags := []string{"tag1", "tag2"}

	params := AddTagsToPostParams{
		PostID: post.ID,
		Tags:   tags,
	}

	err := testStore.AddTagsToPost(context.Background(), params)
	require.NoError(t, err)

	postTags, err := testQueries.GetTagsOfPost(context.Background(), post.ID)
	require.NoError(t, err)
	require.Equal(t, len(tags), len(postTags))
	for _, tag := range postTags {
		require.Contains(t, tags, tag.Name)
	}
}
