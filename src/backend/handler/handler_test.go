package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"go-bookmark-manager/src/backend/models"
)

// MockStore is a mock implementation of the store.Store interface for testing
type MockStore struct {
	bookmarks map[string]*models.Bookmark
	err       error
}

func NewMockStore() *MockStore {
	return &MockStore{
		bookmarks: make(map[string]*models.Bookmark),
	}
}

func (m *MockStore) Create(ctx context.Context, bookmark *models.Bookmark) error {
	if m.err != nil {
		return m.err
	}
	m.bookmarks[bookmark.ID] = bookmark
	return nil
}

func (m *MockStore) GetByID(ctx context.Context, id string) (*models.Bookmark, error) {
	if m.err != nil {
		return nil, m.err
	}
	bookmark, exists := m.bookmarks[id]
	if !exists {
		return nil, &NotFoundError{Message: "bookmark not found"}
	}
	copy := *bookmark
	return &copy, nil
}

func (m *MockStore) GetAll(ctx context.Context, page, pageSize int) ([]models.Bookmark, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	bookmarks := make([]models.Bookmark, 0, len(m.bookmarks))
	for _, b := range m.bookmarks {
		copy := *b
		bookmarks = append(bookmarks, copy)
	}
	return bookmarks, len(bookmarks), nil
}

func (m *MockStore) Update(ctx context.Context, id string, bookmark *models.Bookmark) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.bookmarks[id]; !exists {
		return &NotFoundError{Message: "bookmark not found"}
	}
	m.bookmarks[id] = bookmark
	return nil
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.bookmarks[id]; !exists {
		return &NotFoundError{Message: "bookmark not found"}
	}
	delete(m.bookmarks, id)
	return nil
}

func (m *MockStore) Search(ctx context.Context, query string, page, pageSize int) ([]models.Bookmark, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	matches := make([]models.Bookmark, 0)
	for _, b := range m.bookmarks {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(b.Description), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(b.URL), strings.ToLower(query)) {
			copy := *b
			matches = append(matches, copy)
		}
	}
	return matches, len(matches), nil
}

func (m *MockStore) GetByTag(ctx context.Context, tag string, page, pageSize int) ([]models.Bookmark, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	matches := make([]models.Bookmark, 0)
	for _, b := range m.bookmarks {
		for _, t := range b.Tags {
			if strings.ToLower(t) == strings.ToLower(tag) {
				copy := *b
				matches = append(matches, copy)
				break
			}
		}
	}
	return matches, len(matches), nil
}

func (m *MockStore) GetAllTags(ctx context.Context) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	tagCounts := make(map[string]int)
	for _, b := range m.bookmarks {
		for _, tag := range b.Tags {
			tagCounts[tag]++
		}
	}
	return tagCounts, nil
}

// NotFoundError represents a 404 error
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// ErrorStore is a mock store that returns errors
type ErrorStore struct{}

func (e *ErrorStore) Create(ctx context.Context, bookmark *models.Bookmark) error {
	return &InternalError{Message: "store error"}
}

func (e *ErrorStore) GetByID(ctx context.Context, id string) (*models.Bookmark, error) {
	return nil, &InternalError{Message: "store error"}
}

func (e *ErrorStore) GetAll(ctx context.Context, page, pageSize int) ([]models.Bookmark, int, error) {
	return nil, 0, &InternalError{Message: "store error"}
}

func (e *ErrorStore) Update(ctx context.Context, id string, bookmark *models.Bookmark) error {
	return &InternalError{Message: "store error"}
}

func (e *ErrorStore) Delete(ctx context.Context, id string) error {
	return &InternalError{Message: "store error"}
}

func (e *ErrorStore) Search(ctx context.Context, query string, page, pageSize int) ([]models.Bookmark, int, error) {
	return nil, 0, &InternalError{Message: "store error"}
}

func (e *ErrorStore) GetByTag(ctx context.Context, tag string, page, pageSize int) ([]models.Bookmark, int, error) {
	return nil, 0, &InternalError{Message: "store error"}
}

func (e *ErrorStore) GetAllTags(ctx context.Context) (map[string]int, error) {
	return nil, &InternalError{Message: "store error"}
}

// InternalError represents a 500 error
type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

