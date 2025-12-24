-- name: GetUserByClerkID :one
SELECT * FROM users
WHERE clerk_user_id = $1 AND deleted = false
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (clerk_user_id, email, name, currency)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE(sqlc.narg('email'), email),
    name = COALESCE(sqlc.narg('name'), name),
    currency = COALESCE(sqlc.narg('currency'), currency),
    updated_at = NOW()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: GetCurrentUser :one
SELECT * FROM users
WHERE id = $1 AND deleted = false
LIMIT 1;

-- name: ListAllUsers :many
SELECT * FROM users
WHERE deleted = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
