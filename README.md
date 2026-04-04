# Go Gin REST API - Learning Project

Welcome to your Go learning project! This repository contains a simple REST API built with the [Gin Web Framework](https://gin-gonic.com/) and [GORM](https://gorm.io/) (connected to PostgreSQL). 

The project has been refactored to follow **Hexagonal Architecture** (also known as Ports and Adapters) to demonstrate how to build maintainable, decoupled, and testable applications in Go.

---

## 🏗 Directory Structure

This project follows Go's standard project layout conventions.

```
/cmd
  /api
    main.go      <-- Application entry point, dependency injection
/internal
  /adapters      <-- Implementations of ports
    /db
      connection.go           <-- DB initialization and migrations
      postgres_repository.go  <-- Outbound Adapter for PostgreSQL
    /web
      gin_handler.go          <-- Inbound Adapter for Gin
  /config
    config.go        <-- Configuration loading and environment variables
  /core
    item_service.go  <-- Core application logic, implements ItemService port
  /domain
    item.go          <-- Core business models (entities)
    ports.go         <-- Port interfaces (ItemService, ItemRepository)
go.mod
README.md
```

### Hexagonal Architecture (Ports & Adapters)

*   **`/cmd`**: Contains the application's entry point (`cmd/api/main.go`). Its only job is to load configuration, wire together the dependencies (Dependency Injection) from the `internal` directory, and start the server.
*   **`/internal`**: This is the heart of the application.
    *   **`/internal/domain`**: The very center of the application. It has **zero dependencies** on external libraries.
    *   **`/internal/core`**: Implements the business logic. It depends only on the domain and is injected with outbound ports.
    *   **`/internal/adapters`**: The bridge between the core logic and the outside world. Components here are injected with their dependencies (e.g., the DB handle).
    *   **`/internal/config`**: Manages application settings, allowing for clean injection of secrets and parameters.

---

## 🏗 Dependency Injection & Clean Code

The project has been refactored to strictly follow Dependency Injection (DI) principles:
*   **Constructor Injection**: All services and adapters are initialized via constructors (e.g., `NewItemService`, `NewItemHandler`) that clearly define their dependencies.
*   **Decoupled DB**: The database connection logic is separated from the repository implementation, allowing the repository to be tested with any `*gorm.DB` handle.
*   **Injected Secrets**: Middleware no longer uses hardcoded secrets; the `apiToken` is injected into the web handler during initialization.
*   **Structured Logging**: The project uses `uber-go/zap` for high-performance, structured logging. The logger is initialized in the entry point and injected into adapters, ensuring consistent and searchable logs.

## 🧪 Testing

The project uses `stretchr/testify` for assertions and mocking.

### Running Tests
To run all unit tests in the project:
```bash
go test ./...
```

For verbose output:
```bash
go test -v ./...
```

### Testing Strategy
*   **Unit Tests**: Located alongside the code (e.g., `item_service_test.go`). These use mocks to isolate the component being tested.
*   **Mocks**: We use `testify/mock` to create mock implementations of our interfaces (`ItemRepository`, `ItemService`), allowing us to test each layer in isolation.

---

## 🚀 Gin Framework Concepts Covered

*   **Routing & Grouping**: We use `router.Group("/api/v1")` to organize endpoints and apply middleware to specific sets of routes.
*   **Middleware**: Functions that run before your main handler.
    *   *Global Middleware*: `CustomLogger` measures how long every single request takes.
    *   *Group Middleware*: `AuthMiddleware` checks for an `X-API-Token` header.
*   **Data Validation**: Using Gin's integration with the `validator` package. In `domain/item.go`, tags like `binding:"required,min=3"` ensure incoming JSON automatically meets our rules.
*   **Path Parameters**: Extracting variables from the URL, like `:id` in `router.GET("/items/:id")`.

---

## 🛠 How to Run the Project

### Prerequisites
1.  **Go**: Make sure Go is installed (`go version`).
2.  **PostgreSQL**: You need a running PostgreSQL database. 

You can easily start a PostgreSQL instance using Docker:
```bash
docker run --name my-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```

### Setup & Run
1.  Download dependencies:
    ```bash
    go mod tidy
    ```
2. Start the Go server:
    ```bash
    go run cmd/api/main.go
    ```
    *You should see a message: "Database connection established and migrations completed." The database tables are automatically created by GORM.*

### Configuration (Environment Variables)
The project uses environment variables for configuration. Default values are provided for local development:
*   `DB_DSN`: PostgreSQL connection string (Default: `host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC`)
*   `API_TOKEN`: Secret token for `X-API-Token` header (Default: `secret123`)
*   `SERVER_ADDR`: Port the server listens on (Default: `:8080`)
*   `LOG_LEVEL`: Logger verbosity: `debug`, `info`, `warn`, `error` (Default: `info`)
*   `APP_ENV`: Application environment: `development` for console-friendly logs, `production` for JSON logs (Default: `production`)

### Testing the API


The API is secured with a simple token. You must include the header `X-API-Token: secret123` in your requests.

**1. Create an Item:**
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "X-API-Token: secret123" \
  -H "Content-Type: application/json" \
  -d '{"title": "Learning Go", "price": 49.99}'
```

**2. Get All Items:**
```bash
curl -H "X-API-Token: secret123" http://localhost:8080/api/v1/items
```

**3. Get a Specific Item (replace 1 with your item ID):**
```bash
curl -H "X-API-Token: secret123" http://localhost:8080/api/v1/items/1
```
