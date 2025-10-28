package usecase

import (
	"context"
	"portal_link/modules/user/domain"
	"portal_link/modules/user/repository"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignInUC_Execute(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	ctx := context.Background()

	// 準備一個測試用戶
	prepareTestUser := func(t *testing.T, email, password string) {
		t.Helper()
		user := &domain.User{
			Name:     "Test User",
			Email:    email,
			Password: password,
		}
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
	}

	tests := []struct {
		name           string
		params         *SignInParams
		setupData      func(t *testing.T) // 準備測試數據
		wantErr        bool
		expectedErrMsg string
		checkResult    func(t *testing.T, result *SignInResult)
	}{
		{
			name: "成功登入",
			params: &SignInParams{
				Email:    "testuser@example.com",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				prepareTestUser(t, "testuser@example.com", "password123")
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
				Email:    "testexample.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 格式錯誤（缺少域名）",
			params: &SignInParams{
				Email:    "test@",
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
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "invalid credentials",
		},
		{
			name: "密碼錯誤",
			params: &SignInParams{
				Email:    "wrongpass@example.com",
				Password: "wrongpassword123",
			},
			setupData: func(t *testing.T) {
				prepareTestUser(t, "wrongpass@example.com", "correctpassword123")
			},
			wantErr:        true,
			expectedErrMsg: "invalid credentials",
		},
		{
			name: "密碼正確但大小寫不同（應失敗）",
			params: &SignInParams{
				Email:    "caseuser@example.com",
				Password: "Password123",
			},
			setupData: func(t *testing.T) {
				prepareTestUser(t, "caseuser@example.com", "password123")
			},
			wantErr:        true,
			expectedErrMsg: "invalid credentials",
		},
		{
			name: "Password 包含特殊字元（有效）",
			params: &SignInParams{
				Email:    "special@example.com",
				Password: "pass@word123!",
			},
			setupData: func(t *testing.T) {
				prepareTestUser(t, "special@example.com", "pass@word123!")
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignInResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			},
		},
		{
			name: "7 字元密碼（無效）",
			params: &SignInParams{
				Email:    "test@example.com",
				Password: "1234567",
			},
			wantErr:        true,
			expectedErrMsg: "password is invalid",
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
