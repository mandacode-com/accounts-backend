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
)

type LocalAuthHandler struct {
	localLogin *login.LocalLoginUsecase
	logger     *zap.Logger
	validator  *validator.Validate
}

func NewLocalAuthHandler(
	localLogin *login.LocalLoginUsecase,
	logger *zap.Logger,
	validator *validator.Validate,
) (*LocalAuthHandler, error) {
	if localLogin == nil {
		return nil, stdErrors.New("localLogin cannot be nil")
	}
	if validator == nil {
		return nil, stdErrors.New("validator cannot be nil")
	}

	return &LocalAuthHandler{
		localLogin: localLogin,
		logger:     logger,
		validator:  validator,
	}, nil
}

func (h *LocalAuthHandler) ValidateRequest(req interface{}) error {
	if req == nil {
		return errors.New("request cannot be nil", "InvalidRequest", errcode.ErrInvalidInput)
	}
	if err := h.validator.Struct(req); err != nil {
		joinedErr := errors.Join(err, "validation failed")
		return errors.Upgrade(joinedErr, "InvalidRequest", errcode.ErrInvalidInput)
	}
	return nil
}

// RegisterRoutes registers the local authentication routes
func (h *LocalAuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.Login)
	rg.POST("/login/code", h.LoginCode)
	rg.GET("/verify/:userID", h.VerifyCode)
}

// Login handles local user login
func (h *LocalAuthHandler) Login(c *gin.Context) {
	responseType := c.Query("response_type")
	if responseType != "" && responseType != "direct" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid response type"})
		return
	}

	var req handlerv1dto.LocalLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.ValidateRequest(&req); err != nil {
		c.Error(err)
		return
	}

	input := logindto.LocalLoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	// If responseType is "direct", return access and refresh tokens directly
	if responseType == "direct" {
		accessToken, refreshToken, err := h.localLogin.Login(c.Request.Context(), input)
		if err != nil {
			c.Error(err)
			return
		}
		response := handlerv1dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// If responseType is not "direct", save the refresh token in the session
	accessToken, refreshToken, err := h.localLogin.Login(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	session := sessions.Default(c)
	session.Set("refresh_token", refreshToken)
	if err := session.Save(); err != nil {
		c.Error(err)
		return
	}
	response := handlerv1dto.AccessTokenResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, response)
}

// LoginCode handles issuing a login code for local user login
func (h *LocalAuthHandler) LoginCode(c *gin.Context) {
	var req handlerv1dto.LocalLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.ValidateRequest(&req); err != nil {
		c.Error(err)
		return
	}

	input := logindto.LocalLoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	code, userID, err := h.localLogin.IssueLoginCode(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	response := handlerv1dto.IssueCodeResponse{
		Code:   code,
		UserID: userID.String(),
	}
	c.JSON(http.StatusOK, response)
}

// VerifyCode handles verification of the login code
func (h *LocalAuthHandler) VerifyCode(c *gin.Context) {
	// Get userID and code from the request
	userID := c.Param("userID")
	if userID == "" {
		c.Error(errors.New("userID is required", "InvalidUserID", errcode.ErrInvalidInput))
		return
	}
	code := c.Query("code")
	if code == "" {
		c.Error(errors.New("code is required", "InvalidCode", errcode.ErrInvalidInput))
		return
	}
	responseType := c.Query("response_type")
	if responseType != "direct" && responseType != "" {
		c.Error(errors.New("invalid response type", "InvalidResponseType", errcode.ErrInvalidInput))
		return
	}

	// Validate userID format
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		c.Error(errors.New("invalid userID format", "InvalidUserIDFormat", errcode.ErrInvalidInput))
		return
	}
	// Verify the login code
	accessToken, refreshToken, err := h.localLogin.VerifyLoginCode(c.Request.Context(), userIDParsed, code)
	if err != nil {
		c.Error(err)
		return
	}

	if responseType == "direct" {
		response := handlerv1dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	session := sessions.Default(c)
	session.Set("refresh_token", refreshToken)
	if err := session.Save(); err != nil {
		c.Error(err)
		return
	}
	response := handlerv1dto.AccessTokenResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, response)
}
