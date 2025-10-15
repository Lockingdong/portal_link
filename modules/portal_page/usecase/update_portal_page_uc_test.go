package usecase

import (
	"context"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatePortalPageUC_Execute(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := repository.NewPortalPageRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		params         *UpdatePortalPageParams
		setupData      func(t *testing.T) int // 準備測試數據，返回已創建的 Portal Page ID
		wantErr        bool
		expectedErrMsg string
		checkResult    func(t *testing.T, result *UpdatePortalPageResult)
	}{
		{
			name: "成功更新 Portal Page 的基本資訊",
			setupData: func(t *testing.T) int {
				// 創建一個 Portal Page
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "john-doe",
					Title:           "John's Page",
					Bio:             "Welcome to my page",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID:          1,
				Title:           "John's Updated Page",
				Bio:             "Updated bio",
				ProfileImageURL: "https://example.com/new-image.jpg",
				Theme:           "dark",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *UpdatePortalPageResult) {
				assert.NotNil(t, result)
				assert.NotZero(t, result.ID)

				// 驗證 Portal Page 已被更新
				portalPage, err := repo.FindByID(ctx, result.ID)
				assert.NoError(t, err)
				assert.Equal(t, "John's Updated Page", portalPage.Title)
				assert.Equal(t, "Updated bio", portalPage.Bio)
				assert.Equal(t, "https://example.com/new-image.jpg", portalPage.ProfileImageURL)
				assert.Equal(t, domain.Theme("dark"), portalPage.Theme)
				// slug 不應該改變
				assert.Equal(t, "john-doe", portalPage.Slug)
			},
		},
		{
			name: "成功更新 Slug",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "old-slug",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "new-slug",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *UpdatePortalPageResult) {
				assert.NotNil(t, result)

				// 驗證 slug 已被更新
				portalPage, err := repo.FindBySlug(ctx, "new-slug")
				assert.NoError(t, err)
				assert.NotNil(t, portalPage)
				assert.Equal(t, result.ID, portalPage.ID)

				// 舊的 slug 應該找不到
				oldPage, err := repo.FindBySlug(ctx, "old-slug")
				assert.Error(t, err)
				assert.Nil(t, oldPage)
			},
		},
		{
			name: "成功更新 Links",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				// 添加一個連結
				page.AddLink(domain.LinkParams{
					Title:        "Original Link",
					URL:          "https://example.com/original",
					Description:  "Original description",
					IconURL:      "https://example.com/icon.png",
					DisplayOrder: 1,
				})

				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "Updated Link",
						URL:          "https://example.com/updated",
						Description:  "Updated description",
						IconURL:      "https://example.com/new-icon.png",
						DisplayOrder: 1,
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *UpdatePortalPageResult) {
				assert.NotNil(t, result)

				// 驗證連結已被更新
				portalPage, err := repo.FindByID(ctx, result.ID)
				assert.NoError(t, err)
				assert.Len(t, portalPage.Links, 1)
				assert.Equal(t, "Updated Link", portalPage.Links[0].Title)
				assert.Equal(t, "https://example.com/updated", portalPage.Links[0].URL)
			},
		},
		{
			name: "頁面不存在",
			setupData: func(t *testing.T) int {
				return 0 // 不創建任何頁面
			},
			params: &UpdatePortalPageParams{
				UserID:       1,
				PortalPageID: 99999, // 不存在的 ID
				Title:        "Test",
			},
			wantErr:        true,
			expectedErrMsg: "portal page not found",
		},
		{
			name: "使用者無權限更新（非頁面擁有者）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "user1-page",
					Title:           "User 1's Page",
					Bio:             "Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 2, // 不同的使用者
				Title:  "Trying to update",
			},
			wantErr:        true,
			expectedErrMsg: "unauthorized",
		},
		{
			name: "Slug 太短（少於 3 字元）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "ab",
			},
			wantErr:        true,
			expectedErrMsg: "slug must be between 3 and 50 characters",
		},
		{
			name: "Slug 太長（超過 50 字元）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   strings.Repeat("a", 51),
			},
			wantErr:        true,
			expectedErrMsg: "slug must be between 3 and 50 characters",
		},
		{
			name: "Slug 格式錯誤（包含大寫字母）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "Test-Page",
			},
			wantErr:        true,
			expectedErrMsg: "slug can only contain lowercase letters",
		},
		{
			name: "Slug 格式錯誤（連字號開頭）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "-test-page",
			},
			wantErr:        true,
			expectedErrMsg: "slug can only contain lowercase letters",
		},
		{
			name: "Slug 為保留字",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "admin",
			},
			wantErr:        true,
			expectedErrMsg: "slug is reserved and cannot be used",
		},
		{
			name: "Slug 已被其他頁面使用",
			setupData: func(t *testing.T) int {
				// 創建第一個頁面
				page1 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "existing-slug",
					Title:           "Existing Page",
					Bio:             "Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page1)
				assert.NoError(t, err)

				// 創建第二個頁面（要更新的）
				page2 := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "my-page",
					Title:           "My Page",
					Bio:             "Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err = repo.Create(ctx, page2)
				assert.NoError(t, err)
				return page2.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "existing-slug", // 嘗試使用已存在的 slug
			},
			wantErr:        true,
			expectedErrMsg: "slug already exists",
		},
		{
			name: "更新為相同的 Slug（應該成功）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "my-slug",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Slug:   "my-slug", // 相同的 slug
				Title:  "Updated Title",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *UpdatePortalPageResult) {
				assert.NotNil(t, result)
				portalPage, err := repo.FindBySlug(ctx, "my-slug")
				assert.NoError(t, err)
				assert.Equal(t, "Updated Title", portalPage.Title)
			},
		},
		{
			name: "Title 為空",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Title:  "",
			},
			wantErr: false, // 空字串表示不更新
			checkResult: func(t *testing.T, result *UpdatePortalPageResult) {
				assert.NotNil(t, result)
				// title 應該保持不變
				portalPage, err := repo.FindByID(ctx, result.ID)
				assert.NoError(t, err)
				assert.Equal(t, "Test Page", portalPage.Title)
			},
		},
		{
			name: "Title 太長（超過 100 字元）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Title:  strings.Repeat("a", 101),
			},
			wantErr:        true,
			expectedErrMsg: "title must be between 1 and 100 characters",
		},
		{
			name: "Bio 太長（超過 500 字元）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Bio:    strings.Repeat("a", 501),
			},
			wantErr:        true,
			expectedErrMsg: "bio must not exceed 500 characters",
		},
		{
			name: "無效的 Profile Image URL",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID:          1,
				ProfileImageURL: "invalid-url",
			},
			wantErr:        true,
			expectedErrMsg: "profile_image_url must be a valid URL",
		},
		{
			name: "無效的 Theme",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Theme:  "invalid-theme",
			},
			wantErr:        true,
			expectedErrMsg: "theme must be one of: light, dark",
		},
		{
			name: "Link title 為空",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "",
						URL:          "https://example.com",
						DisplayOrder: 1,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link title must be between 1 and 100 characters",
		},
		{
			name: "Link title 太長（超過 100 字元）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        strings.Repeat("a", 101),
						URL:          "https://example.com",
						DisplayOrder: 1,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link title must be between 1 and 100 characters",
		},
		{
			name: "Link URL 格式錯誤",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "Test Link",
						URL:          "invalid-url",
						DisplayOrder: 1,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link URL must be a valid URL",
		},
		{
			name: "Link description 太長（超過 500 字元）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "Test Link",
						URL:          "https://example.com",
						Description:  strings.Repeat("a", 501),
						DisplayOrder: 1,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link description must not exceed 500 characters",
		},
		{
			name: "Link icon_url 格式錯誤",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "Test Link",
						URL:          "https://example.com",
						IconURL:      "invalid-url",
						DisplayOrder: 1,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link icon_url must be a valid URL",
		},
		{
			name: "Link display_order 為 0（無效）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "Test Link",
						URL:          "https://example.com",
						DisplayOrder: 0,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link display_order must be a positive integer",
		},
		{
			name: "Link display_order 為負數（無效）",
			setupData: func(t *testing.T) int {
				page := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "test-page",
					Title:           "Test Page",
					Bio:             "Test Bio",
					ProfileImageURL: "https://example.com/image.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, page)
				assert.NoError(t, err)
				return page.ID
			},
			params: &UpdatePortalPageParams{
				UserID: 1,
				Links: []Link{
					{
						Title:        "Test Link",
						URL:          "https://example.com",
						DisplayOrder: -1,
					},
				},
			},
			wantErr:        true,
			expectedErrMsg: "link display_order must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 每個測試前清理數據庫
			cleanupTestDB(t, db)

			// 準備測試數據，獲取創建的 Portal Page ID
			var portalPageID int
			if tt.setupData != nil {
				portalPageID = tt.setupData(t)
			}

			// 設置 PortalPageID
			if tt.params != nil {
				if tt.params.PortalPageID == 0 {
					tt.params.PortalPageID = portalPageID
				}
			}

			uc := NewUpdatePortalPageUC(repo)
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
