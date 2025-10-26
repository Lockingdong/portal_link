package repository

import (
	"context"
	"database/sql"
	"portal_link/modules/user/domain"
	"sync"
)

var _ domain.UserRepository = (*InMemoryUserRepository)(nil)

// InMemoryUserRepository is an in-memory implementation of UserRepository for testing
type InMemoryUserRepository struct {
	mu      sync.RWMutex
	users   map[int]*domain.User
	emails  map[string]int // email -> user ID mapping
	nextID  int
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*domain.User),
		emails: make(map[string]int),
		nextID: 1,
	}
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	if _, exists := r.emails[user.Email]; exists {
		return domain.ErrEmailExists
	}

	// Assign ID if not set
	if user.ID == 0 {
		user.ID = r.nextID
		r.nextID++
	} else {
		// If ID is provided, update nextID if necessary
		if user.ID >= r.nextID {
			r.nextID = user.ID + 1
		}
	}

	// Store user
	r.users[user.ID] = user
	r.emails[user.Email] = user.ID

	return nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return nil, sql.ErrNoRows
	}

	user := r.users[userID]
	return user, nil
}

// Find retrieves a user by ID
func (r *InMemoryUserRepository) Find(ctx context.Context, id int) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

// Reset clears all data (useful for testing)
func (r *InMemoryUserRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users = make(map[int]*domain.User)
	r.emails = make(map[string]int)
	r.nextID = 1
}
