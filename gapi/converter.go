package gapi

import (
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertUser converts a db.User object to a UserResponse object
func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
