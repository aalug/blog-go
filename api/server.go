package api

import (
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/utils"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP  requests for the service
type Server struct {
	config utils.Config
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setups routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()

	return server, nil
}

// setupRouter sets up the HTTP routing
func (server *Server) setupRouter() {
	router := gin.Default()

	// --- users ---
	router.POST("/users", server.createUser)

	server.router = router
}

// Start runs the HTTP server on a given address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
