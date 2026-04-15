package store

import (
	"context"

	"go-bookmark-manager/src/backend/models"
)

// Store defines the interface for bookmark data persistence
type Store interface {
	// Create adds a new bookmark to the store
	Create(ctx context.Context, bookmark *models.Bookmark) error

	// GetByID retrieves a bookmark by its ID
	GetByID(ctx context.Context, id string) (*models.Bookmark, error)

	// GetAll retrieves all bookmarks with pagination
	// Returns the bookmarks and total count
	GetAll(ctx context.Context, page, pageSize int) ([]models.Bookmark, int, error)

	// Update modifies an existing bookmark
	Update(ctx context.Context, id string, bookmark *models.Bookmark) error

	// Delete removes a bookmark by its ID
	Delete(ctx context.Context, id string) error

	// Search finds bookmarks matching a query in title, description, or URL
	Search(ctx context.Context, query string, page, pageSize int) ([]models.Bookmark, int, error)

	// GetByTag retrieves bookmarks filtered by a specific tag
	GetByTag(ctx context.Context, tag string, page, pageSize int) ([]models.Bookmark, int, error)

	// GetAllTags returns all tags with their bookmark counts
	GetAllTags(ctx context.Context) (map[string]int, error)
}
