package domain

import "context"

// PortalPageRepository Portal Page Repository
type PortalPageRepository interface {
	// Create 建立 Portal Page
	Create(ctx context.Context, portalPage *PortalPage) error

	// Update 更新 Portal Page
	// 流程：
	// 1. 查找現有的 Portal Page
	// 2. 更新 Portal Page 的欄位
	// 3. 更新 Portal Page 的 Links
	// 4. 刪除不存在於新的 Links 中的舊 Links
	Update(ctx context.Context, portalPage *PortalPage) error

	// FindBySlug 根據 Slug 查找 Portal Page
	// 依照 display_order 升冪排序
	FindBySlug(ctx context.Context, slug string) (*PortalPage, error)

	// ListByUserID 根據 UserID 查找 Portal Page
	// 不包含 Links
	ListByUserID(ctx context.Context, userID int) ([]*PortalPage, error)

	// FindByID 根據 ID 查找 Portal Page
	// 依照 display_order 升冪排序
	FindByID(ctx context.Context, id int) (*PortalPage, error)
}
