# Project Configuration
APP_NAME := koda-b8-backend
MAIN := cmd/main.go
BINARY := bin/$(APP_NAME)

MIGRATION_PATH := migrations
DATABASE_URL := postgres://postgres2:1@localhost:5432/weekly?sslmode=disable

# Help
.PHONY: help

help:
	@echo ""
	@echo "========= AVAILABLE COMMANDS ========="
	@echo ""
	@echo "Development"
	@echo "  make run               Run application"
	@echo "  make dev               Run application (alias)"
	@echo ""
	@echo "Build"
	@echo "  make build             Build application"
	@echo "  make clean             Remove binary"
	@echo ""
	@echo "Database"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make migrate-version"
	@echo "  make migrate-force v=1"
	@echo "  make migration name=create_users"
	@echo ""
	@echo "Swagger"
	@echo "  make swagger"
	@echo ""
	@echo "Go"
	@echo "  make tidy"
	@echo "  make fmt"
	@echo "  make vet"
	@echo "  make test"
	@echo ""
	@echo "Docker"
	@echo "  make docker-up"
	@echo "  make docker-down"
	@echo "  make docker-logs"
	@echo ""

# Install dependencies
install:
	@echo "Installing development tools..."

	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Development
run:
	@echo "Starting server..."
	go run $(MAIN)

dev: run

# Build
build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o $(BINARY) $(MAIN)

clean:
	@echo "Cleaning build..."
	rm -rf bin

# Go
tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

# Swagger
swagger:
	@echo "Generating Swagger..."
	swag init -g $(MAIN)

# Migration
migration:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

migrate-up:
	migrate -path $(MIGRATION_PATH) \
	-database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATION_PATH) \
	-database "$(DATABASE_URL)" down

migrate-force:
	migrate -path $(MIGRATION_PATH) \
	-database "$(DATABASE_URL)" force $(v)

migrate-version:
	migrate -path $(MIGRATION_PATH) \
	-database "$(DATABASE_URL)" version

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f
