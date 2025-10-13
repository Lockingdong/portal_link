package repository

import (
	"context"
	"database/sql"
	"portal_link/models"
	"portal_link/modules/portal_page/domain"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

var _ domain.PortalPageRepository = (*PortalPageRepository)(nil)

type PortalPageRepository struct {
	db *sql.DB
}

// NewPortalPageRepository 建立新的 Portal Page Repository
func NewPortalPageRepository(db *sql.DB) *PortalPageRepository {
	return &PortalPageRepository{db: db}
}

// Create 建立 Portal Page
func (r *PortalPageRepository) Create(ctx context.Context, portalPage *domain.PortalPage) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 建立 Portal Page
	m := &models.PortalPage{
		UserID:          portalPage.UserID,
		Slug:            portalPage.Slug,
		Title:           portalPage.Title,
		Bio:             null.StringFrom(portalPage.Bio),
		ProfileImageURL: null.StringFrom(portalPage.ProfileImageURL),
		Theme:           null.StringFrom(portalPage.Theme),
		CreatedAt:       null.TimeFrom(portalPage.CreatedAt),
		UpdatedAt:       null.TimeFrom(portalPage.UpdatedAt),
	}

	err = m.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}

	// 回填生成的 ID 和時間戳記
	portalPage.ID = m.ID
	portalPage.CreatedAt = m.CreatedAt.Time
	portalPage.UpdatedAt = m.UpdatedAt.Time

	// 建立關聯的 Links
	for _, link := range portalPage.Links {
		linkModel := &models.Link{
			PortalPageID: m.ID,
			Title:        link.Title,
			URL:          link.URL,
			Description:  null.StringFrom(link.Description),
			IconURL:      null.StringFrom(link.IconURL),
			DisplayOrder: link.DisplayOrder,
			CreatedAt:    null.TimeFrom(link.CreatedAt),
			UpdatedAt:    null.TimeFrom(link.UpdatedAt),
		}

		err = linkModel.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		// 回填生成的 ID 和時間戳記
		link.ID = linkModel.ID
		link.PortalPageID = linkModel.PortalPageID
		link.CreatedAt = linkModel.CreatedAt.Time
		link.UpdatedAt = linkModel.UpdatedAt.Time
	}

	return tx.Commit()
}

// Update 更新 Portal Page
func (r *PortalPageRepository) Update(ctx context.Context, portalPage *domain.PortalPage) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 更新 Portal Page
	m, err := models.FindPortalPage(ctx, tx, portalPage.ID)
	if err != nil {
		return err
	}

	m.UserID = portalPage.UserID
	m.Slug = portalPage.Slug
	m.Title = portalPage.Title
	m.Bio = null.StringFrom(portalPage.Bio)
	m.ProfileImageURL = null.StringFrom(portalPage.ProfileImageURL)
	m.Theme = null.StringFrom(portalPage.Theme)
	m.UpdatedAt = null.TimeFrom(portalPage.UpdatedAt)

	_, err = m.Update(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}

	// 載入現有的 Links
	existingLinks, err := m.Links().All(ctx, tx)
	if err != nil {
		return err
	}

	// 建立現有 Links 的 ID 對照表
	existingLinkMap := make(map[int]*models.Link)
	for _, link := range existingLinks {
		existingLinkMap[link.ID] = link
	}

	// 建立新 Links 的 ID 集合
	newLinkIDs := make(map[int]bool)

	// 處理新的 Links：新增或更新
	for _, link := range portalPage.Links {
		newLinkIDs[link.ID] = true

		if link.ID > 0 {
			// 更新現有的 Link
			if existingLink, exists := existingLinkMap[link.ID]; exists {
				existingLink.Title = link.Title
				existingLink.URL = link.URL
				existingLink.Description = null.StringFrom(link.Description)
				existingLink.IconURL = null.StringFrom(link.IconURL)
				existingLink.DisplayOrder = link.DisplayOrder
				existingLink.UpdatedAt = null.TimeFrom(link.UpdatedAt)

				_, err = existingLink.Update(ctx, tx, boil.Infer())
				if err != nil {
					return err
				}

				// 回填時間戳記
				link.UpdatedAt = existingLink.UpdatedAt.Time
			}
		} else {
			// 新增新的 Link
			linkModel := &models.Link{
				PortalPageID: m.ID,
				Title:        link.Title,
				URL:          link.URL,
				Description:  null.StringFrom(link.Description),
				IconURL:      null.StringFrom(link.IconURL),
				DisplayOrder: link.DisplayOrder,
				CreatedAt:    null.TimeFrom(link.CreatedAt),
				UpdatedAt:    null.TimeFrom(link.UpdatedAt),
			}

			err = linkModel.Insert(ctx, tx, boil.Infer())
			if err != nil {
				return err
			}

			// 回填生成的 ID 和時間戳記
			link.ID = linkModel.ID
			link.PortalPageID = linkModel.PortalPageID
			link.CreatedAt = linkModel.CreatedAt.Time
			link.UpdatedAt = linkModel.UpdatedAt.Time
		}
	}

	// 刪除不存在於新 Links 中的舊 Links
	for _, existingLink := range existingLinks {
		if !newLinkIDs[existingLink.ID] {
			_, err = existingLink.Delete(ctx, tx)
			if err != nil {
				return err
			}
		}
	}

	// 更新 Portal Page 的時間戳記
	portalPage.UpdatedAt = m.UpdatedAt.Time

	return tx.Commit()
}

