package usecase

import (
	"context"
	"database/sql"
	"portal_link/modules/user/domain"
	"regexp"

	"github.com/cockroachdb/errors"
)

// SignInParams 登入用例的輸入參數
type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInResult 登入用例的輸出結果
type SignInResult struct {
	AccessToken string `json:"access_token"`
}

// SignInUC 登入用例
type SignInUC struct {
	userRepository domain.UserRepository
}

func NewSignInUC(userRepository domain.UserRepository) *SignInUC {
	return &SignInUC{userRepository: userRepository}
}

func (s *SignInUC) Execute(ctx context.Context, signInParams *SignInParams) (*SignInResult, error) {
	// 1. 驗證輸入參數格式
	if err := s.validateParams(signInParams); err != nil {
		return nil, err
	}

	// 2. 根據電子郵件地址查詢使用者
	user, err := s.userRepository.GetByEmail(ctx, signInParams.Email)
	if err != nil {
		// 如果是找不到使用者，返回 ErrInvalidCredentials（不透露具體原因）
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		// 其他錯誤直接返回
		return nil, err
	}

	// 3. 驗證密碼是否正確（暫時以明文比對）
	if user.Password != signInParams.Password {
		return nil, domain.ErrInvalidCredentials
	}

	// 4. 產生該 User 的 access_token
	accessToken := generateAccessToken(user.ID)

	// 5. 返回 access_token
	return &SignInResult{
		AccessToken: accessToken,
	}, nil
}

// validateParams 驗證輸入參數
func (s *SignInUC) validateParams(params *SignInParams) error {
	// 驗證 email
	if len(params.Email) < 1 || len(params.Email) > 255 {
		return errors.Wrap(domain.ErrInvalidParams, "email is invalid")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(params.Email) {
		return errors.Wrap(domain.ErrInvalidParams, "email is invalid")
	}

	// 驗證 password：最少 8 字元
	if len(params.Password) < 8 {
		return errors.Wrap(domain.ErrInvalidParams, "password is invalid")
	}

	return nil
}
