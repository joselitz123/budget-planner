-- name: ListUserBudgets :many
SELECT * FROM budgets
WHERE user_id = $1 AND deleted = false
ORDER BY month DESC;

-- name: GetBudgetByID :one
SELECT * FROM budgets
WHERE id = $1 AND deleted = false
LIMIT 1;

-- name: GetBudgetByMonth :one
SELECT * FROM budgets
WHERE user_id = $1 AND month = $2 AND deleted = false
LIMIT 1;

-- name: CreateBudget :one
INSERT INTO budgets (user_id, name, month, total_limit)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateBudget :one
UPDATE budgets
SET
    name = COALESCE(sqlc.narg('name'), name),
    total_limit = COALESCE(sqlc.narg('total_limit'), total_limit),
    updated_at = NOW()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: DeleteBudget :exec
UPDATE budgets
SET deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: GetBudgetCategories :many
SELECT bc.*, c.name, c.icon, c.color
FROM budget_categories bc
JOIN categories c ON bc.category_id = c.id
WHERE bc.budget_id = $1;

-- name: AddBudgetCategory :one
INSERT INTO budget_categories (budget_id, category_id, limit_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBudgetCategory :one
UPDATE budget_categories
SET limit_amount = COALESCE(sqlc.narg('limit_amount'), limit_amount), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: RemoveBudgetCategory :exec
DELETE FROM budget_categories
WHERE id = $1;

-- name: GetBudgetSpent :one
SELECT COALESCE(SUM(t.amount), 0) as total_spent
FROM transactions t
WHERE t.budget_id = $1 
  AND t.type = 'expense' 
  AND t.deleted = false
  AND t.transaction_date >= (SELECT month FROM budgets WHERE id = $1)
  AND t.transaction_date < ((SELECT month FROM budgets WHERE id = $1) + INTERVAL '1 month');
