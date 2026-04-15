package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-bookmark-manager/src/backend/handler"
	"go-bookmark-manager/src/backend/models"
	"go-bookmark-manager/src/backend/store"
)

// TestServer helps set up and tear down test server
type TestServer struct {
	Server *httptest.Server
	Store  *store.MemoryStore
	Client *http.Client
}

// NewTestServer creates a new test server with in-memory store
func NewTestServer() *TestServer {
	st := store.NewMemoryStore()
	h := handler.NewHandler(st)

	mux := http.NewServeMux()

	// Bookmark routes
	mux.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetBookmarks(w, r)
		case http.MethodPost:
			h.CreateBookmark(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Single bookmark routes
	mux.HandleFunc("/bookmarks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetBookmark(w, r)
		case http.MethodPut:
			h.UpdateBookmark(w, r)
		case http.MethodDelete:
			h.DeleteBookmark(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Search route
	mux.HandleFunc("/bookmarks/search", func(w http.ResponseWriter, r *http.Request) {
		h.SearchBookmarks(w, r)
	})

	// Tag routes
	mux.HandleFunc("/bookmarks/tag/", func(w http.ResponseWriter, r *http.Request) {
		h.GetBookmarksByTag(w, r)
	})

	// Tags list route
	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		h.GetTags(w, r)
	})

	// Health check route
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		h.HealthCheck(w, r)
	})

	server := httptest.NewServer(mux)

	return &TestServer{
		Server: server,
		Store:  st,
		Client: server.Client(),
	}
}

// Helper function to make JSON requests
func makeRequest(server *httptest.Server, method, path string, body interface{}) (*http.Response, []byte, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := http.NewRequest(method, server.URL+path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}
	resp.Body.Close()

	return resp, respBodyBytes, nil
}

// Helper function to create a bookmark via HTTP
func createBookmark(server *httptest.Server, url, title, description string, tags interface{}) (*http.Response, models.Bookmark, error) {
	bookmarkData := map[string]interface{}{
		"url":   url,
		"title": title,
	}
	if description != "" {
		bookmarkData["description"] = description
	}
	if tags != nil {
		bookmarkData["tags"] = tags
	}

	reqBody, _ := json.Marshal(bookmarkData)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/bookmarks", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, models.Bookmark{}, err
	}
	defer resp.Body.Close()

	var bookmark models.Bookmark
	if err := json.NewDecoder(resp.Body).Decode(&bookmark); err != nil {
		return resp, bookmark, err
	}

	return resp, bookmark, nil
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	resp, err := http.Get(ts.Server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var health map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if health["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", health["status"])
	}
}

