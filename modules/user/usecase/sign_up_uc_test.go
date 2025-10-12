package usecase

import (
	"context"
	"database/sql"
	"portal_link/modules/user/domain"
	"portal_link/modules/user/repository"
	"portal_link/pkg"
	"portal_link/pkg/config"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupTestDB 設置測試數據庫連接
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	config.Init()

	// 使用 config 獲取資料庫連接
	db := pkg.NewPG(config.GetDBConfig().DSN())

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	return db
}

// cleanupTestDB 清理測試數據
func cleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM portal_link.users")
	if err != nil {
		t.Logf("failed to cleanup test data: %v", err)
	}
}

func TestSignUpUC_Execute(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		params         *SignUpParams
		setupData      func(t *testing.T) // 準備測試數據
		wantErr        bool
		expectedErrMsg string
		checkResult    func(t *testing.T, result *SignUpResult)
	}{
		{
			name: "成功註冊",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignUpResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)

				// 驗證用戶已被創建
				user, err := repo.GetByEmail(ctx, "john@example.com")
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john@example.com", user.Email)
				assert.NotZero(t, user.ID)
			},
		},
		{
			name: "Name 為空",
			params: &SignUpParams{
				Name:     "",
				Email:    "test1@example.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "name is invalid",
		},
		{
			name: "Name 太長（超過 255 字元）",
			params: &SignUpParams{
				Name:     strings.Repeat("a", 256),
				Email:    "test2@example.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "name is invalid",
		},
		{
			name: "Email 為空",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 格式錯誤（缺少 @）",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "johnexample.com",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 格式錯誤（缺少域名）",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "john@",
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Email 太長（超過 255 字元）",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    strings.Repeat("a", 247) + "@test.com", // 247 + 9 = 256 字元
				Password: "password123",
			},
			wantErr:        true,
			expectedErrMsg: "email is invalid",
		},
		{
			name: "Password 太短（少於 8 字元）",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "test3@example.com",
				Password: "pass123",
			},
			wantErr:        true,
			expectedErrMsg: "password is invalid",
		},
		{
			name: "Password 沒有字母",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "test4@example.com",
				Password: "12345678",
			},
			wantErr:        true,
			expectedErrMsg: "password is invalid",
		},
		{
			name: "Password 沒有數字",
			params: &SignUpParams{
				Name:     "John Doe",
				Email:    "test5@example.com",
				Password: "password",
			},
			wantErr:        true,
			expectedErrMsg: "password is invalid",
		},
		{
			name: "Email 已存在",
			params: &SignUpParams{
				Name:     "Another User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				// 先創建一個使用者
				existingUser := &domain.User{
					Name:     "Existing User",
					Email:    "existing@example.com",
					Password: "oldpassword",
				}
				err := repo.Create(ctx, existingUser)
				assert.NoError(t, err)
			},
			wantErr:        true,
			expectedErrMsg: "email already exists",
		},
		{
			name: "Password 包含特殊字元（有效）",
			params: &SignUpParams{
				Name:     "Special User",
				Email:    "special@example.com",
				Password: "pass@word123!",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignUpResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)

				// 驗證用戶已被創建
				user, err := repo.GetByEmail(ctx, "special@example.com")
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "Special User", user.Name)
			},
		},
		{
			name: "Email 有大寫字母（有效）",
			params: &SignUpParams{
				Name:     "Upper Case User",
				Email:    "Upper.Case@Example.COM",
				Password: "password123",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *SignUpResult) {
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)

				// 驗證用戶已被創建（注意：email 應該存為原始格式）
				user, err := repo.GetByEmail(ctx, "Upper.Case@Example.COM")
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "Upper Case User", user.Name)
				assert.Equal(t, "Upper.Case@Example.COM", user.Email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 每個測試前清理數據庫
			cleanupTestDB(t, db)

			// 準備測試數據
			if tt.setupData != nil {
				tt.setupData(t)
			}

			uc := NewSignUpUC(repo)
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
