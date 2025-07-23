package httphandlerv1

import (
	stdErrors "errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	handlerv1dto "mandacode.com/accounts/auth/internal/handler/v1/http/dto"
	"mandacode.com/accounts/auth/internal/usecase/login"
	logindto "mandacode.com/accounts/auth/internal/usecase/login/dto"

	"mandacode.com/accounts/auth/internal/util"
)

type OAuthHandler struct {
	oauthLogin *login.OAuthLoginUsecase
	logger     *zap.Logger
	validator  *validator.Validate
}

// NewOAuthHandler creates a new OAuthHandler instance
func NewOAuthHandler(
	oauthLogin *login.OAuthLoginUsecase,
	logger *zap.Logger,
	validator *validator.Validate,
) (*OAuthHandler, error) {
	if oauthLogin == nil {
		return nil, stdErrors.New("oauthLogin cannot be nil")
	}
	if logger == nil {
		return nil, stdErrors.New("logger cannot be nil")
	}
	if validator == nil {
		return nil, stdErrors.New("validator cannot be nil")
	}

	return &OAuthHandler{
		oauthLogin: oauthLogin,
		logger:     logger,
		validator:  validator,
	}, nil
}

func (h *OAuthHandler) ValidateRequest(req interface{}) error {
	if req == nil {
		return errors.New("request cannot be nil", "InvalidRequest", errcode.ErrInvalidInput)
	}
	if err := h.validator.Struct(req); err != nil {
		joinedErr := errors.Join(err, "validation failed")
		return errors.Upgrade(joinedErr, "InvalidRequest", errcode.ErrInvalidInput)
	}
	return nil
}

func (h *OAuthHandler) LogError(err error) {
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error(errors.Trace(appErr))
		} else {
			h.logger.Error("error occurred", zap.Error(err))
		}
	}
}

// RegisterRoutes registers the OAuth routes
func (h *OAuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/login/:provider", h.Login)
	rg.POST("/m/login/:provider", h.MobileLogin)
	rg.GET("/callback/:provider", h.Callback)
	rg.GET("/verify/:user_id", h.VerifyCode)
}

func (h *OAuthHandler) Login(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	ctx := c.Request.Context()

	// Get Login URL from the use case
	loginURL, err := h.oauthLogin.GetLoginURL(ctx, provider)
	if err != nil {
		h.LogError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get login URL"})
		return
	}

	c.Redirect(http.StatusFound, loginURL)
}

func (h *OAuthHandler) MobileLogin(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	ctx := c.Request.Context()

	var req handlerv1dto.MobileOAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.LogError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.ValidateRequest(&req); err != nil {
		h.LogError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Exchange code for access token and user info
	providerEnum, err := util.ConvertToEnt(provider)
	if err != nil {
		h.LogError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}
	input := logindto.OAuthLoginInput{
		Provider:    providerEnum,
		AccessToken: req.AccessToken,
		Code:        "",
	}
	accessToken, refreshToken, err := h.oauthLogin.Login(ctx, input)

	if err != nil {
		h.LogError(err)
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(errcode.MapCodeToHTTP(appErr.Code()), gin.H{"error": appErr.Public()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login with OAuth"})
		}
		return
	}

	c.JSON(http.StatusOK, handlerv1dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *OAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	ctx := c.Request.Context()

	// Extract code from query parameters
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	// Exchange code for access token and user info
	providerEnum, err := util.ConvertToEnt(provider)
	if err != nil {
		h.LogError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}
	input := logindto.OAuthLoginInput{
		Provider:    providerEnum,
		Code:        code,
		AccessToken: "",
	}
	code, userID, err := h.oauthLogin.IssueLoginCode(ctx, input)
	if err != nil {
		h.LogError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login with OAuth"})
		return
	}

	response := handlerv1dto.OAuthCallbackResponse{
		Code:   code,
		UserID: userID.String(),
	}
	c.JSON(http.StatusOK, response)
}

func (h *OAuthHandler) VerifyCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}
	userID := c.Param("user_id")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		h.LogError(errors.New("invalid user_id format", "InvalidUserIDFormat", errcode.ErrInvalidInput))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	ctx := c.Request.Context()

	accessToken, refreshToken, err := h.oauthLogin.VerifyLoginCode(ctx, userUID, code)
	if err != nil {
		h.LogError(err)
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(errcode.MapCodeToHTTP(appErr.Code()), gin.H{"error": appErr.Public()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify login code"})
		}
		return
	}

	session := sessions.Default(c)
	session.Set("refresh_token", refreshToken)

	c.JSON(http.StatusOK, handlerv1dto.AccessTokenResponse{
		AccessToken: accessToken,
	})
}
