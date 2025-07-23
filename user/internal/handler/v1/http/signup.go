package httphandlerv1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mandacode.com/accounts/user/internal/models/provider"
	"mandacode.com/accounts/user/internal/usecase/signup"
	signupdto "mandacode.com/accounts/user/internal/usecase/signup/dto"
)

type SignupHandler struct {
	signup      *signup.SingupUsecase
	verifyEmail *signup.VerifyEmailUsecase
	validator   *validator.Validate
	logger      *zap.Logger
}

// NewSignupHandler creates a new SignupHandler with the provided use case and logger.
func NewSignupHandler(
	signup *signup.SingupUsecase,
	verifyEmail *signup.VerifyEmailUsecase,
	validator *validator.Validate,
	logger *zap.Logger,
) *SignupHandler {
	if signup == nil {
		logger.Error("signup use case cannot be nil")
		return nil
	}
	if verifyEmail == nil {
		logger.Error("verifyEmail use case cannot be nil")
		return nil
	}
	if validator == nil {
		logger.Error("validator cannot be nil")
		return nil
	}
	if logger == nil {
		logger.Error("logger cannot be nil")
		return nil
	}
	return &SignupHandler{
		signup:      signup,
		verifyEmail: verifyEmail,
		validator:   validator,
		logger:      logger,
	}
}

// RegisterRoutes registers the user routes with the provided router.
func (h *SignupHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/", h.LocalSignup)
	router.GET("/verify-email", h.VerifyEmail)
	router.POST("/verify-email/resend/:user_id", h.ResendVerificationEmail)
	router.GET("/o/:provider", h.OAuthSignup)
}

// LocalSignup handles the local signup process.
func (h *SignupHandler) LocalSignup(ctx *gin.Context) {
	var req signupdto.LocalSignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON for local signup", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation failed for local signup", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	signupRes, err := h.signup.LocalSignup(ctx, &req)
	if err != nil {
		h.logger.Error("Local signup failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Local signup failed"})
		return
	}
	if err := h.verifyEmail.SendVerificationEmail(ctx, signupRes.UserID, signupRes.Email); err != nil {
		h.logger.Error("Failed to send verification email", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":    signupRes.UserID,
		"email":      signupRes.Email,
		"created_at": signupRes.CreatedAt,
	})
}

// VerifyEmail handles the email verification process.
func (h *SignupHandler) VerifyEmail(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		h.logger.Error("Email verification token is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email verification token is required"})
		return
	}
	userID, err := h.verifyEmail.VerifyEmailToken(ctx, token)
	if err != nil {
		h.logger.Error("Email verification failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Email verification failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Email verified successfully",
	})
}

// OAuthSignup handles the OAuth signup process.
func (h *SignupHandler) OAuthSignup(ctx *gin.Context) {
	providerParam := ctx.Param("provider")
	if providerParam == "" {
		h.logger.Error("OAuth provider is required", zap.String("provider", providerParam))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "OAuth provider is required"})
		return
	}
	accessToken := ctx.GetHeader("Authorization")
	if accessToken == "" {
		h.logger.Error("Authorization header is required for OAuth signup", zap.String("provider", providerParam))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	if accessToken == "" {
		h.logger.Error("Access token is required for OAuth signup", zap.String("provider", providerParam))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Access token is required"})
		return
	}

	providerEnum, err := provider.FromString(providerParam)
	if err != nil {
		h.logger.Error("Invalid OAuth provider", zap.String("provider", providerParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth provider"})
		return
	}

	signupRes, err := h.signup.OAuthSignup(ctx, &signupdto.OAuthSignupRequest{
		Provider:    providerEnum,
		AccessToken: accessToken,
	})
	if err != nil {
		h.logger.Error("OAuth signup failed", zap.String("provider", providerParam), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth signup failed"})
		return
	}

	ctx.JSON(http.StatusOK, signupRes)
}

// ResendVerificationEmail handles the request to resend the email verification.
func (h *SignupHandler) ResendVerificationEmail(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		h.logger.Error("User ID is required for resending verification email")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Error("Invalid User ID format for resend verification email", zap.String("user_id", userID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	var req signupdto.SendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON for resend verification email", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation failed for resend verification email", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	if err := h.verifyEmail.SendVerificationEmail(ctx, userUID, req.Email); err != nil {
		h.logger.Error("Failed to resend verification email", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend verification email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Verification email resent successfully"})
}
