## Overview

The go-bookmark-manager is a full-stack Go and React application for managing web bookmarks with RESTful API, in-memory storage, pagination, search, tag filtering, and CRUD operations. The implementation follows clean architecture with testable components, including a bookmark store interface, HTTP handler layer, and React components with validation and state management.

## Tasks

### 1. Project Setup and Backend Structure
**Description:** Initialize the Go backend project structure with proper directory layout, go.mod file, and base configuration. Set up the project skeleton including src/backend directory with subdirectories for store, handler, models, and routes. Create the Bookmark model struct with fields: ID (string), URL (string), Title (string), Description (string), Tags ([]string), CreatedAt (time.Time), UpdatedAt (time.Time). Initialize the Go module with proper dependencies.

**Files to create:**
- go.mod
- src/backend/models/bookmark.go
- src/backend/store/store.go
- src/backend/store/memory_store.go
- src/backend/handler/handler.go
- src/backend/main.go

**Files to modify:**
- None

**Complexity:** Low
**Dependencies:** None

### 2. Bookmark Store Implementation
**Description:** Implement the in-memory store that satisfies the Store interface defined in store.go. The store should use a sync.RWMutex for thread-safe access to an in-memory map[string]Bookmark. Implement all store methods: Create (assign UUID, set timestamps), GetByID, GetAll (with pagination support returning bookmarks slice and total count), Update, Delete, Search (case-insensitive search in title, description, URL), GetByTag, GetAllTags (return map of tag->count). Initialize the store in main.go.

**Files to create:**
- src/backend/store/memory_store.go

**Files to modify:**
- src/backend/store/store.go (if interface needs refinement)

**Complexity:** Medium
**Dependencies:** 1

### 3. HTTP Handler Implementation
**Description:** Implement all handler methods in handler.go. Each method should: parse JSON request body, validate input, call store methods, handle errors appropriately, and return JSON responses with proper HTTP status codes. Specific endpoints: POST /bookmarks (create), GET /bookmarks/:id (get single), GET /bookmarks (list with query params page, pageSize), PUT /bookmarks/:id (update), DELETE /bookmarks/:id (delete), GET /bookmarks/search (query param q), GET /bookmarks/tag/:tag (filter by tag), GET /tags (list all tags with counts), GET /health (health check returning {"status": "ok"}). Include validation: URL must be valid format, title is required, tags are optional comma-separated.

**Files to create:**
- src/backend/handler/handler.go

**Files to modify:**
- src/backend/main.go

**Complexity:** Medium
**Dependencies:** 2

