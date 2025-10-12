package usecase

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"portal_link/modules/user/domain"
	"regexp"
	"time"

	"github.com/cockroachdb/errors"
)

// SignUpParams 註冊用例的輸入參數
type SignUpParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResult 註冊用例的輸出結果
type SignUpResult struct {
	AccessToken string `json:"access_token"`
}

// SignUpUC 註冊用例
type SignUpUC struct {
	userRepository domain.UserRepository
}

func NewSignUpUC(userRepository domain.UserRepository) *SignUpUC {
	return &SignUpUC{userRepository: userRepository}
}

func (s *SignUpUC) Execute(ctx context.Context, signUpParams *SignUpParams) (*SignUpResult, error) {
	// 1. 驗證輸入參數格式
	if err := s.validateParams(signUpParams); err != nil {
		return nil, err
	}

	// 2. 檢查電子郵件地址是否已被註冊
	existingUser, err := s.userRepository.GetByEmail(ctx, signUpParams.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrEmailExists
	}
	// 如果 err 不是 sql.ErrNoRows，則返回錯誤
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// 3. 建立新的 User 實體
	user, err := domain.NewUser(domain.UserParams{
		Name:     signUpParams.Name,
		Email:    signUpParams.Email,
		Password: signUpParams.Password, // 暫時以明文存儲
	})
	if err != nil {
		return nil, err
	}

	// 4. 將使用者資訊存入資料庫
	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	// 5. 產生該 User 的 access_token
	accessToken := generateAccessToken(user.ID)

	// 6. 返回 access_token
	return &SignUpResult{
		AccessToken: accessToken,
	}, nil
}

// validateParams 驗證輸入參數
func (s *SignUpUC) validateParams(params *SignUpParams) error {
	// 驗證 name
	if len(params.Name) < 1 || len(params.Name) > 255 {
		return errors.Wrap(domain.ErrInvalidParams, "name is invalid")
	}

	// 驗證 email
	if len(params.Email) < 1 || len(params.Email) > 255 {
		return errors.Wrap(domain.ErrInvalidParams, "email is invalid")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(params.Email) {
		return errors.Wrap(domain.ErrInvalidParams, "email is invalid")
	}

	// 驗證 password：最少 8 字元，需包含英文和數字
	if len(params.Password) < 8 {
		return errors.Wrap(domain.ErrInvalidParams, "password is invalid")
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(params.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(params.Password)
	if !hasLetter || !hasNumber {
		return errors.Wrap(domain.ErrInvalidParams, "password is invalid")
	}

	return nil
}

// generateAccessToken 產生 access token
// 使用 user id + 過期時間 timestamp，再進行 base64 encode
// 過期時間：1 天
func generateAccessToken(userID int) string {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	tokenString := fmt.Sprintf("%d:%d", userID, expiresAt)
	return base64.StdEncoding.EncodeToString([]byte(tokenString))
}