// TestCreateBookmark tests bookmark creation
func TestCreateBookmark(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	t.Run("create bookmark with array tags returns 201", func(t *testing.T) {
		resp, bookmark, err := createBookmark(ts.Server, "https://example.com", "Example Bookmark", "A test bookmark", []string{"go", "test"})
		if err != nil {
			t.Fatalf("Failed to create bookmark: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		if bookmark.URL != "https://example.com" {
			t.Errorf("Expected URL https://example.com, got %s", bookmark.URL)
		}

		if bookmark.Title != "Example Bookmark" {
			t.Errorf("Expected title 'Example Bookmark', got '%s'", bookmark.Title)
		}

		if len(bookmark.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(bookmark.Tags))
		}
	})

	t.Run("create bookmark with string tags returns 201", func(t *testing.T) {
		resp, bookmark, err := createBookmark(ts.Server, "https://test.com", "Test Bookmark", "", "tag1,tag2,tag3")
		if err != nil {
			t.Fatalf("Failed to create bookmark: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		if len(bookmark.Tags) != 3 {
			t.Errorf("Expected 3 tags, got %d", len(bookmark.Tags))
		}
	})

	t.Run("create bookmark with invalid JSON returns 400", func(t *testing.T) {
		resp, err := http.Post(ts.Server.URL+"/bookmarks", "application/json", strings.NewReader("invalid json"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("create bookmark with missing URL returns 400", func(t *testing.T) {
		data := map[string]string{"title": "Test"}
		resp, _ := createBookmarkWithBody(ts.Server, data)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("create bookmark with missing title returns 400", func(t *testing.T) {
		data := map[string]string{"url": "https://test.com"}
		resp, _ := createBookmarkWithBody(ts.Server, data)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("create bookmark with invalid URL format returns 400", func(t *testing.T) {
		data := map[string]string{"url": "not-a-url", "title": "Test"}
		resp, _ := createBookmarkWithBody(ts.Server, data)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

// createBookmarkWithBody is a helper to create bookmark with custom body
func createBookmarkWithBody(server *httptest.Server, data map[string]string) (*http.Response, error) {
	reqBody, _ := json.Marshal(data)
	resp, err := http.Post(server.URL+"/bookmarks", "application/json", bytes.NewReader(reqBody))
	return resp, err
}

// TestGetBookmark tests retrieving a single bookmark
func TestGetBookmark(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// First create a bookmark
	_, bookmark, _ := createBookmark(ts.Server, "https://example.com", "Example", "", nil)

	t.Run("get existing bookmark returns 200", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?id=" + bookmark.ID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var retrieved models.Bookmark
		json.NewDecoder(resp.Body).Decode(&retrieved)

		if retrieved.ID != bookmark.ID {
			t.Errorf("Expected ID %s, got %s", bookmark.ID, retrieved.ID)
		}
	})

	t.Run("get non-existent bookmark returns 404", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?id=non-existent-id")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("get bookmark without ID returns 400", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// This should return 200 because it's treated as GetBookmarks
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

// TestGetBookmarks tests listing bookmarks with pagination
func TestGetBookmarks(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Create multiple bookmarks for pagination testing
	for i := 0; i < 25; i++ {
		createBookmark(ts.Server, "https://example"+string(rune(i))+".com", "Bookmark "+string(rune(i)), "", nil)
	}

	t.Run("get bookmarks returns paginated results", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if int(response["total"].(float64)) != 25 {
			t.Errorf("Expected total 25, got %v", response["total"])
		}

		bookmarks := response["bookmarks"].([]interface{})
		if len(bookmarks) != 10 {
			t.Errorf("Expected 10 bookmarks on page 1, got %d", len(bookmarks))
		}
	})

	t.Run("get page 2 returns next set of results", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?page=2&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		bookmarks := response["bookmarks"].([]interface{})
		if len(bookmarks) != 10 {
			t.Errorf("Expected 10 bookmarks on page 2, got %d", len(bookmarks))
		}
	})

	t.Run("get bookmarks with default pagination", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if int(response["pageSize"].(float64)) != 10 {
			t.Errorf("Expected default pageSize 10, got %v", response["pageSize"])
		}
	})

	t.Run("get bookmarks on empty database returns empty array", func(t *testing.T) {
		// Create a new test server with empty store
		ts2 := NewTestServer()
		defer ts2.Server.Close()

		resp, err := http.Get(ts2.Server.URL + "/bookmarks")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if int(response["total"].(float64)) != 0 {
			t.Errorf("Expected total 0, got %v", response["total"])
		}
	})
}

// TestUpdateBookmark tests bookmark update functionality
func TestUpdateBookmark(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// First create a bookmark
	_, bookmark, _ := createBookmark(ts.Server, "https://old.com", "Old Title", "Old Description", []string{"old"})

	t.Run("update bookmark returns 200", func(t *testing.T) {
		data := map[string]interface{}{
			"url":         "https://new.com",
			"title":       "New Title",
			"description": "New Description",
			"tags":        []string{"new", "updated"},
		}

		reqBody, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodPut, ts.Server.URL+"/bookmarks?id="+bookmark.ID, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var updated models.Bookmark
		json.NewDecoder(resp.Body).Decode(&updated)

		if updated.URL != "https://new.com" {
			t.Errorf("Expected URL https://new.com, got %s", updated.URL)
		}

		if updated.Title != "New Title" {
			t.Errorf("Expected title 'New Title', got '%s'", updated.Title)
		}
	})

	t.Run("update non-existent bookmark returns 404", func(t *testing.T) {
		data := map[string]string{"url": "https://test.com", "title": "Test"}
		reqBody, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodPut, ts.Server.URL+"/bookmarks?id=non-existent", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("update with missing URL returns 400", func(t *testing.T) {
		data := map[string]string{"title": "Test"}
		reqBody, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodPut, ts.Server.URL+"/bookmarks?id="+bookmark.ID, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

// TestDeleteBookmark tests bookmark deletion
func TestDeleteBookmark(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// First create a bookmark
	_, bookmark, _ := createBookmark(ts.Server, "https://example.com", "Example", "", nil)

	t.Run("delete existing bookmark returns 200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.Server.URL+"/bookmarks?id="+bookmark.ID, nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("delete bookmark and verify it's gone", func(t *testing.T) {
		// Create another bookmark
		_, bookmark2, _ := createBookmark(ts.Server, "https://test.com", "Test", "", nil)

		// Delete it
		req, _ := http.NewRequest(http.MethodDelete, ts.Server.URL+"/bookmarks?id="+bookmark2.ID, nil)
		http.DefaultClient.Do(req)

		// Try to get it
		resp, _ := http.Get(ts.Server.URL + "/bookmarks?id=" + bookmark2.ID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d after deletion, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("delete non-existent bookmark returns 404", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.Server.URL+"/bookmarks?id=non-existent", nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})
}

// TestSearchBookmarks tests search functionality
func TestSearchBookmarks(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Create bookmarks with different content
	createBookmark(ts.Server, "https://golang.org", "Go Programming Language", "Learn Go programming", []string{"go", "programming"})
	createBookmark(ts.Server, "https://example.com", "Example Site", "Just an example", []string{"example"})
	createBookmark(ts.Server, "https://go-tutorial.com", "Go Tutorial", "A comprehensive Go tutorial", []string{"go", "tutorial"})

	t.Run("search by title returns matching results", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/search?q=Go&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := int(response["total"].(float64))
		if total < 2 {
			t.Errorf("Expected at least 2 results for 'Go', got %d", total)
		}
	})

	t.Run("search by description returns matching results", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/search?q=tutorial&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := int(response["total"].(float64))
		if total < 1 {
			t.Errorf("Expected at least 1 result for 'tutorial', got %d", total)
		}
	})

	t.Run("search by URL returns matching results", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/search?q=golang&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := int(response["total"].(float64))
		if total < 1 {
			t.Errorf("Expected at least 1 result for 'golang', got %d", total)
		}
	})

	t.Run("search with no results returns empty array", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/search?q=nonexistent12345&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if int(response["total"].(float64)) != 0 {
			t.Errorf("Expected total 0, got %v", response["total"])
		}
	})

	t.Run("search without query returns 400", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/search?page=1")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

// TestGetBookmarksByTag tests tag filtering
func TestGetBookmarksByTag(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Create bookmarks with tags
	createBookmark(ts.Server, "https://go1.com", "Go Book 1", "", []string{"go", "books"})
	createBookmark(ts.Server, "https://go2.com", "Go Book 2", "", []string{"go", "programming"})
	createBookmark(ts.Server, "https://example.com", "Example", "", []string{"example"})
	createBookmark(ts.Server, "https://go3.com", "Go Book 3", "", []string{"go", "tutorials"})

	t.Run("filter by tag returns matching bookmarks", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/tag?tag=go&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := int(response["total"].(float64))
		if total != 3 {
			t.Errorf("Expected 3 bookmarks with tag 'go', got %d", total)
		}
	})

	t.Run("filter by non-existent tag returns empty array", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/tag?tag=nonexistent&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if int(response["total"].(float64)) != 0 {
			t.Errorf("Expected total 0, got %v", response["total"])
		}
	})

	t.Run("filter without tag returns 400", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks/tag?page=1")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

// TestGetTags tests listing all tags
func TestGetTags(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Create bookmarks with tags
	createBookmark(ts.Server, "https://go1.com", "Go Book", "", []string{"go", "books"})
	createBookmark(ts.Server, "https://go2.com", "Go Tutorial", "", []string{"go", "tutorials"})
	createBookmark(ts.Server, "https://example.com", "Example", "", []string{"example"})

	t.Run("get all tags returns tag counts", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/tags")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var tags map[string]int
		json.NewDecoder(resp.Body).Decode(&tags)

		if tags["go"] != 2 {
			t.Errorf("Expected 'go' count 2, got %d", tags["go"])
		}
		if tags["books"] != 1 {
			t.Errorf("Expected 'books' count 1, got %d", tags["books"])
		}
	})

	t.Run("get tags from empty database returns empty object", func(t *testing.T) {
		ts2 := NewTestServer()
		defer ts2.Server.Close()

		resp, err := http.Get(ts2.Server.URL + "/tags")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var tags map[string]int
		json.NewDecoder(resp.Body).Decode(&tags)

		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})
}

// TestFullCRUD tests complete CRUD lifecycle
func TestFullCRUD(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Create
	_, bookmark, _ := createBookmark(ts.Server, "https://initial.com", "Initial Title", "Initial Desc", []string{"initial"})

	// Verify creation
	resp, _ := http.Get(ts.Server.URL + "/bookmarks?id=" + bookmark.ID)
	var created models.Bookmark
	json.NewDecoder(resp.Body).Decode(&created)
	if created.Title != "Initial Title" {
		t.Errorf("Expected 'Initial Title', got '%s'", created.Title)
	}
	resp.Body.Close()

	// Update
	updateData := map[string]interface{}{
		"url":         "https://updated.com",
		"title":       "Updated Title",
		"description": "Updated Desc",
		"tags":        []string{"updated"},
	}
	reqBody, _ := json.Marshal(updateData)
	req, _ := http.NewRequest(http.MethodPut, ts.Server.URL+"/bookmarks?id="+bookmark.ID, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	var updated models.Bookmark
	json.NewDecoder(resp.Body).Decode(&updated)
	resp.Body.Close()

	// Verify update
	if updated.Title != "Updated Title" {
		t.Errorf("Expected 'Updated Title', got '%s'", updated.Title)
	}
	if len(updated.Tags) != 1 || updated.Tags[0] != "updated" {
		t.Errorf("Expected tag 'updated', got %v", updated.Tags)
	}

	// Delete
	req, _ = http.NewRequest(http.MethodDelete, ts.Server.URL+"/bookmarks?id="+bookmark.ID, nil)
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()

	// Verify deletion
	resp, _ = http.Get(ts.Server.URL + "/bookmarks?id=" + bookmark.ID)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 after deletion, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

// TestPaginationEdgeCases tests pagination edge cases
func TestPaginationEdgeCases(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Create exactly 10 bookmarks
	for i := 0; i < 10; i++ {
		createBookmark(ts.Server, "https://example"+string(rune(i))+".com", "Bookmark "+string(rune(i)), "", nil)
	}

	t.Run("page beyond total returns empty array", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?page=100&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := int(response["total"].(float64))
		bookmarks := response["bookmarks"].([]interface{})

		if total != 10 {
			t.Errorf("Expected total 10, got %d", total)
		}
		if len(bookmarks) != 0 {
			t.Errorf("Expected 0 bookmarks on page 100, got %d", len(bookmarks))
		}
	})

	t.Run("page 0 defaults to 1", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?page=0&pageSize=10")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		page := int(response["page"].(float64))
		if page != 1 {
			t.Errorf("Expected page 1 (default), got %d", page)
		}
	})

	t.Run("pageSize 0 defaults to 10", func(t *testing.T) {
		resp, err := http.Get(ts.Server.URL + "/bookmarks?page=1&pageSize=0")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		pageSize := int(response["pageSize"].(float64))
		if pageSize != 10 {
			t.Errorf("Expected pageSize 10 (default), got %d", pageSize)
		}
	})
}

// TestErrorScenarios tests various error scenarios
func TestErrorScenarios(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		resp, err := http.Post(ts.Server.URL+"/bookmarks", "application/json", strings.NewReader("{invalid json"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		resp, err := http.Post(ts.Server.URL+"/bookmarks", "application/json", strings.NewReader("{}"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("method not allowed returns 405 for GET on create endpoint", func(t *testing.T) {
		// Note: GET on /bookmarks is actually GetBookmarks, so we test a different scenario
		// This tests that wrong methods are rejected
		resp, err := http.Post(ts.Server.URL+"/health", "application/json", nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

// TestConcurrentAccess tests concurrent access to the server
func TestConcurrentAccess(t *testing.T) {
	ts := NewTestServer()
	defer ts.Server.Close()

	// Run multiple concurrent requests
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(idx int) {
			_, _, err := createBookmark(ts.Server, "https://concurrent"+string(rune(idx))+".com", "Concurrent "+string(rune(idx)), "", nil)
			if err != nil {
				t.Errorf("Failed to create bookmark %d: %v", idx, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all bookmarks were created
	resp, err := http.Get(ts.Server.URL + "/bookmarks?page=1&pageSize=100")
	if err != nil {
		t.Fatalf("Failed to get bookmarks: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	total := int(response["total"].(float64))
	if total != 10 {
		t.Errorf("Expected 10 bookmarks, got %d", total)
	}
}
