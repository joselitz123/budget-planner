# SQL Schema Directory

This directory contains SQL migration files for the Budget Planner database.

## Structure

- Migration files should be named with a timestamp prefix (e.g., `001_initial_schema.sql`)
- The init-db.sql file in `.devcontainer/` contains the complete initial schema
- When adding new migrations, run: `migrate create -ext sql -dir sql/schema -seq <migration_name>`

## Current Schema

The database schema includes:
- Users
- Categories
- Budgets
- Budget Categories
- Payment Methods
- Transactions
- Reflections
- Reflection Questions
- Reflection Templates
- Template Questions
- Share Invitations
- Share Access
- Sync Operations
- Activity Log

See `.devcontainer/init-db.sql` for the complete schema definition.
