run: build start

# ===========================================
# Migrations
# ===========================================
mig-build:
	@echo ">> Building migration..."
	@go build -o bin/migration ./cmd/migration

mig-up: mig-build
	@echo ">> executing migration..."
	@./bin/migration migrate up
	@echo ">> finished executing migration..."

mig-down-all: mig-build
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

deploy:
	@echo ">> Deploying changes..."
	@docker-compose stop && docker-compose up -d --build
	@echo ">> Running schema migration"
	mig-up
	@echo ">> Changes are Deployed"