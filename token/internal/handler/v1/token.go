package handlerv1

import (
	"context"

	tokenv1 "github.com/mandacode-com/accounts-proto/token/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	"mandacode.com/accounts/token/internal/usecase/token"
	"mandacode.com/accounts/token/internal/util"
)

type TokenHandler struct {
	tokenv1.UnimplementedTokenServiceServer
	token  *token.TokenUsecase
	logger *zap.Logger
}

func NewTokenHandler(
	token *token.TokenUsecase,
	logger *zap.Logger,
) (tokenv1.TokenServiceServer, error) {
	if token == nil {
		return nil, errors.New("token usecase cannot be nil", "Token Handler Error", errcode.ErrDependencyFailure)
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil", "Token Handler Error", errcode.ErrDependencyFailure)
	}
	return &TokenHandler{
		token:  token,
		logger: logger,
	}, nil
}

func (h *TokenHandler) logError(err error) {
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("application error", zap.String("message", appErr.Error()), zap.String("code", appErr.Code()), zap.String("trace", errors.Trace(err)))
		} else {
			h.logger.Error("unexpected error", zap.Error(err))
		}
	}
}

func (h *TokenHandler) GenerateAccessToken(ctx context.Context, req *tokenv1.GenerateAccessTokenRequest) (*tokenv1.GenerateAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		err = errors.Upgrade(err, "Invalid Access Token Request", errcode.ErrInvalidInput)
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	token, expiresAt, err := h.token.GenerateAccessToken(req.UserId)
	if err != nil {
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	return &tokenv1.GenerateAccessTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *TokenHandler) VerifyAccessToken(ctx context.Context, req *tokenv1.VerifyAccessTokenRequest) (*tokenv1.VerifyAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		err = errors.Upgrade(err, "Invalid Input", errcode.ErrInvalidInput)
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	userId, err := h.token.VerifyAccessToken(req.Token)
	if err != nil {
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	return &tokenv1.VerifyAccessTokenResponse{
		Valid:  true,
		UserId: userId,
	}, nil
}

func (h *TokenHandler) GenerateRefreshToken(ctx context.Context, req *tokenv1.GenerateRefreshTokenRequest) (*tokenv1.GenerateRefreshTokenResponse, error) {
	if err := req.Validate(); err != nil {
		err = errors.Upgrade(err, "Invalid Refresh Token Request", errcode.ErrInvalidInput)
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	token, expiresAt, err := h.token.GenerateRefreshToken(req.UserId)
	if err != nil {
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	return &tokenv1.GenerateRefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *TokenHandler) VerifyRefreshToken(ctx context.Context, req *tokenv1.VerifyRefreshTokenRequest) (*tokenv1.VerifyRefreshTokenResponse, error) {
	if err := req.Validate(); err != nil {
		err = errors.Upgrade(err, "Invalid Refresh Token Request", errcode.ErrInvalidInput)
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	userId, err := h.token.VerifyRefreshToken(req.Token)
	if err != nil {
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	return &tokenv1.VerifyRefreshTokenResponse{
		Valid:  true,
		UserId: userId,
	}, nil
}

func (h *TokenHandler) GenerateEmailVerificationToken(ctx context.Context, req *tokenv1.GenerateEmailVerificationTokenRequest) (*tokenv1.GenerateEmailVerificationTokenResponse, error) {
	if err := req.Validate(); err != nil {
		err = errors.Upgrade(err, "Invalid Email Verification Token Request", errcode.ErrInvalidInput)
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	token, expiresAt, err := h.token.GenerateEmailVerificationToken(req.UserId, req.Email, req.Code)
	if err != nil {
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	return &tokenv1.GenerateEmailVerificationTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *TokenHandler) VerifyEmailVerificationToken(ctx context.Context, req *tokenv1.VerifyEmailVerificationTokenRequest) (*tokenv1.VerifyEmailVerificationTokenResponse, error) {
	if err := req.Validate(); err != nil {
		err = errors.Upgrade(err, "Invalid Email Verification Token Request", errcode.ErrInvalidInput)
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	userID, email, code, err := h.token.VerifyEmailVerificationToken(req.Token)
	if err != nil {
		h.logError(err)
		return nil, util.NewGRPCError(err)
	}

	return &tokenv1.VerifyEmailVerificationTokenResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
		Code:   code,
	}, nil
}
