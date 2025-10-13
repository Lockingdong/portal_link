package domain

import (
	"time"

	"github.com/cockroachdb/errors"
)

type PortalPageParams PortalPage

// PortalPage 實體代表使用者的個人化連結整合頁面
// PortalPage 是聚合根（Aggregate Root），負責管理其內部的所有 Link 實體
type PortalPage struct {
	ID              int
	UserID          int
	Slug            string
	Title           string
	Bio             string
	ProfileImageURL string
	Theme           Theme
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Links           []*Link // 聚合內的 Links
}

// NewPortalPage 建立新的 PortalPage 實體
func NewPortalPage(params PortalPageParams) *PortalPage {
	now := time.Now()

	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}

	if params.UpdatedAt.IsZero() {
		params.UpdatedAt = now
	}

	if params.Links == nil {
		params.Links = make([]*Link, 0)
	}

	portalPage := &PortalPage{
		ID:              params.ID,
		UserID:          params.UserID,
		Slug:            params.Slug,
		Title:           params.Title,
		Bio:             params.Bio,
		ProfileImageURL: params.ProfileImageURL,
		Theme:           params.Theme,
		CreatedAt:       params.CreatedAt,
		UpdatedAt:       params.UpdatedAt,
		Links:           params.Links,
	}

	return portalPage
}

// AddLink 新增 Link 到 PortalPage（聚合根方法）
func (p *PortalPage) AddLink(params LinkParams) {
	link := newLink(p.ID, params)
	p.Links = append(p.Links, link)
	p.UpdatedAt = time.Now()
}

// RemoveLink 從 PortalPage 移除 Link（聚合根方法）
func (p *PortalPage) RemoveLink(linkID int) error {
	for i, link := range p.Links {
		if link.ID == linkID {
			// 移除 Link
			p.Links = append(p.Links[:i], p.Links[i+1:]...)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.Wrapf(ErrLinkNotFound, "link ID: %d", linkID)
}

// UpdateLink 更新 PortalPage 中的 Link（聚合根方法）
func (p *PortalPage) UpdateLink(linkID int, params LinkParams) error {
	for _, link := range p.Links {
		if link.ID == linkID {
			link.Update(params)
			p.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.Wrapf(ErrLinkNotFound, "link ID: %d", linkID)
}
