.PHONY: all build run test clean migrate sqlc swag docker docker-run format lint pre-commit imports

# Project variables
BINARY_NAME=cloudsprint_server
DB_URL=postgresql://neondb_owner:J1Kmrk5PNRqg@ep-white-tooth-a1uuj3k9.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
MIGRATION_URL=file://db/migration

all: clean sqlc swag build format

build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) cmd/api/main.go

run: build
	@echo "Running application..."
	./$(BINARY_NAME)

watch:
	@echo "Watching for file changes..."
	go install github.com/air-verse/air@latest
ifeq ($(OS),Windows_NT)
	air
else
	air -c .air-linux.toml
endif

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	go clean

sqlc:
	@echo "Generating SQLC code..."
ifeq ($(OS),Windows_NT)
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	sqlc generate -f db/sqlc.yaml
else
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	sqlc generate -f db/sqlc.yaml
endif

swag:
	@echo "Generating Swagger documentation..."
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go -o docs/swagger

postgres:
	@echo "Starting PostgreSQL container..."
	docker stop postgres
	docker rm postgres
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres:14-alpine
	@sleep 10

createdb:
	@echo "Creating database..."
	docker exec -it postgres createdb --username=postgres --owner=postgres postgres

dropdb:
	@echo "Dropping database..."
	docker exec -it postgres dropdb --username=postgres postgres

migrate-up:
	@echo "Running migrations up..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -path db/migration -database $(DB_URL) -verbose up

migrate-down:
	@echo "Running migrations down..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -path db/migration -database $(DB_URL) -verbose down

migrate-create:
	@echo "Creating migration files..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate create -ext sql -dir db/migration -seq $(name)

docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

docker-run: docker
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name $(BINARY_NAME) -d $(BINARY_NAME)

format:
	@echo "Formatting code..."

imports:
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w ./

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

pre-commit: format imports lint

setup: clean postgres migrate-up sqlc swag build format