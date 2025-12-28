package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
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
		Currency:    utils.PgText("PHP"),
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

// CreateTestPaymentMethod creates a test payment method for the given user
func CreateTestPaymentMethod(t *testing.T, ctx context.Context, userID, name, methodType string) string {
	t.Helper()

	paymentMethod, err := TestQueries.CreatePaymentMethod(ctx, models.CreatePaymentMethodParams{
		UserID:   utils.PgUUID(userID),
		Name:     name,
		Type:     methodType,
		IsDefault: utils.PgBool(false),
		IsActive:  utils.PgBool(true),
	})
	if err != nil {
		t.Fatalf("Failed to create test payment method: %v", err)
	}
	return paymentMethod.ID
}

// setAuthContext sets up auth context for http.Request in tests
// This also sets up Chi router context for path parameter extraction
func setAuthContext(req *http.Request, userID string) {
	ctx := auth.SetUserIDInContext(req, userID)

	// Set up Chi router context for path parameter extraction
	// This is needed for chi.URLParam to work in tests
	rctx := chi.NewRouteContext()

	// Parse the URL path to extract path parameters
	pathParts := splitPath(req.URL.Path)

	// Common patterns:
	// /api/categories/{id} -> ["api", "categories", "{id}"]
	// /api/budgets/{id} -> ["api", "budgets", "{id}"]
	// /api/analytics/dashboard/{month} -> ["api", "analytics", "dashboard", "{month}"]
	// /api/shares/invitations/{id} -> ["api", "shares", "invitations", "{id}"]
	// /api/payment-methods/{id} -> ["api", "payment-methods", "{id}"]
	// /api/reflections/{month} -> ["api", "reflections", "{month}"]
	// /api/budgets/{budgetId}/categories/{categoryId} -> ["api", "budgets", "{budgetId}", "categories", "{categoryId}"]

	if len(pathParts) >= 1 && pathParts[0] == "api" {
		switch len(pathParts) {
		case 3: // /api/{resource}/{id}
			resource := pathParts[1]
			// Handle special case for reflections and budgets which can use {month} or {id}
			if resource == "reflections" || resource == "budgets" {
				// Check if the third part looks like a date (YYYY-MM) or a UUID
				// Date format: YYYY-MM (7 chars like "2025-01")
				// UUID format: 36 chars with hyphens
				if len(pathParts[2]) <= 7 && strings.Contains(pathParts[2], "-") {
					// This looks like a month parameter (YYYY-MM)
					rctx.URLParams.Add("month", pathParts[2])
				} else {
					// This is a UUID, use id parameter
					rctx.URLParams.Add("id", pathParts[2])
				}
			} else {
				// Most common pattern: /api/categories/{id}, /api/budgets/{id}, etc.
				rctx.URLParams.Add("id", pathParts[2])
			}

		case 4: // /api/{resource}/{sub-resource}/{id-or-value}
			resource := pathParts[1]
			switch resource {
			case "analytics":
				// /api/analytics/dashboard/{month}
				// /api/analytics/spending/{month}
				// /api/analytics/category/{categoryId}
				if pathParts[2] == "category" {
					rctx.URLParams.Add("categoryId", pathParts[3])
				} else {
					rctx.URLParams.Add("month", pathParts[3])
				}
			case "budgets":
				// /api/budgets/{id}/categories
				if pathParts[3] == "categories" {
					rctx.URLParams.Add("id", pathParts[2])
				}
			case "shares", "sharing":
				// /api/shares/invitations/{id}
				// /api/sharing/budgets/{budgetId}
				// /api/sharing/access/{id}
				if pathParts[2] == "budgets" {
					rctx.URLParams.Add("budgetId", pathParts[3])
				} else if pathParts[2] == "access" {
					rctx.URLParams.Add("id", pathParts[3])
				} else {
					rctx.URLParams.Add("id", pathParts[3])
				}
			case "reflections":
				// /api/reflections/month/{month}
				if pathParts[2] == "month" {
					rctx.URLParams.Add("month", pathParts[3])
				}
			case "payment-methods":
				// This should be length 3, but handle if tests use different pattern
				rctx.URLParams.Add("id", pathParts[3])
			default:
				// Handle unexpected patterns
				rctx.URLParams.Add("id", pathParts[3])
			}

		case 5: // /api/{resource}/{id}/{sub-resource}/{sub-id}
			resource := pathParts[1]
			if resource == "budgets" && pathParts[3] == "categories" {
				// /api/budgets/{budgetId}/categories/{categoryId}
				rctx.URLParams.Add("budgetId", pathParts[2])
				rctx.URLParams.Add("categoryId", pathParts[4])
			} else if resource == "budgets" {
				// /api/budgets/{id}/categories - for AddBudgetCategory
				rctx.URLParams.Add("id", pathParts[2])
			} else if resource == "analytics" && pathParts[3] == "category" {
				// /api/analytics/category/{categoryId}
				rctx.URLParams.Add("categoryId", pathParts[4])
			} else {
				// General pattern: /api/{resource}/{id}/{sub-resource}/{sub-id}
				rctx.URLParams.Add("id", pathParts[2])
				rctx.URLParams.Add(pathParts[3]+"Id", pathParts[4])
			}

		case 6: // /api/budgets/{id}/categories/{categoryId}/...
			// /api/budgets/{budgetId}/categories/{categoryId}
			rctx.URLParams.Add("budgetId", pathParts[2])
			rctx.URLParams.Add("categoryId", pathParts[4])
		}
	}

	*req = *req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
}

// splitPath splits URL path into parts
func splitPath(path string) []string {
	if path == "" || path == "/" {
		return []string{}
	}

	// Remove leading slash and split
	path = path[1:]
	if path == "" {
		return []string{}
	}

	parts := make([]string, 0)
	for _, part := range strings.Split(path, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}
