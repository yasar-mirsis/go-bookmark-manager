# AGENTS.md — go-bookmark-manager

This file describes the project for AI agents working on implementation issues.

## Project Context

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
- Given tags exist in the system, when I view the tag

[... truncated for brevity ...]

## Architecture

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

#### 2. Bookmark Card Component (`src/frontend/co

[... truncated for brevity ...]

## Working Guidelines

- Read this file and README.md before starting any work
- Follow existing code patterns and conventions
- Write clean, production-quality code with proper error handling
- Create or update tests if a testing setup exists
- Do NOT run git commands — the pipeline handles commits and pushes
- Do NOT ask questions — you are running in an automated pipeline