# go-bookmark-manager - System Architecture

## System Overview

The go-bookmark-manager is a full-stack web application designed for users to save, organize, and retrieve web bookmarks. The system follows a clean architecture pattern with a Go backend and React frontend, providing RESTful API communication.

**Key Characteristics:**
- RESTful API architecture
- Clean architecture with separation of concerns
- In-memory storage with interface abstraction for testability
- Single-page application (SPA) frontend
- CORS-enabled for cross-origin requests
- Health check endpoints for monitoring

## Components

### Backend Components (src/backend/)

#### 1. Bookmark Store (`src/backend/store/store.go`)
- **Responsibility:** Data persistence and retrieval operations
- **Interfaces:**
  ```go
  type Store interface {
      Create(ctx context.Context, bookmark *Bookmark) error
      GetByID(ctx context.Context, id string) (*Bookmark, error)
      GetAll(ctx context.Context, page, pageSize int) ([]Bookmark, int, error)
      Update(ctx context.Context, id string, bookmark *Bookmark) error
      Delete(ctx context.Context, id string) error
      Search(ctx context.Context, query string, page, pageSize int) ([]Bookmark, int, error)
      GetByTag(ctx context.Context, tag string, page, pageSize int) ([]Bookmark, int, error)
      GetAllTags(ctx context.Context) (map[string]int, error)
  }
  ```
- **Testing Strategy:** Mock implementations for unit tests

#### 2. Bookmark Handler (`src/backend/handler/handler.go`)
- **Responsibility:** HTTP request handling, validation, and response formatting
- **Interfaces:**
  ```go
  type Handler struct {
      store store.Store
  }
  ```
- **Methods:**
  - `CreateBookmark(w http.ResponseWriter, r *http.Request)`
  - `GetBookmark(w http.ResponseWriter, r *http.Request)`
  - `GetBookmarks(w http.ResponseWriter, r *http.Request)`
  - `UpdateBookmark(w http.ResponseWriter, r *http.Request)`
  - `DeleteBookmark(w http.ResponseWriter, r *http.Request)`
  - `SearchBookmarks(w http.ResponseWriter, r *http.Request)`
  - `GetBookmarksByTag(w http.ResponseWriter, r *http.Request)`
  - `GetTags(w http.ResponseWriter, r *http.Request)`
  - `HealthCheck(w http.ResponseWriter, r *http.Request)`

#### 3. Route Configuration (`src/backend/main.go`)
- **Responsibility:** Application entry point, routing setup, CORS middleware, and server initialization
- **Components:**
  - HTTP router configuration
  - CORS middleware with configurable origins
  - Health check endpoint registration
  - Graceful shutdown handling

### Frontend Components (src/frontend/)

#### 1. Bookmark List Component (`src/frontend/components/BookmarkList.tsx`)
- **Responsibility:** Display paginated list of bookmarks
- **State:** Bookmarks array, pagination state, loading state
- **Interfaces:**
  - Props: `bookmarks: Bookmark[]`, `onBookmarkClick: (id: string) => void`
  - Events: `onPageChange: (page: number) => void`

#### 2. Bookmark Card Component (`src/frontend/components/BookmarkCard.tsx`)
- **Responsibility:** Display individual bookmark with title, URL preview, tags, and description
- **State:** None (presentational component)
- **Interfaces:**
  - Props: `bookmark: Bookmark`, `onEdit: () => void`, `onDelete: () => void`, `onClick: () => void`

#### 3. Bookmark Form Component (`src/frontend/components/BookmarkForm.tsx`)
- **Responsibility:** Create and edit bookmark forms with validation
- **State:** Form fields (url, title, description, tags), validation errors, submission state
- **Interfaces:**
  - Props: `initialValues?: Bookmark`, `onSubmit: (data: BookmarkFormData) => void`, `mode: 'create' | 'edit'`
  - Validation: URL format, required title field

