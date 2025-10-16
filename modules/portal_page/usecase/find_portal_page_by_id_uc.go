package usecase

import (
	"context"
	"database/sql"
	"portal_link/modules/portal_page/domain"

	"github.com/cockroachdb/errors"
)

// FindPortalPageByIDParams 依照 ID 取得單一 Portal Page（含 Links）的輸入參數
type FindPortalPageByIDParams struct {
	UserID int `json:"user_id"`
	ID     int `json:"id"`
}

// FindPortalPageByIDResult 單一 Portal Page 的完整輸出（含 Links）
type FindPortalPageByIDResult struct {
	ID              int                            `json:"id"`
	Slug            string                         `json:"slug"`
	Title           string                         `json:"title"`
	Bio             string                         `json:"bio"`
	ProfileImageURL string                         `json:"profile_image_url"`
	Theme           string                         `json:"theme"`
	Links           []FindPortalPageByIDResultLink `json:"links"`
}

// FindPortalPageByIDResultLink 連結的輸出參數
type FindPortalPageByIDResultLink struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Description  string `json:"description"`
	IconURL      string `json:"icon_url"`
	DisplayOrder int    `json:"display_order"`
}

// FindPortalPageByIDUC 依照 ID 取得單一 Portal Page（含 Links）的用例
type FindPortalPageByIDUC struct {
	portalPageRepository domain.PortalPageRepository
}

func NewFindPortalPageByIDUC(portalPageRepository domain.PortalPageRepository) *FindPortalPageByIDUC {
	return &FindPortalPageByIDUC{portalPageRepository: portalPageRepository}
}

func (u *FindPortalPageByIDUC) Execute(ctx context.Context, params *FindPortalPageByIDParams) (*FindPortalPageByIDResult, error) {
	// 1. 驗證輸入參數
	if params == nil || params.UserID <= 0 || params.ID <= 0 {
		return nil, domain.ErrInvalidParams
	}

	// 2. 透過 Repository 取得指定 ID 的 Portal Page（含 Links）
	// Links 依照 display_order 升冪排序
	portalPage, err := u.portalPageRepository.FindByID(ctx, params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPortalPageNotFound
		}
		return nil, err
	}

	// 3. 驗證擁有權
	if portalPage.UserID != params.UserID {
		return nil, domain.ErrUnauthorized
	}

	// 4. 組裝輸出結果（返回完整的 Portal Page 資料）
	links := make([]FindPortalPageByIDResultLink, 0, len(portalPage.Links))
	for _, l := range portalPage.Links {
		links = append(links, FindPortalPageByIDResultLink{
			ID:           l.ID,
			Title:        l.Title,
			URL:          l.URL,
			Description:  l.Description,
			IconURL:      l.IconURL,
			DisplayOrder: l.DisplayOrder,
		})
	}

	return &FindPortalPageByIDResult{
		ID:              portalPage.ID,
		Slug:            portalPage.Slug,
		Title:           portalPage.Title,
		Bio:             portalPage.Bio,
		ProfileImageURL: portalPage.ProfileImageURL,
		Theme:           string(portalPage.Theme),
		Links:           links,
	}, nil
}
