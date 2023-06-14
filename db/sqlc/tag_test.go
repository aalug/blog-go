package db

import (
	"context"
	"github.com/aalug/blog-go/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomTag creates and returns a random tag
func createRandomTag(t *testing.T) Tag {
	name := utils.RandomString(4)
	tag, err := testQueries.CreateTag(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, tag)
	require.Equal(t, tag.Name, name)
	require.NotZero(t, tag.ID)

	return tag
}

// TestQueries_CreateTag tests the create tag function
func TestQueries_CreateTag(t *testing.T) {
	createRandomUser(t)
}

// TestQueries_ListTags tests the list tags function
func TestQueries_ListTags(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTag(t)
	}

	params := ListTagsParams{
		Limit:  5,
		Offset: 5,
	}

	tags, err := testQueries.ListTags(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.Len(t, tags, 5)

	for _, tag := range tags {
		require.NotEmpty(t, tag)
		require.NotZero(t, tag.ID)
	}
}

// TestQueries_DeleteTag tests the delete tag function
func TestQueries_DeleteTag(t *testing.T) {
	tag := createRandomTag(t)

	err := testQueries.DeleteTag(context.Background(), tag.Name)
	require.NoError(t, err)
}

// TestQueries_UpdateTag tests the update tag function
func TestQueries_UpdateTag(t *testing.T) {
	tag := createRandomTag(t)

	params := UpdateTagParams{
		Name:   tag.Name,
		Name_2: utils.RandomString(5),
	}

	updatedTag, err := testQueries.UpdateTag(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTag)
	require.Equal(t, updatedTag.Name, params.Name_2)
	require.NotZero(t, updatedTag.ID)
}

// TestQueries_GetOrCreateTags tests the get or create tags function
func TestQueries_GetOrCreateTags(t *testing.T) {
	tagNames := []string{"tag1", "tag2", "tag3"}
	ids, err := testQueries.GetOrCreateTags(context.Background(), tagNames)
	require.NoError(t, err)

	require.Len(t, ids, 3)
}
