package handlers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// TestDB holds the test database connection pool
var TestDB *pgxpool.Pool

// TestQueries holds the sqlc queries for testing
var TestQueries *models.Queries

// TestJWTClient holds the test JWT client
var TestJWTClient *auth.JWTClient

// TestUserID holds a valid test user ID
var TestUserID string

// TestClerkID holds the test Clerk user ID
var TestClerkID string

// SetupTestDB initializes the test database connection
// This is called once before all tests run
func SetupTestDB(t *testing.T) {
	ctx := context.Background()
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	var err error
	TestDB, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	TestQueries = models.New(TestDB)

	// Initialize JWT client
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret-key"
	}
	TestJWTClient, err = auth.NewJWTClient(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to create JWT client: %v", err)
	}
}

// CleanupTestDB closes the test database connection
func CleanupTestDB(t *testing.T) {
	if TestDB != nil {
		TestDB.Close()
	}
}

// CreateTestUser creates a test user in the database and returns its ID
// This should be called at the beginning of each test that needs a user
func CreateTestUser(t *testing.T, ctx context.Context) string {
	t.Helper()

	// Generate unique identifiers using timestamp and test name
	uniqueID := fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixNano())
	TestClerkID = fmt.Sprintf("test_clerk_%s", uniqueID)
	userEmail := fmt.Sprintf("test_%s@example.com", uniqueID)

	user, err := TestQueries.CreateUser(ctx, models.CreateUserParams{
		ClerkUserID: TestClerkID,
		Email:       userEmail,
		Name:        utils.PgText("Test User"),
		Currency:    utils.PgText("USD"),
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	TestUserID = user.ID
	return TestUserID
}

// CleanupTestUser deletes the test user and all related data
func CleanupTestUser(t *testing.T, ctx context.Context, userID string) {
	t.Helper()

	// For now, just use soft delete which will cascade delete related records
	// In production, ON DELETE CASCADE is set up in the database schema
	TestQueries.DeleteUser(ctx, userID)
}

// GenerateTestToken creates a valid JWT token for the test user
func GenerateTestToken(t *testing.T, clerkUserID string) string {
	t.Helper()

	token, err := TestJWTClient.GenerateToken(clerkUserID)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}
	return token
}

// TruncateTables truncates all test tables (useful for cleanup between tests)
func TruncateTables(t *testing.T, ctx context.Context) {
	t.Helper()

	tables := []string{
		"activity_log",
		"sync_operations",
		"share_access",
		"share_invitations",
		"reflection_questions",
		"reflections",
		"template_questions",
		"reflection_templates",
		"transactions",
		"budget_categories",
		"budgets",
		"payment_methods",
		"categories",
		"users",
	}

	for _, table := range tables {
		_, err := TestDB.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// CreateTestCategory creates a test category for the given user
func CreateTestCategory(t *testing.T, ctx context.Context, userID, name, icon, color string) string {
	t.Helper()

	category, err := TestQueries.CreateCategory(ctx, models.CreateCategoryParams{
		UserID:     utils.PgUUID(userID),
		Name:       name,
		Icon:       utils.PgTextPtr(&icon),
		Color:      utils.PgText(color),
		IsSystem:   utils.PgBool(false),
		DefaultLimit: utils.PgNumericPtr(nil),
	})
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	return category.ID
}

// CreateTestBudget creates a test budget for the given user
func CreateTestBudget(t *testing.T, ctx context.Context, userID, month string, totalLimit float64) string {
	t.Helper()

	// Parse the month string to time.Time
	parsedMonth, err := time.Parse("2006-01-02", month)
	if err != nil {
		t.Fatalf("Failed to parse month: %v", err)
	}

	budget, err := TestQueries.CreateBudget(ctx, models.CreateBudgetParams{
		UserID:     utils.PgUUID(userID),
		Month:      utils.PgDate(parsedMonth),
		TotalLimit: utils.PgNumeric(totalLimit),
		Name:       utils.PgTextPtr(nil),
	})
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}
	return budget.ID
}
