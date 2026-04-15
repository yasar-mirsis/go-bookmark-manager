package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	// Get allowed origins from environment or default to localhost:3000
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
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

	// Add logging middleware to all routes
	loggingHandler := withLogging(mux)

	// Add CORS middleware to all routes
	corsHandler := withCORS(loggingHandler, allowedOrigins)

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

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
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

// withLogging adds request/response logging middleware
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf(
			"%s %s %d %v",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
