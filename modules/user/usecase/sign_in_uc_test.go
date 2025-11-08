package usecase

import (
	"context"
	"portal_link/modules/user/domain"
	"portal_link/modules/user/repository"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
)

func TestSignInUC_Execute(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	ctx := context.Background()

	tests := []struct {
		name           string
		params         *SignInParams
		setupData      func(t *testing.T) // 準備測試數據
		wantErr        bool
		expectedErr    error
		expectedErrMsg string
		checkResult    func(t *testing.T, result *SignInResult)
	}{
		{
			name: "成功登入",
			params: &SignInParams{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				// 先創建一個使用者
				existingUser, err := domain.NewUser(domain.UserParams{
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: "password123",
				})
				assert.NoError(t, err)
				err = repo.Create(ctx, existingUser)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignInResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			},
		},
		{
			name: "Email 為空",
			params: &SignInParams{
				Email:    "",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 格式錯誤（缺少 @）",
			params: &SignInParams{
				Email:    "johnexample.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 格式錯誤（缺少域名）",
			params: &SignInParams{
				Email:    "john@",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 太長（超過 255 字元）",
			params: &SignInParams{
				Email:    strings.Repeat("a", 247) + "@test.com", // 247 + 9 = 256 字元
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Password 太短（少於 8 字元）",
			params: &SignInParams{
				Email:    "test@example.com",
				Password: "pass123",
			},
			wantErr:        true,
			expectedErrMsg: "password is invalid",
		},
		{
			name: "使用者不存在",
			params: &SignInParams{
				Email:    "notexist@example.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErr:    domain.ErrInvalidCredentials,
			expectedErrMsg: "invalid credentials",
		},
		{
			name: "密碼錯誤",
			params: &SignInParams{
				Email:    "wrongpass@example.com",
				Password: "wrongpassword123",
			},
			setupData: func(t *testing.T) {
				// 先創建一個使用者
				existingUser, err := domain.NewUser(domain.UserParams{
					Name:     "Wrong Pass User",
					Email:    "wrongpass@example.com",
					Password: "correctpassword123",
				})
				assert.NoError(t, err)
				err = repo.Create(ctx, existingUser)
				assert.NoError(t, err)
			},
			wantErr:        true,
			expectedErr:    domain.ErrInvalidCredentials,
			expectedErrMsg: "invalid credentials",
		},
		{
			name: "Email 有大寫字母（有效）",
			params: &SignInParams{
				Email:    "Upper.Case@Example.COM",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				// 先創建一個使用者
				existingUser, err := domain.NewUser(domain.UserParams{
					Name:     "Upper Case User",
					Email:    "Upper.Case@Example.COM",
					Password: "password123",
				})
				assert.NoError(t, err)
				err = repo.Create(ctx, existingUser)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignInResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			},
		},
		{
			name: "Password 包含特殊字元（有效）",
			params: &SignInParams{
				Email:    "special@example.com",
				Password: "pass@word123!",
			},
			setupData: func(t *testing.T) {
				// 先創建一個使用者
				existingUser, err := domain.NewUser(domain.UserParams{
					Name:     "Special User",
					Email:    "special@example.com",
					Password: "pass@word123!",
				})
				assert.NoError(t, err)
				err = repo.Create(ctx, existingUser)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignInResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			},
		},
		{
			name: "密碼長度剛好 8 字元（有效）",
			params: &SignInParams{
				Email:    "minpass@example.com",
				Password: "12345678",
			},
			setupData: func(t *testing.T) {
				// 先創建一個使用者
				existingUser, err := domain.NewUser(domain.UserParams{
					Name:     "Min Pass User",
					Email:    "minpass@example.com",
					Password: "12345678",
				})
				assert.NoError(t, err)
				err = repo.Create(ctx, existingUser)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignInResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 準備測試數據
			if tt.setupData != nil {
				tt.setupData(t)
			}

			uc := NewSignInUC(repo)
			result, err := uc.Execute(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.True(t, errors.Is(err, tt.expectedErr))
				}
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}
		})
	}
}
