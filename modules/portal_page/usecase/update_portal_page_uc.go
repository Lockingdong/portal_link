package usecase

import (
	"context"
	"database/sql"
	"net/url"
	"portal_link/modules/portal_page/domain"
	"strings"

	"github.com/cockroachdb/errors"
)

// UpdatePortalPageParams 更新 Portal Page 用例的輸入參數
type UpdatePortalPageParams struct {
	UserID          int    `json:"user_id"`
	PortalPageID    int    `json:"portal_page_id"`
	Slug            string `json:"slug,omitempty"`
	Title           string `json:"title,omitempty"`
	Bio             string `json:"bio,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Theme           string `json:"theme,omitempty"`
	Links           []Link `json:"links"`
}

// Link 連結的輸入參數
type Link struct {
	ID           int    `json:"id,omitempty"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Description  string `json:"description,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	DisplayOrder int    `json:"display_order"`
}

// UpdatePortalPageResult 更新 Portal Page 用例的輸出結果
type UpdatePortalPageResult struct {
	ID int `json:"id"`
}

// UpdatePortalPageUC 更新 Portal Page 用例
type UpdatePortalPageUC struct {
	portalPageRepository domain.PortalPageRepository
}

func NewUpdatePortalPageUC(portalPageRepository domain.PortalPageRepository) *UpdatePortalPageUC {
	return &UpdatePortalPageUC{portalPageRepository: portalPageRepository}
}

func (u *UpdatePortalPageUC) Execute(ctx context.Context, params *UpdatePortalPageParams) (*UpdatePortalPageResult, error) {
	// 1. 從資料庫讀取現有的 PortalPage 實體
	portalPage, err := u.portalPageRepository.FindByID(ctx, params.PortalPageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPortalPageNotFound
		}
		return nil, err
	}

	// 2. 驗證使用者權限
	if portalPage.UserID != params.UserID {
		return nil, domain.ErrUnauthorized
	}

	// 3. 驗證輸入參數格式
	if err := u.validateParams(params, portalPage); err != nil {
		return nil, err
	}

	// 4. 如果要更新 slug，檢查新的 slug 是否已被使用（排除當前頁面）
	if params.Slug != "" && params.Slug != portalPage.Slug {
		existingPortalPage, err := u.portalPageRepository.FindBySlug(ctx, params.Slug)
		if err == nil && existingPortalPage != nil && existingPortalPage.ID != params.PortalPageID {
			return nil, domain.ErrSlugExists
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		portalPage.Slug = strings.ToLower(params.Slug)
	}

	// 5. 更新 PortalPage 實體的欄位
	if params.Title != "" {
		portalPage.Title = params.Title
	}
	if params.Bio != "" {
		portalPage.Bio = params.Bio
	}
	if params.ProfileImageURL != "" {
		portalPage.ProfileImageURL = params.ProfileImageURL
	}
	if params.Theme != "" {
		portalPage.Theme = domain.Theme(params.Theme)
	}

	// 6. 處理 links 的更新
	portalPage.Links = make([]*domain.Link, 0, len(params.Links))
	for _, link := range params.Links {
		portalPage.AddLink(domain.LinkParams{
			ID:           link.ID,
			Title:        link.Title,
			URL:          link.URL,
			Description:  link.Description,
			IconURL:      link.IconURL,
			DisplayOrder: link.DisplayOrder,
		})
	}

	// 7. 透過 Repository 將更新後的 Portal Page 存入資料庫
	if err := u.portalPageRepository.Update(ctx, portalPage); err != nil {
		return nil, err
	}

	// 8. 返回更新後的 PortalPage 實體
	return &UpdatePortalPageResult{
		ID: portalPage.ID,
	}, nil
}

// validateParams 驗證輸入參數
func (u *UpdatePortalPageUC) validateParams(params *UpdatePortalPageParams, existingPage *domain.PortalPage) error {
	// 驗證 slug（如有更新）
	if params.Slug != "" {
		if err := validateSlug(params.Slug); err != nil {
			return err
		}
	}

	// 驗證 title（如有更新）
	if params.Title != "" {
		if err := validateTitle(params.Title); err != nil {
			return err
		}
	}

	// 驗證 bio（如有更新）
	if params.Bio != "" {
		if err := validateBio(params.Bio); err != nil {
			return errors.Wrap(domain.ErrInvalidParams, "bio must not exceed 500 characters")
		}
	}

	// 驗證 profile_image_url（如有更新）
	if params.ProfileImageURL != "" {
		if _, err := url.ParseRequestURI(params.ProfileImageURL); err != nil {
			return errors.Wrap(domain.ErrInvalidParams, "profile_image_url must be a valid URL")
		}
	}

	// 驗證 theme（如有更新）
	if params.Theme != "" {
		theme := domain.Theme(params.Theme)
		if theme != domain.ThemeLight && theme != domain.ThemeDark {
			return errors.Wrap(domain.ErrInvalidParams, "theme must be one of: light, dark")
		}
	}

	// 驗證 links（必填）
	if err := u.validateLinks(params.Links); err != nil {
		return err
	}

	return nil
}

// validateLinks 驗證連結清單
func (u *UpdatePortalPageUC) validateLinks(links []Link) error {
	for _, link := range links {
		// 驗證 title
		if len(link.Title) < 1 || len(link.Title) > 100 {
			return errors.Wrap(domain.ErrInvalidParams, "link title must be between 1 and 100 characters")
		}

		// 驗證 URL
		if _, err := url.ParseRequestURI(link.URL); err != nil {
			return errors.Wrap(domain.ErrInvalidParams, "link URL must be a valid URL")
		}

		// 驗證 description（如有提供）
		if link.Description != "" && len(link.Description) > 500 {
			return errors.Wrap(domain.ErrInvalidParams, "link description must not exceed 500 characters")
		}

		// 驗證 icon_url（如有提供）
		if link.IconURL != "" {
			if _, err := url.ParseRequestURI(link.IconURL); err != nil {
				return errors.Wrap(domain.ErrInvalidParams, "link icon_url must be a valid URL")
			}
		}

		// 驗證 display_order（必須為正整數）
		if link.DisplayOrder < 1 {
			return errors.Wrap(domain.ErrInvalidParams, "link display_order must be a positive integer")
		}
	}

	return nil
}
