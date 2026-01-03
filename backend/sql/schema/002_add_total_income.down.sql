-- Rollback: Remove total_income column from budgets table

BEGIN;

-- Drop index if it exists
DROP INDEX IF EXISTS idx_budgets_total_income;

-- Remove total_income column
ALTER TABLE budgets DROP COLUMN IF EXISTS total_income;

COMMIT;
