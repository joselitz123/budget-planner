-- Drop tables in reverse order of creation

DROP TABLE IF EXISTS activity_log;
DROP TABLE IF EXISTS sync_operations;
DROP TABLE IF EXISTS share_access;
DROP TABLE IF EXISTS share_invitations;
DROP TABLE IF EXISTS template_questions;
DROP TABLE IF EXISTS reflection_templates;
DROP TABLE IF EXISTS reflection_questions;
DROP TABLE IF EXISTS reflections;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS budget_categories;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pgcrypto";
