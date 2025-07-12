package grpchandlerv1

import (
	"context"

	"github.com/google/uuid"
	userv1 "github.com/mandacode-com/accounts-proto/go/user/user/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/user/internal/usecase/user"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	userUsecase *user.UserUsecase
	logger      *zap.Logger
}

// DeleteUser implements userv1.UserServiceServer.
func (u *UserHandler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if err := req.Validate(); err != nil {
		u.logger.Error("DeleteUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		u.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = u.userUsecase.DeleteUser(ctx, userID)
	if err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &userv1.DeleteUserResponse{
		UserId:    req.UserId,
		DeletedAt: timestamppb.Now(),
	}, nil
}

// InitUser implements userv1.UserServiceServer.
func (u *UserHandler) InitUser(ctx context.Context, req *userv1.InitUserRequest) (*userv1.InitUserResponse, error) {
	if err := req.Validate(); err != nil {
		u.logger.Error("InitUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		u.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user, err := u.userUsecase.CreateUser(ctx, userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to initialize user: %v", err)
	}

	return &userv1.InitUserResponse{
		UserId:        user.ID.String(),
		InitializedAt: timestamppb.Now(),
	}, nil

}

// IsActive implements userv1.UserServiceServer.
func (u *UserHandler) IsActive(ctx context.Context, req *userv1.IsActiveRequest) (*userv1.IsActiveResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user, err := u.userUsecase.GetUserByID(ctx, userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to check user active status: %v", err)
	}

	return &userv1.IsActiveResponse{
		UserId:   user.ID.String(),
		IsActive: user.IsActive,
	}, nil
}

// IsBlocked implements userv1.UserServiceServer.
func (u *UserHandler) IsBlocked(ctx context.Context, req *userv1.IsBlockedRequest) (*userv1.IsBlockedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user, err := u.userUsecase.GetUserByID(ctx, userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to check user blocked status: %v", err)
	}

	return &userv1.IsBlockedResponse{
		UserId:    user.ID.String(),
		IsBlocked: user.IsBlocked,
	}, nil
}

// NewUserHandler creates a new UserSystemHandler with the provided use case.
func NewUserHandler(userUsecase *user.UserUsecase) userv1.UserServiceServer {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}
