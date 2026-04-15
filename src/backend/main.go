package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go-bookmark-manager/src/backend/handler"
	"go-bookmark-manager/src/backend/store"
)

func main() {
	// Initialize store
	st := store.NewMemoryStore()

	// Initialize handler
	h := handler.NewHandler(st)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get allowed origins from environment or default to all
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	// Setup routes
	mux := http.NewServeMux()

	// Bookmark routes
	mux.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(r.Header.Get("Origin"), allowedOrigins))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

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
		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(r.Header.Get("Origin"), allowedOrigins))
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

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
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(r.Header.Get("Origin"), allowedOrigins))
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		h.GetTags(w, r)
	})

	// Health check route
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		h.HealthCheck(w, r)
	})

	// Add CORS middleware to all routes
	corsHandler := withCORS(mux, allowedOrigins)

	// Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsHandler,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down")
}

// getAllowedOrigin returns the allowed origin or default
func getAllowedOrigin(requestedOrigin, allowedOrigins string) string {
	if allowedOrigins == "*" {
		return "*"
	}

	// Check if the requested origin is in the allowed list
	origins := strings.Split(allowedOrigins, ",")
	for _, origin := range origins {
		if strings.TrimSpace(origin) == requestedOrigin {
			return requestedOrigin
		}
	}

	return "*"
}

// withCORS adds CORS headers to all responses
func withCORS(next http.Handler, allowedOrigins string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(r.Header.Get("Origin"), allowedOrigins))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
