# go-bookmark-manager - Requirements Analysis

## Stakeholders

| Stakeholder | Role | Interest |
|-------------|------|----------|
| End Users | Bookmark managers | Save, organize, and retrieve bookmarks efficiently |
| Developer | Go backend developer | Clean architecture, testable code, proper error handling |
| Developer | React frontend developer | Simple state management, type-safe API integration |
| DevOps | Deployment engineer | Health checks, logging, CORS configuration |
| QA Engineer | Testing | Unit tests for store and handler layers |

## User Stories (with acceptance criteria)

### US-1: Create Bookmark
**As a** user, **I want to** create a new bookmark with URL, title, description, and tags, **so that** I can save web resources for later reference.

**Acceptance Criteria:**
- Given I am on the bookmark list page, when I click "Add Bookmark", then a modal form appears with fields: URL, Title, Description (optional), Tags (comma-separated).
- Given I enter a valid URL and title, when I submit the form, then the bookmark is created and appears in the list.
- Given I enter an invalid URL format, when I submit the form, then an error message is displayed and the bookmark is not created.
- Given I enter tags separated by commas, when I submit the form, then each tag is stored individually and displayed as chips on the bookmark card.
- Given I leave the description field empty, when I submit the form, then the bookmark is created without a description.

### US-2: View Bookmarks List
**As a** user, **I want to** view all my bookmarks in a list, **so that** I can browse my saved resources.

**Acceptance Criteria:**
- Given bookmarks exist, when I view the bookmark list page, then I see bookmark cards with title, truncated URL, tag chips, and description preview.
- Given there are many bookmarks, when I view the list, then pagination controls (prev/next + page indicator) are displayed.
- Given I am on page 1, when I click "Next", then I navigate to page 2 with the next set of bookmarks.
- Given I am on the last page, when I click "Next", then the button is disabled or has no effect.

### US-3: Search Bookmarks
**As a** user, **I want to** search bookmarks by title, description, or URL, **so that** I can quickly find specific bookmarks.

**Acceptance Criteria:**
- Given I type a search term in the search bar, when I press Enter or click search, then only matching bookmarks are displayed.
- Given a bookmark's title, description, or URL contains the search term, when I search, then that bookmark appears in results.
- Given no bookmarks match the search term, when I search, then an empty state message is displayed.
- Given I clear the search bar, when I press Enter, then all bookmarks are displayed again.

### US-4: Filter by Tag
**As a** user, **I want to** filter bookmarks by tag using a dropdown or sidebar, **so that** I can view bookmarks in a specific category.

**Acceptance Criteria:**
- Given tags exist in the system, when I view the tag sidebar or filter dropdown, then all unique tags are displayed with their bookmark counts.
- Given I click on a tag in the sidebar, then only bookmarks with that tag are displayed.
- Given I am viewing filtered bookmarks, when I click "Clear Filter" or select "All Tags", then all bookmarks are displayed again.
- Given I have both a search term and a tag filter active, when I search, then results match both criteria.

### US-5: View Single Bookmark
**As a** user, **I want to** view details of a single bookmark, **so that** I can see the full description and all tags.

**Acceptance Criteria:**
- Given I click on a bookmark card, when the bookmark details view opens, then I see the full title, complete URL, full description, and all tags.
- Given I am on the bookmark detail view, when I click "Back", then I return to the bookmark list.

### US-6: Edit Bookmark
**As a** user, **I want to** edit an existing bookmark, **so that** I can update incorrect or outdated information.

**Acceptance Criteria:**
- Given I am viewing a bookmark, when I click "Edit", then a modal opens pre-filled with the bookmark's current data.
- Given I modify any field and save, when the form is submitted, then the bookmark is updated and the list refreshes.
- Given I try to change the URL to an invalid format, when I save, then an error is displayed and the bookmark is not updated.
- Given I modify tags, when I save, then the new tags replace the old tags on the bookmark.

### US-7: Delete Bookmark
**As a** user, **I want to** delete a bookmark, **so that** I can remove resources I no longer need.

**Acceptance Criteria:**
- Given I am viewing a bookmark, when I click "Delete", then a confirmation dialog appears.
- Given I confirm the deletion, when I click "Yes", then the bookmark is removed from the list.
- Given I cancel the deletion, when I click "No", then the bookmark remains in the list.
- Given I delete a bookmark that is paginated, when the bookmark is removed, then the pagination adjusts if necessary.

### US-8: View All Tags
**As a** user, **I want to** see all available tags with counts, **so that** I can understand how my bookmarks are categorized.

**Acceptance Criteria:**
- Given I load the application, when I view the tag sidebar, then all unique tags are displayed with the count of bookmarks per tag.
- Given tags are added or removed from bookmarks, when I refresh the tag list, then the counts are updated accordingly.
- Given no bookmarks exist, when I view the tag list, then it shows "No tags" or is empty.

### US-9: System Health Check
**As a** developer/ops engineer, **I want to** check the system health, **so that** I can verify the backend is running properly.

**Acceptance Criteria:**
- Given the backend is running, when I call GET /api/health, then I receive { status: "ok", uptime: <seconds> }.
- Given the backend is not running, when I call GET /api/health, then the request fails or times out.

