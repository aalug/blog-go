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
		Image:       utils.RandomString(3) + ".png",
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

// TestQueries_ListPostsByCategory tests the list posts by category function
func TestQueries_ListPostsByCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	category2 := createRandomCategory(t)
	user := createRandomUser(t)

	for i := 0; i < 10; i++ {
		var categoryID int64
		if i%2 == 1 {
			categoryID = category1.ID
		} else {
			categoryID = category2.ID
		}
		params := CreatePostParams{
			Title:       utils.RandomString(7),
			Description: utils.RandomString(8),
			Content:     utils.RandomString(9),
			AuthorID:    int32(user.ID),
			CategoryID:  int32(categoryID),
			Image:       "test.jpg",
		}
		_, err := testQueries.CreatePost(context.Background(), params)
		require.NoError(t, err)
	}

	params := ListPostsByCategoryParams{
		ID:     category1.ID,
		Limit:  10,
		Offset: 0,
	}
	var posts []ListPostsByCategoryRow
	var err error

	posts, err = testQueries.ListPostsByCategory(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, posts)
	// because of limit and offset, it may return up to 10 rows, but only 5 are in the Category1
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.NotZero(t, post.CreatedAt)
		require.Equal(t, post.CategoryName, category1.Name)
	}
}

// TestQueries_ListPostsByAuthor tests the list posts by category function
func TestQueries_ListPostsByAuthor(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	for i := 0; i < 10; i++ {
		var authorID int32
		if i%2 == 1 {
			authorID = int32(user1.ID)
		} else {
			authorID = int32(user2.ID)
		}
		categoryID := int32(createRandomCategory(t).ID)
		params := CreatePostParams{
			Title:       utils.RandomString(7),
			Description: utils.RandomString(8),
			Content:     utils.RandomString(9),
			AuthorID:    authorID,
			CategoryID:  categoryID,
			Image:       "test.jpg",
		}
		_, err := testQueries.CreatePost(context.Background(), params)
		require.NoError(t, err)
	}

	params := ListPostsByAuthorParams{
		AuthorID: int32(user1.ID),
		Limit:    10,
		Offset:   0,
	}
	var posts []ListPostsByAuthorRow
	var err error

	posts, err = testQueries.ListPostsByAuthor(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, posts)
	// because of limit and offset, it may return up to 10 rows, but only 5 are written by user1
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.NotZero(t, post.CreatedAt)
		require.Equal(t, post.AuthorUsername, user1.Username)
	}
}

// TestQueries_ListPostsByTags tests the list posts by category function
func TestQueries_ListPostsByTags(t *testing.T) {
	tag1 := createRandomTag(t)
	tag2 := createRandomTag(t)
	for i := 0; i < 10; i++ {
		post := createRandomPost(t)
		if i%2 == 1 {
			params := AddTagToPostParams{
				PostID: post.ID,
				TagID:  tag1.ID,
			}
			err := testQueries.AddTagToPost(context.Background(), params)
			require.NoError(t, err)
		}
		if i == 5 {
			params2 := AddTagToPostParams{
				PostID: post.ID,
				TagID:  tag2.ID,
			}
			err := testQueries.AddTagToPost(context.Background(), params2)
			require.NoError(t, err)
		}
	}

	params := ListPostsByTagsParams{
		TagIds: []int32{tag1.ID},
		Limit:  10,
		Offset: 0,
	}
	var posts []ListPostsByTagsRow
	var err error

	posts, err = testQueries.ListPostsByTags(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, posts)
	// because of limit and offset, it may return up to 10 rows, but only 5 have tag1
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.NotZero(t, post.CreatedAt)
	}

	params2 := ListPostsByTagsParams{
		TagIds: []int32{tag2.ID},
		Limit:  10,
		Offset: 0,
	}
	var posts2 []ListPostsByTagsRow
	posts2, err = testQueries.ListPostsByTags(context.Background(), params2)

	require.NoError(t, err)
	require.NotEmpty(t, posts2)
	// because of limit and offset, it may return up to 10 rows, but only 1 have tag2
	require.Len(t, posts2, 1)

	for _, post := range posts2 {
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
