package api

import (
	"database/sql"
	"errors"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/token"
	"github.com/gin-gonic/gin"
	"net/http"
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

	err = server.store.DeletePost(ctx, int64(request.ID))
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
