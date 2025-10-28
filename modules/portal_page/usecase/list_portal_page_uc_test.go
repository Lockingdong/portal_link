package usecase

import (
	"context"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPortalPagesUC_Execute(t *testing.T) {
	repo := repository.NewInMemoryPortalPageRepository()
	ctx := context.Background()

	tests := []struct {
		name        string
		params      *ListPortalPagesParams
		setupData   func(t *testing.T) // 準備測試數據
		wantErr     bool
		checkResult func(t *testing.T, result *ListPortalPagesResult)
	}{
		{
			name: "成功列出用戶的多個 Portal Pages",
			params: &ListPortalPagesParams{
				UserID: 1,
			},
			setupData: func(t *testing.T) {
				// 創建多個 Portal Pages
				page1 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "page-one",
					Title:           "First Page",
					Bio:             "This is the first page",
					ProfileImageURL: "https://example.com/image1.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page1)
				assert.NoError(t, err)

				page2 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "page-two",
					Title:           "Second Page",
					Bio:             "This is the second page",
					ProfileImageURL: "https://example.com/image2.jpg",
					Theme:           domain.Theme("dark"),
				})
				err = repo.Create(ctx, page2)
				assert.NoError(t, err)

				page3 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "page-three",
					Title:           "Third Page",
					Bio:             "This is the third page",
					ProfileImageURL: "https://example.com/image3.jpg",
					Theme:           domain.Theme("light"),
				})
				err = repo.Create(ctx, page3)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *ListPortalPagesResult) {
				assert.NotNil(t, result)
				assert.Len(t, result.PortalPages, 3)

				// 驗證返回的數據只包含摘要信息 (ID, Slug, Title)
				for _, pp := range result.PortalPages {
					assert.NotZero(t, pp.ID)
					assert.NotEmpty(t, pp.Slug)
					assert.NotEmpty(t, pp.Title)
				}

				// 驗證具體的頁面資訊
				slugs := make(map[string]bool)
				titles := make(map[string]bool)
				for _, pp := range result.PortalPages {
					slugs[pp.Slug] = true
					titles[pp.Title] = true
				}
				assert.True(t, slugs["page-one"])
				assert.True(t, slugs["page-two"])
				assert.True(t, slugs["page-three"])
				assert.True(t, titles["First Page"])
				assert.True(t, titles["Second Page"])
				assert.True(t, titles["Third Page"])
			},
		},
		{
			name: "用戶沒有任何 Portal Pages 時返回空列表",
			params: &ListPortalPagesParams{
				UserID: 999, // 不存在的用戶 ID
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *ListPortalPagesResult) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.PortalPages)
				assert.Len(t, result.PortalPages, 0)
			},
		},
		{
			name: "用戶只有一個 Portal Page",
			params: &ListPortalPagesParams{
				UserID: 2,
			},
			setupData: func(t *testing.T) {
				// 創建一個 Portal Page
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          2,
					Slug:            "single-page",
					Title:           "Only Page",
					Bio:             "This is the only page",
					ProfileImageURL: "https://example.com/single.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *ListPortalPagesResult) {
				assert.NotNil(t, result)
				assert.Len(t, result.PortalPages, 1)
				assert.Equal(t, "single-page", result.PortalPages[0].Slug)
				assert.Equal(t, "Only Page", result.PortalPages[0].Title)
				assert.NotZero(t, result.PortalPages[0].ID)
			},
		},
		{
			name: "只返回指定用戶的 Portal Pages（不包含其他用戶的）",
			params: &ListPortalPagesParams{
				UserID: 3,
			},
			setupData: func(t *testing.T) {
				// 創建用戶 3 的 Portal Pages
				page1 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          3,
					Slug:            "user3-page1",
					Title:           "User 3 Page 1",
					Bio:             "User 3's first page",
					ProfileImageURL: "https://example.com/user3-1.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page1)
				assert.NoError(t, err)

				page2 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          3,
					Slug:            "user3-page2",
					Title:           "User 3 Page 2",
					Bio:             "User 3's second page",
					ProfileImageURL: "https://example.com/user3-2.jpg",
					Theme:           domain.Theme("dark"),
				})
				err = repo.Create(ctx, page2)
				assert.NoError(t, err)

				// 創建其他用戶的 Portal Pages（不應該被列出）
				otherUserPage := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          4,
					Slug:            "user4-page",
					Title:           "User 4 Page",
					Bio:             "User 4's page",
					ProfileImageURL: "https://example.com/user4.jpg",
					Theme:           domain.Theme("light"),
				})
				err = repo.Create(ctx, otherUserPage)
				assert.NoError(t, err)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *ListPortalPagesResult) {
				assert.NotNil(t, result)
				assert.Len(t, result.PortalPages, 2)

				// 驗證只返回用戶 3 的頁面
				for _, pp := range result.PortalPages {
					assert.Contains(t, []string{"user3-page1", "user3-page2"}, pp.Slug)
					assert.NotContains(t, pp.Slug, "user4")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 準備測試數據
			if tt.setupData != nil {
				tt.setupData(t)
			}

			uc := NewListPortalPagesUC(repo)
			result, err := uc.Execute(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}
		})
	}
}
