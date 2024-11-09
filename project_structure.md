# Project Structure

## Directory Layout

livemusic-go/
├── cmd/                    # Application entry points
│   ├── cli/               # Command line tools
│   ├── server/            # Server entry point
│   └── test/              # Test programs
├── conf/                   # Configuration files
├── docs/                   # Documentation
│   └── swagger/           # API documentation
└── internal/              # Private application code
    ├── cache/             # Cache related
    ├── chrome/            # Chrome browser interaction
    ├── config/            # Configuration handling
    ├── crawl/             # Crawler implementation
    ├── http/              # HTTP client
    ├── instance/          # Instance management
    ├── model/             # Data models
    └── router/            # Router handlers
        └── api/           # API routes

## Project Overview
This appears to be a web crawler project that interacts with Chrome browser instances. The project is structured following Go best practices with clear separation of concerns.

## Key Components

### cmd Layer
- **cli**: Command line interface tools for manual operations
- **server**: Main application server that handles HTTP requests
- **test**: Testing utilities and integration tests

### internal Layer
- **cache**: Handles data caching mechanisms
- **chrome**: Manages Chrome browser automation and control
- **config**: Handles application configuration and environment settings
- **crawl**: Contains core crawling logic and strategies
- **http**: HTTP client for external communications
- **instance**: Manages Chrome browser instances lifecycle
- **model**: Data models and database interactions
- **router**: HTTP routing and API endpoint definitions
    - **api**: RESTful API route handlers

## Core Components Relationships
1. **Chrome Instance Management Flow**:
   - `model/chromeinstance.go` → Defines data structure
   - `chrome/` → Implements browser control logic
   - `instance/` → Manages instance lifecycle

2. **Crawling Flow**:
   - `crawl/` → Implements crawling logic
   - `chrome/` → Controls browser
   - `cache/` → Caches crawled data

3. **API Flow**:
   - `router/api/` → Defines endpoints
   - `model/` → Handles data persistence
   - `instance/` → Manages resources

## Design Patterns
1. **Layered Architecture**:
   - Models layer (`internal/model`)
   - Business logic layer (`internal/chrome`, `internal/crawl`)
   - API layer (`internal/router`)

2. **Resource Management**:
   - Instance pooling
   - Connection management
   - Cache strategies

## Common Development Workflows
1. Adding new features:
   - Add model definitions in `internal/model`
   - Implement business logic in respective domains
   - Add API endpoints in `internal/router/api`

2. Testing:
   - Unit tests alongside source files
   - Integration tests in `cmd/test`
   - API tests via Swagger

## Configuration
- Configuration files stored in `conf/`
- Runtime configs managed by `internal/config`

---
*Note: This document should be updated when making significant architectural changes.*