package repository

import (
	"context"
	"database/sql"
	"portal_link/modules/user/domain"
	"portal_link/pkg"
	"portal_link/pkg/config"
	"testing"
	"time"

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

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewUserRepository(db)
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
		errMsg  string
		check   func(t *testing.T, user *domain.User)
	}{
		{
			name: "成功建立使用者",
			user: &domain.User{
				Name:      "測試用戶",
				Email:     "test@example.com",
				Password:  "hashed_password",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
			check: func(t *testing.T, user *domain.User) {
				assert.NotZero(t, user.ID)
				assert.Equal(t, user.Name, "測試用戶")
				assert.Equal(t, user.Email, "test@example.com")
			},
		},
		{
			name: "成功建立另一個使用者",
			user: &domain.User{
				Name:      "另一個用戶",
				Email:     "another@example.com",
				Password:  "another_password",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
			check: func(t *testing.T, user *domain.User) {
				assert.Equal(t, user.Name, "另一個用戶")
			},
		},
		{
			name: "失敗 - Email 重複",
			user: &domain.User{
				Name:      "重複用戶",
				Email:     "test@example.com", // 與第一個測試重複
				Password:  "password",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "duplicate key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if err != nil {
				assert.Error(t, err)
				return
			}

			if tt.check != nil {
				tt.check(t, tt.user)
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewUserRepository(db)
	ctx := context.Background()
	now := time.Now()

	// 準備測試數據
	testUser := &domain.User{
		Name:      "測試用戶",
		Email:     "getbyemail@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(ctx, testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
		check   func(t *testing.T, user *domain.User)
	}{
		{
			name:    "成功根據 Email 獲取使用者",
			email:   "getbyemail@example.com",
			wantErr: false,
			check: func(t *testing.T, user *domain.User) {
				assert.NotNil(t, user)
				assert.Equal(t, user.Email, "getbyemail@example.com")
				assert.Equal(t, user.Name, "測試用戶")
				assert.NotZero(t, user.ID)
			},
		},
		{
			name:    "失敗 - Email 不存在",
			email:   "notexist@example.com",
			wantErr: true,
			check: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByEmail(ctx, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.check != nil {
				tt.check(t, user)
			}
		})
	}
}

func TestUserRepository_Find(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewUserRepository(db)
	ctx := context.Background()
	now := time.Now()

	// 準備測試數據
	testUser := &domain.User{
		Name:      "查找用戶",
		Email:     "find@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(ctx, testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
		check   func(t *testing.T, user *domain.User)
	}{
		{
			name:    "成功根據 ID 獲取使用者",
			id:      testUser.ID,
			wantErr: false,
			check: func(t *testing.T, user *domain.User) {
				assert.NotNil(t, user)
				assert.Equal(t, user.ID, testUser.ID)
				assert.Equal(t, user.Email, "find@example.com")
				assert.Equal(t, user.Name, "查找用戶")
			},
		},
		{
			name:    "失敗 - ID 不存在",
			id:      99999,
			wantErr: true,
			check: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.Find(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.check != nil {
				tt.check(t, user)
			}
		})
	}
}
