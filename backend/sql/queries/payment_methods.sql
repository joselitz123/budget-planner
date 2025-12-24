-- name: ListPaymentMethods :many
SELECT * FROM payment_methods
WHERE user_id = $1 AND deleted = false
ORDER BY is_default DESC, created_at DESC;

-- name: GetPaymentMethodByID :one
SELECT * FROM payment_methods
WHERE id = $1 AND deleted = false
LIMIT 1;

-- name: CreatePaymentMethod :one
INSERT INTO payment_methods (
    user_id, name, type, last_four, brand,
    is_default, is_active, credit_limit, current_balance
)
VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9
)
RETURNING *;

-- name: UpdatePaymentMethod :one
UPDATE payment_methods
SET
    name = COALESCE(sqlc.narg('name'), name),
    type = COALESCE(sqlc.narg('type'), type),
    last_four = COALESCE(sqlc.narg('last_four'), last_four),
    brand = COALESCE(sqlc.narg('brand'), brand),
    is_default = COALESCE(sqlc.narg('is_default'), is_default),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    credit_limit = COALESCE(sqlc.narg('credit_limit'), credit_limit),
    current_balance = COALESCE(sqlc.narg('current_balance'), current_balance),
    updated_at = NOW()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: DeletePaymentMethod :exec
UPDATE payment_methods
SET deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: SetDefaultPaymentMethod :exec
UPDATE payment_methods
SET is_default = false
WHERE user_id = $1 AND deleted = false;

UPDATE payment_methods
SET is_default = true, updated_at = NOW()
WHERE id = $2 AND deleted = false;
