package grpchandlerv1

import (
	"context"

	"github.com/google/uuid"
	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/auth/internal/usecase/authuser"
)

type LocalUserHandler struct {
	authv1.UnimplementedLocalUserServiceServer
	userUsecase authuser.LocalUserUsecase
	logger      *zap.Logger
}

// UpdateEmailVerification implements authv1.LocalUserServiceServer.
func (l *LocalUserHandler) UpdateEmailVerification(ctx context.Context, req *authv1.UpdateEmailVerificationRequest) (*authv1.UpdateEmailVerificationResponse, error) {
	if err := req.Validate(); err != nil {
		l.logger.Error("UpdateEmailVerification request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		l.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = l.userUsecase.UpdateLocalEmailVerificationStatus(ctx, userID, req.Verified)
	if err != nil {
		l.logger.Error("Failed to update email verification status", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to update email verification status: %v", err)
	}

	return &authv1.UpdateEmailVerificationResponse{
		UserId:    req.UserId,
		Verified:  req.Verified,
		UpdatedAt: timestamppb.Now(),
	}, nil
}

// CreateLocalUser implements authv1.LocalUserServiceServer.
func (l *LocalUserHandler) CreateLocalUser(ctx context.Context, req *authv1.CreateLocalUserRequest) (*authv1.CreateLocalUserResponse, error) {
	if err := req.Validate(); err != nil {
		l.logger.Error("CreateLocalUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		l.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	createdUser, err := l.userUsecase.CreateLocalAuthUser(ctx, userID, req.Email, req.Password)
	if err != nil {
		l.logger.Error("Failed to create local user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to create local user: %v", err)
	}

	return &authv1.CreateLocalUserResponse{
		UserId:    createdUser.UserID.String(),
		CreatedAt: timestamppb.Now(),
	}, nil
}

// DeleteLocalUser implements authv1.LocalUserServiceServer.
func (l *LocalUserHandler) DeleteLocalUser(ctx context.Context, req *authv1.DeleteLocalUserRequest) (*authv1.DeleteLocalUserResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		l.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = l.userUsecase.DeleteAuthUser(ctx, userID)
	if err != nil {
		l.logger.Error("Failed to delete local user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete local user: %v", err)
	}

	return &authv1.DeleteLocalUserResponse{
		UserId:    req.UserId,
		DeletedAt: timestamppb.Now(),
	}, nil
}

// UpdateLocalUserEmail implements authv1.LocalUserServiceServer.
func (l *LocalUserHandler) UpdateLocalUserEmail(ctx context.Context, req *authv1.UpdateLocalUserEmailRequest) (*authv1.UpdateLocalUserEmailResponse, error) {
	if err := req.Validate(); err != nil {
		l.logger.Error("UpdateLocalUserEmail request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		l.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	updatedUser, err := l.userUsecase.UpdateAuthUserEmail(ctx, userID, req.NewEmail)
	if err != nil {
		l.logger.Error("Failed to update local user email", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to update local user email: %v", err)
	}

	return &authv1.UpdateLocalUserEmailResponse{
		UserId:       updatedUser.UserID.String(),
		UpdatedEmail: updatedUser.Email,
		UpdatedAt:    timestamppb.Now(),
	}, nil
}

// NewUserHandler creates a new UserHandler with the provided use case and logger.
func NewLocalUserHandler(userUsecase authuser.LocalUserUsecase, logger *zap.Logger) authv1.LocalUserServiceServer {
	return &LocalUserHandler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}
