package usecase

import (
	"context"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPortalPageByIDUC_Execute(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := repository.NewPortalPageRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		params         *FindPortalPageByIDParams
		setupData      func(t *testing.T) int // 準備測試數據並返回頁面 ID
		wantErr        bool
		expectedErrMsg string
		checkResult    func(t *testing.T, result *FindPortalPageByIDResult)
	}{
		{
			name: "成功取得頁面（含 Links，排序升冪）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "john-page",
					Title:           "John's Page",
					Bio:             "hello",
					ProfileImageURL: "https://example.com/p.png",
					Theme:           domain.Theme("light"),
				})
				// 添加亂序 Links，驗證輸出順序
				page.AddLink(domain.LinkParams{Title: "B", URL: "https://b.com", DisplayOrder: 2})
				page.AddLink(domain.LinkParams{Title: "A", URL: "https://a.com", DisplayOrder: 1})
				page.AddLink(domain.LinkParams{Title: "C", URL: "https://c.com", DisplayOrder: 3})

				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params:  &FindPortalPageByIDParams{UserID: 1},
			wantErr: false,
			checkResult: func(t *testing.T, result *FindPortalPageByIDResult) {
				assert.NotNil(t, result)
				assert.Equal(t, "john-page", result.Slug)
				assert.Equal(t, "John's Page", result.Title)
				assert.Equal(t, "light", result.Theme)
				// Links 應依 display_order 升冪
				assert.Len(t, result.Links, 3)
				assert.Equal(t, 1, result.Links[0].DisplayOrder)
				assert.Equal(t, "A", result.Links[0].Title)
				assert.Equal(t, 2, result.Links[1].DisplayOrder)
				assert.Equal(t, "B", result.Links[1].Title)
				assert.Equal(t, 3, result.Links[2].DisplayOrder)
				assert.Equal(t, "C", result.Links[2].Title)
			},
		},
		{
			name:           "頁面不存在",
			setupData:      func(t *testing.T) int { return 0 },
			params:         &FindPortalPageByIDParams{UserID: 1, ID: 999999},
			wantErr:        true,
			expectedErrMsg: "portal page not found",
		},
		{
			name: "非擁有者存取返回 Unauthorized",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          2,
					Slug:            "user2-page",
					Title:           "User2 Page",
					Bio:             "",
					ProfileImageURL: "",
					Theme:           domain.Theme("dark"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params:         &FindPortalPageByIDParams{UserID: 1},
			wantErr:        true,
			expectedErrMsg: "unauthorized",
		},
		{
			name: "無效參數（user_id <= 0）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID: 1, Slug: "p", Title: "t", Theme: domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params:         &FindPortalPageByIDParams{UserID: 0},
			wantErr:        true,
			expectedErrMsg: "invalid parameters",
		},
		{
			name: "無連結時返回空陣列",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID: 1, Slug: "nolinks", Title: "No Links", Theme: domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params:  &FindPortalPageByIDParams{UserID: 1},
			wantErr: false,
			checkResult: func(t *testing.T, result *FindPortalPageByIDResult) {
				assert.NotNil(t, result)
				assert.Equal(t, "nolinks", result.Slug)
				assert.Len(t, result.Links, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 每個測試前清理數據庫
			cleanupTestDB(t, db)

			// 準備資料，獲取 ID
			var pageID int
			if tt.setupData != nil {
				pageID = tt.setupData(t)
			}

			// 設置查詢的 ID
			if tt.params != nil {
				if tt.params.ID == 0 {
					tt.params.ID = pageID
				}
			}

			uc := NewFindPortalPageByIDUC(repo)
			result, err := uc.Execute(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}
		})
	}
}
