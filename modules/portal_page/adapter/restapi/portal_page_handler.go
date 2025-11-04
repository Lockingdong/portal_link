package restapi

import (
	user_domain "portal_link/modules/user/domain"

	"github.com/gin-gonic/gin"
)

// PortalPageHandler 個人頁面處理器
type PortalPageHandler struct {
}

// NewInMemPortalPageHandler 建立新的個人頁面處理器 (in-memory version)
func NewInMemPortalPageHandler(e *gin.Engine, userRepo user_domain.UserRepository) error {
	return nil
}
