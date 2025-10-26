package repository

import (
	"context"
	"database/sql"
	"fmt"
	"portal_link/modules/portal_page/domain"
	"sort"
	"sync"
	"time"
)

// Ensure InMemoryPortalPageRepository implements domain.PortalPageRepository
var _ domain.PortalPageRepository = (*InMemoryPortalPageRepository)(nil)

// InMemoryPortalPageRepository provides an in-memory implementation of PortalPageRepository
// This is useful for local testing without requiring a database connection
type InMemoryPortalPageRepository struct {
	mu             sync.RWMutex
	portalPages    map[int]*domain.PortalPage
	nextPortalID   int
	nextLinkID     int
	slugIndex      map[string]int // slug -> portal page ID
	userIndex      map[int][]int  // user ID -> []portal page IDs
}

// NewInMemoryPortalPageRepository creates a new in-memory repository instance
func NewInMemoryPortalPageRepository() *InMemoryPortalPageRepository {
	return &InMemoryPortalPageRepository{
		portalPages: make(map[int]*domain.PortalPage),
		nextPortalID: 1,
		nextLinkID:   1,
		slugIndex:   make(map[string]int),
		userIndex:   make(map[int][]int),
	}
}

// Create creates a new Portal Page in memory
func (r *InMemoryPortalPageRepository) Create(ctx context.Context, portalPage *domain.PortalPage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if slug already exists
	if _, exists := r.slugIndex[portalPage.Slug]; exists {
		return fmt.Errorf("duplicate key value violates unique constraint \"portal_pages_slug_key\"")
	}

	// Create a deep copy to avoid reference issues
	newPortalPage := r.deepCopyPortalPage(portalPage)
	
	// Assign auto-increment ID
	newPortalPage.ID = r.nextPortalID
	r.nextPortalID++

	// Set timestamps if not already set
	now := time.Now().UTC()
	if newPortalPage.CreatedAt.IsZero() {
		newPortalPage.CreatedAt = now
	}
	if newPortalPage.UpdatedAt.IsZero() {
		newPortalPage.UpdatedAt = now
	}

	// Process links
	for _, link := range newPortalPage.Links {
		link.ID = r.nextLinkID
		r.nextLinkID++
		link.PortalPageID = newPortalPage.ID
		
		if link.CreatedAt.IsZero() {
			link.CreatedAt = now
		}
		if link.UpdatedAt.IsZero() {
			link.UpdatedAt = now
		}
	}

	// Store the portal page
	r.portalPages[newPortalPage.ID] = newPortalPage
	
	// Update indexes
	r.slugIndex[newPortalPage.Slug] = newPortalPage.ID
	r.userIndex[newPortalPage.UserID] = append(r.userIndex[newPortalPage.UserID], newPortalPage.ID)

	// Update the original object with generated values
	portalPage.ID = newPortalPage.ID
	portalPage.CreatedAt = newPortalPage.CreatedAt
	portalPage.UpdatedAt = newPortalPage.UpdatedAt
	for i, link := range portalPage.Links {
		link.ID = newPortalPage.Links[i].ID
		link.PortalPageID = newPortalPage.Links[i].PortalPageID
		link.CreatedAt = newPortalPage.Links[i].CreatedAt
		link.UpdatedAt = newPortalPage.Links[i].UpdatedAt
	}

	return nil
}

// Update updates an existing Portal Page in memory
func (r *InMemoryPortalPageRepository) Update(ctx context.Context, portalPage *domain.PortalPage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if portal page exists
	existing, exists := r.portalPages[portalPage.ID]
	if !exists {
		return sql.ErrNoRows
	}

	// Check slug uniqueness (only if slug has changed)
	if existing.Slug != portalPage.Slug {
		if existingID, slugExists := r.slugIndex[portalPage.Slug]; slugExists && existingID != portalPage.ID {
			return fmt.Errorf("duplicate key value violates unique constraint \"portal_pages_slug_key\"")
		}
	}

	// Create a deep copy
	updated := r.deepCopyPortalPage(portalPage)
	updated.UpdatedAt = time.Now().UTC()

	// Build map of existing links by ID
	existingLinkMap := make(map[int]*domain.Link)
	for _, link := range existing.Links {
		existingLinkMap[link.ID] = link
	}

	// Build set of new link IDs
	newLinkIDs := make(map[int]bool)

	// Process links: update existing or create new
	for _, link := range updated.Links {
		newLinkIDs[link.ID] = true

		if link.ID > 0 {
			// Update existing link
			if existingLink, exists := existingLinkMap[link.ID]; exists {
				// Copy over creation time from existing link
				link.CreatedAt = existingLink.CreatedAt
				link.UpdatedAt = time.Now().UTC()
				link.PortalPageID = updated.ID
			}
		} else {
			// New link - assign ID
			link.ID = r.nextLinkID
			r.nextLinkID++
			link.PortalPageID = updated.ID
			
			now := time.Now().UTC()
			if link.CreatedAt.IsZero() {
				link.CreatedAt = now
			}
			if link.UpdatedAt.IsZero() {
				link.UpdatedAt = now
			}
		}
	}

	// Update indexes if slug changed
	if existing.Slug != updated.Slug {
		delete(r.slugIndex, existing.Slug)
		r.slugIndex[updated.Slug] = updated.ID
	}

	// Store updated portal page
	r.portalPages[updated.ID] = updated

	// Update the original object with generated values
	portalPage.UpdatedAt = updated.UpdatedAt
	for i, link := range portalPage.Links {
		link.ID = updated.Links[i].ID
		link.PortalPageID = updated.Links[i].PortalPageID
		link.CreatedAt = updated.Links[i].CreatedAt
		link.UpdatedAt = updated.Links[i].UpdatedAt
	}

	return nil
}

