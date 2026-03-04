.PHONY: build run run-api test clean vendor swagger docker-build docker-up docker-down docker-logs help

# Docker image configuration
DOCKER_IMAGE_NAME ?= transactions-platform
DOCKER_IMAGE_TAG ?= latest

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
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Show help
help:
	@echo "Available targets:"
	@echo "  make build              - Build the application binary (generates swagger docs)"
	@echo "  make run                - Run the API server"
	@echo "  make test               - Run tests"
	@echo "  make swagger            - Generate Swagger documentation"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make vendor             - Update vendor directory"
	@echo "  make docker-build       - Alias for docker-build-image"
	@echo "  make docker-up          - Start Docker containers (requires pre-built image)"
	@echo "  make docker-down        - Stop Docker containers"
	@echo "  make docker-logs        - Show Docker container logs"
	@echo "  make help               - Show this help message"
