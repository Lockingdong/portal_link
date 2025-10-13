package domain

import "context"

// PortalPageRepository Portal Page Repository
type PortalPageRepository interface {
	// Create 建立 Portal Page
	Create(ctx context.Context, portalPage *PortalPage) error

	// Update 更新 Portal Page
	Update(ctx context.Context, portalPage *PortalPage) error

	// FindBySlug 根據 Slug 查找 Portal Page
	FindBySlug(ctx context.Context, slug string) (*PortalPage, error)

	// FindByUserID 根據 UserID 查找 Portal Page
	FindByUserID(ctx context.Context, userID int) ([]*PortalPage, error)
}
