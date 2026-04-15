# Backend Documentation

## Overview

The backend is a Go-based REST API server that provides bookmark management functionality. It uses the standard library's `net/http` package for HTTP handling and an in-memory store for data persistence.

## Project Structure

```
src/backend/
├── handler/          # HTTP request handlers
│   ├── handler.go    # Main handler implementation
│   └── handler_test.go  # Unit tests for handlers
├── store/            # Data persistence layer
│   ├── store.go      # Store interface definition
│   ├── memory_store.go  # In-memory store implementation
│   └── memory_store_test.go  # Unit tests for store
├── models/           # Data models
│   └── bookmark.go   # Bookmark model definition
├── main.go           # Application entry point and routing
├── integration_test.go  # Integration tests for full HTTP stack
└── README.md         # This file
```

## API Endpoints

### Bookmarks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/bookmarks` | List all bookmarks (paginated) |
| POST | `/bookmarks` | Create a new bookmark |
| GET | `/bookmarks?id={id}` | Get a specific bookmark |
| PUT | `/bookmarks?id={id}` | Update a bookmark |
| DELETE | `/bookmarks?id={id}` | Delete a bookmark |

### Search

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/bookmarks/search?q={query}` | Search bookmarks by title, description, or URL |

### Tags

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/bookmarks/tag?tag={tag}` | Get bookmarks filtered by tag |
| GET | `/tags` | List all tags with counts |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

## Query Parameters

### Pagination
- `page` (default: 1) - Page number
- `pageSize` (default: 10) - Number of items per page

### Search
- `q` (required) - Search query string

### Tag Filtering
- `tag` (required) - Tag name to filter by

## Running Tests

### Run All Tests

```bash
go test -v ./...
```

### Run Unit Tests Only

```bash
go test -v ./store ./handler
```

### Run Integration Tests Only

```bash
go test -v -run TestHealth ./...
go test -v -run TestCreateBookmark ./...
go test -v -run TestFullCRUD ./...
```

### Run Tests with Coverage

```bash
go test -v -cover ./...
```

### Run Tests with Race Detection

```bash
go test -v -race ./...
```

## Integration Tests

The integration tests (`integration_test.go`) provide end-to-end testing of the HTTP server using `httptest.StartServer()`. They cover:

### Health Check Verification
- Verifies `/health` endpoint returns `{"status": "ok"}` with status 200

### Full CRUD Operations
- Create bookmark with array tags
- Create bookmark with string tags (comma-separated)
- Get bookmark by ID
- Update bookmark
- Delete bookmark and verify removal

### Pagination
- Paginated results retrieval
- Page/pageSize defaults (page=1, pageSize=10)
- Edge cases: page 0, pageSize 0, page beyond total

### Search Functionality
- Search by title
- Search by description
- Search by URL
- Empty results handling

### Tag Filtering
- Filter bookmarks by tag
- List all tags with counts
- Non-existent tag handling

### Error Scenarios
- Invalid JSON body (400)
- Missing required fields (400)
- Invalid URL format (400)
- Non-existent resources (404)
- Method not allowed (405)

### Concurrent Access
- Thread-safety verification with concurrent bookmark creation

## Environment Variables

- `PORT` - Server port (default: 8080)
- `CORS_ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins (default: http://localhost:3000)

## Building

```bash
go build -o bookmark-manager ./main.go
```

## Running

```bash
./bookmark-manager
```

Or with environment variables:

```bash
PORT=3000 CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:8080" ./bookmark-manager
```
