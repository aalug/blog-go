package gapi

import (
	"context"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/pb"
	"github.com/aalug/blog-go/utils"
	"github.com/aalug/blog-go/validation"
	"github.com/aalug/blog-go/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// CreateUser creates a new user
func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	params := db.CreateUserParams{
		Username:       request.GetUsername(),
		HashedPassword: hashedPassword,
		Email:          request.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "email already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// send confirmation email
	taskPayload := &worker.PayloadSendVerificationEmail{
		Email: user.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = server.taskDistributor.DistributeTaskSendVerificationEmail(ctx, taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send verification email: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

// validateCreateUserRequest validates all the fields of the create user request.
func validateCreateUserRequest(request *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidateEmail(request.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validation.ValidatePassword(request.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
