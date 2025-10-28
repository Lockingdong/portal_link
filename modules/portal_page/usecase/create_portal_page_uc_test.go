package usecase

import (
	"context"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePortalPageUC_Execute(t *testing.T) {
	repo := repository.NewInMemoryPortalPageRepository()
	ctx := context.Background()

	tests := []struct {
		name           string
		params         *CreatePortalPageParams
		setupData      func(t *testing.T) // 準備測試數據
		wantErr        bool
		expectedErrMsg string
		checkResult    func(t *testing.T, result *CreatePortalPageResult)
	}{
		{
			name: "成功建立 Portal Page",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "john-doe",
				Title:           "John's Page",
				Bio:             "Welcome to my page",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *CreatePortalPageResult) {
				assert.NotNil(t, result)
				assert.NotZero(t, result.ID)

				// 驗證 Portal Page 已被創建
				portalPage, err := repo.FindBySlug(ctx, "john-doe")
				assert.NoError(t, err)
				assert.NotNil(t, portalPage)
				assert.Equal(t, "John's Page", portalPage.Title)
				assert.Equal(t, "Welcome to my page", portalPage.Bio)
				assert.Equal(t, "https://example.com/image.jpg", portalPage.ProfileImageURL)
				assert.Equal(t, domain.Theme("light"), portalPage.Theme)
			},
		},
		{
			name: "Slug 太短（少於 3 字元）",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "ab",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "slug must be between 3 and 50 characters",
		},
		{
			name: "Slug 太長（超過 50 字元）",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            strings.Repeat("a", 51),
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "slug must be between 3 and 50 characters",
		},
		{
			name: "Slug 格式錯誤（包含大寫字母）",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "John-Doe",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "slug can only contain lowercase letters",
		},
		{
			name: "Slug 格式錯誤（連字號開頭）",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "-test-page",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "slug can only contain lowercase letters",
		},
		{
			name: "Slug 為保留字",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "admin",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "slug is reserved and cannot be used",
		},
		{
			name: "Slug 已存在",
			params: &CreatePortalPageParams{
				UserID:          2,
				Slug:            "existing-page",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			setupData: func(t *testing.T) {
				// 先創建一個 Portal Page
				existingPage := domain.NewPortalPage(domain.PortalPageParams{
					UserID:          1,
					Slug:            "existing-page",
					Title:           "Existing Page",
					Bio:             "Existing Bio",
					ProfileImageURL: "https://example.com/existing.jpg",
					Theme:           domain.Theme("light"),
				})
				err := repo.Create(ctx, existingPage)
				assert.NoError(t, err)
			},
			wantErr:        true,
			expectedErrMsg: "slug already exists",
		},
		{
			name: "Title 為空",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "test-page",
				Title:           "",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "title must be between 1 and 100 characters",
		},
		{
			name: "Title 太長（超過 100 字元）",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "test-page",
				Title:           strings.Repeat("a", 101),
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "title must be between 1 and 100 characters",
		},
		{
			name: "Bio 太長（超過 500 字元）",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "test-page",
				Title:           "Test Page",
				Bio:             strings.Repeat("a", 501),
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "bio must not exceed 500 characters",
		},
		{
			name: "無效的 Profile Image URL",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "test-page",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "invalid-url",
				Theme:           "light",
			},
			wantErr:        true,
			expectedErrMsg: "profile_image_url must be a valid URL",
		},
		{
			name: "無效的 Theme",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "test-page",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "invalid-theme",
			},
			wantErr:        true,
			expectedErrMsg: "theme must be one of: light, dark",
		},
		{
			name: "Theme 為空時使用預設值",
			params: &CreatePortalPageParams{
				UserID:          1,
				Slug:            "test-page",
				Title:           "Test Page",
				Bio:             "Test Bio",
				ProfileImageURL: "https://example.com/image.jpg",
				Theme:           "",
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *CreatePortalPageResult) {
				assert.NotNil(t, result)
				portalPage, err := repo.FindBySlug(ctx, "test-page")
				assert.NoError(t, err)
				assert.Equal(t, domain.ThemeLight, portalPage.Theme)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 準備測試數據
			if tt.setupData != nil {
				tt.setupData(t)
			}

			uc := NewCreatePortalPageUC(repo)
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
