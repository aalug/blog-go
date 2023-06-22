package api

import (
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
