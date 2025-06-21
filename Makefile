.PHONY: help build run test clean swagger deps migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  swagger  - Generate Swagger documentation"
	@echo "  deps     - Install dependencies"
	@echo "  migrate  - Run database migrations"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o bin/amartha main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Generate Swagger documentation
swagger:
	@echo "Installing swag if not already installed..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	swag init -g main.go -o docs
	@echo "Swagger documentation generated successfully!"

# Development mode (run with hot reload)
dev:
	@echo "Starting development server..."
	@echo "Note: Install air for hot reload: go install github.com/cosmtrek/air@latest"
	air

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

mig-status: mig-build
	@echo ">> Migration Status"
	@./bin/migration migrate status

mig-create: mig-build
	@echo ">> Create Migration"
	@./bin/migration migrate create $(name) go

setup: deps swagger mig-up seed
	@echo "Project setup completed!"

deploy:
	@echo ">> Deploying changes..."
	@docker-compose stop && docker-compose up -d --build
	@echo ">> Running schema migration"
	mig-up
	@echo ">> Changes are Deployed"