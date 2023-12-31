package api

import (
	"database/sql"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

type createCategoryRequest struct {
	Name string `json:"name" binding:"required,alpha"`
}

type createCategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// createCategory handles creating a new category
func (server *Server) createCategory(ctx *gin.Context) {
	var request createCategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.CreateCategory(ctx, request.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, res)
}

type deleteCategoryRequest struct {
	Name string `uri:"name" binding:"required,alpha"`
}

// deleteCategory handles deleting a category
func (server *Server) deleteCategory(ctx *gin.Context) {
	var request deleteCategoryRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteCategory(ctx, request.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type listCategoriesRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5"`
}

// listCategories handles listing categories
func (server *Server) listCategories(ctx *gin.Context) {
	var request listCategoriesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListCategoriesParams{
		Limit:  request.PageSize,
		Offset: (request.Page - 1) * request.PageSize,
	}

	categories, err := server.store.ListCategories(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

type updateCategoryRequest struct {
	OldName string `json:"old_name" binding:"required,alpha"`
	NewName string `json:"new_name" binding:"required,alpha"`
}

// updateCategory handles updating a category
func (server *Server) updateCategory(ctx *gin.Context) {
	var request updateCategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.UpdateCategoryParams{
		Name:   request.OldName,
		Name_2: request.NewName,
	}

	category, err := server.store.UpdateCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}
