package domain

import "github.com/cockroachdb/errors"

var (
	// ErrInvalidParams 參數驗證失敗
	ErrInvalidParams = errors.New("invalid parameters")

	// ErrSlugExists Slug 已被使用
	ErrSlugExists = errors.New("slug already exists")

	// ErrLinkNotFound 找不到指定的 Link
	ErrLinkNotFound = errors.New("link not found")
)
