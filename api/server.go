package api

import (
	"fmt"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/token"
	"github.com/aalug/blog-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

// Server serves HTTP  requests for the service
type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setups routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("slug", isValidSlug)
		if err != nil {
			log.Fatal("failed to register validation")
		}
	}

	server.setupRouter()

	return server, nil
}

// setupRouter sets up the HTTP routing
func (server *Server) setupRouter() {
	router := gin.Default()

	// --- users ---
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// --- categories ---
	router.GET("/category", server.listCategories)

	// --- posts ---
	router.GET("/posts/id/:id", server.getPostByID)
	router.GET("/posts/title/:slug", server.getPostByTitle)
	router.GET("/posts/all", server.listPosts)

	// ===== routes that require authentication =====
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// --- categories ---
	authRoutes.POST("/category", server.createCategory)
	authRoutes.DELETE("/category/:name", server.deleteCategory)
	authRoutes.PATCH("/category", server.updateCategory)

	// --- posts ---
	authRoutes.POST("/posts", server.createPost)
	authRoutes.DELETE("/posts/:id", server.deletePost)

	server.router = router
}

// Start runs the HTTP server on a given address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
