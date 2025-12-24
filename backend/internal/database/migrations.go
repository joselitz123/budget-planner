package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migrate runs all pending database migrations from the file system
func Migrate(ctx context.Context, db *DB, schemaPath string) error {
	// If schemaPath is empty, use a default relative path
	if schemaPath == "" {
		// Get the working directory or use the sql/schema path
		schemaPath = "sql/schema"
	}

	// Get migration files from file system
	files, err := os.ReadDir(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema directory: %w", err)
	}

	// Sort files to ensure they run in order
	var migrations []string
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".up.sql") {
			migrations = append(migrations, name)
		}
	}
	sort.Strings(migrations)

	// Create migrations table if it doesn't exist
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run each migration that hasn't been applied yet
	for _, migration := range migrations {
		if err := runMigration(ctx, db, filepath.Join(schemaPath, migration), migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration, err)
		}
	}

	log.Printf("Migrations completed successfully")
	return nil
}

// ensureMigrationsTable creates the migrations tracking table if it doesn't exist
func ensureMigrationsTable(ctx context.Context, db *DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		);
	`
	_, err := db.Pool.Exec(ctx, query)
	return err
}

// runMigration executes a single migration if it hasn't been applied yet
func runMigration(ctx context.Context, db *DB, fullPath, filename string) error {
	// Check if migration has already been applied
	var exists bool
	err := db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name = $1)", filename).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if exists {
		log.Printf("Migration %s already applied, skipping", filename)
		return nil
	}

	// Read the migration file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	// Execute the migration within a transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Run the migration SQL
	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record the migration
	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (name) VALUES ($1)", filename); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	log.Printf("Applied migration: %s", filename)
	return nil
}
