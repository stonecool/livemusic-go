# Project Structure

## Directory Layout

livemusic-go/
├── cmd/                    # Application entry points
│   ├── cli/               # Command line tools
│   ├── server/            # Server entry point
│   └── test/              # Test programs
├── conf/                  # Configuration files
├── docs/                  # Documentation
│   └── swagger/           # API documentation
└── internal/              # Private application code
    ├── account/           # Account management
    ├── cache/             # Cache related
    ├── chrome/            # Chrome browser interaction
    ├── client/            # Client implementations
    ├── config/            # Configuration handling
    ├── database/          # Database operations
    ├── router/            # Router handlers
    │   └── api/           # API routes
    ├── scheduler/         # Task scheduling
    └── task/              # Task definitions and handlers

## Project Overview
This is a web automation project built with Go, featuring Chrome browser control, task scheduling, and API services.

## Key Components

### cmd Layer
- **cli**: Command line interface tools
- **server**: Main application server
- **test**: Testing utilities

### internal Layer
- **account**: User account management
- **cache**: Data caching mechanisms
- **chrome**: Chrome browser automation
- **client**: Client implementations for various services
- **config**: Application configuration
- **database**: Database operations and models
- **router**: HTTP routing and API handlers
- **scheduler**: Task scheduling system
- **task**: Task definitions and execution logic

## Core Components Relationships
1. **Task Management Flow**:
   - `scheduler/` → Manages task scheduling
   - `task/` → Defines task types and execution logic
   - `chrome/` → Handles browser automation

2. **Data Flow**:
   - `database/` → Manages data persistence
   - `cache/` → Handles data caching
   - `client/` → External service communication

3. **API Flow**:
   - `router/api/` → API endpoints
   - `account/` → User authentication
   - `database/` → Data operations

## Design Patterns
1. **Layered Architecture**:
   - Data layer (`internal/database`)
   - Business logic layer (`internal/task`, `internal/scheduler`)
   - API layer (`internal/router`)

2. **Resource Management**:
   - Task scheduling
   - Browser automation
   - Cache management

## Common Development Workflows
1. Adding new features:
   - Add database models in `internal/database`
   - Implement business logic in respective packages
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