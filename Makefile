.PHONY: build run run-api test clean vendor swagger docker-build docker-up docker-down docker-logs migrate-up migrate-down migrate-status migrate-create help

# Docker image configuration
DOCKER_IMAGE_NAME ?= transactions-platform
DOCKER_IMAGE_TAG ?= latest

# Database configuration
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= transactions
DB_PASSWORD ?= transactions_password
DB_NAME ?= transactions_platform
DB_SSLMODE ?= disable
DB_DSN ?= "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)"

# Build the application binary
build: swagger
	@echo "Building application..."
	go build -o bin/transactions-platform .

# Run the API server
run:
	@echo "Starting API server..."
	go run . api

# Alias for run
run-api: run

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Update vendor directory
vendor:
	@echo "Updating vendor directory..."
	go mod tidy
	go mod vendor

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@command -v swag >/dev/null 2>&1 || { echo >&2 "swag is not installed. Installing..."; go install github.com/swaggo/swag/cmd/swag@latest; }
	@swag init --parseDependency --parseInternal

# Docker commands
docker-build: 
	@echo "Building Docker image: $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)"
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

docker-up: docker-build
	@echo "Starting Docker containers..."
	docker-compose up -d --wait
	@echo "Database is ready!"
	@echo "Running database migrations..."
	@$(MAKE) migrate-up || echo "⚠️  Migration failed. Run 'make migrate-up' manually after fixing the issue."

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Database migration commands
migrate-up:
	@echo "Running migrations..."
	@command -v goose >/dev/null 2>&1 || { echo >&2 "goose is not installed. Installing..."; go install github.com/pressly/goose/v3/cmd/goose@latest; }
	@goose -dir ./migrations postgres $(DB_DSN) up

migrate-down:
	@echo "Rolling back last migration..."
	@command -v goose >/dev/null 2>&1 || { echo >&2 "goose is not installed. Installing..."; go install github.com/pressly/goose/v3/cmd/goose@latest; }
	goose -dir ./migrations postgres $(DB_DSN) down

migrate-status:
	@echo "Migration status..."
	@command -v goose >/dev/null 2>&1 || { echo >&2 "goose is not installed. Installing..."; go install github.com/pressly/goose/v3/cmd/goose@latest; }
	goose -dir ./migrations postgres $(DB_DSN) status

migrate-create:
	@echo "Creating new migration..."
	@command -v goose >/dev/null 2>&1 || { echo >&2 "goose is not installed. Installing..."; go install github.com/pressly/goose/v3/cmd/goose@latest; }
	@read -p "Enter migration name: " name; \
	goose -dir ./migrations create $$name sql

# Show help
help:
	@echo "Available targets:"
	@echo "  make build              - Build the application binary (generates swagger docs)"
	@echo "  make run                - Run the API server"
	@echo "  make test               - Run tests"
	@echo "  make swagger            - Generate Swagger documentation"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make vendor             - Update vendor directory"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-up          - Start Docker containers and run migrations"
	@echo "  make docker-down        - Stop Docker containers"
	@echo "  make docker-logs        - Show Docker container logs"
	@echo "  make migrate-up         - Run all pending migrations"
	@echo "  make migrate-down       - Rollback the last migration"
	@echo "  make migrate-status     - Show migration status"
	@echo "  make migrate-create     - Create a new migration file"
	@echo "  make help               - Show this help message"
