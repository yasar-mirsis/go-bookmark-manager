package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-bookmark-manager/src/backend/models"
	"go-bookmark-manager/src/backend/store"
)

// Handler handles HTTP requests for bookmark operations
type Handler struct {
	store store.Store
}

// NewHandler creates a new Handler with the given store
func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// writeJSON writes a JSON response with the given status code and data
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// writeError writes an error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// CreateBookmark handles POST /bookmarks
func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var bookmarkData struct {
		URL         string   `json:"url"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bookmarkData); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if bookmarkData.URL == "" {
		writeError(w, http.StatusBadRequest, "URL is required")
		return
	}

	if bookmarkData.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	// Validate URL format
	if !strings.HasPrefix(bookmarkData.URL, "http://") && !strings.HasPrefix(bookmarkData.URL, "https://") {
		writeError(w, http.StatusBadRequest, "Invalid URL format. URL must start with http:// or https://")
		return
	}

	// Parse tags if provided as comma-separated string
	tags := bookmarkData.Tags
	if len(tags) == 0 {
		// Check if tags were sent as a string
		// This handles the case where frontend sends tags as comma-separated string
	} else {
		// Clean up tags
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// Generate unique ID
	id := generateID()

	bookmark := models.NewBookmark(id, bookmarkData.URL, bookmarkData.Title, bookmarkData.Description, tags)

	if err := h.store.Create(r.Context(), bookmark); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create bookmark")
		return
	}

	writeJSON(w, http.StatusCreated, bookmark)
}

// GetBookmark handles GET /bookmarks/{id}
func (h *Handler) GetBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		// Try to extract from URL path
		path := strings.TrimPrefix(r.URL.Path, "/bookmarks/")
		id = strings.TrimPrefix(path, "/")
	}

	if id == "" {
		writeError(w, http.StatusBadRequest, "Bookmark ID is required")
		return
	}

	bookmark, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Bookmark not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get bookmark")
		return
	}

	writeJSON(w, http.StatusOK, bookmark)
}

// GetBookmarks handles GET /bookmarks
func (h *Handler) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bookmarks, total, err := h.store.GetAll(r.Context(), page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get bookmarks")
		return
	}

	response := struct {
		Bookmarks []models.Bookmark `json:"bookmarks"`
		Total     int               `json:"total"`
		Page      int               `json:"page"`
		PageSize  int               `json:"pageSize"`
	}{
		Bookmarks: bookmarks,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}

	writeJSON(w, http.StatusOK, response)
}

// UpdateBookmark handles PUT /bookmarks/{id}
func (h *Handler) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		path := strings.TrimPrefix(r.URL.Path, "/bookmarks/")
		id = strings.TrimPrefix(path, "/")
	}

	if id == "" {
		writeError(w, http.StatusBadRequest, "Bookmark ID is required")
		return
	}

	var bookmarkData struct {
		URL         string   `json:"url"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bookmarkData); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if bookmarkData.URL == "" {
		writeError(w, http.StatusBadRequest, "URL is required")
		return
	}

	if bookmarkData.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	// Validate URL format
	if !strings.HasPrefix(bookmarkData.URL, "http://") && !strings.HasPrefix(bookmarkData.URL, "https://") {
		writeError(w, http.StatusBadRequest, "Invalid URL format. URL must start with http:// or https://")
		return
	}

	// Clean up tags
	tags := bookmarkData.Tags
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	bookmark := models.NewBookmark(id, bookmarkData.URL, bookmarkData.Title, bookmarkData.Description, tags)

	if err := h.store.Update(r.Context(), id, bookmark); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Bookmark not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update bookmark")
		return
	}

	writeJSON(w, http.StatusOK, bookmark)
}

// DeleteBookmark handles DELETE /bookmarks/{id}
func (h *Handler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		path := strings.TrimPrefix(r.URL.Path, "/bookmarks/")
		id = strings.TrimPrefix(path, "/")
	}

	if id == "" {
		writeError(w, http.StatusBadRequest, "Bookmark ID is required")
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Bookmark not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete bookmark")
		return
	}

	writeJSON(w, http.StatusOK, SuccessResponse{Message: "Bookmark deleted successfully"})
}

// SearchBookmarks handles GET /bookmarks/search
func (h *Handler) SearchBookmarks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bookmarks, total, err := h.store.Search(r.Context(), query, page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to search bookmarks")
		return
	}

	response := struct {
		Bookmarks []models.Bookmark `json:"bookmarks"`
		Total     int               `json:"total"`
		Page      int               `json:"page"`
		PageSize  int               `json:"pageSize"`
	}{
		Bookmarks: bookmarks,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}

	writeJSON(w, http.StatusOK, response)
}

// GetBookmarksByTag handles GET /bookmarks/tag/{tag}
func (h *Handler) GetBookmarksByTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		// Try to extract from URL path
		path := strings.TrimPrefix(r.URL.Path, "/bookmarks/tag/")
		tag = strings.TrimPrefix(path, "/")
	}

	if tag == "" {
		writeError(w, http.StatusBadRequest, "Tag is required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bookmarks, total, err := h.store.GetByTag(r.Context(), tag, page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get bookmarks by tag")
		return
	}

	response := struct {
		Bookmarks []models.Bookmark `json:"bookmarks"`
		Total     int               `json:"total"`
		Page      int               `json:"page"`
		PageSize  int               `json:"pageSize"`
	}{
		Bookmarks: bookmarks,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}

	writeJSON(w, http.StatusOK, response)
}

// GetTags handles GET /tags
func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tags, err := h.store.GetAllTags(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get tags")
		return
	}

	writeJSON(w, http.StatusOK, tags)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// generateID generates a unique ID for bookmarks
func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
