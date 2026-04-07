# JWT Authentication Implementation Plan

This document outlines the steps to add JWT-based authentication to the Go Gin REST API.

## 1. Infrastructure & Configuration

- [x] Install JWT library: `go get github.com/golang-jwt/jwt/v5`
- [x] Update `internal/config/config.go`:
    - Add `JWTSecret` (string) and `JWTExpiration` (time.Duration).
- [x] Update `.env` and `example.env` with default values.

## 2. Domain Layer (Ports & Entities)

- [x] Create `internal/domain/user.go`:
    - Define `User` struct (ID, Username, PasswordHash, CreatedAt).
    - Define `RegisterRequest` and `LoginRequest` structs.
- [x] Update `internal/domain/ports.go`:
    - Add `UserRepository` (Outbound Port): `GetByUsername`, `Create`.
    - Add `AuthService` (Inbound Port): `Register`, `Login`, `ValidateToken`.

## 3. Core Layer (Business Logic)

- [x] Implement `internal/core/auth_service.go`:
    - **Register**: Hash password using `bcrypt`, save user.
    - **Login**: Verify password, generate signed JWT.
    - **ValidateToken**: Parse and verify JWT claims.
- [x] Add unit tests in `internal/core/auth_service_test.go`.

## 4. Adapter Layer (Persistence & Web)

- [x] Implement `internal/adapters/db/postgres_user_repository.go`:
    - GORM implementation for user persistence.
    - Add `User` to `AutoMigrate` in `connection.go`.
- [x] Implement `internal/adapters/web/auth_handler.go`:
    - Handlers for `POST /auth/register` and `POST /auth/login`.
- [x] Implement `JWTAuthMiddleware` in `internal/adapters/web/middleware.go` (or update existing handler):
    - Extract `Authorization: Bearer <token>`.
    - Validate token and set `userID` in Gin context.

## 5. Integration & Wiring

- [x] Update `cmd/api/main.go`:
    - Initialize `UserRepository` and `AuthService`.
    - Register new auth routes.
    - Apply `JWTAuthMiddleware` to protected `/api/v1` routes.
- [x] Run `task mock` to generate mocks for new interfaces.
- [x] Final end-to-end verification.
