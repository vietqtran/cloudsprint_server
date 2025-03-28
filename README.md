# Go Postgres API

A production-ready RESTful API built with Go, Fiber, PostgreSQL, sqlc, and JWT authentication.

## Features

- RESTful API architecture
- Structured project layout
- PostgreSQL database with migrations
- JWT-based authentication and authorization
- SQL query generation with sqlc
- Swagger API documentation
- Request validation
- Request/error logging
- Graceful shutdown
- Configuration management
- Containerization with Docker
- Makefile for common operations

## Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Docker (optional)

## Project Structure

```
cloudsprint_server/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   ├── config.go
│   └── config.yaml
├── db/
│   ├── migration/
│   │   ├── 000001_init_schema.up.sql
│   │   └── 000001_init_schema.down.sql
│   ├── query/
│   │   └── user.sql
│   └── sqlc.yaml
├── docs/
│   └── swagger/
│       └── swagger.yaml
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── request/
│   │   ├── response/
│   │   └── router.go
│   ├── db/
│   │   └── sqlc/
│   ├── logger/
│   │   └── logger.go
│   └── token/
│       └── jwt_maker.go
├── pkg/
│   └── util/
│       ├── password.go
│       └── random.go
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Getting Started

### Manual Setup

1. Clone the repository:

```bash
git clone https://cloud-sprint.git
cd cloudsprint_server
```

2. Start a PostgreSQL instance:

```bash
make postgres
```

3. Create a database:

```bash
make createdb
```

4. Run database migrations:

```bash
make migrate-up
```

5. Generate SQL code with sqlc:

```bash
make sqlc
```

6. Generate Swagger documentation:

```bash
make swag
```

7. Build and run the application:

```bash
make run
```

### Docker Setup

1. Clone the repository:

```bash
git clone https://cloud-sprint.git
cd cloudsprint_server
```

2. Build and run with Docker:

```bash
make docker-run
```

## Quick Setup

For a one-step setup:

```bash
make setup
```

This will:
1. Start PostgreSQL in Docker
2. Create the database
3. Run migrations
4. Generate sqlc code
5. Generate Swagger documentation
6. Build the application

## Configuration

The application is configured using the `config/config.yaml` file. You can override these settings with environment variables.

## API Documentation

Swagger UI is available at `http://localhost:8080/swagger/`

## Common Commands

- `make run`: Build and run the application
- `make build`: Build the application
- `make test`: Run tests
- `make migrate-up`: Run database migrations
- `make migrate-down`: Rollback migrations
- `make sqlc`: Generate SQL code
- `make swag`: Generate Swagger documentation
- `make docker`: Build Docker image
- `make docker-run`: Run the application in Docker
- `make setup`: Setup everything from scratch

## License

This project is licensed under the MIT License - see the LICENSE file for details.l
