package api

import (
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
