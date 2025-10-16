package usecase

import (
	"context"
	"database/sql"
	"portal_link/modules/portal_page/domain"

	"github.com/cockroachdb/errors"
)

// FindPortalPageBySlugParams 依照 Slug 取得單一 Portal Page（含 Links）的輸入參數
type FindPortalPageBySlugParams struct {
	Slug string `json:"slug"`
}

// FindPortalPageBySlugResult 單一 Portal Page 的完整輸出（含 Links）
type FindPortalPageBySlugResult struct {
	ID              int                              `json:"id"`
	Slug            string                           `json:"slug"`
	Title           string                           `json:"title"`
	Bio             string                           `json:"bio"`
	ProfileImageURL string                           `json:"profile_image_url"`
	Theme           string                           `json:"theme"`
	Links           []FindPortalPageBySlugResultLink `json:"links"`
}

// FindPortalPageBySlugResultLink 連結的輸出參數
type FindPortalPageBySlugResultLink struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Description  string `json:"description"`
	IconURL      string `json:"icon_url"`
	DisplayOrder int    `json:"display_order"`
}

// FindPortalPageBySlugUC 依照 Slug 取得單一 Portal Page（含 Links）的用例
type FindPortalPageBySlugUC struct {
	portalPageRepository domain.PortalPageRepository
}

func NewFindPortalPageBySlugUC(portalPageRepository domain.PortalPageRepository) *FindPortalPageBySlugUC {
	return &FindPortalPageBySlugUC{portalPageRepository: portalPageRepository}
}

func (u *FindPortalPageBySlugUC) Execute(ctx context.Context, params *FindPortalPageBySlugParams) (*FindPortalPageBySlugResult, error) {
	// 1. 驗證輸入參數
	if params == nil || params.Slug == "" {
		return nil, domain.ErrInvalidParams
	}

	// 2. 透過 Repository 取得指定 Slug 的 Portal Page（含 Links）
	// Links 依照 display_order 升冪排序
	portalPage, err := u.portalPageRepository.FindBySlug(ctx, params.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPortalPageNotFound
		}
		return nil, err
	}

	// 3. 組裝輸出結果（返回完整的 Portal Page 資料）
	links := make([]FindPortalPageBySlugResultLink, 0, len(portalPage.Links))
	for _, l := range portalPage.Links {
		links = append(links, FindPortalPageBySlugResultLink{
			ID:           l.ID,
			Title:        l.Title,
			URL:          l.URL,
			Description:  l.Description,
			IconURL:      l.IconURL,
			DisplayOrder: l.DisplayOrder,
		})
	}

	return &FindPortalPageBySlugResult{
		ID:              portalPage.ID,
		Slug:            portalPage.Slug,
		Title:           portalPage.Title,
		Bio:             portalPage.Bio,
		ProfileImageURL: portalPage.ProfileImageURL,
		Theme:           string(portalPage.Theme),
		Links:           links,
	}, nil
}
