package usecase

import (
	"context"
	"database/sql"
	"net/url"
	"portal_link/modules/portal_page/domain"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
)

// CreatePortalPageParams 建立 Portal Page 用例的輸入參數
type CreatePortalPageParams struct {
	UserID          int    `json:"user_id"`
	Slug            string `json:"slug"`
	Title           string `json:"title"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
	Theme           string `json:"theme"` // 從 API 接收字串，內部轉換為 Theme enum
}

// CreatePortalPageResult 建立 Portal Page 用例的輸出結果
type CreatePortalPageResult struct {
	ID int `json:"id"`
}

// CreatePortalPageUC 建立 Portal Page 用例
type CreatePortalPageUC struct {
	portalPageRepository domain.PortalPageRepository
}

// 保留字清單 - 系統保留的 slug 不可使用
var reservedSlugs = map[string]bool{
	"admin":   true,
	"api":     true,
	"static":  true,
	"public":  true,
	"auth":    true,
	"login":   true,
	"signup":  true,
	"help":    true,
	"about":   true,
	"terms":   true,
	"privacy": true,
}

func NewCreatePortalPageUC(portalPageRepository domain.PortalPageRepository) *CreatePortalPageUC {
	return &CreatePortalPageUC{portalPageRepository: portalPageRepository}
}

func (c *CreatePortalPageUC) Execute(ctx context.Context, params *CreatePortalPageParams) (*CreatePortalPageResult, error) {
	// 1. 驗證輸入參數格式
	if err := c.validateParams(params); err != nil {
		return nil, err
	}

	// 2. 檢查 slug 是否已被使用
	existingPortalPage, err := c.portalPageRepository.FindBySlug(ctx, params.Slug)
	if err == nil && existingPortalPage != nil {
		return nil, domain.ErrSlugExists
	}
	// 如果 err 不是 sql.ErrNoRows，則返回錯誤
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// 3. 建立新的 PortalPage 實體
	portalPage := domain.NewPortalPage(domain.PortalPageParams{
		UserID:          params.UserID,
		Slug:            strings.ToLower(params.Slug), // 轉換為小寫
		Title:           params.Title,
		Bio:             params.Bio,
		ProfileImageURL: params.ProfileImageURL,
		Theme:           c.getThemeOrDefault(params.Theme),
	})

	// 4. 透過 Repository 將 Portal Page 存入資料庫
	if err := c.portalPageRepository.Create(ctx, portalPage); err != nil {
		return nil, err
	}

	// 5. 返回建立成功的 PortalPage 實體
	return &CreatePortalPageResult{
		ID: portalPage.ID,
	}, nil
}

// validateParams 驗證輸入參數
func (c *CreatePortalPageUC) validateParams(params *CreatePortalPageParams) error {
	// 驗證 slug
	if err := c.validateSlug(params.Slug); err != nil {
		return err
	}

	// 驗證 title
	if len(params.Title) < 1 || len(params.Title) > 100 {
		return errors.Wrap(domain.ErrInvalidParams, "title must be between 1 and 100 characters")
	}

	// 驗證 bio（選填）
	if len(params.Bio) > 500 {
		return errors.Wrap(domain.ErrInvalidParams, "bio must not exceed 500 characters")
	}

	// 驗證 profile_image_url（選填）
	if params.ProfileImageURL != "" {
		if _, err := url.ParseRequestURI(params.ProfileImageURL); err != nil {
			return errors.Wrap(domain.ErrInvalidParams, "profile_image_url must be a valid URL")
		}
	}

	// 驗證 theme（選填）
	if params.Theme != "" {
		theme := domain.Theme(params.Theme)
		if !theme.IsValid() {
			return errors.Wrap(domain.ErrInvalidParams, "theme must be one of: light, dark")
		}
	}

	return nil
}

// validateSlug 驗證 slug 格式
func (c *CreatePortalPageUC) validateSlug(slug string) error {
	// 長度限制：3-50 字元
	if len(slug) < 3 || len(slug) > 50 {
		return errors.Wrap(domain.ErrInvalidParams, "slug must be between 3 and 50 characters")
	}

	// 格式驗證：只能包含小寫英文字母、數字和連字號（-）
	// 不可以連字號開頭或結尾
	// 不可包含連續的連字號
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	if !slugRegex.MatchString(slug) {
		return errors.Wrap(domain.ErrInvalidParams, "slug can only contain lowercase letters, numbers, and hyphens (not at start/end or consecutive)")
	}

	// 檢查是否為保留字
	if reservedSlugs[strings.ToLower(slug)] {
		return errors.Wrap(domain.ErrInvalidParams, "slug is reserved and cannot be used")
	}

	return nil
}

// getThemeOrDefault 取得 theme 或返回預設值
func (c *CreatePortalPageUC) getThemeOrDefault(theme string) domain.Theme {
	if theme == "" {
		return domain.GetDefaultTheme()
	}
	return domain.Theme(theme)
}
