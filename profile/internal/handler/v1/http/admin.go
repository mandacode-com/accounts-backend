package httphandlerv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/profile/internal/usecase/admin"
	"mandacode.com/accounts/profile/internal/usecase/dto"
)

type AdminProfileHandler struct {
	profile   *admin.ProfileUsecase
	uidHeader string
}

func NewAdminProfileHandler(profile *admin.ProfileUsecase, uidHeader string) (*AdminProfileHandler, error) {
	return &AdminProfileHandler{
		profile:   profile,
		uidHeader: uidHeader,
	}, nil
}

func (h *AdminProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:user_id", h.GetProfile)
	router.PUT("/:user_id", h.UpdateProfile)
}

func (h *AdminProfileHandler) isAdmin(c *gin.Context) bool {
	return false
}

func (h *AdminProfileHandler) GetProfile(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		err := errors.New("user ID is required", "InvalidRequest", errcode.ErrInvalidInput)
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

func (h *AdminProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		err := errors.New("user ID is required", "InvalidRequest", errcode.ErrInvalidInput)
		c.Error(err)
		return
	}

	userUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.Upgrade(err, "Invalid user ID format", errcode.ErrInvalidInput)
		c.Error(err)
		return
	}

	var updateData dto.UpdateProfileData
	if err := c.ShouldBindJSON(&updateData); err != nil {
		err = errors.Upgrade(err, "Failed to bind update profile data", errcode.ErrInvalidInput)
		c.Error(err)
		return
	}

	updateData.UserID = userUID

	profile, err := h.profile.UpdateProfile(c.Request.Context(), &updateData)
	if err != nil {
		if errors.Is(err, errcode.ErrNotFound) {
			err = errors.Join(err, "Update profile handler failed")
			c.Error(err)
			return
		}
		err = errors.Upgrade(err, "Failed to update profile", errcode.ErrInternalFailure)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, profile)
}
