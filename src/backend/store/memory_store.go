package store

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"go-bookmark-manager/src/backend/models"
)

// MemoryStore is an in-memory implementation of the Store interface
type MemoryStore struct {
	mu        sync.RWMutex
	bookmarks map[string]*models.Bookmark
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookmarks: make(map[string]*models.Bookmark),
	}
}

// Create adds a new bookmark to the store
func (s *MemoryStore) Create(ctx context.Context, bookmark *models.Bookmark) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if bookmark.ID == "" {
		return errors.New("bookmark ID cannot be empty")
	}

	if bookmark.URL == "" {
		return errors.New("bookmark URL cannot be empty")
	}

	if bookmark.Title == "" {
		return errors.New("bookmark title cannot be empty")
	}

	s.bookmarks[bookmark.ID] = bookmark
	return nil
}

// GetByID retrieves a bookmark by its ID
func (s *MemoryStore) GetByID(ctx context.Context, id string) (*models.Bookmark, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	bookmark, exists := s.bookmarks[id]
	if !exists {
		return nil, errors.New("bookmark not found")
	}

	// Return a copy to prevent external modification
	copy := *bookmark
	return &copy, nil
}

// GetAll retrieves all bookmarks with pagination
func (s *MemoryStore) GetAll(ctx context.Context, page, pageSize int) ([]models.Bookmark, int, error) {
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	bookmarks := make([]models.Bookmark, 0, len(s.bookmarks))
	for _, b := range s.bookmarks {
		copy := *b
		bookmarks = append(bookmarks, copy)
	}

	total := len(bookmarks)

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []models.Bookmark{}, total, nil
	}

	if end > total {
		end = total
	}

	return bookmarks[start:end], total, nil
}

// Update modifies an existing bookmark
func (s *MemoryStore) Update(ctx context.Context, id string, bookmark *models.Bookmark) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.bookmarks[id]; !exists {
		return errors.New("bookmark not found")
	}

	// Update the timestamp to reflect when the update occurred
	bookmark.UpdatedAt = time.Now()
	s.bookmarks[id] = bookmark
	return nil
}

// Delete removes a bookmark by its ID
func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.bookmarks[id]; !exists {
		return errors.New("bookmark not found")
	}

	delete(s.bookmarks, id)
	return nil
}

// Search finds bookmarks matching a query in title, description, or URL
func (s *MemoryStore) Search(ctx context.Context, query string, page, pageSize int) ([]models.Bookmark, int, error) {
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	matches := make([]models.Bookmark, 0)

	for _, b := range s.bookmarks {
		if strings.Contains(strings.ToLower(b.Title), query) ||
			strings.Contains(strings.ToLower(b.Description), query) ||
			strings.Contains(strings.ToLower(b.URL), query) {
			copy := *b
			matches = append(matches, copy)
		}
	}

	total := len(matches)

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []models.Bookmark{}, total, nil
	}

	if end > total {
		end = total
	}

	return matches[start:end], total, nil
}

// GetByTag retrieves bookmarks filtered by a specific tag
func (s *MemoryStore) GetByTag(ctx context.Context, tag string, page, pageSize int) ([]models.Bookmark, int, error) {
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	tag = strings.ToLower(tag)
	matches := make([]models.Bookmark, 0)

	for _, b := range s.bookmarks {
		for _, t := range b.Tags {
			if strings.ToLower(t) == tag {
				copy := *b
				matches = append(matches, copy)
				break
			}
		}
	}

	total := len(matches)

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []models.Bookmark{}, total, nil
	}

	if end > total {
		end = total
	}

	return matches[start:end], total, nil
}

// GetAllTags returns all tags with their bookmark counts
func (s *MemoryStore) GetAllTags(ctx context.Context) (map[string]int, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	tagCounts := make(map[string]int)

	for _, b := range s.bookmarks {
		for _, tag := range b.Tags {
			tagCounts[tag]++
		}
	}

	return tagCounts, nil
}
