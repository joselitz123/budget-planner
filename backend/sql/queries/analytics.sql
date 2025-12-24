-- name: GetDashboardSummary :one
WITH budget_month AS (
    SELECT month FROM budgets WHERE id = $1
),
spent AS (
    SELECT COALESCE(SUM(amount), 0) as total
    FROM transactions t
    WHERE t.budget_id = $1 
      AND t.type = 'expense' 
      AND t.deleted = false
      AND t.transaction_date >= (SELECT month FROM budget_month)
      AND t.transaction_date < ((SELECT month FROM budget_month) + INTERVAL '1 month')
),
income AS (
    SELECT COALESCE(SUM(amount), 0) as total
    FROM transactions t
    WHERE t.budget_id = $1 
      AND t.type = 'income' 
      AND t.deleted = false
      AND t.transaction_date >= (SELECT month FROM budget_month)
      AND t.transaction_date < ((SELECT month FROM budget_month) + INTERVAL '1 month')
),
transactions_count AS (
    SELECT COUNT(*) as total
    FROM transactions t
    WHERE t.budget_id = $1 AND t.deleted = false
      AND t.transaction_date >= (SELECT month FROM budget_month)
      AND t.transaction_date < ((SELECT month FROM budget_month) + INTERVAL '1 month')
)
SELECT 
    b.*,
    s.total as total_spent,
    i.total as total_income,
    tc.total as transaction_count
FROM budgets b
CROSS JOIN spent s
CROSS JOIN income i
CROSS JOIN transactions_count tc
WHERE b.id = $1;

-- name: GetSpendingByCategory :many
SELECT 
    c.id,
    c.name,
    c.icon,
    c.color,
    COALESCE(SUM(t.amount), 0) as total_spent,
    COALESCE(SUM(t.amount), 0) / bc.limit_amount * 100 as percentage
FROM budget_categories bc
JOIN categories c ON bc.category_id = c.id
LEFT JOIN transactions t ON t.category_id = c.id 
    AND t.budget_id = $1 
    AND t.type = 'expense' 
    AND t.deleted = false
    AND t.transaction_date >= (SELECT month FROM budgets WHERE id = $1)
    AND t.transaction_date < ((SELECT month FROM budgets WHERE id = $1) + INTERVAL '1 month')
WHERE bc.budget_id = $1
GROUP BY c.id, c.name, c.icon, c.color, bc.limit_amount
ORDER BY total_spent DESC;

-- name: GetSpendingTrends :many
SELECT 
    DATE_TRUNC('month', transaction_date) as month,
    SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expenses,
    SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income
FROM transactions
WHERE user_id = $1 
  AND deleted = false
  AND transaction_date >= $2
  AND transaction_date <= $3
GROUP BY DATE_TRUNC('month', transaction_date)
ORDER BY month ASC;

-- name: GetCategoryReport :many
SELECT 
    DATE_TRUNC('day', transaction_date) as date,
    COALESCE(SUM(amount), 0) as total,
    COUNT(*) as transaction_count
FROM transactions
WHERE user_id = $1 
  AND category_id = $2
  AND type = 'expense'
  AND deleted = false
  AND transaction_date >= $3
  AND transaction_date <= $4
GROUP BY DATE_TRUNC('day', transaction_date)
ORDER BY date ASC;

-- name: GetRecentTransactions :many
SELECT t.*, c.name as category_name, c.icon as category_icon, c.color as category_color,
       pm.name as payment_method_name, pm.type as payment_method_type
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
WHERE t.user_id = $1 AND t.deleted = false
ORDER BY t.transaction_date DESC, t.created_at DESC
LIMIT $2 OFFSET $3;
