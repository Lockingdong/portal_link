package domain

import "context"

// UserRepository 使用者 Repository
type UserRepository interface {
	// Create 建立使用者
	Create(ctx context.Context, user *User) error

	// GetByEmail 根據 Email 獲取使用者
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Find 根據 ID 獲取使用者
	Find(ctx context.Context, id int) (*User, error)
}
