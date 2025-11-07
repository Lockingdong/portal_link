package restapi

import (
	"database/sql"
	"errors"
	"net/http"
	"portal_link/modules/user/domain"
	"portal_link/modules/user/repository"
	"portal_link/modules/user/usecase"
	"portal_link/pkg/http_error"

	"github.com/gin-gonic/gin"
)

// UserHandler 用戶處理器
type UserHandler struct {
	signUpUC *usecase.SignUpUC
	// signInUC *usecase.SignInUC
}

// NewInMemUserHandler 建立新的用戶處理器 (in-memory version)
func NewInMemUserHandler(e *gin.Engine, userRepo domain.UserRepository) error {
	handler := &UserHandler{
		signUpUC: usecase.NewSignUpUC(userRepo),
		// signInUC: usecase.NewSignInUC(userRepo),
	}

	router := e.Group("/api/v1/user")
	{
		router.POST("/signup", handler.SignUp)
		router.POST("/signin", handler.SignIn)
	}
	return nil
}

// NewUserHandler 建立新的用戶處理器
func NewUserHandler(e *gin.Engine, db *sql.DB) error {
	userRepo := repository.NewInMemoryUserRepository()
	handler := &UserHandler{
		signUpUC: usecase.NewSignUpUC(userRepo),
		// signInUC: usecase.NewSignInUC(userRepo),
	}

	router := e.Group("/api/v1/user")
	{
		router.POST("/signup", handler.SignUp)
		router.POST("/signin", handler.SignIn)
	}
	return nil
}

// SignUp 處理用戶註冊請求
func (h *UserHandler) SignUp(c *gin.Context) {
	var req usecase.SignUpParams

	// 綁定並驗證請求體
	if err := c.ShouldBindJSON(&req); err != nil {
		http_error.ResponseBadRequest(c, nil)
		return
	}

	// 執行註冊用例
	result, err := h.signUpUC.Execute(c.Request.Context(), &usecase.SignUpParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidParams) || errors.Is(err, domain.ErrEmailExists) {
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
	c.JSON(http.StatusOK, &usecase.SignUpResult{
		AccessToken: result.AccessToken,
	})
}

// SignIn 處理用戶登入請求
// func (h *UserHandler) SignIn(c *gin.Context) {
// 	var req usecase.SignInParams

// 	// 綁定並驗證請求體
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		http_error.ResponseBadRequest(c, nil)
// 		return
// 	}

// 	// 執行登入用例
// 	result, err := h.signInUC.Execute(c.Request.Context(), &usecase.SignInParams{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	})

// 	if err != nil {
// 		if errors.Is(err, domain.ErrInvalidParams) || errors.Is(err, domain.ErrInvalidCredentials) {
// 			http_error.ResponseBadRequest(c, &http_error.ErrorResponse{
// 				Message: err.Error(),
// 			})
// 			return
// 		}
// 		http_error.ResponseInternalServerError(c, &http_error.ErrorResponse{
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	// 返回成功響應
// 	c.JSON(http.StatusOK, &usecase.SignInResult{
// 		AccessToken: result.AccessToken,
// 	})
// }
