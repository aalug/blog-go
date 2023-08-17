package gapi

import (
	"fmt"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/mail"
	"github.com/aalug/blog-go/pb"
	"github.com/aalug/blog-go/token"
	"github.com/aalug/blog-go/utils"
	"github.com/aalug/blog-go/worker"
)

// Server serves gRPC requests for the service
type Server struct {
	pb.UnimplementedBlogGoServer
	config          utils.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
	MailChan        chan mail.Data
}

// NewServer creates a new gRPC server
func NewServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:  %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
