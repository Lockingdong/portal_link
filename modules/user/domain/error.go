package domain

import "github.com/cockroachdb/errors"

var (
	// ErrInvalidParams 參數錯誤
	ErrInvalidParams = errors.New("invalid parameters")

	// ErrEmailExists Email 已存在於系統
	ErrEmailExists = errors.New("email already exists")

	// ErrInvalidCredentials 登入憑證錯誤（帳號或密碼錯誤）
	ErrInvalidCredentials = errors.New("invalid credentials")
)