#### 4. Search Component (`src/frontend/components/SearchBar.tsx`)
- **Responsibility:** Search input with debounced query handling
- **State:** Search query, debounce timer
- **Interfaces:**
  - Props: `onSearch: (query: string) => void`
  - Events: `onClear: () => void`

#### 5. Tag Filter Component (`src/frontend/components/TagFilter.tsx`)
- **Responsibility:** Display available tags and allow filtering by tag
- **State:** Selected tag, available tags with counts
- **Interfaces:**
  - Props: `tags: TagCount[]`, `selectedTag: string | null`, `onTagSelect: (tag: string | null) => void`

#### 6. Pagination Component (`src/frontend/components/Pagination.tsx`)
- **Responsibility:** Page navigation controls
- **State:** Current page, total pages
- **Interfaces:**
  - Props: `currentPage: number`, `totalPages: number`, `onPageChange: (page: number) => void`

#### 7. Modal Component (`src/frontend/components/Modal.tsx`)
- **Responsibility:** Reusable modal wrapper for forms
- **State:** Open/closed state
- **Interfaces:**
  - Props: `isOpen: boolean`, `onClose: () => void`, `title: string`, `children: ReactNode`

#### 8. API Service (`src/frontend/services/api.ts`)
- **Responsibility:** Type-safe API client for all backend endpoints
- **Interfaces:**
  ```typescript
  interface BookmarkAPI {
      createBookmark(data: BookmarkFormData): Promise<Bookmark>
      getBookmark(id: string): Promise<Bookmark>
      getBookmarks(page?: number, pageSize?: number): Promise<PaginatedResponse<Bookmark>>
      updateBookmark(id: string, data: BookmarkFormData): Promise<Bookmark>
      deleteBookmark(id: string): Promise<void>
      searchBookmarks(query: string, page?: number): Promise<PaginatedResponse<Bookmark>>
      getBookmarksByTag(tag: string, page?: number): Promise<PaginatedResponse<Bookmark>>
      getTags(): Promise<TagCount[]>
  }
  ```

## Data Model

### Entities

#### Bookmark
```typescript
interface Bookmark {
    id: string;           // UUID v4
    url: string;          // Valid URL format
    title: string;        // Required, non-empty
    description?: string; // Optional
    tags: string[];       // Array of tag strings
    createdAt: string;    // ISO 8601 timestamp
    updatedAt: string;    // ISO 8601 timestamp
}
```

#### TagCount
```typescript
interface TagCount {
    tag: string;
    count: number;
}
```

#### PaginatedResponse
```typescript
interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
}
```

#### BookmarkFormData
```typescript
interface BookmarkFormData {
    url: string;
    title: string;
    description?: string;
    tags: string; // Comma-separated string for UI input
}
```

### Relationships

- **Bookmark ↔ Tag:** Many-to-many relationship (implemented as array of strings in Bookmark)
  - Each bookmark can have multiple tags
  - Each tag can be associated with multiple bookmarks
  - Tags are stored as a string array within the Bookmark entity

### Entity Relationships Diagram

```
┌─────────────────┐
│    Bookmark     │
├─────────────────┤
│ id (PK)         │
│ url             │
│ title           │
│ description     │
│ tags[]          │───┐
│ createdAt       │   │
│ updatedAt       │   │
└─────────────────┘   │
                      │
                      ▼
                ┌──────────┐
                │   Tag    │
                ├──────────┤
                │ name     │
                │ count    │
                └──────────┘
```

## API Contracts

### Base URL
- Development: `http://localhost:8080/api`
- Production: `https://api.example.com/api`

### Endpoints

#### 1. Health Check
- **Method:** `GET`
- **Path:** `/api/health`
- **Request:** None
- **Response:**
  ```json
  {
    "status": "ok",
    "timestamp": "2026-04-15T10:30:00Z"
  }
  ```
- **Status Codes:**
  - `200 OK`: Service is healthy

#### 2. Create Bookmark
- **Method:** `POST`
- **Path:** `/api/bookmarks`
- **Request Body:**
  ```json
  {
    "url": "https://example.com",
    "title": "Example Site",
    "description": "An example website",
    "tags": "example,demo"
  }
  ```
