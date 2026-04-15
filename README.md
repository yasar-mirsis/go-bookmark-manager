# go-bookmark-manager

## Overview

The go-bookmark-manager is a full-stack web application designed for users to save, organize, and retrieve web bookmarks. The system follows a clean architecture pattern with a Go backend and React frontend, providing RESTful API communication.

**Key Characteristics:**
- RESTful API architecture
- Clean architecture with separation of concerns
- In-memory storage with interface abstraction for testability
- Single-page application (SPA) frontend
- CORS-enabled for cross-origin requests
- Health check endpoints for monitoring


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


## Running Tests

### Backend Tests

Run all backend tests (unit tests and integration tests):

```bash
cd src/backend
go test -v ./...
```

Run only unit tests:

```bash
cd src/backend
go test -v ./store ./handler
```

Run only integration tests:

```bash
cd src/backend
go test -v -run ^Test.* ./integration_test.go ./main.go ./handler/handler.go ./store/memory_store.go ./models/bookmark.go
```

Or run integration tests as a package:

```bash
cd src/backend
go test -v -run TestHealth ./...
```

Run tests with coverage:

```bash
cd src/backend
go test -v -cover ./...
```

Run tests with race detection:

```bash
cd src/backend
go test -v -race ./...
```

### Integration Test Coverage

The integration tests (`integration_test.go`) cover:

- **Health Check Endpoint**: Verifies `/health` returns `{"status": "ok"}`
- **Full CRUD Operations**: Create, Read, Update, Delete bookmarks end-to-end
- **Pagination**: Tests page/pageSize parameters, edge cases (page 0, pageSize 0, beyond total)
- **Search Functionality**: Search by title, description, and URL
- **Tag Filtering**: Filter bookmarks by tag, list all tags with counts
- **Error Scenarios**: Invalid JSON, missing required fields, invalid URL format, 404s, 405s
- **Concurrent Access**: Tests thread-safety with concurrent bookmark creation

### Test Commands Summary

| Command | Description |
|---------|-------------|
| `go test -v ./...` | Run all tests with verbose output |
| `go test -v -run TestHealth` | Run only health check tests |
| `go test -v -cover ./...` | Run tests with coverage report |
| `go test -v -race ./...` | Run tests with race detection |
| `go test -v -count=1 ./...` | Run tests without caching |

---

This project is managed by the SDLC Pipeline. Implementation tasks are tracked as GitHub/GitLab issues.
Each issue is solved by an autonomous agent on its own branch with a pull request.