// FindBySlug 根據 Slug 查找 Portal Page
func (r *PortalPageRepository) FindBySlug(ctx context.Context, slug string) (*domain.PortalPage, error) {
	m, err := models.PortalPages(
		models.PortalPageWhere.Slug.EQ(slug),
		qm.Load(models.PortalPageRels.Links),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return r.toDomainPortalPage(m), nil
}

// FindByUserID 根據 UserID 查找 Portal Page
func (r *PortalPageRepository) FindByUserID(ctx context.Context, userID int) ([]*domain.PortalPage, error) {
	ms, err := models.PortalPages(
		models.PortalPageWhere.UserID.EQ(userID),
		qm.Load(models.PortalPageRels.Links),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	portalPages := make([]*domain.PortalPage, len(ms))
	for i, m := range ms {
		portalPages[i] = r.toDomainPortalPage(m)
	}

	return portalPages, nil
}

// toDomainPortalPage 將 models.PortalPage 轉換為 domain.PortalPage
func (r *PortalPageRepository) toDomainPortalPage(m *models.PortalPage) *domain.PortalPage {
	portalPage := &domain.PortalPage{
		ID:              m.ID,
		UserID:          m.UserID,
		Slug:            m.Slug,
		Title:           m.Title,
		Bio:             m.Bio.String,
		ProfileImageURL: m.ProfileImageURL.String,
		Theme:           m.Theme.String,
		CreatedAt:       m.CreatedAt.Time,
		UpdatedAt:       m.UpdatedAt.Time,
		Links:           make([]*domain.Link, 0),
	}

	// 轉換關聯的 Links
	if m.R != nil && m.R.Links != nil {
		for _, linkModel := range m.R.Links {
			link := &domain.Link{
				ID:           linkModel.ID,
				PortalPageID: linkModel.PortalPageID,
				Title:        linkModel.Title,
				URL:          linkModel.URL,
				Description:  linkModel.Description.String,
				IconURL:      linkModel.IconURL.String,
				DisplayOrder: linkModel.DisplayOrder,
				CreatedAt:    linkModel.CreatedAt.Time,
				UpdatedAt:    linkModel.UpdatedAt.Time,
			}
			portalPage.Links = append(portalPage.Links, link)
		}
	}

	return portalPage
}
