package api

import (
	"database/sql"
	"errors"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/token"
	"github.com/aalug/blog-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type createPostRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	Category    string   `json:"category" binding:"required,alpha"`
	Image       string   `json:"image" binding:"required"`
}

type createPostResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Author      string   `json:"author"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Image       string   `json:"image"`
}

// createPost creates a new post
func (server *Server) createPost(ctx *gin.Context) {
	var request createPostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get or create category and get the id
	categoryID, err := server.store.GetOrCreateCategory(ctx, request.Category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	authUser, err := server.store.GetUser(ctx, authPayload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := db.CreatePostParams{
		Title:       request.Title,
		Description: request.Description,
		Content:     request.Content,
		AuthorID:    int32(authUser.ID),
		CategoryID:  int32(categoryID),
		Image:       request.Image,
	}

	post, err := server.store.CreatePost(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	postTagsParams := db.AddTagsToPostParams{
		PostID: post.ID,
		Tags:   request.Tags,
	}
	err = server.store.AddTagsToPost(ctx, postTagsParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createPostResponse{
		Title:       post.Title,
		Description: post.Description,
		Content:     post.Content,
		Author:      authUser.Username,
		Category:    request.Category,
		Tags:        request.Tags,
		Image:       post.Image,
	}

	ctx.JSON(http.StatusCreated, res)
}

type deletePostRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deletePost deletes a post. Checks if the authenticated user is
// the author of the post, and if so, deletes the post.
func (server *Server) deletePost(ctx *gin.Context) {
	var request deletePostRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	authUser, err := server.store.GetUser(ctx, authPayload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// special function to get the bare minimum post data to validate the user
	// (is the logged-in user an author of this post)
	post, err := server.store.GetMinimalPostData(ctx, int64(request.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if post.AuthorID != int32(authUser.ID) {
		err := errors.New("post does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeletePost(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type getPostByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getPostResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Author      string   `json:"author"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Image       string   `json:"image"`
}