// FindBySlug finds a Portal Page by slug with associated links sorted by display_order ASC
func (r *InMemoryPortalPageRepository) FindBySlug(ctx context.Context, slug string) (*domain.PortalPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	portalPageID, exists := r.slugIndex[slug]
	if !exists {
		return nil, sql.ErrNoRows
	}

	portalPage, exists := r.portalPages[portalPageID]
	if !exists {
		return nil, sql.ErrNoRows
	}

	// Return a deep copy with links sorted by DisplayOrder
	result := r.deepCopyPortalPage(portalPage)
	r.sortLinksByDisplayOrder(result.Links)
	
	return result, nil
}

// ListByUserID returns all Portal Pages for a user (without Links) ordered by created_at ASC
func (r *InMemoryPortalPageRepository) ListByUserID(ctx context.Context, userID int) ([]*domain.PortalPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	portalPageIDs := r.userIndex[userID]
	if len(portalPageIDs) == 0 {
		return []*domain.PortalPage{}, nil
	}

	result := make([]*domain.PortalPage, 0, len(portalPageIDs))
	for _, id := range portalPageIDs {
		if portalPage, exists := r.portalPages[id]; exists {
			// Create copy without links
			copy := &domain.PortalPage{
				ID:              portalPage.ID,
				UserID:          portalPage.UserID,
				Slug:            portalPage.Slug,
				Title:           portalPage.Title,
				Bio:             portalPage.Bio,
				ProfileImageURL: portalPage.ProfileImageURL,
				Theme:           portalPage.Theme,
				CreatedAt:       portalPage.CreatedAt,
				UpdatedAt:       portalPage.UpdatedAt,
				Links:           nil, // No links in list view
			}
			result = append(result, copy)
		}
	}

	// Sort by CreatedAt ASC
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

// FindByID finds a Portal Page by ID with associated links sorted by display_order ASC
func (r *InMemoryPortalPageRepository) FindByID(ctx context.Context, id int) (*domain.PortalPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	portalPage, exists := r.portalPages[id]
	if !exists {
		return nil, sql.ErrNoRows
	}

	// Return a deep copy with links sorted by DisplayOrder
	result := r.deepCopyPortalPage(portalPage)
	r.sortLinksByDisplayOrder(result.Links)
	
	return result, nil
}

// Helper methods

// deepCopyPortalPage creates a deep copy of a PortalPage
func (r *InMemoryPortalPageRepository) deepCopyPortalPage(original *domain.PortalPage) *domain.PortalPage {
	copy := &domain.PortalPage{
		ID:              original.ID,
		UserID:          original.UserID,
		Slug:            original.Slug,
		Title:           original.Title,
		Bio:             original.Bio,
		ProfileImageURL: original.ProfileImageURL,
		Theme:           original.Theme,
		CreatedAt:       original.CreatedAt,
		UpdatedAt:       original.UpdatedAt,
		Links:           make([]*domain.Link, len(original.Links)),
	}

	// Deep copy links
	for i, link := range original.Links {
		copy.Links[i] = &domain.Link{
			ID:           link.ID,
			PortalPageID: link.PortalPageID,
			Title:        link.Title,
			URL:          link.URL,
			Description:  link.Description,
			IconURL:      link.IconURL,
			DisplayOrder: link.DisplayOrder,
			CreatedAt:    link.CreatedAt,
			UpdatedAt:    link.UpdatedAt,
		}
	}

	return copy
}

// sortLinksByDisplayOrder sorts links by DisplayOrder in ascending order
func (r *InMemoryPortalPageRepository) sortLinksByDisplayOrder(links []*domain.Link) {
	sort.Slice(links, func(i, j int) bool {
		return links[i].DisplayOrder < links[j].DisplayOrder
	})
}