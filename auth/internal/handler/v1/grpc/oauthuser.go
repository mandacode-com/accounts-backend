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
	"mandacode.com/accounts/auth/internal/util"
)

type OAuthUserHandler struct {
	authv1.UnimplementedOAuthUserServiceServer
	userUsecase authuser.OAuthUserUsecase
	logger      *zap.Logger
}

// CreateOAuthUser implements authv1.OAuthUserServiceServer.
func (o *OAuthUserHandler) CreateOAuthUser(ctx context.Context, req *authv1.CreateOAuthUserRequest) (*authv1.CreateOAuthUserResponse, error) {
	if err := req.Validate(); err != nil {
		o.logger.Error("CreateOAuthUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		o.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	entProvider, err := util.FromProtoToEnt(req.Provider)
	if err != nil {
		o.logger.Error("Failed to convert provider from proto to ent", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Errorf(codes.InvalidArgument, "invalid provider: %v", err)
	}

	createdUser, err := o.userUsecase.CreateOAuthUser(ctx, userID, entProvider, req.AccessToken, req.Code)
	if err != nil {
		o.logger.Error("Failed to create OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to create OAuth user: %v", err)
	}

	return &authv1.CreateOAuthUserResponse{
		UserId:    createdUser.UserID.String(),
		Provider:  req.Provider,
		Email:     createdUser.Email,
		CreatedAt: timestamppb.Now(),
	}, nil
}

// DeleteOAuthUser implements authv1.OAuthUserServiceServer.
func (o *OAuthUserHandler) DeleteOAuthUser(ctx context.Context, req *authv1.DeleteOAuthUserRequest) (*authv1.DeleteOAuthUserResponse, error) {
	if err := req.Validate(); err != nil {
		o.logger.Error("DeleteOAuthUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		o.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = o.userUsecase.DeleteOAuthUser(ctx, userID)
	if err != nil {
		o.logger.Error("Failed to delete OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete OAuth user: %v", err)
	}

	return &authv1.DeleteOAuthUserResponse{
		UserId:    req.UserId,
		DeletedAt: timestamppb.Now(),
	}, nil
}

// SyncOAuthUser implements authv1.OAuthUserServiceServer.
func (o *OAuthUserHandler) SyncOAuthUser(ctx context.Context, req *authv1.SyncOAuthUserRequest) (*authv1.SyncOAuthUserResponse, error) {
	if err := req.Validate(); err != nil {
		o.logger.Error("SyncOAuthUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		o.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	entProvider, err := util.FromProtoToEnt(req.Provider)
	if err != nil {
		o.logger.Error("Failed to convert provider from proto to ent", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Errorf(codes.InvalidArgument, "invalid provider: %v", err)
	}

	updatedUser, err := o.userUsecase.SyncOAuthUser(ctx, userID, entProvider, req.AccessToken, req.Code)
	if err != nil {
		o.logger.Error("Failed to sync OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to sync OAuth user: %v", err)
	}
	return &authv1.SyncOAuthUserResponse{
		UserId:   updatedUser.UserID.String(),
		Provider: req.Provider,
		SyncedAt: timestamppb.Now(),
	}, nil
}

// NewUserHandler creates a new UserHandler with the provided use case and logger.
func NewOAuthUserHandler(userUsecase authuser.OAuthUserUsecase, logger *zap.Logger) authv1.OAuthUserServiceServer {
	return &OAuthUserHandler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}