### 4. Backend Routing and CORS Configuration
**Description:** Set up HTTP routing in main.go with middleware. Configure CORS middleware to allow requests from frontend origin (http://localhost:3000). Register all handler routes: GET /health, GET /bookmarks, POST /bookmarks, GET /bookmarks/:id, PUT /bookmarks/:id, DELETE /bookmarks/:id, GET /bookmarks/search, GET /bookmarks/tag/:tag, GET /tags. Add graceful shutdown handling for SIGINT/SIGTERM signals. Start HTTP server on port 8080. Include logging middleware for request/response logging.

**Files to create:**
- src/backend/main.go

**Files to modify:**
- src/backend/handler/handler.go

**Complexity:** Medium
**Dependencies:** 2

### 5. React Frontend Setup
**Description:** Initialize React frontend project with Vite or Create React App. Install dependencies: react, react-dom, axios for API calls, react-router-dom for routing. Set up project structure with src/components, src/pages, src/services, src/types. Create TypeScript configuration (tsconfig.json) with strict mode enabled. Configure API base URL environment variable for backend connection.

**Files to create:**
- src/frontend/package.json
- src/frontend/tsconfig.json
- src/frontend/vite.config.ts (or package.json scripts if using CRA)
- src/frontend/index.html
- src/frontend/src/main.tsx
- src/frontend/src/App.tsx
- src/frontend/src/types/index.ts

**Files to modify:**
- None

**Complexity:** Low
**Dependencies:** None

### 6. Frontend Types and API Service Layer
**Description:** Define TypeScript interfaces in src/types/index.ts: Bookmark (matching backend model), BookmarkFormData (for form submission), PaginatedResponse (for API responses with data and total), SearchParams, TagInfo. Create API service in src/services/api.ts using axios: createBookmark, getBookmark, getBookmarks (with pagination params), updateBookmark, deleteBookmark, searchBookmarks, getBookmarksByTag, getTags, healthCheck. Include error handling and response type casting.

**Files to create:**
- src/frontend/src/types/index.ts
- src/frontend/src/services/api.ts

**Files to modify:**
- src/frontend/src/App.tsx

**Complexity:** Low
**Dependencies:** 5

### 7. Bookmark Card Component
**Description:** Create BookmarkCard component in src/components/BookmarkCard.tsx. Display bookmark title (truncated if long), URL preview (first 50 chars), tag chips (clickable for filtering), description preview (truncated to 100 chars). Include action buttons: Edit (opens edit modal), Delete (with confirmation), View (navigates to detail if implemented). Props: bookmark: Bookmark, onEdit: () => void, onDelete: (id: string) => void, onClick: (id: string) => void, onTagClick: (tag: string) => void. Style with CSS classes for card layout, tag chips, button group.

**Files to create:**
- src/frontend/src/components/BookmarkCard.tsx
- src/frontend/src/components/BookmarkCard.module.css

**Files to modify:**
- None

**Complexity:** Low
**Dependencies:** 6

### 8. Bookmark List Component with Pagination
**Description:** Create BookmarkList component in src/components/BookmarkList.tsx. Display array of BookmarkCard components. Implement pagination controls: Previous/Next buttons, page indicator (e.g., "Page 2 of 5"). Handle loading state with skeleton or spinner. Handle empty state with message "No bookmarks found". Props: bookmarks: Bookmark[], currentPage: number, totalPages: number, onPageChange: (page: number) => void, loading: boolean. Include error state display.

**Files to create:**
- src/frontend/src/components/BookmarkList.tsx
- src/frontend/src/components/BookmarkList.module.css

**Files to modify:**
- None

**Complexity:** Low
**Dependencies:** 7

### 9. Bookmark Form Component with Validation
**Description:** Create BookmarkForm component in src/components/BookmarkForm.tsx supporting both create and edit modes (mode: 'create' | 'edit'). Fields: URL (required, validated format), Title (required), Description (optional), Tags (comma-separated, converted to array). Include validation: show error messages below fields, disable submit during validation errors. Props: initialValues?: Bookmark, onSubmit: (data: BookmarkFormData) => void, mode: 'create' | 'edit', onSubmitError?: (error: string) => void. Include loading state during submission.

**Files to create:**
- src/frontend/src/components/BookmarkForm.tsx
- src/frontend/src/components/BookmarkForm.module.css

**Files to modify:**
- None

**Complexity:** Medium
**Dependencies:** 6

### 10. Bookmark Modal and Main Page
**Description:** Create BookmarkModal component for displaying form in modal overlay. Implement BookmarkPage in src/pages/BookmarkPage.tsx as main page with: Add Bookmark button (opens modal), BookmarkList with pagination, SearchBar component, TagFilter component. State management: bookmarks array, current page, total pages, search query, selected tag filter, modal open state, loading state. Integrate API calls for fetching, creating, updating, deleting bookmarks. Handle success/error notifications.

**Files to create:**
- src/frontend/src/components/BookmarkModal.tsx
- src/frontend/src/components/BookmarkModal.module.css
- src/frontend/src/pages/BookmarkPage.tsx
- src/frontend/src/pages/BookmarkPage.module.css

**Files to modify:**
- src/frontend/src/App.tsx (add routing)

**Complexity:** High
**Dependencies:** 7, 8, 9

### 11. Search and Tag Filter Components
**Description:** Create SearchBar component in src/components/SearchBar.tsx with debounced search (300ms delay). Props: onSearch: (query: string) => void, onClear: () => void. Create TagFilter component in src/components/TagFilter.tsx displaying clickable tag chips from getAllTags API. Props: tags: TagInfo[], selectedTag?: string, onTagClick: (tag: string) => void. Include "Clear Filter" button when a tag is selected.

**Files to create:**
- src/frontend/src/components/SearchBar.tsx
- src/frontend/src/components/SearchBar.module.css
- src/frontend/src/components/TagFilter.tsx
- src/frontend/src/components/TagFilter.module.css

**Files to modify:**
- src/frontend/src/pages/BookmarkPage.tsx

**Complexity:** Low
**Dependencies:** 6

### 12. Frontend Styling and Layout
**Description:** Create global CSS in src/index.css with base styles, CSS variables for colors (primary, secondary, background, text). Create layout wrapper component with consistent header (app title, Add Bookmark button), main content area, footer. Style all components with responsive design: card grid layout for bookmarks, modal centered overlay, tag chips with hover effects, pagination buttons. Ensure consistent spacing, typography, and color scheme across all components.

**Files to create:**
- src/frontend/src/index.css
- src/frontend/src/components/Layout.tsx
- src/frontend/src/components/Layout.module.css

**Files to modify:**
- src/frontend/src/App.tsx

**Complexity:** Medium
**Dependencies:** 7, 8, 9, 10, 11

### 13. Backend Unit Tests
**Description:** Write unit tests for store layer in src/backend/store/memory_store_test.go. Test cases: Create (verify ID generation, timestamps), GetByID (existing and non-existent), GetAll (pagination, empty result), Update (modify fields, verify timestamp), Delete (remove from store), Search (case-insensitive, multiple matches), GetByTag (filter by tag), GetAllTags (count per tag). Test thread safety with concurrent access. Write unit tests for handler layer in src/backend/handler/handler_test.go using httptest: valid requests return 200/201, invalid input returns 400, non-existent resource returns 404, store errors return 500.

**Files to create:**
- src/backend/store/memory_store_test.go
- src/backend/handler/handler_test.go

**Files to modify:**
- None

**Complexity:** Medium
**Dependencies:** 2, 3

### 14. Integration Tests and Health Check Verification
**Description:** Create integration tests that start the full HTTP server and test all endpoints end-to-end. Use httptest.StartServer() to start backend, make HTTP requests to verify: full CRUD operations, pagination works correctly, search returns correct results, tag filtering works, health endpoint returns {"status": "ok"}. Test error scenarios: invalid JSON, missing required fields, duplicate requests. Document test commands in README.md.

**Files to create:**
- src/backend/integration_test.go
- src/backend/README.md

**Files to modify:**
- src/backend/main.go

**Complexity:** High
**Dependencies:** 4, 13

### 15. Frontend Component Tests
**Description:** Set up testing library (React Testing Library) in frontend. Write tests for BookmarkCard: renders correctly with all props, button click events fire. Test BookmarkForm: form submission with valid data, validation errors display correctly, disabled state during submit. Test BookmarkList: renders bookmarks, pagination buttons change page. Test SearchBar: debounce behavior, clear button works. Test TagFilter: displays tags, click events fire. Include snapshot tests for key components.

**Files to create:**
- src/frontend/package.json (update with testing dependencies)
- src/frontend/src/components/__tests__/BookmarkCard.test.tsx
- src/frontend/src/components/__tests__/BookmarkForm.test.tsx
- src/frontend/src/components/__tests__/BookmarkList.test.tsx
- src/frontend/src/components/__tests__/SearchBar.test.tsx
- src/frontend/src/components/__tests__/TagFilter.test.tsx

**Files to modify:**
- src/frontend/tsconfig.json

**Complexity:** Medium
**Dependencies:** 7, 8, 9, 11

## File Structure

```
go-bookmark-manager/
├── go.mod
├── src/
│   ├── backend/
│   │   ├── main.go
│   │   ├── models/
│   │   │   └── bookmark.go
│   │   ├── store/
│   │   │   ├── store.go
│   │   │   ├── memory_store.go
│   │   │   └── memory_store_test.go
│   │   ├── handler/
│   │   │   ├── handler.go
│   │   │   └── handler_test.go
│   │   └── integration_test.go
│   └── frontend/
│       ├── package.json
│       ├── tsconfig.json
│       ├── vite.config.ts
│       ├── index.html
│       ├── src/
│       │   ├── main.tsx
│       │   ├── App.tsx
│       │   ├── index.css
│       │   ├── types/
│       │   │   └── index.ts
│       │   ├── services/
│       │   │   └── api.ts
│       │   ├── components/
│       │   │   ├── Layout.tsx
│       │   │   ├── Layout.module.css
│       │   │   ├── BookmarkCard.tsx
│       │   │   ├── BookmarkCard.module.css
│       │   │   ├── BookmarkList.tsx
│       │   │   ├── BookmarkList.module.css
│       │   │   ├── BookmarkForm.tsx
│       │   │   ├── BookmarkForm.module.css
│       │   │   ├── BookmarkModal.tsx
│       │   │   ├── BookmarkModal.module.css
│       │   │   ├── SearchBar.tsx
│       │   │   ├── SearchBar.module.css
│       │   │   ├── TagFilter.tsx
│       │   │   ├── TagFilter.module.css
│       │   │   └── __tests__/
│       │   │       ├── BookmarkCard.test.tsx
│       │   │       ├── BookmarkForm.test.tsx
│       │   │       ├── BookmarkList.test.tsx
│       │   │       ├── SearchBar.test.tsx
│       │   │       └── TagFilter.test.tsx
│       │   └── pages/
│       │       └── BookmarkPage.tsx
│       │       └── BookmarkPage.module.css
```

## Testing Strategy

**Backend Testing:**
- Unit tests for memory_store.go: Test each store method with mock data, verify CRUD operations, pagination logic, search accuracy, tag counting. Use table-driven tests for efficiency. Test concurrent access with sync.Mutex.
- Unit tests for handler.go: Use httptest to mock HTTP requests, validate JSON parsing, test all HTTP status codes (200, 201, 400, 404, 500), verify response structure.
- Integration tests: Start real HTTP server, test full request/response cycles, verify data persistence in memory store across requests.

**Frontend Testing:**
- Unit tests for components: Verify rendering with React Testing Library, test prop handling, event firing for buttons and inputs.
- Form validation tests: Test invalid URL format, required field validation, empty tags handling.
- Integration tests: Test component interaction, API service calls with mock axios.
- Coverage goals: 80%+ for critical components (form, list, API service).

**Manual Testing:**
- Start backend on port 8080, frontend on port 3000.
- Verify CORS allows frontend origin.
- Test all CRUD operations via UI.
- Test pagination with 50+ bookmarks.
- Test search with various queries.
- Test tag filtering and multi-tag bookmarks.
- Verify URL validation with invalid formats.

## Risks

1. **URL Validation**: Go's url.Parse may not catch all invalid URL formats. Risk of accepting malformed URLs. Resolution: Use additional validation library or regex for stricter URL format checking.

2. **Memory Store Limits**: In-memory storage will not persist across restarts and may grow unbounded. Risk of memory exhaustion with large datasets. Resolution: Document this limitation, consider adding cleanup or migration path to database.

3. **Concurrent Access**: Go routines accessing shared map without proper locking could cause race conditions. Resolution: Use sync.RWMutex (already planned) and verify with race detector tests.

4. **CORS Configuration**: Incorrect CORS setup may block frontend requests or create security vulnerabilities. Resolution: Explicitly whitelist frontend origin, avoid wildcard (*) in production.

5. **Frontend Type Safety**: TypeScript interface mismatch with backend JSON could cause runtime errors. Resolution: Generate types from backend schema or use OpenAPI spec, add runtime validation in API layer.
