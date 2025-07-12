package httphandlerv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/profile/internal/usecase/dto"
	"mandacode.com/accounts/profile/internal/usecase/user"
)

type UserProfileHandler struct {
	profile   *user.ProfileUsecase
	uidHeader string
}

func NewUserProfileHandler(profile *user.ProfileUsecase, uidHeader string) (*UserProfileHandler, error) {
	return &UserProfileHandler{
		profile:   profile,
		uidHeader: uidHeader,
	}, nil
}

func (h *UserProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", h.GetProfile)
	router.PUT("/", h.UpdateProfile)
}

func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	if userID == "" {
		err := errors.New("user ID header is required", "InvalidRequest", errcode.ErrInvalidInput)
		c.Error(err)
		return
	}

	userUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.Upgrade(err, "Invalid user ID format", errcode.ErrInvalidInput)
		c.Error(err)
		return
	}

	profile, err := h.profile.GetProfile(c.Request.Context(), userUID)
	if err != nil {
		if errors.Is(err, errcode.ErrNotFound) {
			err = errors.Join(err, "Get profile handler failed")
			c.Error(err)
			return
		}
		err = errors.Upgrade(err, "Failed to get profile", errcode.ErrInternalFailure)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	if userID == "" {
		c.JSON(http.StatusBadRequest, errors.New("user ID header is required", "InvalidRequest", errcode.ErrInvalidInput))
		return
	}

	var updateData dto.UpdateProfileData
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, errors.Upgrade(err, "Invalid request data", errcode.ErrInvalidInput))
		return
	}

	profile, err := h.profile.UpdateProfile(c.Request.Context(), &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Upgrade(err, "Failed to update profile", errcode.ErrInternalFailure))
		return
	}

	c.JSON(http.StatusOK, profile)
}
