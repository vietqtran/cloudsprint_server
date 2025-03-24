.PHONY: all build run test clean migrate sqlc swag docker docker-run

# Project variables
BINARY_NAME=cloudsprint_server
DB_URL=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
MIGRATION_URL=file://db/migration

all: clean sqlc swag build

build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) cmd/api/main.go

run: build
	@echo "Running application..."
	./$(BINARY_NAME)

watch:
	@echo "Watching for file changes..."
	go install github.com/air-verse/air@latest
	air

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	go clean

sqlc:
	@echo "Generating SQLC code..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.20.0
	sqlc generate -f db/sqlc.yaml

swag:
	@echo "Generating Swagger documentation..."
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go -o docs/swagger

postgres:
	@echo "Starting PostgreSQL container..."
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres:14-alpine

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

setup: postgres createdb migrate-up sqlc swag build