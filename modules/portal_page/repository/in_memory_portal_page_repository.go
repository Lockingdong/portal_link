package repository

import (
	"portal_link/modules/portal_page/domain"
)

type InMemoryPortalPageRepository struct {
	portalPages map[int]*domain.PortalPage
}

func NewInMemoryPortalPageRepository() *InMemoryPortalPageRepository {
	return &InMemoryPortalPageRepository{portalPages: make(map[int]*domain.PortalPage)}
}
