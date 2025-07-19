package httphandlerv1

import (
	"github.com/gin-gonic/gin"
	"mandacode.com/accounts/user/internal/usecase/admin"
	manage "mandacode.com/accounts/user/internal/usecase/management"
)

type AdminHandler struct {
	adminUsecase  *admin.AdminUsecase
	manageUsecase *manage.AdminManageUsecase
}

// NewAdminHandler creates a new AdminHandler with the provided use cases.
func NewAdminHandler(adminUsecase *admin.AdminUsecase, manageUsecase *manage.AdminManageUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase:  adminUsecase,
		manageUsecase: manageUsecase,
	}
}

// RegisterRoutes registers the admin routes with the provided router.
func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Define your admin routes here
	// Example:
	// router.GET("/admin/users", h.GetUsers)
	// router.POST("/admin/users", h.CreateUser)
	// router.DELETE("/admin/users/:id", h.DeleteUser)
}
