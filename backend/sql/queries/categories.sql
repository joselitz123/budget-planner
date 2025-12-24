-- name: GetUserCategories :many
SELECT * FROM categories
WHERE user_id = $1 AND deleted = false
ORDER BY created_at ASC;

-- name: GetSystemCategories :many
SELECT * FROM categories
WHERE is_system = true AND deleted = false
ORDER BY name ASC;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 AND deleted = false
LIMIT 1;

-- name: CreateCategory :one
INSERT INTO categories (user_id, name, icon, color, is_system, default_limit)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET
    name = COALESCE(sqlc.narg('name'), name),
    icon = COALESCE(sqlc.narg('icon'), icon),
    color = COALESCE(sqlc.narg('color'), color),
    default_limit = COALESCE(sqlc.narg('default_limit'), default_limit),
    updated_at = NOW()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: DeleteCategory :exec
UPDATE categories
SET deleted = true, updated_at = NOW()
WHERE id = $1;
