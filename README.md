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

## How to Test

This section provides a complete walkthrough for testing the API locally using Docker.

### Step 1: Start the Application with Docker

Start both the application and PostgreSQL database:

```bash
# Copy environment file (if not already done)
cp .env.example .env

# Start all services (PostgreSQL + API)
make docker-up
```

Wait for the services to start. You should see logs indicating:
```
✅ PostgreSQL is ready
✅ Running database migrations
✅ Server listening on port 8080
```

**Alternative:** Run only the database in Docker and the app locally:
```bash
# Start only PostgreSQL
make docker-up-deps

# In another terminal, run the app locally
make run
```

### Step 2: Verify the API is Running

Check the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2026-03-05T18:15:30.123456Z"
}
```

### Step 3: Access Swagger Documentation

Open your browser and navigate to:

**🔗 http://localhost:8080/swagger/index.html**

The Swagger UI provides:
- ✅ Interactive API documentation
- ✅ Try-it-out functionality to test endpoints directly
- ✅ Request/response examples
- ✅ Schema definitions for all models

You can test all endpoints directly from the Swagger UI without writing any code!

### Step 4: Test with curl

#### 4.1 Create an Account

```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "document_number": "12345678900"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "document_number": "12345678900"
}
```

💡 **Save the `id` value** - you'll need it for creating transactions!

#### 4.2 Get Account by ID

```bash
# Replace {account_id} with the ID from previous step
curl http://localhost:8080/accounts/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "document_number": "12345678900"
}
```

#### 4.3 Create a Purchase Transaction

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation_type": "NORMAL_PURCHASE",
    "amount": 123.45
  }'
```

**Response:**
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440001",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": "-123.45",
  "event_date": "2026-03-05T18:20:00Z",
  "operation_type": "NORMAL_PURCHASE",
  "created_at": "2026-03-05T18:20:00Z"
}
```

**Note:** The amount is automatically converted to negative for debit operations (purchases, withdrawals).

#### 4.4 Create a Credit Transaction

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation_type": "CREDIT_VOUCHER",
    "amount": 500.00
  }'
```

**Response:**
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440002",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": "500",
  "event_date": "2026-03-05T18:21:00Z",
  "operation_type": "CREDIT_VOUCHER",
  "created_at": "2026-03-05T18:21:00Z"
}
```

**Note:** Credit voucher amounts remain positive.

#### 4.5 Test Error Handling

**Invalid operation type:**
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation_type": "INVALID_TYPE",
    "amount": 100.00
  }'
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid operation type: INVALID_TYPE"
}
```

**Account not found:**
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "non-existent-id",
    "operation_type": "NORMAL_PURCHASE",
    "amount": 100.00
  }'
```

**Response (404 Not Found):**
```json
{
  "error": "account not found"
}
```

**Invalid amount (zero or negative):**
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "operation_type": "NORMAL_PURCHASE",
    "amount": 0
  }'
```

**Response (422 Unprocessable Entity):**
```json
{
  "error": "amount must be positive, got: 0"
}
```

### Step 5: View Application Logs

To see structured logs in real-time:

```bash
# View all container logs
make docker-logs

# Or view only the app logs
docker-compose logs -f app

# For local development
LOG_LEVEL=debug make run
```

**Example log output:**
```json
{"level":"info","msg":"Creating new account","document_number":"12345678900","time":"2026-03-05T18:15:30Z"}
{"level":"info","msg":"Account created successfully","account_id":"550e...","document_number":"12345678900","time":"2026-03-05T18:15:30Z"}
{"level":"info","method":"POST","path":"/accounts","status":201,"latency":15,"ip":"172.18.0.1","msg":"HTTP request","time":"2026-03-05T18:15:30Z"}
```

### Step 6: Cleanup

When done testing:

```bash
# Stop all containers
make docker-down

# Stop and remove volumes (WARNING: deletes database data)
docker-compose down -v
```

### Quick Test Script

Here's a complete test script you can run:

