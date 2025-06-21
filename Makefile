.PHONY: help deps build dev test clean swagger mig-up mig-down mig-reset seed setup test-unit test-coverage test-verbose test-file test-services test-helpers generate-mocks generate-swagger migrate seed fmt lint

# Default target
help:
	@echo "Available commands:"
	@echo "  make help      # Show available commands"
	@echo "  make deps      # Install dependencies"
	@echo "  make build     # Build the application"
	@echo "  make dev       # Run the application"
	@echo "  make test # Run unit tests only"
	@echo "  make test-coverage # Run tests with coverage report"
	@echo "  make mocks # Generate mocks using mockery"
	@echo "  make swagger # Generate Swagger documentation"
	@echo "  make mig-up # Run database migrations"
	@echo "  make mig-down # Rollback database migrations"
	@echo "  make mig-reset # Reset database migrations"
	@echo "  make seed # Run database seeders"
	@echo "  make clean     # Clean build artifacts"
	@echo "  make setup     # Complete project setup"
	@echo "  make fmt       # Format code"
	@echo "  make lint      # Lint code"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	go build -o bin/amartha-loan-engine main.go

# Run the application
dev:
	@echo "Starting development server..."Add commentMore actions
	@echo "Note: Install air for hot reload: go install github.com/cosmtrek/air@latest"
	air

# Run unit tests only
test:
	go test -v ./...

# Run tests with coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...

# Generate mocks
mocks:
	mockery --dir=repositories --output=mock --outpkg=mock --all

# Generate Swagger documentation
swagger:
	swag init -g main.go -o docs

seed:
	@echo ">> Seeding data..."
	@go build -o bin/migration ./cmd/migration
	@./bin/migration seed
	@echo ">> finished seeding data..."

mig-build:
	@echo "Installing goose if not already installed..."
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo ">> Building migration..."
	@go build -o bin/migration ./cmd/migration

mig-up: mig-build
	@echo ">> executing migration..."
	@./bin/migration migrate up
	@echo ">> finished executing migration..."

mig-reset: mig-build
	@echo ">> resetting migration..."
	@./bin/migration migrate reset
	@echo ">> finished resetting migration..."

mig-down: mig-build
	@echo ">> Rolling back migration 1 version..."
	@./bin/migration migrate down
	@echo ">> finished rolling bank migration 1 version..."


# Complete project setup
setup: deps mig-up
	@echo "Project setup completed!"

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out