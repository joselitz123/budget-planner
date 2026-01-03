-- Migration: Add total_income column to budgets table
-- This migration adds support for tracking total income, which enables
-- savings calculation (income - limit) and better budget management

BEGIN;

-- Add total_income column (nullable for backward compatibility)
ALTER TABLE budgets 
ADD COLUMN total_income DECIMAL(12, 2);

-- Add comment for documentation
COMMENT ON COLUMN budgets.total_income IS 'Total income for the budget period. Used to calculate savings (income - limit).';

-- Add index for queries filtering by income ranges (optional, for analytics)
CREATE INDEX idx_budgets_total_income ON budgets(total_income) WHERE NOT deleted;

COMMIT;
