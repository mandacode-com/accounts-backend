package httphandlerv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	"mandacode.com/accounts/user/internal/usecase/user"
)

type UserHandler struct {
	userUsecase *user.UserUsecase
	uidHeader   string
	logger      *zap.Logger
}

// NewUserHandler creates a new UserHandler with the provided use case.
func NewUserHandler(userUsecase *user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// RegisterRoutes registers the user routes with the provided router.
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", h.GetUser)
	router.DELETE("/", h.DeleteUser)
}

// GetUser handles the retrieval of a user by their ID.
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userID := ctx.GetHeader(h.uidHeader)
	if userID == "" {
		err := errors.New("user ID is required", "Unauthorized", errcode.ErrUnauthorized)
		h.logger.Error("GetUser request failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  errcode.ErrInvalidInput,
		})
		return
	}
	// Convert userID to UUID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
			"code":  errcode.ErrInvalidInput,
		})
		return
	}

	user, err := h.userUsecase.GetUserByID(ctx, parsedID)
	if err != nil {
		ctx.Error(err)
		if appErr, ok := err.(*errors.AppError); ok {
			ctx.JSON(errcode.MapCodeToHTTP(appErr.Code()), gin.H{
				"error": appErr.Public(),
				"code":  appErr.Code(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  errcode.ErrInternalFailure,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// DeleteUser handles the deletion of a user by their ID.
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userID := ctx.GetHeader(h.uidHeader)
	if userID == "" {
		err := errors.New("User ID is required", "Unauthorized", errcode.ErrUnauthorized)
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  errcode.ErrInvalidInput,
		})
		return

	}
	// Convert userID to UUID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
			"code":  errcode.ErrInvalidInput,
		})
		return
	}

	archivedUser, err := h.userUsecase.ArchiveUser(ctx, parsedID)
	if err != nil {
		ctx.Error(err)
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Failed to archive user", zap.Error(err))
			ctx.JSON(errcode.MapCodeToHTTP(appErr.Code()), gin.H{
				"error": appErr.Public(),
				"code":  appErr.Code(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  errcode.ErrInternalFailure,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User archived successfully",
		"user":    archivedUser,
	})
}