## Functional Requirements

### Backend (Go)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-BE-001 | Implement RESTful API endpoints for bookmark CRUD operations | Must |
| FR-BE-002 | Implement GET /api/bookmarks with pagination (page, limit params) | Must |
| FR-BE-003 | Implement search functionality on title, description, and URL fields | Must |
| FR-BE-004 | Implement tag filtering on bookmark list endpoint | Must |
| FR-BE-005 | Implement GET /api/tags endpoint returning unique tags with counts | Must |
| FR-BE-006 | Implement GET /api/health endpoint returning status and uptime | Must |
| FR-BE-007 | Validate URL format on bookmark creation and update | Must |
| FR-BE-008 | Implement CORS middleware allowing localhost:5173 | Must |
| FR-BE-009 | Implement request logging middleware (method, path, duration) | Must |
| FR-BE-010 | Implement panic recovery middleware | Must |
| FR-BE-011 | Use SQLite database with go-sqlite3 driver | Must |
| FR-BE-012 | Implement bookmark_tags junction table for many-to-many relationship | Must |
| FR-BE-013 | Write unit tests for store layer | Must |
| FR-BE-014 | Write unit tests for handler layer | Must |
| FR-BE-015 | Follow project structure: cmd/server/main.go, internal/handler/, internal/model/, internal/store/, internal/middleware/ | Should |

### Frontend (React + TypeScript)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-FE-001 | Implement BookmarkList component displaying bookmarks as cards | Must |
| FR-FE-002 | Display bookmark title, truncated URL, tag chips, and description preview | Must |
| FR-FE-003 | Implement search bar for filtering bookmarks | Must |
| FR-FE-004 | Implement tag filter dropdown/sidebar | Must |
| FR-FE-005 | Implement AddBookmarkForm modal component | Must |
| FR-FE-006 | Implement EditBookmarkModal component with pre-filled data | Must |
| FR-FE-007 | Implement TagSidebar component showing tags with counts | Must |
| FR-FE-008 | Implement Pagination component with prev/next and page indicator | Must |
| FR-FE-009 | Use fetch API for all HTTP requests (no axios) | Must |
| FR-FE-010 | Use useState/useEffect for state management (no Redux) | Must |
| FR-FE-011 | Use CSS modules for styling (no UI library) | Must |
| FR-FE-012 | Configure Vite dev server to proxy /api to Go backend on port 8080 | Must |
| FR-FE-013 | Use Vite + React + TypeScript stack | Must |

### API Specification

| Endpoint | Method | Request Body | Response | Priority |
|----------|--------|--------------|----------|----------|
| /api/bookmarks | POST | { url, title, description?, tags?: string[] } | Bookmark object | Must |
| /api/bookmarks | GET | ?search=&tag=&page=&limit= | Paginated bookmark list | Must |
| /api/bookmarks/:id | GET | - | Single bookmark object | Must |
| /api/bookmarks/:id | PUT | { url?, title?, description?, tags??: string[] } | Updated bookmark object | Must |
| /api/bookmarks/:id | DELETE | - | 204 No Content or { success: true } | Must |
| /api/tags | GET | - | { tags: [{ name, count }] } | Must |
| /api/health | GET | - | { status: "ok", uptime: number } | Must |

### Database Schema

**bookmarks table:**
- id: INTEGER PRIMARY KEY AUTOINCREMENT
- url: TEXT NOT NULL
- title: TEXT NOT NULL
- description: TEXT (nullable)
- created_at: DATETIME
- updated_at: DATETIME

**bookmark_tags table:**
- bookmark_id: INTEGER (FOREIGN KEY references bookmarks.id)
- tag: TEXT
- PRIMARY KEY (bookmark_id, tag)

## Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-001 | Backend API response time should be under 500ms for standard operations | Must |
| NFR-002 | Frontend should render initial page load within 2 seconds | Must |
| NFR-003 | All API endpoints must return appropriate HTTP status codes | Must |
| NFR-004 | Application must handle concurrent requests without data corruption | Must |
| NFR-005 | URL validation must follow RFC 3986 standards | Must |
| NFR-006 | Code must have unit test coverage for store and handler layers | Must |
| NFR-007 | Frontend must be type-safe using TypeScript | Must |
| NFR-008 | Request logging must include method, path, and duration | Must |
| NFR-009 | Panic recovery must prevent server crashes from unhandled panics | Must |
| NFR-010 | CORS must be configured to allow only localhost:5173 in development | Must |
| NFR-011 | Default pagination limit should be 20 items per page | Should |
| NFR-012 | Database connection should be properly managed and closed | Should |
| NFR-013 | Frontend error handling should display user-friendly messages | Should |
| NFR-014 | URL truncation in UI should show meaningful portion (e.g., domain + path start) | Should |
| NFR-015 | Tag input should trim whitespace and normalize case (optional) | Could |

## Edge Cases

