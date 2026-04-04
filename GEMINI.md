# GEMINI.md

This file provides context and instructions for the **Go Gin REST API - Learning Project**.

## Project Overview
A learning project demonstrating a modern REST API built with **Go (1.26)** using **Hexagonal Architecture** (Ports and Adapters). The project is designed to be highly maintainable, decoupled, and testable.

### Core Technologies
- **Web Framework:** [Gin](https://gin-gonic.com/)
- **ORM:** [GORM](https://gorm.io/) (with PostgreSQL driver)
- **Logging:** [zap](https://github.com/uber-go/zap) (Structured, high-performance logging)
- **Testing:** [testify](https://github.com/stretchr/testify) (Assertions and Mocking)
- **Architecture:** Hexagonal (Ports & Adapters)
- **Mocking:** [mockery](https://github.com/vektra/mockery) (Automated mock generation)

## Building and Running

### Prerequisites
- **Go 1.26+**
- **PostgreSQL** (can be run via Docker):
  ```bash
  docker run --name my-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
  ```

### Commands
- **Install Dependencies:** `go mod tidy`
- **Run the API:** `go run cmd/api/main.go`
- **Run Tests:** `go test ./...` (use `-v` for verbose output)
- **Generate Mocks:** `mockery` (Uses configuration in `.mockery.yaml`)

### Environment Variables
| Variable | Description | Default Value |
| :--- | :--- | :--- |
| `DB_DSN` | PostgreSQL connection string | `host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC` |
| `API_TOKEN` | Secret token for `X-API-Token` header | `secret123` |
| `SERVER_ADDR` | API listen address | `:8080` |
| `LOG_LEVEL` | Logger verbosity (`debug`, `info`, `warn`, `error`) | `info` |

## Development Conventions

### Hexagonal Architecture
The project strictly follows the Ports and Adapters pattern to decouple business logic from external concerns:
- **`internal/domain`**: Contains core business entities (`Item`) and interface definitions (**Ports**).
- **`internal/core`**: Implements the business logic (Inbound Ports). It depends only on the domain.
- **`internal/adapters`**: Implements the **Adapters** for external systems:
    - `db`: Outbound adapter for PostgreSQL using GORM.
    - `web`: Inbound adapter for the Gin framework.

### Dependency Injection (DI)
All components are wired together using **Constructor Injection**. The application entry point (`cmd/api/main.go`) is responsible for:
1. Loading configuration.
2. Initializing the logger.
3. Establishing the database connection.
4. Injecting dependencies into adapters and services.

### Logging
Always use the injected `*zap.Logger` for structured logging. Avoid using `fmt.Println` or standard `log` package.

### Testing Strategy
- **Unit Tests:** Located alongside the code (e.g., `item_service_test.go`).
- **Mocking:** Use [mockery](https://github.com/vektra/mockery) for automated mock generation. Mocks are stored in `internal/mocks`.
- **Assertions:** Use `testify/assert` or `testify/require`. For mock expectations, prefer the type-safe `EXPECT()` API provided by `mockery`.

### API Security
Endpoints under `/api/v1` are protected by `AuthMiddleware`. Requests must include the `X-API-Token` header.

## Directory Structure
- `cmd/api/`: Application entry point and dependency wiring.
- `internal/adapters/`: External system implementations (DB, Web).
- `internal/config/`: Configuration loading and environment variables.
- `internal/core/`: Application services and business logic.
- `internal/domain/`: Core entities and port interfaces.
- `internal/logger/`: Centralized logger configuration.
