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
	@echo "  make migrate # Run database migrations"
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
	go run main.go

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

# Run database migrations
mig-up:
	go run cmd/migration/main.go migrate up

# Rollback database migrations
mig-down:
	go run cmd/migration/main.go migrate down

# DANGEROUS - Reset migration
mig-reset:
	go run cmd/migration/main.go migrate reset

# Insert seed data (for development)
seed:
	go run database/seeder/main.go

# Complete project setup
setup: deps swagger mig-up seed
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