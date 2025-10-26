package repository

import (
	"context"
	"fmt"
	"portal_link/modules/portal_page/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryPortalPageRepository_Create(t *testing.T) {
	repo := NewInMemoryPortalPageRepository()
	ctx := context.Background()
	now := time.Now()

	testUserID := 1

	tests := []struct {
		name       string
		portalPage *domain.PortalPage
		wantErr    bool
		errMsg     string
		check      func(t *testing.T, portalPage *domain.PortalPage)
	}{
		{
			name: "成功建立沒有 Links 的 PortalPage",
			portalPage: &domain.PortalPage{
				UserID:          testUserID,
				Slug:            "test-user",
				Title:           "測試用戶的頁面",
				Bio:             "這是測試用戶的自我介紹",
				ProfileImageURL: "https://example.com/avatar.jpg",
				Theme:           "light",
				CreatedAt:       now,
				UpdatedAt:       now,
				Links:           []*domain.Link{},
			},
			wantErr: false,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.NotZero(t, portalPage.ID)
				assert.Equal(t, testUserID, portalPage.UserID)
				assert.Equal(t, "test-user", portalPage.Slug)
				assert.Equal(t, "測試用戶的頁面", portalPage.Title)
				assert.Equal(t, "這是測試用戶的自我介紹", portalPage.Bio)
				assert.Equal(t, "https://example.com/avatar.jpg", portalPage.ProfileImageURL)
				assert.Equal(t, "light", string(portalPage.Theme))
				assert.NotZero(t, portalPage.CreatedAt)
				assert.NotZero(t, portalPage.UpdatedAt)
				assert.Empty(t, portalPage.Links)
			},
		},
		{
			name: "成功建立有 Links 的 PortalPage",
			portalPage: &domain.PortalPage{
				UserID:          testUserID,
				Slug:            "user-with-links",
				Title:           "有連結的用戶頁面",
				Bio:             "我有很多連結",
				ProfileImageURL: "https://example.com/profile.jpg",
				Theme:           "dark",
				CreatedAt:       now,
				UpdatedAt:       now,
				Links: []*domain.Link{
					{
						Title:        "GitHub",
						URL:          "https://github.com/testuser",
						Description:  "我的 GitHub 頁面",
						IconURL:      "https://github.com/favicon.ico",
						DisplayOrder: 1,
						CreatedAt:    now,
						UpdatedAt:    now,
					},
					{
						Title:        "Twitter",
						URL:          "https://twitter.com/testuser",
						Description:  "我的 Twitter",
						IconURL:      "https://twitter.com/favicon.ico",
						DisplayOrder: 2,
						CreatedAt:    now,
						UpdatedAt:    now,
					},
				},
			},
			wantErr: false,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.NotZero(t, portalPage.ID)
				assert.Equal(t, "user-with-links", portalPage.Slug)
				assert.Len(t, portalPage.Links, 2)

				// 檢查第一個 Link
				link1 := portalPage.Links[0]
				assert.NotZero(t, link1.ID)
				assert.Equal(t, portalPage.ID, link1.PortalPageID)
				assert.Equal(t, "GitHub", link1.Title)
				assert.Equal(t, "https://github.com/testuser", link1.URL)
				assert.Equal(t, "我的 GitHub 頁面", link1.Description)
				assert.Equal(t, "https://github.com/favicon.ico", link1.IconURL)
				assert.Equal(t, 1, link1.DisplayOrder)
				assert.NotZero(t, link1.CreatedAt)
				assert.NotZero(t, link1.UpdatedAt)

				// 檢查第二個 Link
				link2 := portalPage.Links[1]
				assert.NotZero(t, link2.ID)
				assert.Equal(t, portalPage.ID, link2.PortalPageID)
				assert.Equal(t, "Twitter", link2.Title)
				assert.Equal(t, 2, link2.DisplayOrder)
			},
		},
		{
			name: "成功建立只有必要欄位的 PortalPage",
			portalPage: &domain.PortalPage{
				UserID:    testUserID,
				Slug:      "minimal-user",
				Title:     "最小化頁面",
				CreatedAt: now,
				UpdatedAt: now,
				Links:     []*domain.Link{},
			},
			wantErr: false,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.NotZero(t, portalPage.ID)
				assert.Equal(t, "minimal-user", portalPage.Slug)
				assert.Equal(t, "最小化頁面", portalPage.Title)
				assert.Empty(t, portalPage.Bio)
				assert.Empty(t, portalPage.ProfileImageURL)
				assert.Empty(t, portalPage.Theme)
			},
		},
		{
			name: "失敗 - Slug 重複",
			portalPage: &domain.PortalPage{
				UserID:    testUserID,
				Slug:      "test-user", // 與第一個測試重複
				Title:     "重複的 Slug",
				CreatedAt: now,
				UpdatedAt: now,
				Links:     []*domain.Link{},
			},
			wantErr: true,
			errMsg:  "duplicate key",
		},
		{
			name: "成功建立沒有時間戳記的 PortalPage - 應自動設定",
			portalPage: &domain.PortalPage{
				UserID: testUserID,
				Slug:   "auto-timestamp",
				Title:  "自動時間戳記頁面",
				Links:  []*domain.Link{},
			},
			wantErr: false,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.NotZero(t, portalPage.ID)
				assert.Equal(t, "auto-timestamp", portalPage.Slug)
				assert.NotZero(t, portalPage.CreatedAt)
				assert.NotZero(t, portalPage.UpdatedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.portalPage)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			assert.NoError(t, err)

			if tt.check != nil {
				tt.check(t, tt.portalPage)
			}
		})
	}
}