func TestHandler_CreateBookmark(t *testing.T) {
	t.Run("valid request returns 201", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url":         "https://example.com",
			"title":       "Example",
			"description": "An example bookmark",
			"tags":        []string{"example", "test"},
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response models.Bookmark
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.URL != "https://example.com" {
			t.Errorf("Expected URL https://example.com, got %s", response.URL)
		}
		if response.Title != "Example" {
			t.Errorf("Expected title Example, got %s", response.Title)
		}
	})

	t.Run("valid request with string tags returns 201", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url":         "https://example.com",
			"title":       "Example",
			"description": "",
			"tags":        "tag1,tag2,tag3",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response models.Bookmark
		json.Unmarshal(w.Body.Bytes(), &response)
		if len(response.Tags) != 3 {
			t.Errorf("Expected 3 tags, got %d", len(response.Tags))
		}
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing URL returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"title": "Example",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing title returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url": "https://example.com",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid URL format returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url":  "not-a-valid-url",
			"title": "Example",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		bookmarkData := map[string]interface{}{
			"url":   "https://example.com",
			"title": "Example",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
		w := httptest.NewRecorder()

		handler.CreateBookmark(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_GetBookmark(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		bookmark := models.NewBookmark("test-id", "https://example.com", "Example", "Description", []string{"tag1"})
		mockStore.bookmarks["test-id"] = bookmark

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.GetBookmark(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response models.Bookmark
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.ID != "test-id" {
			t.Errorf("Expected ID test-id, got %s", response.ID)
		}
	})

	t.Run("non-existent bookmark returns 404", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?id=non-existent", nil)
		w := httptest.NewRecorder()

		handler.GetBookmark(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("missing ID returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
		w := httptest.NewRecorder()

		handler.GetBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.GetBookmark(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.GetBookmark(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_GetBookmarks(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["1"] = models.NewBookmark("1", "https://example1.com", "Example 1", "", nil)
		mockStore.bookmarks["2"] = models.NewBookmark("2", "https://example2.com", "Example 2", "", nil)

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["total"].(float64)) != 2 {
			t.Errorf("Expected total 2, got %v", response["total"])
		}
	})

	t.Run("empty result returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["total"].(float64)) != 0 {
			t.Errorf("Expected total 0, got %v", response["total"])
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarks(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/bookmarks", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarks(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_UpdateBookmark(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["test-id"] = models.NewBookmark("test-id", "https://old.com", "Old Title", "", nil)

		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url":   "https://new.com",
			"title": "New Title",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPut, "/bookmarks?id=test-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateBookmark(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response models.Bookmark
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.URL != "https://new.com" {
			t.Errorf("Expected URL https://new.com, got %s", response.URL)
		}
	})

	t.Run("non-existent bookmark returns 404", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url":   "https://example.com",
			"title": "Example",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPut, "/bookmarks?id=non-existent", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateBookmark(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("missing ID returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"url":   "https://example.com",
			"title": "Example",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPut, "/bookmarks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing URL returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["test-id"] = models.NewBookmark("test-id", "https://old.com", "Old Title", "", nil)

		handler := NewHandler(mockStore)

		bookmarkData := map[string]interface{}{
			"title": "New Title",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPut, "/bookmarks?id=test-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		bookmarkData := map[string]interface{}{
			"url":   "https://example.com",
			"title": "Example",
		}

		body, _ := json.Marshal(bookmarkData)
		req := httptest.NewRequest(http.MethodPut, "/bookmarks?id=test-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateBookmark(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.UpdateBookmark(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_DeleteBookmark(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["test-id"] = models.NewBookmark("test-id", "https://example.com", "Example", "", nil)

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodDelete, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.DeleteBookmark(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("non-existent bookmark returns 404", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodDelete, "/bookmarks?id=non-existent", nil)
		w := httptest.NewRecorder()

		handler.DeleteBookmark(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("missing ID returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodDelete, "/bookmarks", nil)
		w := httptest.NewRecorder()

		handler.DeleteBookmark(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		req := httptest.NewRequest(http.MethodDelete, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.DeleteBookmark(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?id=test-id", nil)
		w := httptest.NewRecorder()

		handler.DeleteBookmark(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_SearchBookmarks(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["1"] = models.NewBookmark("1", "https://example.com", "Go Programming", "Learn Go", nil)

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/search?q=go&page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.SearchBookmarks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["total"].(float64)) != 1 {
			t.Errorf("Expected total 1, got %v", response["total"])
		}
	})

	t.Run("missing query returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/search?page=1", nil)
		w := httptest.NewRecorder()

		handler.SearchBookmarks(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("no results returns 200 with empty array", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/search?q=nonexistent&page=1", nil)
		w := httptest.NewRecorder()

		handler.SearchBookmarks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["total"].(float64)) != 0 {
			t.Errorf("Expected total 0, got %v", response["total"])
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/search?q=test&page=1", nil)
		w := httptest.NewRecorder()

		handler.SearchBookmarks(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/bookmarks/search?q=test", nil)
		w := httptest.NewRecorder()

		handler.SearchBookmarks(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_GetBookmarksByTag(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["1"] = models.NewBookmark("1", "https://example.com", "Example", "", []string{"go", "programming"})

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/tag/go?tag=go&page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarksByTag(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["total"].(float64)) != 1 {
			t.Errorf("Expected total 1, got %v", response["total"])
		}
	})

	t.Run("missing tag returns 400", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/tag?page=1", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarksByTag(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("no results returns 200 with empty array", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/tag?tag=nonexistent&page=1", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarksByTag(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["total"].(float64)) != 0 {
			t.Errorf("Expected total 0, got %v", response["total"])
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks/tag?tag=test&page=1", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarksByTag(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/bookmarks/tag?tag=go", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarksByTag(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_GetTags(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		mockStore.bookmarks["1"] = models.NewBookmark("1", "https://example.com", "Example 1", "", []string{"go", "programming"})
		mockStore.bookmarks["2"] = models.NewBookmark("2", "https://example.com", "Example 2", "", []string{"go", "docs"})

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/tags", nil)
		w := httptest.NewRecorder()

		handler.GetTags(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]int
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["go"] != 2 {
			t.Errorf("Expected go count 2, got %d", response["go"])
		}
		if response["programming"] != 1 {
			t.Errorf("Expected programming count 1, got %d", response["programming"])
		}
	})

	t.Run("empty tags returns 200 with empty object", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/tags", nil)
		w := httptest.NewRecorder()

		handler.GetTags(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]int
		json.Unmarshal(w.Body.Bytes(), &response)
		if len(response) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(response))
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		errorStore := &ErrorStore{}
		handler := NewHandler(errorStore)

		req := httptest.NewRequest(http.MethodGet, "/tags", nil)
		w := httptest.NewRecorder()

		handler.GetTags(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/tags", nil)
		w := httptest.NewRecorder()

		handler.GetTags(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestHandler_HealthCheck(t *testing.T) {
	t.Run("valid request returns 200", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.HealthCheck(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["status"] != "ok" {
			t.Errorf("Expected status ok, got %s", response["status"])
		}
	})

	t.Run("method not allowed returns 405", func(t *testing.T) {
		mockStore := NewMockStore()
		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		w := httptest.NewRecorder()

		handler.HealthCheck(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

// Helper function to create test bookmarks with specific IDs
func createTestBookmark(id, url, title, description string, tags []string) *models.Bookmark {
	return models.NewBookmark(id, url, title, description, tags)
}

// Test pagination edge cases
func TestHandler_PaginationEdgeCases(t *testing.T) {
	t.Run("page 0 defaults to 1", func(t *testing.T) {
		mockStore := NewMockStore()
		for i := 0; i < 25; i++ {
			id := strconv.Itoa(i)
			mockStore.bookmarks[id] = models.NewBookmark(id, "https://example"+id+".com", "Title "+id, "", nil)
		}

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?page=0&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["page"].(float64)) != 1 {
			t.Errorf("Expected page 1 (default), got %v", response["page"])
		}
	})

	t.Run("pageSize 0 defaults to 10", func(t *testing.T) {
		mockStore := NewMockStore()
		for i := 0; i < 25; i++ {
			id := strconv.Itoa(i)
			mockStore.bookmarks[id] = models.NewBookmark(id, "https://example"+id+".com", "Title "+id, "", nil)
		}

		handler := NewHandler(mockStore)

		req := httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=0", nil)
		w := httptest.NewRecorder()

		handler.GetBookmarks(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if int(response["pageSize"].(float64)) != 10 {
			t.Errorf("Expected pageSize 10 (default), got %v", response["pageSize"])
		}
	})
}
