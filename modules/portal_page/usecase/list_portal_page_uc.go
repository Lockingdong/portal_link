package usecase

import (
	"context"
	"portal_link/modules/portal_page/domain"
)

// ListPortalPagesParams 列出 Portal Pages 用例的輸入參數
type ListPortalPagesParams struct {
	UserID int `json:"user_id"`
}

// PortalPage Portal Page 摘要資訊
type PortalPage struct {
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// ListPortalPagesResult 列出 Portal Pages 用例的輸出結果
type ListPortalPagesResult struct {
	PortalPages []*PortalPage `json:"portal_pages"`
}

// ListPortalPagesUC 列出 Portal Pages 用例
type ListPortalPagesUC struct {
	portalPageRepository domain.PortalPageRepository
}

func NewListPortalPagesUC(portalPageRepository domain.PortalPageRepository) *ListPortalPagesUC {
	return &ListPortalPagesUC{portalPageRepository: portalPageRepository}
}

func (u *ListPortalPagesUC) Execute(ctx context.Context, params *ListPortalPagesParams) (*ListPortalPagesResult, error) {
	// 1. 透過 Repository 查詢該使用者的所有 Portal Pages
	portalPages, err := u.portalPageRepository.ListByUserID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	// 2. 轉換為摘要格式（只返回 id, slug, title）
	summaries := make([]*PortalPage, 0, len(portalPages))
	for _, pp := range portalPages {
		summaries = append(summaries, &PortalPage{
			ID:    pp.ID,
			Slug:  pp.Slug,
			Title: pp.Title,
		})
	}

	// 3. 返回完整的 Portal Pages 列表
	return &ListPortalPagesResult{
		PortalPages: summaries,
	}, nil
}
