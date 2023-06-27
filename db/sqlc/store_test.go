package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSQLStore_AddAndRemoveTagsToPost tests the AddTagsToPost
// and the RemoveTagsFromPost methods
func TestSQLStore_AddAndRemoveTagsToPost(t *testing.T) {
	post := createRandomPost(t)
	tags := []string{"tag1", "tag2", "tag3"}

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

	// RemoveTagsFromPost
	removeParams := RemoveTagsFromPostParams{
		PostID: post.ID,
		Tags:   []string{postTags[0].Name},
	}

	err = testStore.RemoveTagsFromPost(context.Background(), removeParams)
	require.NoError(t, err)
	postTags, err = testQueries.GetTagsOfPost(context.Background(), post.ID)
	require.NoError(t, err)
	require.Equal(t, len(tags)-1, len(postTags))

	err = testStore.RemoveAllTagsFromPost(context.Background(), post.ID)
	require.NoError(t, err)
	postTags, err = testQueries.GetTagsOfPost(context.Background(), post.ID)
	require.NoError(t, err)
	require.Equal(t, 0, len(postTags))
}