```bash
#!/bin/bash

# Start services
make docker-up

# Wait for services to be ready
sleep 5

# Test health check
echo "Testing health endpoint..."
curl http://localhost:8080/health
echo -e "\n"

# Create account
echo "Creating account..."
ACCOUNT_RESPONSE=$(curl -s -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"document_number": "12345678900"}')
echo $ACCOUNT_RESPONSE | jq .
ACCOUNT_ID=$(echo $ACCOUNT_RESPONSE | jq -r '.id')
echo -e "\n"

# Create purchase transaction
echo "Creating purchase transaction..."
curl -s -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d "{
    \"account_id\": \"$ACCOUNT_ID\",
    \"operation_type\": \"NORMAL_PURCHASE\",
    \"amount\": 123.45
  }" | jq .
echo -e "\n"

# Create credit transaction
echo "Creating credit transaction..."
curl -s -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d "{
    \"account_id\": \"$ACCOUNT_ID\",
    \"operation_type\": \"CREDIT_VOUCHER\",
    \"amount\": 500.00
  }" | jq .
echo -e "\n"

echo "✅ All tests completed!"
```

Save as `test.sh`, make executable with `chmod +x test.sh`, and run with `./test.sh`.

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

| Method | Path              | Description              | Swagger Docs |
|--------|-------------------|--------------------------|--------------|
| GET    | /health           | Health check             | [View](http://localhost:8080/swagger/index.html#/health) |
| POST   | /accounts         | Create a new account     | [View](http://localhost:8080/swagger/index.html#/accounts) |
| GET    | /accounts/:id     | Get account by ID        | [View](http://localhost:8080/swagger/index.html#/accounts) |
| POST   | /transactions     | Create a new transaction | [View](http://localhost:8080/swagger/index.html#/transactions) |
| GET    | /swagger/*any     | Swagger documentation    | - |

### Operation Types

Transactions support the following operation types:

| Operation Type              | Description                    | Amount Behavior |
|-----------------------------|--------------------------------|-----------------|
| `NORMAL_PURCHASE`           | Regular purchase               | Stored as negative (debit) |
| `PURCHASE_WITH_INSTALLMENTS`| Installment purchase           | Stored as negative (debit) |
| `WITHDRAWAL`                | Cash withdrawal                | Stored as negative (debit) |
| `CREDIT_VOUCHER`            | Payment/credit                 | Stored as positive (credit) |

**Note:** All amounts are provided as positive values in the request. The system automatically adjusts the sign based on the operation type.

## Logging

The application uses [zerolog](https://github.com/rs/zerolog) for structured, high-performance logging.

### Log Levels

Logs can be filtered by level (from most to least verbose):
- `debug` - Detailed debugging information
- `info` - General informational messages (default)
- `warn` - Warning messages
- `error` - Error messages
- `fatal` - Fatal errors (application exits)

### Configuration

Set the log level via environment variable:
```bash
LOG_LEVEL=debug ./bin/transactions-platform api
```

Set the environment mode for log formatting:
```bash
# Pretty console output for development
APP_ENV=development ./bin/transactions-platform api

# JSON output for production (default)
APP_ENV=production ./bin/transactions-platform api
```

### Log Output Examples

**Development (pretty console):**
```
2026-03-05T18:00:00Z INF Starting application port=8080 environment=development
2026-03-05T18:00:00Z INF Database connection established
2026-03-05T18:00:00Z INF Server listening port=8080 address=:8080
2026-03-05T18:00:01Z INF HTTP request method=POST path=/accounts status=201 latency=15ms
```

**Production (JSON):**
```json
{"level":"info","msg":"Starting application","port":"8080","environment":"production","time":"2026-03-05T18:00:00Z"}
{"level":"info","msg":"Database connection established","time":"2026-03-05T18:00:00Z"}
{"level":"info","msg":"Server listening","port":"8080","address":":8080","time":"2026-03-05T18:00:00Z"}
{"level":"info","method":"POST","path":"/accounts","status":201,"latency":15,"time":"2026-03-05T18:00:01Z"}
```

### Logged Events

The application logs:
- **Application lifecycle**: startup, shutdown, database connections
- **HTTP requests**: method, path, status code, latency, IP address, user agent
- **Business operations**: account creation, transaction creation with context
- **Errors**: detailed error information with stack traces
- **Recovery**: panic recovery with context

## Environment Variables

| Variable            | Description              | Default                  |
|---------------------|--------------------------|--------------------------|
| PORT                | HTTP server port         | 8080                     |
| GIN_MODE            | Gin mode (debug/release) | release                  |
| LOG_LEVEL           | Logging level            | info                     |
| APP_ENV             | Application environment  | production               |
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
