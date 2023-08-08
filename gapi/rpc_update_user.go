package gapi

import (
	"context"
	"database/sql"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/pb"
	"github.com/aalug/blog-go/utils"
	"github.com/aalug/blog-go/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// UpdateUser creates a new user
func (server *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// check if the user is updating their own account
	if authPayload.Email != request.GetEmail() {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this user")
	}

	params := db.UpdateUserParams{
		Username: sql.NullString{
			String: request.GetUsername(),
			Valid:  request.Username != nil,
		},
		Email: request.GetEmail(),
	}

	if request.Password != nil {
		hashedPassword, err := utils.HashPassword(request.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		// set new password to update
		params.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		// update the password changed at value
		params.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	res := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

// validateUpdateUserRequest validates all the fields of the request.
func validateUpdateUserRequest(request *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateEmail(request.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if request.Username != nil {
		if err := validation.ValidateUsername(request.GetUsername()); err != nil {
			violations = append(violations, fieldViolation("username", err))
		}
	}

	if request.Password != nil {
		if err := validation.ValidatePassword(request.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}
