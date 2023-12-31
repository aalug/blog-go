package api

import (
	"database/sql"
	"errors"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type createCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  int32  `json:"post_id" binding:"required,min=1"`
}

// createComment creates a comment for a post. As a comment author
// sets the authenticated user
func (server *Server) createComment(ctx *gin.Context) {
	var request createCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	authUser, err := server.store.GetUser(ctx, authPayload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := db.CreateCommentParams{
		Content: request.Content,
		UserID:  int32(authUser.ID),
		PostID:  request.PostID,
	}

	comment, err := server.store.CreateComment(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

type deleteCommentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteComment deletes a comment. Checks if the authenticated user
// is the author of the comment, and if so, deletes the comment.
func (server *Server) deleteComment(ctx *gin.Context) {
	var request deleteCommentRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the comment id to check if it exists and
	// to check if the authenticated user is the author
	comment, err := server.store.GetComment(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	authUser, err := server.store.GetUser(ctx, authPayload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if the authenticated user is the author
	if comment.UserID != int32(authUser.ID) {
		err := errors.New("comment does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteComment(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type updateCommentUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// updateComment updates a comment. Checks if the authenticated user
// is the author of the comment, and if so, updates the comment.
func (server *Server) updateComment(ctx *gin.Context) {
	var uriRequest updateCommentUriRequest
	if err := ctx.ShouldBindUri(&uriRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the comment id to check if it exists and
	// to check if the authenticated user is the author
	comment, err := server.store.GetComment(ctx, uriRequest.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	authUser, err := server.store.GetUser(ctx, authPayload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if the authenticated user is the author
	if comment.UserID != int32(authUser.ID) {
		err := errors.New("comment does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var request updateCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.UpdateCommentParams{
		ID:      uriRequest.ID,
		Content: request.Content,
	}

	updatedComment, err := server.store.UpdateComment(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedComment)
}

type listCommentsUriRequest struct {
	PostID int32 `uri:"post_id" binding:"required,min=1"`
}

type listCommentsRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=15"`
}

// listComments lists comments for a post.
// Get the post ID from the URI and the pagination parameters from the query parameters.
func (server *Server) listComments(ctx *gin.Context) {
	var uriRequest listCommentsUriRequest
	if err := ctx.ShouldBindUri(&uriRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var request listCommentsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListCommentsForPostParams{
		PostID: uriRequest.PostID,
		Limit:  request.PageSize,
		Offset: (request.Page - 1) * request.PageSize,
	}

	comments, err := server.store.ListCommentsForPost(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
