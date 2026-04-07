# Database Migration Plan: golang-migrate CLI

## Objective

Replace GORM's automatic schema migration (`AutoMigrate`) with explicit, version-controlled SQL migrations using the
`golang-migrate/migrate` CLI. This adheres to industry standards, provides better control over schema changes in
production, and integrates migrations cleanly into our `Taskfile.yml`.

## Scope & Impact

- **Removed:** GORM's `db.AutoMigrate` from application startup sequence (`InitDB`).
- **Added:** Raw SQL migration files (`up` and `down`) for the existing `items` and `users` tables.
- **Added:** New commands in `Taskfile.yml` to create, run, and rollback migrations via the `migrate` CLI.
- **Added:** `DB_URL` environment variable for `golang-migrate` compatibility (since it requires a URL format instead of
  DSN).

## Proposed Solution & Implementation Steps

### 1. Configure Environment Variables

- Add `DB_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"` to `.env` and `example.env`. The
  `migrate` CLI requires a proper URL format rather than a key=value DSN string.

### 2. Create Initial Migration Files

Create a new directory: `internal/adapters/db/migrations/` and generate the initial SQL scripts to match the current DB
models (`Item` and `domain.User`).

* **`000001_create_users_table.up.sql`**:
  ```sql
  CREATE TABLE users (
      id SERIAL PRIMARY KEY,
      username VARCHAR(255) NOT NULL UNIQUE,
      password_hash VARCHAR(255) NOT NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );
  ```
* **`000001_create_users_table.down.sql`**:
  ```sql
  DROP TABLE IF EXISTS users;
  ```

* **`000002_create_items_table.up.sql`**:
  ```sql
  CREATE TABLE items (
      id SERIAL PRIMARY KEY,
      title VARCHAR(255) NOT NULL,
      price NUMERIC(10, 2) NOT NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      deleted_at TIMESTAMP WITH TIME ZONE
  );
  CREATE INDEX idx_items_deleted_at ON items(deleted_at);
  ```
* **`000002_create_items_table.down.sql`**:
  ```sql
  DROP TABLE IF EXISTS items;
  ```

### 3. Remove GORM AutoMigrate

In `internal/adapters/db/connection.go`:

- Remove `err = db.AutoMigrate(&Item{}, &domain.User{})` and its associated error check from the `InitDB` function.

### 4. Update Taskfile.yml

Introduce `dotenv: ['.env']` (if not already at the top-level) and add the following `migrate` tasks:

```yaml
  migrate:create:
    desc:
      Create a new database migration file (Usage: task migrate:create -- <name>)
    cmds:
      - migrate create -ext sql -dir internal/adapters/db/migrations -seq {{.CLI_ARGS}}

  migrate:up:
    desc: Run all pending database migrations
    cmds:
      - migrate -path internal/adapters/db/migrations -database "$DB_URL" up

  migrate:down:
    desc: Revert the last database migration
    cmds:
      - migrate -path internal/adapters/db/migrations -database "$DB_URL" down 1
```

## Verification & Testing

1. Install the `golang-migrate` CLI (e.g., `brew install golang-migrate` or via `go install`).
2. Run `task db:start` to ensure the fresh Postgres container is up.
3. Apply migrations using `task migrate:up`.
4. Run `task run` to start the application and ensure no runtime panics occur.
5. Create a user and an item via the API to verify schema compatibility.
6. Test rollbacks via `task migrate:down`.