| ID | Edge Case | Handling Strategy |
|----|-----------|-------------------|
| EC-001 | Empty search query | Return all bookmarks (treat as no filter) |
| EC-002 | Search with no results | Display empty state message |
| EC-003 | Tag filter with no matching bookmarks | Display empty state |
| EC-004 | Invalid page number (negative, zero, non-numeric) | Return page 1 or 400 error |
| EC-005 | Invalid limit value (negative, zero, extremely large) | Clamp to default (20) or max (100) |
| EC-006 | Duplicate URL submission | Return 409 Conflict or allow (business decision) |
| EC-007 | Bookmark with no tags | Display without tag chips, include in "untagged" filter if needed |
| EC-008 | Bookmark with many tags (overflow UI) | Show "N more" indicator or truncate display |
| EC-009 | Tag name with special characters | Escape properly in UI and database |
| EC-010 | Very long URL | Truncate in UI, store full URL in database |
| EC-011 | Very long title/description | Truncate in list view, show full in detail view |
| EC-012 | Deleting non-existent bookmark | Return 404 Not Found |
| EC-013 | Updating non-existent bookmark | Return 404 Not Found |
| EC-014 | GET single bookmark for non-existent ID | Return 404 Not Found |
| EC-015 | Malformed JSON in request body | Return 400 Bad Request |
| EC-016 | Missing required fields (url, title) | Return 400 Bad Request with validation errors |
| EC-017 | Invalid URL format | Return 400 Bad Request |
| EC-018 | Database connection failure | Return 500 Internal Server Error |
| EC-019 | Concurrent updates to same bookmark | Last write wins (SQLite handles locking) |
| EC-020 | Empty tag string in comma-separated input | Ignore empty tags after split |
| EC-021 | Tag filter dropdown with many tags | Add scroll or search within dropdown |
| EC-022 | Network failure during API call | Show error message, allow retry |
| EC-023 | Backend unavailable when frontend loads | Show connection error, allow retry |
| EC-024 | Total bookmarks less than page size | Show only available items, disable next button |
| EC-025 | Bookmark count changes during pagination | Current page remains valid, counts may be stale |

## Assumptions

| ID | Assumption |
|----|------------|
| AS-001 | Single-user/local development scenario is the primary use case (no authentication/authorization required) |
| AS-002 | SQLite is sufficient for the expected data volume (hundreds to low thousands of bookmarks) |
| AS-003 | Tags are case-sensitive (e.g., "Work" and "work" are different tags) |
| AS-004 | URL validation uses a standard regex pattern that accepts most valid HTTP/HTTPS URLs |
| AS-005 | The Go server runs on port 8080 and frontend dev server on port 5173 |
| AS-006 | Tags input accepts comma-separated values and splits on comma character |
| AS-007 | Pagination is 1-indexed (page=1 is the first page) |
| AS-008 | The "uptime" in health check is seconds since server start |
| AS-009 | Bookmark deletion is permanent (no soft delete or trash) |
| AS-010 | No image previews or metadata extraction from bookmarked URLs |
| AS-011 | Bookmarks are sorted by created_at descending (newest first) by default |
| AS-012 | The application runs locally; no cloud deployment or containerization is required initially |
| AS-013 | go-sqlite3 is acceptable despite CGO requirement (no pure Go SQLite alternative needed) |
| AS-014 | CSS modules are sufficient for styling needs (no responsive framework required) |
| AS-015 | No real-time updates or WebSocket connections needed |
| AS-016 | Tag counts are calculated at query time, not cached |
| AS-017 | No import/export functionality for bookmarks is required |
| AS-018 | No bookmark validation (e.g., checking if URL is still accessible) is required |
| AS-019 | No duplicate detection beyond optional URL uniqueness constraint |
| AS-020 | The search is case-insensitive for better UX |

## Open Questions

| ID | Question | Impact |
|----|----------|--------|
| OQ-001 | Should URLs be unique (prevent duplicate bookmarks with same URL)? | Database constraint, API validation |
| OQ-002 | What is the maximum length for title, description, and tag fields? | Database schema design |
| OQ-003 | Should tag names have a maximum length or character restrictions? | Input validation, UI design |
| OQ-004 | Should the search support exact match vs. partial match? | Search algorithm complexity |
| OQ-005 | Should search support multiple keywords (AND/OR logic)? | Search query parsing |
| OQ-006 | What is the maximum number of tags allowed per bookmark? | Input validation |
| OQ-007 | Should tags be normalized (trimmed, lowercased) on save? | Data consistency |
| OQ-008 | Should the API support bulk operations (create/update/delete multiple)? | Additional endpoints |
| OQ-009 | Should deleted bookmarks be soft-deleted (recoverable) or hard-deleted? | Database schema, API design |
| OQ-010 | What HTTP status code for successful DELETE (200 vs 204)? | API specification |
| OQ-011 | Should the health endpoint include database connectivity check? | Health check depth |
| OQ-012 | Should request logging write to file or console only? | Logging configuration |
| OQ-013 | What is the expected max concurrent users? | Performance requirements |
| OQ-014 | Should the frontend support keyboard navigation/accessibility features? | UX requirements |
| OQ-015 | Should there be a maximum page limit (e.g., prevent page=999999)? | API security |