// getPostByID gets post details by id
func (server *Server) getPostByID(ctx *gin.Context) {
	var request getPostByIDRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := server.store.GetPostByID(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get tags for this post
	tags, err := server.store.GetTagsOfPost(ctx, post.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	res := getPostResponse{
		Title:       post.Title,
		Description: post.Description,
		Content:     post.Content,
		Author:      post.AuthorUsername,
		Category:    post.CategoryName,
		Tags:        tagNames,
		Image:       post.Image,
	}

	ctx.JSON(http.StatusOK, res)
}

type getPostByTitleRequest struct {
	Slug string `uri:"slug" binding:"required,slug"`
}

// getPostByID gets post details by title
func (server *Server) getPostByTitle(ctx *gin.Context) {
	var request getPostByTitleRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	title := strings.ReplaceAll(request.Slug, "-", " ")

	post, err := server.store.GetPostByTitle(ctx, title)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get tags for this post
	tags, err := server.store.GetTagsOfPost(ctx, post.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	res := getPostResponse{
		Title:       post.Title,
		Description: post.Description,
		Content:     post.Content,
		Author:      post.AuthorUsername,
		Category:    post.CategoryName,
		Tags:        tagNames,
		Image:       post.Image,
	}

	ctx.JSON(http.StatusOK, res)
}

type listPostsRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=15"`
}

// listPosts lists all posts
func (server *Server) listPosts(ctx *gin.Context) {
	var request listPostsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListPostsParams{
		Limit:  request.PageSize,
		Offset: (request.Page - 1) * request.PageSize,
	}

	posts, err := server.store.ListPosts(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

type listPostsByAuthorRequest struct {
	Page     int32  `form:"page" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=15"`
	Author   string `form:"author" binding:"required"`
}

// listPostsByAuthor lists all posts created by user with username or email
// containing the string in the request
func (server *Server) listPostsByAuthor(ctx *gin.Context) {
	var request listPostsByAuthorRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get all authors containing the string in the request
	authors, err := server.store.ListUsersContainingString(ctx, request.Author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if any authors were found
	if len(authors) == 0 {
		err := errors.New("no users found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	authorIDs := make([]int32, len(authors))
	for i, author := range authors {
		authorIDs[i] = int32(author.ID)
	}

	var allPosts []db.ListPostsByAuthorRow

	// get posts by author
	for _, authorID := range authorIDs {
		params := db.ListPostsByAuthorParams{
			AuthorID: authorID,
			Limit:    request.PageSize,
			Offset:   (request.Page - 1) * request.PageSize,
		}

		posts, err := server.store.ListPostsByAuthor(ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		allPosts = append(allPosts, posts...)
	}

	ctx.JSON(http.StatusOK, allPosts)
}

type listPostsByCategoryRequest struct {
	Page       int32 `form:"page" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=15"`
	CategoryID int64 `form:"category_id" binding:"required,min=1"`
}

// listPostsByCategory  lists posts from the given category
func (server *Server) listPostsByCategory(ctx *gin.Context) {
	var request listPostsByCategoryRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListPostsByCategoryParams{
		ID:     request.CategoryID,
		Limit:  request.PageSize,
		Offset: (request.Page - 1) * request.PageSize,
	}

	posts, err := server.store.ListPostsByCategory(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(posts) == 0 {
		err := errors.New("no posts found in the given category")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

// listPostsByTagsRequest represents the request to list posts by tags
// where tag_ids is a comma separated list of tag ids
type listPostsByTagsRequest struct {
	Page     int32  `form:"page" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=15"`
	TagIDs   string `form:"tag_ids" binding:"required,tags"`
}

// listPostsByTags lists posts with the given tags
func (server *Server) listPostsByTags(ctx *gin.Context) {
	var request listPostsByTagsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get tag ids from the request and convert into a slice of int32
	tagIDs, err := utils.TagsToIntSlice(request.TagIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListPostsByTagsParams{
		Limit:  request.PageSize,
		Offset: (request.Page - 1) * request.PageSize,
		TagIds: tagIDs,
	}

	posts, err := server.store.ListPostsByTags(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(posts) == 0 {
		err := errors.New("no posts found with given tags")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

type updatePostUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updatePostRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	Category    string   `json:"category"`
	Image       string   `json:"image"`
}

type updatePostResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Image       string   `json:"image"`
}

// updatePost updates a post with provided details
func (server *Server) updatePost(ctx *gin.Context) {
	var uriRequest updatePostUriRequest
	if err := ctx.ShouldBindUri(&uriRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	authUser, err := server.store.GetUser(ctx, authPayload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if the user making the request is the author of the post
	p, err := server.store.GetMinimalPostData(ctx, uriRequest.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if p.AuthorID != int32(authUser.ID) {
		err := errors.New("post does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var request updatePostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if request.Title == "" &&
		request.Description == "" &&
		request.Content == "" &&
		len(request.Tags) == 0 &&
		request.Category == "" &&
		request.Image == "" {
		err := errors.New("no fields to update")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := server.store.GetPostByID(ctx, uriRequest.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var categoryName string
	if request.Category != "" {
		categoryName = request.Category
	} else {
		categoryName = post.CategoryName
	}
	categoryID, err := server.store.GetOrCreateCategory(ctx, categoryName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := db.UpdatePostParams{
		ID: uriRequest.ID,
		Title: func() string {
			if request.Title == "" {
				return post.Title
			}
			return request.Title
		}(),
		Description: func() string {
			if request.Description == "" {
				return post.Description
			}
			return request.Description
		}(),
		Content: func() string {
			if request.Content == "" {
				return post.Content
			}
			return request.Content
		}(),
		CategoryID: int32(categoryID),
		Image: func() string {
			if request.Image == "" {
				return post.Image
			}
			return request.Image
		}(),
		UpdatedAt: time.Now(),
	}

	updatedPost, err := server.store.UpdatePost(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(request.Tags) > 0 {
		postTags, err := server.store.GetTagsOfPost(ctx, uriRequest.ID)
		postTagNames := make([]string, len(postTags))
		for i, postTag := range postTags {
			postTagNames[i] = postTag.Name
		}

		// s1 - tags that are in post but not in request - should be removed
		// s2 - tags that are in request but not in post - should be added
		s1, s2 := utils.CompareTagLists(postTagNames, request.Tags)

		if len(s2) > 0 {
			paramsToAdd := db.AddTagsToPostParams{
				PostID: uriRequest.ID,
				Tags:   s2,
			}
			err = server.store.AddTagsToPost(ctx, paramsToAdd)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		if len(s1) > 0 {
			paramsToRemove := db.RemoveTagsFromPostParams{
				PostID: uriRequest.ID,
				Tags:   s1,
			}
			err = server.store.RemoveTagsFromPost(ctx, paramsToRemove)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	tags, err := server.store.GetTagsOfPost(ctx, uriRequest.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	tagsAfterUpdate := make([]string, len(tags))
	for i, tag := range tags {
		tagsAfterUpdate[i] = tag.Name
	}

	res := updatePostResponse{
		Title:       updatedPost.Title,
		Description: updatedPost.Description,
		Content:     updatedPost.Content,
		Category:    categoryName,
		Tags:        tagsAfterUpdate,
		Image:       updatedPost.Image,
	}

	ctx.JSON(http.StatusOK, res)
}
