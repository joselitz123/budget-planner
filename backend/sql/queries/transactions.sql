-- name: ListTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
  AND deleted = false
  AND ($2::date IS NULL OR transaction_date >= $2)
  AND ($3::date IS NULL OR transaction_date <= $3)
  AND ($4 = '' OR category_id = $4::uuid)
  AND ($5 = '' OR budget_id = $5::uuid)
ORDER BY transaction_date DESC
LIMIT $6 OFFSET $7;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = $1 AND deleted = false
LIMIT 1;

-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, budget_id, category_id, payment_method_id, 
    amount, type, is_transfer, transfer_to_account_id,
    description, transaction_date, is_recurring, recurrence_pattern
)
VALUES (
    $1, $2, $3, $4, 
    $5, $6, $7, $8, 
    $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateTransaction :one
UPDATE transactions
SET
    budget_id = COALESCE(sqlc.narg('budget_id'), budget_id),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    payment_method_id = COALESCE(sqlc.narg('payment_method_id'), payment_method_id),
    amount = COALESCE(sqlc.narg('amount'), amount),
    type = COALESCE(sqlc.narg('type'), type),
    is_transfer = COALESCE(sqlc.narg('is_transfer'), is_transfer),
    transfer_to_account_id = COALESCE(sqlc.narg('transfer_to_account_id'), transfer_to_account_id),
    description = COALESCE(sqlc.narg('description'), description),
    transaction_date = COALESCE(sqlc.narg('transaction_date'), transaction_date),
    is_recurring = COALESCE(sqlc.narg('is_recurring'), is_recurring),
    recurrence_pattern = COALESCE(sqlc.narg('recurrence_pattern'), recurrence_pattern),
    updated_at = NOW()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: DeleteTransaction :exec
UPDATE transactions
SET deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: GetTransactionsByBudget :many
SELECT * FROM transactions
WHERE budget_id = $1 AND deleted = false
ORDER BY transaction_date DESC;

-- name: GetCategorySpent :one
SELECT COALESCE(SUM(t.amount), 0) as total_spent
FROM transactions t
WHERE t.budget_id = $1
  AND t.category_id = $2
  AND t.type = 'expense'
  AND t.deleted = false
  AND t.transaction_date >= (SELECT month FROM budgets WHERE id = $1)
  AND t.transaction_date < ((SELECT month FROM budgets WHERE id = $1) + INTERVAL '1 month');
