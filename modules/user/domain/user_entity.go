package domain

import "time"

type UserParams User

// User 實體代表使用 Portal Link 的使用者
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser 建立新的 User 實體
func NewUser(params UserParams) (*User, error) {
	now := time.Now()

	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}

	if params.UpdatedAt.IsZero() {
		params.UpdatedAt = now
	}

	user := &User{
		ID:        params.ID,
		Name:      params.Name,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	}

	return user, nil
}
