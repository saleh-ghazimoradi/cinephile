# Variables for migration path and database URL
MIGRATE_PATH = ./scripts/migrations
DATABASE_URL = ${DB_SOURCE}

# Default target to format the code
format:
	@echo "Applying go fmt to the project"
	go fmt ./...

# Check for errors
vet:
	@echo "Checking for errors with vet"
	go vet ./...

# Docker commands to start and stop the application
dockerup:
	docker compose --env-file app.env up -d

dockerdown:
	docker compose --env-file app.env down

# Migration commands
migrate-create:
	@if [ -z "$$name" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration $$name..."
	migrate create -seq -ext=.sql -dir=$(MIGRATE_PATH) $$name

migrate-up:
	@echo "Applying migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" down

migrate-drop:
	@echo "Dropping all migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" drop

# Run the HTTP server
http:
	go run . http

# Declare targets that are not files
.PHONY: format vet dockerup dockerdown migrate-create migrate-up migrate-down migrate-drop http
