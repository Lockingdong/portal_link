package repository

import (
	"context"
	"database/sql"
	"portal_link/modules/portal_page/domain"
	"portal_link/pkg"
	"portal_link/pkg/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// setupTestDB 設置測試數據庫連接
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	config.Init()

	// 使用 config 獲取資料庫連接
	db := pkg.NewPG(config.GetDBConfig().DSN())

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	return db
}

// cleanupTestDB 清理測試數據
func cleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	// 先刪除 links (因為有外鍵約束)
	_, err := db.Exec("DELETE FROM portal_link.links")
	if err != nil {
		t.Logf("failed to cleanup links: %v", err)
	}
	// 再刪除 portal_pages
	_, err = db.Exec("DELETE FROM portal_link.portal_pages")
	if err != nil {
		t.Logf("failed to cleanup portal_pages: %v", err)
	}
	// 清理 users
	_, err = db.Exec("DELETE FROM portal_link.users")
	if err != nil {
		t.Logf("failed to cleanup users: %v", err)
	}
}

func TestPortalPageRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewPortalPageRepository(db)
	ctx := context.Background()
	now := time.Now()

	// 準備測試用的 User ID
	// 假設已經有 user_id = 1 的使用者
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
				assert.Equal(t, "light", portalPage.Theme)
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
			errMsg:  "violates foreign key constraint",
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

			if err != nil {
				assert.Error(t, err)
				return
			}

			if tt.check != nil {
				tt.check(t, tt.portalPage)
			}
		})
	}
}

func TestPortalPageRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewPortalPageRepository(db)
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
				assert.Equal(t, "dark", updated.Theme)
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
			name: "成功更新並同時新增、更新、刪除 Links（綜合場景）",
			setup: func(t *testing.T) *domain.PortalPage {
				portalPage := &domain.PortalPage{
					UserID:    testUserID,
					Slug:      "update-test-5",
					Title:     "測試頁面",
					CreatedAt: now,
					UpdatedAt: now,
					Links: []*domain.Link{
						{
							Title:        "Keep and Update",
							URL:          "https://example.com/keep",
							DisplayOrder: 1,
							CreatedAt:    now,
							UpdatedAt:    now,
						},
						{
							Title:        "Will be Deleted",
							URL:          "https://example.com/delete",
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
					Title:     "更新後的標題",
					UpdatedAt: time.Now(),
					Links: []*domain.Link{
						{
							// 更新現有 Link
							ID:           original.Links[0].ID,
							PortalPageID: original.ID,
							Title:        "Updated Link",
							URL:          "https://example.com/updated",
							DisplayOrder: 1,
							UpdatedAt:    time.Now(),
						},
						{
							// 新增 Link
							Title:        "New Link",
							URL:          "https://example.com/new",
							DisplayOrder: 2,
							CreatedAt:    time.Now(),
							UpdatedAt:    time.Now(),
						},
						// 不包含 Links[1]，所以會被刪除
					},
				}
				return updated
			},
			wantErr: false,
			check: func(t *testing.T, original, updated *domain.PortalPage) {
				assert.Equal(t, "更新後的標題", updated.Title)
				assert.Len(t, updated.Links, 2)
				assert.Equal(t, "Updated Link", updated.Links[0].Title)
				assert.Equal(t, original.Links[0].ID, updated.Links[0].ID) // 應該保留原始 ID
				assert.Equal(t, "New Link", updated.Links[1].Title)
				assert.NotZero(t, updated.Links[1].ID)                        // 新增的 Link 應該有新的 ID
				assert.NotEqual(t, original.Links[1].ID, updated.Links[1].ID) // 不應該等於被刪除的 Link ID
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

func TestPortalPageRepository_FindBySlug(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewPortalPageRepository(db)
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
				assert.Equal(t, "light", portalPage.Theme)
				assert.Equal(t, testUserID, portalPage.UserID)
				assert.NotZero(t, portalPage.CreatedAt)
				assert.NotZero(t, portalPage.UpdatedAt)

				// 驗證 Links
				assert.Len(t, portalPage.Links, 2)
				assert.Equal(t, "GitHub", portalPage.Links[0].Title)
				assert.Equal(t, "https://github.com/testuser", portalPage.Links[0].URL)
				assert.Equal(t, "我的 GitHub", portalPage.Links[0].Description)
				assert.Equal(t, 1, portalPage.Links[0].DisplayOrder)
				assert.Equal(t, "Twitter", portalPage.Links[1].Title)
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

func TestPortalPageRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		cleanupTestDB(t, db)
		db.Close()
	})

	repo := NewPortalPageRepository(db)
	ctx := context.Background()
	now := time.Now()
	testUserID := 1
	anotherUserID := 2

	// 準備測試數據 - 為 testUserID 建立多個 Portal Pages
	portalPage1 := &domain.PortalPage{
		UserID:          testUserID,
		Slug:            "user-page-1",
		Title:           "第一個頁面",
		Bio:             "第一個測試頁面",
		ProfileImageURL: "https://example.com/avatar1.jpg",
		Theme:           "light",
		CreatedAt:       now,
		UpdatedAt:       now,
		Links: []*domain.Link{
			{
				Title:        "GitHub",
				URL:          "https://github.com/user1",
				DisplayOrder: 1,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}
	err := repo.Create(ctx, portalPage1)
	assert.NoError(t, err)

	portalPage2 := &domain.PortalPage{
		UserID:    testUserID,
		Slug:      "user-page-2",
		Title:     "第二個頁面",
		Bio:       "第二個測試頁面",
		Theme:     "dark",
		CreatedAt: now,
		UpdatedAt: now,
		Links: []*domain.Link{
			{
				Title:        "Twitter",
				URL:          "https://twitter.com/user1",
				DisplayOrder: 1,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				Title:        "LinkedIn",
				URL:          "https://linkedin.com/in/user1",
				DisplayOrder: 2,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
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
		Links:     []*domain.Link{},
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

				// 驗證第一個 Portal Page
				var page1, page2 *domain.PortalPage
				for _, p := range portalPages {
					if p.Slug == "user-page-1" {
						page1 = p
					} else if p.Slug == "user-page-2" {
						page2 = p
					}
				}

				assert.NotNil(t, page1)
				assert.Equal(t, testUserID, page1.UserID)
				assert.Equal(t, "第一個頁面", page1.Title)
				assert.Equal(t, "第一個測試頁面", page1.Bio)
				assert.Equal(t, "light", page1.Theme)
				assert.Len(t, page1.Links, 1)
				assert.Equal(t, "GitHub", page1.Links[0].Title)

				assert.NotNil(t, page2)
				assert.Equal(t, testUserID, page2.UserID)
				assert.Equal(t, "第二個頁面", page2.Title)
				assert.Equal(t, "dark", page2.Theme)
				assert.Len(t, page2.Links, 2)
				assert.Equal(t, "Twitter", page2.Links[0].Title)
				assert.Equal(t, "LinkedIn", page2.Links[1].Title)
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
				assert.Empty(t, portalPages[0].Links)
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
			portalPages, err := repo.FindByUserID(ctx, tt.userID)

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