- **Response:**
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "url": "https://example.com",
    "title": "Example Site",
    "description": "An example website",
    "tags": ["example", "demo"],
    "createdAt": "2026-04-15T10:30:00Z",
    "updatedAt": "2026-04-15T10:30:00Z"
  }
  ```
- **Status Codes:**
  - `201 Created`: Bookmark successfully created
  - `400 Bad Request`: Invalid URL format or missing required fields
  - `500 Internal Server Error`: Database error

#### 3. Get All Bookmarks
- **Method:** `GET`
- **Path:** `/api/bookmarks`
- **Query Parameters:**
  - `page` (optional, default: 1): Page number
  - `pageSize` (optional, default: 10): Items per page
- **Response:**
  ```json
  {
    "data": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "url": "https://example.com",
        "title": "Example Site",
        "description": "An example website",
        "tags": ["example", "demo"],
        "createdAt": "2026-04-15T10:30:00Z",
        "updatedAt": "2026-04-15T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10,
    "totalPages": 1
  }
  ```
- **Status Codes:**
  - `200 OK`: Successful retrieval

#### 4. Get Single Bookmark
- **Method:** `GET`
- **Path:** `/api/bookmarks/{id}`
- **Path Parameters:**
  - `id` (required): Bookmark UUID
- **Response:**
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "url": "https://example.com",
    "title": "Example Site",
    "description": "An example website",
    "tags": ["example", "demo"],
    "createdAt": "2026-04-15T10:30:00Z",
    "updatedAt": "2026-04-15T10:30:00Z"
  }
  ```
- **Status Codes:**
  - `200 OK`: Bookmark found
  - `404 Not Found`: Bookmark does not exist

#### 5. Update Bookmark
- **Method:** `PUT`
- **Path:** `/api/bookmarks/{id}`
- **Path Parameters:**
  - `id` (required): Bookmark UUID
- **Request Body:**
  ```json
  {
    "url": "https://updated-example.com",
    "title": "Updated Site",
    "description": "Updated description",
    "tags": "updated,example"
  }
  ```
