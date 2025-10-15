package restapi

import (
	"database/sql"
	"errors"
	"net/http"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"portal_link/modules/portal_page/usecase"
	"portal_link/pkg"
	"portal_link/pkg/http_error"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PortalPageHandler 個人頁面處理器
type PortalPageHandler struct {
	createPortalPageUC *usecase.CreatePortalPageUC
	updatePortalPageUC *usecase.UpdatePortalPageUC
	listPortalPagesUC  *usecase.ListPortalPagesUC
}

// NewPortalPageHandler 建立新的個人頁面處理器
func NewPortalPageHandler(e *gin.Engine, db *sql.DB) error {
	portalPageRepo := repository.NewPortalPageRepository(db)
	handler := &PortalPageHandler{
		createPortalPageUC: usecase.NewCreatePortalPageUC(portalPageRepo),
		updatePortalPageUC: usecase.NewUpdatePortalPageUC(portalPageRepo),
		listPortalPagesUC:  usecase.NewListPortalPagesUC(portalPageRepo),
	}

	e.GET("/api/v1/me/portal-pages", pkg.AuthMiddleware(db), handler.ListPortalPages)
	e.POST("/api/v1/me/portal-pages", pkg.AuthMiddleware(db), handler.CreatePortalPage)
	e.PUT("/api/v1/me/portal-pages/:id", pkg.AuthMiddleware(db), handler.UpdatePortalPage)

	return nil
}

// CreatePortalPage 處理創建個人頁面請求
func (h *PortalPageHandler) CreatePortalPage(c *gin.Context) {
	// 取得已驗證的使用者 ID
	userID, err := pkg.GetUserIDFromContext(c)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 綁定請求參數
	var params usecase.CreatePortalPageParams
	if err := c.ShouldBindJSON(&params); err != nil {
		http_error.ResponseBadRequest(c, nil)
		return
	}

	// 執行創建頁面邏輯
	// 將 userID 從字符串轉換為整數
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}
	params.UserID = userIDInt
	result, err := h.createPortalPageUC.Execute(c.Request.Context(), &params)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidParams) || errors.Is(err, domain.ErrSlugExists) {
			http_error.ResponseBadRequest(c, &http_error.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusCreated, &usecase.CreatePortalPageResult{
		ID: result.ID,
	})
}

// UpdatePortalPage 處理更新個人頁面請求
func (h *PortalPageHandler) UpdatePortalPage(c *gin.Context) {
	// 取得已驗證的使用者 ID
	userID, err := pkg.GetUserIDFromContext(c)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 取得路徑參數中的 Portal Page ID
	portalPageIDStr := c.Param("id")
	portalPageID, err := strconv.Atoi(portalPageIDStr)
	if err != nil {
		http_error.ResponseBadRequest(c, &http_error.ErrorResponse{
			Message: "Invalid portal page ID",
		})
		return
	}

	// 綁定請求參數
	var params usecase.UpdatePortalPageParams
	if err := c.ShouldBindJSON(&params); err != nil {
		http_error.ResponseBadRequest(c, nil)
		return
	}

	// 將 userID 從字符串轉換為整數
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 設置 UserID 和 PortalPageID
	params.UserID = userIDInt
	params.PortalPageID = portalPageID

	// 執行更新頁面邏輯
	result, err := h.updatePortalPageUC.Execute(c.Request.Context(), &params)
	if err != nil {
		// 處理各種錯誤情況
		if errors.Is(err, domain.ErrPortalPageNotFound) {
			http_error.ResponseBadRequest(c, &http_error.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, domain.ErrUnauthorized) {
			http_error.ResponseForbidden(c, &http_error.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, domain.ErrInvalidParams) || errors.Is(err, domain.ErrSlugExists) {
			http_error.ResponseBadRequest(c, &http_error.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, &usecase.UpdatePortalPageResult{
		ID: result.ID,
	})
}

// ListPortalPages 處理列出個人頁面請求
func (h *PortalPageHandler) ListPortalPages(c *gin.Context) {
	// 取得已驗證的使用者 ID
	userID, err := pkg.GetUserIDFromContext(c)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 將 userID 從字符串轉換為整數
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 執行列出頁面邏輯
	params := &usecase.ListPortalPagesParams{
		UserID: userIDInt,
	}
	result, err := h.listPortalPagesUC.Execute(c.Request.Context(), params)
	if err != nil {
		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, result)
}
