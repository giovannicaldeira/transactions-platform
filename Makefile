.PHONY: build run run-api test clean vendor swagger help

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

# Show help
help:
	@echo "Available targets:"
	@echo "  make build    - Build the application binary (generates swagger docs)"
	@echo "  make run      - Run the API server"
	@echo "  make test     - Run tests"
	@echo "  make swagger  - Generate Swagger documentation"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make vendor   - Update vendor directory"
	@echo "  make help     - Show this help message"
