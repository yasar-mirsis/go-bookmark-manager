package models

import (
	"time"
)

// Bookmark represents a saved web resource
type Bookmark struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewBookmark creates a new Bookmark with timestamps
func NewBookmark(id, url, title, description string, tags []string) *Bookmark {
	now := time.Now()
	return &Bookmark{
		ID:          id,
		URL:         url,
		Title:       title,
		Description: description,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