- **Response:**
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "url": "https://updated-example.com",
    "title": "Updated Site",
    "description": "Updated description",
    "tags": ["updated", "example"],
    "createdAt": "2026-04-15T10:30:00Z",
    "updatedAt": "2026-04-15T11:00:00Z"
  }
  ```
- **Status Codes:**
  - `200 OK`: Bookmark successfully updated
  - `400 Bad Request`: Invalid data format
  - `404 Not Found`: Bookmark does not exist

#### 6. Delete Bookmark
- **Method:** `DELETE`
- **Path:** `/api/bookmarks/{id}`
- **Path Parameters:**
  - `id` (required): Bookmark UUID
- **Response:** Empty body
- **Status Codes:**
  - `204 No Content`: Bookmark successfully deleted
  - `404 Not Found`: Bookmark does not exist

#### 7. Search Bookmarks
- **Method:** `GET`
- **Path:** `/api/bookmarks/search`
- **Query Parameters:**
  - `q` (required): Search query string
  - `page` (optional, default: 1): Page number
  - `pageSize` (optional, default: 10): Items per page
- **Response:** Same format as Get All Bookmarks
- **Status Codes:**
  - `200 OK`: Search completed
  - `400 Bad Request`: Missing search query

#### 8. Get Bookmarks by Tag
- **Method:** `GET`
- **Path:** `/api/bookmarks/tags/{tag}`
- **Path Parameters:**
  - `tag` (required): Tag name
- **Query Parameters:**
  - `page` (optional, default: 1): Page number
  - `pageSize` (optional, default: 10): Items per page
- **Response:** Same format as Get All Bookmarks
- **Status Codes:**
  - `200 OK`: Tags found
  - `404 Not Found`: Tag does not exist (no bookmarks with this tag)

#### 9. Get All Tags
- **Method:** `GET`
- **Path:** `/api/tags`
- **Request:** None
- **Response:**
  ```json
  [
    {"tag": "example", "count": 5},
    {"tag": "demo", "count": 3},
    {"tag": "tutorial", "count": 2}
  ]
  ```
- **Status Codes:**
  - `200 OK`: Tags retrieved successfully

## Technology Stack

### Backend

| Technology | Version | Justification |
|------------|---------|---------------|
| Go | 1.21+ | Built-in HTTP server, excellent performance, strong typing, simple deployment with single binary |
| net/http (standard library) | - | No external dependencies required, built-in routing, production-ready, fits the simple REST API requirements |
| testing (standard library) | - | Built-in testing framework, zero configuration needed |

**Why Go?**
- High performance with low resource consumption
- Excellent concurrency model for future scaling
- Strong static typing reduces runtime errors
- Single binary deployment simplifies DevOps
- Large ecosystem and community support

**Why standard library HTTP?**
- Zero external dependencies = smaller attack surface
- Faster startup time and lower memory footprint
- Sufficient for REST API requirements
- Easier testing without complex mocking frameworks
- Aligns with the simplicity requirements of the project

### Frontend

| Technology | Version | Justification |
|------------|---------|---------------|
| React | 18.x | Industry standard, component-based architecture, large ecosystem, excellent developer experience |
| TypeScript | 5.x | Type safety, better IDE support, catches errors at compile time, improves code maintainability |
| Vite | 5.x | Fast development server, instant HMR, optimized builds, modern tooling |
| React Router | 6.x | Client-side routing, type-safe route definitions, lazy loading support |

**Why React?**
- Component-based architecture matches the UI requirements
- Large talent pool and community support
- Excellent ecosystem for form handling and state management
- Strong TypeScript integration

**Why TypeScript?**
- Type safety for API integration
- Better refactoring support
- Self-documenting code
- Catches errors before runtime

### Development Tools

| Tool | Purpose |
|------|---------|
| go test | Backend unit testing |
| Jest | Frontend unit testing |
| React Testing Library | Frontend component testing |
| go vet | Go static analysis |
| ESLint | TypeScript linting |
| Prettier | Code formatting |

## Data Flow

### Create Bookmark Flow

```
User → Frontend
  │
  ├──► 1. Opens "Add Bookmark" modal
  │     └─► BookmarkForm component renders
  │
  ├──► 2. Enters URL, title, description, tags
  │     └─► Form validation checks URL format, required fields
  │
  ├──► 3. Submits form
  │     └─► API.createBookmark() called
  │           └─► POST /api/bookmarks
  │                 │
  │                 ▼
  │            Backend Handler
  │                 │
  │                 ├──► Parse JSON body
  │                 ├──► Validate URL format
  │                 ├──► Generate UUID
  │                 ├──► Parse tags (comma-separated → array)
  │                 ├──► Set timestamps
  │                 └─► Store.Create()
  │                       │
  │                       ▼
  │                 In-Memory Store
  │                       │
  │                       ▼
  │                 201 Created response
  │
  └─► 4. Frontend receives response
        └─► Close modal
        └─► Refresh bookmark list
        └─► Show success notification
```

### Search Bookmarks Flow

```
User → Frontend
  │
  ├──► 1. Types search term in SearchBar
  │     └─► Debounced input (300ms)
  │
  ├──► 2. Presses Enter or search icon clicked
  │     └─► API.searchBookmarks(query) called
  │           └─► GET /api/bookmarks/search?q={query}
  │                 │
  │                 ▼
  │            Backend Handler
  │                 │
  │                 ├──► Parse query parameter
  │                 └─► Store.Search(query, page, pageSize)
  │                       │
  │                       ▼
  │                 In-Memory Store
  │                       │
  │                       ├──► Iterate bookmarks
  │                       ├──► Check title, description, URL, tags
  │                       └─► Return matches with pagination
  │
  └─► 3. Frontend receives results
        └─► Update BookmarkList with filtered results
        └─► Show empty state if no matches
