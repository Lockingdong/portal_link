package restapi

import (
	"database/sql"
	"net/http"
	"portal_link/modules/user/domain"
	"portal_link/modules/user/repository"
	"portal_link/modules/user/usecase"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
)

// UserHandler 用戶處理器
type UserHandler struct {
	signUpUC *usecase.SignUpUC
	signInUC *usecase.SignInUC
}

// NewUserHandler 建立新的用戶處理器
func NewUserHandler(e *gin.Engine, db *sql.DB) error {
	userRepo := repository.NewUserRepository(db)
	handler := &UserHandler{
		signUpUC: usecase.NewSignUpUC(userRepo),
		signInUC: usecase.NewSignInUC(userRepo),
	}

	router := e.Group("/api/v1/user")
	{
		router.POST("/signup", handler.SignUp)
		router.POST("/signin", handler.SignIn)
	}
	return nil
}

// ErrorResponse API 錯誤響應結構
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SignUp 處理用戶註冊請求
// POST /api/v1/user/signup
func (h *UserHandler) SignUp(c *gin.Context) {
	var req usecase.SignUpParams

	// 綁定並驗證請求體
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, err)
		return
	}

	// 執行註冊用例
	result, err := h.signUpUC.Execute(c.Request.Context(), &usecase.SignUpParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		h.handleError(c, err)
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, &usecase.SignUpResult{
		AccessToken: result.AccessToken,
	})
}

// SignIn 處理用戶登入請求
// POST /api/v1/user/signin
func (h *UserHandler) SignIn(c *gin.Context) {
	var req usecase.SignInParams

	// 綁定並驗證請求體
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, err)
		return
	}

	// 執行登入用例
	result, err := h.signInUC.Execute(c.Request.Context(), &usecase.SignInParams{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		h.handleError(c, err)
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, &usecase.SignInResult{
		AccessToken: result.AccessToken,
	})
}

// handleError 處理用例層錯誤並映射到 HTTP 狀態碼
func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidParams):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ErrInvalidParams",
			Message: "Invalid request parameters",
		})
	case errors.Is(err, domain.ErrEmailExists):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ErrEmailExists",
			Message: "Email already exists",
		})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "ErrInvalidCredentials",
			Message: "Invalid email or password",
		})
	default:
		// 記錄未預期的錯誤（實際應用中應使用 logger）
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "ErrInternal",
			Message: "Internal server error",
		})
	}
}