func TestInMemoryPortalPageRepository_Update(t *testing.T) {
	repo := NewInMemoryPortalPageRepository()
	ctx := context.Background()
	now := time.Now()
	testUserID := 1

	tests := []struct {
		name    string
		setup   func(t *testing.T) *domain.PortalPage
		update  func(t *testing.T, original *domain.PortalPage) *domain.PortalPage
		wantErr bool
		errMsg  string
		check   func(t *testing.T, original, updated *domain.PortalPage)
	}{
		{
			name: "成功更新 PortalPage 基本信息（不改變 Links）",
			setup: func(t *testing.T) *domain.PortalPage {
				portalPage := &domain.PortalPage{
					UserID:          testUserID,
					Slug:            "update-test-1",
					Title:           "原始標題",
					Bio:             "原始介紹",
					ProfileImageURL: "https://example.com/old.jpg",
					Theme:           "light",
					CreatedAt:       now,
					UpdatedAt:       now,
					Links: []*domain.Link{
						{
							Title:        "Link 1",
							URL:          "https://example.com/1",
							DisplayOrder: 1,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
					},
				}
				err := repo.Create(ctx, portalPage)
				assert.NoError(t, err)
				return portalPage
			},
			update: func(t *testing.T, original *domain.PortalPage) *domain.PortalPage {
				updated := &domain.PortalPage{
					ID:              original.ID,
					UserID:          original.UserID,
					Slug:            "update-test-1-modified",
					Title:           "更新後的標題",
					Bio:             "更新後的介紹",
					ProfileImageURL: "https://example.com/new.jpg",
					Theme:           "dark",
					UpdatedAt:       time.Now(),
					Links:           original.Links, // 保持 Links 不變
				}
				return updated
			},
			wantErr: false,
			check: func(t *testing.T, original, updated *domain.PortalPage) {
				assert.Equal(t, original.ID, updated.ID)
				assert.Equal(t, "update-test-1-modified", updated.Slug)
				assert.Equal(t, "更新後的標題", updated.Title)
				assert.Equal(t, "更新後的介紹", updated.Bio)
				assert.Equal(t, "https://example.com/new.jpg", updated.ProfileImageURL)
				assert.Equal(t, "dark", string(updated.Theme))
				assert.Len(t, updated.Links, 1)
				assert.Equal(t, "Link 1", updated.Links[0].Title)
			},
		},
		{
			name: "成功更新並新增 Links",
			setup: func(t *testing.T) *domain.PortalPage {
				portalPage := &domain.PortalPage{
					UserID:    testUserID,
					Slug:      "update-test-2",
					Title:     "測試頁面",
					CreatedAt: now,
					UpdatedAt: now,
					Links: []*domain.Link{
						{
							Title:        "Original Link",
							URL:          "https://example.com/original",
							DisplayOrder: 1,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
					},
				}
				err := repo.Create(ctx, portalPage)
				assert.NoError(t, err)
				return portalPage
			},
			update: func(t *testing.T, original *domain.PortalPage) *domain.PortalPage {
				updated := &domain.PortalPage{
					ID:        original.ID,
					UserID:    original.UserID,
					Slug:      original.Slug,
					Title:     original.Title,
					UpdatedAt: time.Now(),
					Links: []*domain.Link{
						original.Links[0], // 保留原始 Link
						{
							// 新增 Link (ID = 0)
							Title:        "New Link",
							URL:          "https://example.com/new",
							Description:  "新增的連結",
							DisplayOrder: 2,
							CreatedAt:    time.Now(),
							UpdatedAt:    time.Now(),
						},
					},
				}
				return updated
			},
			wantErr: false,
			check: func(t *testing.T, original, updated *domain.PortalPage) {
				assert.Len(t, updated.Links, 2)
				assert.Equal(t, "Original Link", updated.Links[0].Title)
				assert.Equal(t, "New Link", updated.Links[1].Title)
				assert.NotZero(t, updated.Links[1].ID) // 新增的 Link 應該有 ID
			},
		},
		{
			name: "成功更新並修改現有 Links",
			setup: func(t *testing.T) *domain.PortalPage {
				portalPage := &domain.PortalPage{
					UserID:    testUserID,
					Slug:      "update-test-3",
					Title:     "測試頁面",
					CreatedAt: now,
					UpdatedAt: now,
					Links: []*domain.Link{
						{
							Title:        "GitHub",
							URL:          "https://github.com/old",
							Description:  "舊的描述",
							DisplayOrder: 1,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
						{
							Title:        "Twitter",
							URL:          "https://twitter.com/old",
							DisplayOrder: 2,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
					},
				}
				err := repo.Create(ctx, portalPage)
				assert.NoError(t, err)
				return portalPage
			},
			update: func(t *testing.T, original *domain.PortalPage) *domain.PortalPage {
				updated := &domain.PortalPage{
					ID:        original.ID,
					UserID:    original.UserID,
					Slug:      original.Slug,
					Title:     original.Title,
					UpdatedAt: time.Now(),
					Links: []*domain.Link{
						{
							ID:           original.Links[0].ID, // 使用現有 Link 的 ID
							PortalPageID: original.ID,
							Title:        "GitHub Updated",
							URL:          "https://github.com/new",
							Description:  "更新後的描述",
							DisplayOrder: 1,
							UpdatedAt:    time.Now(),
						},
						{
							ID:           original.Links[1].ID,
							PortalPageID: original.ID,
							Title:        "Twitter Updated",
							URL:          "https://twitter.com/new",
							Description:  "新增的描述",
							DisplayOrder: 2,
							UpdatedAt:    time.Now(),
						},
					},
				}
				return updated
			},
			wantErr: false,
			check: func(t *testing.T, original, updated *domain.PortalPage) {
				assert.Len(t, updated.Links, 2)
				assert.Equal(t, "GitHub Updated", updated.Links[0].Title)
				assert.Equal(t, "https://github.com/new", updated.Links[0].URL)
				assert.Equal(t, "更新後的描述", updated.Links[0].Description)
				assert.Equal(t, "Twitter Updated", updated.Links[1].Title)
				assert.Equal(t, "新增的描述", updated.Links[1].Description)
			},
		},
		{
			name: "成功更新並刪除 Links",
			setup: func(t *testing.T) *domain.PortalPage {
				portalPage := &domain.PortalPage{
					UserID:    testUserID,
					Slug:      "update-test-4",
					Title:     "測試頁面",
					CreatedAt: now,
					UpdatedAt: now,
					Links: []*domain.Link{
						{
							Title:        "Link 1",
							URL:          "https://example.com/1",
							DisplayOrder: 1,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
						{
							Title:        "Link 2",
							URL:          "https://example.com/2",
							DisplayOrder: 2,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
						{
							Title:        "Link 3",
							URL:          "https://example.com/3",
							DisplayOrder: 3,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
					},
				}
				err := repo.Create(ctx, portalPage)
				assert.NoError(t, err)
				return portalPage
			},
			update: func(t *testing.T, original *domain.PortalPage) *domain.PortalPage {
				updated := &domain.PortalPage{
					ID:        original.ID,
					UserID:    original.UserID,
					Slug:      original.Slug,
					Title:     original.Title,
					UpdatedAt: time.Now(),
					Links: []*domain.Link{
						original.Links[0], // 只保留第一個 Link
					},
				}
				return updated
			},
			wantErr: false,
			check: func(t *testing.T, original, updated *domain.PortalPage) {
				assert.Len(t, updated.Links, 1)
				assert.Equal(t, "Link 1", updated.Links[0].Title)
			},
		},
		{
			name: "失敗 - PortalPage 不存在",
			setup: func(t *testing.T) *domain.PortalPage {
				// 不建立任何 Portal Page
				return nil
			},
			update: func(t *testing.T, original *domain.PortalPage) *domain.PortalPage {
				return &domain.PortalPage{
					ID:        99999, // 不存在的 ID
					UserID:    testUserID,
					Slug:      "non-existent",
					Title:     "不存在的頁面",
					UpdatedAt: time.Now(),
					Links:     []*domain.Link{},
				}
			},
			wantErr: true,
			errMsg:  "sql: no rows in result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: 建立初始資料
			original := tt.setup(t)

			// Update: 執行更新
			updated := tt.update(t, original)
			err := repo.Update(ctx, updated)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)

			if tt.check != nil {
				tt.check(t, original, updated)
			}
		})
	}
}

func TestInMemoryPortalPageRepository_FindBySlug(t *testing.T) {
	repo := NewInMemoryPortalPageRepository()
	ctx := context.Background()
	now := time.Now()
	testUserID := 1

	// 準備測試數據
	testPortalPage := &domain.PortalPage{
		UserID:          testUserID,
		Slug:            "find-by-slug-test",
		Title:           "測試查找頁面",
		Bio:             "這是用於測試的頁面",
		ProfileImageURL: "https://example.com/avatar.jpg",
		Theme:           "light",
		CreatedAt:       now,
		UpdatedAt:       now,
		Links: []*domain.Link{
			{
				Title:        "GitHub",
				URL:          "https://github.com/testuser",
				Description:  "我的 GitHub",
				IconURL:      "https://github.com/favicon.ico",
				DisplayOrder: 2, // 故意設定為 2
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				Title:        "Twitter",
				URL:          "https://twitter.com/testuser",
				Description:  "我的 Twitter",
				IconURL:      "https://twitter.com/favicon.ico",
				DisplayOrder: 1, // 故意設定為 1，應該排在前面
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}
	err := repo.Create(ctx, testPortalPage)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		slug    string
		wantErr bool
		check   func(t *testing.T, portalPage *domain.PortalPage)
	}{
		{
			name:    "成功根據 Slug 查找 PortalPage",
			slug:    "find-by-slug-test",
			wantErr: false,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.NotNil(t, portalPage)
				assert.Equal(t, testPortalPage.ID, portalPage.ID)
				assert.Equal(t, "find-by-slug-test", portalPage.Slug)
				assert.Equal(t, "測試查找頁面", portalPage.Title)
				assert.Equal(t, "這是用於測試的頁面", portalPage.Bio)
				assert.Equal(t, "https://example.com/avatar.jpg", portalPage.ProfileImageURL)
				assert.Equal(t, "light", string(portalPage.Theme))
				assert.Equal(t, testUserID, portalPage.UserID)
				assert.NotZero(t, portalPage.CreatedAt)
				assert.NotZero(t, portalPage.UpdatedAt)

				// 驗證 Links 按照 DisplayOrder 排序
				assert.Len(t, portalPage.Links, 2)
				assert.Equal(t, "Twitter", portalPage.Links[0].Title) // DisplayOrder = 1
				assert.Equal(t, 1, portalPage.Links[0].DisplayOrder)
				assert.Equal(t, "GitHub", portalPage.Links[1].Title) // DisplayOrder = 2
				assert.Equal(t, 2, portalPage.Links[1].DisplayOrder)
			},
		},
		{
			name:    "失敗 - Slug 不存在",
			slug:    "non-existent-slug",
			wantErr: true,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.Nil(t, portalPage)
			},
		},
		{
			name:    "失敗 - 空白 Slug",
			slug:    "",
			wantErr: true,
			check: func(t *testing.T, portalPage *domain.PortalPage) {
				assert.Nil(t, portalPage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portalPage, err := repo.FindBySlug(ctx, tt.slug)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.check != nil {
				tt.check(t, portalPage)
			}
		})
	}
}

func TestInMemoryPortalPageRepository_ListByUserID(t *testing.T) {
	repo := NewInMemoryPortalPageRepository()
	ctx := context.Background()
	now := time.Now()
	testUserID := 1
	anotherUserID := 2

	// 準備測試數據 - 為 testUserID 建立多個 Portal Pages
	// 故意調整創建時間順序來測試排序
	portalPage1 := &domain.PortalPage{
		UserID:          testUserID,
		Slug:            "user-page-1",
		Title:           "第一個頁面",
		Bio:             "第一個測試頁面",
		ProfileImageURL: "https://example.com/avatar1.jpg",
		Theme:           "light",
		CreatedAt:       now.Add(-2 * time.Hour), // 較早的時間
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, portalPage1)
	assert.NoError(t, err)

	portalPage2 := &domain.PortalPage{
		UserID:    testUserID,
		Slug:      "user-page-2",
		Title:     "第二個頁面",
		Bio:       "第二個測試頁面",
		Theme:     "dark",
		CreatedAt: now.Add(-1 * time.Hour), // 較晚的時間
		UpdatedAt: now,
	}
	err = repo.Create(ctx, portalPage2)
	assert.NoError(t, err)

	// 為另一個用戶建立一個 Portal Page
	portalPage3 := &domain.PortalPage{
		UserID:    anotherUserID,
		Slug:      "another-user-page",
		Title:     "另一個用戶的頁面",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = repo.Create(ctx, portalPage3)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		userID  int
		wantErr bool
		check   func(t *testing.T, portalPages []*domain.PortalPage)
	}{
		{
			name:    "成功查找擁有多個 PortalPages 的用戶",
			userID:  testUserID,
			wantErr: false,
			check: func(t *testing.T, portalPages []*domain.PortalPage) {
				assert.NotNil(t, portalPages)
				assert.Len(t, portalPages, 2)

				// 驗證按 CreatedAt 升序排序
				assert.Equal(t, "user-page-1", portalPages[0].Slug) // 較早創建的
				assert.Equal(t, "user-page-2", portalPages[1].Slug) // 較晚創建的
				
				// 驗證第一個 Portal Page
				page1 := portalPages[0]
				assert.Equal(t, testUserID, page1.UserID)
				assert.Equal(t, "第一個頁面", page1.Title)
				assert.Equal(t, "第一個測試頁面", page1.Bio)
				assert.Equal(t, "light", string(page1.Theme))
				assert.Nil(t, page1.Links) // ListByUserID 不包含 Links

				// 驗證第二個 Portal Page
				page2 := portalPages[1]
				assert.Equal(t, testUserID, page2.UserID)
				assert.Equal(t, "第二個頁面", page2.Title)
				assert.Equal(t, "dark", string(page2.Theme))
				assert.Nil(t, page2.Links) // ListByUserID 不包含 Links
			},
		},
		{
			name:    "成功查找擁有單個 PortalPage 的用戶",
			userID:  anotherUserID,
			wantErr: false,
			check: func(t *testing.T, portalPages []*domain.PortalPage) {
				assert.NotNil(t, portalPages)
				assert.Len(t, portalPages, 1)
				assert.Equal(t, anotherUserID, portalPages[0].UserID)
				assert.Equal(t, "another-user-page", portalPages[0].Slug)
				assert.Equal(t, "另一個用戶的頁面", portalPages[0].Title)
				assert.Nil(t, portalPages[0].Links)
			},
		},
		{
			name:    "查找沒有 PortalPage 的用戶 - 返回空列表",
			userID:  99999,
			wantErr: false,
			check: func(t *testing.T, portalPages []*domain.PortalPage) {
				assert.NotNil(t, portalPages)
				assert.Empty(t, portalPages)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portalPages, err := repo.ListByUserID(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.check != nil {
				tt.check(t, portalPages)
			}
		})
	}
}

func TestInMemoryPortalPageRepository_FindByID(t *testing.T) {
	repo := NewInMemoryPortalPageRepository()
	ctx := context.Background()
	now := time.Now()
	testUserID := 1

	portalPage := &domain.PortalPage{
		UserID:          testUserID,
		Slug:            "find-by-id-test",
		Title:           "測試查找頁面",
		Bio:             "這是用於測試的頁面",
		ProfileImageURL: "https://example.com/avatar.jpg",
		Theme:           "light",
		CreatedAt:       now,
		UpdatedAt:       now,
		Links: []*domain.Link{
			{
				Title:        "Link 2",
				URL:          "https://example.com/2",
				DisplayOrder: 2,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				Title:        "Link 1",
				URL:          "https://example.com/1",
				DisplayOrder: 1,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}
	err := repo.Create(ctx, portalPage)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      int
		wantErr bool
		check   func(t *testing.T, found *domain.PortalPage)
	}{
		{
			name:    "成功根據 ID 查找 PortalPage",
			id:      portalPage.ID,
			wantErr: false,
			check: func(t *testing.T, found *domain.PortalPage) {
				assert.NotNil(t, found)
				assert.Equal(t, portalPage.ID, found.ID)
				assert.Equal(t, "find-by-id-test", found.Slug)
				assert.Equal(t, "測試查找頁面", found.Title)
				assert.Equal(t, "這是用於測試的頁面", found.Bio)
				assert.Equal(t, "https://example.com/avatar.jpg", found.ProfileImageURL)
				assert.Equal(t, "light", string(found.Theme))
				assert.Equal(t, testUserID, found.UserID)

				// 驗證 Links 按照 DisplayOrder 排序
				assert.Len(t, found.Links, 2)
				assert.Equal(t, "Link 1", found.Links[0].Title) // DisplayOrder = 1
				assert.Equal(t, 1, found.Links[0].DisplayOrder)
				assert.Equal(t, "Link 2", found.Links[1].Title) // DisplayOrder = 2
				assert.Equal(t, 2, found.Links[1].DisplayOrder)
			},
		},
		{
			name:    "失敗 - ID 不存在",
			id:      99999,
			wantErr: true,
			check: func(t *testing.T, found *domain.PortalPage) {
				assert.Nil(t, found)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.check != nil {
				tt.check(t, found)
			}
		})
	}
}

func TestInMemoryPortalPageRepository_Concurrency(t *testing.T) {
	repo := NewInMemoryPortalPageRepository()
	ctx := context.Background()

	// 測試並發安全性
	const numGoroutines = 10
	const itemsPerGoroutine = 5

	done := make(chan bool, numGoroutines)

	// 並發創建 Portal Pages
	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer func() { done <- true }()

			for j := 0; j < itemsPerGoroutine; j++ {
				portalPage := &domain.PortalPage{
					UserID: routineID,
					Slug:   fmt.Sprintf("test-routine-%d-item-%d", routineID, j),
					Title:  fmt.Sprintf("Test Page %d-%d", routineID, j),
					Links:  []*domain.Link{},
				}
				err := repo.Create(ctx, portalPage)
				assert.NoError(t, err)
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// 驗證所有項目都已創建
	for i := 0; i < numGoroutines; i++ {
		pages, err := repo.ListByUserID(ctx, i)
		assert.NoError(t, err)
		assert.Len(t, pages, itemsPerGoroutine)
	}
}