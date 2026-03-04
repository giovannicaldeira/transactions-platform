# transactions-platform

A Go-based API platform built with [Gin](https://github.com/gin-gonic/gin) web framework and [Cobra](https://github.com/spf13/cobra) CLI.

## Project Structure

```
transactions-platform/
├── cmd/              # CLI commands (Cobra)
│   ├── root.go      # Root command
│   └── api.go       # API server command
├── internal/        # Private application code
│   ├── app/         # Application logic
│   │   └── api.go   # API server build and run functions
│   └── handlers/    # HTTP handlers
│       └── health.go # Health check endpoint
├── main.go          # Application entry point
├── go.mod           # Go module definition
└── Makefile         # Build automation
```

## Features

- ✅ Gin web framework for high-performance HTTP handling
- ✅ Cobra CLI for command management
- ✅ Graceful shutdown handling
- ✅ Health check endpoint
- ✅ Structured logging with middleware
- ✅ Panic recovery middleware
- ✅ Configurable port via environment variable
- ✅ Swagger/OpenAPI documentation

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/transactions-platform/transactions-platform.git
cd transactions-platform

# Install dependencies
go mod download
```

### Running the Application

```bash
# Run directly
make run
# or
go run . api

# Build and run binary
make build
./bin/transactions-platform api

# With custom port
PORT=3000 ./bin/transactions-platform api
```

### Development

```bash
# Run tests
make test

# Generate Swagger documentation
make swagger

# Update vendor directory
make vendor

# Clean build artifacts
make clean

# Show available commands
make help
```

## API Documentation

The API documentation is automatically generated using Swagger/OpenAPI and is available at:

**Swagger UI:** `http://localhost:8080/swagger/index.html`

The Swagger documentation provides:
- Interactive API documentation
- Request/response examples
- Schema definitions
- Try-it-out functionality

### Regenerating Swagger Docs

After modifying API endpoints or adding new handlers with Swagger annotations:

```bash
make swagger
```

Or manually:

```bash
swag init --parseDependency --parseInternal
```

## API Endpoints

| Method | Path       | Description              |
|--------|-----------|--------------------------|
| GET    | /health   | Health check             |
| GET    | /swagger/*any | Swagger documentation |

### Example Request

```bash
curl http://localhost:8080/health
```

### Example Response

```json
{
  "status": "healthy",
  "timestamp": "2026-03-03T14:33:50.098537Z"
}
```

## Environment Variables

| Variable  | Description              | Default |
|-----------|--------------------------|---------|
| PORT      | HTTP server port         | 8080    |
| GIN_MODE  | Gin mode (debug/release) | release |

## Commands

```bash
# Show available commands
./bin/transactions-platform --help

# Run API server
./bin/transactions-platform api
```
