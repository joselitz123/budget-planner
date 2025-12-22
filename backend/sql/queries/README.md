# SQL Queries Directory

This directory contains SQL query files that will be used by sqlc to generate Go code.

## How It Works

1. Write SQL queries in `.sql` files in this directory
2. Run `sqlc generate` to generate Go code in `internal/models/`
3. Import and use the generated models and queries in your Go code

## Query File Format

Each query file should follow this naming convention:
- Use descriptive names like `users.sql`, `budgets.sql`, `transactions.sql`
- Use comments (`-- name: <name> <operation>`) to name your queries
- Example:

```sql
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: ListBudgets :many
SELECT * FROM budgets
WHERE user_id = $1
AND NOT deleted
ORDER BY month DESC;
```

## Operations

- `:one` - Returns a single row (returns *Model or error)
- `:many` - Returns multiple rows (returns []Model, error)
- `:exec` - Executes a query without returning rows (returns sql.Result, error)
- `:execrows` - Executes and returns number of rows affected (returns int64, error)

## Generated Code

After running `sqlc generate`, you'll find:
- `internal/models/models.go` - Generated struct definitions
- `internal/models/<table>.go` - Query functions for each query file

## Running Generation

From the backend directory:
```bash
sqlc generate
```

Or with air (if watching is enabled):
```bash
air
