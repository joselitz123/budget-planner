# Budget Planner Backend - Development Guide

This file provides comprehensive guidance for Claude Code when working with the Go backend of the Budget Planner project.

## Backend Overview

**Tech Stack:**
- Go 1.23
- Chi v5 router for HTTP routing
- PostgreSQL 16 with pgx/v5 driver
- sqlc for type-safe SQL query generation
- JWT-based authentication (Clerk-compatible)
- Docker Compose for local development

**Current Status:** All 10 handlers implemented with 100% test coverage (48/48 tests passing)

---

## Quick Commands

```bash
cd /workspace/budget-planner/backend

# Run with hot reload (Air)
air

# Generate Go code from SQL queries (sqlc)
/home/vscode/go/bin/sqlc generate

# Run database migrations
migrate -path sql/schema -database "$DATABASE_URL" up

# Build
go build -o ./tmp/main ./cmd/api

# Run all tests
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/ -v

# Run specific test
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/auth_test.go -v
```

---

## Architecture Patterns

### sqlc Query Generation

**Pattern:** SQL queries are written in `.sql` files, then Go code is generated.

```sql
-- sql/queries/users.sql
-- name: GetUserByID :one
SELECT id, email, name FROM users WHERE id = $1;
```

```bash
/home/vscode/go/bin/sqlc generate
# Generates: internal/models/users.sql.go with GetUserByID() function
```

**Query Naming Convention:**
- `:one` - Returns single row
- `:many` - Returns multiple rows
- `:exec` - Executes without returning data
- `:execrows` - Executes and returns rows affected

**Always prefer using sqlc-generated queries over raw SQL in handlers.**

### Database Schema Conventions

All tables follow these conventions:
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
user_id UUID REFERENCES users(id) ON DELETE CASCADE
created_at TIMESTAMPTZ DEFAULT NOW()
updated_at TIMESTAMPTZ DEFAULT NOW()
deleted BOOLEAN DEFAULT FALSE
```

### Type Conversion Patterns (pgtype ↔ Go)

The backend uses `pgtype` types from `github.com/jackc/pgx/v5/pgtype`. **Always use helpers from `internal/utils/types.go`:**

```go
// Go → pgtype (for database operations)
utils.PgUUID("user-id-123")
utils.PgText("user name")
utils.PgNumeric(123.45)
utils.PgDate(time.Now())
utils.PgBool(true)

// pgtype → Go (from database results)
utils.UUIDToString(uuidValue)   // pgtype.UUID → string
utils.TextToString(textValue)   // pgtype.Text → string
utils.NumericToFloat64(numericValue) // pgtype.Numeric → float64
```

**Important:** Some sqlc-generated model fields use plain Go types, not pgtype:
- `CreatePaymentMethodParams.Name` is `string`, not `pgtype.Text`
- `CreatePaymentMethodParams.Type` is `string`, not `pgtype.Text`

### Chi Router Path Parameter Handling

**CRITICAL:** All handlers must use `chi.URLParam(r, "paramName")` to extract path values:

```go
// ✅ Correct
id := chi.URLParam(r, "id")

// ❌ Wrong - returns empty in tests without Chi context
id := r.PathValue("id")
```

---

## Handler Reference

All 10 handlers are implemented and tested (48/48 tests passing):

| # | Handler | Endpoints | Purpose |
|---|---------|-----------|---------|
| 1 | Auth | `POST /api/auth/login`, `POST /api/auth/logout`, `POST /api/auth/onboarding`, `GET /api/auth/me` | Authentication and user onboarding |
| 2 | Users | `GET /api/users/me`, `PUT /api/users/me`, `DELETE /api/users/me` | User profile management |
| 3 | Categories | `GET /api/categories`, `GET /api/categories/system`, `POST /api/categories`, `PUT /api/categories/{id}`, `DELETE /api/categories/{id}` | Category CRUD |
| 4 | Budgets | `GET /api/budgets`, `GET /api/budgets/{month}`, `GET /api/budgets/id/{id}`, `POST /api/budgets`, `PUT /api/budgets/{id}`, `DELETE /api/budgets/{id}`, `GET /api/budgets/{id}/categories`, `POST /api/budgets/{id}/categories` | Budget management |
| 5 | Transactions | `GET /api/transactions`, `GET /api/transactions/{id}`, `POST /api/transactions`, `PUT /api/transactions/{id}`, `DELETE /api/transactions/{id}`, `GET /api/budgets/{budgetId}/transactions` | Transaction CRUD |
| 6 | Payment Methods | `GET /api/payment-methods`, `GET /api/payment-methods/{id}`, `POST /api/payment-methods`, `PUT /api/payment-methods/{id}`, `DELETE /api/payment-methods/{id}` | Payment method CRUD |
| 7 | Sync | `POST /api/sync/push`, `POST /api/sync/pull`, `GET /api/sync/status`, `POST /api/sync/resolve-conflict`, `GET /api/sync/conflict-history` | Offline sync operations |
| 8 | Reflections | `GET /api/reflections/month/{month}`, `POST /api/reflections`, `PUT /api/reflections/{id}`, `DELETE /api/reflections/{id}`, `GET /api/reflections/templates` | Monthly reflections |
| 9 | Sharing | `POST /api/sharing/invitations`, `GET /api/sharing/invitations`, `POST /api/sharing/invitations/{id}/respond`, `GET /api/sharing/budgets/{budgetId}`, `GET /api/sharing/access/{budgetId}`, `DELETE /api/sharing/access/{id}`, `DELETE /api/sharing/invitations/{id}` | Budget sharing |
| 10 | Analytics | `GET /api/analytics/dashboard/{month}`, `GET /api/analytics/spending/{month}`, `GET /api/analytics/trends`, `GET /api/analytics/category/{categoryId}` | Reports and analytics |

---

## Testing Guide

### Test Infrastructure

Tests use `internal/handlers/test_setup.go` which provides:

**Helper Functions:**
- `CreateTestUser(t, ctx)` - Creates a test user with unique ID
- `CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)` - Creates a test budget
- `CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")` - Creates a test category
- `CreateTestPaymentMethod(t, ctx, userID, "Cash", "cash")` - Creates a test payment method
- `CleanupTestUser(t, ctx, userID)` - Soft-deletes test user (cascades to related records)
- `setAuthContext(req, userID)` - Sets authenticated user + Chi router context
- `GenerateTestToken(t, clerkUserID)` - Creates a valid JWT for testing

