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


---

This project is managed by the SDLC Pipeline. Implementation tasks are tracked as GitHub/GitLab issues.
Each issue is solved by an autonomous agent on its own branch with a pull request.