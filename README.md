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
      postgres_repository.go  <-- Outbound Adapter for Items
      postgres_user_repository.go <-- Outbound Adapter for Users
    /web
      gin_handler.go          <-- Inbound Adapter for Items
      auth_handler.go         <-- Inbound Adapter for Authentication
      middleware.go           <-- JWT and API Token middlewares
  /config
    config.go        <-- Configuration loading and environment variables
  /core
    item_service.go  <-- Core item logic
    auth_service.go  <-- Core authentication logic
  /domain
    item.go          <-- Item business models
    user.go          <-- User and Auth business models
    ports.go         <-- Port interfaces (Services, Repositories)
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

* **Constructor Injection**: All services and adapters are initialized via constructors (e.g., `NewItemService`,
  `NewAuthService`, `NewItemHandler`) that clearly define their dependencies.
*   **Decoupled DB**: The database connection logic is separated from the repository implementation, allowing the repository to be tested with any `*gorm.DB` handle.
* **JWT Authentication**: The application uses JSON Web Tokens (JWT) for secure authentication. Passwords are hashed
  using `bcrypt` before being stored in the database.
*   **Structured Logging**: The project uses `uber-go/zap` for high-performance, structured logging. The logger is initialized in the entry point and injected into adapters, ensuring consistent and searchable logs.

## 🧪 Testing

The project uses `stretchr/testify` for assertions and mocking.

### Running Tests
To run all unit tests in the project:
```bash
task test
```

For verbose output:
```bash
task test-v
```

### Testing Strategy

* **Unit Tests**: Located alongside the code (e.g., `item_service_test.go`, `auth_service_test.go`). These use mocks to
  isolate the component being tested.
* **Mocks**: We use [mockery](https://github.com/vektra/mockery) to automatically generate mock implementations of our
  interfaces (`ItemRepository`, `UserRepository`, `ItemService`, `AuthService`). Mocks are stored in `internal/mocks`.

#### Generating Mocks

To generate or update mocks:
```bash
task mock
```
The configuration for mockery is defined in `.mockery.yaml`.

---

## 🚀 Gin Framework Concepts Covered

* **Routing & Grouping**: We use `router.Group("/api/v1")` for protected resources and `router.Group("/auth")` for
  public authentication endpoints.
*   **Middleware**: Functions that run before your main handler.
    *   *Global Middleware*: `CustomLogger` measures how long every single request takes.
    * *Auth Middleware*: `JWTAuthMiddleware` validates the Bearer token in the `Authorization` header.
* **Data Validation**: Using Gin's integration with the `validator` package. In `domain/item.go` and `domain/user.go`,
  tags like `binding:"required,min=6"` ensure incoming JSON automatically meets our rules.
*   **Path Parameters**: Extracting variables from the URL, like `:id` in `router.GET("/items/:id")`.

---

## 🛠 How to Run the Project

### Prerequisites
1.  **Go**: Make sure Go is installed (`go version`).
2.  **PostgreSQL**: You need a running PostgreSQL database. 

You can easily start a PostgreSQL instance using Docker:
```bash
task db:start
```

### Setup & Run
The project uses [Task](https://taskfile.dev/) as a task runner for a simplified developer experience.

1.  List available tasks:
    ```bash
    task --list
    ```
2.  Start the database:
    ```bash
    task db:start
    ```
3.  Download dependencies:
    ```bash
    task tidy
    ```
4.  Start the Go server:
    ```bash
    task run
    ```

### Configuration (Environment Variables)
The project uses environment variables for configuration. For local development, these are loaded from a `.env` file in the project root.

Default values are provided for local development if neither `.env` nor system variables are set:
*   `DB_DSN`: PostgreSQL connection string (Default: `host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC`)
* `API_TOKEN`: Legacy secret token for `X-API-Token` header (Default: `secret123`)
* `JWT_SECRET`: Secret key used for signing JWTs (Default: `super-secret-key`)
* `JWT_EXPIRATION`: Duration before a token expires (Default: `24h`)
*   `SERVER_ADDR`: Port the server listens on (Default: `:8080`)
*   `LOG_LEVEL`: Logger verbosity: `debug`, `info`, `warn`, `error` (Default: `info`)
*   `APP_ENV`: Application environment: `development` for console-friendly logs, `production` for JSON logs (Default: `production`)

### Testing the API

The API is secured with JWT. You must first register and login to get a token.

**1. Register a User:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "password": "securepassword"}'
```

**2. Login to get a Token:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "password": "securepassword"}'
```

*Note the `"token": "..."` in the response.*

**3. Create an Item (using the token):**
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"title": "Learning Go JWT", "price": 59.99}'
```

**4. Get All Items:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" http://localhost:8080/api/v1/items
```