### Running Tests

```bash
# Run all handler tests
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/ -v

# Run specific test file
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/auth_test.go -v

# Run with coverage
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./... -cover
```

### Test Status

**Overall: 48 out of 48 tests passing (100%)** ✅

| Handler | Passing | Total | Status |
|---------|---------|-------|--------|
| Auth | 5 | 5 | ✅ All passing |
| Budgets | 8 | 8 | ✅ All passing |
| Categories | 6 | 6 | ✅ All passing |
| Payment Methods | 6 | 6 | ✅ All passing |
| Reflections | 6 | 6 | ✅ All passing |
| Sharing | 7 | 7 | ✅ All passing |
| Sync | 5 | 5 | ✅ All passing |
| Transactions | 7 | 7 | ✅ All passing |
| Analytics | 4 | 4 | ✅ All passing |

**See `internal/utils/types.go` for complete type conversion helper list.**

---

## Development Workflow

### Hot Reload with Air

- Backend uses Air (configured in `.air.toml`)
- Watches `.go` and `.sql` files
- Automatically rebuilds and restarts on changes
- `sqlc generate` is manual - run after modifying `.sql` files

### Database Migrations

```bash
# Run migrations
migrate -path sql/schema -database "$DATABASE_URL" up

# Rollback migrations
migrate -path sql/schema -database "$DATABASE_URL" down

# Create new migration (numbered sequence)
# 1. Create sql/schema/002_<name>.up.sql
# 2. Create sql/schema/002_<name>.down.sql
# 3. Test both up and down
```

---

## Project Structure

```
backend/
├── cmd/api/
│   └── main.go              # Entry point, route registration
├── internal/
│   ├── auth/
│   │   ├── clerk.go         # JWT client (GenerateToken for testing)
│   │   └── middleware.go    # Auth middleware (SetUserIDInContext helper)
│   ├── config/
│   │   └── config.go        # Environment variable loading
│   ├── database/
│   │   ├── connection.go    # PostgreSQL connection pool
│   │   └── migrations.go    # Schema migration runner
│   ├── handlers/
│   │   ├── *.go             # HTTP request handlers (one per domain)
│   │   ├── *_test.go        # Unit tests for each handler
│   │   └── test_setup.go    # Common test utilities
│   ├── middleware/
│   │   ├── error.go         # Panic recovery and error handling
│   │   ├── logging.go       # Request logging
│   │   └── permission.go    # Budget sharing permission checks
│   ├── models/              # GENERATED by sqlc - do not edit
│   │   ├── models.go        # Table struct definitions
│   │   └── *.sql.go         # Generated query functions
│   └── utils/
│       ├── response.go      # JSON response helpers
│       └── types.go         # Type conversion helpers (PgUUID, PgText, etc.)
├── sql/
│   ├── schema/
│   │   ├── *.up.sql         # Migration files (up)
│   │   └── *.down.sql       # Migration files (down)
│   └── queries/
│       └── *.sql            # SQL queries for sqlc
├── sqlc.yaml                # sqlc configuration
├── .air.toml                # Air hot reload config
├── .env                     # Environment variables (create from .env.example)
└── docker-compose.yml       # Local PostgreSQL
```

---

## Important Files

- `todo.md` - Implementation status and session history (renamed from TASK.md)
- `starting-point.md` (root) - Full project specification
- `CLAUDE.md` (root) - Project-level documentation
- `internal/handlers/*.go` - All 10 handlers
- `internal/utils/types.go` - Type conversion helpers
- `internal/handlers/test_setup.go` - Test infrastructure

---

## Environment Configuration

Required environment variables (`.env`):

```bash
DATABASE_URL=postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable
PORT=8080
JWT_SECRET=dev-secret-key-change-in-production
```

For Clerk authentication in development:
```bash
CLERK_SECRET_KEY=your_clerk_secret_key
CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
```

---

## Authentication Flow

1. Frontend uses Clerk SDK for authentication
2. Backend validates JWT tokens via `internal/auth/middleware.go`
3. User context attached to requests by `auth.SetUserIDInContext()`
4. Clerk user ID stored in `users.clerk_user_id` for lookups

---

## Permission Model for Sharing

Budgets can be shared with two permission levels:
- `view` - Read-only access
- `edit` - Full edit access

Permission checks occur via the `share_access` table. Always check:
1. Is user the owner? → Full access
2. Is user in share_access with edit permission? → Edit access
3. Is user in share_access with view permission? → Read-only
4. Otherwise → 404/403