```

### Filter by Tag Flow

```
User → Frontend
  │
  ├──► 1. Clicks tag in TagFilter sidebar
  │     └─► API.getBookmarksByTag(tag) called
  │           └─► GET /api/bookmarks/tags/{tag}
  │                 │
  │                 ▼
  │            Backend Handler
  │                 │
  │                 └─► Store.GetByTag(tag, page, pageSize)
  │                       │
  │                       ▼
  │                 In-Memory Store
  │                       │
  │                       ├──► Filter bookmarks by tags array
  │                       └─► Return matches with pagination
  │
  └─► 3. Frontend receives results
        └─► Update BookmarkList with filtered results
        └─► Highlight selected tag
```

### Combined Search + Filter Flow

```
User → Frontend
  │
  ├──► 1. Has active search term AND tag filter
  │
  ├──► 2. Search or filter changes
  │     └─► API.getBookmarks() with both params
  │           └─► GET /api/bookmarks?q={query}&tag={tag}
  │                 │
  │                 ▼
  │            Backend Handler
  │                 │
  │                 └─► Store.Search() then filter by tag
  │                       (or vice versa, depending on implementation)
  │
  └─► 3. Frontend receives combined results
```

## Security Considerations

### Input Validation
- **URL Validation:** Strict URL format validation using Go's `url.Parse()` and protocol whitelist (http/https only)
- **Title Validation:** Required field, non-empty, trimmed whitespace
- **Tag Sanitization:** Trim whitespace, lowercase, reject special characters
- **Description Length:** Maximum 1000 characters to prevent abuse

### CORS Configuration
- **Allowed Origins:** Configurable via environment variable (`CORS_ALLOWED_ORIGINS`)
- **Allowed Methods:** GET, POST, PUT, DELETE, OPTIONS
- **Allowed Headers:** Content-Type, Authorization
- **Credentials:** Disabled (no cookies/sessions required)

### HTTP Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security` (production only)

### Rate Limiting (Future)
- Implement rate limiting on write operations (POST, PUT, DELETE)
- Consider token bucket algorithm for fine-grained control

### Data Sanitization
- All user input is sanitized before storage
- HTML entities encoded in frontend display to prevent XSS
- No user input directly rendered without sanitization

### Error Handling
- No internal error details exposed to client
- Generic error messages for 500 responses
- Structured logging for debugging (server-side only)

## Scalability Notes

### Current Architecture Limitations
- **In-Memory Storage:** Data lost on restart, not suitable for production
- **Single Instance:** No horizontal scaling support with current store

### Future Scalability Enhancements

#### 1. Database Migration
- **Short-term:** SQLite file-based database for persistence
- **Long-term:** PostgreSQL for production with:
  - Connection pooling
  - Indexes on title, description, tags columns
  - Full-text search capabilities

#### 2. Caching Layer
- **Redis Integration:**
  - Cache frequently accessed bookmarks
  - Cache tag list for filter dropdown
  - Cache search results with TTL

#### 3. Horizontal Scaling
- **Stateless Design:** Backend is already stateless (except store)
- **Load Balancer:** Ready for multiple instances once database is externalized
- **Session Management:** Not required (no user authentication in current scope)

#### 4. Performance Optimizations
- **Indexing:** Database indexes on frequently queried fields
- **Pagination:** Already implemented to limit result sets
- **Query Optimization:** Full-text search for better search performance
- **CDN:** Static assets (frontend) served via CDN

#### 5. Monitoring & Observability
- **Health Checks:** Already implemented at `/api/health`
- **Metrics:** Prometheus metrics for request counts, latency
- **Logging:** Structured JSON logging for log aggregation
- **Tracing:** Distributed tracing with OpenTelemetry (future)

### Resource Considerations
- **Memory:** In-memory store requires sufficient RAM for bookmark dataset
- **CPU:** Go's goroutines handle concurrent requests efficiently
- **Network:** Stateless API allows easy horizontal scaling

### Backward Compatibility
- API versioning not required for initial release
- Design endpoints for future versioning if needed (`/api/v1/bookmarks`)
