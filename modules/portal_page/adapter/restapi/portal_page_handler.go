package restapi

import (
	"database/sql"
	"errors"
	"net/http"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"portal_link/modules/portal_page/usecase"
	user_domain "portal_link/modules/user/domain"
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
	findByIDUC         *usecase.FindPortalPageByIDUC
	findBySlugUC       *usecase.FindPortalPageBySlugUC
}

// NewPortalPageHandler 建立新的個人頁面處理器
func NewPortalPageHandler(e *gin.Engine, db *sql.DB, userRepo user_domain.UserRepository) error {
	portalPageRepo := repository.NewPortalPageRepository(db)
	handler := &PortalPageHandler{
		createPortalPageUC: usecase.NewCreatePortalPageUC(portalPageRepo),
		updatePortalPageUC: usecase.NewUpdatePortalPageUC(portalPageRepo),
		listPortalPagesUC:  usecase.NewListPortalPagesUC(portalPageRepo),
		findByIDUC:         usecase.NewFindPortalPageByIDUC(portalPageRepo),
		findBySlugUC:       usecase.NewFindPortalPageBySlugUC(portalPageRepo),
	}

	e.GET("/api/v1/me/portal-pages", pkg.AuthMiddleware(userRepo), handler.ListPortalPages)
	e.GET("/api/v1/me/portal-pages/:id", pkg.AuthMiddleware(userRepo), handler.FindPortalPageByID)
	e.POST("/api/v1/me/portal-pages", pkg.AuthMiddleware(userRepo), handler.CreatePortalPage)
	e.PUT("/api/v1/me/portal-pages/:id", pkg.AuthMiddleware(userRepo), handler.UpdatePortalPage)

	// Public endpoint: find portal page by slug (no auth)
	e.GET("/api/v1/portal-pages/:slug", handler.FindPortalPageBySlug)

	return nil
}

// NewInMemPortalPageHandler 建立新的個人頁面處理器 (in-memory version)
func NewInMemPortalPageHandler(e *gin.Engine, userRepo user_domain.UserRepository) error {
	portalPageRepo := repository.NewInMemoryPortalPageRepository()
	handler := &PortalPageHandler{
		createPortalPageUC: usecase.NewCreatePortalPageUC(portalPageRepo),
		updatePortalPageUC: usecase.NewUpdatePortalPageUC(portalPageRepo),
		listPortalPagesUC:  usecase.NewListPortalPagesUC(portalPageRepo),
		findByIDUC:         usecase.NewFindPortalPageByIDUC(portalPageRepo),
		findBySlugUC:       usecase.NewFindPortalPageBySlugUC(portalPageRepo),
	}

	e.GET("/api/v1/me/portal-pages", pkg.AuthMiddleware(userRepo), handler.ListPortalPages)
	e.GET("/api/v1/me/portal-pages/:id", pkg.AuthMiddleware(userRepo), handler.FindPortalPageByID)
	e.POST("/api/v1/me/portal-pages", pkg.AuthMiddleware(userRepo), handler.CreatePortalPage)
	e.PUT("/api/v1/me/portal-pages/:id", pkg.AuthMiddleware(userRepo), handler.UpdatePortalPage)

	// Public endpoint: find portal page by slug (no auth)
	e.GET("/api/v1/portal-pages/:slug", handler.FindPortalPageBySlug)

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
	result, err := h.createPortalPageUC.Execute(c, &params)
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
	result, err := h.updatePortalPageUC.Execute(c, &params)
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

// FindPortalPageByID 取得單一 Portal Page（含 Links）
func (h *PortalPageHandler) FindPortalPageByID(c *gin.Context) {
	// 取得已驗證的使用者 ID
	userID, err := pkg.GetUserIDFromContext(c)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 取得路徑參數中的 Portal Page ID
	portalPageIDStr := c.Param("id")
	portalPageID, err := strconv.Atoi(portalPageIDStr)
	if err != nil || portalPageID <= 0 {
		http_error.ResponseBadRequest(c, nil)
		return
	}

	// 將 userID 從字符串轉換為整數
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http_error.ResponseInternalServerError(c, nil)
		return
	}

	// 執行查詢邏輯
	params := &usecase.FindPortalPageByIDParams{
		UserID: userIDInt,
		ID:     portalPageID,
	}
	result, err := h.findByIDUC.Execute(c, params)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidParams) {
			http_error.ResponseBadRequest(c, &http_error.ErrorResponse{Message: err.Error()})
			return
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			http_error.ResponseForbidden(c, &http_error.ErrorResponse{Message: err.Error()})
			return
		}
		if errors.Is(err, domain.ErrPortalPageNotFound) {
			http_error.ResponseNotFound(c, &http_error.ErrorResponse{Message: err.Error()})
			return
		}

		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{Message: err.Error()})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, result)
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
	result, err := h.listPortalPagesUC.Execute(c, params)
	if err != nil {
		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, result)
}

// FindPortalPageBySlug 透過 slug 取得公開的 Portal Page（含 Links）
func (h *PortalPageHandler) FindPortalPageBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// 執行查詢邏輯
	params := &usecase.FindPortalPageBySlugParams{Slug: slug}
	result, err := h.findBySlugUC.Execute(c, params)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidParams) {
			http_error.ResponseBadRequest(c, &http_error.ErrorResponse{Message: err.Error()})
			return
		}
		if errors.Is(err, domain.ErrPortalPageNotFound) {
			http_error.ResponseNotFound(c, &http_error.ErrorResponse{Message: err.Error()})
			return
		}

		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{Message: err.Error()})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, result)
}
