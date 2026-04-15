package store

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"go-bookmark-manager/src/backend/models"
)

func TestMemoryStore_Create(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	t.Run("Create with valid bookmark", func(t *testing.T) {
		bookmark := models.NewBookmark(
			"test-id-1",
			"https://example.com",
			"Example Title",
			"Example Description",
			[]string{"tag1", "tag2"},
		)

		err := store.Create(ctx, bookmark)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify bookmark was created
		retrieved, err := store.GetByID(ctx, "test-id-1")
		if err != nil {
			t.Errorf("Expected no error retrieving bookmark, got %v", err)
		}

		if retrieved.ID != bookmark.ID {
			t.Errorf("Expected ID %s, got %s", bookmark.ID, retrieved.ID)
		}
		if retrieved.URL != bookmark.URL {
			t.Errorf("Expected URL %s, got %s", bookmark.URL, retrieved.URL)
		}
		if retrieved.Title != bookmark.Title {
			t.Errorf("Expected Title %s, got %s", bookmark.Title, retrieved.Title)
		}
		if retrieved.Description != bookmark.Description {
			t.Errorf("Expected Description %s, got %s", bookmark.Description, retrieved.Description)
		}
		if len(retrieved.Tags) != len(bookmark.Tags) {
			t.Errorf("Expected %d tags, got %d", len(bookmark.Tags), len(retrieved.Tags))
		}
	})

	t.Run("Verify ID generation", func(t *testing.T) {
		id1 := strconv.FormatInt(time.Now().UnixNano(), 36)
		bookmark1 := models.NewBookmark(id1, "https://example1.com", "Title 1", "", nil)

		id2 := strconv.FormatInt(time.Now().UnixNano(), 36)
		bookmark2 := models.NewBookmark(id2, "https://example2.com", "Title 2", "", nil)

		err1 := store.Create(ctx, bookmark1)
		err2 := store.Create(ctx, bookmark2)

		if err1 != nil || err2 != nil {
			t.Errorf("Expected no errors, got err1=%v, err2=%v", err1, err2)
		}

		if id1 == id2 {
			t.Error("Expected unique IDs, got duplicate")
		}
	})

	t.Run("Verify timestamps are set", func(t *testing.T) {
		id := strconv.FormatInt(time.Now().UnixNano(), 36)
		bookmark := models.NewBookmark(id, "https://example.com", "Title", "", nil)

		err := store.Create(ctx, bookmark)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		retrieved, _ := store.GetByID(ctx, id)

		if retrieved.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set, got zero value")
		}
		if retrieved.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set, got zero value")
		}
		if !retrieved.CreatedAt.Equal(retrieved.UpdatedAt) {
			t.Error("Expected CreatedAt and UpdatedAt to be equal on creation")
		}
	})

	t.Run("Create with empty ID", func(t *testing.T) {
		bookmark := models.NewBookmark("", "https://example.com", "Title", "", nil)
		err := store.Create(ctx, bookmark)
		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})

	t.Run("Create with empty URL", func(t *testing.T) {
		bookmark := models.NewBookmark("test-id", "", "Title", "", nil)
		err := store.Create(ctx, bookmark)
		if err == nil {
			t.Error("Expected error for empty URL, got nil")
		}
	})

	t.Run("Create with empty title", func(t *testing.T) {
		bookmark := models.NewBookmark("test-id", "https://example.com", "", "", nil)
		err := store.Create(ctx, bookmark)
		if err == nil {
			t.Error("Expected error for empty title, got nil")
		}
	})

	t.Run("Create with empty description (should succeed)", func(t *testing.T) {
		bookmark := models.NewBookmark("test-id-2", "https://example.com", "Title", "", nil)
		err := store.Create(ctx, bookmark)
		if err != nil {
			t.Errorf("Expected no error for empty description, got %v", err)
		}
	})

	t.Run("Create with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		bookmark := models.NewBookmark("test-id-3", "https://example.com", "Title", "", nil)
		err := store.Create(ctx, bookmark)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_GetByID(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Setup
	bookmark := models.NewBookmark("test-id", "https://example.com", "Title", "Description", []string{"tag1"})
	store.Create(ctx, bookmark)

	t.Run("GetByID existing bookmark", func(t *testing.T) {
		retrieved, err := store.GetByID(ctx, "test-id")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if retrieved.ID != "test-id" {
			t.Errorf("Expected ID test-id, got %s", retrieved.ID)
		}
	})

	t.Run("GetByID non-existent bookmark", func(t *testing.T) {
		retrieved, err := store.GetByID(ctx, "non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent bookmark, got nil")
		}
		if retrieved != nil {
			t.Error("Expected nil bookmark, got non-nil")
		}
		if err.Error() != "bookmark not found" {
			t.Errorf("Expected 'bookmark not found' error, got %v", err)
		}
	})

	t.Run("GetByID returns copy (not original)", func(t *testing.T) {
		retrieved, _ := store.GetByID(ctx, "test-id")
		originalTitle := retrieved.Title

		// Try to modify the returned bookmark
		retrieved.Title = "Modified Title"

		// Get again and verify original is unchanged
		retrievedAgain, _ := store.GetByID(ctx, "test-id")
		if retrievedAgain.Title != originalTitle {
			t.Error("Expected returned bookmark to be a copy, modifications affected original")
		}
	})

	t.Run("GetByID with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := store.GetByID(ctx, "test-id")
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_GetAll(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	t.Run("GetAll with empty store", func(t *testing.T) {
		bookmarks, total, err := store.GetAll(ctx, 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 0 {
			t.Errorf("Expected 0 bookmarks, got %d", len(bookmarks))
		}
		if total != 0 {
			t.Errorf("Expected total 0, got %d", total)
		}
	})

	t.Run("GetAll with pagination", func(t *testing.T) {
		// Create 25 bookmarks
		for i := 0; i < 25; i++ {
			id := strconv.Itoa(i)
			bookmark := models.NewBookmark(id, "https://example"+id+".com", "Title "+id, "", nil)
			store.Create(ctx, bookmark)
		}

		// Get page 1 with page size 10
		bookmarks, total, err := store.GetAll(ctx, 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 10 {
			t.Errorf("Expected 10 bookmarks on page 1, got %d", len(bookmarks))
		}
		if total != 25 {
			t.Errorf("Expected total 25, got %d", total)
		}

		// Get page 2 with page size 10
		bookmarks, total, err = store.GetAll(ctx, 2, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 10 {
			t.Errorf("Expected 10 bookmarks on page 2, got %d", len(bookmarks))
		}

		// Get page 3 with page size 10
		bookmarks, total, err = store.GetAll(ctx, 3, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 5 {
			t.Errorf("Expected 5 bookmarks on page 3, got %d", len(bookmarks))
		}
	})

	t.Run("GetAll with page 0 (should default to 1)", func(t *testing.T) {
		bookmarks, _, err := store.GetAll(ctx, 0, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 25 {
			t.Errorf("Expected 25 bookmarks with page 0 (defaulting to 1), got %d", len(bookmarks))
		}
	})

	t.Run("GetAll with pageSize 0 (should default to 10)", func(t *testing.T) {
		bookmarks, _, err := store.GetAll(ctx, 1, 0)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 10 {
			t.Errorf("Expected 10 bookmarks with pageSize 0 (defaulting to 10), got %d", len(bookmarks))
		}
	})

	t.Run("GetAll with page beyond available data", func(t *testing.T) {
		bookmarks, total, err := store.GetAll(ctx, 100, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(bookmarks) != 0 {
			t.Errorf("Expected 0 bookmarks for page beyond data, got %d", len(bookmarks))
		}
		if total != 25 {
			t.Errorf("Expected total 25, got %d", total)
		}
	})

	t.Run("GetAll with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := store.GetAll(ctx, 1, 10)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_Update(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Setup
	bookmark := models.NewBookmark("test-id", "https://example.com", "Original Title", "Original Description", []string{"tag1"})
	store.Create(ctx, bookmark)

	t.Run("Update existing bookmark", func(t *testing.T) {
		updated := models.NewBookmark("test-id", "https://updated.com", "Updated Title", "Updated Description", []string{"tag2", "tag3"})

		err := store.Update(ctx, "test-id", updated)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		retrieved, _ := store.GetByID(ctx, "test-id")
		if retrieved.URL != "https://updated.com" {
			t.Errorf("Expected URL https://updated.com, got %s", retrieved.URL)
		}
		if retrieved.Title != "Updated Title" {
			t.Errorf("Expected Title Updated Title, got %s", retrieved.Title)
		}
		if retrieved.Description != "Updated Description" {
			t.Errorf("Expected Description Updated Description, got %s", retrieved.Description)
		}
		if len(retrieved.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(retrieved.Tags))
		}
	})

	t.Run("Update modifies timestamp", func(t *testing.T) {
		// Get original timestamps
		original, _ := store.GetByID(ctx, "test-id")
		originalUpdatedAt := original.UpdatedAt

		// Small delay to ensure timestamp change
		time.Sleep(10 * time.Millisecond)

		updated := models.NewBookmark("test-id", "https://updated.com", "Updated Title", "", nil)
		err := store.Update(ctx, "test-id", updated)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		retrieved, _ := store.GetByID(ctx, "test-id")
		if !retrieved.UpdatedAt.After(originalUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated after modification")
		}
	})

	t.Run("Update non-existent bookmark", func(t *testing.T) {
		updated := models.NewBookmark("non-existent", "https://example.com", "Title", "", nil)
		err := store.Update(ctx, "non-existent", updated)
		if err == nil {
			t.Error("Expected error for non-existent bookmark, got nil")
		}
		if err.Error() != "bookmark not found" {
			t.Errorf("Expected 'bookmark not found' error, got %v", err)
		}
	})

	t.Run("Update with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		updated := models.NewBookmark("test-id", "https://example.com", "Title", "", nil)
		err := store.Update(ctx, "test-id", updated)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Setup
	bookmark := models.NewBookmark("test-id", "https://example.com", "Title", "", nil)
	store.Create(ctx, bookmark)

	t.Run("Delete existing bookmark", func(t *testing.T) {
		err := store.Delete(ctx, "test-id")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		_, err = store.GetByID(ctx, "test-id")
		if err == nil {
			t.Error("Expected error after delete, got nil")
		}
	})

	t.Run("Delete non-existent bookmark", func(t *testing.T) {
		err := store.Delete(ctx, "non-existent")
		if err == nil {
			t.Error("Expected error for non-existent bookmark, got nil")
		}
		if err.Error() != "bookmark not found" {
			t.Errorf("Expected 'bookmark not found' error, got %v", err)
		}
	})

	t.Run("Delete with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := store.Delete(ctx, "test-id")
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_Search(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Setup test data
	bookmarks := []struct {
		id          string
		url         string
		title       string
		description string
		tags        []string
	}{
		{"1", "https://golang.org", "Go Programming Language", "Official Go documentation", []string{"programming", "docs"}},
		{"2", "https://github.com", "GitHub", "Git repository hosting", []string{"git", "programming"}},
		{"3", "https://stackoverflow.com", "Stack Overflow", "Programming Q&A", []string{"programming", "qanda"}},
		{"4", "https://example.com", "Example Site", "This is an example description", []string{"example"}},
	}

	for _, b := range bookmarks {
		store.Create(ctx, models.NewBookmark(b.id, b.url, b.title, b.description, b.tags))
	}

	t.Run("Search by title", func(t *testing.T) {
		results, total, err := store.Search(ctx, "Go", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 1 {
			t.Errorf("Expected 1 result, got %d", total)
		}
		if results[0].Title != "Go Programming Language" {
			t.Errorf("Expected 'Go Programming Language', got %s", results[0].Title)
		}
	})

	t.Run("Search by description", func(t *testing.T) {
		results, total, err := store.Search(ctx, "documentation", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 1 {
			t.Errorf("Expected 1 result, got %d", total)
		}
	})

	t.Run("Search by URL", func(t *testing.T) {
		results, total, err := store.Search(ctx, "github", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 1 {
			t.Errorf("Expected 1 result, got %d", total)
		}
	})

	t.Run("Search case-insensitive", func(t *testing.T) {
		results, total, err := store.Search(ctx, "GO", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 1 {
			t.Errorf("Expected 1 result for case-insensitive search, got %d", total)
		}

		resultsLower, _ := store.Search(ctx, "go", 1, 10)
		if len(results) != len(resultsLower) {
			t.Error("Case-insensitive search should return same results")
		}
	})

	t.Run("Search with multiple matches", func(t *testing.T) {
		results, total, err := store.Search(ctx, "programming", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 3 {
			t.Errorf("Expected 3 results, got %d", total)
		}
	})

	t.Run("Search with no matches", func(t *testing.T) {
		results, total, err := store.Search(ctx, "nonexistent", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 0 {
			t.Errorf("Expected 0 results, got %d", total)
		}
		if len(results) != 0 {
			t.Errorf("Expected empty results, got %d", len(results))
		}
	})

	t.Run("Search with pagination", func(t *testing.T) {
		results, total, err := store.Search(ctx, "programming", 1, 2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 3 {
			t.Errorf("Expected total 3, got %d", total)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results on page 1, got %d", len(results))
		}

		results2, _, _ := store.Search(ctx, "programming", 2, 2)
		if len(results2) != 1 {
			t.Errorf("Expected 1 result on page 2, got %d", len(results2))
		}
	})

	t.Run("Search with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := store.Search(ctx, "programming", 1, 10)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_GetByTag(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Setup test data
	bookmarks := []struct {
		id    string
		title string
		tags  []string
	}{
		{"1", "Go Docs", []string{"go", "programming", "docs"}},
		{"2", "GitHub", []string{"git", "programming"}},
		{"3", "React Docs", []string{"react", "programming", "frontend"}},
		{"4", "Node.js", []string{"nodejs", "javascript", "backend"}},
	}

	for _, b := range bookmarks {
		store.Create(ctx, models.NewBookmark(b.id, "https://example.com/"+b.id, b.title, "", b.tags))
	}

	t.Run("GetByTag with existing tag", func(t *testing.T) {
		results, total, err := store.GetByTag(ctx, "programming", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 3 {
			t.Errorf("Expected 3 results, got %d", total)
		}
	})

	t.Run("GetByTag case-insensitive", func(t *testing.T) {
		resultsUpper, totalUpper, _ := store.GetByTag(ctx, "GO", 1, 10)
		resultsLower, totalLower, _ := store.GetByTag(ctx, "go", 1, 10)

		if totalUpper != totalLower {
			t.Error("Case-insensitive tag search should return same results")
		}
		if len(resultsUpper) != len(resultsLower) {
			t.Error("Case-insensitive tag search should return same number of results")
		}
	})

	t.Run("GetByTag with non-existent tag", func(t *testing.T) {
		results, total, err := store.GetByTag(ctx, "nonexistent", 1, 10)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 0 {
			t.Errorf("Expected 0 results, got %d", total)
		}
		if len(results) != 0 {
			t.Errorf("Expected empty results, got %d", len(results))
		}
	})

	t.Run("GetByTag with pagination", func(t *testing.T) {
		results, total, err := store.GetByTag(ctx, "programming", 1, 2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if total != 3 {
			t.Errorf("Expected total 3, got %d", total)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results on page 1, got %d", len(results))
		}
	})

	t.Run("GetByTag with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := store.GetByTag(ctx, "programming", 1, 10)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_GetAllTags(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	t.Run("GetAllTags with empty store", func(t *testing.T) {
		tags, err := store.GetAllTags(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("GetAllTags with bookmarks", func(t *testing.T) {
		// Setup
		store.Create(ctx, models.NewBookmark("1", "https://example.com", "Title 1", "", []string{"go", "programming"}))
		store.Create(ctx, models.NewBookmark("2", "https://example.com", "Title 2", "", []string{"go", "docs"}))
		store.Create(ctx, models.NewBookmark("3", "https://example.com", "Title 3", "", []string{"programming", "web"}))

		tags, err := store.GetAllTags(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedCounts := map[string]int{
			"go":         2,
			"programming": 2,
			"docs":       1,
			"web":        1,
		}

		for tag, count := range expectedCounts {
			if tags[tag] != count {
				t.Errorf("Expected tag %s to have count %d, got %d", tag, count, tags[tag])
			}
		}
	})

	t.Run("GetAllTags with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := store.GetAllTags(ctx)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
	})
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	t.Run("Concurrent Create and Get", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100

		// Concurrent creates
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				id := strconv.Itoa(idx)
				bookmark := models.NewBookmark(id, "https://example.com/"+id, "Title "+id, "", nil)
				err := store.Create(ctx, bookmark)
				if err != nil {
					t.Errorf("Create failed: %v", err)
				}
			}(i)
		}

		wg.Wait()

		// Concurrent gets
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				id := strconv.Itoa(idx)
				_, err := store.GetByID(ctx, id)
				if err != nil {
					t.Errorf("GetByID failed: %v", err)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Concurrent Update", func(t *testing.T) {
		// Setup
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i)
			bookmark := models.NewBookmark(id, "https://example.com/"+id, "Title "+id, "", nil)
			store.Create(ctx, bookmark)
		}

		var wg sync.WaitGroup
		numGoroutines := 50

		// Concurrent updates
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				id := strconv.Itoa(idx % 10)
				updated := models.NewBookmark(id, "https://updated.com/"+id, "Updated "+id, "", nil)
				err := store.Update(ctx, id, updated)
				if err != nil {
					// Expected for some updates as IDs are shared
					return
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Concurrent Delete", func(t *testing.T) {
		// Setup
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i + 100)
			bookmark := models.NewBookmark(id, "https://example.com/"+id, "Title "+id, "", nil)
			store.Create(ctx, bookmark)
		}

		var wg sync.WaitGroup

		// Concurrent deletes
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				id := strconv.Itoa(idx + 100)
				_ = store.Delete(ctx, id)
			}(i)
		}

		wg.Wait()

		// Verify all deleted
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i + 100)
			_, err := store.GetByID(ctx, id)
			if err == nil {
				t.Errorf("Bookmark %s should have been deleted", id)
			}
		}
	})

	t.Run("Concurrent Read and Write", func(t *testing.T) {
		// Setup
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i + 200)
			bookmark := models.NewBookmark(id, "https://example.com/"+id, "Title "+id, "", []string{"tag1"})
			store.Create(ctx, bookmark)
		}

		var wg sync.WaitGroup
		duration := time.After(100 * time.Millisecond)

		// Concurrent reads
		go func() {
			for {
				select {
				case <-duration:
					return
				default:
					store.GetAll(ctx, 1, 10)
					store.Search(ctx, "title", 1, 10)
					store.GetByTag(ctx, "tag1", 1, 10)
				}
			}
		}()

		// Concurrent writes
		go func() {
			i := 300
			for {
				select {
				case <-duration:
					return
				default:
					id := strconv.Itoa(i)
					bookmark := models.NewBookmark(id, "https://example.com/"+id, "Title "+id, "", nil)
					store.Create(ctx, bookmark)
					i++
				}
			}
		}()

		wg.Wait()
	})
}
