package repository

import (
	"context"
	"database/sql"
	"portal_link/models"
	"portal_link/modules/user/domain"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 建立新的使用者 Repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 建立使用者
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	m := &models.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: null.TimeFrom(user.CreatedAt),
		UpdatedAt: null.TimeFrom(user.UpdatedAt),
	}

	err := m.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}

	// 回填生成的 ID 和時間戳記
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt.Time
	user.UpdatedAt = m.UpdatedAt.Time

	return nil
}

// GetByEmail 根據 Email 獲取使用者
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

// Find 根據 ID 獲取使用者
func (r *UserRepository) Find(ctx context.Context, id int) (*domain.User, error) {
	user, err := models.FindUser(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}
