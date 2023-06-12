package db

import (
	"context"
	"database/sql"
	"github.com/aalug/blog-go/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// createRandomPost creates and returns a random post
func createRandomPost(t *testing.T) Post {
	user := createRandomUser(t)
	category := createRandomCategory(t)

	params := CreatePostParams{
		Title:       utils.RandomString(6),
		Description: utils.RandomString(7),
		Content:     utils.RandomString(10),
		AuthorID:    int32(user.ID),
		CategoryID:  int32(category.ID),
		Image:       "image.jpg",
	}

	post, err := testQueries.CreatePost(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.NotZero(t, post.ID)
	require.Equal(t, post.Title, params.Title)
	require.Equal(t, post.Description, params.Description)
	require.Equal(t, post.Content, params.Content)
	require.Equal(t, post.AuthorID, params.AuthorID)
	require.Equal(t, post.CategoryID, params.CategoryID)
	require.Equal(t, post.Image, params.Image)
	require.NotZero(t, post.CreatedAt)
	require.NotZero(t, post.UpdatedAt)

	return post
}

// TestQueries_CreatePost tests the create post function
func TestQueries_CreatePost(t *testing.T) {
	createRandomPost(t)
}

// TestQueries_GetPostByID tests the get post by id function
func TestQueries_GetPostByID(t *testing.T) {
	post := createRandomPost(t)
	post2, err := testQueries.GetPostByID(context.Background(), post.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post2)
	require.Equal(t, post.ID, post2.ID)
	require.Equal(t, post.Title, post2.Title)
	require.Equal(t, post.Description, post2.Description)
	require.Equal(t, post.Content, post2.Content)
	require.Equal(t, post.AuthorID, post2.AuthorID)
	require.Equal(t, post.CategoryID, post2.CategoryID)
	require.Equal(t, post.Image, post2.Image)
	require.WithinDuration(t, post.CreatedAt, post2.CreatedAt, time.Second)
	require.WithinDuration(t, post.UpdatedAt, post2.UpdatedAt, time.Second)
}

// TestQueries_GetPostByTitle tests the get post by title function
func TestQueries_GetPostByTitle(t *testing.T) {
	post := createRandomPost(t)
	post2, err := testQueries.GetPostByTitle(context.Background(), post.Title)
	require.NoError(t, err)
	require.NotEmpty(t, post2)
	require.Equal(t, post.ID, post2.ID)
	require.Equal(t, post.Title, post2.Title)
	require.Equal(t, post.Description, post2.Description)
	require.Equal(t, post.Content, post2.Content)
	require.Equal(t, post.AuthorID, post2.AuthorID)
	require.Equal(t, post.CategoryID, post2.CategoryID)
	require.Equal(t, post.Image, post2.Image)
	require.WithinDuration(t, post.CreatedAt, post2.CreatedAt, time.Second)
	require.WithinDuration(t, post.UpdatedAt, post2.UpdatedAt, time.Second)
}

// TestQueries_ListPosts tests the list posts function
func TestQueries_ListPosts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPost(t)
	}

	params := ListPostsParams{
		Limit:  5,
		Offset: 5,
	}
	var posts []ListPostsRow
	var err error

	posts, err = testQueries.ListPosts(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, posts)
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.NotZero(t, post.CreatedAt)
	}
}

// TestQueries_DeletePost tests the delete post function
func TestQueries_DeletePost(t *testing.T) {
	post := createRandomPost(t)

	err := testQueries.DeletePost(context.Background(), post.ID)
	require.NoError(t, err)

	post2, err := testQueries.GetPostByID(context.Background(), post.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, post2)
}

// TestQueries_UpdatePost tests the update post function
func TestQueries_UpdatePost(t *testing.T) {
	post := createRandomPost(t)

	params := UpdatePostParams{
		ID:          post.ID,
		Title:       "new title",
		Description: "new description",
		Content:     "nwe content",
		CategoryID:  post.CategoryID,
		Image:       post.Image,
		UpdatedAt:   time.Now(),
	}

	updatedPost, err := testQueries.UpdatePost(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedPost)
	require.Equal(t, updatedPost.ID, post.ID)
	require.Equal(t, updatedPost.Title, params.Title)
	require.Equal(t, updatedPost.Description, params.Description)
	require.WithinDuration(t, updatedPost.UpdatedAt, params.UpdatedAt, time.Second)
}
