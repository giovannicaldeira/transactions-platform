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

- Go 1.25 or higher

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

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make coverage-html

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

| Method | Path              | Description              |
|--------|-------------------|--------------------------|
| GET    | /health           | Health check             |
| POST   | /accounts         | Create a new account     |
| GET    | /accounts/:id     | Get account by ID        |
| POST   | /transactions     | Create a new transaction |
| GET    | /swagger/*any     | Swagger documentation    |

### Example Requests

**Health Check:**
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2026-03-03T14:33:50.098537Z"
}
```

**Create Account:**
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"document_number": "12345678900"}'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "document_number": "12345678900",
  "created_at": "2026-03-04T12:00:00Z",
  "updated_at": "2026-03-04T12:00:00Z"
}
```

**Get Account:**
```bash
curl http://localhost:8080/accounts/550e8400-e29b-41d4-a716-446655440000
```

**Create Transaction:**
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation_type": "NORMAL_PURCHASE",
    "amount": 123.45
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": -123.45,
  "event_date": "2026-03-05T12:00:00Z",
  "operation_type": "NORMAL_PURCHASE",
  "created_at": "2026-03-05T12:00:00Z"
}
```

**Operation Types:**
- `NORMAL_PURCHASE` - Regular purchase (negative amount)
- `PURCHASE_WITH_INSTALLMENTS` - Installment purchase (negative amount)
- `WITHDRAWAL` - Cash withdrawal (negative amount)
- `CREDIT_VOUCHER` - Payment/credit (positive amount)

## Environment Variables

| Variable            | Description              | Default                  |
|---------------------|--------------------------|--------------------------|
| PORT                | HTTP server port         | 8080                     |
| GIN_MODE            | Gin mode (debug/release) | release                  |
| DATABASE_HOST       | PostgreSQL host          | postgres                 |
| DATABASE_PORT       | PostgreSQL port          | 5432                     |
| DATABASE_USER       | PostgreSQL username      | transactions             |
| DATABASE_PASSWORD   | PostgreSQL password      | transactions_password    |
| DATABASE_NAME       | PostgreSQL database name | transactions_platform    |
| DATABASE_SSLMODE    | PostgreSQL SSL mode      | disable                  |

## Docker

The application includes Docker support with PostgreSQL database.

### Quick Start with Docker

```bash
# Copy environment file
cp .env.example .env

# Build and start ALL containers (app + PostgreSQL, runs migrations)
make docker-up

# OR start ONLY dependencies for local development
make docker-up-deps

# View logs
make docker-logs

# Stop containers
make docker-down
```

**Note:** Both `make docker-up` and `make docker-up-deps` automatically run database migrations after the containers start. If migrations fail, you can run them manually with `make migrate-up`.

### Development Workflows

**Local development (run app locally, use Docker for DB):**
```bash
# 1. Start only PostgreSQL
make docker-up-deps

# 2. Run app locally
make run

# 3. App connects to PostgreSQL on localhost:5432
```

**Full Docker deployment:**
```bash
# Start everything (app + PostgreSQL)
make docker-up
```

### Manual Docker Commands

```bash
# Build the Docker image
docker-compose build

# Start all services
docker-compose up -d

# Start only PostgreSQL
docker-compose up -d postgres

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Remove volumes (WARNING: deletes database data)
docker-compose down -v
```

### Services

- **app**: API server (http://localhost:8080)
- **postgres**: PostgreSQL database (localhost:5432)

### Database Migrations

The project uses [Goose](https://github.com/pressly/goose) for database migrations.

#### Running Migrations

```bash
# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check migration status
make migrate-status

# Create a new migration
make migrate-create
```

#### Migration Workflow

1. **Create a new migration:**
   ```bash
   make migrate-create
   # Enter migration name when prompted (e.g., "create_users_table")
   ```

2. **Edit the migration file** in `migrations/` directory:
   ```sql
   -- +goose Up
   CREATE TABLE users (
       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
       email VARCHAR(255) UNIQUE NOT NULL,
       created_at TIMESTAMP DEFAULT NOW()
   );

   -- +goose Down
   DROP TABLE users;
   ```

3. **Run the migration:**
   ```bash
   make migrate-up
   ```

#### Environment Variables for Migrations

Set these in your `.env` file or export them:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=transactions
DB_PASSWORD=transactions_password
DB_NAME=transactions_platform
DB_SSLMODE=disable
```

Or use a custom DSN:
```bash
DB_DSN="host=localhost port=5432 user=myuser password=mypass dbname=mydb sslmode=disable" make migrate-up
```

## Commands

```bash
# Show available commands
./bin/transactions-platform --help

# Run API server
./bin/transactions-platform api
```
