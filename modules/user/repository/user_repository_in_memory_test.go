package repository

import (
	"context"
	"database/sql"
	"portal_link/modules/user/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := NewInMemoryUserRepository()
	ctx := context.Background()

	t.Run("successfully creates user with auto-generated ID", func(t *testing.T) {
		repo.Reset()

		user := &domain.User{
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Equal(t, 1, user.ID, "ID should be auto-generated")

		// Verify user can be retrieved
		retrieved, err := repo.Find(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Name, retrieved.Name)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("successfully creates user with provided ID", func(t *testing.T) {
		repo.Reset()

		user := &domain.User{
			ID:        100,
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Equal(t, 100, user.ID, "ID should remain as provided")

		// Next auto-generated ID should be 101
		user2 := &domain.User{
			Name:      "Test User 2",
			Email:     "test2@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err = repo.Create(ctx, user2)
		require.NoError(t, err)
		assert.Equal(t, 101, user2.ID)
	})

	t.Run("returns error when email already exists", func(t *testing.T) {
		repo.Reset()

		user1 := &domain.User{
			Name:      "User 1",
			Email:     "duplicate@example.com",
			Password:  "password1",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &domain.User{
			Name:      "User 2",
			Email:     "duplicate@example.com",
			Password:  "password2",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err = repo.Create(ctx, user2)
		assert.ErrorIs(t, err, domain.ErrEmailExists)
	})
}

func TestInMemoryUserRepository_GetByEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()
	ctx := context.Background()

	t.Run("successfully retrieves user by email", func(t *testing.T) {
		repo.Reset()

		user := &domain.User{
			Name:      "Test User",
			Email:     "findme@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetByEmail(ctx, "findme@example.com")
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Name, retrieved.Name)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("returns error when email not found", func(t *testing.T) {
		repo.Reset()

		_, err := repo.GetByEmail(ctx, "notfound@example.com")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestInMemoryUserRepository_Find(t *testing.T) {
	repo := NewInMemoryUserRepository()
	ctx := context.Background()

	t.Run("successfully retrieves user by ID", func(t *testing.T) {
		repo.Reset()

		user := &domain.User{
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.Find(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Name, retrieved.Name)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("returns error when ID not found", func(t *testing.T) {
		repo.Reset()

		_, err := repo.Find(ctx, 9999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestInMemoryUserRepository_Reset(t *testing.T) {
	repo := NewInMemoryUserRepository()
	ctx := context.Background()

	t.Run("clears all data", func(t *testing.T) {
		// Create some users
		user1 := &domain.User{
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "password1",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &domain.User{
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  "password2",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err = repo.Create(ctx, user2)
		require.NoError(t, err)

		// Reset
		repo.Reset()

		// Verify users are gone
		_, err = repo.Find(ctx, user1.ID)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		_, err = repo.GetByEmail(ctx, user1.Email)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		// Verify ID counter is reset
		user3 := &domain.User{
			Name:      "User 3",
			Email:     "user3@example.com",
			Password:  "password3",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err = repo.Create(ctx, user3)
		require.NoError(t, err)
		assert.Equal(t, 1, user3.ID, "ID should start from 1 again")
	})
}

func TestInMemoryUserRepository_Concurrency(t *testing.T) {
	repo := NewInMemoryUserRepository()
	ctx := context.Background()

	t.Run("handles concurrent operations safely", func(t *testing.T) {
		repo.Reset()

		// Create users concurrently
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func(index int) {
				user := &domain.User{
					Name:      "User",
					Email:     "user" + string(rune('0'+index)) + "@example.com",
					Password:  "password",
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				}
				err := repo.Create(ctx, user)
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify all users were created
		assert.Equal(t, 10, len(repo.users))
		assert.Equal(t, 10, len(repo.emails))
	})
}
