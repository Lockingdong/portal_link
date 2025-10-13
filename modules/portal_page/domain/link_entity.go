package domain

import "time"

// Link 實體代表使用者在 Portal Page 中展示的個別連結項目
// Link 是 Portal Page 聚合內的實體，必須透過 Portal Page（聚合根）來管理
type Link struct {
	ID           int
	PortalPageID int
	Title        string
	URL          string
	Description  string
	IconURL      string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// LinkParams 用於建立或更新 Link 的參數
type LinkParams Link

// newLink 建立新的 Link 實體（私有方法，只能透過 PortalPage 聚合根調用）
func newLink(portalPageID int, params LinkParams) *Link {
	now := time.Now()

	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}

	if params.UpdatedAt.IsZero() {
		params.UpdatedAt = now
	}

	return &Link{
		ID:           params.ID,
		PortalPageID: portalPageID,
		Title:        params.Title,
		URL:          params.URL,
		Description:  params.Description,
		IconURL:      params.IconURL,
		DisplayOrder: params.DisplayOrder,
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.UpdatedAt,
	}
}

// Update 更新 Link 的資訊
func (l *Link) Update(params LinkParams) {
	if params.Title != "" {
		l.Title = params.Title
	}
	if params.URL != "" {
		l.URL = params.URL
	}
	if params.Description != "" {
		l.Description = params.Description
	}
	if params.IconURL != "" {
		l.IconURL = params.IconURL
	}
	if params.DisplayOrder > 0 {
		l.DisplayOrder = params.DisplayOrder
	}
	l.UpdatedAt = time.Now()
}
