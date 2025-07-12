package grpchandlerv1

import (
	"context"
	"strings"

	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/go/profile/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/profile/internal/usecase/dto"
	"mandacode.com/accounts/profile/internal/usecase/system"
)

type ProfileHandler struct {
	profilev1.UnimplementedProfileServiceServer
	profile *system.ProfileUsecase
	logger  *zap.Logger
}

// UpdateEmail implements profilev1.ProfileServiceServer.
func (u *ProfileHandler) UpdateEmail(ctx context.Context, req *profilev1.UpdateEmailRequest) (*profilev1.UpdateEmailResponse, error) {
	if err := req.Validate(); err != nil {
		u.logger.Error("UpdateEmail request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		u.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}
	profile, err := u.profile.UpdateProfile(ctx, &dto.UpdateProfileData{
		UserID: userID,
		Email:  &req.NewEmail,
	})
	if err != nil {
		u.logger.Error("Failed to update email", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to update email: %v", err)
	}

	return &profilev1.UpdateEmailResponse{
		UserId:       profile.UserID.String(),
		UpdatedEmail: profile.Email,
		UpdatedAt:    timestamppb.Now(),
	}, nil
}

// DeleteUser implements userv1.UserServiceServer.
func (u *ProfileHandler) DeleteUser(ctx context.Context, req *profilev1.DeleteUserRequest) (*profilev1.DeleteUserResponse, error) {
	if err := req.Validate(); err != nil {
		u.logger.Error("DeleteUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		u.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = u.profile.DeleteProfile(ctx, userID)
	if err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err), zap.String("user_id", req.UserId))
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &profilev1.DeleteUserResponse{
		UserId:    req.UserId,
		DeletedAt: timestamppb.Now(),
	}, nil
}

// InitUser implements userv1.UserServiceServer.
func (u *ProfileHandler) InitUser(ctx context.Context, req *profilev1.InitUserRequest) (*profilev1.InitUserResponse, error) {
	if err := req.Validate(); err != nil {
		u.logger.Error("InitUser request validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		u.logger.Error("Invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	// parse initial nickname from email by splitting at '@' and taking the first part
	var nickname string
	if req.Email != "" {
		parts := strings.Split(req.Email, "@")
		if len(parts) > 0 {
			nickname = parts[0]
		}
	} else {
		nickname = "user_" + userID.String() // fallback nickname
	}

	user, err := u.profile.CreateProfile(ctx, &dto.CreateProfileData{
		UserID:   userID,
		Email:    req.Email,
		Nickname: nickname,
	})
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, status.Errorf(errcode.MapCodeToGRPC(appErr.Code()), appErr.Public())
		}
		return nil, status.Errorf(codes.Internal, "failed to initialize user: %v", err)
	}

	return &profilev1.InitUserResponse{
		UserId:        user.UserID.String(),
		InitializedAt: timestamppb.Now(),
	}, nil

}

// NewProfileHandler creates a new UserSystemHandler with the provided use case.
func NewProfileHandler(profile *system.ProfileUsecase, logger *zap.Logger) profilev1.ProfileServiceServer {
	return &ProfileHandler{
		profile: profile,
		logger:  logger,
	}
}
