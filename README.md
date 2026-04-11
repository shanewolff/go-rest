# Go Gin REST API - Hexagonal Architecture

Welcome to the **Go Gin REST API** learning project. This repository serves as a definitive guide and production-ready
template for building maintainable, decoupled, and highly testable Go applications using **Hexagonal Architecture** (
Ports and Adapters).

---

## 📖 Table of Contents

- [Project Overview](#-project-overview)
- [Architecture Deep Dive](#-architecture-deep-dive)
- [Business Features](#-business-features)
- [Getting Started](#-getting-started)
- [Developer Maintenance & Workflow](#-developer-maintenance--workflow)
    - [Task Runner (Taskfile)](#task-runner-taskfile)
    - [Testing & Coverage Strategy](#testing--coverage-strategy)
    - [Mocks & Dependency Injection](#mocks--dependency-injection)
    - [Git Hooks & Quality Gates](#git-hooks--quality-gates)
- [Security & Authentication](#-security--authentication)
- [Database Management](#-database-management)
- [API Reference](#-api-reference)

---

## 🌟 Project Overview

This project demonstrates a modern RESTful API built with:

- **Framework:** [Gin Gonic](https://gin-gonic.com/) (High-performance HTTP web framework)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL
- **Logging:** [Zap](https://github.com/uber-go/zap) (Structured, lightning-fast logging)
- **Validation:** [Go Playground Validator](https://github.com/go-playground/validator)
- **Architecture:** Hexagonal (Ports & Adapters) for maximum decoupling.

---

## 🏗 Architecture Deep Dive

The project is structured to separate business logic from technical implementation details (like databases or web
frameworks).

### The Layers

1. **Domain (`internal/domain`)**: The core. Contains business entities (`Item`, `User`) and **Port Interfaces**. It has
   zero dependencies on external libraries or other layers.
2. **Core/Application (`internal/core`)**: Implements the business logic (Services). It coordinates tasks and delegates
   data persistence to the Domain Ports.
3. **Adapters (`internal/adapters`)**:
    * **Inbound (Primary)**: The Entry points. In this project, it's the **Web Adapter** (Gin handlers) that translates
      HTTP requests into Domain calls.
    * **Outbound (Secondary)**: External integrations. Here, it's the **DB Adapter** (GORM) that translates Domain calls
      into SQL queries.
4. **CMD (`cmd/api`)**: The "Main" entry point. Its sole responsibility is **Dependency Injection (DI)**—wiring the
   adapters to the services and starting the engine.

---

## 🚀 Business Features

The API provides a complete flow for managing items within a secured environment:

### Authentication & User Management

- **User Registration**: Secure sign-up with password hashing (Bcrypt).
- **JWT Login**: Issue JSON Web Tokens for stateless authentication.
- **Identity Verification**: Middleware to protect sensitive routes.

### Item Management

- **Create**: Add new items with title and price.
- **List**: Retrieve all items.
- **Detail**: Fetch a single item by its ID.
- **Delete**: Remove items from the catalog.

---

## 🏁 Getting Started

### 1. Prerequisites

- **Go 1.26+**
- **Docker** (for running PostgreSQL)
- **Task** runner: `brew install go-task` (or `go install github.com/go-task/task/v3/cmd/task@latest`)
- **Migrate** tool: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### 2. Setup Environment

Copy the example environment file:

```bash
cp example.env .env
```

### 3. Start the Engine

Use the automated tasks to get running in seconds:
```bash
task db:start      # Starts PostgreSQL in Docker
task migrate:up    # Runs database migrations
task tidy          # Installs Go dependencies
task run           # Starts the API server
```

---

## 🛠 Developer Maintenance & Workflow

### Task Runner (Taskfile)

We use `Task` instead of `Make` for a better DX. Key commands:

- `task run`: Start local development server.
- `task test`: Run all tests.
- `task mock`: Regenerate all mocks using Mockery.
- `task lint`: Run golangci-lint.
- `task clean`: Remove build artifacts and coverage files.

### Testing & Coverage Strategy

We maintain a strict **80% Minimum Coverage** threshold.

- **Unit Tests**: Found in `*_test.go` files next to the source code.
- **Centralized Coverage**: We use a custom script (`scripts/calculate_coverage.sh`) that excludes "noise" like mocks,
  entry points (`main.go`), and boilerplate.
- **Commands**:
    - `task test:coverage`: Quick CLI summary.
    - `task test:coverage-out`: Detailed function-level breakdown.
    - `task test:coverage-html`: Visual report in your browser.

### Mocks & Dependency Injection

We use **Mockery** to generate type-safe mocks for our Port interfaces.

- Interfaces are defined in `internal/domain/ports.go`.
- Mocks are generated into `internal/mocks/`.
- This allows us to test the `Core` logic without needing a real database.

### Git Hooks & Quality Gates

We use **Lefthook** to ensure no "bad code" is committed. The `pre-commit` hook automatically runs:

1. **Linter**: `golangci-lint` to check for code smells.
2. **Security Scan**: `gosec` to find potential security vulnerabilities.
3. **Tests**: Runs `go test -short`.
4. **Coverage Check**: Ensures the 80% threshold is met.
5. **Formatters**: Runs `go fmt` and `goimports`.

Set up hooks with: `task hooks:setup`

---

## 🔒 Security & Authentication

- **Password Hashing**: We use `bcrypt` with a default cost of 10. Never store plain-text passwords!
- **JWT (JSON Web Token)**:
    - Signed with a `JWT_SECRET` from your `.env`.
    - Included in requests via the `Authorization: Bearer <token>` header.
- **Middleware**:
    - `JWTAuthMiddleware`: Validates user identity for `/api/v1/*` routes.
    - `APITokenMiddleware`: A secondary layer (legacy) for `X-API-Token` header validation.

---

## 🗄 Database Management

We avoid GORM's `AutoMigrate` in favor of **Explicit Versioned Migrations**. This prevents unexpected schema changes in
production.

- **Migrations Path**: `internal/adapters/db/migrations/`
- **Create New Migration**: `task migrate:create -- name_of_migration`
- **Apply/Rollback**: `task migrate:up` or `task migrate:down`.

---

## 📡 API Reference

### Auth Endpoints

| Method | Endpoint         | Description          |
|:-------|:-----------------|:---------------------|
| `POST` | `/auth/register` | Create a new account |
| `POST` | `/auth/login`    | Get a JWT token      |

### Protected Endpoints (Requires Bearer Token)

| Method   | Endpoint            | Description       |
|:---------|:--------------------|:------------------|
| `GET`    | `/api/v1/items`     | List all items    |
| `GET`    | `/api/v1/items/:id` | Get item details  |
| `POST`   | `/api/v1/items`     | Create a new item |
| `DELETE` | `/api/v1/items/:id` | Remove an item    |

#### Example Request (Create Item)
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Go Concurrency in Practice", "price": 45.00}'
```
