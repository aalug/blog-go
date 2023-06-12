package db

import (
	"context"
	"github.com/aalug/blog-go/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomComment creates and returns a random comment
func createRandomComment(t *testing.T) Comment {
	user := createRandomUser(t)
	post := createRandomPost(t)

	params := CreateCommentParams{
		Content: utils.RandomString(10),
		UserID:  int32(user.ID),
		PostID:  int32(post.ID),
	}

	comment, err := testQueries.CreateComment(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	require.NotZero(t, comment.ID)
	require.Equal(t, comment.Content, params.Content)
	require.Equal(t, comment.UserID, params.UserID)
	require.Equal(t, comment.PostID, params.PostID)

	return comment
}

// TestQueries_CreateComment tests the create comment function
func TestQueries_CreateComment(t *testing.T) {
	createRandomComment(t)
}

// TestQueries_ListCommentsForPost tests the list comments for post function
func TestQueries_ListCommentsForPost(t *testing.T) {
	user := createRandomUser(t)
	post := createRandomPost(t)

	for i := 0; i < 10; i++ {
		params := CreateCommentParams{
			Content: utils.RandomString(10),
			UserID:  int32(user.ID),
			PostID:  int32(post.ID),
		}
		_, err := testQueries.CreateComment(context.Background(), params)
		require.NoError(t, err)
	}

	params := ListCommentsForPostParams{
		Limit:  5,
		Offset: 5,
		PostID: int32(post.ID),
	}

	comments, err := testQueries.ListCommentsForPost(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Len(t, comments, 5)

	for _, comment := range comments {
		require.NotEmpty(t, comment)
		require.NotZero(t, comment.ID)
		require.NotZero(t, comment.CreatedAt)
		require.Equal(t, int64(comment.PostID), post.ID)
		require.Equal(t, int64(comment.UserID), user.ID)
		require.NotEmpty(t, comment.Content)
	}
}

// TestQueries_DeleteComment tests the delete comment function
func TestQueries_DeleteComment(t *testing.T) {
	comment := createRandomComment(t)

	err := testQueries.DeleteComment(context.Background(), comment.ID)
	require.NoError(t, err)
}

// TestQueries_UpdateComment tests the update comment function
func TestQueries_UpdateComment(t *testing.T) {
	comment := createRandomComment(t)

	params := UpdateCommentParams{
		ID:      comment.ID,
		Content: "new content",
	}

	updatedComment, err := testQueries.UpdateComment(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment)
	require.Equal(t, updatedComment.Content, params.Content)
	require.Equal(t, updatedComment.ID, comment.ID)
	require.Equal(t, updatedComment.UserID, comment.UserID)
	require.Equal(t, updatedComment.PostID, comment.PostID)
	require.NotZero(t, updatedComment.CreatedAt)
}
